//go:generate mockgen -package service -source cart.go -destination service_mocks_test.go
package service_test

import (
	"context"
	"testing"

	"github.com/cubny/cart/internal/service"
	"github.com/cubny/cart/internal/storage"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateCart(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		adjust        func(db *service.MockStorage)
		expectedError error
	}{
		{
			name:          "ok",
			expectedError: nil,
			userID:        int64(1),
			adjust: func(db *service.MockStorage) {
				db.EXPECT().
					CreateCart(gomock.Any(), &cart.Cart{UserID: int64(1)}).
					Return(nil)
			},
		},
		{
			name:          "invalid userID",
			userID:        int64(0),
			expectedError: cart.ErrInvalidUserID,
			adjust: func(db *service.MockStorage) {
			},
		},
		{
			name:          "repo error - bubbles up",
			expectedError: assert.AnError,
			userID:        int64(1),
			adjust: func(db *service.MockStorage) {
				db.EXPECT().
					CreateCart(gomock.Any(), &cart.Cart{UserID: int64(1)}).
					Return(assert.AnError)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbMock := service.NewMockStorage(ctrl)
			test.adjust(dbMock)
			svc, err := service.New(dbMock)
			assert.Nil(t, err)
			_, err = svc.CreateCart(context.TODO(), test.userID)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestService_AddItem(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		item          *cart.Item
		expectedError error
		adjust        func(db *service.MockStorage, item *cart.Item, userID int64)
	}{
		{
			name:   "ok",
			userID: 1,
			item: &cart.Item{
				CartID:    1,
				ProductID: 1,
				Price:     cart.Price(10.00),
				Quantity:  1,
			},
			expectedError: nil,
			adjust: func(db *service.MockStorage, item *cart.Item, userID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, item.CartID).Return(&cart.Cart{}, nil)
				db.EXPECT().FindItemByProductID(gomock.Any(), item.CartID, item.ProductID).Return(nil, nil)
				db.EXPECT().CreateItem(gomock.Any(), item).Return(nil)
			},
		},
		{
			name:   "cart does not belong to the user - ErrTimerNotFound",
			userID: 1,
			item: &cart.Item{
				CartID:    2,
				ProductID: 1,
				Price:     cart.Price(10.00),
				Quantity:  1,
			},
			expectedError: service.ErrTimerNotFound,
			adjust: func(db *service.MockStorage, item *cart.Item, userID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, item.CartID).Return(nil, storage.ErrRecordNotFound)
			},
		},
		{
			name:   "repo returns error on GetCart - error",
			userID: 1,
			item: &cart.Item{
				CartID:    2,
				ProductID: 1,
				Price:     cart.Price(10.00),
				Quantity:  1,
			},
			expectedError: assert.AnError,
			adjust: func(db *service.MockStorage, item *cart.Item, userID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, item.CartID).Return(nil, assert.AnError)
			},
		},
		{
			name:   "product already exists - ErrProductAlreadyInCart",
			userID: 1,
			item: &cart.Item{
				CartID:    1,
				ProductID: 1,
				Price:     cart.Price(10.00),
				Quantity:  1,
			},
			expectedError: service.ErrProductAlreadyInCart,
			adjust: func(db *service.MockStorage, item *cart.Item, userID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, item.CartID).Return(&cart.Cart{}, nil)
				db.EXPECT().FindItemByProductID(gomock.Any(), item.CartID, item.ProductID).
					Return(&cart.Item{}, nil)
			},
		},
		{
			name:   "repo returns error on createItem - error",
			userID: 1,
			item: &cart.Item{
				CartID:    1,
				ProductID: 1,
				Price:     cart.Price(10.00),
				Quantity:  1,
			},
			expectedError: assert.AnError,
			adjust: func(db *service.MockStorage, item *cart.Item, userID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, item.CartID).Return(&cart.Cart{}, nil)
				db.EXPECT().FindItemByProductID(gomock.Any(), item.CartID, item.ProductID).
					Return(nil, nil)
				db.EXPECT().CreateItem(gomock.Any(), item).Return(assert.AnError)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbMock := service.NewMockStorage(ctrl)
			test.adjust(dbMock, test.item, test.userID)
			svc, err := service.New(dbMock)
			assert.Nil(t, err)
			err = svc.AddItem(context.TODO(), test.userID, test.item)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestService_RemoveItem(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		itemID        int64
		expectedError error
		adjust        func(db *service.MockStorage, userID, itemID int64)
	}{
		{
			name:          "ok",
			userID:        1,
			itemID:        1,
			expectedError: nil,
			adjust: func(db *service.MockStorage, userID, itemID int64) {
				db.EXPECT().GetItem(gomock.Any(), itemID).Return(&cart.Item{CartID: 1, ID: itemID}, nil)
				db.EXPECT().GetCart(gomock.Any(), userID, int64(1)).Return(&cart.Cart{}, nil)
				db.EXPECT().RemoveItem(gomock.Any(), itemID).Return(nil)
			},
		},
		{
			name:          "item not found - ErrItemNotFound",
			userID:        1,
			itemID:        1,
			expectedError: service.ErrItemNotFound,
			adjust: func(db *service.MockStorage, userID, itemID int64) {
				db.EXPECT().GetItem(gomock.Any(), itemID).Return(nil, storage.ErrRecordNotFound)
			},
		},
		{
			name:          "repo getItem returns error - error",
			userID:        1,
			itemID:        1,
			expectedError: assert.AnError,
			adjust: func(db *service.MockStorage, userID, itemID int64) {
				db.EXPECT().GetItem(gomock.Any(), itemID).Return(nil, assert.AnError)
			},
		},
		{
			name:          "repo getCart returns record not found - ErrTimerNotFound",
			userID:        1,
			itemID:        1,
			expectedError: service.ErrTimerNotFound,
			adjust: func(db *service.MockStorage, userID, itemID int64) {
				db.EXPECT().GetItem(gomock.Any(), itemID).Return(&cart.Item{CartID: 1, ID: itemID}, nil)
				db.EXPECT().GetCart(gomock.Any(), userID, int64(1)).Return(nil, storage.ErrRecordNotFound)
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbMock := service.NewMockStorage(ctrl)
			test.adjust(dbMock, test.userID, test.itemID)
			svc, err := service.New(dbMock)
			assert.Nil(t, err)
			err = svc.RemoveItem(context.TODO(), test.userID, test.itemID)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestService_EmptyCart(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		cartID        int64
		expectedError error
		adjust        func(db *service.MockStorage, userID, cartID int64)
	}{
		{
			name:          "ok",
			userID:        1,
			cartID:        1,
			expectedError: nil,
			adjust: func(db *service.MockStorage, userID, cartID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, cartID).Return(&cart.Cart{}, nil)
				db.EXPECT().RemoveItemsByCartID(gomock.Any(), cartID).Return(nil)
			},
		},
		{
			name:          "repo getCart returns ErrTimerNotFound - ErrTimerNotFound",
			userID:        1,
			cartID:        1,
			expectedError: service.ErrTimerNotFound,
			adjust: func(db *service.MockStorage, userID, cartID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, cartID).Return(nil, storage.ErrRecordNotFound)
			},
		},
		{
			name:          "repo getCart returns error - error",
			userID:        1,
			cartID:        1,
			expectedError: assert.AnError,
			adjust: func(db *service.MockStorage, userID, cartID int64) {
				db.EXPECT().GetCart(gomock.Any(), userID, cartID).Return(nil, assert.AnError)
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbMock := service.NewMockStorage(ctrl)
			test.adjust(dbMock, test.userID, test.cartID)
			svc, err := service.New(dbMock)
			assert.Nil(t, err)
			err = svc.EmptyCart(context.TODO(), test.userID, test.cartID)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
