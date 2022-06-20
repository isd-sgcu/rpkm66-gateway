package main

import (
	"context"
	"fmt"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/health-check"
	usrHdr "github.com/isd-sgcu/rnkm65-gateway/src/app/handler/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/router"
	usrSrv "github.com/isd-sgcu/rnkm65-gateway/src/app/service/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/config"
	_ "github.com/isd-sgcu/rnkm65-gateway/src/docs"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/isd-sgcu/rnkm65-gateway/src/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

// @tag.name user
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
		log.Fatal("Cannot load config", err.Error())
	}

	v, err := validator.NewValidator()
	if err != nil {
		log.Fatal(err.Error())
	}

	backendConn, err := grpc.Dial(conf.Service.Backend, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot connect to rnkm65 backend service: ", err.Error())
	}

	r := router.NewFiberRouter()

	hc := health_check.NewHandler()

	uClient := proto.NewUserServiceClient(backendConn)
	uSrv := usrSrv.NewService(uClient)
	uHdr := usrHdr.NewHandler(uSrv, v)

	r.GetHealthCheck("/", hc.HealthCheck)

	r.GetUser("/", uHdr.FindOne)
	r.PostUser("/", uHdr.Create)
	r.PutUser("/:id", uHdr.Update)
	r.PutUser("/", uHdr.CreateOrUpdate)
	r.DeleteUser("/:id", uHdr.Delete)

	go func() {
		if err := r.Listen(fmt.Sprintf(":%v", conf.App.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
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

		log.Printf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
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

				log.Printf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Printf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}
