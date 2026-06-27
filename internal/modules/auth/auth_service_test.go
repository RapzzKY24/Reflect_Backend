package auth

import (
	"testing"

	"reflect-backend/internal/config"
	"reflect-backend/internal/modules/user"
	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) FindAll() ([]user.User, error) {
	args := m.Called()
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *mockUserRepository) FindByID(id uuid.UUID) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepository) FindByEmail(email string) (user.User, error) {
	args := m.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepository) Create(u user.User) (user.User, error) {
	args := m.Called(u)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepository) Update(u user.User) (user.User, error) {
	args := m.Called(u)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(mockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}

	mockRepo.On("FindByEmail", "new@test.com").Return(user.User{}, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.MatchedBy(func(u user.User) bool {
		return u.Email == "new@test.com" && u.Name == "New User" && u.Password != ""
	})).Return(user.User{
		ID:    uuid.New(),
		Name:  "New User",
		Email: "new@test.com",
		Role:  user.RoleUser,
	}, nil)

	service := NewAuthService(mockRepo, cfg)
	resp, err := service.Register(RegisterRequest{
		Name:     "New User",
		Email:    "new@test.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "New User", resp.User.Name)
	assert.Equal(t, "new@test.com", resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}

	existingUser := user.User{
		ID:    uuid.New(),
		Name:  "Existing",
		Email: "existing@test.com",
		Role:  user.RoleUser,
	}

	mockRepo.On("FindByEmail", "existing@test.com").Return(existingUser, nil)

	service := NewAuthService(mockRepo, cfg)
	resp, err := service.Register(RegisterRequest{
		Name:     "New User",
		Email:    "existing@test.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Email already exists", err.Error())
	assert.Empty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}

func TestRegister_CreateFails(t *testing.T) {
	mockRepo := new(mockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}

	mockRepo.On("FindByEmail", "new@test.com").Return(user.User{}, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.Anything).Return(user.User{}, assert.AnError)

	service := NewAuthService(mockRepo, cfg)
	resp, err := service.Register(RegisterRequest{
		Name:     "New User",
		Email:    "new@test.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Empty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(mockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}

	hash, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	userID := uuid.New()
	now := "2024-01-01T00:00:00Z"

	_ = now

	existingUser := user.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@test.com",
		Password: string(hash),
		Role:     user.RoleUser,
	}

	mockRepo.On("FindByEmail", "test@test.com").Return(existingUser, nil)

	service := NewAuthService(mockRepo, cfg)
	resp, err := service.Login(LoginRequest{
		Email:    "test@test.com",
		Password: "Password123!",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "Test User", resp.User.Name)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockRepo := new(mockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}

	mockRepo.On("FindByEmail", "notfound@test.com").Return(user.User{}, gorm.ErrRecordNotFound)

	service := NewAuthService(mockRepo, cfg)
	resp, err := service.Login(LoginRequest{
		Email:    "notfound@test.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Invalid email or password", err.Error())
	assert.Empty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(mockUserRepository)
	cfg := &config.Config{JWTSecret: "test-secret"}

	hash, _ := bcrypt.GenerateFromPassword([]byte("CorrectPass1!"), bcrypt.DefaultCost)

	existingUser := user.User{
		ID:       uuid.New(),
		Name:     "Test User",
		Email:    "test@test.com",
		Password: string(hash),
		Role:     user.RoleUser,
	}

	mockRepo.On("FindByEmail", "test@test.com").Return(existingUser, nil)

	service := NewAuthService(mockRepo, cfg)
	resp, err := service.Login(LoginRequest{
		Email:    "test@test.com",
		Password: "WrongPass1!",
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Invalid email or password", err.Error())
	assert.Empty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}
