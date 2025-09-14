package handlers

import (
	"net/http"

	"espazeBackend/domain/entities"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

// WarehouseHandler handles HTTP requests for warehouse operations
type WarehouseHandler struct {
	warehouseUseCase *usecase.WarehouseUseCase
}

// NewWarehouseHandler creates a new warehouse handler
func NewWarehouseHandler(warehouseUseCase *usecase.WarehouseUseCase) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseUseCase: warehouseUseCase,
	}
}

func (h *WarehouseHandler) GetAllWarehouses(c *gin.Context) {
	response, err := h.warehouseUseCase.GetAllWarehouses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.Data,
	})
}

func (h *WarehouseHandler) GetWarehouseById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad request",
			"message": "Warehouse ID is required",
		})
		return
	}

	warehouse, err := h.warehouseUseCase.GetWarehouseById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Not found",
			"message": "Warehouse not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    warehouse,
	})
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var warehouse entities.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := h.warehouseUseCase.CreateWarehouse(c.Request.Context(), &warehouse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": response.Success,
		"message": response.Message,
	})
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad request",
			"message": "Warehouse ID is required",
		})
		return
	}

	var warehouse entities.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := h.warehouseUseCase.UpdateWarehouse(c.Request.Context(), id, &warehouse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": response.Success,
		"message": response.Message,
	})
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad request",
			"message": "Warehouse ID is required",
		})
		return
	}

	err := h.warehouseUseCase.DeleteWarehouse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "Failed to delete warehouse",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Warehouse deleted successfully",
	})
}
