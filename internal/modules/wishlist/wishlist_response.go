package wishlist

import (
	"time"

	"github.com/google/uuid"
)

type WishlistItemResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Image        string    `json:"image"`
	Price        int       `json:"price"`
	Stock        int       `json:"stock"`
	StockStatus  string    `json:"stock_status"`
	IsNewArrival bool      `json:"is_new_arrival"`
	IsFeatured   bool      `json:"is_featured"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToWishlistItemResponse(item WishlistItem) WishlistItemResponse {
	return WishlistItemResponse{
		ID:           item.ID,
		ProductID:    item.ProductID,
		Name:         item.Product.Name,
		Slug:         item.Product.Slug,
		Image:        item.Product.Image,
		Price:        item.Product.Price,
		Stock:        item.Product.Stock,
		StockStatus:  string(item.Product.StockStatus),
		IsNewArrival: item.Product.IsNewArrival,
		IsFeatured:   item.Product.IsFeatured,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func ToWishlistResponses(items []WishlistItem) []WishlistItemResponse {
	responses := make([]WishlistItemResponse, 0)

	for _, item := range items {
		responses = append(responses, ToWishlistItemResponse(item))
	}

	return responses
}