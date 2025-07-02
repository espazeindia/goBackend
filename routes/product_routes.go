package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configures product-related routes using onion architecture
func SetupProductRoutes(router *gin.RouterGroup) {
	// Initialize dependencies
	database := db.GetDatabase()

	// Infrastructure layer
	var productRepo repositories.ProductRepository = mongodb.NewProductRepositoryMongoDB(database)

	// Use case layer
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Handler layer
	productHandler := handlers.NewProductHandler(productUseCase)

	// Product routes

	// GET /api/products - Get all products with pagination
	router.GET("/products", productHandler.GetProducts)
}
