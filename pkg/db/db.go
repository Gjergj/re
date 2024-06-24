package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Persistence interface {
	FetchProduct(ctx context.Context) (Product, error)
	InsertProduct(ctx context.Context, p Product) error
}

type Product struct {
	ID        int       `db:"id"`
	Date      time.Time `db:"date"`
	Active    int       `db:"active"`
	PackSizes PackSizes `db:"pack_sizes"`
}

type PackSizes []int

func (pc *PackSizes) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &pc)
		return nil
	case []int:
		*pc = v
		return nil
	case []string:
		for _, s := range v {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			*pc = append(*pc, i)
		}
		return nil
	case string:
		json.Unmarshal([]byte(v), &pc)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (pc *PackSizes) Value() (driver.Value, error) {
	return json.Marshal(pc)
}

func New(username, password, host, dbName string) (*sqlx.DB, error) {
	connStr := "%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true"
	connStr = fmt.Sprintf(connStr, username, password, host, dbName)

	db, err := sqlx.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDb(dbName string, d *sql.DB) error {
	driver, err := mysql.WithInstance(d, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}
