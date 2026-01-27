package service

import (
	"context"
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductRepository interface {
	List(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]model.Product, error) {
	select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
	}

	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}
