package auth

import (
	"net/http"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.authService.Register(request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"Register successful",
		response,
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.authService.Login(request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Login successful",
		response,
	)
}