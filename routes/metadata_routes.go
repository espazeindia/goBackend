package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupMetadataRoutes(router *gin.RouterGroup) {

	database := db.GetDatabase()

	var metadataRepo repositories.MetadataRepository = mongodb.NewMetadataRepositoryMongoDB(database)

	var metadataUseCase *usecase.MetadataUseCase = usecase.NewMetadataUseCase(metadataRepo)

	var metadataHandler *handlers.MetadataHandler = handlers.NewMetadataHandler(metadataUseCase)

	router.GET("/getMetadata", metadataHandler.GetMetadata)
	router.GET("/getMetadata/:id", metadataHandler.GetMetadataByID)

	router.POST("/createMetadata", metadataHandler.CreateMetadata)

	router.PUT("/updateMetadata/:id", metadataHandler.UpdateMetadata)

	router.DELETE("/deleteMetadata/:id", metadataHandler.DeleteMetadata)

	router.POST("/add_review", metadataHandler.AddReview)

	router.GET("/getMetadataForSeller", metadataHandler.GetMetadataForSeller)

}
