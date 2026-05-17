package auth

import (
	"reflect-backend/internal/modules/user"
)

type AuthResponse struct {
	Token string                `json:"token"`
	User  user.UserResponse     `json:"user"`
}