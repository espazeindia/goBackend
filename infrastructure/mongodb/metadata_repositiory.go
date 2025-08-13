package mongodb

import (
	"context"
	"fmt"
	"time"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MetadataRepositoryMongoDB implements the MetadataRepository interface using MongoDB
type MetadataRepositoryMongoDB struct {
	db *mongo.Database
}

// NewMetadataRepositoryMongoDB creates a new MongoDB metadata repository
func NewMetadataRepositoryMongoDB(db *mongo.Database) repositories.MetadataRepository {
	return &MetadataRepositoryMongoDB{
		db: db,
	}
}

// GetAllMetadata retrieves all metadata with pagination
func (r *MetadataRepositoryMongoDB) GetAllMetadata(ctx context.Context, limit, offset int64, search string) ([]*entities.GetAllMetadata, int64, error) {
	metadataColl := r.db.Collection("metadata")

	// Base match filter
	match := bson.M{}
	if search != "" && search != `""` {
		regex := bson.M{"$regex": search, "$options": "i"}
		match["$or"] = []bson.M{
			{"metadata_name": regex},
			{"subcategory_info.subcategory_name": regex}, // Will work after $lookup
		}
	}

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		// convert string IDs to ObjectIDs for lookups
		{{Key: "$addFields", Value: bson.M{
			"category_oid":    bson.M{"$toObjectId": "$metadata_category_id"},
			"subcategory_oid": bson.M{"$toObjectId": "$metadata_subcategory_id"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "category_oid",
			"foreignField": "_id",
			"as":           "category_info",
		}}},
		{{Key: "$unwind", Value: "$category_info"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "subcategories",
			"localField":   "subcategory_oid",
			"foreignField": "_id",
			"as":           "subcategory_info",
		}}},
		{{Key: "$unwind", Value: "$subcategory_info"}},
	}

	// Apply search filter if provided
	if len(match) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: match}})
	}

	// Count stage
	countStage := append(pipeline, bson.D{{Key: "$count", Value: "total"}})
	countCursor, err := metadataColl.Aggregate(ctx, countStage)
	if err != nil {
		return nil, 0, err
	}
	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	var total int64
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Pagination + sorting
	pipeline = append(pipeline,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "metadata_created_at", Value: -1}}}},
		bson.D{{Key: "$skip", Value: offset * limit}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	// Projection
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.M{
		"_id":              bson.M{"$toString": "$_id"},
		"hsn_code":         "$hsn_code",
		"name":             "$metadata_name",
		"description":      "$metadata_description",
		"image":            "$metadata_image",
		"category_id":      "$metadata_category_id",
		"category_name":    "$category_info.category_name",
		"subcategory_id":   "$metadata_subcategory_id",
		"subcategory_name": "$subcategory_info.subcategory_name",
		"mrp":              "$metadata_mrp",
		"created_at":       "$metadata_created_at",
		"updated_at":       "$metadata_updated_at",
	}}})

	// Execute aggregation
	cursor, err := metadataColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*entities.GetAllMetadata
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetAllMetadataForSeller retrieves all metadata not already in seller's inventory with optimized aggregation pipeline
func (r *MetadataRepositoryMongoDB) GetAllMetadataForSeller(ctx context.Context, limit, offset int64, search, sellerID string) ([]*entities.GetAllMetadata, int64, error) {
	metadataColl := r.db.Collection("metadata")
	inventoryColl := r.db.Collection("inventory")
	inventoryProductColl := r.db.Collection("inventory_product")

	// 1️⃣ Get seller's inventory ID
	var inventory *entities.Inventory
	err := inventoryColl.FindOne(ctx, bson.M{"seller_id": sellerID}).Decode(&inventory)
	if err == mongo.ErrNoDocuments {
		inventory = nil
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, 0, err
	}
	fmt.Print("hello", inventory, "\n")
	// 2️⃣ Collect product IDs to exclude
	var excludeIDs []primitive.ObjectID
	if inventory != nil {
		cursor, err := inventoryProductColl.Find(ctx, bson.M{"inventory_id": inventory.InventoryID})
		if err != nil {
			return nil, 0, err
		}
		defer cursor.Close(ctx)

		var products []struct {
			MetadataProductID primitive.ObjectID `bson:"metadata_product_id"`
		}
		if err := cursor.All(ctx, &products); err != nil {
			return nil, 0, err
		}
		fmt.Print(products)
		excludeIDs = make([]primitive.ObjectID, len(products))
		for i, p := range products {
			excludeIDs[i] = p.MetadataProductID
		}
	}

	// 3️⃣ Base filter for initial match (only exclusions)
	initialMatch := bson.M{}
	if len(excludeIDs) > 0 {
		initialMatch["_id"] = bson.M{"$nin": excludeIDs}
	}

	// 4️⃣ Aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: initialMatch}},
		{{Key: "$addFields", Value: bson.M{
			"category_oid":    bson.M{"$toObjectId": "$metadata_category_id"},
			"subcategory_oid": bson.M{"$toObjectId": "$metadata_subcategory_id"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "category_oid",
			"foreignField": "_id",
			"as":           "category_info",
		}}},
		{{Key: "$unwind", Value: "$category_info"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "subcategories",
			"localField":   "subcategory_oid",
			"foreignField": "_id",
			"as":           "subcategory_info",
		}}},
		{{Key: "$unwind", Value: "$subcategory_info"}},
	}

	// 5️⃣ Add search filter for both metadata name and subcategory name after lookups
	if search != "" && search != `""` {
		searchMatch := bson.M{
			"$or": []bson.M{
				{"metadata_name": bson.M{"$regex": search, "$options": "i"}},
				{"subcategory_info.subcategory_name": bson.M{"$regex": search, "$options": "i"}},
			},
		}
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: searchMatch}})
	}

	// Common sort
	sortStage := bson.D{
		{Key: "metadata_created_at", Value: -1},
		{Key: "_id", Value: 1}, // stable sort
	}

	// 6️⃣ Count (clone pipeline to avoid rebuilding)
	countPipeline := append(append(mongo.Pipeline{}, pipeline...), bson.D{{Key: "$count", Value: "total"}})
	countCursor, err := metadataColl.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	var total int64
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// 7️⃣ Data fetch
	dataPipeline := append(pipeline,
		bson.D{{Key: "$sort", Value: sortStage}},
		bson.D{{Key: "$skip", Value: offset * limit}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":              bson.M{"$toString": "$_id"},
			"hsn_code":         "$hsn_code",
			"name":             "$metadata_name",
			"description":      "$metadata_description",
			"image":            "$metadata_image",
			"category_id":      "$metadata_category_id",
			"category_name":    "$category_info.category_name",
			"subcategory_id":   "$metadata_subcategory_id",
			"subcategory_name": "$subcategory_info.subcategory_name",
			"mrp":              "$metadata_mrp",
			"created_at":       "$metadata_created_at",
			"updated_at":       "$metadata_updated_at",
		}}},
	)

	cursor, err := metadataColl.Aggregate(ctx, dataPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*entities.GetAllMetadata
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetMetadataByID retrieves a metadata by ID
func (r *MetadataRepositoryMongoDB) GetMetadataByID(ctx context.Context, id string) (*entities.MetadataResponse, error) {
	reviews := r.db.Collection("reviews")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var metadata entities.Metadata

	err = r.db.Collection("metadata").FindOne(ctx, filter).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	review := entities.Review{}
	err = reviews.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}

	categoryCollection := r.db.Collection("categories")

	subCategoryCollection := r.db.Collection("subcategories")

	var category *entities.Category

	var subcategory *entities.Subcategory

	categoryObjectId, err := primitive.ObjectIDFromHex(metadata.MetadataCategoryID)
	if err != nil {
		return nil, err
	}

	categoryFilter := bson.M{"_id": categoryObjectId}

	err = categoryCollection.FindOne(ctx, categoryFilter).Decode(&category)
	if err != nil {
		return nil, err
	}

	subcategoryObjectId, err := primitive.ObjectIDFromHex(metadata.MetadataSubcategoryID)
	if err != nil {
		return nil, err
	}

	subcategoryFilter := bson.M{"_id": subcategoryObjectId}

	err = subCategoryCollection.FindOne(ctx, subcategoryFilter).Decode(&subcategory)
	if err != nil {
		return nil, err
	}

	metadataResponse := &entities.MetadataResponse{
		ID:              metadata.MetadataProductID,
		Name:            metadata.MetadataName,
		Description:     metadata.MetadataDescription,
		Image:           metadata.MetadataImage,
		CategoryID:      metadata.MetadataCategoryID,
		SubcategoryID:   metadata.MetadataSubcategoryID,
		CategoryName:    category.CategoryName,
		SubCategoryName: subcategory.SubcategoryName,
		MRP:             metadata.MetadataMRP,
		CreatedAt:       metadata.MetadataCreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       metadata.MetadataUpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		TotalStars:      review.TotalStars,
		TotalReviews:    review.TotalReviews,
		HsnCode:         metadata.MetadataHSNCode,
	}
	return metadataResponse, nil
}

