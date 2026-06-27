package order

import (
	"errors"
	"testing"

	"reflect-backend/internal/modules/cart"
	"reflect-backend/internal/modules/product"
	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockOrderRepository struct {
	mock.Mock
}

func (m *mockOrderRepository) FindByUserID(userID uuid.UUID) ([]Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]Order), args.Error(1)
}

func (m *mockOrderRepository) FindByID(id uuid.UUID) (Order, error) {
	args := m.Called(id)
	return args.Get(0).(Order), args.Error(1)
}

func (m *mockOrderRepository) FindByOrderNumber(orderNumber string) (Order, error) {
	args := m.Called(orderNumber)
	return args.Get(0).(Order), args.Error(1)
}

func (m *mockOrderRepository) FindCartItems(userID uuid.UUID) ([]cart.CartItem, error) {
	args := m.Called(userID)
	return args.Get(0).([]cart.CartItem), args.Error(1)
}

func (m *mockOrderRepository) FindProductForUpdate(tx *gorm.DB, productID uuid.UUID) (product.Product, error) {
	args := m.Called(productID)
	return args.Get(0).(product.Product), args.Error(1)
}

func (m *mockOrderRepository) CreateOrderWithTx(tx *gorm.DB, order Order) (Order, error) {
	args := m.Called(order)
	return args.Get(0).(Order), args.Error(1)
}

func (m *mockOrderRepository) CreateOrderItemWithTx(tx *gorm.DB, item OrderItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *mockOrderRepository) CreateTrackingWithTx(tx *gorm.DB, tracking OrderTracking) error {
	args := m.Called(tracking)
	return args.Error(0)
}

func (m *mockOrderRepository) UpdateProductStockWithTx(tx *gorm.DB, p product.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockOrderRepository) ClearCartWithTx(tx *gorm.DB, userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *mockOrderRepository) UpdateOrder(order Order) (Order, error) {
	args := m.Called(order)
	return args.Get(0).(Order), args.Error(1)
}

func (m *mockOrderRepository) Transaction(fn func(tx *gorm.DB) error) error {
	args := m.Called(fn)

	if args.Error(0) == nil {
		return fn(nil)
	}

	return args.Error(0)
}

func TestCheckout_Success(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()

	cartItems := []cart.CartItem{
		{
			ID:        uuid.New(),
			UserID:    userID,
			ProductID: productID,
			Quantity:  2,
			Product: product.Product{
				ID:    productID,
				Name:  "Test Product",
				Slug:  "test-product",
				Price: 50000,
				Stock: 10,
			},
		},
	}

	expectedOrder := Order{
		ID:          uuid.New(),
		OrderNumber: "BLCK1234567890",
		UserID:      userID,
		AddressID:   addressID,
		Subtotal:    100000,
		Shipping:    15000,
		Total:       115000,
		OrderStatus: OrderPending,
	}

	expectedOrderFull := expectedOrder
	expectedOrderFull.OrderItems = []OrderItem{
		{
			ID:        uuid.New(),
			OrderID:   expectedOrder.ID,
			ProductID: productID,
			Price:     50000,
			Quantity:  2,
			Subtotal:  100000,
		},
	}

	mockRepo.On("FindCartItems", userID).Return(cartItems, nil)
	mockRepo.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
	mockRepo.On("FindProductForUpdate", productID).Return(cartItems[0].Product, nil).Twice()
	mockRepo.On("CreateOrderWithTx", mock.MatchedBy(func(o Order) bool {
		return o.Subtotal == 100000 && o.Shipping == 15000 && o.Total == 115000
	})).Return(expectedOrder, nil)
	mockRepo.On("CreateOrderItemWithTx", mock.Anything).Return(nil)
	mockRepo.On("CreateTrackingWithTx", mock.Anything).Return(nil)
	mockRepo.On("UpdateProductStockWithTx", mock.Anything).Return(nil)
	mockRepo.On("ClearCartWithTx", userID).Return(nil)
	mockRepo.On("FindByID", expectedOrder.ID).Return(expectedOrderFull, nil)

	service := NewOrderService(mockRepo)
	resp, err := service.Checkout(userID.String(), CheckoutRequest{
		AddressID: addressID.String(),
	})

	assert.NoError(t, err)
	assert.Equal(t, 100000, resp.Subtotal)
	assert.Equal(t, 15000, resp.Shipping)
	assert.Equal(t, 115000, resp.Total)
	mockRepo.AssertExpectations(t)
}

