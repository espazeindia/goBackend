package repositories

import (
	"context"

	"espazeBackend/domain/entities"
)

// ProductRepository defines the interface for product data access operations
type ProductRepository interface {
	// GetAllProducts retrieves all products with pagination
	GetAllProducts(ctx context.Context, limit, offset int64) ([]*entities.Product, int64, error)
}
