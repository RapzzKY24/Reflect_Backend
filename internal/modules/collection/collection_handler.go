package collection

import (
	"net/http"
	"strconv"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	collectionService CollectionService
}

func NewCollectionHandler(collectionService CollectionService) *CollectionHandler {
	return &CollectionHandler{collectionService: collectionService}
}

// @Summary      Get all collections
// @Description  Retrieve a paginated list of collections with search and sorting
// @Tags         Collections
// @Produce      json
// @Param        page   query int    false "Page number (default: 1)"
// @Param        limit  query int    false "Items per page (default: 12)"
// @Param        search query string false "Search collections by name (partial match)"
// @Param        sort   query string false "Sort by (name_asc, name_desc)"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=utils.PaginatedData{items=[]collection.CollectionResponse}}
// @Router       /collections [get]
func (h *CollectionHandler) GetAllCollections(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort", "name_asc")

	filter := CollectionFilter{
		Search: search,
		SortBy: sortBy,
		Page:   page,
		Limit:  limit,
	}

	result, err := h.collectionService.GetAllCollections(filter)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessPaginatedResponse(
		c,
		http.StatusOK,
		"Collections retrieved successfully",
		result.Items,
		result.TotalItems,
		result.Page,
		result.Limit,
	)
}

// @Summary      Get active collections
// @Description  Retrieve only active collections
// @Tags         Collections
// @Produce      json
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]collection.CollectionResponse}
// @Router       /collections/active [get]
func (h *CollectionHandler) GetActiveCollections(c *gin.Context) {
	collections, err := h.collectionService.GetActiveCollections()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Active collections retrieved successfully", collections)
}

// @Summary      Get featured collections
// @Description  Retrieve only featured collections
// @Tags         Collections
// @Produce      json
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]collection.CollectionResponse}
// @Router       /collections/featured [get]
func (h *CollectionHandler) GetFeaturedCollections(c *gin.Context) {
	collections, err := h.collectionService.GetFeaturedCollections()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Featured collections retrieved successfully", collections)
}

// @Summary      Get collection by ID
// @Description  Retrieve a single collection by its UUID
// @Tags         Collections
// @Produce      json
// @Param        id path string true "Collection ID"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=collection.CollectionResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /collections/{id} [get]
func (h *CollectionHandler) GetCollectionByID(c *gin.Context) {
	id := c.Param("id")

	collection, err := h.collectionService.GetCollectionByID(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection retrieved successfully", collection)
}

// @Summary      Get collection by slug
// @Description  Retrieve a single collection by its URL slug
// @Tags         Collections
// @Produce      json
// @Param        slug path string true "Collection Slug"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=collection.CollectionResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /collections/slug/{slug} [get]
func (h *CollectionHandler) GetCollectionBySlug(c *gin.Context) {
	slug := c.Param("slug")

	collection, err := h.collectionService.GetCollectionBySlug(slug)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection retrieved successfully", collection)
}

// @Summary      Create collection (Admin)
// @Description  Create a new collection (admin only)
// @Tags         Admin Collections
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateCollectionRequest true "Create Collection Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=collection.CollectionResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Router       /admin/collections [post]
func (h *CollectionHandler) CreateCollection(c *gin.Context) {
	var request CreateCollectionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	collection, err := h.collectionService.CreateCollection(request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Collection created successfully", collection)
}

// @Summary      Update collection (Admin)
// @Description  Update an existing collection by ID (admin only)
// @Tags         Admin Collections
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Collection ID"
// @Param        request body UpdateCollectionRequest true "Update Collection Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=collection.CollectionResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/collections/{id} [put]
func (h *CollectionHandler) UpdateCollection(c *gin.Context) {
	id := c.Param("id")

	var request UpdateCollectionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	collection, err := h.collectionService.UpdateCollection(id, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection updated successfully", collection)
}

// @Summary      Delete collection (Admin)
// @Description  Delete a collection by ID (admin only)
// @Tags         Admin Collections
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Collection ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/collections/{id} [delete]
func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	id := c.Param("id")

	err := h.collectionService.DeleteCollection(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection deleted successfully", nil)
}