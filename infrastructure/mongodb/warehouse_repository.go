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

// WarehouseRepositoryMongoDB implements the WarehouseRepository interface using MongoDB
type WarehouseRepositoryMongoDB struct {
	db *mongo.Database
}

// NewWarehouseRepositoryMongoDB creates a new MongoDB warehouse repository
func NewWarehouseRepositoryMongoDB(db *mongo.Database) repositories.WarehouseRepository {
	return &WarehouseRepositoryMongoDB{db: db}
}

func (r *WarehouseRepositoryMongoDB) GetAllWarehouses(ctx context.Context) (*entities.MessageResponse, error) {
	collection := r.db.Collection("warehouses")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "No warehouse found",
			Message: "Not able to fetch warehousees",
		}, err
	}
	defer cursor.Close(ctx)

	var warehouses []*entities.Warehouse

	if err := cursor.All(ctx, &warehouses); err != nil {
		return &entities.MessageResponse{
			Success: false,
			Error:   "DB error",
			Message: "database error",
		}, err
	}

	return &entities.MessageResponse{
		Success: true,
		Data:    warehouses,
	}, nil
}

func (r *WarehouseRepositoryMongoDB) GetWarehouseById(ctx context.Context, id string) (*entities.Warehouse, error) {
	collection := r.db.Collection("warehouses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var warehouse entities.Warehouse
	err = collection.FindOne(ctx, filter).Decode(&warehouse)
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *WarehouseRepositoryMongoDB) CreateWarehouse(ctx context.Context, warehouse *entities.CreateWarehouseRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("warehouses")

	now := time.Now()

	// Create document with ObjectID
	warehouseData := &entities.Warehouse{
		WarehouseName:             warehouse.WarehouseName,
		WarehouseAddress:          warehouse.WarehouseAddress,
		WarehouseLeaseDetails:     warehouse.WarehouseLeaseDetails,
		WarehouseStorageCapacity:  warehouse.WarehouseStorageCapacity,
		WarehouseOperationalGuyID: "",
		OwnerName:                 warehouse.OwnerName,
		OwnerAddress:              warehouse.OwnerAddress,
		OwnerPhoneNumber:          warehouse.OwnerPhoneNumber,
		WarehouseCreatedAt:        now,
		WarehouseUpdatedAt:        now,
	}

	response, err := collection.InsertOne(ctx, warehouseData)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "Db error",
		}, err
	}
	_, ok := response.InsertedID.(primitive.ObjectID)
	if !ok {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "Db error",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Warehouse Created Successfully",
	}, nil
}

func (r *WarehouseRepositoryMongoDB) UpdateWarehouse(ctx context.Context, id string, warehouse *entities.UpdateWarehouseRequest) (*entities.MessageResponse, error) {
	collection := r.db.Collection("warehouses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Internal Server Error",
			Error:   "Error in converting objectIdFromHex of warehouse",
		}, err
	}
	now := time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"warehouse_name":               warehouse.WarehouseName,
			"warehouse_address":            warehouse.WarehouseAddress,
			"warehouse_storage_capacity":   warehouse.WarehouseStorageCapacity,
			"warehouse_operational_guy_id": warehouse.WarehouseOperationalGuyID,
			"warehouse_updated_at":         now,
		},
	}

	response, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &entities.MessageResponse{
			Success: false,
			Message: "Database Error",
			Error:   "Db error",
		}, err
	}
	if response.MatchedCount == 0 {
		return &entities.MessageResponse{
			Success: false,
			Message: "Warehouse Not Found",
			Error:   "No matching document",
		}, err
	}
	return &entities.MessageResponse{
		Success: true,
		Message: "Warehouse Updated Successfully",
	}, nil
}

func (r *WarehouseRepositoryMongoDB) DeleteWarehouse(ctx context.Context, id string) error {
	collection := r.db.Collection("warehouses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}
