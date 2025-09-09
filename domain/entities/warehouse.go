package entities

import "time"

type Warehouse struct {
	ID                        string    `json:"id" bson:"_id,omitempty"`
	WarehouseName             string    `json:"warehouseName" bson:"warehouseName"`
	WarehouseAddress          string    `json:"warehouseAddress" bson:"warehouseAddress"`
	WarehouseLeaseDetails     string    `json:"warehouseLeaseDetails" bson:"warehouseLeaseDetails"`
	WarehouseStorageCapacity  int       `json:"warehouse_storage_capacity" bson:"warehouse_storage_capacity"`
	WarehouseOperationalGuyID string    `json:"warehouse_operational_guy_id" bson:"warehouse_operational_guy_id"`
	WarehouseCreatedAt        time.Time `json:"warehouse_created_at" bson:"warehouse_created_at"`
	WarehouseUpdatedAt        time.Time `json:"warehouse_updated_at" bson:"warehouse_updated_at"`
	OwnerName                 string    `json:"ownerName" bson:"ownerName"`
	OwnerAddress              string    `json:"ownerAddress" bson:"ownerAddress"`
	OwnerPhoneNumber          string    `json:"ownerPhoneNumber" binding:"required,min=10"`
}

type CreateWarehouseRequest struct {
	WarehouseName            string `json:"warehouseName" bson:"warehouseName"`
	WarehouseAddress         string `json:"warehouseAddress" bson:"warehouseAddress"`
	WarehouseLeaseDetails    string `json:"warehouseLeaseDetails" bson:"warehouseLeaseDetails"`
	WarehouseStorageCapacity int    `json:"storage_capacity" bson:"storage_capacity"`
	OwnerName                string `json:"ownerName" bson:"ownerName"`
	OwnerAddress             string `json:"ownerAddress" bson:"ownerAddress"`
	OwnerPhoneNumber         string `json:"ownerPhoneNumber" binding:"required,min=10"`
}
type UpdateWarehouseRequest struct {
	WarehouseName             string `json:"name" bson:"name"`
	WarehouseAddress          string `json:"address" bson:"address"`
	WarehouseStorageCapacity  int    `json:"storage_capacity" bson:"storage_capacity"`
	WarehouseOperationalGuyID string `json:"operational_guy_id" bson:"operational_guy_id"`
}
