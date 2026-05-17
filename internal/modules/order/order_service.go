package order

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	Checkout(userID string, request CheckoutRequest) (OrderResponse, error)
	GetMyOrders(userID string) ([]OrderResponse, error)
	GetMyOrderByID(userID string, orderID string) (OrderResponse, error)
	GetMyOrderByNumber(userID string, orderNumber string) (OrderResponse, error)
	UpdateOrderStatus(orderID string, request UpdateOrderStatusRequest) (OrderResponse, error)
}

type orderService struct {
	orderRepository OrderRepository
}

func NewOrderService(orderRepository OrderRepository) OrderService {
	return &orderService{orderRepository: orderRepository}
}

func (s *orderService) Checkout(userID string, request CheckoutRequest) (OrderResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return OrderResponse{}, utils.BadRequest("Invalid user ID")
	}

	addressID, err := uuid.Parse(request.AddressID)

	if err != nil {
		return OrderResponse{}, utils.BadRequest("Invalid address ID")
	}

	cartItems, err := s.orderRepository.FindCartItems(parsedUserID)

	if err != nil {
		return OrderResponse{}, utils.InternalServerError("Failed to get cart items")
	}

	if len(cartItems) == 0 {
		return OrderResponse{}, utils.BadRequest("Cart is empty")
	}

	var createdOrder Order

	err = s.orderRepository.Transaction(func(tx *gorm.DB) error {
		subtotal := 0

		for _, item := range cartItems {
			productData, err := s.orderRepository.FindProductForUpdate(tx, item.ProductID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.NotFound("Product not found")
				}

				return utils.InternalServerError("Failed to get product")
			}

			if item.Quantity > productData.Stock {
				return utils.BadRequest("Product stock is not enough")
			}

			subtotal += productData.Price * item.Quantity
		}

		shipping := 15000
		total := subtotal + shipping

		newOrder := Order{
			OrderNumber:    generateOrderNumber(),
			UserID:         parsedUserID,
			AddressID:      addressID,
			Subtotal:       subtotal,
			Shipping:       shipping,
			Total:          total,
			OrderStatus:    OrderPending,
			PaymentStatus:  PaymentPending,
			ShippingStatus: ShippingPending,
			TrackingNumber: generateTrackingNumber(),
		}

		orderData, err := s.orderRepository.CreateOrderWithTx(tx, newOrder)

		if err != nil {
			return utils.InternalServerError("Failed to create order")
		}

		for _, item := range cartItems {
			productData, err := s.orderRepository.FindProductForUpdate(tx, item.ProductID)

			if err != nil {
				return utils.InternalServerError("Failed to get product")
			}

			productData.Stock -= item.Quantity

			if productData.Stock <= 0 {
				productData.StockStatus = "out_of_stock"
			} else if productData.Stock <= 5 {
				productData.StockStatus = "low_stock"
			} else {
				productData.StockStatus = "in_stock"
			}

			err = s.orderRepository.UpdateProductStockWithTx(tx, productData)

			if err != nil {
				return utils.InternalServerError("Failed to update product stock")
			}

			orderItem := OrderItem{
				OrderID:      orderData.ID,
				ProductID:    item.ProductID,
				ProductName:  item.Product.Name,
				ProductSlug:  item.Product.Slug,
				ProductImage: item.Product.Image,
				Price:        item.Product.Price,
				Quantity:     item.Quantity,
				Subtotal:     item.Product.Price * item.Quantity,
			}

			err = s.orderRepository.CreateOrderItemWithTx(tx, orderItem)

			if err != nil {
				return utils.InternalServerError("Failed to create order item")
			}
		}

		tracking := OrderTracking{
			OrderID:     orderData.ID,
			Status:      "Order Created",
			Description: "Your order has been created and waiting for payment.",
			Location:    "Reflect Store",
		}

		err = s.orderRepository.CreateTrackingWithTx(tx, tracking)

		if err != nil {
			return utils.InternalServerError("Failed to create tracking")
		}

		err = s.orderRepository.ClearCartWithTx(tx, parsedUserID)

		if err != nil {
			return utils.InternalServerError("Failed to clear cart")
		}

		createdOrder = orderData

		return nil
	})

	if err != nil {
		return OrderResponse{}, err
	}

	orderData, err := s.orderRepository.FindByID(createdOrder.ID)

	if err != nil {
		return OrderResponse{}, utils.InternalServerError("Failed to get created order")
	}

	return ToOrderResponse(orderData), nil
}

