package entities

import "time"

type Warehouse struct {
	ID                        string    `json:"id" bson:"_id,omitempty"`
	WarehouseName             string    `json:"warehouse_name" bson:"warehouse_name"`
	WarehouseAddress          string    `json:"warehouse_address" bson:"warehouse_address"`
	WarehouseCoordinates      string    `json:"warehouse_coordinates" bson:"warehouse_coordinates"`
	WarehouseStorageCapacity  int       `json:"warehouse_storage_capacity" bson:"warehouse_storage_capacity"`
	WarehouseOperationalGuyID string    `json:"warehouse_operational_guy_id" bson:"warehouse_operational_guy_id"`
	WarehouseCreatedAt        time.Time `json:"warehouse_created_at" bson:"warehouse_created_at"`
	WarehouseUpdatedAt        time.Time `json:"warehouse_updated_at" bson:"warehouse_updated_at"`
}

type CreateWarehouseRequest struct {
	WarehouseName            string `json:"name" bson:"name"`
	WarehouseAddress         string `json:"address" bson:"address"`
	WarehouseCoordinates     string `json:"coordinates" bson:"coordinates"`
	WarehouseStorageCapacity int    `json:"storage_capacity" bson:"storage_capacity"`
}
type UpdateWarehouseRequest struct {
	WarehouseName             string `json:"name" bson:"name"`
	WarehouseAddress          string `json:"address" bson:"address"`
	WarehouseCoordinates      string `json:"coordinates" bson:"coordinates"`
	WarehouseStorageCapacity  int    `json:"storage_capacity" bson:"storage_capacity"`
	WarehouseOperationalGuyID string `json:"operational_guy_id" bson:"operational_guy_id"`
}
