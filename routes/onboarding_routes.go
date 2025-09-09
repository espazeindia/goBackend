package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupOnboardingRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()

	var onboardingRepo repositories.OnboardingRepository = mongodb.NewOnboardingRepositoryMongoDB(database)

	var onboardingUseCase *usecase.OnboardingUseCaseInterface = usecase.NewOnboardingUseCase(onboardingRepo)

	var onboardingHandler *handlers.OnboardingHandler = handlers.NewOnboardingHandler(onboardingUseCase)

	router.POST("/seller/addBasicDetail", onboardingHandler.AddBasicDetail)
	router.GET("/seller/getBasicDetails", onboardingHandler.GetBasicDetail)

	router.POST("/admin/onboarding", onboardingHandler.OnboardingAdmin)

	router.POST("/operational_guy/register", onboardingHandler.RegisterOperationalGuy)
	router.GET("/operational_guy/getOperations", onboardingHandler.GetOperationalGuy)
	router.PUT("/operational_guy/editOperations", onboardingHandler.EditOperationalGuy)
	router.PUT("/operational_guy/onboarding", onboardingHandler.OnboardingOperationalGuy)

}
