package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"
	"fmt"
	"log"

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
	err := collectionInventory.FindOne(ctx, bson.M{"seller_id": sellerID}).Decode(&inventory)
	if err == mongo.ErrNoDocuments {
		return nil, 0, nil
	}
	if err != nil {
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
	case "asc":
		sortStage = bson.D{{Key: "product_price", Value: 1}, {Key: "_id", Value: 1}}
	case "desc":
		sortStage = bson.D{{Key: "product_price", Value: -1}, {Key: "_id", Value: 1}}
	case "mrp_asc":
		sortStage = bson.D{{Key: "metadata_info.metadata_mrp", Value: 1}, {Key: "_id", Value: 1}}
	case "mrp_desc":
		sortStage = bson.D{{Key: "metadata_info.metadata_mrp", Value: -1}, {Key: "_id", Value: 1}}

	}

	// 6. Data fetch with pagination and projection
	dataPipeline := append(pipeline,
		bson.D{{Key: "$sort", Value: sortStage}},
		bson.D{{Key: "$skip", Value: offset * limit}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"inventory_id":            "$inventory_id",
			"inventory_product_id":    bson.M{"$toString": "$_id"},
			"metadata_product_id":     bson.M{"$toString": "$metadata_info._id"},
			"product_visibility":      "$product_visibility",
			"product_price":           "$product_price",
			"metadata_name":           "$metadata_info.metadata_name",
			"metadata_description":    "$metadata_info.metadata_description",
			"metadata_image":          "$metadata_info.metadata_image",
			"metadata_category_id":    "$metadata_info.metadata_category_id",
			"metadata_subcategory_id": "$metadata_info.metadata_subcategory_id",
			"metadata_mrp":            "$metadata_info.metadata_mrp",
			"metadata_hsn_code":       "$metadata_info.hsn_code",
			"product_quantity":        "$product_quantity",
			"product_expiry_date": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$product_expiry_date"}, "date"}},
				"then": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M:%S", "date": "$product_expiry_date"}},
				"else": "$product_expiry_date",
			}},
			"product_manufacturing_date": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$product_manufacturing_date"}, "date"}},
				"then": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M:%S", "date": "$product_manufacturing_date"}},
				"else": "$product_manufacturing_date",
			}},
			"metadata_created_at": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$metadata_info.metadata_created_at"}, "date"}},
				"then": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M:%S", "date": "$metadata_info.metadata_created_at"}},
				"else": "$metadata_info.metadata_created_at",
			}},
			"metadata_category_name":    "$category_info.category_name",
			"metadata_subcategory_name": "$subcategory_info.subcategory_name",
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

func (r *InventoryRepositoryMongoDB) UpdateInventory(ctx context.Context, inventoryRequest entities.UpdateInventoryRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("inventory_product")

	objectId, err := primitive.ObjectIDFromHex(inventoryRequest.InventoryProductID)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating object from productId",
			Error:   "ObjectIdFromHex Error",
		}, err
	}

	filter := bson.M{
		"_id": objectId,
	}

	// Parse date strings to time.Time
	var expiryDate, manufacturingDate time.Time
	var err1, err2 error

	if inventoryRequest.ProductExpiryDate != "" {
		expiryDate, err1 = time.Parse("2006-01-02 15:04:05", inventoryRequest.ProductExpiryDate)
		if err1 != nil {
			// Try alternative format
			expiryDate, err1 = time.Parse("2006-01-02", inventoryRequest.ProductExpiryDate)
		}
	}

	if inventoryRequest.ProductManufacturingDate != "" {
		manufacturingDate, err2 = time.Parse("2006-01-02 15:04:05", inventoryRequest.ProductManufacturingDate)
		if err2 != nil {
			// Try alternative format
			manufacturingDate, err2 = time.Parse("2006-01-02", inventoryRequest.ProductManufacturingDate)
		}
	}

	update := bson.M{
		"$set": bson.M{
			"product_visibility": inventoryRequest.ProductVisibility,
			"product_quantity":   inventoryRequest.ProductQuantity,
			"product_price":      inventoryRequest.ProductPrice,
		},
	}

	// Only add date fields if they were successfully parsed
	if err1 == nil && !expiryDate.IsZero() {
		update["$set"].(bson.M)["product_expiry_date"] = expiryDate
	}
	if err2 == nil && !manufacturingDate.IsZero() {
		update["$set"].(bson.M)["product_manufacturing_date"] = manufacturingDate
	}

	response, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "Db Error",
		}, err
	}

	if response.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "No Matching Document ",
		}, err
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Product Updated Successfully",
	}, nil
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

