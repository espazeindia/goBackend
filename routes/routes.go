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
	// Create a protected route group with authentication middleware
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		metadata := protected.Group("/metadata")
		{
			SetupMetadataRoutes(metadata)
		}

		inventory := protected.Group("/inventory")
		{
			SetupInventoryRoutes(inventory)
		}

		products := protected.Group("/products")
		{
			SetupProductRoutes(products)
		}

		store := protected.Group("/store")
		{
			SetupStoreRoutes(store)
		}

		warehouse := protected.Group("/warehouse")
		{
			SetupWarehouseRoutes(warehouse)
		}

		location := protected.Group("/location")
		{
			SetupLocationRoutes(location)
		}

		category_subcategory := protected.Group("/category")
		{
			SetupCategorySubcategoryRoutes(category_subcategory)
		}

		order := protected.Group("/order")
		{
			SetupOrderRoutes(order)
		}
	}
}
