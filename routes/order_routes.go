package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()

	var orderRepository repositories.OrderRepository = mongodb.NewOrderRepositoryMongoDB(database)
	var orderUsecase *usecase.OrderUsecase = usecase.NewOrderUsecase(orderRepository)
	var orderHandler *handlers.OrderHandler = handlers.NewOrderHandler(orderUsecase)

	router.GET("/getAllOrders", orderHandler.GetAllOrders)
	router.POST("/createOrder", orderHandler.CreateOrder)
	router.GET("/getOrderByOrderID", orderHandler.GetOrderByOrderID)
	router.GET("/getOrderByUserID", orderHandler.GetOrderByUserID)
	router.GET("/getOrderBySellerID", orderHandler.GetOrderBySellerID)

}
