package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (r *CategorySubcategoryRepositoryMongoDB) GetCategories(ctx context.Context, limit, offset int64, search *string) ([]*entities.Category, int64, error) {
	collection := r.db.Collection("categories")
	filter := bson.M{}
	if *search != "" {
		filter = bson.M{
			"category_name": bson.M{"$regex": search, "$options": "i"},
		}
	}
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset * limit).
		SetSort(bson.D{{Key: "category_created_at", Value: -1}}) // Sort by creation date descending

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []*entities.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) GetAllCategories(ctx context.Context) ([]*entities.Category, error) {
	collection := r.db.Collection("categories")
	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
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

// Subcategory operations
func (r *CategorySubcategoryRepositoryMongoDB) GetAllSubcategories(ctx context.Context, categoryId, search string) ([]*entities.Subcategory, error) {
	collection := r.db.Collection("subcategories")
	filter := bson.M{"category_id": categoryId,
		"subcategory_name": bson.M{"$regex": search, "$options": "i"},
	}

	cursor, err := collection.Find(ctx, filter)
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

func (r *CategorySubcategoryRepositoryMongoDB) GetSubcategories(ctx context.Context, limit, offset int64, search *string) ([]*entities.Subcategory, int64, error) {
	collection := r.db.Collection("subcategories")
	filter := bson.M{}
	if *search != "" {
		filter = bson.M{
			"subcategory_name": bson.M{"$regex": search, "$options": "i"},
		}
	}
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset * limit).
		SetSort(bson.D{{Key: "subcategory_created_at", Value: -1}}) // Sort by creation date descending

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var subcategories []*entities.Subcategory
	if err := cursor.All(ctx, &subcategories); err != nil {
		return nil, 0, err
	}

	return subcategories, total, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) GetSubcategoryByCategoryId(ctx context.Context, categoryId *string, limit, offset int64, search *string) ([]*entities.Subcategory, int64, error) {
	collection := r.db.Collection("subcategories")
	filter := bson.M{"category_id": categoryId}
	if *search != "" {
		filter = bson.M{
			"subcategory_name": bson.M{"$regex": search, "$options": "i"},
			"category_id":      categoryId,
		}
	}
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset * limit).
		SetSort(bson.D{{Key: "subcategory_created_at", Value: -1}}) // Sort by creation date descending

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var subcategories []*entities.Subcategory
	if err := cursor.All(ctx, &subcategories); err != nil {
		return nil, 0, err
	}

	return subcategories, total, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) CreateCategory(ctx context.Context, category *entities.Category) (*entities.MessageResponse, error) {
	collection := r.db.Collection("categories")
	filter := bson.M{"category_name": category.CategoryName}
	var existingCategory *entities.Category
	err := collection.FindOne(ctx, filter).Decode(&existingCategory)
	if err == nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "This category already exists",
			Error:   "Category already exists",
		}, err
	} else if err != mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating category",
			Error:   "Database error",
		}, err
	}

	result, err := collection.InsertOne(ctx, category)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating category",
			Error:   "Database error",
		}, err
	}
	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating category",
			Error:   "InsertedId error",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Category Created",
	}, nil

}
func (r *CategorySubcategoryRepositoryMongoDB) CreateSubcategory(ctx context.Context, subcategory *entities.Subcategory) (*entities.MessageResponse, error) {
	collection := r.db.Collection("subcategories")

	filter := bson.M{"subcategory_name": subcategory.SubcategoryName}
	var existingSubCategory *entities.Category
	err := collection.FindOne(ctx, filter).Decode(&existingSubCategory)
	if err == nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "This sub-category already exists",
			Error:   "Sub-Category already exists",
		}, err
	} else if err != mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating sub-category",
			Error:   "Database error",
		}, err
	}

	result, err := collection.InsertOne(ctx, subcategory)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating sub-category",
			Error:   "Database error",
		}, err
	}
	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating sub-category",
			Error:   "InsertedId error",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Sub-Category Created",
	}, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) UpdateCategory(ctx context.Context, categoryId *string, request *entities.UpdateCategoryRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("categories")

	objectID, err := primitive.ObjectIDFromHex(*categoryId)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating object from Category Id",
			Error:   "ObjectId from Hex error",
		}, err
	}
	updateDocs := bson.M{}
	if request.CategoryName != "" {
		updateDocs["category_name"] = request.CategoryName
	}
	if request.CategoryImage != "" {
		updateDocs["category_image"] = request.CategoryImage
	}
	updateDocs["category_updated_at"] = time.Now()

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateDocs},
	)
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "No matching document found",
			Error:   "No Document Matched",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Category Updated Successfully",
	}, nil

}

