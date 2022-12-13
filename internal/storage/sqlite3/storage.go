package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cubny/cart"
	"github.com/cubny/cart/internal/storage"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3 struct {
	db *sql.DB
}

func New(dbfile *os.File) (*Sqlite3, error) {
	db, err := sql.Open("sqlite3", dbfile.Name())
	if err != nil {
		return nil, err
	}
	return &Sqlite3{db: db}, nil
}

func (s *Sqlite3) Close() error {
	return s.db.Close()
}

func (s *Sqlite3) CreateCart(ctx context.Context, cart *cart.Cart) error {
	stmt, err := s.db.Prepare(queryInsertCart)
	if err != nil {
		return err
	}

	now := time.Now()
	cart.CreatedAt = now
	cart.UpdatedAt = now

	res, err := stmt.ExecContext(ctx, cart.UserID, cart.CreatedAt, cart.UpdatedAt)
	if err != nil {
		return err
	}

	cart.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3) GetCart(ctx context.Context, userID, cartID int64) (*cart.Cart, error) {
	stmt, err := s.db.Prepare(queryCartsByIDAndUserID)
	if err != nil {
		return nil, err
	}

	c := &cart.Cart{}

	rows, err := stmt.QueryContext(ctx, cartID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&c.ID, &c.UserID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("sqlite3: GetCart result scan error, %s", err)
		}
		return c, nil
	}

	return nil, storage.ErrRecordNotFound
}

func (s *Sqlite3) FindItemByProductID(ctx context.Context, cartID, productID int64) (*cart.Item, error) {
	stmt, err := s.db.Prepare(queryItemsByCartIDAndProductID)
	if err != nil {
		return nil, err
	}

	item := &cart.Item{}

	rows, err := stmt.QueryContext(ctx, cartID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("sqlite3: FindItemByProductID result scan error, %s", err)
		}
		return item, nil
	}

	return nil, storage.ErrRecordNotFound
}

func (s *Sqlite3) CreateItem(ctx context.Context, item *cart.Item) error {
	stmt, err := s.db.Prepare(queryInsertItem)
	if err != nil {
		return err
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	res, err := stmt.ExecContext(ctx,
		item.CartID,
		item.ProductID,
		item.Quantity,
		item.Price,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return err
	}

	item.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3) GetItem(ctx context.Context, itemID int64) (*cart.Item, error) {
	stmt, err := s.db.Prepare(queryItemByID)
	if err != nil {
		return nil, err
	}

	item := &cart.Item{}

	rows, err := stmt.QueryContext(ctx, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("sqlite3: GetItem result scan error, %s", err)
		}
		return item, nil
	}

	return nil, storage.ErrRecordNotFound
}

func (s *Sqlite3) RemoveItem(ctx context.Context, itemID int64) error {
	_, err := s.db.ExecContext(ctx, queryRemoveItem, itemID)
	return err
}

func (s *Sqlite3) RemoveItemsByCartID(ctx context.Context, cartID int64) error {
	_, err := s.db.ExecContext(ctx, queryRemoveItemsByCartID, cartID)
	return err
}