func (r *InventoryRepositoryMongoDB) GetInventoryById(ctx context.Context, InventoryId string) (*entities.GetInventoryByIdResponse, error) {
	collectionProduct := r.db.Collection("inventory_product")
	objectId, err := primitive.ObjectIDFromHex(InventoryId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectId,
	}
	var ProductDetails *entities.InventoryProduct
	err = collectionProduct.FindOne(ctx, filter).Decode(&ProductDetails)
	if err != nil {
		return nil, err
	}

	// Create metadata repository to get metadata details
	metadataRepo := NewMetadataRepositoryMongoDB(r.db)
	metadataResponse, err := metadataRepo.GetMetadataByID(ctx, ProductDetails.MetadataProductID)
	if err != nil {
		return nil, err
	}

	result := &entities.GetInventoryByIdResponse{
		InventoryProductId:       ProductDetails.InventoryProductID,
		MetadataProductId:        metadataResponse.ID,
		ProductVisibility:        ProductDetails.ProductVisibility,
		MetadataName:             metadataResponse.Name,
		MetadataDescription:      metadataResponse.Description,
		MetadataImage:            metadataResponse.Image,
		MetadataCategoryId:       metadataResponse.CategoryID,
		MetadataSubcategoryId:    metadataResponse.SubcategoryID,
		MetadataMrp:              metadataResponse.MRP,
		ProductQuantity:          ProductDetails.ProductQuantity,
		ProductPrice:             ProductDetails.ProductPrice,
		ProductExpiryDate:        ProductDetails.ProductExpiryDate,
		ProductManufacturingDate: ProductDetails.ProductManufacturingDate,
		MetadataCreatedAt:        metadataResponse.CreatedAt,
		MetadataCategoryName:     metadataResponse.CategoryName,
		MetadataSubcategoryName:  metadataResponse.SubCategoryName,
		MetadataHSNCode:          metadataResponse.HsnCode,
		TotalReviews:             metadataResponse.TotalReviews,
		TotalStars:               metadataResponse.TotalStars,
	}

	return result, nil
}