func (r *CategorySubcategoryRepositoryMongoDB) UpdateSubcategory(ctx context.Context, subcategoryID string, request *entities.UpdateSubcategoryRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("subcategories")

	objectID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating object from Sub-Category Id",
			Error:   "ObjectId from Hex error",
		}, err
	}
	updateDocs := bson.M{}
	if request.SubcategoryName != "" {
		updateDocs["subcategory_name"] = request.SubcategoryName
	}
	if request.SubcategoryImage != "" {
		updateDocs["category_image"] = request.SubcategoryImage
	}
	if request.CategoryID != "" {
		updateDocs["category_id"] = request.CategoryID
	}
	updateDocs["subcategory_updated_at"] = time.Now()

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateDocs},
	)
	if result.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "No matching document found",
			Error:   "No Document Matched",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Category Updated Successfully",
	}, nil

}

// func (r *CategorySubcategoryRepositoryMongoDB) GetCategoryById(ctx context.Context, categoryID string) (*entities.Category, error) {
// 	collection := r.db.Collection("categories")

// 	objectID, err := primitive.ObjectIDFromHex(categoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var category entities.Category
// 	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &category, nil
// }

func (r *CategorySubcategoryRepositoryMongoDB) DeleteCategory(ctx context.Context, categoryID string) (*entities.MessageResponse, error) {
	categoryCollection := r.db.Collection("categories")
	SubCollection := r.db.Collection("subcategories")

	filter := bson.M{"category_id": categoryID}

	var subcategory *entities.Subcategory

	err := SubCollection.FindOne(ctx, filter).Decode(&subcategory)
	if err == nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "This category contains some subcategories so DELETION is not possible",
			Error:   "Subcategoy exists for category",
		}, nil
	}

	if err != mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database error",
			Error:   "error fetching suubcategories",
		}, err
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating object from category Id",
			Error:   "object ID from hex error",
		}, nil
	}

	result, err := categoryCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database error",
			Error:   "error deleting category",
		}, nil
	}

	if result.DeletedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "No matching entry found",
			Error:   "no matching document in db ",
		}, nil
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Category deleted successfully",
	}, nil

}

func (r *CategorySubcategoryRepositoryMongoDB) DeleteSubcategory(ctx context.Context, subcategoryID string) (*entities.MessageResponse, error) {
	metadataCollection := r.db.Collection("metadata")
	SubCollection := r.db.Collection("subcategories")

	filter := bson.M{"metadata_subcategory_id": subcategoryID}

	var metadata *entities.Metadata

	err := metadataCollection.FindOne(ctx, filter).Decode(&metadata)
	if err == nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "This sub-category contains some metadata so DELETION is not possible",
			Error:   "metadata exists for sub-category",
		}, nil
	}

	if err != mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database error",
			Error:   "error fetching metadata",
		}, err
	}

	objectID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating object from sub-category Id",
			Error:   "object ID from hex error",
		}, nil
	}

	result, err := SubCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database error",
			Error:   "error deleting category",
		}, nil
	}

	if result.DeletedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "No matching entry found",
			Error:   "no matching document in db ",
		}, nil
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Sub-Category deleted successfully",
	}, nil

}

