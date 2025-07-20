package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryUseCase *usecase.InventoryUseCaseInterface
}

func NewInventoryHandler(inventoryUseCase *usecase.InventoryUseCaseInterface) *InventoryHandler {
	return &InventoryHandler{inventoryUseCase: inventoryUseCase}
}

func (h *InventoryHandler) GetAllInventory(c *gin.Context) {
	var inventoryRequest entities.GetAllInventoryRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory, err := h.inventoryUseCase.GetAllInventory(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) AddInventory(c *gin.Context) {
	var inventoryRequest entities.AddInventoryRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.inventoryUseCase.AddInventory(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inventory added successfully"})
}

func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	var inventoryRequest entities.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.inventoryUseCase.UpdateInventory(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully"})
}

func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	var inventoryRequest entities.DeleteInventoryRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.inventoryUseCase.DeleteInventory(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inventory deleted successfully"})
}

func (h *InventoryHandler) GetInventoryById(c *gin.Context) {
	var inventoryRequest entities.GetInventoryByIdRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory, err := h.inventoryUseCase.GetInventoryById(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inventory)
}
