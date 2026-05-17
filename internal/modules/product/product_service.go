package product

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProducts() ([]ProductResponse, error)
	GetProductByID(id string) (ProductResponse, error)
	GetProductBySlug(slug string) (ProductResponse, error)
	GetNewArrivals() ([]ProductResponse, error)
	GetFeaturedProducts() ([]ProductResponse, error)
	GetProductsByCategory(category string) ([]ProductResponse, error)
	CreateProduct(request CreateProductRequest) (ProductResponse, error)
	UpdateProduct(id string, request UpdateProductRequest) (ProductResponse, error)
	DeleteProduct(id string) error
}

type productService struct {
	productRepository ProductRepository
}

func NewProductService(productRepository ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) GetAllProducts() ([]ProductResponse, error) {
	products, err := s.productRepository.FindAll()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get products")
	}

	return ToProductResponses(products), nil
}

func (s *productService) GetProductByID(id string) (ProductResponse, error) {
	productID, err := uuid.Parse(id)

	if err != nil {
		return ProductResponse{}, utils.BadRequest("Invalid product ID")
	}

	productData, err := s.productRepository.FindByID(productID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProductResponse{}, utils.NotFound("Product not found")
		}

		return ProductResponse{}, utils.InternalServerError("Failed to get product")
	}

	return ToProductResponse(productData), nil
}

func (s *productService) GetProductBySlug(slug string) (ProductResponse, error) {
	productData, err := s.productRepository.FindBySlug(slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProductResponse{}, utils.NotFound("Product not found")
		}

		return ProductResponse{}, utils.InternalServerError("Failed to get product")
	}

	return ToProductResponse(productData), nil
}

func (s *productService) GetNewArrivals() ([]ProductResponse, error) {
	products, err := s.productRepository.FindNewArrivals()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get new arrivals")
	}

	return ToProductResponses(products), nil
}

func (s *productService) GetFeaturedProducts() ([]ProductResponse, error) {
	products, err := s.productRepository.FindFeatured()

	if err != nil {
		return nil, utils.InternalServerError("Failed to get featured products")
	}

	return ToProductResponses(products), nil
}

func (s *productService) GetProductsByCategory(category string) ([]ProductResponse, error) {
	products, err := s.productRepository.FindByCategory(category)

	if err != nil {
		return nil, utils.InternalServerError("Failed to get products by category")
	}

	return ToProductResponses(products), nil
}

func (s *productService) CreateProduct(request CreateProductRequest) (ProductResponse, error) {
	stockStatus := generateStockStatus(request.Stock)

	newProduct := Product{
		Name:         request.Name,
		Slug:         request.Slug,
		Description:  request.Description,
		Category:     request.Category,
		Collection:   request.Collection,
		Image:        request.Image,
		Price:        request.Price,
		Stock:        request.Stock,
		StockStatus:  stockStatus,
		IsNewArrival: request.IsNewArrival,
		IsFeatured:   request.IsFeatured,
	}

	productData, err := s.productRepository.Create(newProduct)

	if err != nil {
		return ProductResponse{}, utils.InternalServerError("Failed to create product")
	}

	return ToProductResponse(productData), nil
}

func (s *productService) UpdateProduct(id string, request UpdateProductRequest) (ProductResponse, error) {
	productID, err := uuid.Parse(id)

	if err != nil {
		return ProductResponse{}, utils.BadRequest("Invalid product ID")
	}

	productData, err := s.productRepository.FindByID(productID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProductResponse{}, utils.NotFound("Product not found")
		}

		return ProductResponse{}, utils.InternalServerError("Failed to get product")
	}

	if request.Name != "" {
		productData.Name = request.Name
	}

	if request.Slug != "" {
		productData.Slug = request.Slug
	}

	if request.Description != "" {
		productData.Description = request.Description
	}

	if request.Category != "" {
		productData.Category = request.Category
	}

	if request.Collection != "" {
		productData.Collection = request.Collection
	}

	if request.Image != "" {
		productData.Image = request.Image
	}

	if request.Price > 0 {
		productData.Price = request.Price
	}

	if request.Stock >= 0 {
		productData.Stock = request.Stock
		productData.StockStatus = generateStockStatus(request.Stock)
	}

	productData.IsNewArrival = request.IsNewArrival
	productData.IsFeatured = request.IsFeatured

	updatedProduct, err := s.productRepository.Update(productData)

	if err != nil {
		return ProductResponse{}, utils.InternalServerError("Failed to update product")
	}

	return ToProductResponse(updatedProduct), nil
}

func (s *productService) DeleteProduct(id string) error {
	productID, err := uuid.Parse(id)

	if err != nil {
		return utils.BadRequest("Invalid product ID")
	}

	_, err = s.productRepository.FindByID(productID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Product not found")
		}

		return utils.InternalServerError("Failed to get product")
	}

	err = s.productRepository.Delete(productID)

	if err != nil {
		return utils.InternalServerError("Failed to delete product")
	}

	return nil
}

func generateStockStatus(stock int) StockStatus {
	if stock <= 0 {
		return StockOutOfStock
	}

	if stock <= 5 {
		return StockLowStock
	}

	return StockInStock
}