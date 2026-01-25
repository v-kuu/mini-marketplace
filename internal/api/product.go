package api

import (
	"encoding/json"
	"net/http"
	"context"
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductLister interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
}

type ProductHandler struct {
	service ProductLister
}

func NewProductHandler(s ProductLister) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
