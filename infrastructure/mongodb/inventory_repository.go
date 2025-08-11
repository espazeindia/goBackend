package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryRepositoryMongoDB struct {
	db *mongo.Database
}

func NewInventoryRepositoryMongoDB(db *mongo.Database) repositories.InventoryRepository {
	return &InventoryRepositoryMongoDB{db: db}
}

func (r *InventoryRepositoryMongoDB) GetAllInventory(ctx context.Context, sellerID string, offset, limit int64, search, sort string) ([]entities.GetAllInventoryResponse, int64, error) {
	collectionInventory := r.db.Collection("inventory")
	collectionProduct := r.db.Collection("inventory_product")

	// 1. Find the seller's inventory
	var inventory entities.Inventory
	if err := collectionInventory.FindOne(ctx, bson.M{"seller_id": sellerID}).Decode(&inventory); err != nil {
		return nil, 0, err
	}

	// 2. Build aggregation pipeline on inventory_product with lookups and search
	pipeline := mongo.Pipeline{
		// Filter for this seller's inventory
		{{Key: "$match", Value: bson.M{"inventory_id": inventory.InventoryID}}},
		// Convert metadata_product_id (string) to ObjectId for lookup
		{{Key: "$addFields", Value: bson.M{
			"metadata_oid": bson.M{"$toObjectId": "$metadata_product_id"},
		}}},
		// Lookup metadata
		{{Key: "$lookup", Value: bson.M{
			"from":         "metadata",
			"localField":   "metadata_oid",
			"foreignField": "_id",
			"as":           "metadata_info",
		}}},
		{{Key: "$unwind", Value: "$metadata_info"}},
		// Convert category/subcategory string IDs to ObjectIds for lookups
		{{Key: "$addFields", Value: bson.M{
			"category_oid":    bson.M{"$toObjectId": "$metadata_info.metadata_category_id"},
			"subcategory_oid": bson.M{"$toObjectId": "$metadata_info.metadata_subcategory_id"},
		}}},
		// Lookup category and subcategory
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

	// 3. Apply search on metadata name and subcategory name
	if search != "" && search != `""` {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{
			"$or": []bson.M{
				{"metadata_info.metadata_name": bson.M{"$regex": search, "$options": "i"}},
				{"subcategory_info.subcategory_name": bson.M{"$regex": search, "$options": "i"}},
			},
		}}})
	}

	// 4. Count total
	countPipeline := append(append(mongo.Pipeline{}, pipeline...), bson.D{{Key: "$count", Value: "total"}})
	countCursor, err := collectionProduct.Aggregate(ctx, countPipeline)
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

	// 5. Sorting
	sortStage := bson.D{{Key: "metadata_info.metadata_created_at", Value: -1}, {Key: "_id", Value: 1}}
	switch sort {
	case "name_asc":
		sortStage = bson.D{{Key: "metadata_info.metadata_name", Value: 1}, {Key: "_id", Value: 1}}
	case "name_desc":
		sortStage = bson.D{{Key: "metadata_info.metadata_name", Value: -1}, {Key: "_id", Value: 1}}
	case "mrp_asc":
		sortStage = bson.D{{Key: "metadata_info.metadata_mrp", Value: 1}, {Key: "_id", Value: 1}}
	case "mrp_desc":
		sortStage = bson.D{{Key: "metadata_info.metadata_mrp", Value: -1}, {Key: "_id", Value: 1}}
	case "oldest":
		sortStage = bson.D{{Key: "metadata_info.metadata_created_at", Value: 1}, {Key: "_id", Value: 1}}
	}

	// 6. Data fetch with pagination and projection
	dataPipeline := append(pipeline,
		bson.D{{Key: "$sort", Value: sortStage}},
		bson.D{{Key: "$skip", Value: offset * limit}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"inventory_id":               "$inventory_id",
			"inventory_product_id":       bson.M{"$toString": "$_id"},
			"metadata_product_id":        bson.M{"$toString": "$metadata_info._id"},
			"product_visibility":         "$product_visibility",
			"product_price":              "$product_price",
			"metadata_name":              "$metadata_info.metadata_name",
			"metadata_description":       "$metadata_info.metadata_description",
			"metadata_image":             "$metadata_info.metadata_image",
			"metadata_category_id":       "$metadata_info.metadata_category_id",
			"metadata_subcategory_id":    "$metadata_info.metadata_subcategory_id",
			"metadata_mrp":               "$metadata_info.metadata_mrp",
			"metadata_hsn_code":          "$metadata_info.hsn_code",
			"product_quantity":           "$product_quantity",
			"product_expiry_date":        bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M:%S", "date": "$product_expiry_date"}},
			"product_manufacturing_date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M:%S", "date": "$product_manufacturing_date"}},
			"metadata_created_at":        bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M:%S", "date": "$metadata_info.metadata_created_at"}},
			"metadata_category_name":     "$category_info.category_name",
			"metadata_subcategory_name":  "$subcategory_info.subcategory_name",
		}}},
	)

	cursor, err := collectionProduct.Aggregate(ctx, dataPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []entities.GetAllInventoryResponse
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *InventoryRepositoryMongoDB) CreateInventory(ctx context.Context, inventoryRequest *entities.AddInventoryRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("inventory")
	inventoryCollection := r.db.Collection("inventory_product")

	filter := bson.M{
		"seller_id": inventoryRequest.SellerID,
	}

	var inventory *entities.Inventory
	err := collection.FindOne(ctx, filter).Decode(&inventory)
	if err != nil && err != mongo.ErrNoDocuments {
		return &entities.MessageResponse{
			Message: "Database Error",
			Error:   "Db Error",
			Success: false,
		}, err
	}
	if err == mongo.ErrNoDocuments {
		sellerCollection := r.db.Collection("sellers")
		objectId, err := primitive.ObjectIDFromHex(inventoryRequest.SellerID)
		if err != nil {
			return &entities.MessageResponse{
				Error:   "Error converting ObjectIdFromHex",
				Message: "Error creating object id from seller id ",
				Success: false,
			}, err
		}
		filter := bson.M{"_id": objectId}
		var sellerData *entities.Seller
		err = sellerCollection.FindOne(ctx, filter).Decode(&sellerData)
		if err != nil {
			return &entities.MessageResponse{
				Error:   "Db Error",
				Message: "Database Error",
				Success: false,
			}, err
		}
		inventoryData := &entities.Inventory{
			SellerID: inventoryRequest.SellerID,
			StoreId:  sellerData.StoreID,
		}
		response, err := collection.InsertOne(ctx, inventoryData)
		if err != nil {
			return &entities.MessageResponse{
				Error:   "Db Error",
				Message: "Database Error",
				Success: false,
			}, err
		}
		inventoryId, ok := response.InsertedID.(primitive.ObjectID)
		if !ok {
			return &entities.MessageResponse{
				Error:   "Db Error",
				Message: "Database Error",
				Success: false,
			}, err
		}
		var AllInventoryProducts []*entities.InventoryProduct
		now := time.Now()
		for _, metadataId := range inventoryRequest.MetadataProductID {
			InventoryProduct := &entities.InventoryProduct{
				InventoryID:              inventoryId.Hex(),
				MetadataProductID:        metadataId,
				ProductVisibility:        false,
				ProductQuantity:          0,
				ProductExpiryDate:        now,
				ProductManufacturingDate: now,
			}
			AllInventoryProducts = append(AllInventoryProducts, InventoryProduct)
		}
		docs := make([]interface{}, len(AllInventoryProducts))
		for i, v := range AllInventoryProducts {
			docs[i] = v
		}
		_, err = inventoryCollection.InsertMany(ctx, docs)
		if err != nil {
			return &entities.MessageResponse{
				Error:   "Db Error",
				Message: "Database Error",
				Success: false,
			}, err
		}
		return &entities.MessageResponse{
			Message: "Inventory Added Successfully",
			Success: true,
		}, err

	}

	var AllInventoryProducts []*entities.InventoryProduct
	now := time.Now()
	for _, metadataId := range inventoryRequest.MetadataProductID {
		InventoryProduct := &entities.InventoryProduct{
			InventoryID:              inventory.InventoryID,
			MetadataProductID:        metadataId,
			ProductVisibility:        false,
			ProductQuantity:          0,
			ProductExpiryDate:        now,
			ProductManufacturingDate: now,
		}
		AllInventoryProducts = append(AllInventoryProducts, InventoryProduct)
	}
	docs := make([]interface{}, len(AllInventoryProducts))
	for i, v := range AllInventoryProducts {
		docs[i] = v
	}

	_, err = inventoryCollection.InsertMany(ctx, docs)
	if err != nil {
		return &entities.MessageResponse{
			Error:   "Db Error",
			Message: "Database Error",
			Success: false,
		}, err
	}
	return &entities.MessageResponse{
		Message: "Inventory Added Successfully",
		Success: true,
	}, err

}

