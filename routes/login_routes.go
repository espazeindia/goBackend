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
	router.GET("/seller/verifyOTP", loginHandler.VerifyOTP)
	router.GET("/seller/getOTP", loginHandler.GetOTP)
	router.GET("/seller/verifyPin", loginHandler.VerifyPin)
	router.GET("/customer/verifyOTP", loginHandler.VerifyOTPForCustomer)
	router.GET("/customer/getOTP", loginHandler.GetOTPForCustomer)
	router.GET("/customer/verifyPin", loginHandler.VerifyPinForCustomer)

	router.POST("/customer/basicSetup", loginHandler.CustomerBasicSetup)

	router.POST("/admin/login", loginHandler.LoginAdmin)
	router.POST("/admin/register", loginHandler.RegisterAdmin)

	//router.POST("/seller/addBasicData", loginHandler.AddBasicData)

	// router.POST("/customer", loginHandler.LoginCustomer)
	// router.POST("/customer/register", loginHandler.RegisterCustomer)
}
