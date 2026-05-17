package cart

import (
	"reflect-backend/internal/modules/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	FindByUserID(userID uuid.UUID) ([]CartItem, error)
	FindItemByID(id uuid.UUID) (CartItem, error)
	FindItemByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (CartItem, error)
	FindProductForUpdate(tx *gorm.DB, productID uuid.UUID) (product.Product, error)
	CreateWithTx(tx *gorm.DB, item CartItem) (CartItem, error)
	UpdateWithTx(tx *gorm.DB, item CartItem) (CartItem, error)
	DeleteByID(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	Transaction(fn func(tx *gorm.DB) error) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(userID uuid.UUID) ([]CartItem, error) {
	var items []CartItem

	err := r.db.
		Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error

	return items, err
}

func (r *cartRepository) FindItemByID(id uuid.UUID) (CartItem, error) {
	var item CartItem

	err := r.db.
		Preload("Product").
		First(&item, "id = ?", id).Error

	return item, err
}

func (r *cartRepository) FindItemByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (CartItem, error) {
	var item CartItem

	err := r.db.
		Preload("Product").
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&item).Error

	return item, err
}

func (r *cartRepository) FindProductForUpdate(tx *gorm.DB, productID uuid.UUID) (product.Product, error) {
	var productData product.Product

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&productData, "id = ?", productID).Error

	return productData, err
}

func (r *cartRepository) CreateWithTx(tx *gorm.DB, item CartItem) (CartItem, error) {
	err := tx.Create(&item).Error
	return item, err
}

func (r *cartRepository) UpdateWithTx(tx *gorm.DB, item CartItem) (CartItem, error) {
	err := tx.Save(&item).Error
	return item, err
}

func (r *cartRepository) DeleteByID(id uuid.UUID) error {
	return r.db.Delete(&CartItem{}, "id = ?", id).Error
}

func (r *cartRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&CartItem{}, "user_id = ?", userID).Error
}

func (r *cartRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}