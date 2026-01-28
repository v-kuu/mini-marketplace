package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"context"
	"strings"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
)

type fakeProductService struct {
	products []model.Product
	err error
}

func (f *fakeProductService) ListProducts(ctx context.Context) ([]model.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func (f *fakeProductService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, p := range f.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, errors.New("not found")
}

func (f *fakeProductService) CreateProduct(ctx context.Context, p model.Product) error {
	if p.ID == "" || p.Name == "" || p.Price <= 0 {
		return service.ErrInvalidProduct
	}

	for _, product := range f.products {
		if product.ID == p.ID {
			return service.ErrProductAlreadyExists
		}
	}

	f.products = append(f.products, p)
	return nil
}

func (f *fakeProductService) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrInvalidProduct
	}

	for i, product := range f.products {
		if product.ID == id {
			f.products = append(f.products[:i], f.products[i+1:]...)
			return nil
		}
	}

	return service.ErrProductNotFound
}

func TestProductHandler_List(t *testing.T) {
	tests := []struct {
		name string
		service *fakeProductService
		wantStatus int
		wantLen int
	}{
		{
			name: "Success",
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantStatus: http.StatusOK,
			wantLen: 2,
		},
		{
			name: "Service error",
			service: &fakeProductService{
				err: errors.New("Failure"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := NewProductHandler(tt.service)

			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()

			handler.Products(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			if tt.wantStatus == http.StatusOK {
				var products []model.Product
				if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}	
				if len(products) != tt.wantLen {
					t.Fatalf("Expected %d products, got %d", tt.wantLen, len(products))
				}
			}
		})
	}
}

func TestProductHandler_Get(t *testing.T) {
	tests := []struct {
		name string
		service *fakeProductService
		wantStatus int
		wantLen int
	}{
		{
			name: "Success",
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantStatus: http.StatusOK,
			wantLen: 2,
		},
		{
			name: "Not found",
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Tea", Price: 499},
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantLen: 1,
		},
		{
			name: "Service error",
			service: &fakeProductService{
				err: errors.New("Failure"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := NewProductHandler(tt.service)

			req := httptest.NewRequest(http.MethodGet, "/products/2", nil)
			rec := httptest.NewRecorder()

			handler.ProductByID(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			if tt.wantStatus == http.StatusOK {
				var product model.Product
				if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if product.Name != "Sandwich" {
					t.Fatalf("Expected Sandwich, got %s", product.Name)
				}
			}
		})
	}
}

func TestProductHandler_Create(t *testing.T) {
	tests := []struct {
		name string
		body string
		service *fakeProductService
		wantStatus int
		wantLen int
	}{
		{
			name: "Success",
			body : `{"id":"3","name":"Tea","price":499}`,
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantStatus: http.StatusCreated,
			wantLen: 3,
		},
		{
			name: "Already exists",
			body : `{"id":"3","name":"Tea","price":499}`,
			service: &fakeProductService{
				products: []model.Product{
					{ID: "3", Name: "Tea", Price: 499},
				},
			},
			wantStatus: http.StatusConflict,
			wantLen: 1,
		},
		{
			name: "Invalid Product",
			body : `{"id":"","name":"","price":0}`,
			service: &fakeProductService{
				products: []model.Product{
					{ID: "3", Name: "Tea", Price: 499},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := NewProductHandler(tt.service)

			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Products(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			if tt.wantStatus == http.StatusOK {
				var products []model.Product
				if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}	
				if len(products) != tt.wantLen {
					t.Fatalf("Expected %d products, got %d", tt.wantLen, len(products))
				}
			}
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	tests := []struct {
		name string
		id string
		service *fakeProductService
		wantStatus int
		wantLen int
	}{
		{
			name: "Success",
			id: "2",
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Coffee", Price: 499},
					{ID: "2", Name: "Sandwich", Price: 899},
				},
			},
			wantStatus: http.StatusNoContent,
			wantLen: 1,
		},
		{
			name: "Not found",
			id: "2",
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Tea", Price: 499},
				},
			},
			wantStatus: http.StatusNotFound,
			wantLen: 1,
		},
		{
			name: "Invalid product",
			id: "",
			service: &fakeProductService{
				products: []model.Product{
					{ID: "1", Name: "Tea", Price: 499},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := NewProductHandler(tt.service)

			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.id, nil)
			rec := httptest.NewRecorder()

			handler.ProductByID(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("Expected status %d, got %d", tt.wantStatus, res.StatusCode)
			}

			if tt.wantLen != len(tt.service.products) {
				t.Fatalf("expected %d elements, got %d", tt.wantLen, len(tt.service.products))
			}
			if tt.wantStatus == http.StatusOK {
				var product model.Product
				if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
			}
		})
	}
}
