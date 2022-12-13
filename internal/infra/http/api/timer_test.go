package api_test

import (
	"github.com/cubny/cart/internal/service"
	"github.com/cubny/cart/internal/tests"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandler_CreateCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := NewMockAuthProvider(ctrl)
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "abc123456").
		Return(&auth.AccessKey{UserID: 1, Key: "abc123456"}, nil).AnyTimes()
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "unauthorised").
		Return(nil, auth.ErrNotFound).AnyTimes()

	serviceMock := NewMockServiceProvider(ctrl)
	serviceMock.EXPECT().CreateCart(gomock.Any(), int64(1)).Return(&cart.Cart{
		ID:     1,
		UserID: 1,
	}, nil)

	tests := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         "/carts",
			AccessKey:      "abc123456",
			ExpectedBody:   `{"id":1, "user_id":1}`,
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "unauthorised - error",
			Method:         http.MethodPost,
			Target:         "/carts",
			AccessKey:      "unauthorised",
			ExpectedBody:   `{"error":{"code":100401, "details": "Unauthorised access - incorrect access_key"}}`,
			ExpectedStatus: http.StatusUnauthorized,
		},
	}

	execHTTPTestCases(t, serviceMock, authMock, tests)
}

func TestHandler_AddItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := NewMockAuthProvider(ctrl)
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "abc123456").
		Return(&auth.AccessKey{UserID: 1, Key: "abc123456"}, nil).AnyTimes()
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "unauthorised").
		Return(nil, auth.ErrNotFound).AnyTimes()

	firstItem := &cart.Item{
		CartID:    int64(1),
		ProductID: int64(1),
		Quantity:  int64(1),
		Price:     cart.Price(100.00),
	}
	secondItem := &cart.Item{
		CartID:    int64(2),
		ProductID: int64(1),
		Quantity:  int64(1),
		Price:     cart.Price(100.00),
	}
	serviceMock := NewMockServiceProvider(ctrl)
	serviceMock.EXPECT().
		AddItem(gomock.Any(), int64(1), firstItem).
		Return(nil)
	serviceMock.EXPECT().
		AddItem(gomock.Any(), int64(1), secondItem).
		Return(service.ErrTimerNotFound)

	testsCases := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         "/carts/1/items",
			AccessKey:      "abc123456",
			ReqBody:        `{"product_id":1, "quantity":1, "price": 100.00}`,
			ExpectedBody:   `{"cart_id":1, "id":0, "price":100, "product_id":1, "quantity":1}`,
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "cart does not exist - error",
			Method:         http.MethodPost,
			Target:         "/carts/2/items",
			AccessKey:      "abc123456",
			ReqBody:        `{"product_id":1, "quantity":1, "price": 100.00}`,
			ExpectedBody:   `{"error":{"code":100404, "details":"Not found - cart does not exist"}}`,
			ExpectedStatus: http.StatusNotFound,
		},
	}

	execHTTPTestCases(t, serviceMock, authMock, testsCases)
}

func TestHandler_RemoveItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := NewMockAuthProvider(ctrl)
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "abc123456").
		Return(&auth.AccessKey{UserID: 1, Key: "abc123456"}, nil).AnyTimes()
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "unauthorised").
		Return(nil, auth.ErrNotFound).AnyTimes()

	serviceMock := NewMockServiceProvider(ctrl)
	serviceMock.EXPECT().RemoveItem(gomock.Any(), int64(1), int64(1)).Return(nil)
	serviceMock.EXPECT().RemoveItem(gomock.Any(), int64(1), int64(2)).Return(service.ErrItemNotFound)

	testsCases := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodDelete,
			Target:         "/items/1",
			AccessKey:      "abc123456",
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Name:           "item not found",
			Method:         http.MethodDelete,
			Target:         "/items/2",
			AccessKey:      "abc123456",
			ExpectedBody:   `{"error":{"code":100404, "details":"Not found - item does not exist"}}`,
			ExpectedStatus: http.StatusNotFound,
		},
	}

	execHTTPTestCases(t, serviceMock, authMock, testsCases)
}

func TestHandler_EmptyCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := NewMockAuthProvider(ctrl)
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "abc123456").
		Return(&auth.AccessKey{UserID: 1, Key: "abc123456"}, nil).AnyTimes()
	authMock.EXPECT().
		VerifyAccessKey(gomock.Any(), "unauthorised").
		Return(nil, auth.ErrNotFound).AnyTimes()

	serviceMock := NewMockServiceProvider(ctrl)
	serviceMock.EXPECT().EmptyCart(gomock.Any(), int64(1), int64(1)).Return(nil)
	serviceMock.EXPECT().EmptyCart(gomock.Any(), int64(1), int64(2)).Return(service.ErrTimerNotFound)
	serviceMock.EXPECT().EmptyCart(gomock.Any(), int64(1), int64(3)).Return(assert.AnError)

	testsCases := []tests.TestCase{
		{
			Name:           "ok - 204",
			Method:         http.MethodDelete,
			Target:         "/carts/1/items",
			AccessKey:      "abc123456",
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Name:           "cart not found - 404",
			Method:         http.MethodDelete,
			Target:         "/carts/2/items",
			AccessKey:      "abc123456",
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name:           "cart id string - invalid param",
			Method:         http.MethodDelete,
			Target:         "/carts/cart/items",
			AccessKey:      "abc123456",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":100422, "details":"Invalid params - cart_id param is not a valid number"}}`,
		},
		{
			Name:           "storage error - 500",
			Method:         http.MethodDelete,
			Target:         "/carts/3/items",
			AccessKey:      "abc123456",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":100500, "details":"Internal error - could not empty cart"}}`,
		},
	}
	execHTTPTestCases(t, serviceMock, authMock, testsCases)
}
