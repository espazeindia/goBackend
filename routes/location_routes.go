package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupLocationRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()
	var locationRepository repositories.LocationRepository = mongodb.NewLocationRepositoryMongoDB(database)
	var locationUseCase *usecase.LocationUseCase = usecase.NewLocationUseCase(locationRepository)
	var locationHandler *handlers.LocationHandler = handlers.NewLocationHandler(locationUseCase)

	// Location routes

	router.GET("/:userId", locationHandler.GetLocationForUserID)
	router.POST("/", locationHandler.CreateLocation)
	router.GET("/", locationHandler.GetLocationByAddress)
}
