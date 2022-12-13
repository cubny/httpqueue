package main

import (
	"context"
	"flag"
	"github.com/cubny/cart/internal/infra/http/api"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cubny/cart/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		optsAddr    = flag.String("addr", ":8080", "HTTP bind address")
		metricsAddr = flag.String("metricsAddr", ":8081", "Metrics HTTP bind address")
		dataPath    = flag.String("data", "/app/data/cart.db", "Path to the sqlite3 data file")
	)
	flag.Parse()

	service, err := service.New()
	if err != nil {
		log.Fatalf("cannot create service, %s", err)
	}

	handler, err := api.New(service)
	if err != nil {
		log.Fatalf("cannot create handler, %s", err)
	}

	srv := http.Server{
		Addr:    *optsAddr,
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Printf("shuting down the http server...")
		idleCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.Shutdown(idleCtx); err != nil {
			panic(err)
		}
		close(idleConnsClosed)
	}()

	go func() {
		log.Debugf("starting metrics server %s", *metricsAddr)
		log.Fatal(http.ListenAndServe(*metricsAddr, promhttp.Handler()))
	}()

	log.Printf("HTTP Server starting %s", *optsAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
	<-idleConnsClosed
}
