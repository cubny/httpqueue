package api

import (
	"context"
	"errors"
	"github.com/cubny/cart/internal/app"
	"github.com/cubny/cart/internal/infra/http/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	api500Count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http-request-scheduler",
			Name:      "error_500_counter",
			Help:      "Counter of 500 responses of http-request-scheduler api",
		}, []string{"method", "reason"})
)

func init() {
	prometheus.MustRegister(api500Count)
}

// ServiceProvider holds all the business logic
type ServiceProvider interface {
	CreateTimer(ctx context.Context, timer *app.Timer) (*app.Timer, error)
	GetTimer(ctx context.Context) error
}

// Router handles http requests
type Router struct {
	service ServiceProvider
	http.Handler
}

// New creates a new handler to handle http requests
func New(service ServiceProvider) (*Router, error) {
	if service == nil {
		return nil, errors.New("service is not set up")
	}

	h := &Router{
		service: service,
	}
	router := httprouter.New()

	chain := middleware.NewChain(middleware.ContentTypeJSON)

	router.GET("/health", h.health)
	router.POST("/timers", chain.Wrap(h.setTimer))
	router.GET("/timers", chain.Wrap(h.getTimer))

	h.Handler = router
	return h, nil
}

func (h *Router) health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Error("failed to compose body of the response")
	}
}
