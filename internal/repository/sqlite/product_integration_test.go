package sqlite

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/v-kuu/mini-marketplace/internal/model"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}

	schema := `
	CREATE TABLE products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price INTEGER NOT NULL
	);
	`
	
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func TestProductRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	ctx := context.Background()

	_, err := db.Exec(
		`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
		"1", "Coffee", 499,
	)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

	products, err := repo.List(ctx)
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
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.List(ctx)
	if err == nil {
		t.Fatalf("Expected error due to cancelled context")
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	ctx := context.Background()

	_, err := db.Exec(
		`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
		"1", "Coffee", 499,
	)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

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
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	ctx := context.Background()

	err := repo.Create(ctx, model.Product{ID: "1", Name: "Coffee", Price: 499})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	product, err := repo.GetByID(ctx, "1")
	if product.Name != "Coffee" {
		t.Fatalf("Expected Coffee, got %s", product.Name)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	ctx := context.Background()

	_, err := db.Exec(
		`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
		"1", "Coffee", 499,
	)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}
	_, err = db.Exec(
		`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
		"2", "Tea", 499,
	)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
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
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	ctx := context.Background()

	_, err := db.Exec(
		`INSERT INTO products (id, name, price) VALUES (?, ?, ?)`,
		"1", "Coffee", 499,
	)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

	err = repo.Update(ctx, model.Product{ID: "1", Name: "Tea", Price: 499})
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
