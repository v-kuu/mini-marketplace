package service

import "github.com/v-kuu/mini-marketplace/internal/model"
import "testing"

type fakeProductRepo struct {
	products []model.Product
}

func (f *fakeProductRepo) List() []model.Product {
	return f.products
}

func TestProductService_ListProducts(t *testing.T) {
	tests := []struct {
		name     string
		products []model.Product
		wantLen  int
	}{
		{
			name: "Returns all products",
			products: []model.Product{
				{ID: "1", Name: "Coffee", Price: 499},
				{ID: "2", Name: "Sandwich", Price: 899},
			},
			wantLen: 2,
		},
		{
			name:     "Returns empty list when no products",
			products: []model.Product{},
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeProductRepo{products: tt.products}
			svc := NewProductService(repo)

			got := svc.ListProducts()

			if len(got) != tt.wantLen {
				t.Fatalf("Expected %d products, got %d", tt.wantLen, len(got))
			}
		})
	}
}

