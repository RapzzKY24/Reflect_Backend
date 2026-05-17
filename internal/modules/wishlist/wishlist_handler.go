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