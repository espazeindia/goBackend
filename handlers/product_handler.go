package handlers

import (
	"net/http"

	"espazeBackend/domain/entities"
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

func (h *ProductHandler) GetProductsForSpecificStore(c *gin.Context) {
	store_id := c.Query("store_id")
	if store_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "store id missing in request",
		})
	}

	response, err := h.productUseCase.GetProductsForSpecificStore(c.Request.Context(), store_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "An unexpected error occurred",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *ProductHandler) GetProductsForAllStores(c *gin.Context) {
	var getProductsForAllStoresRequest entities.GetProductsForAllStoresRequest
	if err := c.ShouldBindJSON(&getProductsForAllStoresRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}

	response, err := h.productUseCase.GetProductsForAllStores(c.Request.Context(), getProductsForAllStoresRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "An unexpected error occurred",
		})
	}

	c.JSON(http.StatusOK, response)
}
