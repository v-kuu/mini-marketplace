package api

import (
	"encoding/json"
	"net/http"
	"github.com/v-kuu/mini-marketplace/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) List(w http.ResponseWriter, _ *http.Request) {
	products := h.service.ListProducts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
