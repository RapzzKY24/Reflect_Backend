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

// @Summary      Get my cart
// @Description  Retrieve the authenticated user's shopping cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse{data=cart.CartResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /cart [get]
func (h *CartHandler) GetMyCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	cart, err := h.cartService.GetMyCart(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart retrieved successfully", cart)
}

// @Summary      Add item to cart
// @Description  Add a product to the authenticated user's cart. If the product already exists, quantities are added together.
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body AddToCartRequest true "Add to Cart Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=cart.CartResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /cart [post]
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

// @Summary      Update cart item quantity
// @Description  Update the quantity of a specific item in the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Cart Item ID"
// @Param        request body UpdateCartItemRequest true "Update Cart Item Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=cart.CartResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /cart/{id} [put]
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

// @Summary      Remove cart item
// @Description  Remove a specific item from the cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Cart Item ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /cart/{id} [delete]
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

// @Summary      Clear cart
// @Description  Remove all items from the authenticated user's cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	err := h.cartService.ClearCart(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart cleared successfully", nil)
}