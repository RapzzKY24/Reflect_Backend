package user

import (
	"net/http"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary      Get all users (Admin)
// @Description  Retrieve a list of all users (admin only)
// @Tags         Admin Users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]user.UserResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}

// @Summary      Get user by ID
// @Description  Retrieve a single user by their UUID (authenticated users can get their own)
// @Tags         Admin Users
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=user.UserResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// @Summary      Delete user (Admin)
// @Description  Delete a user by ID (admin only)
// @Tags         Admin Users
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}