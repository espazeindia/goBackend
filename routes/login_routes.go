package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupLoginRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()

	var loginRepo repositories.LoginRepository = mongodb.NewLoginRepositoryMongoDB(database)

	var loginUseCase *usecase.LoginUseCaseInterface = usecase.NewLoginUseCase(loginRepo)

	var loginHandler *handlers.LoginHandler = handlers.NewLoginHandler(loginUseCase)

	// Public routes (no authentication required)
	router.POST("/operational_guy/login", loginHandler.LoginOperationalGuy)
	router.POST("/operational_guy/register", loginHandler.RegisterOperationalGuy)
	router.POST("/seller/register", loginHandler.RegisterSeller)
	router.GET("/seller/verifyOTP", loginHandler.VerifyOTP)
	// router.POST("/customer", loginHandler.LoginCustomer)
	// router.POST("/customer/register", loginHandler.RegisterCustomer)
}
