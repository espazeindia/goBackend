package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"espazeBackend/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *ProductRepositoryMongoDB) GetProductsForSpecificStore(ctx context.Context, store_id string) ([]*entities.GetProductsForSpecificStoreResponse, error) {
	collection := r.db.Collection("inventory")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "store_id", Value: store_id}}}},
		bson.D{{Key: "$addFields", Value: bson.D{{Key: "_id_str", Value: bson.D{{Key: "$toString", Value: "$_id"}}}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "inventory_product"},
			{Key: "let", Value: bson.D{{Key: "invId", Value: "$_id_str"}}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$and", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$inventory_id", "$$invId"}}},
							bson.D{{Key: "$eq", Value: bson.A{"$product_visibility", true}}},
						}},
					}},
				}}},
			}},
			{Key: "as", Value: "ip"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$ip"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "metadata"},
			{Key: "let", Value: bson.D{{Key: "mdId", Value: "$ip.metadata_product_id"}}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{
							"$_id",
							bson.D{{Key: "$toObjectId", Value: "$$mdId"}},
						}},
					}},
				}}},
			}},
			{Key: "as", Value: "md"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$md"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "inventory_id", Value: "$_id_str"},
			{Key: "inventory_product_id", Value: bson.D{{Key: "$toString", Value: "$ip._id"}}},
			{Key: "metadata_product_id", Value: bson.D{{Key: "$toString", Value: "$md._id"}}},
			{Key: "product_visibility", Value: "$ip.product_visibility"},
			{Key: "product_price", Value: "$ip.product_price"},
			{Key: "metadata_name", Value: "$md.metadata_name"},
			{Key: "metadata_description", Value: "$md.metadata_description"},
			{Key: "metadata_image", Value: "$md.metadata_image"},
			{Key: "metadata_category_id", Value: "$md.metadata_category_id"},
			{Key: "metadata_subcategory_id", Value: "$md.metadata_subcategory_id"},
			{Key: "metadata_mrp", Value: "$md.metadata_mrp"},
			{Key: "product_quantity", Value: "$ip.product_quantity"},
			{Key: "product_expiry_date", Value: "$ip.product_expiry_date"},
			{Key: "product_manufacturing_date", Value: "$ip.product_manufacturing_date"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type aggResult struct {
		InventoryId              string    `bson:"inventory_id"`
		InventoryProductId       string    `bson:"inventory_product_id"`
		MetadataProductId        string    `bson:"metadata_product_id"`
		ProductVisibility        bool      `bson:"product_visibility"`
		ProductPrice             float64   `bson:"product_price"`
		MetadataName             string    `bson:"metadata_name"`
		MetadataDescription      string    `bson:"metadata_description"`
		MetadataImage            string    `bson:"metadata_image"`
		MetadataCategoryId       string    `bson:"metadata_category_id"`
		MetadataSubcategoryId    string    `bson:"metadata_subcategory_id"`
		MetadataMrp              float64   `bson:"metadata_mrp"`
		ProductQuantity          int       `bson:"product_quantity"`
		ProductExpiryDate        time.Time `bson:"product_expiry_date"`
		ProductManufacturingDate time.Time `bson:"product_manufacturing_date"`
	}

	var results []aggResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var responseData []*entities.GetProductsForSpecificStoreResponse
	for _, rdoc := range results {
		// build response inline so we don't create mismatched anonymous struct types
		resp := &entities.GetProductsForSpecificStoreResponse{
			InventoryId:              rdoc.InventoryId,
			InventoryProductId:       rdoc.InventoryProductId,
			MetadataProductId:        rdoc.MetadataProductId,
			ProductVisibility:        rdoc.ProductVisibility,
			ProductPrice:             rdoc.ProductPrice,
			MetadataName:             rdoc.MetadataName,
			MetadataDescription:      rdoc.MetadataDescription,
			MetadataImage:            rdoc.MetadataImage,
			MetadataCategoryId:       rdoc.MetadataCategoryId,
			MetadataSubcategoryId:    rdoc.MetadataSubcategoryId,
			MetadataMrp:              rdoc.MetadataMrp,
			ProductQuantity:          rdoc.ProductQuantity,
			ProductExpiryDate:        rdoc.ProductExpiryDate,
			ProductManufacturingDate: rdoc.ProductManufacturingDate,
			ProductCategoryName:      "",
			ProductSubCategoryName:   "",
			TotalStars:               "",
			TotalReviews:             "",
		}

		responseData = append(responseData, resp)
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

func (r *ProductRepositoryMongoDB) GetProductsForStoreSubcategory(ctx context.Context, storeId, subcategoryId string) ([]*entities.GetProductsForStoreSubcategory, error) {
	inventoryCollection := r.db.Collection("inventory")
	inventoryProductsCollection := r.db.Collection("inventory_product")
	storeCollection := r.db.Collection("stores")

	objectId, err := primitive.ObjectIDFromHex(storeId)
	if err != nil {
		return nil, err
	}
	var storeData entities.Store
	err = storeCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&storeData)
	if err != nil {
		return nil, err
	}

	var inventoryData entities.Inventory
	err = inventoryCollection.FindOne(ctx, bson.M{"store_id": storeId}).Decode(&inventoryData)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$eq": bson.A{
					"$inventory_id",
					bson.M{"$toString": inventoryData.InventoryID},
				},
			},
		}}},
		{{Key: "$match", Value: bson.M{
			"product_visibility": true,
		}}},
		{{Key: "$addFields", Value: bson.M{
			"metadata_product_objectId": bson.M{"$toObjectId": "$metadata_product_id"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "metadata",
			"localField":   "metadata_product_objectId",
			"foreignField": "_id",
			"as":           "metadata",
		}}},
		{{Key: "$unwind", Value: "$metadata"}},
		{{Key: "$match", Value: bson.M{
			"metadata.metadata_subcategory_id": subcategoryId,
		}}},
		{{Key: "$addFields", Value: bson.M{
			"metadata_category_objectId": bson.M{"$toObjectId": "$metadata.metadata_category_id"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "metadata_category_objectId",
			"foreignField": "_id",
			"as":           "category",
		}}},
		{{Key: "$unwind", Value: "$category"}},
		{{Key: "$addFields", Value: bson.M{
			"metadata_subcategory_objectId": bson.M{"$toObjectId": "$metadata.metadata_subcategory_id"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "subcategories",
			"localField":   "metadata_subcategory_objectId",
			"foreignField": "_id",
			"as":           "subcategory",
		}}},
		{{Key: "$unwind", Value: "$subcategory"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "reviews",
			"localField":   "metadata_product_id",
			"foreignField": "_id",
			"as":           "review",
		}}},
		{{Key: "$unwind", Value: "$review"}},
		{{Key: "$project", Value: bson.M{
			"metadata_id":                "$metadata._id",
			"metadata_name":              "$metadata.metadata_name",
			"metadata_description":       "$metadata.metadata_description",
			"metadata_image":             "$metadata.metadata_image",
			"metadata_category_id":       "$metadata.metadata_category_id",
			"metadata_subcategory_id":    "$metadata.metadata_subcategory_id",
			"metadata_mrp":               "$metadata.metadata_mrp",
			"category_name":              "$category.category_name",
			"subcategory_name":           "$subcategory.subcategory_name",
			"total_stars":                "$review.total_stars",
			"total_reviews":              "$review.total_reviews",
			"_id":                        1,
			"inventory_id":               1,
			"product_quantity":           1,
			"product_price":              1,
			"product_expiry_date":        1,
			"product_manufacturing_date": 1,
		}}},
	}

	cursor, err := inventoryProductsCollection.Aggregate(ctx, pipeline)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type aggResult struct {
		MetadataProductId        string    `bson:"metadata_id"`
		MetadataName             string    `bson:"metadata_name"`
		MetadataDescription      string    `bson:"metadata_description"`
		MetadataImage            string    `bson:"metadata_image"`
		MetadataCategoryId       string    `bson:"metadata_category_id"`
		MetadataSubcategoryId    string    `bson:"metadata_subcategory_id"`
		MetadataMrp              float64   `bson:"metadata_mrp"`
		ProductCategoryName      string    `bson:"category_name"`
		ProductSubCategoryName   string    `bson:"subcategory_name"`
		TotalStars               int       `bson:"total_stars"`
		TotalReviews             int       `bson:"total_reviews"`
		InventoryProductId       string    `bson:"_id"`
		InventoryId              string    `bson:"inventory_id"`
		ProductPrice             float64   `bson:"product_price"`
		ProductQuantity          int       `bson:"product_quantity"`
		ProductExpiryDate        time.Time `bson:"product_expiry_date"`
		ProductManufacturingDate time.Time `bson:"product_manufacturing_date"`
	}
	var cursorResults []*aggResult
	err = cursor.All(ctx, &cursorResults)
	if err != nil {
		return nil, err
	}
	var results []*entities.GetProductsForStoreSubcategory
	for _, metadataData := range cursorResults {
		results = append(results, &entities.GetProductsForStoreSubcategory{
			MetadataProductId:        metadataData.MetadataProductId,
			MetadataName:             metadataData.MetadataName,
			MetadataDescription:      metadataData.MetadataDescription,
			MetadataImage:            metadataData.MetadataImage,
			MetadataCategoryId:       metadataData.MetadataCategoryId,
			MetadataMrp:              metadataData.MetadataMrp,
			MetadataSubcategoryId:    metadataData.MetadataSubcategoryId,
			ProductCategoryName:      metadataData.ProductCategoryName,
			ProductSubCategoryName:   metadataData.ProductSubCategoryName,
			TotalStars:               metadataData.TotalStars,
			TotalReviews:             metadataData.TotalReviews,
			InventoryId:              metadataData.InventoryId,
			InventoryProductId:       metadataData.InventoryProductId,
			ProductPrice:             metadataData.ProductPrice,
			ProductQuantity:          metadataData.ProductQuantity,
			ProductExpiryDate:        metadataData.ProductExpiryDate,
			ProductManufacturingDate: metadataData.ProductManufacturingDate,
			MetadataRating: func() float64 {
				if metadataData.TotalReviews == 0 {
					return 0
				}
				return float64(metadataData.TotalStars) / float64(metadataData.TotalReviews)
			}(),
			StoreName: storeData.StoreName,
		},
		)

	}

	return results, nil
}

