package repository

import "github.com/v-kuu/mini-marketplace/internal/model"

type ProductRepository struct {}

func (r *ProductRepository) List() ([]model.Product, error) {
	return []model.Product {
		{ID: "1", Name: "Coffee", Price: 499},
		{ID: "2", Name: "Sandwich", Price: 899},
	},
	nil
}
