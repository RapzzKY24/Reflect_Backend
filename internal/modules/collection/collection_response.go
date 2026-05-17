package collection

import (
	"time"

	"github.com/google/uuid"
)

type CollectionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	IsFeatured  bool      `json:"is_featured"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToCollectionResponse(collection Collection) CollectionResponse {
	return CollectionResponse{
		ID:          collection.ID,
		Name:        collection.Name,
		Slug:        collection.Slug,
		Description: collection.Description,
		Image:       collection.Image,
		IsFeatured:  collection.IsFeatured,
		IsActive:    collection.IsActive,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
	}
}

func ToCollectionResponses(collections []Collection) []CollectionResponse {
	responses := make([]CollectionResponse, 0)

	for _, collection := range collections {
		responses = append(responses, ToCollectionResponse(collection))
	}

	return responses
}