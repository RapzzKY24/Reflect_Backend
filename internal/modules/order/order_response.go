package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderTrackingResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderItemResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductSlug  string    `json:"product_slug"`
	ProductImage string    `json:"product_image"`
	Price        int       `json:"price"`
	Quantity     int       `json:"quantity"`
	Subtotal     int       `json:"subtotal"`
}

type OrderResponse struct {
	ID uuid.UUID `json:"id"`

	OrderNumber string `json:"order_number"`

	Subtotal int `json:"subtotal"`
	Shipping int `json:"shipping"`
	Total    int `json:"total"`

	OrderStatus    OrderStatus    `json:"order_status"`
	PaymentStatus  PaymentStatus  `json:"payment_status"`
	ShippingStatus ShippingStatus `json:"shipping_status"`

	TrackingNumber string `json:"tracking_number"`

	Items    []OrderItemResponse     `json:"items"`
	Tracking []OrderTrackingResponse `json:"tracking"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToOrderResponse(order Order) OrderResponse {
	items := make([]OrderItemResponse, 0)

	for _, item := range order.OrderItems {
		items = append(items, OrderItemResponse{
			ID:           item.ID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductSlug:  item.ProductSlug,
			ProductImage: item.ProductImage,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Subtotal:     item.Subtotal,
		})
	}

	tracking := make([]OrderTrackingResponse, 0)

	for _, item := range order.Tracking {
		tracking = append(tracking, OrderTrackingResponse{
			ID:          item.ID,
			Status:      item.Status,
			Description: item.Description,
			Location:    item.Location,
			CreatedAt:   item.CreatedAt,
		})
	}

	return OrderResponse{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		Subtotal:       order.Subtotal,
		Shipping:       order.Shipping,
		Total:          order.Total,
		OrderStatus:    order.OrderStatus,
		PaymentStatus:  order.PaymentStatus,
		ShippingStatus: order.ShippingStatus,
		TrackingNumber: order.TrackingNumber,
		Items:          items,
		Tracking:       tracking,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}
}

func ToOrderResponses(orders []Order) []OrderResponse {
	responses := make([]OrderResponse, 0)

	for _, order := range orders {
		responses = append(responses, ToOrderResponse(order))
	}

	return responses
}