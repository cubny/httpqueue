package api_test

import (
	"encoding/json"
	"github.com/cubny/httpqueue/internal/app/timer"
	"github.com/cubny/httpqueue/internal/infra/http/api"
	mocks "github.com/cubny/httpqueue/internal/mocks/app/timer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type spec struct {
	Name           string
	ReqBody        string
	ExpectedStatus int
	ExpectedBody   string
	Method         string
	Target         string
	MockFn         func(s *mocks.Service)
}

func (s *spec) execHTTPTestCases(sp *mocks.Service) func(t *testing.T) {
	return func(t *testing.T) {
		s.MockFn(sp)
		handler, err := api.New(sp)
		assert.Nil(t, err)
		s.HandlerTest(t, handler)
	}
}

// HandlerTest is a helper method to run http test cases
func (s *spec) HandlerTest(t *testing.T, h *api.Router) {
	t.Helper()

	req := httptest.NewRequest(s.Method, s.Target, strings.NewReader(s.ReqBody))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	switch {
	case s.ExpectedBody != "" && isJSON(s.ExpectedBody):
		assert.JSONEq(t, s.ExpectedBody, string(body))
	case s.ExpectedBody != "" && !isJSON(s.ExpectedBody):
		assert.Equal(t, s.ExpectedBody, strings.TrimSpace(string(body)))
	}

	assert.Equal(t, s.ExpectedStatus, resp.StatusCode)
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func TestRouter_setTimer(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewService(ctrl)

	specs := []spec{
		{
			Name:   "ok",
			Method: http.MethodPost,
			Target: "/timers",
			MockFn: func(s *mocks.Service) {
				s.EXPECT().CreateTimer(gomock.Any(), gomock.Any()).Return(&timer.Timer{
					ID: "1",
				}, nil)
			},
			ReqBody:        `{"url":"http://valid.url","hours":0,"minutes":0,"seconds":0}`,
			ExpectedBody:   `{"id":"1"}`,
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:   "service returns error",
			Method: http.MethodPost,
			Target: "/timers",
			MockFn: func(s *mocks.Service) {
				s.EXPECT().CreateTimer(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			ReqBody:        `{"url":"http://valid.url","hours":0,"minutes":0,"seconds":0}`,
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - failed to set timers due to server internal error"}}`,
			ExpectedStatus: http.StatusInternalServerError,
		},
		{
			Name:           "invalid url",
			Method:         http.MethodPost,
			MockFn:         func(s *mocks.Service) {},
			Target:         "/timers",
			ReqBody:        `{"url":"invalid.url","hours":0,"minutes":0,"seconds":0}`,
			ExpectedBody:   `{"error":{"code":422, "details":"Invalid params - invalid param: invalid 'POST' field 'url'"}}`,
			ExpectedStatus: http.StatusUnprocessableEntity,
		},
		{
			Name:           "empty body",
			Method:         http.MethodPost,
			MockFn:         func(s *mocks.Service) {},
			Target:         "/timers",
			ReqBody:        ``,
			ExpectedBody:   `{"error":{"code":400, "details":"Bad Request - cannot set timers, bad request payload"}}`,
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(service))
	}
}

func TestRouter_getTimer(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewService(ctrl)
	now := time.Now()

	specs := []spec{
		{
			Name:   "ok",
			Method: http.MethodGet,
			Target: "/timers/1",
			MockFn: func(s *mocks.Service) {
				s.EXPECT().GetTimer(gomock.Any(), "1").Return(&timer.Timer{
					ID:     "1",
					FireAt: now.Add(2 * time.Second),
				}, nil)
			},
			ExpectedBody:   `{"ID":"1", "time_left":1}`,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:   "not found",
			Method: http.MethodGet,
			Target: "/timers/1",
			MockFn: func(s *mocks.Service) {
				s.EXPECT().GetTimer(gomock.Any(), "1").Return(nil, timer.ErrTimerNotFound)
			},
			ExpectedBody:   `{"error":{"code":404, "details":"Not found - timer does not exist"}}`,
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name:   "timer is archived",
			Method: http.MethodGet,
			Target: "/timers/1",
			MockFn: func(s *mocks.Service) {
				s.EXPECT().GetTimer(gomock.Any(), "1").Return(nil, timer.ErrTimerArchived)
			},
			ExpectedBody:   `{"ID":"1", "time_left":0}`,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:   "service unknown error",
			Method: http.MethodGet,
			Target: "/timers/1",
			MockFn: func(s *mocks.Service) {
				s.EXPECT().GetTimer(gomock.Any(), "1").Return(nil, assert.AnError)
			},
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - failed to get timers due to server internal error"}}`,
			ExpectedStatus: http.StatusInternalServerError,
		},
	}
	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(service))
	}
}
