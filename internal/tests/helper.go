package tests

import (
	"encoding/json"
	"github.com/cubny/cart/internal/infra/http/api"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase is meant to be used in test tables for testing http handlers
type TestCase struct {
	// name of the test
	Name string
	// reqBody is the request body, usually in json
	ReqBody string
	// expectedStatus is the expected return status code
	ExpectedStatus int
	// expectedBody is the expected return body, usually in json
	ExpectedBody string
	// method is the http method the http test server needs to be called with
	Method string
	// accessKey is the auth header value to send the reqeust with
	AccessKey string
	// target is the route of the handler to be tested
	Target string
}

// HandlerTest is a helper method to run http test cases
func HandlerTest(t *testing.T, h *api.Router, tc *TestCase) {
	t.Helper()

	req := httptest.NewRequest(tc.Method, tc.Target, strings.NewReader(tc.ReqBody))

	auth.AddKeyToRequest(req, tc.AccessKey)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	switch {
	case tc.ExpectedBody != "" && isJSON(tc.ExpectedBody):
		assert.JSONEq(t, tc.ExpectedBody, string(body))
	case tc.ExpectedBody != "" && !isJSON(tc.ExpectedBody):
		assert.Equal(t, tc.ExpectedBody, strings.TrimSpace(string(body)))
	}

	assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