func (r *InventoryRepositoryMongoDB) UpdateInventory(ctx context.Context, inventoryRequest entities.UpdateInventoryRequest) error {
	collection := r.db.Collection("inventory_product")

	filter := bson.M{
		"inventory_id":         inventoryRequest.InventoryID,
		"inventory_product_id": inventoryRequest.InventoryProductID,
	}

	update := bson.M{
		"$set": bson.M{
			"product_visibility":         inventoryRequest.ProductVisibility,
			"product_quantity":           inventoryRequest.ProductQuantity,
			"product_price":              inventoryRequest.ProductPrice,
			"product_expiry_date":        inventoryRequest.ProductExpiryDate,
			"product_manufacturing_date": inventoryRequest.ProductManufacturingDate,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *InventoryRepositoryMongoDB) DeleteInventory(ctx context.Context, inventoryRequest entities.DeleteInventoryRequest) error {
	collection := r.db.Collection("inventory_product")

	filter := bson.M{
		"inventory_id":         inventoryRequest.InventoryID,
		"inventory_product_id": inventoryRequest.InventoryProductID,
	}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *InventoryRepositoryMongoDB) GetInventoryById(ctx context.Context, inventoryRequest entities.GetInventoryByIdRequest) (*entities.GetInventoryByIdResponse, error) {
	collection := r.db.Collection("inventory")
	collectionProduct := r.db.Collection("inventory_product")
	collectionMetadata := r.db.Collection("metadata")

	// First, get the inventory
	filterInventory := bson.M{
		"inventory_id": inventoryRequest.InventoryID,
	}

	var inventory entities.Inventory
	err := collection.FindOne(ctx, filterInventory).Decode(&inventory)
	if err != nil {
		return nil, err
	}

	// Get the inventory product
	filterProduct := bson.M{
		"inventory_id": inventoryRequest.InventoryID,
	}

	var inventoryProduct entities.InventoryProduct
	err = collectionProduct.FindOne(ctx, filterProduct).Decode(&inventoryProduct)
	if err != nil {
		return nil, err
	}

	// Get the metadata
	filterMetadata := bson.M{
		"metadata_product_id": inventoryProduct.MetadataProductID,
	}

	var metadata entities.Metadata
	err = collectionMetadata.FindOne(ctx, filterMetadata).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	response := &entities.GetInventoryByIdResponse{
		InventoryId:              inventory.InventoryID,
		InventoryProductId:       inventoryProduct.InventoryProductID,
		MetadataProductId:        metadata.MetadataProductID,
		ProductVisibility:        inventoryProduct.ProductVisibility,
		ProductPrice:             inventoryProduct.ProductPrice,
		MetadataName:             metadata.MetadataName,
		MetadataDescription:      metadata.MetadataDescription,
		MetadataImage:            metadata.MetadataImage,
		MetadataCategoryId:       metadata.MetadataCategoryID,
		MetadataSubcategoryId:    metadata.MetadataSubcategoryID,
		MetadataMrp:              metadata.MetadataMRP,
		ProductQuantity:          inventoryProduct.ProductQuantity,
		ProductExpiryDate:        inventoryProduct.ProductExpiryDate,
		ProductManufacturingDate: inventoryProduct.ProductManufacturingDate,
		MetadataCreatedAt:        metadata.MetadataCreatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}
