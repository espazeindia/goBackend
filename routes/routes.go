package routes

import (
	"espazeBackend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {

	// Health check route
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "healthy",
			"message": "backend is running",
		})
	})

	// Login routes (no authentication required)
	login := router.Group("/login")
	{
		SetupLoginRoutes(login)
	}

	// Protected routes (authentication required)
	// Apply authentication middleware to all routes below this point
	router.Use(middlewares.AuthMiddleware())

	metadata := router.Group("/metadata")
	{
		SetupMetadataRoutes(metadata)
	}

	inventory := router.Group("/inventory")
	{
		SetupInventoryRoutes(inventory)
	}

	products := router.Group("/products")
	{
		SetupProductRoutes(products)
	}

	store := router.Group("/store")
	{
		SetupStoreRoutes(store)
	}

	warehouse := router.Group("/warehouse")
	{
		SetupWarehouseRoutes(warehouse)
	}

	location := router.Group("/location")
	{
		SetupLocationRoutes(location)
	}

	category_subcategory := router.Group("/category")
	{
		SetupCategorySubcategoryRoutes(category_subcategory)
	}

	order := router.Group("/order")
	{
		SetupOrderRoutes(order)
	}
}