func (r *ProductRepositoryMongoDB) GetProductsForAllStoresSubcategory(ctx context.Context, warehouseId, subcategoryId string) ([]*entities.GetProductsForStoreSubcategory, error) {

	storesCollection := r.db.Collection("stores")

	var storesData []*entities.Store
	cursor, err := storesCollection.Find(ctx, bson.M{"warehouse_id": warehouseId})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &storesData)
	if err != nil {
		return nil, err
	}

	groupedByMetadata := make(map[string][]*entities.GetProductsForStoreSubcategory)

	for _, store := range storesData {
		apiResult, err := r.GetProductsForStoreSubcategory(ctx, store.StoreID, subcategoryId)
		if err != nil {
			return nil, err
		}
		for _, item := range apiResult {
			key := item.MetadataProductId
			if existing, ok := groupedByMetadata[key]; ok {
				// If new item has lower price, replace existing entries
				if existing[0].ProductPrice > item.ProductPrice {
					groupedByMetadata[key] = []*entities.GetProductsForStoreSubcategory{item}
				} else if existing[0].ProductPrice == item.ProductPrice {
					// Same price - check expiry date
					if existing[0].ProductExpiryDate.Compare(item.ProductExpiryDate) == -1 {
						// New item has longer expiry date, replace existing
						groupedByMetadata[key] = []*entities.GetProductsForStoreSubcategory{item}
					} else if existing[0].ProductExpiryDate.Equal(item.ProductExpiryDate) {
						// Same price and same expiry date, add to list for random selection
						groupedByMetadata[key] = append(groupedByMetadata[key], item)
					}
					// If existing has longer expiry date, ignore new item
				}
				// If existing has lower price, ignore new item
			} else {
				// Initialize entry with current item
				groupedByMetadata[key] = []*entities.GetProductsForStoreSubcategory{item}
			}
		}
	}
	var result []*entities.GetProductsForStoreSubcategory
	for _, v := range groupedByMetadata {
		if len(v) > 1 {
			data, err := utils.GenerateRandomIndex(v)
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		} else {
			result = append(result, v[0])
		}
	}
	return result, nil
}

