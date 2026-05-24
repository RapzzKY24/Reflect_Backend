package category

import (
	"net/http"
	"strconv"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// @Summary      Get all categories
// @Description  Retrieve a paginated list of categories with search and sorting
// @Tags         Categories
// @Produce      json
// @Param        page   query int    false "Page number (default: 1)"
// @Param        limit  query int    false "Items per page (default: 12)"
// @Param        search query string false "Search categories by name (partial match)"
// @Param        sort   query string false "Sort by (name_asc, name_desc)"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=utils.PaginatedData{items=[]category.CategoryResponse}}
// @Router       /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort", "name_asc")

	filter := CategoryFilter{
		Search: search,
		SortBy: sortBy,
		Page:   page,
		Limit:  limit,
	}

	result, err := h.categoryService.GetAllCategories(filter)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessPaginatedResponse(
		c,
		http.StatusOK,
		"Categories retrieved successfully",
		result.Items,
		result.TotalItems,
		result.Page,
		result.Limit,
	)
}

// @Summary      Get active categories
// @Description  Retrieve only active categories
// @Tags         Categories
// @Produce      json
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]category.CategoryResponse}
// @Router       /categories/active [get]
func (h *CategoryHandler) GetActiveCategories(c *gin.Context) {
	categories, err := h.categoryService.GetActiveCategories()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Active categories retrieved successfully", categories)
}

// @Summary      Get category by ID
// @Description  Retrieve a single category by its UUID
// @Tags         Categories
// @Produce      json
// @Param        id path string true "Category ID"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=category.CategoryResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetCategoryByID(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// @Summary      Get category by slug
// @Description  Retrieve a single category by its URL slug
// @Tags         Categories
// @Produce      json
// @Param        slug path string true "Category Slug"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=category.CategoryResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryService.GetCategoryBySlug(slug)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// @Summary      Create category (Admin)
// @Description  Create a new category (admin only)
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateCategoryRequest true "Create Category Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=category.CategoryResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Router       /admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var request CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

// @Summary      Update category (Admin)
// @Description  Update an existing category by ID (admin only)
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Category ID"
// @Param        request body UpdateCategoryRequest true "Update Category Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=category.CategoryResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var request UpdateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(id, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

// @Summary      Delete category (Admin)
// @Description  Delete a category by ID (admin only)
// @Tags         Admin Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Category ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryService.DeleteCategory(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}