// CreateMetadata creates a new metadata
func (r *MetadataRepositoryMongoDB) CreateMetadata(ctx context.Context, metadata *entities.Metadata) (*entities.MetadataApiResponse, error) {
	collection := r.db.Collection("metadata")

	filter := bson.M{"hsn_code": metadata.MetadataHSNCode}
	var Metadata *entities.Metadata
	err := collection.FindOne(ctx, filter).Decode(&Metadata)
	if err == nil {
		return &entities.MetadataApiResponse{Success: false, Message: "Metadata for this HSN Code already exists", Error: "Metadata Already Exists"}, err
	}
	if err != mongo.ErrNoDocuments {
		return &entities.MetadataApiResponse{Success: false, Message: "Internal Server Error", Error: "DataBase Error"}, err
	}

	result, err := r.db.Collection("metadata").InsertOne(ctx, metadata)
	if err != nil {
		return &entities.MetadataApiResponse{Success: false, Message: "Internal Server Error", Error: "DataBase Error"}, err
	}
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MetadataApiResponse{Success: false, Message: "Error Creating Metadata", Error: "ObjectId fetch error"}, nil
	}
	stringID := objectID.Hex()
	return &entities.MetadataApiResponse{Success: true, Message: "Metadata Created Successfully", Id: stringID}, nil
}

