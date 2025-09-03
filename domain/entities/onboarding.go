package entities

type SellerBasicDetail struct {
	Name        string `json:"name" bson:"name"`
	ShopAddress string `json:"address" bson:"address"`
	Gstin       string `json:"gstin" bson:"Gstin"`
	Pan         string `json:"pan" bson:"pan"`
	CompanyName string `json:"companyName" bson:"companyName"`
	ShopName    string `json:"shopName" bson:"ShopName"`
	PIN         int    `json:"pin" bson:"pin,min=6"`
}
