package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "your-module-name/db"
)

func main() {
    // 📦 Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("❌ Error loading .env file")
    }

    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("❌ MONGO_URI is not set in .env")
    }

    // 🔌 Connect to MongoDB
    db.ConnectMongoDB(mongoURI)

    // 🚀 Setup Gin
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "MongoDB is connected!",
        })
    })

    r.Run(":8080")
}