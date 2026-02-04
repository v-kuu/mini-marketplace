package api

import (
	"encoding/json"
	"net/http"
	"context"
	"strings"
	"log"
	"time"
	"errors"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
)

type ProductService interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	CreateProduct(ctx context.Context, p model.Product) error
	UpdateProduct(ctx context.Context, id string, p model.Product) error
	PatchProduct(ctx context.Context, id string, patch service.ProductPatch) error
	DeleteProduct(ctx context.Context, id string) error
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
			writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	products, err := h.service.ListProducts(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			writeJSONError(w, "Request timeout", http.StatusRequestTimeout)
		} else {
			log.Printf("ListProducts: %v", err)
			writeJSONError(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSONError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	err := h.service.CreateProduct(ctx, p)
	if err != nil {
		switch err {
			case service.ErrInvalidProduct:
				writeJSONError(w, err.Error(), http.StatusBadRequest)
			case service.ErrProductAlreadyExists:
				writeJSONError(w, err.Error(), http.StatusConflict)
			case context.DeadlineExceeded:
				writeJSONError(w, "Request timeout", http.StatusRequestTimeout)
			default:
				log.Printf("CreateProduct: %v", err)
				writeJSONError(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}

func (h *ProductHandler) ProductByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	switch r.Method {
		case http.MethodGet:
			h.getProduct(w, r, id)
		case http.MethodPut:
			h.updateProduct(w, r, id)
		case http.MethodPatch:
			h.patchProduct(w, r, id)
		case http.MethodDelete:
			h.deleteProduct(w, r, id)
		default:
			writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (h *ProductHandler) getProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			writeJSONError(w, "Request timeout", http.StatusRequestTimeout)
		} else {
			log.Printf("GetProduct: %v", err)
			writeJSONError(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	if product == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}

func (h *ProductHandler) updateProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateProduct(ctx, id, p)
	if err != nil {
		switch err {
		case service.ErrInvalidProduct:
			writeJSONError(w, err.Error(), http.StatusBadRequest)
		case service.ErrIDMismatch:
			writeJSONError(w, err.Error(), http.StatusBadRequest)
		case service.ErrProductNotFound:
			writeJSONError(w, err.Error(), http.StatusNotFound)
		case context.DeadlineExceeded:
			writeJSONError(w, "Request timeout", http.StatusRequestTimeout)
		default:
			log.Printf("UpdateProduct: %v", err)
			writeJSONError(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	updated := p
	updated.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}

func (h *ProductHandler) patchProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	var patch service.ProductPatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeJSONError(w, "Invalid json", http.StatusBadRequest)
		return
	}

	err := h.service.PatchProduct(ctx, id ,patch)
	if err != nil {
		switch err {
			case service.ErrInvalidProduct:
				writeJSONError(w, err.Error(), http.StatusBadRequest)
			case service.ErrProductNotFound:
				writeJSONError(w, err.Error(), http.StatusNotFound)
			case context.DeadlineExceeded:
				writeJSONError(w, "Request timeout", http.StatusRequestTimeout)
			default:
				log.Printf("PatchProduct: %v", err)
				writeJSONError(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	updated, err := h.service.GetProduct(ctx, id)
	if err != nil {
		writeJSONError(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}

func (h *ProductHandler) deleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	err := h.service.DeleteProduct(ctx, id)
	if err != nil {
		switch err {
			case service.ErrInvalidProduct:
				writeJSONError(w, err.Error(), http.StatusBadRequest)
			case service.ErrProductNotFound:
				writeJSONError(w, err.Error(), http.StatusNotFound)
			case context.DeadlineExceeded:
				writeJSONError(w, "Request timeout", http.StatusRequestTimeout)
			default:
				log.Printf("DeleteProduct: %v", err)
				writeJSONError(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorResponse{message}); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}
