package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id uuid.UUID) (Product, error)
	FindBySlug(slug string) (Product, error)
	FindNewArrivals() ([]Product, error)
	FindFeatured() ([]Product, error)
	FindRelatedProducts(category string, excludeID uuid.UUID, limit int) ([]Product, error)
	FindByCategory(category string) ([]Product, error)
	Create(product Product) (Product, error)
	Update(product Product) (Product, error)
	Delete(id uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]Product, error) {
	var products []Product

	err := r.db.Order("created_at DESC").Find(&products).Error

	return products, err
}

func (r *productRepository) FindByID(id uuid.UUID) (Product, error) {
	var product Product

	err := r.db.First(&product, "id = ?", id).Error

	return product, err
}

func (r *productRepository) FindRelatedProducts(
	category string,
	excludeID uuid.UUID,
	limit int,
) ([]Product, error) {
	var products []Product

	err := r.db.
		Where("category = ? AND id != ?", category, excludeID).
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error

	return products, err
}

func (r *productRepository) FindBySlug(slug string) (Product, error) {
	var product Product

	err := r.db.First(&product, "slug = ?", slug).Error

	return product, err
}

func (r *productRepository) FindNewArrivals() ([]Product, error) {
	var products []Product

	err := r.db.
		Where("is_new_arrival = ?", true).
		Order("created_at DESC").
		Find(&products).Error

	return products, err
}

func (r *productRepository) FindFeatured() ([]Product, error) {
	var products []Product

	err := r.db.
		Where("is_featured = ?", true).
		Order("created_at DESC").
		Find(&products).Error

	return products, err
}

func (r *productRepository) FindByCategory(category string) ([]Product, error) {
	var products []Product

	err := r.db.
		Where("category = ?", category).
		Order("created_at DESC").
		Find(&products).Error

	return products, err
}

func (r *productRepository) Create(product Product) (Product, error) {
	err := r.db.Create(&product).Error

	return product, err
}

func (r *productRepository) Update(product Product) (Product, error) {
	err := r.db.Save(&product).Error

	return product, err
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Product{}, "id = ?", id).Error
}