package mongodb

import (
	"context"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepositoryMongoDB implements the ProductRepository interface using MongoDB
type ProductRepositoryMongoDB struct {
	collection *mongo.Collection
}

// NewProductRepositoryMongoDB creates a new MongoDB product repository
func NewProductRepositoryMongoDB(db *mongo.Database) repositories.ProductRepository {
	return &ProductRepositoryMongoDB{
		collection: db.Collection("products"),
	}
}

// GetAllProducts retrieves all products with pagination from MongoDB
func (r *ProductRepositoryMongoDB) GetAllProducts(ctx context.Context, limit, offset int64) ([]*entities.Product, int64, error) {
	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Set up pagination options
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*entities.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
