package cart

import (
	"time"

	"github.com/google/uuid"
)

type CartItemResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
	Price       int       `json:"price"`
	Quantity    int       `json:"quantity"`
	Subtotal    int       `json:"subtotal"`
	Stock       int       `json:"stock"`
	StockStatus string    `json:"stock_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalItems int                `json:"total_items"`
	TotalPrice int                `json:"total_price"`
}

func ToCartItemResponse(item CartItem) CartItemResponse {
	return CartItemResponse{
		ID:          item.ID,
		ProductID:   item.ProductID,
		Name:        item.Product.Name,
		Slug:        item.Product.Slug,
		Image:       item.Product.Image,
		Price:       item.Product.Price,
		Quantity:    item.Quantity,
		Subtotal:    item.Product.Price * item.Quantity,
		Stock:       item.Product.Stock,
		StockStatus: string(item.Product.StockStatus),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ToCartResponse(items []CartItem) CartResponse {
	responses := make([]CartItemResponse, 0)
	totalItems := 0
	totalPrice := 0

	for _, item := range items {
		itemResponse := ToCartItemResponse(item)

		responses = append(responses, itemResponse)
		totalItems += item.Quantity
		totalPrice += itemResponse.Subtotal
	}

	return CartResponse{
		Items:      responses,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}
}