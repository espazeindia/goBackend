package mongodb

import (
	"context"
	"time"

	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (r *MetadataRepositoryMongoDB) GetAllMetadata(ctx context.Context, limit, offset int64, search string) ([]*entities.Metadata, int64, error) {
	// Build filter based on search parameter

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"metadata_name": bson.M{"$regex": search, "$options": "i"},
		}
	}

	// Get total count
	total, err := r.db.Collection("metadata").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Set up pagination options
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset * limit).
		SetSort(bson.D{{Key: "metadata_created_at", Value: -1}}) // Sort by creation date descending

	// Execute query
	cursor, err := r.db.Collection("metadata").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var metadata []*entities.Metadata
	if err = cursor.All(ctx, &metadata); err != nil {
		return nil, 0, err
	}

	return metadata, total, nil
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
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	review := entities.Review{}
	err = reviews.FindOne(ctx, filter).Decode(&review)
	if err != nil {
		return nil, err
	}

	metadataResponse := &entities.MetadataResponse{
		ID:            metadata.MetadataProductID,
		Name:          metadata.MetadataName,
		Description:   metadata.MetadataDescription,
		Image:         metadata.MetadataImage,
		CategoryID:    metadata.MetadataCategoryID,
		SubcategoryID: metadata.MetadataSubcategoryID,
		MRP:           metadata.MetadataMRP,
		CreatedAt:     metadata.MetadataCreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     metadata.MetadataUpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		TotalStars:    review.TotalStars,
		TotalReviews:  review.TotalReviews,
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
