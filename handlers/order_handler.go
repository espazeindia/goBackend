package handlers

import (
	"espazeBackend/domain/entities"
	"espazeBackend/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	OrderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{OrderUsecase: orderUsecase}

}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	var requestData entities.GetAllOrdersRequest
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.OrderUsecase.GetAllOrders(c.Request.Context(), &requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var requestOrder entities.CreateOrderRequest
	err := c.ShouldBindJSON(&requestOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Print("request order", requestOrder)
	err = h.OrderUsecase.CreateNewOrder(c.Request.Context(), &requestOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "New order created!"})
}

func (h *OrderHandler) GetOrderByOrderID(c *gin.Context) {
	orderId := c.Query("orderId")
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order Id is empty"})
		return
	}

	order, err := h.OrderUsecase.GetOrderByOrderID(c.Request.Context(), &orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, order)
}

func (h *OrderHandler) GetOrderByUserID(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Id is empty"})
		return
	}

	order, err := h.OrderUsecase.GetOrderByUserID(c.Request.Context(), &userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, order)
}

func (h *OrderHandler) GetOrderBySellerID(c *gin.Context) {
	sellerId := c.Query("sellerId")
	if sellerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seller Id is empty"})
		return
	}

	order, err := h.OrderUsecase.GetOrderBySellerID(c.Request.Context(), &sellerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, order)
}
