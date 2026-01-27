package api

import (
	"encoding/json"
	"net/http"
	"context"
	"strings"

	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductService interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(s ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) Products(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) ProductByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
		case http.MethodGet:
			h.getProduct(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (h *ProductHandler) getProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(product)
}
