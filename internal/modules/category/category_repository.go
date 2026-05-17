package category

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]Category, error)
	FindActive() ([]Category, error)
	FindByID(id uuid.UUID) (Category, error)
	FindBySlug(slug string) (Category, error)
	Create(category Category) (Category, error)
	Update(category Category) (Category, error)
	Delete(id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]Category, error) {
	var categories []Category

	err := r.db.Order("created_at DESC").Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) FindActive() ([]Category, error) {
	var categories []Category

	err := r.db.
		Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) FindByID(id uuid.UUID) (Category, error) {
	var category Category

	err := r.db.First(&category, "id = ?", id).Error

	return category, err
}

func (r *categoryRepository) FindBySlug(slug string) (Category, error) {
	var category Category

	err := r.db.First(&category, "slug = ?", slug).Error

	return category, err
}

func (r *categoryRepository) Create(category Category) (Category, error) {
	err := r.db.Create(&category).Error

	return category, err
}

func (r *categoryRepository) Update(category Category) (Category, error) {
	err := r.db.Save(&category).Error

	return category, err
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Category{}, "id = ?", id).Error
}