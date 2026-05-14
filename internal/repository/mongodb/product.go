package mongodb

import (
	"context"
	"log"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/semaphore"

	"github.com/v-kuu/mini-marketplace/internal/model"
	"github.com/v-kuu/mini-marketplace/internal/service"
	"github.com/v-kuu/mini-marketplace/internal/config"
)

type ProductRepository struct {
	client *mongo.Client
	collection *mongo.Collection
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
	coll := client.Database("mini-marketplace").Collection("products")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	coll.Indexes().CreateMany(context.Background(), indexes)

	return &ProductRepository{
		client: client,
		collection: coll,
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
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	filter := bson.M{"_id": id}
	var p model.Product
	err := r.collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p model.Product) error {
	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return service.ErrProductAlreadyExists
		}
		return err
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return service.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, p model.Product) error {
	filter := bson.M{"_id": p.ID}

	fields := bson.M{}
	if p.Name != "" {
		fields["name"] = p.Name
	}
	if p.Price > 0 {
		fields["price"] = p.Price
	}
	if len(fields) == 0 {
		return nil
	}
	update := bson.M{"$set": fields}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return service.ErrProductAlreadyExists
		}
		return err
	}
	if result.MatchedCount == 0 {
		return service.ErrProductNotFound
	}
	return nil
}
