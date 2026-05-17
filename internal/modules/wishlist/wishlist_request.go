package wishlist

type AddWishlistRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}