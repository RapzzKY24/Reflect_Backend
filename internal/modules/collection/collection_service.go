package collection

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionService interface {
	GetAllCollections() ([]CollectionResponse, error)
	GetActiveCollections() ([]CollectionResponse, error)
	GetFeaturedCollections() ([]CollectionResponse, error)
	GetCollectionByID(id string) (CollectionResponse, error)
	GetCollectionBySlug(slug string) (CollectionResponse, error)
	CreateCollection(request CreateCollectionRequest) (CollectionResponse, error)
	UpdateCollection(id string, request UpdateCollectionRequest) (CollectionResponse, error)
	DeleteCollection(id string) error
}

type collectionService struct {
	collectionRepository CollectionRepository
}

func NewCollectionService(collectionRepository CollectionRepository) CollectionService {
	return &collectionService{collectionRepository: collectionRepository}
}

func (s *collectionService) GetAllCollections() ([]CollectionResponse, error) {
	collections, err := s.collectionRepository.FindAll()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get collections")
	}

	return ToCollectionResponses(collections), nil
}

func (s *collectionService) GetActiveCollections() ([]CollectionResponse, error) {
	collections, err := s.collectionRepository.FindActive()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get active collections")
	}

	return ToCollectionResponses(collections), nil
}

func (s *collectionService) GetFeaturedCollections() ([]CollectionResponse, error) {
	collections, err := s.collectionRepository.FindFeatured()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get featured collections")
	}

	return ToCollectionResponses(collections), nil
}

func (s *collectionService) GetCollectionByID(id string) (CollectionResponse, error) {
	collectionID, err := uuid.Parse(id)

	if err != nil {
		return CollectionResponse{}, utils.BadRequest("Invalid collection ID")
	}

	collectionData, err := s.collectionRepository.FindByID(collectionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CollectionResponse{}, utils.NotFound("Collection not found")
		}

		return CollectionResponse{}, utils.InternalServerError("Failed to get collection")
	}

	return ToCollectionResponse(collectionData), nil
}

func (s *collectionService) GetCollectionBySlug(slug string) (CollectionResponse, error) {
	collectionData, err := s.collectionRepository.FindBySlug(slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CollectionResponse{}, utils.NotFound("Collection not found")
		}

		return CollectionResponse{}, utils.InternalServerError("Failed to get collection")
	}

	return ToCollectionResponse(collectionData), nil
}

func (s *collectionService) CreateCollection(request CreateCollectionRequest) (CollectionResponse, error) {
	_, err := s.collectionRepository.FindBySlug(request.Slug)

	if err == nil {
		return CollectionResponse{}, utils.Conflict("Collection slug already exists")
	}

	newCollection := Collection{
		Name:        request.Name,
		Slug:        request.Slug,
		Description: request.Description,
		Image:       request.Image,
		IsFeatured:  request.IsFeatured,
		IsActive:    request.IsActive,
	}

	collectionData, err := s.collectionRepository.Create(newCollection)

	if err != nil {
		return CollectionResponse{}, utils.InternalServerError("Failed to create collection")
	}

	return ToCollectionResponse(collectionData), nil
}

func (s *collectionService) UpdateCollection(id string, request UpdateCollectionRequest) (CollectionResponse, error) {
	collectionID, err := uuid.Parse(id)

	if err != nil {
		return CollectionResponse{}, utils.BadRequest("Invalid collection ID")
	}

	collectionData, err := s.collectionRepository.FindByID(collectionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CollectionResponse{}, utils.NotFound("Collection not found")
		}

		return CollectionResponse{}, utils.InternalServerError("Failed to get collection")
	}

	if request.Name != "" {
		collectionData.Name = request.Name
	}

	if request.Slug != "" {
		collectionData.Slug = request.Slug
	}

	if request.Description != "" {
		collectionData.Description = request.Description
	}

	if request.Image != "" {
		collectionData.Image = request.Image
	}

	if request.IsFeatured != nil {
		collectionData.IsFeatured = *request.IsFeatured
	}

	if request.IsActive != nil {
		collectionData.IsActive = *request.IsActive
	}

	updatedCollection, err := s.collectionRepository.Update(collectionData)

	if err != nil {
		return CollectionResponse{}, utils.InternalServerError("Failed to update collection")
	}

	return ToCollectionResponse(updatedCollection), nil
}

func (s *collectionService) DeleteCollection(id string) error {
	collectionID, err := uuid.Parse(id)

	if err != nil {
		return utils.BadRequest("Invalid collection ID")
	}

	_, err = s.collectionRepository.FindByID(collectionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Collection not found")
		}

		return utils.InternalServerError("Failed to get collection")
	}

	err = s.collectionRepository.Delete(collectionID)

	if err != nil {
		return utils.InternalServerError("Failed to delete collection")
	}

	return nil
}