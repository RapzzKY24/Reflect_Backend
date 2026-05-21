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

// @Summary      Checkout
// @Description  Convert cart items into an order, deduct stock, and clear cart
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CheckoutRequest true "Checkout Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=order.OrderResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /checkout [post]
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

// @Summary      Get my orders
// @Description  Retrieve all orders for the authenticated user
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]order.OrderResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /orders [get]
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

// @Summary      Get order by ID
// @Description  Retrieve a specific order by its ID for the authenticated user
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order ID"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=order.OrderResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /orders/{id} [get]
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

// @Summary      Get order by order number
// @Description  Retrieve a specific order by its order number for the authenticated user
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        orderNumber path string true "Order Number"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=order.OrderResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /orders/number/{orderNumber} [get]
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

// @Summary      Update order status (Admin)
// @Description  Update order status, payment status, and shipping status (admin only)
// @Tags         Admin Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order ID"
// @Param        request body UpdateOrderStatusRequest true "Update Order Status Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=order.OrderResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/orders/{id}/status [put]
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