package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupInventoryRoutes(router *gin.RouterGroup) {

	database := db.GetDatabase()

	var inventoryRepo repositories.InventoryRepository = mongodb.NewInventoryRepositoryMongoDB(database)

	var inventoryUseCase *usecase.InventoryUseCaseInterface = usecase.NewInventoryUseCase(inventoryRepo)

	var inventoryHandler *handlers.InventoryHandler = handlers.NewInventoryHandler(inventoryUseCase)

	router.GET("/getAllInventory", inventoryHandler.GetAllInventory)
	router.POST("/addInventory", inventoryHandler.AddInventory)
	router.PUT("/updateInventory", inventoryHandler.UpdateInventory)
	router.DELETE("/deleteInventory", inventoryHandler.DeleteInventory)
	router.GET("/getInventoryById", inventoryHandler.GetInventoryById)
	router.POST("/addInventoryByExcel", inventoryHandler.AddInventoryByExcel)
	router.GET("/getAllInventoryRequests", inventoryHandler.GetAllInventoryRequests)
	router.GET("/acceptProduct", inventoryHandler.AcceptVisibility)

}
