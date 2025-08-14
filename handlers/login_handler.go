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
			"error":   "Invalid request body",
			"message": "Request body is invalid",
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
			"error":   "Invalid request body",
			"message": "Request body is invalid",
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

func (h *LoginHandler) GetOTP(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")

	if len(phoneNumber) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Phone Number",
			"message": "Entered phone Number is less than 10"})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.GetOTP(c.Request.Context(), phoneNumber)
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
	})

}

func (h *LoginHandler) VerifyOTP(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")
	otpStr := c.Query("otp")
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

func (h *LoginHandler) VerifyPin(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")
	pinStr := c.Query("pin")
	pin, err := strconv.ParseInt(pinStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalide Pin",
			"message": "Invalide Pin Format"})
		return
	}

	if len(phoneNumber) < 10 || pin < 100000 || pin > 999999 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Data Inconsistent",
			"message": "Api Data Invalid",
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.VerifyPin(c.Request.Context(), &phoneNumber, &pin)
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
func (h *LoginHandler) GetOTPForCustomer(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")

	if len(phoneNumber) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Phone Number",
			"message": "Entered phone Number is less than 10"})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.GetOTPForCustomer(c.Request.Context(), phoneNumber)
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
	})

}

func (h *LoginHandler) VerifyOTPForCustomer(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")
	otpStr := c.Query("otp")
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
	response, err := h.loginUseCase.VerifyOTPForCustomer(c.Request.Context(), &phoneNumber, &otp)
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

func (h *LoginHandler) VerifyPinForCustomer(c *gin.Context) {
	phoneNumber := c.Query("phonenumber")
	pinStr := c.Query("pin")
	pin, err := strconv.ParseInt(pinStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalide Pin",
			"message": "Invalide Pin Format"})
		return
	}

	if len(phoneNumber) < 10 || pin < 100000 || pin > 999999 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Data Inconsistent",
			"message": "Api Data Invalid",
		})
		return
	}

	// Call the use case
	response, err := h.loginUseCase.VerifyPinForCustomer(c.Request.Context(), &phoneNumber, &pin)
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
