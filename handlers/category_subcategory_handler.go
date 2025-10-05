package handlers

import (
	"net/http"
	"strconv"

	"espazeBackend/domain/entities"
	"espazeBackend/usecase"

	"github.com/gin-gonic/gin"
)

// CategorySubcategoryHandler handles HTTP requests for category and subcategory operations
type CategorySubcategoryHandler struct {
	categorySubcategoryUseCase *usecase.CategorySubcategoryUseCase
}

// NewCategorySubcategoryHandler creates a new category subcategory handler
func NewCategorySubcategoryHandler(categorySubcategoryUseCase *usecase.CategorySubcategoryUseCase) *CategorySubcategoryHandler {
	return &CategorySubcategoryHandler{
		categorySubcategoryUseCase: categorySubcategoryUseCase,
	}
}

// Category handlers
func (h *CategorySubcategoryHandler) GetCategories(c *gin.Context) {
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

	categories, err := h.categorySubcategoryUseCase.GetCategories(c.Request.Context(), limit, offset, &search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}

func (h *CategorySubcategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categorySubcategoryUseCase.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}

func (h *CategorySubcategoryHandler) GetAllSubcategories(c *gin.Context) {
	categoryID := c.Param("id")
	search := c.Query("search")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is invalid",
			"message": "Category ID is required",
		})
		return
	}

	sub_categories, err := h.categorySubcategoryUseCase.GetAllSubcategories(c.Request.Context(), categoryID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sub_categories,
	})
}

func (h *CategorySubcategoryHandler) GetSubcategories(c *gin.Context) {
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
	subcategories, err := h.categorySubcategoryUseCase.GetSubcategories(c.Request.Context(), limit, offset, &search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subcategories,
	})
}

func (h *CategorySubcategoryHandler) CreateCategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	var request entities.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "Request body is invalid",
		})
		return
	}

	result, err := h.categorySubcategoryUseCase.CreateCategory(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": result.Success,
			"error":   result.Error,
			"message": result.Message,
		})
		return
	}
	if result.Success {
		c.JSON(http.StatusCreated, gin.H{
			"success": result.Success,
			"message": result.Message,
		})
		return
	}

}

func (h *CategorySubcategoryHandler) CreateSubcategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	var request entities.CreateSubcategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := h.categorySubcategoryUseCase.CreateSubcategory(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": response.Message,
	})
}

func (h *CategorySubcategoryHandler) GetSubcategoryByCategoryId(c *gin.Context) {
	CategoryID := c.Param("id")
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
	if CategoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid ca",
			"message": "Subcategory ID is required",
		})
		return
	}

	subcategory, err := h.categorySubcategoryUseCase.GetSubcategoryByCategoryId(c.Request.Context(), &CategoryID, limit, offset, &search)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Not found",
			"message": "Subcategory not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subcategory,
	})
}
func (h *CategorySubcategoryHandler) UpdateCategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is invalid",
			"message": "Category ID is required",
		})
		return
	}

	var request entities.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "Request body is invalid",
		})
		return
	}

	response, err := h.categorySubcategoryUseCase.UpdateCategory(c.Request.Context(), &categoryID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}
	if response.Success {
		c.JSON(http.StatusAccepted, gin.H{
			"success": response.Success,
			"message": response.Message,
		})
		return
	}
}

func (h *CategorySubcategoryHandler) UpdateSubcategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	subcategoryID := c.Param("id")
	if subcategoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad request",
			"message": "Subcategory ID is required",
		})
		return
	}

	var request entities.UpdateSubcategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := h.categorySubcategoryUseCase.UpdateSubcategory(c.Request.Context(), subcategoryID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
	})
}

// func (h *CategorySubcategoryHandler) GetCategoryById(c *gin.Context) {
// 	categoryID := c.Param("id")
// 	if categoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Category ID is invalid",
// 		})
// 		return
// 	}

// 	category, err := h.categorySubcategoryUseCase.GetCategoryById(c.Request.Context(), categoryID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"error":   "Not found",
// 			"message": "Some Internal Server Error Occured",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    category,
// 	})
// }

//

func (h *CategorySubcategoryHandler) DeleteCategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is invalid",
			"message": "Category ID is required",
		})
		return
	}

	response, _ := h.categorySubcategoryUseCase.DeleteCategory(c.Request.Context(), categoryID)
	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
	})
}
func (h *CategorySubcategoryHandler) DeleteSubcategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "operations" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	subcategoryID := c.Param("id")
	if subcategoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Sub-Category ID is invalid",
			"message": "Subcategory ID is required",
		})
		return
	}

	response, _ := h.categorySubcategoryUseCase.DeleteSubcategory(c.Request.Context(), subcategoryID)
	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
	})
}

func (h *CategorySubcategoryHandler) GetCategorySubCategoryForSpecificStore(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "customer" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	storeID := c.Query("storeId")
	warehouseId := c.Query("warehouseId")
	if warehouseId == "" || storeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Store or Warehouse ID is invalid",
			"message": "StoreID or Warehouse ID is required",
		})
		return
	}

	response, err := h.categorySubcategoryUseCase.GetCategorySubCategoryForSpecificStore(c.Request.Context(), storeID, warehouseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *CategorySubcategoryHandler) GetSubCategoryForStoreCategory(c *gin.Context) {
	role, isPresent := c.Get("role")
	if role != "customer" || !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unautherised User",
		})
		return
	}
	storeID := c.Query("storeId")
	warehouseId := c.Query("warehouseId")
	categoryId := c.Query("categoryId")
	if warehouseId == "" || storeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Store or Warehouse ID is invalid",
			"message": "StoreID or Warehouse ID is required",
		})
		return
	}
	if categoryId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Category ID is invalid",
			"message": "Category ID is required",
		})
		return
	}

	response, err := h.categorySubcategoryUseCase.GetSubCategoryForStoreCategory(c.Request.Context(), storeID, warehouseId, categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
