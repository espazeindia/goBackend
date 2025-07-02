package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func ConnectMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB Connection Error: ", err)
	}

	// Ping to ensure connection works
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB Ping Error: ", err)
	}

	log.Println("✅ Connected to MongoDB")
	Client = client

	// Set default database (you can make this configurable via environment variable)
	Database = client.Database("espaze_db")
}

// GetDatabase returns the database instance
func GetDatabase() *mongo.Database {
	return Database
}

// GetClient returns the MongoDB client
func GetClient() *mongo.Client {
	return Client
}
