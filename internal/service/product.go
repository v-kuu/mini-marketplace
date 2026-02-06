package service

import (
	"context"

	"github.com/google/uuid"

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

func (s *ProductService) CreateProduct(ctx context.Context, name string, price int64) (string, error) {
	select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
	}

	var id string
	existing := &model.Product{}
	for existing != nil {
		select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
		}
		id = uuid.New().String()
		new, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return id, err
		}
		existing = new
	}

	p := model.Product{ID: id, Name: name, Price: price}
	return id, s.repo.Create(ctx, p)
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

func (s *ProductService) UpdateProduct(ctx context.Context, id string, name string, price int64) error {
	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	if id == "" {
		return ErrInvalidProduct
	}

	p := model.Product{ID: id, Name: name, Price: price}
	return s.repo.Update(ctx, p)
}

func (s *ProductService) PatchProduct(ctx context.Context, id string, name *string, price *int64) error {
	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	if id == "" {
		return ErrInvalidProduct
	}
	p := model.Product{ID: id}
	if name != nil {
		p.Name = *name
	}
	if price != nil {
		p.Price = *price
	}
	err := s.repo.Update(ctx, p)
	return err
}
