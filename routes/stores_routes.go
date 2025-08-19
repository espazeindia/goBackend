package routes

import (
	db "espazeBackend/config"
	"espazeBackend/domain/repositories"
	"espazeBackend/handlers"
	"espazeBackend/infrastructure/mongodb"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

func SetupStoreRoutes(router *gin.RouterGroup) {
	database := db.GetDatabase()
	var storeRepository repositories.StoreRepository = mongodb.NewStoreRepositoryMongoDB(database)
	var storeUseCase *usecase.StoreUseCase = usecase.NewStoreUseCase(storeRepository)
	var storeHandler *handlers.StoreHandler = handlers.NewStoreHandler(storeUseCase)

	// Store management routes
	router.GET("/", storeHandler.GetAllStores)                                  // GET /stores?warehouse_id=xxx&limit=10&offset=0&search=xxx
	router.GET("/getAllStoresForCutomer", storeHandler.GetAllStoresForCustomer) // GET /stores?warehouse_id=xxx&limit=10&offset=0&search=xxx
	router.GET("/:id", storeHandler.GetStoreById)                               // GET /stores/:id
	router.POST("/createStore", storeHandler.CreateStore)                       // POST /stores
	router.PUT("/:id", storeHandler.UpdateStore)                                // PUT /stores/:id
	router.DELETE("/:id", storeHandler.DeleteStore)                             // DELETE /stores/:id

	// Additional store routes
	router.GET("/seller/:seller_id", storeHandler.GetStoreBySellerId) // GET /stores/seller/:seller_id
	router.PATCH("/:id/racks", storeHandler.UpdateStoreRacks)         // PATCH /stores/:id/racks
}
