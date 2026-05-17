package collection

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionRepository interface {
	FindAll() ([]Collection, error)
	FindActive() ([]Collection, error)
	FindFeatured() ([]Collection, error)
	FindByID(id uuid.UUID) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	Create(collection Collection) (Collection, error)
	Update(collection Collection) (Collection, error)
	Delete(id uuid.UUID) error
}

type collectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

func (r *collectionRepository) FindAll() ([]Collection, error) {
	var collections []Collection

	err := r.db.Order("created_at DESC").Find(&collections).Error

	return collections, err
}

func (r *collectionRepository) FindActive() ([]Collection, error) {
	var collections []Collection

	err := r.db.
		Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&collections).Error

	return collections, err
}

func (r *collectionRepository) FindFeatured() ([]Collection, error) {
	var collections []Collection

	err := r.db.
		Where("is_featured = ? AND is_active = ?", true, true).
		Order("created_at DESC").
		Find(&collections).Error

	return collections, err
}

func (r *collectionRepository) FindByID(id uuid.UUID) (Collection, error) {
	var collection Collection

	err := r.db.First(&collection, "id = ?", id).Error

	return collection, err
}

func (r *collectionRepository) FindBySlug(slug string) (Collection, error) {
	var collection Collection

	err := r.db.First(&collection, "slug = ?", slug).Error

	return collection, err
}

func (r *collectionRepository) Create(collection Collection) (Collection, error) {
	err := r.db.Create(&collection).Error

	return collection, err
}

func (r *collectionRepository) Update(collection Collection) (Collection, error) {
	err := r.db.Save(&collection).Error

	return collection, err
}

func (r *collectionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Collection{}, "id = ?", id).Error
}