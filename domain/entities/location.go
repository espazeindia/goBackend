package entities

type Location struct {
	UserID          string `json:"user_id" bson:"user_id"`
	LocationAddress string `json:"location_address" bson:"location_address"`
	Coordinates     string `json:"coordinates" bson:"coordinates"`
}
