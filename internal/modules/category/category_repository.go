package category

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryFilter struct {
	Search string
	SortBy string
	Page   int
	Limit  int
}

type CategoryRepository interface {
	FindAll() ([]Category, error)
	FindAllPaginated(filter CategoryFilter) ([]Category, int64, error)
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

func (r *categoryRepository) FindAllPaginated(filter CategoryFilter) ([]Category, int64, error) {
	var categories []Category
	var total int64

	query := r.db.Model(&Category{})

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	orderBy := "created_at DESC"
	switch filter.SortBy {
	case "name_asc":
		orderBy = "name ASC"
	case "name_desc":
		orderBy = "name DESC"
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

	err = query.Order(orderBy).Offset((page - 1) * limit).Limit(limit).Find(&categories).Error

	return categories, total, err
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