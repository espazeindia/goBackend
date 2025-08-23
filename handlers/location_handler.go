package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationUseCase *usecase.LocationUseCase
}

func NewLocationHandler(locationUseCase *usecase.LocationUseCase) *LocationHandler {
	return &LocationHandler{
		locationUseCase: locationUseCase,
	}
}

func (h *LocationHandler) GetLocationForUserID(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User Id parameter is required", "error": "User Id is empty"})
		return
	}

	response, err := h.locationUseCase.GetLocationForUserID(c, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": response.Success, "message": response.Message, "error": response.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": response.Success, "data": response.Data})
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var location entities.Location
	err := c.ShouldBindJSON(&location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.locationUseCase.CreateLocation(&location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, location)
}

func (h *LocationHandler) GetLocationByAddress(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	location, err := h.locationUseCase.GetLocationByAddress(address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	c.JSON(http.StatusOK, location)
}
