package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	infrastructure "espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupWarehouseRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()
	var warehouseRepository repositories.WarehouseRepository = infrastructure.NewWarehouseRepositoryMongoDB(database)
	var warehouseUseCase *usecase.WarehouseUseCase = usecase.NewWarehouseUseCase(warehouseRepository)
	var warehouseHandler *handlers.WarehouseHandler = handlers.NewWarehouseHandler(warehouseUseCase)

	router.GET("/getAllWarehouse", warehouseHandler.GetAllWarehouses)
	router.GET("/:id", warehouseHandler.GetWarehouseById)
	router.POST("/createWarehouse", warehouseHandler.CreateWarehouse)
	router.PUT("/updateWarehouse/:id", warehouseHandler.UpdateWarehouse)
	router.DELETE("/:id", warehouseHandler.DeleteWarehouse)
}
