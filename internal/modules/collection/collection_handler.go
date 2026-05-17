package collection

import (
	"net/http"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	collectionService CollectionService
}

func NewCollectionHandler(collectionService CollectionService) *CollectionHandler {
	return &CollectionHandler{collectionService: collectionService}
}

func (h *CollectionHandler) GetAllCollections(c *gin.Context) {
	collections, err := h.collectionService.GetAllCollections()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collections retrieved successfully", collections)
}

func (h *CollectionHandler) GetActiveCollections(c *gin.Context) {
	collections, err := h.collectionService.GetActiveCollections()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Active collections retrieved successfully", collections)
}

func (h *CollectionHandler) GetFeaturedCollections(c *gin.Context) {
	collections, err := h.collectionService.GetFeaturedCollections()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Featured collections retrieved successfully", collections)
}

func (h *CollectionHandler) GetCollectionByID(c *gin.Context) {
	id := c.Param("id")

	collection, err := h.collectionService.GetCollectionByID(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection retrieved successfully", collection)
}

func (h *CollectionHandler) GetCollectionBySlug(c *gin.Context) {
	slug := c.Param("slug")

	collection, err := h.collectionService.GetCollectionBySlug(slug)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection retrieved successfully", collection)
}

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

func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	id := c.Param("id")

	err := h.collectionService.DeleteCollection(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Collection deleted successfully", nil)
}