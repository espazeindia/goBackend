package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OnboardingHandler struct {
	OnboardingUseCase *usecase.OnboardingUseCaseInterface
}

func NewOnboardingHandler(OnboardingUseCase *usecase.OnboardingUseCaseInterface) *OnboardingHandler {
	return &OnboardingHandler{OnboardingUseCase: OnboardingUseCase}
}

func (h *OnboardingHandler) AddBasicDetail(c *gin.Context) {
	var request *entities.SellerBasicDetail
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "Request Body is invalide"})
		return

	}

	sellerId, isPresent := c.Get("user_id")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "User ID not present in token"})
		return
	}

	sellerIdString, ok := sellerId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "User ID not present in token"})
		return
	}

	// Call the use case
	response, err := h.OnboardingUseCase.AddBasicDetail(c.Request.Context(), request, sellerIdString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	// Return response based on success status
	c.JSON(http.StatusCreated, gin.H{
		"success": response.Success,
		"message": response.Message,
		"token":   response.Token,
	})

}

func (h *OnboardingHandler) GetBasicDetail(c *gin.Context) {
	userId, isPresent := c.Get("user_id")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID doesn't exist",
			"message": "User ID doesn't exist in token",
		})
		return
	}

	userIdString, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is not string",
			"message": "User ID is not a valid string in token",
		})
		return
	}

	response, err := h.OnboardingUseCase.GetBasicDetail(c, userIdString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"data":    response.Data,
	})
}

func (h *OnboardingHandler) OnboardingAdmin(c *gin.Context) {
	var request *entities.AdminOnboaring
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "Request Body is invalide"})
		return
	}

	response, err := h.OnboardingUseCase.OnboardingAdmin(c.Request.Context(), request)
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
		"token":   response.Token,
	})

}

func (h *OnboardingHandler) RegisterOperationalGuy(c *gin.Context) {
	var registrationRequest entities.OperationalGuyRegistrationRequest

	if err := c.ShouldBindJSON(&registrationRequest); err != nil {
		// Log the error for debugging
		println("Validation error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
			"message": "Request body is invalid",
		})
		return
	}

	// Call the use case
	response, err := h.OnboardingUseCase.RegisterOperationalGuy(c.Request.Context(), &registrationRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "An unexpected error occurred",
		})
		return
	}

	// Return response based on success status
	if response.Success {
		c.JSON(http.StatusCreated, gin.H{
			"success": response.Success,
			"message": response.Message,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

func (h *OnboardingHandler) GetOperationalGuy(c *gin.Context) {
	userId, isPresent := c.Get("user_id")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID doesn't exist",
			"message": "User ID doesn't exist in token",
		})
		return
	}

	userIdString, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is not string",
			"message": "User ID is not a valid string in token",
		})
		return
	}

	response, err := h.OnboardingUseCase.GetOperationalGuy(c, userIdString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "Some Internal Server Error Occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
		"token":   response.Token,
	})
}

func (h *OnboardingHandler) EditOperationalGuy(c *gin.Context) {
	var request *entities.OperationalGuyGetRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "Request Body is invalide"})
		return

	}

	operationsId, isPresent := c.Get("user_id")
	if !isPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "User ID not present in token"})
		return
	}

	operationsIdString, ok := operationsId.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": "User ID not present in token"})
		return
	}

	// Call the use case
	response, err := h.OnboardingUseCase.EditOperationalGuy(c.Request.Context(), request, operationsIdString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
		return
	}

	// Return response based on success status
	c.JSON(http.StatusCreated, gin.H{
		"success": response.Success,
		"message": response.Message,
		"token":   response.Token,
	})

}
