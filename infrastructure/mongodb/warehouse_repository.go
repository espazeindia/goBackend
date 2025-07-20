package mongodb

import (
	"context"
	"espazeBackend/domain/entities"
	"espazeBackend/domain/repositories"

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

func (r *WarehouseRepositoryMongoDB) GetAllWarehouses(ctx context.Context) ([]*entities.Warehouse, error) {
	collection := r.db.Collection("warehouses")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var warehouses []*entities.Warehouse
	if err := cursor.All(ctx, &warehouses); err != nil {
		return nil, err
	}
	return warehouses, nil
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

func (r *WarehouseRepositoryMongoDB) CreateWarehouse(ctx context.Context, warehouse *entities.Warehouse) error {
	collection := r.db.Collection("warehouses")

	// Convert string ID to ObjectID for MongoDB
	objectID, err := primitive.ObjectIDFromHex(warehouse.ID)
	if err != nil {
		return err
	}

	// Create document with ObjectID
	doc := bson.M{
		"_id":                          objectID,
		"warehouse_name":               warehouse.WarehouseName,
		"warehouse_address":            warehouse.WarehouseAddress,
		"warehouse_coordinates":        warehouse.WarehouseCoordinates,
		"warehouse_storage_capacity":   warehouse.WarehouseStorageCapacity,
		"warehouse_operational_guy_id": warehouse.WarehouseOperationalGuyID,
		"warehouse_created_at":         warehouse.WarehouseCreatedAt,
		"warehouse_updated_at":         warehouse.WarehouseUpdatedAt,
	}

	_, err = collection.InsertOne(ctx, doc)
	return err
}

func (r *WarehouseRepositoryMongoDB) UpdateWarehouse(ctx context.Context, id string, warehouse *entities.Warehouse) error {
	collection := r.db.Collection("warehouses")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"warehouse_name":               warehouse.WarehouseName,
			"warehouse_address":            warehouse.WarehouseAddress,
			"warehouse_coordinates":        warehouse.WarehouseCoordinates,
			"warehouse_storage_capacity":   warehouse.WarehouseStorageCapacity,
			"warehouse_operational_guy_id": warehouse.WarehouseOperationalGuyID,
			"warehouse_updated_at":         warehouse.WarehouseUpdatedAt,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
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
