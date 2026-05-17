package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id uuid.UUID) (User, error)
	FindByEmail(email string) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]User, error) {
	var users []User

	err := r.db.Order("created_at DESC").Find(&users).Error

	return users, err
}

func (r *userRepository) FindByID(id uuid.UUID) (User, error) {
	var user User

	err := r.db.First(&user, "id = ?", id).Error

	return user, err
}

func (r *userRepository) FindByEmail(email string) (User, error) {
	var user User

	err := r.db.First(&user, "email = ?", email).Error

	return user, err
}

func (r *userRepository) Create(user User) (User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *userRepository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	return user, err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}