package profile

import (
	"net/http"

	"reflect-backend/internal/middleware"
	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService ProfileService
}

func NewProfileHandler(
	profileService ProfileService,
) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// @Summary      Get my profile
// @Description  Retrieve the authenticated user's profile
// @Tags         Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /me [get]
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	profile, err := h.profileService.GetMyProfile(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Profile retrieved successfully",
		profile,
	)
}

// @Summary      Update my profile
// @Description  Update the authenticated user's profile information
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateProfileRequest true "Update Profile Request"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /me [put]
func (h *ProfileHandler) UpdateMyProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var request UpdateProfileRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	profile, err := h.profileService.UpdateMyProfile(userID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Profile updated successfully",
		profile,
	)
}