package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/config"
	asynqTimer "github.com/cubny/httpqueue/internal/infra/asynq/timer"
	"github.com/cubny/httpqueue/internal/infra/http/api"
	internalHttpClient "github.com/cubny/httpqueue/internal/infra/http/client/timer"
	"github.com/cubny/httpqueue/internal/infra/redis"
	repo "github.com/cubny/httpqueue/internal/infra/redis/timer"
)

type App struct {
	ctx context.Context
	cfg *config.Config

	apiServer   *http.Server
	consumer    *asynqTimer.Consumer
	relay       *asynqTimer.Relay
	service     timer.Service
	redisClient *redis.Client
	db          *repo.DB

	err error
}

type AppMode string

const (
	AppModeAPI     AppMode = "api"
	AppModeRelay   AppMode = "relay"
	AppModeWorkers AppMode = "workers"
	AppModeAll     AppMode = "all"
)

func Init(ctx context.Context) (*App, error) {
	a := &App{ctx: ctx}
	a.initConfig()
	a.initRepo()
	a.initService()
	a.initPromHandler()

	switch AppMode(a.cfg.AppMode) {
	case AppModeAll:
		a.initRelay()
		a.initConsumer()
		a.initAPIServer()
	case AppModeWorkers:
		a.initConsumer()
	case AppModeRelay:
		a.initRelay()
	case AppModeAPI:
		a.initAPIServer()
	}

	return a, a.err
}

func (a *App) ifNoError(fn func() *App) *App {
	if a.err != nil {
		return a
	}
	return fn()
}

func (a *App) initRepo() *App {
	return a.ifNoError(func() *App {
		a.redisClient = redis.NewRedis(&a.cfg.Redis)
		a.db = repo.NewDB(a.redisClient, &a.cfg.DB)

		return a
	})
}

func (a *App) initService() *App {
	return a.ifNoError(func() *App {
		service, err := timer.NewService(a.db)
		if err != nil {
			a.err = err
			return a
		}

		a.service = service
		return a
	})
}

func (a *App) initConfig() *App {
	return a.ifNoError(
		func() *App {
			cfg, err := config.New(a.ctx)
			if err != nil {
				log.Fatalf("failed to initiate config: %v", err)
			}

			a.cfg = cfg
			return a
		},
	)
}

func (a *App) initRelay() *App {
	return a.ifNoError(
		func() *App {
			aClient := asynq.NewClient(a.redisClient)
			producer := asynqTimer.NewProducer(aClient, &a.cfg.Producer)
			relay, err := asynqTimer.NewRelay(&a.cfg.Relay, a.db, producer)
			if err != nil {
				a.err = fmt.Errorf("faild to initiate the relay, %v", err)
				return a
			}

			a.relay = relay
			go a.relay.Start(a.ctx)

			return a
		},
	)
}

func (a *App) initPromHandler() *App {
	return a.ifNoError(func() *App {
		log.Debugf("starting metrics server %d", a.cfg.HTTP.MetricsPort)
		go func() {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.HTTP.MetricsPort), promhttp.Handler()); err != nil {
				a.err = fmt.Errorf("failed to start the prometheus handler, %v", err)
			}
		}()

		return a
	})
}

func (a *App) stopAPIServer() *App {
	log.Info("shutting down HTTP component")
	tctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := a.apiServer.Shutdown(tctx); err != nil {
		a.err = fmt.Errorf("failed to shut down api server, %v", err)
		return a
	}
	log.Infof("api server shut down successfully")
	return a
}

func (a *App) initAPIServer() *App {
	return a.ifNoError(func() *App {
		handler, err := api.New(a.service)
		if err != nil {
			a.err = fmt.Errorf("cannot create handler, %v", err)
			return a
		}
		a.apiServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", a.cfg.HTTP.Port),
			Handler: handler,
		}

		go func() {
			log.Infof("starting API server %d", a.cfg.HTTP.Port)
			if err = a.apiServer.ListenAndServe(); err != nil {
				a.err = err
			}
		}()

		return a
	})
}

func (a *App) initConsumer() *App {
	return a.ifNoError(func() *App {
		srv := asynq.NewServer(
			a.redisClient,
			asynq.Config{
				// number of concurrent workers
				Concurrency: a.cfg.ConsumerConcurrency,
			},
		)

		httpClient := internalHttpClient.NewClient()

		processor, err := asynqTimer.NewProcessor(a.service, httpClient)
		if err != nil {
			log.Fatalf("failed to initiate the timer task processor")
		}

		consumer, err := asynqTimer.NewConsumer(srv, processor)
		if err != nil {
			log.Fatalf("failed to initiate the timer task consumer")
		}

		a.consumer = consumer

		go func() {
			if err := a.consumer.Run(); err != nil {
				a.err = fmt.Errorf("failed to start the consumer")
			}
		}()

		return a
	})
}

func (a *App) Stop() error {
	a.err = a.ctx.Err()
	if a.relay != nil {
		a.relay.Stop()
	}

	if a.apiServer != nil {
		a.stopAPIServer()
	}

	return a.redisClient.Close()
}

func WaitTermination() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
	<-sigs
}
