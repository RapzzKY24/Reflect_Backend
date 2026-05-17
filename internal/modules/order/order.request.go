package order

type CheckoutRequest struct {
	AddressID string `json:"address_id" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	OrderStatus    OrderStatus    `json:"order_status" binding:"required"`
	PaymentStatus  PaymentStatus  `json:"payment_status"`
	ShippingStatus ShippingStatus `json:"shipping_status"`
	TrackingStatus string         `json:"tracking_status"`
	Description    string         `json:"description"`
	Location       string         `json:"location"`
}