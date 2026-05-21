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

// @Summary      Upload image
// @Description  Upload an image file to Cloudinary and return the URL
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image formData file true "Image file to upload"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=map[string]string}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /upload [post]
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