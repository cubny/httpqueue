package tests_test

import (
	"fmt"
	"github.com/cubny/cart/internal/tests"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestCreateCart_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testsCases := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         "/carts",
			AccessKey:      "abcdef123456",
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

	for _, test := range testsCases {
		tests.HandlerTest(t, a, &test)
	}
}

func TestAddItems_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// create a cart
	userID := int64(1)
	cartID, _ := testDB.Seed1Cart(userID)

	testsCases := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         fmt.Sprintf("/carts/%d/items", cartID),
			AccessKey:      "abcdef123456",
			ReqBody:        `{"product_id":1, "quantity":1, "price": 100.00}`,
			ExpectedBody:   `{"cart_id":` + strconv.Itoa(int(cartID)) + `, "id":1, "price":100, "product_id":1, "quantity":1}`,
			ExpectedStatus: http.StatusCreated,
		},
	}

	for _, test := range testsCases {
		test := test
		tests.HandlerTest(t, a, &test)
	}
}

func TestRemoveItems_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// create an item
	userID := int64(1)
	cartID, _ := testDB.Seed1Cart(userID)
	itemID, _ := testDB.Seed1Item(userID, cartID)

	testsCases := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodDelete,
			Target:         fmt.Sprintf("/items/%d", itemID),
			AccessKey:      "abcdef123456",
			ExpectedStatus: http.StatusNoContent,
		},
	}

	for _, test := range testsCases {
		test := test
		tests.HandlerTest(t, a, &test)
	}
}

func TestEmptyCart_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// create an item
	userID := int64(1)
	cartID, _ := testDB.Seed1Cart(userID)
	err := testDB.Seed5Items(userID, cartID)
	assert.Nil(t, err)

	testsCases := []tests.TestCase{
		{
			Name:           "ok",
			Method:         http.MethodDelete,
			Target:         fmt.Sprintf("/carts/%d/items", cartID),
			AccessKey:      "abcdef123456",
			ExpectedStatus: http.StatusNoContent,
		},
	}

	for _, test := range testsCases {
		test := test
		tests.HandlerTest(t, a, &test)
	}
}
