package sqlite

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
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
