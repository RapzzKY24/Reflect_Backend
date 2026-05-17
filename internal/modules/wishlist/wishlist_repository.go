package wishlist

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistRepository interface {
	FindByUserID(userID uuid.UUID) ([]WishlistItem, error)
	FindByID(id uuid.UUID) (WishlistItem, error)
	FindByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (WishlistItem, error)
	Create(item WishlistItem) (WishlistItem, error)
	Delete(id uuid.UUID) error
	DeleteByUserAndProduct(userID uuid.UUID, productID uuid.UUID) error
}

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) FindByUserID(userID uuid.UUID) ([]WishlistItem, error) {
	var items []WishlistItem

	err := r.db.
		Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error

	return items, err
}

func (r *wishlistRepository) FindByID(id uuid.UUID) (WishlistItem, error) {
	var item WishlistItem

	err := r.db.
		Preload("Product").
		First(&item, "id = ?", id).Error

	return item, err
}

func (r *wishlistRepository) FindByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (WishlistItem, error) {
	var item WishlistItem

	err := r.db.
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&item).Error

	return item, err
}

func (r *wishlistRepository) Create(item WishlistItem) (WishlistItem, error) {
	err := r.db.Create(&item).Error

	return item, err
}

func (r *wishlistRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&WishlistItem{}, "id = ?", id).Error
}

func (r *wishlistRepository) DeleteByUserAndProduct(userID uuid.UUID, productID uuid.UUID) error {
	return r.db.
		Delete(&WishlistItem{}, "user_id = ? AND product_id = ?", userID, productID).
		Error
}