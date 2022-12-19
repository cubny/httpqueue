// Package api httpqueue
//
// Documentation of the httpqueue service.
// It is a service to schedule webhooks.
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//	Host: httpqueue
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package api

import (
	"errors"
	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/infra/http/api/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	api500Count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "httpqueue",
			Name:      "error_500_counter",
			Help:      "Counter of 500 responses of httpqueue api",
		}, []string{"method", "reason"})
)

func init() {
	prometheus.MustRegister(api500Count)
}

// Router handles http requests
type Router struct {
	service timer.Service
	http.Handler
}

// New creates a new handler to handle http requests
func New(service timer.Service) (*Router, error) {
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
	router.GET("/timers/:id", chain.Wrap(h.getTimer))

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
