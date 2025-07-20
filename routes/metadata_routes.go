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

	router.GET("/metadata", metadataHandler.GetMetadata)
	router.GET("/metadata/:id", metadataHandler.GetMetadataByID)

	router.POST("/metadata", metadataHandler.CreateMetadata)

	router.PUT("/metadata/:id", metadataHandler.UpdateMetadata)

	router.DELETE("/metadata/:id", metadataHandler.DeleteMetadata)

	router.POST("/metadata/add_review", metadataHandler.AddReview)

}
