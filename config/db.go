package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

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
}