package upload

import (
	"context"

	"reflect-backend/internal/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func NewCloudinary(cfg *config.Config) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)

	if err != nil {
		return nil, err
	}

	cld.Config.URL.Secure = true

	return cld, nil
}

func UploadImage(
	cld *cloudinary.Cloudinary,
	file interface{},
	folder string,
) (string, error) {

	ctx := context.Background()

	result, err := cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			Folder: folder,
		},
	)

	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}