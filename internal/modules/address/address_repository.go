package address

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository interface {
	FindByUserID(userID uuid.UUID) ([]Address, error)
	FindByID(id uuid.UUID) (Address, error)
	Create(address Address) (Address, error)
	Update(address Address) (Address, error)
	Delete(id uuid.UUID) error
	ClearPrimaryAddresses(userID uuid.UUID) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) FindByUserID(userID uuid.UUID) ([]Address, error) {
	var addresses []Address

	err := r.db.
		Where("user_id = ?", userID).
		Order("is_primary DESC").
		Order("created_at DESC").
		Find(&addresses).Error

	return addresses, err
}

func (r *addressRepository) FindByID(id uuid.UUID) (Address, error) {
	var address Address

	err := r.db.First(&address, "id = ?", id).Error

	return address, err
}

func (r *addressRepository) Create(address Address) (Address, error) {
	err := r.db.Create(&address).Error

	return address, err
}

func (r *addressRepository) Update(address Address) (Address, error) {
	err := r.db.Save(&address).Error

	return address, err
}

func (r *addressRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Address{}, "id = ?", id).Error
}

func (r *addressRepository) ClearPrimaryAddresses(userID uuid.UUID) error {
	return r.db.
		Model(&Address{}).
		Where("user_id = ?", userID).
		Update("is_primary", false).
		Error
}