package sqlite

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/config"
)

func setupTestDB(t *testing.T) (*ProductRepository, func()) {
	t.Helper()

	cfg := config.Load()
	cfg.TESTING = "YES"
	repo, cleanup, err := NewProductRepository(cfg)
	if err != nil {
		t.Fatalf("failed to create a repo: %v", err)
	}
	err = repo.Create(t.Context(), model.Product{ID: "1", Name: "Coffee", Price: 499})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	return repo, cleanup
}

func TestProductRepository_List(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	products, err := repo.List(t.Context())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("Expected 1 product, got %d", len(products))
	}

	if products[0].Name != "Coffee" {
		t.Fatalf("Unexpected product name: %s", products[0].Name)
	}
}

func TestProductRepository_List_ContextCancelled(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()
	
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := repo.List(ctx)
	if err == nil {
		t.Fatalf("Expected error due to cancelled context")
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := t.Context()

	product, err := repo.GetByID(ctx, "1")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if product.Name != "Coffee" {
		t.Fatalf("Expected Coffee, got %s", product.Name)
	}

	product, err = repo.GetByID(ctx, "3")
	if product != nil && err == nil {
		t.Fatalf("GetByID should have failed")
	}
}

func TestProductRepository_Create(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := t.Context()

	err := repo.Create(ctx, model.Product{ID: "2", Name: "Candy", Price: 499})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	product, err := repo.GetByID(ctx, "2")
	if err == nil && product.Name != "Candy" {
		t.Fatalf("Expected Candy, got %s", product.Name)
	} else if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := t.Context()

	err := repo.Create(ctx, model.Product{ID: "2", Name: "Candy", Price: 499})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Delete(ctx, "2")
	if err != nil {
		t.Fatalf("Delete failed")
	}

	product, err := repo.GetByID(ctx, "2")
	if err != nil || product != nil {
		t.Fatalf("Element was not deleted")
	}

	err = repo.Delete(ctx, "2")
	if err == nil {
		t.Fatalf("Delete should have failed")
	}
}

func TestProductRepository_Update(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := t.Context()

	err := repo.Update(ctx, model.Product{ID: "1", Name: "Tea", Price: 499})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	product, err := repo.GetByID(ctx, "1")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if product.Name != "Tea" {
		t.Fatalf("Expected Tea, got %s", product.Name)
	}

	err = repo.Update(ctx, model.Product{ID: "", Name: "", Price: 0})
	if err == nil {
		t.Fatalf("Update should have failed")
	}
}
