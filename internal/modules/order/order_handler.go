package order

import (
	"net/http"

	"reflect-backend/internal/middleware"
	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var request CheckoutRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	order, err := h.orderService.Checkout(userID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"Checkout successful",
		order,
	)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	orders, err := h.orderService.GetMyOrders(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Orders retrieved successfully",
		orders,
	)
}

func (h *OrderHandler) GetMyOrderByID(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	orderID := c.Param("id")

	order, err := h.orderService.GetMyOrderByID(userID, orderID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Order retrieved successfully",
		order,
	)
}

func (h *OrderHandler) GetMyOrderByNumber(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	orderNumber := c.Param("orderNumber")

	order, err := h.orderService.GetMyOrderByNumber(userID, orderNumber)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Order retrieved successfully",
		order,
	)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var request UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	order, err := h.orderService.UpdateOrderStatus(orderID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Order status updated successfully",
		order,
	)
}