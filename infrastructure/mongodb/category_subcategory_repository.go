package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CategorySubcategoryRepositoryMongoDB implements the CategorySubcategoryRepository interface using MongoDB
type CategorySubcategoryRepositoryMongoDB struct {
	db *mongo.Database
}

// NewCategorySubcategoryRepositoryMongoDB creates a new MongoDB category subcategory repository
func NewCategorySubcategoryRepositoryMongoDB(db *mongo.Database) repositories.CategorySubcategoryRepository {
	return &CategorySubcategoryRepositoryMongoDB{db: db}
}

// Category operations
func (r *CategorySubcategoryRepositoryMongoDB) GetAllCategories(ctx context.Context) ([]*entities.Category, error) {
	collection := r.db.Collection("categories")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*entities.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) GetCategoryById(ctx context.Context, categoryID string) (*entities.Category, error) {
	collection := r.db.Collection("categories")

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, err
	}

	var category entities.Category
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) CreateCategory(ctx context.Context, category *entities.Category) error {
	collection := r.db.Collection("categories")

	// Generate new ObjectID
	objectID := primitive.NewObjectID()
	category.CategoryID = objectID.Hex()

	// Set timestamps
	now := primitive.NewDateTimeFromTime(category.CategoryCreatedAt)
	category.CategoryCreatedAt = now.Time()
	category.CategoryUpdatedAt = now.Time()

	_, err := collection.InsertOne(ctx, category)
	return err
}

func (r *CategorySubcategoryRepositoryMongoDB) UpdateCategory(ctx context.Context, category *entities.Category) error {
	collection := r.db.Collection("categories")

	objectID, err := primitive.ObjectIDFromHex(category.CategoryID)
	if err != nil {
		return err
	}

	// Update timestamp
	category.CategoryUpdatedAt = primitive.NewDateTimeFromTime(category.CategoryUpdatedAt).Time()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": category},
	)
	return err
}

func (r *CategorySubcategoryRepositoryMongoDB) DeleteCategory(ctx context.Context, categoryID string) error {
	collection := r.db.Collection("categories")

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Subcategory operations
func (r *CategorySubcategoryRepositoryMongoDB) GetAllSubcategories(ctx context.Context) ([]*entities.Subcategory, error) {
	collection := r.db.Collection("subcategories")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subcategories []*entities.Subcategory
	if err := cursor.All(ctx, &subcategories); err != nil {
		return nil, err
	}
	return subcategories, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) GetSubcategoryById(ctx context.Context, subcategoryID string) (*entities.Subcategory, error) {
	collection := r.db.Collection("subcategories")

	objectID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return nil, err
	}

	var subcategory entities.Subcategory
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&subcategory)
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) GetSubcategoriesByCategoryId(ctx context.Context, categoryID string) ([]*entities.Subcategory, error) {
	collection := r.db.Collection("subcategories")

	cursor, err := collection.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subcategories []*entities.Subcategory
	if err := cursor.All(ctx, &subcategories); err != nil {
		return nil, err
	}
	return subcategories, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) CreateSubcategory(ctx context.Context, subcategory *entities.Subcategory) error {
	collection := r.db.Collection("subcategories")

	// Generate new ObjectID
	objectID := primitive.NewObjectID()
	subcategory.SubcategoryID = objectID.Hex()

	// Set timestamps
	now := primitive.NewDateTimeFromTime(subcategory.SubcategoryCreatedAt)
	subcategory.SubcategoryCreatedAt = now.Time()
	subcategory.SubcategoryUpdatedAt = now.Time()

	_, err := collection.InsertOne(ctx, subcategory)
	return err
}

func (r *CategorySubcategoryRepositoryMongoDB) UpdateSubcategory(ctx context.Context, subcategory *entities.Subcategory) error {
	collection := r.db.Collection("subcategories")

	objectID, err := primitive.ObjectIDFromHex(subcategory.SubcategoryID)
	if err != nil {
		return err
	}

	// Update timestamp
	subcategory.SubcategoryUpdatedAt = primitive.NewDateTimeFromTime(subcategory.SubcategoryUpdatedAt).Time()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": subcategory},
	)
	return err
}

func (r *CategorySubcategoryRepositoryMongoDB) DeleteSubcategory(ctx context.Context, subcategoryID string) error {
	collection := r.db.Collection("subcategories")

	objectID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
