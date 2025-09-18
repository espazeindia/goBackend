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

type StoreRepositoryMongoDB struct {
	db *mongo.Database
}

func NewStoreRepositoryMongoDB(db *mongo.Database) repositories.StoreRepository {
	return &StoreRepositoryMongoDB{
		db: db,
	}
}

func (r *StoreRepositoryMongoDB) GetAllStores(ctx context.Context, request entities.GetAllStoresRequest) (*entities.PaginatedStoresResponse, error) {
	collection := r.db.Collection("stores")

	// Build filter
	filter := bson.M{"warehouse_id": request.WarehouseID}
	if request.Search != "" {
		filter["store_name"] = bson.M{"$regex": request.Search, "$options": "i"}
	}

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Set up pagination options
	limit := int64(10) // default limit
	if request.Limit > 0 {
		limit = request.Limit
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(request.Offset).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stores []entities.Store
	if err := cursor.All(ctx, &stores); err != nil {
		return nil, err
	}

	// Calculate pagination info
	hasNext := request.Offset+limit < total
	hasPrevious := request.Offset > 0

	return &entities.PaginatedStoresResponse{
		Stores:      stores,
		Total:       total,
		Limit:       limit,
		Offset:      request.Offset,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}, nil
}

func (r *StoreRepositoryMongoDB) GetAllStoresForCustomer(ctx context.Context, warehouseId, search string) ([]*entities.Store, error) {
	collection := r.db.Collection("stores")

	// Build filter
	filter := bson.M{"warehouse_id": warehouseId}
	if search != "" {
		filter["store_name"] = bson.M{"$regex": search, "$options": "i"}
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stores []*entities.Store
	var ans []*entities.Store
	if err := cursor.All(ctx, &stores); err != nil {
		return nil, err
	}
	if len(stores) != 0 {
		now := time.Now()
		allStoreOption := &entities.Store{
			StoreID:       "0",
			SellerID:      "dummy",
			WarehouseID:   warehouseId,
			StoreName:     "All Stores",
			StoreAddress:  "dummy",
			StoreContact:  "dummy",
			NumberOfRacks: 0,
			OccupiedRacks: 0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		ans = append(ans, allStoreOption)

	}
	ans = append(ans, stores...)
	return ans, nil
}

func (r *StoreRepositoryMongoDB) GetStoreById(ctx context.Context, storeId string) (*entities.Store, error) {
	collection := r.db.Collection("stores")
	filter := bson.M{"store_id": storeId}

	var store entities.Store
	err := collection.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &store, nil
}

func (r *StoreRepositoryMongoDB) CreateStore(ctx context.Context, request *entities.CreateStoreRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("stores")

	now := time.Now()
	store := &entities.Store{
		StoreName:     request.StoreName,
		StoreAddress:  request.StoreAddress,
		StoreContact:  request.StoreContact,
		WarehouseID:   request.WarehouseID,
		SellerID:      request.SellerID,
		NumberOfRacks: request.NumberOfRacks,
		OccupiedRacks: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	response, err := collection.InsertOne(ctx, store)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Db error",
			Message: "Database Error",
		}, err
	}
	_, ok := response.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Error:   "Db error",
			Message: "Database Error",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Store Creatd Successfully",
	}, nil

}

func (r *StoreRepositoryMongoDB) UpdateStore(ctx context.Context, storeId string, store *entities.Store) error {
	collection := r.db.Collection("stores")

	// Update timestamp
	store.UpdatedAt = time.Now()

	filter := bson.M{"store_id": storeId}
	update := bson.M{
		"$set": bson.M{
			"store_name":      store.StoreName,
			"store_address":   store.StoreAddress,
			"store_contact":   store.StoreContact,
			"number_of_racks": store.NumberOfRacks,
			"occupied_racks":  store.OccupiedRacks,
			"updated_at":      store.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *StoreRepositoryMongoDB) DeleteStore(ctx context.Context, storeId string) error {
	collection := r.db.Collection("stores")
	filter := bson.M{"store_id": storeId}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *StoreRepositoryMongoDB) GetStoreBySellerId(ctx context.Context, sellerId string) (*entities.Store, error) {
	collection := r.db.Collection("stores")
	filter := bson.M{"seller_id": sellerId}

	var store entities.Store
	err := collection.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &store, nil
}

func (r *StoreRepositoryMongoDB) GetStoresByWarehouseId(ctx context.Context, warehouseId string) ([]entities.Store, error) {
	collection := r.db.Collection("stores")
	filter := bson.M{"warehouse_id": warehouseId}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stores []entities.Store
	if err := cursor.All(ctx, &stores); err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *StoreRepositoryMongoDB) UpdateStoreRacks(ctx context.Context, storeId string, occupiedRacks int) error {
	collection := r.db.Collection("stores")

	filter := bson.M{"store_id": storeId}
	update := bson.M{
		"$set": bson.M{
			"occupied_racks": occupiedRacks,
			"updated_at":     time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
