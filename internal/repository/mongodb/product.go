package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/semaphore"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/metrics"
	"github.com/v-kuu/mini-marketplace/internal/config"
)

type ProductRepository struct {
	client *mongo.Client
	sem *semaphore.Weighted
}

func openDB(databaseURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(databaseURI))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewProductRepository (cfg *config.Config) (*ProductRepository, func(), error) {
	client, err := openDB(cfg.MONGODB_URI)
	if err != nil {
		return nil, nil, err
	}
	return &ProductRepository{
		client: client,
		sem: semaphore.NewWeighted(cfg.SEM_MAX),
	},
	func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Printf("mongodb disconnect: %v", err)
		}
	},
	nil
}

func (r *ProductRepository) List(ctx context.Context) ([]model.Product, error) {

}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {

}

func (r *ProductRepository) Create(ctx context.Context, p model.Product) error {

}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {

}

func (r *ProductRepository) Update(ctx context.Context, p model.Product) error {

}
