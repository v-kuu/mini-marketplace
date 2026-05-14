package mongodb

import (
	"testing"
	"errors"
	"time"
	"context"

	"github.com/v-kuu/mini-marketplace/internal/config"
	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testProducts = []model.Product{
	{ID: "1", Name: "Chair", Price: 49},
	{ID: "2", Name: "Table", Price: 199},
	{ID: "3", Name: "Lamp", Price: 29},
}

func setupTestDB(t *testing.T) (*ProductRepository, func()) {
	t.Helper()
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:6",
		testcontainers.WithWaitStrategy(wait.ForLog("Waiting for connections").
			WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start mongodb container: %v", err)
	}

	cfg := config.Load()
	cfg.MONGODB_URI, err = container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	repo, cleanup, err := NewProductRepository(cfg)
	if err != nil {
		t.Fatalf("failed to create a repo: %v", err)
	}
	for _, p := range testProducts {
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("failed to seed product: %v", err)
		}
	}

	return repo, func() {
		cleanup()
		if err := container.Terminate(t.Context()); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}
}

func TestList(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	products, err := repo.List(t.Context())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != len(testProducts) {
		t.Fatalf("expected %d products, got %d", len(testProducts), len(products))
	}
}

func TestGetByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()
	
	product, err := repo.GetByID(t.Context(), "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.Name != testProducts[0].Name {
		t.Fatalf("expected %s, got %s", testProducts[0].Name, product.Name)
	}

	product, err = repo.GetByID(t.Context(), "5")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product != nil {
		t.Fatalf("expected nil product")
	}
}

func TestCreate(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	p := model.Product{ID: "4", Name: "Plate", Price: 49}
	err := repo.Create(t.Context(), p)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreate_Duplicate(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	p := model.Product{ID: "4", Name: "Chair", Price: 39}
	err := repo.Create(t.Context(), p)
	if !errors.Is(err, service.ErrProductAlreadyExists) {
		t.Fatalf("expected ErrProductAlreadyExists, got %v", err)
	}

	p = model.Product{ID: "1", Name: "Chair", Price: 39}
	err = repo.Create(t.Context(), p)
	if !errors.Is(err, service.ErrProductAlreadyExists) {
		t.Fatalf("expected ErrProductAlreadyExists, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	err := repo.Delete(t.Context(), "5")
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}

	err = repo.Delete(t.Context(), "1")
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}
	product, err := repo.GetByID(t.Context(), "1")
	if product != nil || err != nil {
		t.Fatalf("expected nil nil, got %+v, %v", product, err)
	}
}

func TestUpdate(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	updated := model.Product{ID: "1", Name: "Updated Chair", Price: 99}
	err := repo.Update(t.Context(), updated)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	product, _ := repo.GetByID(t.Context(), "1")
	if *product != updated {
		t.Fatalf("expected %+v, got %+v", updated, product)
	}

	updated = model.Product{ID: "2", Name: "Updated Table"}
	err = repo.Update(t.Context(), updated)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	product, _ = repo.GetByID(t.Context(), "2")
	if product.Name != updated.Name {
		t.Fatalf("expected %s, got %s", updated.Name, product.Name)
	}
	if product.Price != testProducts[1].Price {
		t.Fatalf("price shouldn't have changed, is now %d", product.Price)
	}

	updated = model.Product{ID: "3", Price: 999}
	err = repo.Update(t.Context(), updated)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	product, _ = repo.GetByID(t.Context(), "3")
	if product.Price != updated.Price {
		t.Fatalf("expected %d, got %d", updated.Price, product.Price)
	}
	if product.Name != testProducts[2].Name {
		t.Fatalf("name shouldn't have changed, is now %s", product.Name)
	}

	updated = model.Product{ID: "6", Name: "New product", Price: 999}
	err = repo.Update(t.Context(), updated)
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}

	updated = model.Product{ID: "1", Name: testProducts[2].Name, Price: 499}
	err = repo.Update(t.Context(), updated)
	if !errors.Is(err, service.ErrProductAlreadyExists) {
		t.Fatalf("expected ErrProductAlreadyExists, got %v", err)
	}

	updated = model.Product{ID: "1"}
	err = repo.Update(t.Context(), updated)
	if err != nil {
		t.Fatalf("expected no error on empty update, got %v", err)
	}
}
