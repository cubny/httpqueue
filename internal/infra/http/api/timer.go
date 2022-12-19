package api

import (
	"encoding/json"
	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// setTimer is the handler for
// swagger:route POST /timers setTimersRequest
//
// Schedule a new timer.
//
// Responses:
//
//	201: setTimers
//	400: invalidRequestBody
//	404: notFoundError
//	422: invalidParams
//	500: serverError
func (h *Router) setTimer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	command, err := toSetTimerCommand(w, r)
	if err != nil {
		// toSetTimerCommand responds with a proper error
		return
	}

	t, err := h.service.CreateTimer(r.Context(), command)
	if err != nil {
		log.WithError(err).Errorf("setTimers: service %s", err)
		api500Count.With(prometheus.Labels{"method": "setTimers", "reason": "service"}).Inc()
		_ = InternalError(w, "failed to set timers due to server internal error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toSetTimerResponse(t)); err != nil {
		log.WithError(err).Errorf("setTimer: encoder %s", err)
		api500Count.With(prometheus.Labels{"method": "setTimer", "reason": "encoder"}).Inc()
		_ = InternalError(w, "cannot encode response")
		return
	}
}

// getTimer is the handler for
// swagger:route GET /timers/{timer_id} getTimerRequest
//
// Responds how much time remains until the timer's webhook is shot.
//
// Responses:
//
//	200: getTimer
//	404: notFoundError
//	422: invalidParams
//	500: serverError
func (h *Router) getTimer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	timerID := p.ByName("id")
	if timerID == "" {
		_ = InvalidParams(w, "invalid param: id is empty")
		return
	}

	var resp GetTimerResponse

	t, err := h.service.GetTimer(r.Context(), timerID)
	switch {
	case err == timer.ErrTimerNotFound:
		_ = NotFound(w, "timer does not exist")
		return
	case err == timer.ErrTimerArchived:
		resp = getArchivedTimerResponse(timerID)
	case err != nil:
		log.WithError(err).Errorf("getTimers: service %s", err)
		api500Count.With(prometheus.Labels{"method": "getTimers", "reason": "service"}).Inc()
		_ = InternalError(w, "failed to get timers due to server internal error")
		return
	default:
		resp = toGetTimersResponse(t)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Errorf("getTimer: encoder %s", err)
		api500Count.With(prometheus.Labels{"method": "getTimer", "reason": "encoder"}).Inc()
		_ = InternalError(w, "cannot encode response")
		return
	}
}
