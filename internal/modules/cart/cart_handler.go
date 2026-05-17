package cart

import (
	"net/http"

	"reflect-backend/internal/middleware"
	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService CartService
}

func NewCartHandler(cartService CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetMyCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	cart, err := h.cartService.GetMyCart(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart retrieved successfully", cart)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var request AddToCartRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	cart, err := h.cartService.AddToCart(userID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Item added to cart successfully", cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	cartItemID := c.Param("id")

	var request UpdateCartItemRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	cart, err := h.cartService.UpdateCartItem(userID, cartItemID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart item updated successfully", cart)
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	cartItemID := c.Param("id")

	err := h.cartService.RemoveCartItem(userID, cartItemID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart item removed successfully", nil)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	err := h.cartService.ClearCart(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart cleared successfully", nil)
}