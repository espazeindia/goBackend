package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupCategorySubcategoryRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()
	var categorySubcategoryRepository repositories.CategorySubcategoryRepository = mongodb.NewCategorySubcategoryRepositoryMongoDB(database)
	var categorySubcategoryUseCase *usecase.CategorySubcategoryUseCase = usecase.NewCategorySubcategoryUseCase(categorySubcategoryRepository)
	var categorySubcategoryHandler *handlers.CategorySubcategoryHandler = handlers.NewCategorySubcategoryHandler(categorySubcategoryUseCase)

	// Category routes
	router.GET("/getCategories", categorySubcategoryHandler.GetAllCategories)
	// router.GET("/with-subcategories", categorySubcategoryHandler.GetAllCategoriesWithSubcategories)
	// router.GET("/:id", categorySubcategoryHandler.GetCategoryById)
	// router.GET("/:id/with-subcategories", categorySubcategoryHandler.GetCategoryWithSubcategories)
	router.POST("/createCategory", categorySubcategoryHandler.CreateCategory)
	router.PUT("/updateCategory/:id", categorySubcategoryHandler.UpdateCategory)
	// router.DELETE("/:id", categorySubcategoryHandler.DeleteCategory)

	// // Subcategory routes
	router.GET("/getSubCategories", categorySubcategoryHandler.GetAllSubcategories)
	router.GET("/getSubcategoryByCategoryId/:id", categorySubcategoryHandler.GetSubcategoryByCategoryId)
	router.POST("/createSubCategory", categorySubcategoryHandler.CreateSubcategory)
	// router.PUT("/subcategory/:id", categorySubcategoryHandler.UpdateSubcategory)
	// router.DELETE("/subcategory/:id", categorySubcategoryHandler.DeleteSubcategory)
}
