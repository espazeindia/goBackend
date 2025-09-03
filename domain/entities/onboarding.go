package entities

type SellerBasicDetail struct {
	Name        string `json:"name" bson:"name"`
	ShopAddress string `json:"address" bson:"address"`
	Gstin       string `json:"Gstin" bson:"Gstin"`
	Pan         string `json:"pan" bson:"pan"`
	CompanyName string `json:"companyName" bson:"companyName"`
	ShopName    string `json:"ShopName" bson:"ShopName"`
	PIN         int    `json:"pin" bson:"pin,min=6"`
}
