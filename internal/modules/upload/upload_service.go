package upload

import (
	"mime/multipart"

	"reflect-backend/internal/utils"

	"github.com/cloudinary/cloudinary-go/v2"
)

type UploadService interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

type uploadService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewUploadService(
	cloudinary *cloudinary.Cloudinary,
) UploadService {
	return &uploadService{
		cloudinary: cloudinary,
	}
}

func (s *uploadService) UploadImage(
	file *multipart.FileHeader,
) (string, error) {

	openedFile, err := file.Open()

	if err != nil {
		return "", utils.InternalServerError("Failed to open image")
	}

	defer openedFile.Close()

	imageURL, err := UploadImage(
		s.cloudinary,
		openedFile,
		"reflect",
	)

	if err != nil {
		return "", utils.InternalServerError("Failed to upload image")
	}

	return imageURL, nil
}