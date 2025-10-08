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

	router.GET("/getProductsForSpecificStore", productHandler.GetProductsForSpecificStore)
	router.GET("/getProductsForAllStores", productHandler.GetProductsForAllStores)
	router.GET("/getAllProductsForSubcategory", productHandler.GetAllProductsForSubcategory)
	router.GET("/getBasicDetailsForProduct", productHandler.GetBasicDetailsForProduct)
	router.GET("/getProductComparisonByStores", productHandler.GetProductComparisonByStore)
}
