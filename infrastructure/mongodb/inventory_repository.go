package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

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

func (r *InventoryRepositoryMongoDB) GetAllInventory(ctx context.Context, inventoryRequest entities.GetAllInventoryRequest) ([]entities.GetAllInventoryResponse, error) {
	collection := r.db.Collection("inventory")
	collectionProduct := r.db.Collection("inventory_product")
	collectionMetadata := r.db.Collection("metadata")

	filter := bson.M{
		"seller_id": inventoryRequest.SellerID,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var inventoryList []entities.Inventory
	if err := cursor.All(ctx, &inventoryList); err != nil {
		return nil, err
	}

	var responseInventory []entities.GetAllInventoryResponse

	for _, inventory := range inventoryList {
		filterProduct := bson.M{
			"inventory_id": inventory.InventoryID,
		}
		cursorProduct, err := collectionProduct.Find(ctx, filterProduct)
		if err != nil {
			return nil, err
		}
		defer cursorProduct.Close(ctx)

		var inventoryProduct []entities.InventoryProduct
		if err := cursorProduct.All(ctx, &inventoryProduct); err != nil {
			return nil, err
		}

		for _, product := range inventoryProduct {
			filterMetadata := bson.M{
				"metadata_product_id": product.MetadataProductID,
			}
			cursorMetadata, err := collectionMetadata.Find(ctx, filterMetadata)
			if err != nil {
				return nil, err
			}
			defer cursorMetadata.Close(ctx)

			var metadata entities.Metadata
			if err := cursorMetadata.Decode(&metadata); err != nil {
				return nil, err
			}

			responseInventory = append(responseInventory, entities.GetAllInventoryResponse{
				InventoryId:              inventory.InventoryID,
				InventoryProductId:       product.InventoryProductID,
				MetadataProductId:        metadata.MetadataProductID,
				ProductVisibility:        product.ProductVisibility,
				MetadataName:             metadata.MetadataName,
				MetadataDescription:      metadata.MetadataDescription,
				MetadataImage:            metadata.MetadataImage,
				MetadataCategoryId:       metadata.MetadataCategoryID,
				MetadataSubcategoryId:    metadata.MetadataSubcategoryID,
				MetadataMrp:              metadata.MetadataMRP,
				ProductQuantity:          product.ProductQuantity,
				ProductExpiryDate:        product.ProductExpiryDate,
				ProductManufacturingDate: product.ProductManufacturingDate,
				MetadataCreatedAt:        metadata.MetadataCreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	return responseInventory, nil
}

func (r *InventoryRepositoryMongoDB) CreateInventory(ctx context.Context, inventoryRequest entities.AddInventoryRequest) (string, error) {
	collection := r.db.Collection("inventory")

	filter := bson.M{
		"seller_id": inventoryRequest.SellerID,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return "", err
	}

	if count > 0 {
		var inventory entities.Inventory
		err := collection.FindOne(ctx, filter).Decode(&inventory)
		if err != nil {
			return "", err
		}

		return inventory.InventoryID, nil
	}

	inventoryId := primitive.NewObjectID()
	inventory := entities.Inventory{
		InventoryID: inventoryId.Hex(),
		SellerID:    inventoryRequest.SellerID,
	}

	_, err = collection.InsertOne(ctx, inventory)
	if err != nil {
		return "", err
	}

	return inventoryId.Hex(), nil
}

func (r *InventoryRepositoryMongoDB) CreateInventoryProduct(ctx context.Context, inventoryId, inventoryProductId primitive.ObjectID, inventoryRequest entities.AddInventoryRequest) error {

	collection := r.db.Collection("inventory_product")

	for _, product := range inventoryRequest.InventoryProducts {

		// Check if sellerId and metadataId combination already exists
		filter := bson.M{
			"inventory_id":        inventoryId.Hex(),
			"metadata_product_id": product.MetadataProductID,
		}

		var existingProduct entities.InventoryProduct
		err := collection.FindOne(ctx, filter).Decode(&existingProduct)
		if err == nil {
			// Product already exists, skip insertion
			continue
		} else if err != mongo.ErrNoDocuments {
			// Some other error occurred
			return err
		}

		// Product doesn't exist, proceed with insertion
		inventoryProduct := entities.InventoryProduct{
			InventoryProductID:       inventoryProductId.Hex(),
			InventoryID:              inventoryId.Hex(),
			MetadataProductID:        product.MetadataProductID,
			ProductVisibility:        product.ProductVisibility,
			ProductQuantity:          product.ProductQuantity,
			ProductExpiryDate:        product.ProductExpiryDate,
			ProductManufacturingDate: product.ProductManufacturingDate,
		}

		_, err = collection.InsertOne(ctx, inventoryProduct)
		if err != nil {
			return err
		}
	}

	return nil
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
