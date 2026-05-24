package category

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryService interface {
	GetAllCategories(filter CategoryFilter) (utils.PaginatedData, error)
	GetActiveCategories() ([]CategoryResponse, error)
	GetCategoryByID(id string) (CategoryResponse, error)
	GetCategoryBySlug(slug string) (CategoryResponse, error)
	CreateCategory(request CreateCategoryRequest) (CategoryResponse, error)
	UpdateCategory(id string, request UpdateCategoryRequest) (CategoryResponse, error)
	DeleteCategory(id string) error
}

type categoryService struct {
	categoryRepository CategoryRepository
}

func NewCategoryService(categoryRepository CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (s *categoryService) GetAllCategories(filter CategoryFilter) (utils.PaginatedData, error) {
	categories, total, err := s.categoryRepository.FindAllPaginated(filter)

	if err != nil {
		return utils.PaginatedData{}, utils.InternalServerError("Failed to get categories")
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}

	limit := filter.Limit
	if limit < 1 {
		limit = 12
	}

	return utils.PaginatedData{
		Items:      ToCategoryResponses(categories),
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *categoryService) GetActiveCategories() ([]CategoryResponse, error) {
	categories, err := s.categoryRepository.FindActive()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get active categories")
	}

	return ToCategoryResponses(categories), nil
}

func (s *categoryService) GetCategoryByID(id string) (CategoryResponse, error) {
	categoryID, err := uuid.Parse(id)

	if err != nil {
		return CategoryResponse{}, utils.BadRequest("Invalid category ID")
	}

	categoryData, err := s.categoryRepository.FindByID(categoryID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CategoryResponse{}, utils.NotFound("Category not found")
		}

		return CategoryResponse{}, utils.InternalServerError("Failed to get category")
	}

	return ToCategoryResponse(categoryData), nil
}

func (s *categoryService) GetCategoryBySlug(slug string) (CategoryResponse, error) {
	categoryData, err := s.categoryRepository.FindBySlug(slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CategoryResponse{}, utils.NotFound("Category not found")
		}

		return CategoryResponse{}, utils.InternalServerError("Failed to get category")
	}

	return ToCategoryResponse(categoryData), nil
}

func (s *categoryService) CreateCategory(request CreateCategoryRequest) (CategoryResponse, error) {
	_, err := s.categoryRepository.FindBySlug(request.Slug)

	if err == nil {
		return CategoryResponse{}, utils.Conflict("Category slug already exists")
	}

	newCategory := Category{
		Name:        request.Name,
		Slug:        request.Slug,
		Description: request.Description,
		Image:       request.Image,
		IsActive:    request.IsActive,
	}

	categoryData, err := s.categoryRepository.Create(newCategory)

	if err != nil {
		return CategoryResponse{}, utils.InternalServerError("Failed to create category")
	}

	return ToCategoryResponse(categoryData), nil
}

func (s *categoryService) UpdateCategory(id string, request UpdateCategoryRequest) (CategoryResponse, error) {
	categoryID, err := uuid.Parse(id)

	if err != nil {
		return CategoryResponse{}, utils.BadRequest("Invalid category ID")
	}

	categoryData, err := s.categoryRepository.FindByID(categoryID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CategoryResponse{}, utils.NotFound("Category not found")
		}

		return CategoryResponse{}, utils.InternalServerError("Failed to get category")
	}

	if request.Name != "" {
		categoryData.Name = request.Name
	}

	if request.Slug != "" {
		categoryData.Slug = request.Slug
	}

	if request.Description != "" {
		categoryData.Description = request.Description
	}

	if request.Image != "" {
		categoryData.Image = request.Image
	}

	if request.IsActive != nil {
		categoryData.IsActive = *request.IsActive
	}

	updatedCategory, err := s.categoryRepository.Update(categoryData)

	if err != nil {
		return CategoryResponse{}, utils.InternalServerError("Failed to update category")
	}

	return ToCategoryResponse(updatedCategory), nil
}

func (s *categoryService) DeleteCategory(id string) error {
	categoryID, err := uuid.Parse(id)

	if err != nil {
		return utils.BadRequest("Invalid category ID")
	}

	_, err = s.categoryRepository.FindByID(categoryID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Category not found")
		}

		return utils.InternalServerError("Failed to get category")
	}

	err = s.categoryRepository.Delete(categoryID)

	if err != nil {
		return utils.InternalServerError("Failed to delete category")
	}

	return nil
}