package main

import (
	"context"
	"fmt"
	authHdr "github.com/isd-sgcu/rnkm65-gateway/src/app/handler/auth"
	fileHdr "github.com/isd-sgcu/rnkm65-gateway/src/app/handler/file"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/health-check"
	usrHdr "github.com/isd-sgcu/rnkm65-gateway/src/app/handler/user"
	guard "github.com/isd-sgcu/rnkm65-gateway/src/app/middleware/auth"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/router"
	authSrv "github.com/isd-sgcu/rnkm65-gateway/src/app/service/auth"
	fileSrv "github.com/isd-sgcu/rnkm65-gateway/src/app/service/file"
	usrSrv "github.com/isd-sgcu/rnkm65-gateway/src/app/service/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	"github.com/isd-sgcu/rnkm65-gateway/src/config"
	"github.com/isd-sgcu/rnkm65-gateway/src/constant/auth"
	_ "github.com/isd-sgcu/rnkm65-gateway/src/docs"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	uClient := proto.NewUserServiceClient(backendConn)
	uSrv := usrSrv.NewService(uClient)
	uHdr := usrHdr.NewHandler(uSrv, v)

	aClient := proto.NewAuthServiceClient(authConn)
	aSrv := authSrv.NewService(aClient)
	aHdr := authHdr.NewHandler(aSrv, uSrv, v)

	fClient := proto.NewFileServiceClient(fileConn)
	fSrv := fileSrv.NewService(fClient)
	fHdr := fileHdr.NewHandler(fSrv, uSrv, conf.App.MaxFileSize)

	authGuard := guard.NewAuthGuard(aSrv, auth.ExcludePath, conf.Guard.Phase)

	r := router.NewFiberRouter(&authGuard, conf.App)

	r.GetHealthCheck("/", hc.HealthCheck)

	r.GetUser("/:id", uHdr.FindOne)
	r.PostUser("/", uHdr.Create)
	r.PutUser("/:id", uHdr.Update)
	r.PutUser("/", uHdr.CreateOrUpdate)
	r.DeleteUser("/:id", uHdr.Delete)

	r.GetAuth("/me", aHdr.Validate)
	r.PostAuth("/verify", aHdr.VerifyTicket)
	r.PostAuth("/refreshToken", aHdr.RefreshToken)

	r.PostFile("/upload", fHdr.Upload)

	r.PostMethod("/vaccine/callback", uHdr.Verify)

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
