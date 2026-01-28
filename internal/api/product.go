package api

import (
	"encoding/json"
	"net/http"
	"context"
	"strings"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
)

type ProductService interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	CreateProduct(ctx context.Context, p model.Product) error
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, id string, p model.Product) error
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(s ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) Products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			h.listProducts(w, r)
		case http.MethodPost:
			h.createProduct(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.service.CreateProduct(ctx, p)
	if err != nil {
		switch err {
			case service.ErrInvalidProduct:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case service.ErrProductAlreadyExists:
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) ProductByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	switch r.Method {
		case http.MethodGet:
			h.getProduct(w, r, id)
		case http.MethodPut:
			h.updateProduct(w, r, id)
		case http.MethodDelete:
			h.deleteProduct(w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func (h *ProductHandler) getProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) updateProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}

	err := h.service.UpdateProduct(ctx, id, p)
	if err != nil {
		switch err {
		case service.ErrInvalidProduct:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case service.ErrIDMismatch:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case service.ErrProductNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}

	updated := p
	updated.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func (h *ProductHandler) deleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	err := h.service.DeleteProduct(ctx, id)
	if err != nil {
		switch err {
			case service.ErrInvalidProduct:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case service.ErrProductNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
