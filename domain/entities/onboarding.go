package entities

type SellerBasicDetail struct {
	Name        string `json:"name" bson:"name"`
	Address     string `json:"address" bson:"address"`
	Gstin       string `json:"gstin" bson:"gstin"`
	Pan         string `json:"pan" bson:"pan"`
	CompanyName string `json:"companyName" bson:"companyName"`
	PIN         int    `json:"pin" bson:"pin,min=6"`
}

type OperationalGuyRegistrationRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Address     string `json:"address" binding:"required,min=10"`
	Email       string `json:"email" binding:"required,email"`
	Pan         string `json:"pan" binding:"required,min=10"`
	Password    string `json:"-" binding:"required,min=5"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
	WarehouseId string `json:"warehouseId" binding:"required"`
}

type OperationalGuyGetRequest struct {
	Address     string `json:"address" binding:"required,min=10"`
	Password    string `json:"-" binding:"required,min=5"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10"`
}

type OperationalOnboarding struct {
	Password string `json:"-" binding:"required,min=5"`
}

type OperationalGuyGetRespone struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	Pan           string `json:"pan" `
	Password      string `json:"password"`
	PhoneNumber   string `json:"phoneNumber"`
	WarehouseId   string `json:"warehouseId"`
	WarehouseName string `json:"warehouse_name" bson:"warehouse_name"`
}
