package order

import (
	"reflect-backend/internal/modules/cart"
	"reflect-backend/internal/modules/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	FindByUserID(userID uuid.UUID) ([]Order, error)
	FindByID(id uuid.UUID) (Order, error)
	FindByOrderNumber(orderNumber string) (Order, error)
	FindCartItems(userID uuid.UUID) ([]cart.CartItem, error)
	FindProductForUpdate(tx *gorm.DB, productID uuid.UUID) (product.Product, error)
	CreateOrderWithTx(tx *gorm.DB, order Order) (Order, error)
	CreateOrderItemWithTx(tx *gorm.DB, item OrderItem) error
	CreateTrackingWithTx(tx *gorm.DB, tracking OrderTracking) error
	UpdateProductStockWithTx(tx *gorm.DB, product product.Product) error
	ClearCartWithTx(tx *gorm.DB, userID uuid.UUID) error
	UpdateOrder(order Order) (Order, error)
	Transaction(fn func(tx *gorm.DB) error) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindByUserID(userID uuid.UUID) ([]Order, error) {
	var orders []Order

	err := r.db.
		Preload("OrderItems").
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *orderRepository) FindByID(id uuid.UUID) (Order, error) {
	var order Order

	err := r.db.
		Preload("Address").
		Preload("OrderItems").
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&order, "id = ?", id).Error

	return order, err
}

func (r *orderRepository) FindByOrderNumber(orderNumber string) (Order, error) {
	var order Order

	err := r.db.
		Preload("Address").
		Preload("OrderItems").
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&order, "order_number = ?", orderNumber).Error

	return order, err
}

func (r *orderRepository) FindCartItems(userID uuid.UUID) ([]cart.CartItem, error) {
	var items []cart.CartItem

	err := r.db.
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&items).Error

	return items, err
}

func (r *orderRepository) FindProductForUpdate(tx *gorm.DB, productID uuid.UUID) (product.Product, error) {
	var productData product.Product

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&productData, "id = ?", productID).Error

	return productData, err
}

func (r *orderRepository) CreateOrderWithTx(tx *gorm.DB, order Order) (Order, error) {
	err := tx.Create(&order).Error
	return order, err
}

func (r *orderRepository) CreateOrderItemWithTx(tx *gorm.DB, item OrderItem) error {
	return tx.Create(&item).Error
}

func (r *orderRepository) CreateTrackingWithTx(tx *gorm.DB, tracking OrderTracking) error {
	return tx.Create(&tracking).Error
}

func (r *orderRepository) UpdateProductStockWithTx(tx *gorm.DB, product product.Product) error {
	return tx.Save(&product).Error
}

func (r *orderRepository) ClearCartWithTx(tx *gorm.DB, userID uuid.UUID) error {
	return tx.Delete(&cart.CartItem{}, "user_id = ?", userID).Error
}

func (r *orderRepository) UpdateOrder(order Order) (Order, error) {
	err := r.db.Save(&order).Error
	return order, err
}

func (r *orderRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}