package handlers

import (
	"net/http"
	"strconv"

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
func (h *CategorySubcategoryHandler) GetAllCategories(c *gin.Context) {
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

	categories, err := h.categorySubcategoryUseCase.GetAllCategories(c.Request.Context(), limit, offset, &search)
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
	subcategories, err := h.categorySubcategoryUseCase.GetAllSubcategories(c.Request.Context(), limit, offset, &search)
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

// func (h *CategorySubcategoryHandler) CreateCategory(c *gin.Context) {
// 	var request entities.CreateCategoryRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": "Request body is invalid",
// 		})
// 		return
// 	}

// 	category := &entities.Category{
// 		CategoryName:  request.CategoryName,
// 		CategoryImage: request.CategoryImage,
// 	}

// 	err := h.categorySubcategoryUseCase.CreateCategory(c.Request.Context(), category)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "Some Internal Server Error Occured",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"message": "Category created successfully",
// 	})
// }

// func (h *CategorySubcategoryHandler) UpdateCategory(c *gin.Context) {
// 	categoryID := c.Param("id")
// 	if categoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Category ID is required",
// 		})
// 		return
// 	}

// 	var request entities.UpdateCategoryRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": "request body is invalid",
// 		})
// 		return
// 	}

// 	category := &entities.Category{
// 		CategoryID:    categoryID,
// 		CategoryName:  request.CategoryName,
// 		CategoryImage: request.CategoryImage,
// 	}

// 	err := h.categorySubcategoryUseCase.UpdateCategory(c.Request.Context(), category)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "Some Internal Server Error Occured",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Category updated successfully",
// 	})
// }

// func (h *CategorySubcategoryHandler) DeleteCategory(c *gin.Context) {
// 	categoryID := c.Param("id")
// 	if categoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Category ID is required",
// 		})
// 		return
// 	}

// 	err := h.categorySubcategoryUseCase.DeleteCategory(c.Request.Context(), categoryID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "Some Internal Server Error Occured",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Category deleted successfully",
// 	})
// }

// // Enhanced category handlers with subcategories
// func (h *CategorySubcategoryHandler) GetAllCategoriesWithSubcategories(c *gin.Context) {
// 	categories, err := h.categorySubcategoryUseCase.GetAllCategoriesWithSubcategories(c.Request.Context())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "Some Internal Server Error Occured",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    categories,
// 	})
// }

// func (h *CategorySubcategoryHandler) GetCategoryWithSubcategories(c *gin.Context) {
// 	categoryID := c.Param("id")
// 	if categoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Category ID is required",
// 		})
// 		return
// 	}

// 	category, err := h.categorySubcategoryUseCase.GetCategoryWithSubcategories(c.Request.Context(), categoryID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"error":   "Not found",
// 			"message": "Category not found",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    category,
// 	})
// }

// func (h *CategorySubcategoryHandler) GetSubcategoryById(c *gin.Context) {
// 	subcategoryID := c.Param("id")
// 	if subcategoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Subcategory ID is required",
// 		})
// 		return
// 	}

// 	subcategory, err := h.categorySubcategoryUseCase.GetSubcategoryById(c.Request.Context(), subcategoryID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"success": false,
// 			"error":   "Not found",
// 			"message": "Subcategory not found",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    subcategory,
// 	})
// }

// func (h *CategorySubcategoryHandler) CreateSubcategory(c *gin.Context) {
// 	var request entities.CreateSubcategoryRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	subcategory := &entities.Subcategory{
// 		SubcategoryName:  request.SubcategoryName,
// 		SubcategoryImage: request.SubcategoryImage,
// 		CategoryID:       request.CategoryID,
// 	}

// 	err := h.categorySubcategoryUseCase.CreateSubcategory(c.Request.Context(), subcategory)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"data":    subcategory,
// 		"message": "Subcategory created successfully",
// 	})
// }

// func (h *CategorySubcategoryHandler) UpdateSubcategory(c *gin.Context) {
// 	subcategoryID := c.Param("id")
// 	if subcategoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Subcategory ID is required",
// 		})
// 		return
// 	}

// 	var request entities.UpdateSubcategoryRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	subcategory := &entities.Subcategory{
// 		SubcategoryID:    subcategoryID,
// 		SubcategoryName:  request.SubcategoryName,
// 		SubcategoryImage: request.SubcategoryImage,
// 		CategoryID:       request.CategoryID,
// 	}

// 	err := h.categorySubcategoryUseCase.UpdateSubcategory(c.Request.Context(), subcategory)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    subcategory,
// 		"message": "Subcategory updated successfully",
// 	})
// }

// func (h *CategorySubcategoryHandler) DeleteSubcategory(c *gin.Context) {
// 	subcategoryID := c.Param("id")
// 	if subcategoryID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Bad request",
// 			"message": "Subcategory ID is required",
// 		})
// 		return
// 	}

// 	err := h.categorySubcategoryUseCase.DeleteSubcategory(c.Request.Context(), subcategoryID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Subcategory deleted successfully",
// 	})
// }
