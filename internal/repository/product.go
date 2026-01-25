package repository

import (
	"context"
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductRepository struct {}

func (r *ProductRepository) List(ctx context.Context) ([]model.Product, error) {
	select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
	}

	return []model.Product {
		{ID: "1", Name: "Coffee", Price: 499},
		{ID: "2", Name: "Sandwich", Price: 899},
	},
	nil
}