func (s *orderService) GetMyOrders(userID string) ([]OrderResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return nil, utils.BadRequest("Invalid user ID")
	}

	orders, err := s.orderRepository.FindByUserID(parsedUserID)

	if err != nil {
		return nil, utils.InternalServerError("Failed to get orders")
	}

	return ToOrderResponses(orders), nil
}

func (s *orderService) GetMyOrderByID(userID string, orderID string) (OrderResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return OrderResponse{}, utils.BadRequest("Invalid user ID")
	}

	parsedOrderID, err := uuid.Parse(orderID)

	if err != nil {
		return OrderResponse{}, utils.BadRequest("Invalid order ID")
	}

	orderData, err := s.orderRepository.FindByID(parsedOrderID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderResponse{}, utils.NotFound("Order not found")
		}

		return OrderResponse{}, utils.InternalServerError("Failed to get order")
	}

	if orderData.UserID != parsedUserID {
		return OrderResponse{}, utils.Forbidden("You are not allowed to access this order")
	}

	return ToOrderResponse(orderData), nil
}

func (s *orderService) GetMyOrderByNumber(userID string, orderNumber string) (OrderResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return OrderResponse{}, utils.BadRequest("Invalid user ID")
	}

	orderData, err := s.orderRepository.FindByOrderNumber(orderNumber)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderResponse{}, utils.NotFound("Order not found")
		}

		return OrderResponse{}, utils.InternalServerError("Failed to get order")
	}

	if orderData.UserID != parsedUserID {
		return OrderResponse{}, utils.Forbidden("You are not allowed to access this order")
	}

	return ToOrderResponse(orderData), nil
}

func (s *orderService) UpdateOrderStatus(orderID string, request UpdateOrderStatusRequest) (OrderResponse, error) {
	parsedOrderID, err := uuid.Parse(orderID)

	if err != nil {
		return OrderResponse{}, utils.BadRequest("Invalid order ID")
	}

	orderData, err := s.orderRepository.FindByID(parsedOrderID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderResponse{}, utils.NotFound("Order not found")
		}

		return OrderResponse{}, utils.InternalServerError("Failed to get order")
	}

	orderData.OrderStatus = request.OrderStatus

	if request.PaymentStatus != "" {
		orderData.PaymentStatus = request.PaymentStatus
	}

	if request.ShippingStatus != "" {
		orderData.ShippingStatus = request.ShippingStatus
	}

	updatedOrder, err := s.orderRepository.UpdateOrder(orderData)

	if err != nil {
		return OrderResponse{}, utils.InternalServerError("Failed to update order")
	}

	if request.TrackingStatus != "" {
		tracking := OrderTracking{
			OrderID:     updatedOrder.ID,
			Status:      request.TrackingStatus,
			Description: request.Description,
			Location:    request.Location,
		}

		err = s.orderRepository.Transaction(func(tx *gorm.DB) error {
			return s.orderRepository.CreateTrackingWithTx(tx, tracking)
		})

		if err != nil {
			return OrderResponse{}, utils.InternalServerError("Failed to create tracking update")
		}
	}

	finalOrder, err := s.orderRepository.FindByID(updatedOrder.ID)

	if err != nil {
		return OrderResponse{}, utils.InternalServerError("Failed to get updated order")
	}

	return ToOrderResponse(finalOrder), nil
}

func generateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf(
		"BLCK%d%d",
		time.Now().Unix(),
		rand.Intn(999),
	)
}

func generateTrackingNumber() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf(
		"TRK%d%d",
		time.Now().Unix(),
		rand.Intn(9999),
	)
}