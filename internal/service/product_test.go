package service

import (
	"testing"
	"errors"
	"context"

	"github.com/v-kuu/mini-marketplace/internal/model"
)

type fakeProductRepo struct {
	products []model.Product
	err error
}

func (f *fakeProductRepo) List(ctx context.Context) ([]model.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func (f *fakeProductRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, product := range f.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, nil
}

func (f *fakeProductRepo) Create(ctx context.Context, p model.Product) error {
	if p.ID == "" || p.Name == "" || p.Price <= 0 {
		return ErrInvalidProduct
	}

	for _, product := range f.products {
		if product.ID == p.ID {
			return ErrProductAlreadyExists
		}
	}

	f.products = append(f.products, p)
	return nil
}

func (f *fakeProductRepo) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidProduct
	}

	for i, product := range f.products {
		if product.ID == id {
			f.products = append(f.products[:i], f.products[i+1:]...)
			return nil
		}
	}

	return ErrProductNotFound
}

func (f *fakeProductRepo) Update(ctx context.Context, p model.Product) error {
	if p.ID == "" || p.Name == "" || p.Price <= 0 {
		return ErrInvalidProduct
	}

	for i, product := range f.products {
		if product.ID == p.ID {
			f.products[i] = p
			return nil
		}
	}
	return ErrProductNotFound
}

func TestProductService_ListProducts(t *testing.T) {

	tests := []struct {
		name string
		repo *fakeProductRepo
		wantLen int
		wantErr bool
	}{
		{
			name: "Returns all products",
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Returns empty list when no products",
			repo: &fakeProductRepo{
				products: []model.Product{},
			},
			wantLen:  0,
			wantErr: false,
		},
		{
			name: "Repository error",
			repo: &fakeProductRepo{
				err: errors.New("repository failure"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewProductService(tt.repo)
			products, err := svc.ListProducts(context.Background())

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && len(products) != tt.wantLen {
				t.Fatalf("expected %d products, got %d", tt.wantLen, len(products))
			}
		})
	}
}

func TestProductService_GetProduct(t *testing.T) {

	tests := []struct {
		name string
		repo *fakeProductRepo
		wantLen int
		wantErr bool
	}{
		{
			name: "Returns Coffee",
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Returns not found",
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen:  2,
			wantErr: false,
		},
		{
			name: "Repository error",
			repo: &fakeProductRepo{
				err: errors.New("repository failure"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewProductService(tt.repo)
			product, err := svc.GetProduct(context.Background(), "1")

			if !tt.wantErr && product.Name != "Coffee" {
				t.Fatalf("wrong result")
			}
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestProductService_Create(t *testing.T) {

	tests := []struct {
		name string
		p model.Product
		repo *fakeProductRepo
		wantLen int
		wantErr bool
	}{
		{
			name: "Success",
			p: model.Product{ID: "3", Name: "Tea", Price: 499},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "Invalid Request",
			p: model.Product{ID: "", Name: "", Price: 0},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: true,
		},
		{
			name: "Already exists",
			p: model.Product{ID: "1", Name: "Coffee", Price: 499},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewProductService(tt.repo)
			err := svc.CreateProduct(context.Background(), tt.p)

			if !tt.wantErr && tt.wantLen != len(tt.repo.products) {
				t.Fatalf("expected %d products, got %d", tt.wantLen, len(tt.repo.products))
			}
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestProductService_Delete(t *testing.T) {

	tests := []struct {
		name string
		id string
		repo *fakeProductRepo
		wantLen int
		wantErr bool
	}{
		{
			name: "Success",
			id: "2",
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Not found",
			id: "3",
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen:  2,
			wantErr: true,
		},
		{
			name: "Invalid id",
			id: "",
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewProductService(tt.repo)
			err := svc.DeleteProduct(context.Background(), tt.id)

			if tt.wantLen != len(tt.repo.products) {
				t.Fatalf("expected %d elements, got %d", tt.wantLen, len(tt.repo.products))
			}
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestProductService_Update(t *testing.T) {

	tests := []struct {
		name string
		id string
		p model.Product
		repo *fakeProductRepo
		wantLen int
		wantErr bool
	}{
		{
			name: "Success",
			id: "1",
			p: model.Product{ID: "1", Name: "Tea", Price: 599},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Not found",
			id: "3",
			p: model.Product{ID: "3", Name: "Tea", Price: 599},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen:  2,
			wantErr: true,
		},
		{
			name: "Invalid id",
			id: "",
			p: model.Product{ID: "", Name: "", Price: 0},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen: 2,
			wantErr: true,
		},
		{
			name: "ID Mismatch",
			id: "1",
			p: model.Product{ID: "2", Name: "Tea", Price: 599},
			repo: &fakeProductRepo{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantLen:  2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewProductService(tt.repo)
			err := svc.UpdateProduct(context.Background(), tt.id, tt.p)

			if tt.wantLen != len(tt.repo.products) {
				t.Fatalf("expected %d elements, got %d", tt.wantLen, len(tt.repo.products))
			}
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.repo.products[0].Name != "Coffee" {
					t.Fatalf("expected coffee, got %s", tt.repo.products[0].Name)
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.repo.products[0].Name != "Tea" {
					t.Fatalf("expected tea, got %s", tt.repo.products[0].Name)
				}
			}
		})
	}
}
