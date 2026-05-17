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