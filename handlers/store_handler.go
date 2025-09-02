package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeUseCase *usecase.StoreUseCase
}

func NewStoreHandler(storeUseCase *usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{
		storeUseCase: storeUseCase,
	}
}

// GetAllStores handles GET /stores - Get all stores with pagination and search
func (h *StoreHandler) GetAllStores(c *gin.Context) {
	// Parse query parameters
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "warehouse_id is required",
			"message": "Warehouse ID must be provided as a query parameter",
		})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	search := c.Query("search")

	request := entities.GetAllStoresRequest{
		WarehouseID: warehouseID,
		Limit:       limit,
		Offset:      offset,
		Search:      search,
	}

	response, err := h.storeUseCase.GetAllStores(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *StoreHandler) GetAllStoresForCustomer(c *gin.Context) {
	// Parse query parameters
	warehouseID := c.Query("warehouse_id")
	if warehouseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "warehouse_id is required",
			"message": "Warehouse ID must be provided as a query parameter",
		})
		return
	}

	search := c.Query("search")

	response, err := h.storeUseCase.GetAllStoresForCustomer(c.Request.Context(), warehouseID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetStoreById handles GET /stores/:id - Get store by ID
func (h *StoreHandler) GetStoreById(c *gin.Context) {
	storeId := c.Param("id")
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "store_id is required",
			"message": "Store ID must be provided in the URL path",
		})
		return
	}

	response, err := h.storeUseCase.GetStoreById(c.Request.Context(), storeId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Store not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateStore handles POST /stores - Create a new store
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var request entities.CreateStoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Request Data ",
			"message": "Request data is invalid",
		})
		return
	}

	response, err := h.storeUseCase.CreateStore(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"success": response.Success,
		"message": response.Message,
	})

}

// UpdateStore handles PUT /stores/:id - Update an existing store
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	storeId := c.Param("id")
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "store_id is required",
			"message": "Store ID must be provided in the URL path",
		})
		return
	}

	var request entities.UpdateStoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := h.storeUseCase.UpdateStore(c.Request.Context(), storeId, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to update store",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteStore handles DELETE /stores/:id - Delete a store
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	storeId := c.Param("id")
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "store_id is required",
			"message": "Store ID must be provided in the URL path",
		})
		return
	}

	response, err := h.storeUseCase.DeleteStore(c.Request.Context(), storeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to delete store",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStoreBySellerId handles GET /stores/seller/:seller_id - Get store by seller ID
func (h *StoreHandler) GetStoreBySellerId(c *gin.Context) {
	sellerId := c.Param("seller_id")
	if sellerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "seller_id is required",
			"message": "Seller ID must be provided in the URL path",
		})
		return
	}

	response, err := h.storeUseCase.GetStoreBySellerId(c.Request.Context(), sellerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Store not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateStoreRacks handles PATCH /stores/:id/racks - Update store occupied racks
func (h *StoreHandler) UpdateStoreRacks(c *gin.Context) {
	storeId := c.Param("id")
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "store_id is required",
			"message": "Store ID must be provided in the URL path",
		})
		return
	}

	var request struct {
		OccupiedRacks int `json:"occupied_racks" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	err := h.storeUseCase.UpdateStoreRacks(c.Request.Context(), storeId, request.OccupiedRacks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to update store racks",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Store racks updated successfully",
	})
}
