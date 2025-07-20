package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUseCase *usecase.LoginUseCaseInterface
}

func NewLoginHandler(loginUseCase *usecase.LoginUseCaseInterface) *LoginHandler {
	return &LoginHandler{loginUseCase: loginUseCase}
}

func (h *LoginHandler) LoginOperationalGuy(c *gin.Context) {
	var loginRequest entities.OperationalGuyLoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.LoginOperationalGuy(c.Request.Context(), loginRequest)
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
		c.JSON(http.StatusOK, gin.H{
			"success": response.Success,
			"message": response.Message,
			"token":   response.Token,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

func (h *LoginHandler) RegisterOperationalGuy(c *gin.Context) {
	var registrationRequest entities.OperationalGuyRegistrationRequest
	if err := c.ShouldBindJSON(&registrationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.RegisterOperationalGuy(c.Request.Context(), registrationRequest)
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
			"user_id": response.UserID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

func (h *LoginHandler) LoginSeller(c *gin.Context) {
	var loginRequest entities.SellerLoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
	}

	// Call the use case
	response, err := h.loginUseCase.LoginSeller(c.Request.Context(), loginRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"message": "An unexpected error occurred",
		})
	}

	// Return response based on success status
	if response.Success {
		c.JSON(http.StatusOK, gin.H{
			"success": response.Success,
			"message": response.Message,
			"token":   response.Token,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

func (h *LoginHandler) RegisterSeller(c *gin.Context) {
	var registrationRequest entities.SellerRegistrationRequest
	if err := c.ShouldBindJSON(&registrationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.RegisterSeller(c.Request.Context(), registrationRequest)
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
			"user_id": response.UserID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

func (h *LoginHandler) LoginCustomer(c *gin.Context) {
	var loginRequest entities.CustomerLoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.LoginCustomer(c.Request.Context(), loginRequest)
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
		c.JSON(http.StatusOK, gin.H{
			"success": response.Success,
			"message": response.Message,
			"token":   response.Token,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

func (h *LoginHandler) RegisterCustomer(c *gin.Context) {
	var registrationRequest entities.CustomerRegistrationRequest
	if err := c.ShouldBindJSON(&registrationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.RegisterCustomer(c.Request.Context(), registrationRequest)
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
			"user_id": response.UserID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}
