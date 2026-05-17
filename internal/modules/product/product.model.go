package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory string
type StockStatus string

const (
	CategoryHoodie      ProductCategory = "hoodie"
	CategoryTShirt      ProductCategory = "t-shirt"
	CategoryOuterwear   ProductCategory = "outerwear"
	CategoryAccessories ProductCategory = "accessories"
	CategoryShoes       ProductCategory = "shoes"
)

const (
	StockInStock    StockStatus = "in_stock"
	StockLowStock   StockStatus = "low_stock"
	StockOutOfStock StockStatus = "out_of_stock"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string          `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string          `gorm:"type:varchar(180);uniqueIndex;not null" json:"slug"`
	Description string          `gorm:"type:text" json:"description"`
	Category    ProductCategory `gorm:"type:varchar(50);not null" json:"category"`
	Collection  string          `gorm:"type:varchar(100)" json:"collection"`
	Image       string          `gorm:"type:text" json:"image"`
	Price       int             `gorm:"not null" json:"price"`
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	StockStatus StockStatus     `gorm:"type:varchar(30);default:'in_stock'" json:"stock_status"`
	IsNewArrival bool            `gorm:"default:false" json:"is_new_arrival"`
	IsFeatured   bool            `gorm:"default:false" json:"is_featured"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return nil
}