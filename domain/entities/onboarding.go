package entities

type SellerBasicDetail struct {
	Name        string `json:"name" bson:"name"`
	Address     string `json:"address" bson:"address"`
	Gstin       string `json:"gstin" bson:"Gstin"`
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
