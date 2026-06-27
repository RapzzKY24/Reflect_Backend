package product

import (
	"errors"
	"testing"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) FindAll() ([]Product, error) {
	args := m.Called()
	return args.Get(0).([]Product), args.Error(1)
}

func (m *mockProductRepository) FindAllPaginated(filter ProductFilter) ([]Product, int64, error) {
	args := m.Called(filter)
	return args.Get(0).([]Product), args.Get(1).(int64), args.Error(2)
}

func (m *mockProductRepository) FindByID(id uuid.UUID) (Product, error) {
	args := m.Called(id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepository) FindBySlug(slug string) (Product, error) {
	args := m.Called(slug)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepository) FindNewArrivals() ([]Product, error) {
	args := m.Called()
	return args.Get(0).([]Product), args.Error(1)
}

func (m *mockProductRepository) FindFeatured() ([]Product, error) {
	args := m.Called()
	return args.Get(0).([]Product), args.Error(1)
}

func (m *mockProductRepository) FindRelatedProducts(category string, excludeID uuid.UUID, limit int) ([]Product, error) {
	args := m.Called(category, excludeID, limit)
	return args.Get(0).([]Product), args.Error(1)
}

func (m *mockProductRepository) FindByCategory(category string) ([]Product, error) {
	args := m.Called(category)
	return args.Get(0).([]Product), args.Error(1)
}

func (m *mockProductRepository) Create(product Product) (Product, error) {
	args := m.Called(product)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepository) Update(product Product) (Product, error) {
	args := m.Called(product)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetProductByID_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	expectedProduct := Product{
		ID:    productID,
		Name:  "Test Product",
		Slug:  "test-product",
		Price: 100000,
		Stock: 10,
	}

	mockRepo.On("FindByID", productID).Return(expectedProduct, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetProductByID(productID.String())

	assert.NoError(t, err)
	assert.Equal(t, "Test Product", resp.Name)
	assert.Equal(t, 100000, resp.Price)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	mockRepo.On("FindByID", productID).Return(Product{}, gorm.ErrRecordNotFound)

	service := NewProductService(mockRepo)
	resp, err := service.GetProductByID(productID.String())

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Product not found", err.Error())
	assert.Equal(t, "", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_InvalidID(t *testing.T) {
	mockRepo := new(mockProductRepository)

	service := NewProductService(mockRepo)
	resp, err := service.GetProductByID("invalid-uuid")

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Invalid product ID", err.Error())
	assert.Equal(t, "", resp.Name)
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestGetProductBySlug_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)

	expectedProduct := Product{
		ID:    uuid.New(),
		Name:  "Slug Product",
		Slug:  "slug-product",
		Price: 50000,
	}

	mockRepo.On("FindBySlug", "slug-product").Return(expectedProduct, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetProductBySlug("slug-product")

	assert.NoError(t, err)
	assert.Equal(t, "Slug Product", resp.Name)
	assert.Equal(t, "slug-product", resp.Slug)
	mockRepo.AssertExpectations(t)
}

func TestGetProductBySlug_NotFound(t *testing.T) {
	mockRepo := new(mockProductRepository)

	mockRepo.On("FindBySlug", "nonexistent").Return(Product{}, gorm.ErrRecordNotFound)

	service := NewProductService(mockRepo)
	resp, err := service.GetProductBySlug("nonexistent")

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Product not found", err.Error())
	assert.Equal(t, "", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	products := []Product{
		{ID: uuid.New(), Name: "Product 1", Price: 10000, Stock: 5},
		{ID: uuid.New(), Name: "Product 2", Price: 20000, Stock: 10},
	}

	filter := ProductFilter{Page: 1, Limit: 10}
	mockRepo.On("FindAllPaginated", filter).Return(products, int64(2), nil)

	service := NewProductService(mockRepo)
	result, err := service.GetAllProducts(filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.TotalItems)
	assert.Equal(t, 2, len(result.Items.([]ProductResponse)))
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	req := CreateProductRequest{
		Name:     "New Product",
		Slug:     "new-product",
		Category: CategoryHoodie,
		Price:    150000,
		Stock:    10,
	}

	mockRepo.On("Create", mock.MatchedBy(func(p Product) bool {
		return p.Name == "New Product" && p.Price == 150000
	})).Return(Product{
		ID:          productID,
		Name:        "New Product",
		Slug:        "new-product",
		Category:    CategoryHoodie,
		Price:       150000,
		Stock:       10,
		StockStatus: StockInStock,
	}, nil)

	service := NewProductService(mockRepo)
	resp, err := service.CreateProduct(req)

	assert.NoError(t, err)
	assert.Equal(t, "New Product", resp.Name)
	assert.Equal(t, 150000, resp.Price)
	assert.Equal(t, StockInStock, resp.StockStatus)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_OutOfStock(t *testing.T) {
	mockRepo := new(mockProductRepository)

	req := CreateProductRequest{
		Name:     "Out of Stock Product",
		Slug:     "out-of-stock",
		Category: CategoryTShirt,
		Price:    50000,
		Stock:    0,
	}

	mockRepo.On("Create", mock.Anything).Return(Product{
		ID:          uuid.New(),
		Name:        "Out of Stock Product",
		Slug:        "out-of-stock",
		Category:    CategoryTShirt,
		Price:       50000,
		Stock:       0,
		StockStatus: StockOutOfStock,
	}, nil)

	service := NewProductService(mockRepo)
	resp, err := service.CreateProduct(req)

	assert.NoError(t, err)
	assert.Equal(t, StockOutOfStock, resp.StockStatus)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	existingProduct := Product{
		ID:    productID,
		Name:  "Old Name",
		Slug:  "old-name",
		Price: 50000,
		Stock: 5,
	}

	req := UpdateProductRequest{
		Name:  "Updated Name",
		Price: 75000,
		Stock: 10,
	}

	mockRepo.On("FindByID", productID).Return(existingProduct, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p Product) bool {
		return p.Name == "Updated Name" && p.Price == 75000 && p.Stock == 10
	})).Return(Product{
		ID:          productID,
		Name:        "Updated Name",
		Slug:        "old-name",
		Price:       75000,
		Stock:       10,
		StockStatus: StockInStock,
	}, nil)

	service := NewProductService(mockRepo)
	resp, err := service.UpdateProduct(productID.String(), req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", resp.Name)
	assert.Equal(t, 75000, resp.Price)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	mockRepo.On("FindByID", productID).Return(Product{}, gorm.ErrRecordNotFound)

	service := NewProductService(mockRepo)
	resp, err := service.UpdateProduct(productID.String(), UpdateProductRequest{
		Name: "Updated Name",
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Product not found", err.Error())
	assert.Equal(t, "", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	mockRepo.On("FindByID", productID).Return(Product{ID: productID}, nil)
	mockRepo.On("Delete", productID).Return(nil)

	service := NewProductService(mockRepo)
	err := service.DeleteProduct(productID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	mockRepo.On("FindByID", productID).Return(Product{}, gorm.ErrRecordNotFound)

	service := NewProductService(mockRepo)
	err := service.DeleteProduct(productID.String())

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Product not found", err.Error())
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestGetNewArrivals_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	products := []Product{
		{ID: uuid.New(), Name: "New Arrival 1", Price: 50000, IsNewArrival: true},
		{ID: uuid.New(), Name: "New Arrival 2", Price: 75000, IsNewArrival: true},
	}

	mockRepo.On("FindNewArrivals").Return(products, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetNewArrivals()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))
	mockRepo.AssertExpectations(t)
}

func TestGetNewArrivals_Empty(t *testing.T) {
	mockRepo := new(mockProductRepository)

	mockRepo.On("FindNewArrivals").Return([]Product{}, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetNewArrivals()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp))
	mockRepo.AssertExpectations(t)
}

func TestGetFeaturedProducts_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	products := []Product{
		{ID: uuid.New(), Name: "Featured 1", Price: 100000, IsFeatured: true},
	}

	mockRepo.On("FindFeatured").Return(products, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetFeaturedProducts()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp))
	mockRepo.AssertExpectations(t)
}

func TestGetRelatedProducts_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()
	category := CategoryHoodie

	mockRepo.On("FindByID", productID).Return(Product{
		ID:       productID,
		Name:     "Original",
		Category: category,
	}, nil)

	related := []Product{
		{ID: uuid.New(), Name: "Related 1", Category: category, Price: 50000},
		{ID: uuid.New(), Name: "Related 2", Category: category, Price: 75000},
	}

	mockRepo.On("FindRelatedProducts", "hoodie", productID, 3).Return(related, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetRelatedProducts(productID.String())

	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))
	mockRepo.AssertExpectations(t)
}

func TestGetRelatedProducts_OriginalNotFound(t *testing.T) {
	mockRepo := new(mockProductRepository)
	productID := uuid.New()

	mockRepo.On("FindByID", productID).Return(Product{}, gorm.ErrRecordNotFound)

	service := NewProductService(mockRepo)
	resp, err := service.GetRelatedProducts(productID.String())

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByCategory_Success(t *testing.T) {
	mockRepo := new(mockProductRepository)
	products := []Product{
		{ID: uuid.New(), Name: "Hoodie 1", Category: CategoryHoodie, Price: 100000},
	}

	mockRepo.On("FindByCategory", "hoodie").Return(products, nil)

	service := NewProductService(mockRepo)
	resp, err := service.GetProductsByCategory("hoodie")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp))
	mockRepo.AssertExpectations(t)
}

func TestServiceError_ReturnsAppError(t *testing.T) {
	mockRepo := new(mockProductRepository)

	mockRepo.On("FindNewArrivals").Return([]Product{}, errors.New("db error"))

	service := NewProductService(mockRepo)
	_, err := service.GetNewArrivals()

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_DefaultPagination(t *testing.T) {
	mockRepo := new(mockProductRepository)

	filter := ProductFilter{}
	mockRepo.On("FindAllPaginated", filter).Return([]Product{}, int64(0), nil)

	service := NewProductService(mockRepo)
	result, err := service.GetAllProducts(filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 12, result.Limit)
	mockRepo.AssertExpectations(t)
}