func (r *ProductRepositoryMongoDB) GetBasicDetailsForProduct(ctx context.Context, inventoryProductID string) (*entities.GetBasicDetailsForProductResponse, error) {
	inventoryProductsCollection := r.db.Collection("inventory_product")

	objectId, err := primitive.ObjectIDFromHex(inventoryProductID)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{

		{{Key: "$match", Value: bson.M{"_id": objectId}}},

		{{Key: "$addFields", Value: bson.M{"metadata_object_id": bson.M{"$toObjectId": "$metadata_product_id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "metadata", "localField": "metadata_object_id", "foreignField": "_id", "as": "metadata"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$metadata", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$addFields", Value: bson.M{"subcategory_objectid": bson.M{"$toObjectId": "$metadata.metadata_subcategory_id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "subcategories", "localField": "subcategory_objectid", "foreignField": "_id", "as": "subcategory"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$subcategory", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$addFields", Value: bson.M{"category_objectid": bson.M{"$toObjectId": "$metadata.metadata_category_id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "categories", "localField": "category_objectid", "foreignField": "_id", "as": "category"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$lookup", Value: bson.M{"from": "reviews", "localField": "metadata_product_id", "foreignField": "_id", "as": "review"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$review", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$addFields", Value: bson.M{"inventory_Object_id": bson.M{"$toObjectId": "$inventory_id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "inventory", "localField": "inventory_Object_id", "foreignField": "_id", "as": "inventory"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$inventory", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$addFields", Value: bson.M{"store_Object_id": bson.M{"$toObjectId": "$inventory.store_id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "stores", "localField": "store_Object_id", "foreignField": "_id", "as": "store"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$store", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$project", Value: bson.M{
			"metadata_id":                "$metadata._id",
			"metadata_name":              "$metadata.metadata_name",
			"metadata_description":       "$metadata.metadata_description",
			"metadata_image":             "$metadata.metadata_image",
			"metadata_category_id":       "$metadata.metadata_category_id",
			"metadata_subcategory_id":    "$metadata.metadata_subcategory_id",
			"metadata_mrp":               "$metadata.metadata_mrp",
			"category_name":              "$category.category_name",
			"subcategory_name":           "$subcategory.subcategory_name",
			"total_stars":                "$review.total_stars",
			"total_reviews":              "$review.total_reviews",
			"_id":                        1,
			"inventory_id":               1,
			"product_quantity":           1,
			"product_price":              1,
			"product_expiry_date":        1,
			"product_manufacturing_date": 1,
			"store_name":                 "$store.store_name",
		}}},
	}

	cursor, err := inventoryProductsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type aggResult struct {
		MetadataProductId        string    `bson:"metadata_id"`
		MetadataName             string    `bson:"metadata_name"`
		MetadataDescription      string    `bson:"metadata_description"`
		MetadataImage            string    `bson:"metadata_image"`
		MetadataCategoryId       string    `bson:"metadata_category_id"`
		MetadataSubcategoryId    string    `bson:"metadata_subcategory_id"`
		MetadataMrp              float64   `bson:"metadata_mrp"`
		ProductCategoryName      string    `bson:"category_name"`
		ProductSubCategoryName   string    `bson:"subcategory_name"`
		TotalStars               int       `bson:"total_stars"`
		TotalReviews             int       `bson:"total_reviews"`
		InventoryProductId       string    `bson:"_id"`
		InventoryId              string    `bson:"inventory_id"`
		ProductPrice             float64   `bson:"product_price"`
		ProductQuantity          int       `bson:"product_quantity"`
		ProductExpiryDate        time.Time `bson:"product_expiry_date"`
		ProductManufacturingDate time.Time `bson:"product_manufacturing_date"`
		StoreName                string    `bson:"store_name"`
	}

	var products []*aggResult
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	metadataData := products[0]
	result := &entities.GetBasicDetailsForProductResponse{
		MetadataProductId:        metadataData.MetadataProductId,
		MetadataName:             metadataData.MetadataName,
		MetadataDescription:      metadataData.MetadataDescription,
		MetadataImage:            metadataData.MetadataImage,
		MetadataCategoryId:       metadataData.MetadataCategoryId,
		MetadataMrp:              metadataData.MetadataMrp,
		MetadataSubcategoryId:    metadataData.MetadataSubcategoryId,
		ProductCategoryName:      metadataData.ProductCategoryName,
		ProductSubCategoryName:   metadataData.ProductSubCategoryName,
		TotalStars:               metadataData.TotalStars,
		TotalReviews:             metadataData.TotalReviews,
		InventoryId:              metadataData.InventoryId,
		InventoryProductId:       metadataData.InventoryProductId,
		ProductPrice:             metadataData.ProductPrice,
		ProductQuantity:          metadataData.ProductQuantity,
		ProductExpiryDate:        metadataData.ProductExpiryDate,
		ProductManufacturingDate: metadataData.ProductManufacturingDate,
		MetadataRating: func() float64 {
			if metadataData.TotalReviews == 0 {
				return 0
			}
			return float64(metadataData.TotalStars) / float64(metadataData.TotalReviews)
		}(),
		StoreName: metadataData.StoreName,
	}

	return result, nil
}

func (r *ProductRepositoryMongoDB) GetProductComparisonByStore(ctx context.Context, warehouse_id string, inventoryProductID string) ([]*entities.GetProductComparisonByStoreResult, error) {
	inventoryProductsCollection := r.db.Collection("inventory_product")
	storeCollection := r.db.Collection("stores")

	productObjectID, err := primitive.ObjectIDFromHex(inventoryProductID)
	if err != nil {
		return nil, err
	}

	var productDetails *entities.InventoryProduct

	err = inventoryProductsCollection.FindOne(ctx, bson.M{"_id": productObjectID}).Decode(&productDetails)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{

		{{Key: "$match", Value: bson.M{"warehouse_id": warehouse_id}}},

		{{Key: "$addFields", Value: bson.M{"storeId": bson.M{"$toString": "$_id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "inventory", "localField": "storeId", "foreignField": "store_id", "as": "inventory"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$inventory", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$addFields", Value: bson.M{"inventoryId": bson.M{"$toString": "$inventory._id"}}}},

		{{Key: "$lookup", Value: bson.M{"from": "inventory_product", "localField": "inventoryId", "foreignField": "inventory_id", "as": "inventoryProduct"}}},

		{{Key: "$unwind", Value: bson.M{"path": "$inventoryProduct", "preserveNullAndEmptyArrays": true}}},

		{{Key: "$match", Value: bson.M{"inventoryProduct.metadata_product_id": productDetails.MetadataProductID, "inventoryProduct.product_visibility": true}}},

		{{Key: "$project", Value: bson.M{
			"storeName":        "$store_name",
			"inventoryProduct": "$inventoryProduct",
		}}},
	}

	cursor, err := storeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []*entities.GetProductComparisonByStoreResult
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil

}
