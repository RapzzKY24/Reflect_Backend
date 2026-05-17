package order

import (
	"time"

	"reflect-backend/internal/modules/address"
	"reflect-backend/internal/modules/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderPaid       OrderStatus = "paid"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
)

type ShippingStatus string

const (
	ShippingPending   ShippingStatus = "pending"
	ShippingPacked    ShippingStatus = "packed"
	ShippingInTransit ShippingStatus = "in_transit"
	ShippingDelivered ShippingStatus = "delivered"
)

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OrderNumber string `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	AddressID uuid.UUID `gorm:"type:uuid;not null" json:"address_id"`

	Subtotal int `gorm:"not null" json:"subtotal"`
	Shipping int `gorm:"not null" json:"shipping"`
	Total    int `gorm:"not null" json:"total"`

	OrderStatus    OrderStatus    `gorm:"type:varchar(30);default:'pending'" json:"order_status"`
	PaymentStatus  PaymentStatus  `gorm:"type:varchar(30);default:'pending'" json:"payment_status"`
	ShippingStatus ShippingStatus `gorm:"type:varchar(30);default:'pending'" json:"shipping_status"`

	TrackingNumber string `gorm:"type:varchar(100)" json:"tracking_number"`

	User    user.User       `gorm:"foreignKey:UserID" json:"-"`
	Address address.Address `gorm:"foreignKey:AddressID" json:"address"`

	OrderItems []OrderItem     `gorm:"foreignKey:OrderID" json:"items"`
	Tracking   []OrderTracking `gorm:"foreignKey:OrderID" json:"tracking"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`

	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`

	ProductName  string `gorm:"type:varchar(200);not null" json:"product_name"`
	ProductSlug  string `gorm:"type:varchar(200)" json:"product_slug"`
	ProductImage string `gorm:"type:text" json:"product_image"`

	Price    int `gorm:"not null" json:"price"`
	Quantity int `gorm:"not null" json:"quantity"`
	Subtotal int `gorm:"not null" json:"subtotal"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderTracking struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`

	Status      string `gorm:"type:varchar(100);not null" json:"status"`
	Description string `gorm:"type:text" json:"description"`
	Location    string `gorm:"type:varchar(150)" json:"location"`

	CreatedAt time.Time `json:"created_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	return nil
}

func (o *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	return nil
}

func (o *OrderTracking) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	return nil
}