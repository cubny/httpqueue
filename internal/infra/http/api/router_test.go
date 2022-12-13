//go:generate mockgen -package handler -source handler.go -destination handler_mocks_test.go
package api_test

import (
	"github.com/cubny/cart/internal/infra/http/api"
	"testing"

	"github.com/cubny/cart/internal/tests"

	"github.com/stretchr/testify/assert"
)

func execHTTPTestCases(t *testing.T, sp api.ServiceProvider, tcs []tests.TestCase) {
	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			handler, err := api.New(sp)
			assert.Nil(t, err)
			tests.HandlerTest(t, handler, &tc)
		})
	}
}
