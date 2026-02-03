package api

import (
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/repository/sqlite"
	"github.com/v-kuu/mini-marketplace/internal/metrics"
	"github.com/v-kuu/mini-marketplace/internal/http/middleware"
)

func AddRoutes() (*http.ServeMux, error) {
	metrics.Register()

	db, err := sqlite.OpenDB("file:products.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	var maxConcurrent int64
	SEM_MAX, ok := os.LookupEnv("SEM_MAX")
	if !ok {
		maxConcurrent = 100
	} else {
		val, err := strconv.Atoi(SEM_MAX)
		if err != nil {
			return nil, err
		}
		maxConcurrent = int64(val);
	}
	repo := sqlite.NewProductRepository(db, maxConcurrent)
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
