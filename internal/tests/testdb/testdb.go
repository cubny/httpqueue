package testdb

import (
	"context"
	"github.com/cubny/cart"
	"github.com/cubny/cart/internal/service"
)

type Storage interface {
	Migrate() error
	TruncateAllTables() error
}

type TestDB struct {
	storage Storage
	service *service.Service
}

func New(storage Storage, service *service.Service) *TestDB {
	return &TestDB{
		storage: storage,
		service: service,
	}
}

// Refresh changes the database to a fresh new database
func (t *TestDB) Refresh() error {
	if err := t.storage.Migrate(); err != nil {
		return err
	}

	if err := t.storage.TruncateAllTables(); err != nil {
		return err
	}

	return nil
}

func (t *TestDB) Seed1Cart(userID int64) (int64, error) {
	c, err := t.service.CreateCart(context.TODO(), userID)
	if err != nil {
		return 0, err
	}
	return c.ID, err
}

func (t *TestDB) Seed1Item(userID, cartID int64) (int64, error) {
	item := &cart.Item{
		ProductID: 1,
		CartID:    cartID,
		Quantity:  1,
		Price:     cart.Price(100.00),
	}
	if err := t.service.AddItem(context.TODO(), userID, item); err != nil {
		return 0, err
	}

	return item.ID, nil
}

func (t *TestDB) Seed5Items(userID, cartID int64) error {
	num := 5
	for i := 1; i <= num; i++ {
		item := &cart.Item{
			ProductID: int64(i),
			CartID:    cartID,
			Quantity:  1,
			Price:     cart.Price(100.00),
		}
		if err := t.service.AddItem(context.TODO(), userID, item); err != nil {
			return err
		}
	}

	return nil
}
