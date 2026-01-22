package service

import (
	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts() []model.Product {
	return s.repo.List()
}