func (r *InventoryRepositoryMongoDB) AddInventoryByExcel(ctx context.Context, inventoryRequest *entities.AddInventoryByExcelRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("inventory")
	inventoryCollection := r.db.Collection("inventory_product")

	session, err := r.db.Client().StartSession()
	if err != nil {
		return &entities.MessageResponse{
			Message: "Database Error",
			Error:   "Db Error",
			Success: false,
		}, err
	}
	defer session.EndSession(ctx)

	var resp *entities.MessageResponse
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		filter := bson.M{
			"seller_id": inventoryRequest.SellerID,
		}

		var inventory *entities.Inventory
		errFind := collection.FindOne(sc, filter).Decode(&inventory)
		if errFind != nil && errFind != mongo.ErrNoDocuments {
			return nil, errFind
		}

		layout := "02-01-2006"

		if errFind == mongo.ErrNoDocuments {
			sellerCollection := r.db.Collection("sellers")
			objectId, errConv := primitive.ObjectIDFromHex(inventoryRequest.SellerID)
			if errConv != nil {
				return nil, errConv
			}
			sellerFilter := bson.M{"_id": objectId}
			var sellerData *entities.Seller
			if err := sellerCollection.FindOne(sc, sellerFilter).Decode(&sellerData); err != nil {
				return nil, err
			}
			inventoryData := &entities.Inventory{
				SellerID: inventoryRequest.SellerID,
				StoreId:  sellerData.StoreID,
			}
			response, err := collection.InsertOne(sc, inventoryData)
			if err != nil {
				return nil, err
			}
			inventoryId, ok := response.InsertedID.(primitive.ObjectID)
			if !ok {
				return nil, fmt.Errorf("failed to get inserted inventory id")
			}
			var allInventoryProducts []*entities.InventoryProduct
			for _, mp := range inventoryRequest.MetadataProducts {
				expiryDate, perr := time.Parse(layout, mp.ProductExpiryDate)
				if perr != nil {
					log.Printf("error parsing expiry date: %v", perr)
					expiryDate = time.Time{}
				}
				manufacturingDate, merr := time.Parse(layout, mp.ProductManufacturingDate)
				if merr != nil {
					log.Printf("error parsing manufacturing date: %v", merr)
					manufacturingDate = time.Time{}
				}
				p := &entities.InventoryProduct{
					InventoryID:              inventoryId.Hex(),
					MetadataProductID:        mp.ProductMetadataId,
					ProductVisibility:        false,
					ProductQuantity:          mp.ProductQuantity,
					ProductExpiryDate:        expiryDate,
					ProductManufacturingDate: manufacturingDate,
				}
				allInventoryProducts = append(allInventoryProducts, p)
			}
			docs := make([]interface{}, len(allInventoryProducts))
			for i, v := range allInventoryProducts {
				docs[i] = v
			}
			if len(docs) > 0 {
				if _, err := inventoryCollection.InsertMany(sc, docs); err != nil {
					return nil, err
				}
			}
			resp = &entities.MessageResponse{Message: "Inventory Added Successfully", Success: true}
			return resp, nil
		}

		// Existing inventory path
		for _, mp := range inventoryRequest.MetadataProducts {
			expiryDate, perr := time.Parse(layout, mp.ProductExpiryDate)
			if perr != nil {
				expiryDate = time.Time{}
			}
			manufacturingDate, merr := time.Parse(layout, mp.ProductManufacturingDate)
			if merr != nil {
				manufacturingDate = time.Time{}
			}

			var inventoryProductData []*entities.InventoryProduct
			cursor, err := inventoryCollection.Find(sc, bson.M{"metadata_product_id": mp.ProductMetadataId, "inventory_id": inventory.InventoryID})
			if err != nil {
				return nil, err
			}
			if err := cursor.All(sc, &inventoryProductData); err != nil {
				return nil, err
			}
			fmt.Print(inventoryProductData, mp)
			updatedOrInserted := false
			for _, inventoryProduct := range inventoryProductData {
				if inventoryProduct.ProductQuantity == 0 && !inventoryProduct.ProductVisibility {
					objectId, err := primitive.ObjectIDFromHex(inventoryProduct.InventoryProductID)
					if err != nil {
						return nil, err
					}
					result, err := inventoryCollection.UpdateByID(sc, objectId, bson.M{"$set": bson.M{"product_quantity": mp.ProductQuantity, "product_price": mp.ProductPrice, "product_expiry_date": expiryDate, "product_manufacturing_date": manufacturingDate}})
					if err != nil {
						return nil, err
					}
					if result.MatchedCount == 0 {
						return nil, fmt.Errorf("no data updated")
					}
					updatedOrInserted = true
					break
				} else if mp.ProductPrice == inventoryProduct.ProductPrice && inventoryProduct.ProductExpiryDate.Compare(expiryDate) == 0 && !inventoryProduct.ProductVisibility {
					objectId, err := primitive.ObjectIDFromHex(inventoryProduct.InventoryProductID)
					if err != nil {
						return nil, err
					}
					result, err := inventoryCollection.UpdateByID(sc, objectId, bson.M{"$set": bson.M{"product_quantity": inventoryProduct.ProductQuantity + mp.ProductQuantity}})
					if err != nil {
						return nil, err
					}
					if result.MatchedCount == 0 {
						return nil, fmt.Errorf("no data updated")
					}
					updatedOrInserted = true
					break
				}
			}
			if !updatedOrInserted {
				newProduct := &entities.InventoryProduct{
					InventoryID:              inventory.InventoryID,
					MetadataProductID:        mp.ProductMetadataId,
					ProductVisibility:        false,
					ProductQuantity:          mp.ProductQuantity,
					ProductExpiryDate:        expiryDate,
					ProductManufacturingDate: manufacturingDate,
					ProductPrice:             mp.ProductPrice,
				}
				if _, err := inventoryCollection.InsertOne(sc, newProduct); err != nil {
					return nil, err
				}
			}
		}

		resp = &entities.MessageResponse{Message: "Inventory Added Successfully", Success: true}
		return resp, nil
	})

	if err != nil {
		return &entities.MessageResponse{
			Message: "Database Error",
			Error:   err.Error(),
			Success: false,
		}, err
	}

	return resp, nil
}

