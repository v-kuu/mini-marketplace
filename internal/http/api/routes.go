package api

import (
	"net/http"
	"log"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/repository/sqlite"
	"github.com/v-kuu/mini-marketplace/internal/repository/mongodb"
	"github.com/v-kuu/mini-marketplace/internal/metrics"
	"github.com/v-kuu/mini-marketplace/internal/config"
	"github.com/v-kuu/mini-marketplace/internal/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes() (*http.ServeMux, func(), error) {
	metrics.Register()

	cfg := config.Load()
	var repo service.ProductRepository
	var cleanup func()
	var err error

	if cfg.MONGODB_URI == "" {
		log.Println("No mongodb URI set, defaulting to sqlite")
		repo, cleanup, err = sqlite.NewProductRepository(cfg)
	} else {
		log.Println("Connecting to mongodb")
		repo, cleanup, err = mongodb.NewProductRepository(cfg)
	}
	if err != nil {
		return nil, nil, err
	}
	defer cleanup()
	svc := service.NewProductService(repo)
	handler := NewProductHandler(svc, cfg)
	ProductsHandler := http.HandlerFunc(handler.Products)
	ProductByIDHandler := http.HandlerFunc(handler.ProductByID)

	mux := http.NewServeMux()
	mux.Handle("/products", middleware.Metrics(ProductsHandler, "/products"))
	mux.Handle("/products/", middleware.Metrics(ProductByIDHandler, "/products/"))
	mux.HandleFunc("/health", HealthHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	return mux, cleanup, nil
}
