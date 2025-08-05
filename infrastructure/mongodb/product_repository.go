package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepositoryMongoDB implements the ProductRepository interface using MongoDB
type ProductRepositoryMongoDB struct {
	db *mongo.Database
}

// NewProductRepositoryMongoDB creates a new MongoDB product repository
func NewProductRepositoryMongoDB(db *mongo.Database) repositories.ProductRepository {
	return &ProductRepositoryMongoDB{db: db}
}

func (r *ProductRepositoryMongoDB) FetchSellerId(ctx context.Context, storeID string) (string, error) {
	collection := r.db.Collection("stores")
	cursor, err := collection.Find(ctx, bson.M{"store_id": storeID})
	if err != nil {
		return "", err
	}
	var store entities.Store
	if err := cursor.Decode(&store); err != nil {
		return "", err
	}
	return store.SellerID, nil
}

func (r *ProductRepositoryMongoDB) GetProductsForSpecificStore(ctx context.Context, sellerId string) ([]*entities.GetProductsForSpecificStoreResponse, error) {
	collection := r.db.Collection("inventory")
	collectionProduct := r.db.Collection("inventory_product")
	collectionMetadata := r.db.Collection("metadata")

	filter := bson.M{
		"seller_id": sellerId,
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

	var responseData []*entities.GetProductsForSpecificStoreResponse

	for _, inventory := range inventoryList {
		filterProduct := bson.M{
			"inventory_id": inventory.InventoryID,
		}
		cursorProduct, err := collectionProduct.Find(ctx, filterProduct)
		if err != nil {
			return nil, err
		}
		defer cursorProduct.Close(ctx)

		var inventoryProductList []entities.InventoryProduct
		if err := cursorProduct.All(ctx, &inventoryProductList); err != nil {
			return nil, err
		}

		for _, inventoryProduct := range inventoryProductList {
			filterMetadata := bson.M{
				"metadata_product_id": inventoryProduct.MetadataProductID,
			}
			cursorMetadata, err := collectionMetadata.Find(ctx, filterMetadata)
			if err != nil {
				return nil, err
			}
			defer cursorMetadata.Close(ctx)

			var metadataList []entities.Metadata
			if err := cursorMetadata.All(ctx, &metadataList); err != nil {
				return nil, err
			}

			for _, metadata := range metadataList {
				storeProduct := struct {
					InventoryId              string  `json:"inventory_id"`
					InventoryProductId       string  `json:"inventory_product_id"`
					MetadataProductId        string  `json:"metadata_product_id"`
					ProductVisibility        string  `json:"product_visibility"`
					MetadataName             string  `json:"metadata_name"`
					MetadataDescription      string  `json:"metadata_description"`
					MetadataImage            string  `json:"metadata_image"`
					MetadataCategoryId       string  `json:"metadata_category_id"`
					MetadataSubcategoryId    string  `json:"metadata_subcategory_id"`
					MetadataMrp              float64 `json:"metadata_mrp"`
					ProductQuantity          int     `json:"product_quantity"`
					ProductExpiryDate        string  `json:"product_expiry_date"`
					ProductManufacturingDate string  `json:"product_manufacturing_date"`
				}{
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
				}

				response := &entities.GetProductsForSpecificStoreResponse{
					StoreProducts: []struct {
						InventoryId              string  `json:"inventory_id"`
						InventoryProductId       string  `json:"inventory_product_id"`
						MetadataProductId        string  `json:"metadata_product_id"`
						ProductVisibility        string  `json:"product_visibility"`
						MetadataName             string  `json:"metadata_name"`
						MetadataDescription      string  `json:"metadata_description"`
						MetadataImage            string  `json:"metadata_image"`
						MetadataCategoryId       string  `json:"metadata_category_id"`
						MetadataSubcategoryId    string  `json:"metadata_subcategory_id"`
						MetadataMrp              float64 `json:"metadata_mrp"`
						ProductQuantity          int     `json:"product_quantity"`
						ProductExpiryDate        string  `json:"product_expiry_date"`
						ProductManufacturingDate string  `json:"product_manufacturing_date"`
					}{storeProduct},
				}
				responseData = append(responseData, response)
			}
		}
	}
	return responseData, nil
}

func (r *ProductRepositoryMongoDB) GetAllStores(ctx context.Context, warehouseID string) (*[]entities.Store, error) {
	collection := r.db.Collection("stores")
	filter := bson.M{
		"warehouse_id": warehouseID,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var storeList []entities.Store
	if err := cursor.All(ctx, &storeList); err != nil {
		return nil, err
	}
	return &storeList, nil
}

func (r *ProductRepositoryMongoDB) GetProductsForAllStores(ctx context.Context, allStores *[]entities.Store) ([]*entities.GetProductsForAllStoresResponse, error) {
	var response []*entities.GetProductsForAllStoresResponse

	for _, store := range *allStores {
		products, err := r.GetProductsForSpecificStore(ctx, store.SellerID)
		if err != nil {
			return nil, err
		}

		// Convert slice of pointers to slice of values
		storeProducts := make([]entities.GetProductsForSpecificStoreResponse, len(products))
		for i, product := range products {
			storeProducts[i] = *product
		}

		response = append(response, &entities.GetProductsForAllStoresResponse{
			AllStoresProducts: []struct {
				StoreID       string                                         `json:"store_id"`
				StoreProducts []entities.GetProductsForSpecificStoreResponse `json:"store_products"`
			}{
				{
					StoreID:       store.StoreID,
					StoreProducts: storeProducts,
				},
			},
		})
	}
	return response, nil
}
