package entities

type Location struct {
	LocationId      string `json:"id" bson:"_id,omitempty"`
	UserID          string `json:"user_id" bson:"user_id"`
	LocationAddress string `json:"location_address" bson:"location_address"`
	Coordinates     string `json:"coordinates" bson:"coordinates"`
	Self            bool   `json:"self" bson:"self"`
	BuildingType    string `json:"building_type" bson:"building_type"`
	PhoneNumber     string `json:"phone" bson:"phone"`
	Name            string `json:"name" bson:"name"`
}
