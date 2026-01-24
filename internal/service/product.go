package service

import (
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductRepository interface {
	List() ([]model.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts() ([]model.Product, error) {
	products, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return products, nil
}
