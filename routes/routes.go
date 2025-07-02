package routes

import (
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

	api := router.Group("/api")
	{
		SetupProductRoutes(api)
	}

}
