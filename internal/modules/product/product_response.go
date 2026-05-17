package product

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Category     ProductCategory `json:"category"`
	Collection   string          `json:"collection"`
	Image        string          `json:"image"`
	Price        int             `json:"price"`
	Stock        int             `json:"stock"`
	StockStatus  StockStatus     `json:"stock_status"`
	IsNewArrival bool            `json:"is_new_arrival"`
	IsFeatured   bool            `json:"is_featured"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func ToProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Slug:         product.Slug,
		Description:  product.Description,
		Category:     product.Category,
		Collection:   product.Collection,
		Image:        product.Image,
		Price:        product.Price,
		Stock:        product.Stock,
		StockStatus:  product.StockStatus,
		IsNewArrival: product.IsNewArrival,
		IsFeatured:   product.IsFeatured,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func ToProductResponses(products []Product) []ProductResponse {
	responses := make([]ProductResponse, 0)

	for _, product := range products {
		responses = append(responses, ToProductResponse(product))
	}

	return responses
}