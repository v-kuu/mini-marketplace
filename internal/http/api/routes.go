package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/repository/sqlite"
	"github.com/v-kuu/mini-marketplace/internal/metrics"
	"github.com/v-kuu/mini-marketplace/internal/http/middleware"
)

func AddRoutes() (*http.ServeMux, error) {
	metrics.Register()

	var maxOpen int64 = 100
	db, err := sqlite.OpenDB(maxOpen)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	repo := sqlite.NewProductRepository(db, maxOpen * 2)
	svc := service.NewProductService(repo)
	handler := NewProductHandler(svc)
	ProductsHandler := http.HandlerFunc(handler.Products)
	ProductByIDHandler := http.HandlerFunc(handler.ProductByID)
	mux.Handle("/products", middleware.Metrics(ProductsHandler, "/products"))
	mux.Handle("/products/", middleware.Metrics(ProductByIDHandler, "/products/"))
	mux.HandleFunc("/health", HealthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	return mux, nil
}