// UpdateMetadata updates an existing metadata
func (r *MetadataRepositoryMongoDB) UpdateMetadata(ctx context.Context, id string, metadata *entities.Metadata) (*entities.MetadataApiResponse, error) {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &entities.MetadataApiResponse{
			Error:   "Error Transforming ObjectIDFromHex",
			Message: "Error in transforming metadata id to object id",
			Success: false,
		}, err
	}

	// Build update document
	updateDoc := bson.M{}
	if metadata.MetadataName != "" {
		updateDoc["metadata_name"] = metadata.MetadataName
	}
	if metadata.MetadataDescription != "" {
		updateDoc["metadata_description"] = metadata.MetadataDescription
	}
	if metadata.MetadataImage != "" {
		updateDoc["metadata_image"] = metadata.MetadataImage
	}
	if metadata.MetadataCategoryID != "" {
		updateDoc["metadata_category_id"] = metadata.MetadataCategoryID
	}
	if metadata.MetadataSubcategoryID != "" {
		updateDoc["metadata_subcategory_id"] = metadata.MetadataSubcategoryID
	}
	if metadata.MetadataHSNCode != "" {
		updateDoc["hsn_code"] = metadata.MetadataHSNCode
	}
	if metadata.MetadataMRP > 0 {
		updateDoc["metadata_mrp"] = metadata.MetadataMRP
	}
	updateDoc["metadata_updated_at"] = time.Now()

	// Execute update
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateDoc}

	result, err := r.db.Collection("metadata").UpdateOne(ctx, filter, update)
	if err != nil {
		return &entities.MetadataApiResponse{
			Error:   "Database Error",
			Message: "Error in database while updating",
			Success: false,
		}, err
	}

	if result.MatchedCount == 0 {
		return &entities.MetadataApiResponse{
			Error:   "No matching document in db",
			Message: "Metadata not found in db",
			Success: false,
		}, err
	}

	return &entities.MetadataApiResponse{
		Message: "Metadata updated successfully",
		Success: true,
	}, nil
}

// DeleteMetadata deletes a metadata by ID
func (r *MetadataRepositoryMongoDB) DeleteMetadata(ctx context.Context, id string) (*entities.MetadataApiResponse, error) {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &entities.MetadataApiResponse{
			Error:   "Error Transforming ObjectIDFromHex",
			Message: "Error in transforming metadata id to object id",
			Success: false,
		}, err
	}

	filter := bson.M{"_id": objectID}
	result, err := r.db.Collection("metadata").DeleteOne(ctx, filter)
	if err != nil {
		return &entities.MetadataApiResponse{
			Error:   "Database Error",
			Message: "Error in database while deleting",
			Success: false,
		}, err
	}

	if result.DeletedCount == 0 {
		return &entities.MetadataApiResponse{
			Error:   "No matching document in db",
			Message: "Metadata not found in db",
			Success: false,
		}, err
	}

	return &entities.MetadataApiResponse{
		Message: "Metadata deleted successfully",
		Success: true,
	}, nil
}

func (r *MetadataRepositoryMongoDB) AddReview(ctx context.Context, req *entities.AddReviewRequest) error {
	collection := r.db.Collection("reviews")
	filter := bson.M{"metadata_product_id": req.MetadataProductID}

	review := entities.Review{}
	err := collection.FindOne(ctx, filter).Decode(&review)
	if err != mongo.ErrNoDocuments {
		return err
	}

	if err != nil {
		review = entities.Review{
			MetadataProductID: req.MetadataProductID,
			TotalStars:        0,
			TotalReviews:      0,
		}
	}
	review.TotalStars += req.Rating
	review.TotalReviews += 1

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": review})
	if err != nil {
		return err
	}

	return nil
}

func (r *MetadataRepositoryMongoDB) CreateReview(ctx context.Context, id string) (*entities.MetadataApiResponse, error) {
	collection := r.db.Collection("reviews")
	review := entities.Review{
		MetadataProductID: id,
		TotalStars:        0,
		TotalReviews:      0,
	}
	_, err := collection.InsertOne(ctx, review)
	if err != nil {
		return &entities.MetadataApiResponse{Success: false, Message: "Internal Server Error", Error: "DataBase Error"}, err
	}
	return &entities.MetadataApiResponse{Success: true, Message: "Metadata Created Successfully"}, nil
}