func (r *CategorySubcategoryRepositoryMongoDB) GetCategorySubCategoryForSpecificStore(
	ctx context.Context,
	storeID string,
) ([]*entities.CategoryWithSubcategoriesResponse, error) {

	collection := r.db.Collection("categories")
	subCategoryCollection := r.db.Collection("subcategories")

	// Get products for this store
	productRepo := ProductRepositoryMongoDB{db: r.db}
	products, err := productRepo.GetProductsForSpecificStore(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Use map for fast lookup: categoryID -> list of subcategoryIDs
	categoryMap := make(map[string][]string)

	for _, product := range products {
		cs := categoryMap[product.MetadataCategoryId]

		// Add subcategory if not already present
		found := false
		for _, sub := range cs {
			if sub == product.MetadataSubcategoryId {
				found = true
				break
			}
		}
		if !found {
			cs = append(cs, product.MetadataSubcategoryId)
		}
		categoryMap[product.MetadataCategoryId] = cs
	}

	// Convert map back to slice
	var categorySubcategoryList []*entities.CategoryWithSubcategoriesResponse

	for categoryID, subcategoryIDs := range categoryMap {
		catObjID, err := primitive.ObjectIDFromHex(categoryID)
		if err != nil {
			return nil, err
		}

		var category entities.Category
		err = collection.FindOne(ctx, bson.M{"_id": catObjID}).Decode(&category)
		if err != nil {
			return nil, err
		}

		var allSubCategories []*entities.Subcategory
		for _, subID := range subcategoryIDs {
			subObjID, err := primitive.ObjectIDFromHex(subID)
			if err != nil {
				return nil, err
			}
			var subcategory entities.Subcategory
			err = subCategoryCollection.FindOne(ctx, bson.M{"_id": subObjID}).Decode(&subcategory)
			if err != nil {
				return nil, err
			}
			allSubCategories = append(allSubCategories, &subcategory)
		}

		response := &entities.CategoryWithSubcategoriesResponse{
			Category:      &category,
			Subcategories: allSubCategories,
		}
		categorySubcategoryList = append(categorySubcategoryList, response)
	}

	return categorySubcategoryList, nil
}

func (r *CategorySubcategoryRepositoryMongoDB) GetCategorySubCategoryForAllStoresInWarehouse(ctx context.Context, warehouseId string) ([]*entities.CategoryWithSubcategoriesResponse, error) {
	collection := r.db.Collection("stores")

	// Fix: Query stores by warehouse_id field, not _id
	cursor, err := collection.Find(ctx, bson.M{"warehouse_id": warehouseId})
	if err != nil {
		return nil, err
	}
	var stores []*entities.Store
	err = cursor.All(ctx, &stores)
	if err != nil {
		return nil, err
	}

	// Use map to merge categories and subcategories from all stores
	categoryMap := make(map[string]*entities.CategoryWithSubcategoriesResponse)

	for _, store := range stores {
		storeResults, err := r.GetCategorySubCategoryForSpecificStore(ctx, store.StoreID)
		if err != nil {
			return nil, err
		}

		// Merge results from this store into the main map
		for _, categoryResult := range storeResults {
			categoryID := categoryResult.Category.CategoryID

			if existingCategory, exists := categoryMap[categoryID]; exists {
				// Category already exists, merge subcategories
				existingSubcategoryMap := make(map[string]bool)

				// Add existing subcategories to map
				for _, sub := range existingCategory.Subcategories {
					existingSubcategoryMap[sub.SubcategoryID] = true
				}

				// Add new subcategories that don't already exist
				for _, newSub := range categoryResult.Subcategories {
					if !existingSubcategoryMap[newSub.SubcategoryID] {
						existingCategory.Subcategories = append(existingCategory.Subcategories, newSub)
					}
				}
			} else {
				// New category, add it to the map
				categoryMap[categoryID] = categoryResult
			}
		}
	}

	// Convert map back to slice
	var result []*entities.CategoryWithSubcategoriesResponse
	for _, categoryResult := range categoryMap {
		result = append(result, categoryResult)
	}

	return result, nil
}

// func (r *CategorySubcategoryRepositoryMongoDB) GetSubcategoryById(ctx context.Context, subcategoryID string) (*entities.Subcategory, error) {
// 	collection := r.db.Collection("subcategories")

// 	objectID, err := primitive.ObjectIDFromHex(subcategoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var subcategory entities.Subcategory
// 	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&subcategory)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &subcategory, nil
// }

// func (r *CategorySubcategoryRepositoryMongoDB) GetSubcategoriesByCategoryId(ctx context.Context, categoryID string) ([]*entities.Subcategory, error) {
// 	collection := r.db.Collection("subcategories")

// 	cursor, err := collection.Find(ctx, bson.M{"category_id": categoryID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var subcategories []*entities.Subcategory
// 	if err := cursor.All(ctx, &subcategories); err != nil {
// 		return nil, err
// 	}
// 	return subcategories, nil
// }

// func (r *CategorySubcategoryRepositoryMongoDB) UpdateSubcategory(ctx context.Context, subcategory *entities.Subcategory) error {
// 	collection := r.db.Collection("subcategories")

// 	objectID, err := primitive.ObjectIDFromHex(subcategory.SubcategoryID)
// 	if err != nil {
// 		return err
// 	}

// 	// Update timestamp
// 	subcategory.SubcategoryUpdatedAt = primitive.NewDateTimeFromTime(subcategory.SubcategoryUpdatedAt).Time()

// 	_, err = collection.UpdateOne(
// 		ctx,
// 		bson.M{"_id": objectID},
// 		bson.M{"$set": subcategory},
// 	)
// 	return err
// }

// func (r *CategorySubcategoryRepositoryMongoDB) DeleteSubcategory(ctx context.Context, subcategoryID string) error {
// 	collection := r.db.Collection("subcategories")

// 	objectID, err := primitive.ObjectIDFromHex(subcategoryID)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	return err
// }
