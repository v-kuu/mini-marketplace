package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/v-kuu/mini-marketplace/internal/api"
	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/repository"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", api.HealthHandler)

	repo := &repository.ProductRepository{}
	svc := service.NewProductService(repo)
	handler := api.NewProductHandler(svc)

	mux.HandleFunc("/products", handler.List)

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