func TestCheckout_EmptyCart(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()

	mockRepo.On("FindCartItems", userID).Return([]cart.CartItem{}, nil)

	service := NewOrderService(mockRepo)
	resp, err := service.Checkout(userID.String(), CheckoutRequest{
		AddressID: uuid.New().String(),
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Cart is empty", err.Error())
	assert.Equal(t, "", resp.OrderNumber)
	mockRepo.AssertExpectations(t)
}

func TestCheckout_InvalidUserID(t *testing.T) {
	mockRepo := new(mockOrderRepository)

	service := NewOrderService(mockRepo)
	resp, err := service.Checkout("invalid-uuid", CheckoutRequest{
		AddressID: uuid.New().String(),
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Invalid user ID", err.Error())
	assert.Equal(t, "", resp.OrderNumber)
}

func TestCheckout_InvalidAddressID(t *testing.T) {
	mockRepo := new(mockOrderRepository)

	service := NewOrderService(mockRepo)
	resp, err := service.Checkout(uuid.New().String(), CheckoutRequest{
		AddressID: "invalid-uuid",
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Invalid address ID", err.Error())
	assert.Equal(t, "", resp.OrderNumber)
}

func TestCheckout_InsufficientStock(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()
	addressID := uuid.New()
	productID := uuid.New()

	cartItems := []cart.CartItem{
		{
			ID:        uuid.New(),
			UserID:    userID,
			ProductID: productID,
			Quantity:  20,
			Product: product.Product{
				ID:    productID,
				Name:  "Low Stock Product",
				Price: 50000,
				Stock: 5,
			},
		},
	}

	mockRepo.On("FindCartItems", userID).Return(cartItems, nil)
	mockRepo.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
	mockRepo.On("FindProductForUpdate", productID).Return(cartItems[0].Product, nil)

	service := NewOrderService(mockRepo)
	resp, err := service.Checkout(userID.String(), CheckoutRequest{
		AddressID: addressID.String(),
	})

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Product stock is not enough", err.Error())
	assert.Equal(t, "", resp.OrderNumber)
	mockRepo.AssertExpectations(t)
}

func TestGetMyOrders_Success(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()

	orders := []Order{
		{
			ID:          uuid.New(),
			OrderNumber: "BLCK111",
			UserID:      userID,
			Subtotal:    50000,
			Total:       65000,
			OrderStatus: OrderPending,
		},
		{
			ID:          uuid.New(),
			OrderNumber: "BLCK222",
			UserID:      userID,
			Subtotal:    100000,
			Total:       115000,
			OrderStatus: OrderDelivered,
		},
	}

	mockRepo.On("FindByUserID", userID).Return(orders, nil)

	service := NewOrderService(mockRepo)
	resp, err := service.GetMyOrders(userID.String())

	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp))
	assert.Equal(t, "BLCK111", resp[0].OrderNumber)
	assert.Equal(t, "BLCK222", resp[1].OrderNumber)
	mockRepo.AssertExpectations(t)
}

func TestGetMyOrders_InvalidUserID(t *testing.T) {
	mockRepo := new(mockOrderRepository)

	service := NewOrderService(mockRepo)
	resp, err := service.GetMyOrders("invalid-uuid")

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Nil(t, resp)
}

func TestGetMyOrderByID_Success(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()
	orderID := uuid.New()

	orderData := Order{
		ID:          orderID,
		OrderNumber: "BLCK123",
		UserID:      userID,
		Subtotal:    75000,
		Total:       90000,
		OrderStatus: OrderPaid,
	}

	mockRepo.On("FindByID", orderID).Return(orderData, nil)

	service := NewOrderService(mockRepo)
	resp, err := service.GetMyOrderByID(userID.String(), orderID.String())

	assert.NoError(t, err)
	assert.Equal(t, "BLCK123", resp.OrderNumber)
	mockRepo.AssertExpectations(t)
}

func TestGetMyOrderByID_Forbidden(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()
	otherUserID := uuid.New()
	orderID := uuid.New()

	orderData := Order{
		ID:          orderID,
		OrderNumber: "BLCK123",
		UserID:      otherUserID,
	}

	mockRepo.On("FindByID", orderID).Return(orderData, nil)

	service := NewOrderService(mockRepo)
	resp, err := service.GetMyOrderByID(userID.String(), orderID.String())

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "You are not allowed to access this order", err.Error())
	assert.Equal(t, "", resp.OrderNumber)
	mockRepo.AssertExpectations(t)
}

func TestGetMyOrderByID_NotFound(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	orderID := uuid.New()

	mockRepo.On("FindByID", orderID).Return(Order{}, gorm.ErrRecordNotFound)

	service := NewOrderService(mockRepo)
	resp, err := service.GetMyOrderByID(uuid.New().String(), orderID.String())

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	assert.Equal(t, "Order not found", err.Error())
	assert.Equal(t, "", resp.OrderNumber)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	orderID := uuid.New()

	existingOrder := Order{
		ID:            orderID,
		OrderNumber:   "BLCK123",
		OrderStatus:   OrderPending,
		PaymentStatus: PaymentPending,
	}

	updatedOrder := existingOrder
	updatedOrder.OrderStatus = OrderShipped
	updatedOrder.PaymentStatus = PaymentPaid
	updatedOrder.ShippingStatus = ShippingInTransit

	mockRepo.On("FindByID", orderID).Return(existingOrder, nil).Once()
	mockRepo.On("UpdateOrder", mock.MatchedBy(func(o Order) bool {
		return o.OrderStatus == OrderShipped && o.PaymentStatus == PaymentPaid
	})).Return(updatedOrder, nil)
	mockRepo.On("FindByID", orderID).Return(updatedOrder, nil).Once()

	service := NewOrderService(mockRepo)
	resp, err := service.UpdateOrderStatus(orderID.String(), UpdateOrderStatusRequest{
		OrderStatus:    OrderShipped,
		PaymentStatus:  PaymentPaid,
		ShippingStatus: ShippingInTransit,
	})

	assert.NoError(t, err)
	assert.Equal(t, OrderShipped, resp.OrderStatus)
	assert.Equal(t, PaymentPaid, resp.PaymentStatus)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_WithTracking(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	orderID := uuid.New()

	existingOrder := Order{
		ID:          orderID,
		OrderNumber: "BLCK123",
		OrderStatus: OrderProcessing,
	}

	updatedOrder := existingOrder
	updatedOrder.OrderStatus = OrderShipped

	mockRepo.On("FindByID", orderID).Return(existingOrder, nil).Once()
	mockRepo.On("UpdateOrder", mock.Anything).Return(updatedOrder, nil)
	mockRepo.On("Transaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
	mockRepo.On("CreateTrackingWithTx", mock.MatchedBy(func(t OrderTracking) bool {
		return t.Status == "Shipped" && t.Location == "Jakarta"
	})).Return(nil)
	mockRepo.On("FindByID", orderID).Return(updatedOrder, nil).Once()

	service := NewOrderService(mockRepo)
	resp, err := service.UpdateOrderStatus(orderID.String(), UpdateOrderStatusRequest{
		OrderStatus:    OrderShipped,
		TrackingStatus: "Shipped",
		Description:    "Package shipped via JNE",
		Location:       "Jakarta",
	})

	assert.NoError(t, err)
	assert.Equal(t, OrderShipped, resp.OrderStatus)
	mockRepo.AssertExpectations(t)
}

func TestServiceError_PropagatesAppError(t *testing.T) {
	mockRepo := new(mockOrderRepository)
	userID := uuid.New()

	mockRepo.On("FindByUserID", userID).Return([]Order{}, errors.New("db error"))

	service := NewOrderService(mockRepo)
	_, err := service.GetMyOrders(userID.String())

	assert.Error(t, err)
	assert.IsType(t, &utils.AppError{}, err)
	mockRepo.AssertExpectations(t)
}
