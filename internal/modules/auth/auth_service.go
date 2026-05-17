package auth

import (
	"errors"

	"reflect-backend/internal/config"
	"reflect-backend/internal/modules/user"
	"reflect-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(request RegisterRequest) (AuthResponse, error)
	Login(request LoginRequest) (AuthResponse, error)
}

type authService struct {
	userRepository user.UserRepository
	config         *config.Config
}

func NewAuthService(
	userRepository user.UserRepository,
	config *config.Config,
) AuthService {
	return &authService{
		userRepository: userRepository,
		config:         config,
	}
}

func (s *authService) Register(request RegisterRequest) (AuthResponse, error) {
	existingUser, err := s.userRepository.FindByEmail(request.Email)

	if err == nil && existingUser.ID.String() != "" {
		return AuthResponse{}, utils.Conflict("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return AuthResponse{}, utils.InternalServerError("Failed to hash password")
	}

	newUser := user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     user.RoleUser,
	}

	createdUser, err := s.userRepository.Create(newUser)

	if err != nil {
		return AuthResponse{}, utils.InternalServerError("Failed to create user")
	}

	token, err := GenerateJWT(
		createdUser.ID,
		createdUser.Email,
		string(createdUser.Role),
		s.config,
	)

	if err != nil {
		return AuthResponse{}, utils.InternalServerError("Failed to generate token")
	}

	return AuthResponse{
		Token: token,
		User:  user.ToUserResponse(createdUser),
	}, nil
}

func (s *authService) Login(request LoginRequest) (AuthResponse, error) {
	existingUser, err := s.userRepository.FindByEmail(request.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthResponse{}, utils.Unauthorized("Invalid email or password")
		}

		return AuthResponse{}, utils.InternalServerError("Failed to login")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(existingUser.Password),
		[]byte(request.Password),
	)

	if err != nil {
		return AuthResponse{}, utils.Unauthorized("Invalid email or password")
	}

	token, err := GenerateJWT(
		existingUser.ID,
		existingUser.Email,
		string(existingUser.Role),
		s.config,
	)

	if err != nil {
		return AuthResponse{}, utils.InternalServerError("Failed to generate token")
	}

	return AuthResponse{
		Token: token,
		User:  user.ToUserResponse(existingUser),
	}, nil
}