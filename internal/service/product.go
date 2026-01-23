package service

import (
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductRepository interface {
	List() []model.Product
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts() []model.Product {
	return s.repo.List()
}
