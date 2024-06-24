package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MySQLPersistence struct {
	db *sqlx.DB
}

func (m *MySQLPersistence) FetchProduct(ctx context.Context) (Product, error) {
	var p []Product
	err := m.db.SelectContext(ctx, &p,
		`select * from products
where active = 1
limit 1`)
	if err != nil {
		return Product{}, err
	}
	if len(p) == 0 {
		return Product{}, fmt.Errorf("no active product found")
	}
	return p[0], err
}

func (m *MySQLPersistence) InsertProduct(ctx context.Context, p Product) error {
	// deactive previous version of the product
	query := "UPDATE products SET active = 0"
	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	// insert new product
	query = "INSERT INTO products (date, active, pack_sizes) VALUES (?, ?, ?)"
	_, err = m.db.ExecContext(ctx, query, &p.Date, &p.Active, &p.PackSizes)
	return err
}

func NewMySQLPersistence(db *sqlx.DB) *MySQLPersistence {
	return &MySQLPersistence{
		db: db,
	}
}
