package address

import (
	"time"

	"reflect-backend/internal/modules/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	RecipientName string `gorm:"type:varchar(120);not null" json:"recipient_name"`
	PhoneNumber   string `gorm:"type:varchar(30);not null" json:"phone_number"`

	Province  string `gorm:"type:varchar(100);not null" json:"province"`
	City      string `gorm:"type:varchar(100);not null" json:"city"`
	District  string `gorm:"type:varchar(100);not null" json:"district"`
	PostalCode string `gorm:"type:varchar(20);not null" json:"postal_code"`

	FullAddress string `gorm:"type:text;not null" json:"full_address"`

	Label string `gorm:"type:varchar(50)" json:"label"`

	IsPrimary bool `gorm:"default:false" json:"is_primary"`

	User user.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	return nil
}