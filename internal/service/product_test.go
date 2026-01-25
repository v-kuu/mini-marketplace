package service

import (
	"github.com/v-kuu/mini-marketplace/internal/model"
	"testing"
	"errors"
	"context"
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

func TestProductService_ListProducts(t *testing.T) {
	repoErr := errors.New("Repository failure")

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
				err: repoErr,
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
				t.Fatalf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !tt.wantErr && len(products) != tt.wantLen {
				t.Fatalf("Expected %d products, got %d", tt.wantLen, len(products))
			}
		})
	}
}

