package sqlite

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) List(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(
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
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, price FROM products WHERE id = ?`,
		id,
	)

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
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
		p.ID, p.Name, p.Price,
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(
		ctx,
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
}

func (r *ProductRepository) Update(ctx context.Context, p model.Product) error {
	res, err := r.db.ExecContext(
		ctx,
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
}
