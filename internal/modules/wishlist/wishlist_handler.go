package wishlist

import (
	"net/http"

	"reflect-backend/internal/middleware"
	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	wishlistService WishlistService
}

func NewWishlistHandler(wishlistService WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// @Summary      Get my wishlist
// @Description  Retrieve the authenticated user's wishlist
// @Tags         Wishlist
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]wishlist.WishlistItemResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /wishlist [get]
func (h *WishlistHandler) GetMyWishlist(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	items, err := h.wishlistService.GetMyWishlist(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Wishlist retrieved successfully",
		items,
	)
}

// @Summary      Add item to wishlist
// @Description  Add a product to the authenticated user's wishlist
// @Tags         Wishlist
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body AddWishlistRequest true "Add to Wishlist Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=[]wishlist.WishlistItemResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /wishlist [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var request AddWishlistRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	items, err := h.wishlistService.AddToWishlist(userID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"Item added to wishlist successfully",
		items,
	)
}

// @Summary      Remove wishlist item
// @Description  Remove a specific item from the wishlist by its ID
// @Tags         Wishlist
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Wishlist Item ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /wishlist/{id} [delete]
func (h *WishlistHandler) RemoveWishlistItem(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	wishlistItemID := c.Param("id")

	err := h.wishlistService.RemoveWishlistItem(userID, wishlistItemID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Wishlist item removed successfully",
		nil,
	)
}

// @Summary      Remove wishlist item by product
// @Description  Remove an item from the wishlist by product ID
// @Tags         Wishlist
// @Produce      json
// @Security     BearerAuth
// @Param        productId path string true "Product ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /wishlist/product/{productId} [delete]
func (h *WishlistHandler) RemoveWishlistByProduct(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	productID := c.Param("productId")

	err := h.wishlistService.RemoveWishlistByProduct(userID, productID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Wishlist item removed successfully",
		nil,
	)
}