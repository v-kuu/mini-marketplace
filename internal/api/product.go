package api

import (
	"encoding/json"
	"net/http"
	"github.com/v-kuu/mini-marketplace/internal/model"
)

type ProductLister interface {
	ListProducts() ([]model.Product, error)
}

type ProductHandler struct {
	service ProductLister
}

func NewProductHandler(s ProductLister) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) List(w http.ResponseWriter, _ *http.Request) {
	products, err := h.service.ListProducts()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
