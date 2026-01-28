package service

import (
	"context"

	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductRepository interface {
	List(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id string) (*model.Product, error)
	Create(ctx context.Context, p model.Product) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, p model.Product) error
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
	select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, p model.Product) error {
	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	if p.ID == "" || p.Name == "" || p.Price <= 0 {
		return ErrInvalidProduct
	}

	existing, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		return err
	}

	if existing != nil {
		return ErrProductAlreadyExists
	}

	return s.repo.Create(ctx, p)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	if id == "" {
		return ErrInvalidProduct
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, p model.Product) error {
	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	if id == "" || p.Name == "" || p.Price <= 0 {
		return ErrInvalidProduct
	}

	if p.ID != "" && p.ID != id {
		return ErrIDMismatch
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}
	
	p.ID = id
	return s.repo.Update(ctx, p)
}

func (s *ProductService) PatchProduct(ctx context.Context, id string, patch ProductPatch) error {
	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	if id == "" {
		return ErrInvalidProduct
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}

	if patch.Name != nil {
		if *patch.Name == "" {
			return ErrInvalidProduct
		}
		existing.Name = *patch.Name
	}

	if patch.Price != nil {
		if *patch.Price <= 0 {
			return ErrInvalidProduct
		}
		existing.Price = *patch.Price
	}

	return s.repo.Update(ctx, *existing)
}
