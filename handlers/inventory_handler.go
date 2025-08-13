package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryUseCase *usecase.InventoryUseCaseInterface
}

func NewInventoryHandler(inventoryUseCase *usecase.InventoryUseCaseInterface) *InventoryHandler {
	return &InventoryHandler{inventoryUseCase: inventoryUseCase}
}

func (h *InventoryHandler) GetAllInventory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "")
	sellerInterface, isPresent := c.Get("user_id")

	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid token",
			"message": "Token is invalid",
		})
		return
	}

	seller, ok := sellerInterface.(string)
	if !ok || seller == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid token",
			"message": "Token is invalid",
		})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid limit parameter",
			"message": "Limit parameter is invalid",
		})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid offset parameter",
			"message": "Offset parameter is invalid",
		})
		return
	}

	inventory, err := h.inventoryUseCase.GetAllInventory(c.Request.Context(), seller, offset, limit, search, sort)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": inventory, "success": true})
}

func (h *InventoryHandler) AddInventory(c *gin.Context) {
	seller_id, isPresent := c.Get("user_id")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid token",
			"message": "Token is invalid",
		})
		return
	}

	seller, ok := seller_id.(string)
	if !ok || seller == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid token",
			"message": "Token is invalid",
		})
		return
	}

	var inventoryRequest *entities.AddInventoryRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Request Body is Invalid",
			"message": "Invalide Request Body",
		})
		return
	}
	inventoryRequest.SellerID = seller

	response, err := h.inventoryUseCase.AddInventory(c.Request.Context(), inventoryRequest)

	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"message": response.Message,
			"error":   response.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
	})
}

func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	var inventoryRequest entities.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false,
			"error":   err.Error(),
			"message": "Invalid Request Body"})
		return
	}

	response, err := h.inventoryUseCase.UpdateInventory(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": response.Message, "success": response.Success})
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
	inventoryRequest := c.Query("id")
	if inventoryRequest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Id not present", "error": "Id required"})
		return
	}

	inventory, err := h.inventoryUseCase.GetInventoryById(c.Request.Context(), inventoryRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false, "message": "Some Internal Server Erorr Occured"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": inventory})
}
