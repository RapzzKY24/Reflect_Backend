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

// @Summary      Register a new user
// @Description  Create a new user account with name, email, and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=auth.AuthResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Router       /auth/register [post]
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

// @Summary      Login user
// @Description  Authenticate user with email and password, returns JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=auth.AuthResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /auth/login [post]
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