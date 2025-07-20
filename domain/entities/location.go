package entities

type Location struct {
	ID              string `json:"id" bson:"_id"`
	LocationAddress string `json:"location_address" bson:"location_address"`
	Coordinates     string `json:"coordinates" bson:"coordinates"`
}
