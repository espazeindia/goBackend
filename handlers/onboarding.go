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

	// Call the use case
	response, err := h.OnboardingUseCase.AddBasicDetail(c.Request.Context(), request)
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
