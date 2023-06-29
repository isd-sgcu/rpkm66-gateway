package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/cfgldr"
	_ "github.com/isd-sgcu/rpkm66-gateway/docs"
	authHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/auth"
	baanHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/baan"
	ciHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/checkin"
	eHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/estamp"
	fileHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/file"
	grpHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/group"
	health_check "github.com/isd-sgcu/rpkm66-gateway/internal/handler/health-check"
	usrHdr "github.com/isd-sgcu/rpkm66-gateway/internal/handler/user"
	guard "github.com/isd-sgcu/rpkm66-gateway/internal/middleware/auth"
	"github.com/isd-sgcu/rpkm66-gateway/internal/router"
	authSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/auth"
	baanSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/baan"
	ciSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/checkin"
	eSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/estamp"
	fileSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/file"
	grpSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/group"
	usrSrv "github.com/isd-sgcu/rpkm66-gateway/internal/service/user"
	"github.com/isd-sgcu/rpkm66-gateway/internal/validator"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title RNKM Backend
// @version 1.0
// @description.markdown

// @contact.name ISD Team
// @contact.email sd.team.sgcu@gmail.com

// @schemes https http

// @securityDefinitions.apikey  AuthToken
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

// @tag.name Chula SSO
// @tag.description Chula SSO documentation
// @tag.docs.url https://account.it.chula.ac.th/wiki/doku.php?id=how_does_it_work
// @tag.docs.description Chula SSO documentation

// @tag.name health check
// @tag.description.markdown

// @tag.name auth
// @tag.description.markdown

// @tag.name user
// @tag.description.markdown

// @tag.name file
// @tag.description.markdown

// @tag.name group
// @tag.description.markdown

// @tag.name baan
// @tag.description.markdown

// @tag.name event
// @tag.description.markdown

func main() {
	conf, err := cfgldr.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "cfgldr").
			Msg("Failed to start service")
	}

	v, err := validator.NewValidator()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "validator").
			Msg("Failed to start service")
	}

	backendConn, err := grpc.Dial(conf.Service.Backend, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "rnkm-backend").
			Msg("Cannot connect to service")
	}

	authConn, err := grpc.Dial(conf.Service.Auth, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "rnkm-auth").
			Msg("Cannot connect to service")
	}

	fileConn, err := grpc.Dial(conf.Service.File, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "rnkm-file").
			Msg("Cannot connect to service")
	}

	hc := health_check.NewHandler()

	usrClient := proto.NewUserServiceClient(backendConn)
	userSrv := usrSrv.NewService(usrClient)
	userHdr := usrHdr.NewHandler(userSrv, v)

	authClient := proto.NewAuthServiceClient(authConn)
	athSrv := authSrv.NewService(authClient)
	athHdr := authHdr.NewHandler(athSrv, userSrv, v)

	fileClient := proto.NewFileServiceClient(fileConn)
	fleSrv := fileSrv.NewService(fileClient)
	fleHdr := fileHdr.NewHandler(fleSrv, userSrv, conf.App.MaxFileSize)

	gClient := proto.NewGroupServiceClient(backendConn)
	gSrv := grpSrv.NewService(gClient)
	gHdr := grpHdr.NewHandler(gSrv, v)

	bnClient := proto.NewBaanServiceClient(backendConn)
	bnSrv := baanSrv.NewService(bnClient)
	bnHdr := baanHdr.NewHandler(bnSrv, userSrv)

	checkinClient := proto.NewCheckinServiceClient(backendConn)
	checkinSrv := ciSrv.NewService(checkinClient)

	estampClient := proto.NewEventServiceClient(backendConn)
	estampSrv := eSrv.NewService(estampClient)
	estampHdr := eHdr.NewHandler(estampSrv, v)

	ciHandler := ciHdr.NewHandler(checkinSrv, v)
	authGuard := guard.NewAuthGuard(athSrv)

	r := router.NewGinRouter(&authGuard, conf.App)

	r.SetHandler("GET /", hc.HealthCheck)
	r.SetHandler("PUT /user", userHdr.CreateOrUpdate)
	r.SetHandler("PATCH /user", userHdr.Update)
	r.SetHandler("GET /auth/me", athHdr.Validate)
	r.SetHandler("POST /auth/verify", athHdr.VerifyTicket)
	r.SetHandler("POST /auth/refreshToken", athHdr.RefreshToken)
	r.SetHandler("PUT /file/upload", fleHdr.Upload)
	r.SetHandler("GET /baan", bnHdr.FindAll)
	r.SetHandler("GET /baan/:id", bnHdr.FindOne)
	r.SetHandler("GET /group", gHdr.FindOne)
	r.SetHandler("GET /group/:token", gHdr.FindByToken)
	r.SetHandler("POST /group/:token", gHdr.Join)
	r.SetHandler("DELETE /group/leave", gHdr.Leave)
	r.SetHandler("PUT /group/select", gHdr.SelectBaan)
	r.SetHandler("DELETE /members/:member_id", gHdr.DeleteMember)
	r.SetHandler("POST /qr/checkin/verify", ciHandler.CheckinVerify)
	r.SetHandler("POST /qr/checkin/confirm", ciHandler.CheckinConfirm)
	r.SetHandler("POST /qr/estamp/verify", estampHdr.VerifyEstamp)
	r.SetHandler("POST /qr/estamp/confirm", userHdr.ConfirmEstamp)
	r.SetHandler("POST /qr/ticket", userHdr.VerifyTicket)
	r.SetHandler("GET /estamp/:id", estampHdr.FindEventByID)
	r.SetHandler("GET /estamp", estampHdr.FindAllEventWithType)
	r.SetHandler("GET /estamp/user", userHdr.GetUserEstamp)
	r.SetHandler("GET /user/:id", userHdr.FindOne)
	r.SetHandler("POST /user", userHdr.Create)
	r.SetHandler("DELETE /user/:id", userHdr.Delete)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", conf.App.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Str("service", "rnkm-gateway").
				Msg("Server not close properly")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"server": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	<-wait
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-s

		log.Info().
			Str("service", "graceful shutdown").
			Msgf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().
				Str("service", "graceful shutdown").
				Msgf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().
						Str("service", "graceful shutdown").
						Err(err).
						Msgf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
