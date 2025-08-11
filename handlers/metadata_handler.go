package handlers

import (
	"net/http"
	"strconv"

	"espazeBackend/domain/entities"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

type MetadataHandler struct {
	metadataUseCase *usecase.MetadataUseCase
}

func NewMetadataHandler(metadataUseCase *usecase.MetadataUseCase) *MetadataHandler {
	return &MetadataHandler{
		metadataUseCase: metadataUseCase,
	}
}

// GetMetadata retrieves all metadata with pagination
func (h *MetadataHandler) GetMetadata(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	search := c.DefaultQuery("search", "")

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

	result, err := h.metadataUseCase.GetAllMetadata(c.Request.Context(), limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve metadata: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

func (h *MetadataHandler) GetMetadataForSeller(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	search := c.DefaultQuery("search", "")
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

	result, err := h.metadataUseCase.GetAllMetadataForSeller(c.Request.Context(), limit, offset, search, seller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve metadata: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetMetadataByID retrieves a metadata by ID
func (h *MetadataHandler) GetMetadataByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Metadata ID is required",
			"message": "Metadata ID is invalid",
		})
		return
	}

	result, err := h.metadataUseCase.GetMetadataByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Metadata not found: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CreateMetadata creates a new metadata
func (h *MetadataHandler) CreateMetadata(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	var req entities.CreateMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
			"message": "Request body is invalid",
		})
		return
	}

	response, err := h.metadataUseCase.CreateMetadata(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create metadata: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
		return
	}
	if response.Success {
		c.JSON(http.StatusCreated, gin.H{
			"success": response.Success,
			"message": response.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"message": response.Message,
			"error":   response.Error,
		})
	}

}

// UpdateMetadata updates an existing metadata
func (h *MetadataHandler) UpdateMetadata(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Metadata ID is required",
			"message": "Metadata ID is empty",
		})
		return
	}

	var req entities.UpdateMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
			"message": "Invalid request body",
		})
		return
	}

	result, err := h.metadataUseCase.UpdateMetadata(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update metadata: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	if result.Success {
		c.JSON(http.StatusAccepted, gin.H{
			"success": result.Success,
			"error":   result.Error,
			"message": result.Message,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": result.Success,
			"error":   result.Error,
			"message": result.Message,
		})
		return
	}

}

// DeleteMetadata deletes a metadata by ID
func (h *MetadataHandler) DeleteMetadata(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Metadata ID is required",
			"message": "Metadata ID is invalid",
		})
		return
	}

	result, err := h.metadataUseCase.DeleteMetadata(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete metadata: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
		return
	}
	if result.Success {
		c.JSON(http.StatusOK, gin.H{
			"success": result.Success,
			"message": result.Message,
		})
		return

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": result.Success,
			"message": result.Message,
			"error":   result.Error,
		})
		return
	}

}

func (h *MetadataHandler) AddReview(c *gin.Context) {
	var req entities.AddReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
			"message": "Request body is invalid",
		})
	}
	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Rating must be between 1 and 5",
			"message": "Star rating is out of the range 1 to 5",
		})
		return
	}
	err := h.metadataUseCase.AddReview(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to add review: " + err.Error(),
			"message": "Some Internal Server Error Occured",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review added successfully",
	})
}
