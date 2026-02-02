package main

import (
	"database/sql"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/repository/sqlite"
	"github.com/v-kuu/mini-marketplace/internal/metrics"
	"github.com/v-kuu/mini-marketplace/internal/http/api"
	"github.com/v-kuu/mini-marketplace/internal/http/middleware"
)

func main() {
	metrics.Register()
	db, err := sql.Open("sqlite3", "file:products.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", api.HealthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	repo := sqlite.NewProductRepository(db)
	svc := service.NewProductService(repo)
	handler := api.NewProductHandler(svc)
	mux.Handle("/products", middleware.Metrics(http.HandlerFunc(handler.Products), "/products"))
	mux.Handle("/products/", middleware.Metrics(http.HandlerFunc(handler.ProductByID), "/products/"))

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
