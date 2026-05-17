package collection

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Collection struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(120);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(150);uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Image       string         `gorm:"type:text" json:"image"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Collection) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	return nil
}