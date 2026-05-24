package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductFilter struct {
	Search   string
	Category string
	SortBy   string
	Page     int
	Limit    int
}

type ProductRepository interface {
	FindAll() ([]Product, error)
	FindAllPaginated(filter ProductFilter) ([]Product, int64, error)
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

func (r *productRepository) FindAllPaginated(filter ProductFilter) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	orderBy := "created_at DESC"
	switch filter.SortBy {
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	case "oldest":
		orderBy = "created_at ASC"
	case "newest":
		orderBy = "created_at DESC"
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}

	limit := filter.Limit
	if limit < 1 {
		limit = 12
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order(orderBy).Offset((page - 1) * limit).Limit(limit).Find(&products).Error

	return products, total, err
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