package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/semaphore"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
)

type ProductRepository struct {
	db *sql.DB
	sem *semaphore.Weighted
}

func OpenDB(maxOpen int64) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:products.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(maxOpen))
	db.SetMaxIdleConns(int(maxOpen) / 2)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func NewProductRepository(db *sql.DB, maxConcurrent int64) *ProductRepository {
	return &ProductRepository{
		db: db,
		sem: semaphore.NewWeighted(maxConcurrent),
	}
}

func (r *ProductRepository) List(ctx context.Context) ([]model.Product, error) {
	rows, err := r.query(
		ctx,
		`SELECT id, name, price FROM products`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	row, err := r.queryRow(
		ctx,
		`SELECT id, name, price FROM products WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}

	var p model.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p model.Product) error {
	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := r.exec(
			ctx,
			tx,
			`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
			p.ID, p.Name, p.Price,
		)
		return err
	})
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		res, err := r.exec(
			ctx,
			tx,
			`DELETE FROM products WHERE id = ?`,
			id,
		)
		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rows == 0 {
			return service.ErrProductNotFound
		}

		return nil
	})
}

func (r *ProductRepository) Update(ctx context.Context, p model.Product) error {
	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(
			ctx,
			`SELECT id, name, price FROM products WHERE id = ?`,
			p.ID,
		)

		var prev model.Product
		if err := row.Scan(&prev.ID, &prev.Name, &prev.Price); errors.Is(err, sql.ErrNoRows) {
			return service.ErrProductNotFound
		} else if err != nil {
			return err
		}
		if p.Name == "" {
			p.Name = prev.Name
		}
		if p.Price <= 0 {
			p.Price = prev.Price
		}

		res, err := r.exec(
			ctx,
			tx,
			`UPDATE products SET name = ?, price = ? WHERE id = ?`,
			p.Name, p.Price, p.ID,
		)

		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rows == 0 {
			return service.ErrProductNotFound
		}

		return nil
	})
}

func (r *ProductRepository) query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if err := r.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer r.sem.Release(1)

	return r.db.QueryContext(ctx, query, args...)
}

func (r *ProductRepository) exec(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	if err := r.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer r.sem.Release(1)

	return tx.ExecContext(ctx, query, args...)
	
}

func (r *ProductRepository) queryRow(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	if err := r.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer r.sem.Release(1)

	return r.db.QueryRowContext(ctx, query, args...), nil
}
