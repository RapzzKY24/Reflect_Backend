package profile

import (
	"errors"

	"reflect-backend/internal/modules/user"
	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService interface {
	GetMyProfile(userID string) (ProfileResponse, error)
	UpdateMyProfile(userID string, request UpdateProfileRequest) (ProfileResponse, error)
}

type profileService struct {
	userRepository user.UserRepository
}

func NewProfileService(
	userRepository user.UserRepository,
) ProfileService {
	return &profileService{
		userRepository: userRepository,
	}
}

func (s *profileService) GetMyProfile(userID string) (ProfileResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return ProfileResponse{}, utils.BadRequest("Invalid user ID")
	}

	userData, err := s.userRepository.FindByID(parsedUserID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProfileResponse{}, utils.NotFound("User not found")
		}

		return ProfileResponse{}, utils.InternalServerError("Failed to get profile")
	}

	return ProfileResponse{
		ID:        userData.ID,
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      string(userData.Role),
		Avatar:    userData.Avatar,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}, nil
}

func (s *profileService) UpdateMyProfile(
	userID string,
	request UpdateProfileRequest,
) (ProfileResponse, error) {

	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return ProfileResponse{}, utils.BadRequest("Invalid user ID")
	}

	userData, err := s.userRepository.FindByID(parsedUserID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProfileResponse{}, utils.NotFound("User not found")
		}

		return ProfileResponse{}, utils.InternalServerError("Failed to get profile")
	}

	if request.Name != "" {
		userData.Name = request.Name
	}

	if request.Avatar != "" {
		userData.Avatar = request.Avatar
	}

	updatedUser, err := s.userRepository.Update(userData)

	if err != nil {
		return ProfileResponse{}, utils.InternalServerError("Failed to update profile")
	}

	return ProfileResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Role:      string(updatedUser.Role),
		Avatar:    updatedUser.Avatar,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}