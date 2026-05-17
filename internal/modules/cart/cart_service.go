package cart

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService interface {
	GetMyCart(userID string) (CartResponse, error)
	AddToCart(userID string, request AddToCartRequest) (CartResponse, error)
	UpdateCartItem(userID string, cartItemID string, request UpdateCartItemRequest) (CartResponse, error)
	RemoveCartItem(userID string, cartItemID string) error
	ClearCart(userID string) error
}

type cartService struct {
	cartRepository CartRepository
}

func NewCartService(cartRepository CartRepository) CartService {
	return &cartService{cartRepository: cartRepository}
}

func (s *cartService) GetMyCart(userID string) (CartResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return CartResponse{}, utils.BadRequest("Invalid user ID")
	}

	items, err := s.cartRepository.FindByUserID(parsedUserID)

	if err != nil {
		return CartResponse{}, utils.InternalServerError("Failed to get cart")
	}

	return ToCartResponse(items), nil
}

func (s *cartService) AddToCart(userID string, request AddToCartRequest) (CartResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return CartResponse{}, utils.BadRequest("Invalid user ID")
	}

	productID, err := uuid.Parse(request.ProductID)

	if err != nil {
		return CartResponse{}, utils.BadRequest("Invalid product ID")
	}

	if request.Quantity <= 0 {
		return CartResponse{}, utils.BadRequest("Quantity must be greater than zero")
	}

	err = s.cartRepository.Transaction(func(tx *gorm.DB) error {
		productData, err := s.cartRepository.FindProductForUpdate(tx, productID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.NotFound("Product not found")
			}

			return utils.InternalServerError("Failed to get product")
		}

		if productData.Stock <= 0 {
			return utils.BadRequest("Product is out of stock")
		}

		existingItem, err := s.cartRepository.FindItemByUserAndProduct(parsedUserID, productID)

		if err == nil {
			newQuantity := existingItem.Quantity + request.Quantity

			if newQuantity > productData.Stock {
				return utils.BadRequest("Requested quantity exceeds available stock")
			}

			existingItem.Quantity = newQuantity

			_, err = s.cartRepository.UpdateWithTx(tx, existingItem)

			if err != nil {
				return utils.InternalServerError("Failed to update cart item")
			}

			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.InternalServerError("Failed to check cart item")
		}

		if request.Quantity > productData.Stock {
			return utils.BadRequest("Requested quantity exceeds available stock")
		}

		newItem := CartItem{
			UserID:    parsedUserID,
			ProductID: productID,
			Quantity:  request.Quantity,
		}

		_, err = s.cartRepository.CreateWithTx(tx, newItem)

		if err != nil {
			return utils.InternalServerError("Failed to add item to cart")
		}

		return nil
	})

	if err != nil {
		return CartResponse{}, err
	}

	return s.GetMyCart(userID)
}

func (s *cartService) UpdateCartItem(userID string, cartItemID string, request UpdateCartItemRequest) (CartResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return CartResponse{}, utils.BadRequest("Invalid user ID")
	}

	itemID, err := uuid.Parse(cartItemID)

	if err != nil {
		return CartResponse{}, utils.BadRequest("Invalid cart item ID")
	}

	if request.Quantity <= 0 {
		return CartResponse{}, utils.BadRequest("Quantity must be greater than zero")
	}

	err = s.cartRepository.Transaction(func(tx *gorm.DB) error {
		item, err := s.cartRepository.FindItemByID(itemID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.NotFound("Cart item not found")
			}

			return utils.InternalServerError("Failed to get cart item")
		}

		if item.UserID != parsedUserID {
			return utils.Forbidden("You are not allowed to update this cart item")
		}

		productData, err := s.cartRepository.FindProductForUpdate(tx, item.ProductID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.NotFound("Product not found")
			}

			return utils.InternalServerError("Failed to get product")
		}

		if request.Quantity > productData.Stock {
			return utils.BadRequest("Requested quantity exceeds available stock")
		}

		item.Quantity = request.Quantity

		_, err = s.cartRepository.UpdateWithTx(tx, item)

		if err != nil {
			return utils.InternalServerError("Failed to update cart item")
		}

		return nil
	})

	if err != nil {
		return CartResponse{}, err
	}

	return s.GetMyCart(userID)
}

func (s *cartService) RemoveCartItem(userID string, cartItemID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return utils.BadRequest("Invalid user ID")
	}

	itemID, err := uuid.Parse(cartItemID)

	if err != nil {
		return utils.BadRequest("Invalid cart item ID")
	}

	item, err := s.cartRepository.FindItemByID(itemID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Cart item not found")
		}

		return utils.InternalServerError("Failed to get cart item")
	}

	if item.UserID != parsedUserID {
		return utils.Forbidden("You are not allowed to remove this cart item")
	}

	err = s.cartRepository.DeleteByID(itemID)

	if err != nil {
		return utils.InternalServerError("Failed to remove cart item")
	}

	return nil
}

func (s *cartService) ClearCart(userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return utils.BadRequest("Invalid user ID")
	}

	err = s.cartRepository.DeleteByUserID(parsedUserID)

	if err != nil {
		return utils.InternalServerError("Failed to clear cart")
	}

	return nil
}