package repositories

import (
	"context"

	"espazeBackend/domain/entities"
)

// CategorySubcategoryRepository defines the interface for category and subcategory data access operations
type CategorySubcategoryRepository interface {
	// Category operations
	GetCategories(ctx context.Context, limit, offset int64, search *string) ([]*entities.Category, int64, error)
	GetAllCategories(ctx context.Context) ([]*entities.Category, error)
	GetSubcategories(ctx context.Context, limit, offset int64, search *string) ([]*entities.Subcategory, int64, error)
	GetAllSubcategories(ctx context.Context, categoryId string) ([]*entities.Subcategory, error)
	CreateCategory(ctx context.Context, category *entities.Category) (*entities.MessageResponse, error)
	CreateSubcategory(ctx context.Context, subcategory *entities.Subcategory) (*entities.MessageResponse, error)
	GetSubcategoryByCategoryId(ctx context.Context, categoryId *string, limit, offset int64, search *string) ([]*entities.Subcategory, int64, error)
	UpdateCategory(ctx context.Context, categoryId *string, request *entities.UpdateCategoryRequest) (*entities.MessageResponse, error)
	UpdateSubcategory(ctx context.Context, subcategoryID string, request *entities.UpdateSubcategoryRequest) (*entities.MessageResponse, error)
	DeleteCategory(ctx context.Context, categoryID string) (*entities.MessageResponse, error)
	DeleteSubcategory(ctx context.Context, subcategoryID string) (*entities.MessageResponse, error)
	GetCategorySubCategoryForAllStoresInWarehouse(ctx context.Context, warehouseId string) ([]*entities.CategoryWithSubcategoriesResponse, error)
	GetCategorySubCategoryForSpecificStore(ctx context.Context, storeID string) ([]*entities.CategoryWithSubcategoriesResponse, error)

	// GetCategoryById(ctx context.Context, categoryID string) (*entities.Category, error)
	// // Subcategory operations
	// GetSubcategoryById(ctx context.Context, subcategoryID string) (*entities.Subcategory, error)
	// GetSubcategoriesByCategoryId(ctx context.Context, categoryID string) ([]*entities.Subcategory, error)
	// DeleteSubcategory(ctx context.Context, subcategoryID string) error
}
