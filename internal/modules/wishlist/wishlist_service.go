package wishlist

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistService interface {
	GetMyWishlist(userID string) ([]WishlistItemResponse, error)
	AddToWishlist(userID string, request AddWishlistRequest) ([]WishlistItemResponse, error)
	RemoveWishlistItem(userID string, wishlistItemID string) error
	RemoveWishlistByProduct(userID string, productID string) error
}

type wishlistService struct {
	wishlistRepository WishlistRepository
}

func NewWishlistService(wishlistRepository WishlistRepository) WishlistService {
	return &wishlistService{
		wishlistRepository: wishlistRepository,
	}
}

func (s *wishlistService) GetMyWishlist(userID string) ([]WishlistItemResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return nil, utils.BadRequest("Invalid user ID")
	}

	items, err := s.wishlistRepository.FindByUserID(parsedUserID)

	if err != nil {
		return nil, utils.InternalServerError("Failed to get wishlist")
	}

	return ToWishlistResponses(items), nil
}

func (s *wishlistService) AddToWishlist(userID string, request AddWishlistRequest) ([]WishlistItemResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return nil, utils.BadRequest("Invalid user ID")
	}

	productID, err := uuid.Parse(request.ProductID)

	if err != nil {
		return nil, utils.BadRequest("Invalid product ID")
	}

	_, err = s.wishlistRepository.FindByUserAndProduct(parsedUserID, productID)

	if err == nil {
		return nil, utils.Conflict("Product already exists in wishlist")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.InternalServerError("Failed to check wishlist item")
	}

	newItem := WishlistItem{
		UserID:    parsedUserID,
		ProductID: productID,
	}

	_, err = s.wishlistRepository.Create(newItem)

	if err != nil {
		return nil, utils.InternalServerError("Failed to add item to wishlist")
	}

	items, err := s.wishlistRepository.FindByUserID(parsedUserID)

	if err != nil {
		return nil, utils.InternalServerError("Failed to get wishlist")
	}

	return ToWishlistResponses(items), nil
}

func (s *wishlistService) RemoveWishlistItem(userID string, wishlistItemID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return utils.BadRequest("Invalid user ID")
	}

	itemID, err := uuid.Parse(wishlistItemID)

	if err != nil {
		return utils.BadRequest("Invalid wishlist item ID")
	}

	item, err := s.wishlistRepository.FindByID(itemID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Wishlist item not found")
		}

		return utils.InternalServerError("Failed to get wishlist item")
	}

	if item.UserID != parsedUserID {
		return utils.Forbidden("You are not allowed to remove this wishlist item")
	}

	err = s.wishlistRepository.Delete(itemID)

	if err != nil {
		return utils.InternalServerError("Failed to remove wishlist item")
	}

	return nil
}

func (s *wishlistService) RemoveWishlistByProduct(userID string, productID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return utils.BadRequest("Invalid user ID")
	}

	parsedProductID, err := uuid.Parse(productID)

	if err != nil {
		return utils.BadRequest("Invalid product ID")
	}

	err = s.wishlistRepository.DeleteByUserAndProduct(parsedUserID, parsedProductID)

	if err != nil {
		return utils.InternalServerError("Failed to remove wishlist item")
	}

	return nil
}