package product

import (
	"net/http"
	"strconv"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(productService ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// @Summary      Get all products
// @Description  Retrieve a paginated list of products with search, category filter, and sorting
// @Tags         Products
// @Produce      json
// @Param        page     query int    false "Page number (default: 1)"
// @Param        limit    query int    false "Items per page (default: 12)"
// @Param        search   query string false "Search products by name (partial match)"
// @Param        category query string false "Filter by category"
// @Param        sort     query string false "Sort by (price_asc, price_desc, newest, oldest)"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=utils.PaginatedData{items=[]product.ProductResponse}}
// @Router       /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := c.Query("search")
	category := c.Query("category")
	sortBy := c.DefaultQuery("sort", "newest")

	filter := ProductFilter{
		Search:   search,
		Category: category,
		SortBy:   sortBy,
		Page:     page,
		Limit:    limit,
	}

	result, err := h.productService.GetAllProducts(filter)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessPaginatedResponse(
		c,
		http.StatusOK,
		"Products retrieved successfully",
		result.Items,
		result.TotalItems,
		result.Page,
		result.Limit,
	)
}

// @Summary      Get product by ID
// @Description  Retrieve a single product by its UUID
// @Tags         Products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=product.ProductResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProductByID(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// @Summary      Get product by slug
// @Description  Retrieve a single product by its URL slug
// @Tags         Products
// @Produce      json
// @Param        slug path string true "Product Slug"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=product.ProductResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.productService.GetProductBySlug(slug)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// @Summary      Get related products
// @Description  Retrieve up to 3 related products from the same category, excluding the current product
// @Tags         Products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]product.ProductResponse}
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /products/{id}/related [get]
func (h *ProductHandler) GetRelatedProducts(c *gin.Context) {
	id := c.Param("id")

	products, err := h.productService.GetRelatedProducts(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Related products retrieved successfully",
		products,
	)
}

// @Summary      Get new arrivals
// @Description  Retrieve products marked as new arrivals
// @Tags         Products
// @Produce      json
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]product.ProductResponse}
// @Router       /products/new-arrivals [get]
func (h *ProductHandler) GetNewArrivals(c *gin.Context) {
	products, err := h.productService.GetNewArrivals()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "New arrivals retrieved successfully", products)
}

// @Summary      Get featured products
// @Description  Retrieve products marked as featured
// @Tags         Products
// @Produce      json
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]product.ProductResponse}
// @Router       /products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	products, err := h.productService.GetFeaturedProducts()

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Featured products retrieved successfully", products)
}

// @Summary      Get products by category
// @Description  Retrieve products filtered by category name
// @Tags         Products
// @Produce      json
// @Param        category query string true "Category name (e.g. hoodie, t-shirt)"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]product.ProductResponse}
// @Failure      400 {object} utils.SwaggerErrorResponse
// @Router       /products/category [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Query("category")

	if category == "" {
		c.Error(utils.BadRequest("Category query is required"))
		return
	}

	products, err := h.productService.GetProductsByCategory(category)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products by category retrieved successfully", products)
}

// @Summary      Create product (Admin)
// @Description  Create a new product (admin only)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateProductRequest true "Create Product Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=product.ProductResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Router       /admin/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	product, err := h.productService.CreateProduct(request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

// @Summary      Update product (Admin)
// @Description  Update an existing product by ID (admin only)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Product ID"
// @Param        request body UpdateProductRequest true "Update Product Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=product.ProductResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var request UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(id, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

// @Summary      Delete product (Admin)
// @Description  Delete a product by ID (admin only)
// @Tags         Admin Products
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Product ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      403 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /admin/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.DeleteProduct(id)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}