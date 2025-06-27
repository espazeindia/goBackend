package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/gin-contrib/cors"
    db "espazeBackend/config"
)
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.Header().Set("X-Frame-Options", "DENY")
        c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
        c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
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
    router := gin.Default()

    router.Use(cors.Default())
    router.Use(SecurityHeaders())
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "MongoDB is connected!",
        })
    })

    router.Run(":8080")
}