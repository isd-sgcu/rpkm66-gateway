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

	vaccineClient "github.com/isd-sgcu/rpkm66-gateway/src/app/client/vaccine"
	authHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/auth"
	baanHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/baan"
	ciHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/checkin"
	eHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/estamp"
	fileHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/file"
	grpHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/group"
	health_check "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/health-check"
	usrHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/user"
	vaccineHdr "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/vaccine"
	guard "github.com/isd-sgcu/rpkm66-gateway/src/app/middleware/auth"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/router"
	authSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/auth"
	baanSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/baan"
	ciSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/checkin"
	eSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/estamp"
	fileSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/file"
	grpSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/group"
	usrSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/user"
	vaccineSrv "github.com/isd-sgcu/rpkm66-gateway/src/app/service/vaccine"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/validator"
	"github.com/isd-sgcu/rpkm66-gateway/src/config"
	"github.com/isd-sgcu/rpkm66-gateway/src/constant/auth"
	_ "github.com/isd-sgcu/rpkm66-gateway/src/docs"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
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

// @tag.name vaccine
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
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "config").
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

	vacClient := vaccineClient.NewClient(conf.Vaccine)
	vacSrv := vaccineSrv.NewService(userSrv, vacClient)
	vacHdr := vaccineHdr.NewHandler(vacSrv, v)
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
	authGuard := guard.NewAuthGuard(athSrv, auth.ExcludePath, conf.App)

	r := router.NewFiberRouter(&authGuard, conf.App)

	r.GetHealthCheck("/", hc.HealthCheck)

	r.PutUser("/", userHdr.CreateOrUpdate)
	r.PatchUser("/", userHdr.Update)

	if conf.App.Debug {
		r.GetUser("/:id", userHdr.FindOne)
		r.PostUser("/", userHdr.Create)
		r.DeleteUser("/:id", userHdr.Delete)
	}

	r.GetAuth("/me", athHdr.Validate)
	r.PostAuth("/verify", athHdr.VerifyTicket)
	r.PostAuth("/refreshToken", athHdr.RefreshToken)

	r.PostFile("/upload", fleHdr.Upload)

	r.PostVaccine("/verify", vacHdr.Verify)

	r.GetBaan("/", bnHdr.FindAll)
	r.GetBaan("/:id", bnHdr.FindOne)

	r.GetGroup("/", gHdr.FindOne)
	r.GetGroup("/:token", gHdr.FindByToken)
	r.PostGroup("/:token", gHdr.Join)
	r.DeleteGroup("/leave", gHdr.Leave)
	r.PutGroup("/select", gHdr.SelectBaan)
	r.DeleteGroup("/members/:member_id", gHdr.DeleteMember)

	r.PostQr("/checkin/verify", ciHandler.CheckinVerify)
	r.PostQr("/checkin/confirm", ciHandler.CheckinConfirm)
	r.PostQr("/estamp/verify", estampHdr.VerifyEstamp)
	r.PostQr("/estamp/confirm", userHdr.ConfirmEstamp)
	r.PostQr("/ticket", userHdr.VerifyTicket)

	r.GetEstamp("/:id", estampHdr.FindEventByID)
	r.GetEstamp("/", estampHdr.FindAllEventWithType)
	r.GetEstamp("/user", userHdr.GetUserEstamp)

	go func() {
		if err := r.Listen(fmt.Sprintf(":%v", conf.App.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Str("service", "rnkm-gateway").
				Msg("Server not close properly")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"server": func(ctx context.Context) error {
			return r.Shutdown()
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
