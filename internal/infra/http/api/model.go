package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cubny/httpqueue/internal/app/timer"
)

// SetTimersRequest is the request model to set a new timer
//
// swagger:model setTimersRequest
type SetTimersRequest struct {
	Hours   int    `json:"hours"`
	Minutes int    `json:"minutes"`
	Seconds int    `json:"seconds"`
	URL     string `json:"url"`
}

// SetTimerResponse is the response model to set a new timer
//
// swagger:model setTimersResponse
type SetTimerResponse struct {
	ID string `json:"id"`
}

func (r *SetTimersRequest) Validate() error {
	if _, err := url.ParseRequestURI(r.URL); err != nil {
		return errors.New("invalid 'POST' field 'url'")
	}

	return nil
}

func toSetTimerCommand(w http.ResponseWriter, req *http.Request) (timer.SetTimerCommand, error) {
	request := &SetTimersRequest{}
	if err := json.NewDecoder(req.Body).Decode(request); err != nil {
		_ = BadRequest(w, "cannot set timers, bad request payload")
		return timer.SetTimerCommand{}, err
	}

	if err := request.Validate(); err != nil {
		_ = InvalidParams(w, fmt.Sprintf("invalid param: %v", err))
		return timer.SetTimerCommand{}, err
	}

	return timer.SetTimerCommand{
		Hours:   request.Hours,
		Minutes: request.Minutes,
		Seconds: request.Seconds,
		URLRaw:  request.URL,
	}, nil
}

func toSetTimerResponse(t *timer.Timer) SetTimerResponse {
	return SetTimerResponse{ID: t.ID}
}

// GetTimerResponse is the response model to get a timer
//
// swagger:model GetTimerResponse
type GetTimerResponse struct {
	ID              string `json:"ID"`
	TimeLeftSeconds int    `json:"time_left"`
}

func getArchivedTimerResponse(timerID string) GetTimerResponse {
	return GetTimerResponse{
		ID: timerID,
	}
}

func toGetTimersResponse(t *timer.Timer) GetTimerResponse {
	timeLeft := t.DelayFromNowSeconds()

	return GetTimerResponse{
		ID:              t.ID,
		TimeLeftSeconds: int(timeLeft),
	}
}
