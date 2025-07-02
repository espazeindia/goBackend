package handlers

import (
	"net/http"
	"strconv"

	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// GetProducts handles GET request to retrieve all products with pagination
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Get query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid limit parameter",
		})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid offset parameter",
		})
		return
	}

	// Get products from use case
	result, err := h.productUseCase.GetAllProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
