package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// setTimers is the handler for
// POST /timers/
func (h *Router) setTimers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := h.service.CreateTimer(r.Context(), 1)
	switch {
	case err == cart.ErrInvalidUserID:
		_ = InvalidParams(w, "user is invalid")
		return
	case err != nil:
		log.WithError(err).Errorf("createCart: service %s", err)
		api500Count.With(prometheus.Labels{"method": "createCart", "reason": "service"}).Inc()
		_ = InternalError(w, "cannot create cart")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		log.WithError(err).Errorf("createCart: encoder %s", err)
		api500Count.With(prometheus.Labels{"method": "createCart", "reason": "encoder"}).Inc()
		_ = InternalError(w, "cannot encode response")
		return
	}
}
