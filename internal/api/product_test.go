package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type fakeProductService struct {
	products []model.Product
	err error
}

func (f *fakeProductService) ListProducts() ([]model.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func testProductHandler_List(t *testing.T) {
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

			handler.List(rec, req)

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
