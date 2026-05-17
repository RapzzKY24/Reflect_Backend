package upload

import (
	"net/http"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService UploadService
}

func NewUploadHandler(
	uploadService UploadService,
) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")

	if err != nil {
		utils.ValidationErrorResponse(c, "Image is required")
		return
	}

	imageURL, err := h.uploadService.UploadImage(file)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Image uploaded successfully",
		gin.H{
			"url": imageURL,
		},
	)
}