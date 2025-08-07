package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"net/http"
	"strconv"

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
	response, err := h.loginUseCase.LoginOperationalGuy(c.Request.Context(), &loginRequest)
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
		// Log the error for debugging
		println("Validation error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.RegisterOperationalGuy(c.Request.Context(), &registrationRequest)
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

// func (h *LoginHandler) LoginSeller(c *gin.Context) {
// 	var loginRequest entities.SellerLoginRequest
// 	if err := c.ShouldBindJSON(&loginRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": err.Error(),
// 		})
// 	}

// 	// Call the use case
// 	response, err := h.loginUseCase.LoginSeller(c.Request.Context(), loginRequest)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "An unexpected error occurred",
// 		})
// 	}

// 	// Return response based on success status
// 	if response.Success {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": response.Success,
// 			"message": response.Message,
// 			"token":   response.Token,
// 		})
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"success": response.Success,
// 			"error":   response.Error,
// 			"message": response.Message,
// 		})
// 	}
// }

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
	response, err := h.loginUseCase.RegisterSeller(c.Request.Context(), &registrationRequest)
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

func (h *LoginHandler) VerifyOTP(c *gin.Context) {
	phoneNumber := c.GetHeader("phonenumber")
	otpStr := c.GetHeader("otp")
	otp, err := strconv.ParseInt(otpStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalide Otp",
			"message": "Invalide OTP Format"})
		return
	}

	if len(phoneNumber) < 10 || otp < 100000 || otp > 999999 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Data Inconsistent",
			"message": "Api Data Invalid",
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.VerifyOTP(c.Request.Context(), &phoneNumber, &otp)
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
			"token":   response.Token,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": response.Success,
			"error":   response.Error,
			"message": response.Message,
		})
	}
}

// func (h *LoginHandler) AddBasicData(c *gin.Context) {
// 	var request *entities.AddBasicData
// 	err := c.ShouldBindJSON(&request)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": "Request Body is invalide"})
// 		return

// 	}

// 	// Call the use case
// 	response, err := h.loginUseCase.AddBasicData(c.Request.Context(), request)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "An unexpected error occurred",
// 		})
// 		return
// 	}

// 	// Return response based on success status
// 	if response.Success {
// 		c.JSON(http.StatusCreated, gin.H{
// 			"success": response.Success,
// 			"message": response.Message,
// 		})
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": response.Success,
// 			"error":   response.Error,
// 			"message": response.Message,
// 		})
// 	}
// }

// func (h *LoginHandler) LoginCustomer(c *gin.Context) {
// 	var loginRequest entities.CustomerLoginRequest
// 	if err := c.ShouldBindJSON(&loginRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	// Call the use case
// 	response, err := h.loginUseCase.LoginCustomer(c.Request.Context(), loginRequest)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "An unexpected error occurred",
// 		})
// 		return
// 	}

// 	// Return response based on success status
// 	if response.Success {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": response.Success,
// 			"message": response.Message,
// 			"token":   response.Token,
// 		})
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"success": response.Success,
// 			"error":   response.Error,
// 			"message": response.Message,
// 		})
// 	}
// }

// func (h *LoginHandler) RegisterCustomer(c *gin.Context) {
// 	var registrationRequest entities.CustomerRegistrationRequest
// 	if err := c.ShouldBindJSON(&registrationRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Validation error",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	// Call the use case
// 	response, err := h.loginUseCase.RegisterCustomer(c.Request.Context(), registrationRequest)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   "Internal server error",
// 			"message": "An unexpected error occurred",
// 		})
// 		return
// 	}

// 	// Return response based on success status
// 	if response.Success {
// 		c.JSON(http.StatusCreated, gin.H{
// 			"success": response.Success,
// 			"message": response.Message,
// 			"user_id": response.UserID,
// 		})
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": response.Success,
// 			"error":   response.Error,
// 			"message": response.Message,
// 		})
// 	}
// }
