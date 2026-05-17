package user

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]UserResponse, error)
	GetUserByID(id string) (UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetAllUsers() ([]UserResponse, error) {
	users, err := s.userRepository.FindAll()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get users")
	}

	return ToUserResponses(users), nil
}

func (s *userService) GetUserByID(id string) (UserResponse, error) {
	userID, err := uuid.Parse(id)

	if err != nil {
		return UserResponse{}, utils.BadRequest("Invalid user ID")
	}

	userData, err := s.userRepository.FindByID(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserResponse{}, utils.NotFound("User not found")
		}

		return UserResponse{}, utils.InternalServerError("Failed to get user")
	}

	return ToUserResponse(userData), nil
}

func (s *userService) DeleteUser(id string) error {
	userID, err := uuid.Parse(id)

	if err != nil {
		return utils.BadRequest("Invalid user ID")
	}

	_, err = s.userRepository.FindByID(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("User not found")
		}

		return utils.InternalServerError("Failed to get user")
	}

	err = s.userRepository.Delete(userID)

	if err != nil {
		return utils.InternalServerError("Failed to delete user")
	}

	return nil
}