func (r *InventoryRepositoryMongoDB) GetAllInventoryRequests(ctx context.Context, operational_id string, offset, limit int64, search string) ([]*entities.GetAllInventoryRequestResponse, int64, error) {
	warehouseCollection := r.db.Collection("warehouses")

	session, err := r.db.Client().StartSession()
	if err != nil {
		return nil, 0, err
	}
	defer session.EndSession(ctx)
	var result []*entities.GetAllInventoryRequestResponse
	var total int64
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: bson.M{"warehouse_operational_guy_id": operational_id}}},

			{{Key: "$addFields", Value: bson.M{"warehouseIdString": bson.M{"$toString": "$_id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "stores", "localField": "warehouseIdString", "foreignField": "warehouse_id", "as": "storeInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$storeInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$addFields", Value: bson.M{"sellerObjectId": bson.M{"$toObjectId": "$storeInfo.seller_id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "sellers", "localField": "sellerObjectId", "foreignField": "_id", "as": "sellerInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$sellerInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$addFields", Value: bson.M{"storeIdString": bson.M{"$toString": "$storeInfo._id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "inventory", "localField": "storeIdString", "foreignField": "store_id", "as": "inventoryInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$inventoryInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$addFields", Value: bson.M{"inventoryIdString": bson.M{"$toString": "$inventoryInfo._id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "inventory_product", "localField": "inventoryIdString", "foreignField": "inventory_id", "as": "productInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$productInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$match", Value: bson.M{"productInfo.product_visibility": false}}},

			{{Key: "$addFields", Value: bson.M{"metadataObjectId": bson.M{"$toObjectId": "$productInfo.metadata_product_id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "metadata", "localField": "metadataObjectId", "foreignField": "_id", "as": "metadataInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$metadataInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$addFields", Value: bson.M{"categoryObjectId": bson.M{"$toObjectId": "$metadataInfo.metadata_category_id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "categories", "localField": "categoryObjectId", "foreignField": "_id", "as": "categoryInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$categoryInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$addFields", Value: bson.M{"subcategoryObjectId": bson.M{"$toObjectId": "$metadataInfo.metadata_subcategory_id"}}}},

			{{Key: "$lookup", Value: bson.M{"from": "subcategories", "localField": "subcategoryObjectId", "foreignField": "_id", "as": "subcategoryInfo"}}},

			{{Key: "$unwind", Value: bson.M{"path": "$subcategoryInfo", "preserveNullAndEmptyArrays": true}}},

			{{Key: "$project", Value: bson.M{
				"store_name":           "$storeInfo.store_name",
				"seller_name":          "$sellerInfo.name",
				"metadata_name":        "$metadataInfo.metadata_name",
				"metadata_image":       "$metadataInfo.metadata_image",
				"metadata_description": "$metadataInfo.metadata_description",
				"metadata_mrp":         "$metadataInfo.metadata_mrp",
				"hsn_code":             "$metadataInfo.hsn_code",
				"category_name":        "$categoryInfo.category_name",
				"subcategory_name":     "$subcategoryInfo.subcategory_name",
				"_id":                  "$productInfo._id",
				"inventory_id":         "$productInfo.inventory_id",
				"visibility":           "$productInfo.product_visibility",
				"quantity":             "$productInfo.product_quantity",
				"price":                "$productInfo.product_price",
				"expiry_date":          "$productInfo.product_expiry_date",
				"manufacturing_date":   "$productInfo.product_manufacturing_date",
			}}},
		}

		// Add search filter only if search is not empty
		if search != "" {
			searchStage := bson.D{{Key: "$match", Value: bson.M{
				"$or": bson.A{
					bson.M{"metadata_name": bson.M{"$regex": search, "$options": "i"}},
					bson.M{"subcategory_name": bson.M{"$regex": search, "$options": "i"}},
					bson.M{"store_name": bson.M{"$regex": search, "$options": "i"}},
					bson.M{"seller_name": bson.M{"$regex": search, "$options": "i"}},
				},
			}}}
			pipeline = append(pipeline, searchStage)
		}
		countPipeline := append(append(mongo.Pipeline{}, pipeline...), bson.D{{Key: "$count", Value: "total"}})

		totalCursor, err := warehouseCollection.Aggregate(ctx, countPipeline)
		if err != nil {
			return nil, err
		}
		var totalData []*struct {
			Total int64 `json:"total" bson:"total"`
		}
		err = totalCursor.All(ctx, &totalData)
		if err != nil {
			return nil, err
		}
		if len(totalData) > 0 {
			total = totalData[0].Total
		} else {
			total = 0
		}

		dataPipeline := append(pipeline, bson.D{{Key: "$skip", Value: limit * offset}}, bson.D{{Key: "$limit", Value: limit}})

		cursor, err := warehouseCollection.Aggregate(ctx, dataPipeline)
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, &result)
		if err != nil {
			return nil, err
		}
		return nil, nil

	})
	if err != nil {
		return nil, 0, err
	}
	return result, total, err
}

func (r *InventoryRepositoryMongoDB) AcceptVisibility(ctx context.Context, productId string) (*entities.MessageResponse, error) {
	collection := r.db.Collection("inventory_product")

	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Error creating object from productId",
			Error:   "ObjectIdFromHex Error",
		}, err
	}

	filter := bson.M{
		"_id": objectId,
	}

	update := bson.M{
		"$set": bson.M{
			"product_visibility": true,
		},
	}

	response, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "Db Error",
		}, err
	}

	if response.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "No Matching Document ",
		}, err
	}

	return &entities.MessageResponse{
		Success: true,
		Message: "Product Accepted Successfully",
	}, nil
}
