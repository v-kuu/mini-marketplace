package api

import (
	"database/sql"
	"net/http"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/repository/sqlite"
	"github.com/v-kuu/mini-marketplace/internal/metrics"
	"github.com/v-kuu/mini-marketplace/internal/http/middleware"
)

func AddRoutes() (*http.ServeMux, error) {
	metrics.Register()

	db, err := sql.Open("sqlite3", "file:products.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	repo := sqlite.NewProductRepository(db)
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
