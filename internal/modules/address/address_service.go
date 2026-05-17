package address

import (
	"errors"

	"reflect-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressService interface {
	GetMyAddresses(userID string) ([]AddressResponse, error)
	CreateAddress(userID string, request CreateAddressRequest) (AddressResponse, error)
	UpdateAddress(userID string, addressID string, request UpdateAddressRequest) (AddressResponse, error)
	DeleteAddress(userID string, addressID string) error
}

type addressService struct {
	addressRepository AddressRepository
}

func NewAddressService(addressRepository AddressRepository) AddressService {
	return &addressService{
		addressRepository: addressRepository,
	}
}

func (s *addressService) GetMyAddresses(userID string) ([]AddressResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return nil, utils.BadRequest("Invalid user ID")
	}

	addresses, err := s.addressRepository.FindByUserID(parsedUserID)

	if err != nil {
		return nil, utils.InternalServerError("Failed to get addresses")
	}

	return ToAddressResponses(addresses), nil
}

func (s *addressService) CreateAddress(userID string, request CreateAddressRequest) (AddressResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return AddressResponse{}, utils.BadRequest("Invalid user ID")
	}

	if request.IsPrimary {
		err = s.addressRepository.ClearPrimaryAddresses(parsedUserID)

		if err != nil {
			return AddressResponse{}, utils.InternalServerError("Failed to update primary address")
		}
	}

	newAddress := Address{
		UserID:        parsedUserID,
		RecipientName: request.RecipientName,
		PhoneNumber:   request.PhoneNumber,
		Province:      request.Province,
		City:          request.City,
		District:      request.District,
		PostalCode:    request.PostalCode,
		FullAddress:   request.FullAddress,
		Label:         request.Label,
		IsPrimary:     request.IsPrimary,
	}

	addressData, err := s.addressRepository.Create(newAddress)

	if err != nil {
		return AddressResponse{}, utils.InternalServerError("Failed to create address")
	}

	return ToAddressResponse(addressData), nil
}

func (s *addressService) UpdateAddress(userID string, addressID string, request UpdateAddressRequest) (AddressResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return AddressResponse{}, utils.BadRequest("Invalid user ID")
	}

	parsedAddressID, err := uuid.Parse(addressID)

	if err != nil {
		return AddressResponse{}, utils.BadRequest("Invalid address ID")
	}

	addressData, err := s.addressRepository.FindByID(parsedAddressID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AddressResponse{}, utils.NotFound("Address not found")
		}

		return AddressResponse{}, utils.InternalServerError("Failed to get address")
	}

	if addressData.UserID != parsedUserID {
		return AddressResponse{}, utils.Forbidden("You are not allowed to update this address")
	}

	if request.IsPrimary != nil && *request.IsPrimary {
		err = s.addressRepository.ClearPrimaryAddresses(parsedUserID)

		if err != nil {
			return AddressResponse{}, utils.InternalServerError("Failed to update primary address")
		}

		addressData.IsPrimary = true
	}

	if request.RecipientName != "" {
		addressData.RecipientName = request.RecipientName
	}

	if request.PhoneNumber != "" {
		addressData.PhoneNumber = request.PhoneNumber
	}

	if request.Province != "" {
		addressData.Province = request.Province
	}

	if request.City != "" {
		addressData.City = request.City
	}

	if request.District != "" {
		addressData.District = request.District
	}

	if request.PostalCode != "" {
		addressData.PostalCode = request.PostalCode
	}

	if request.FullAddress != "" {
		addressData.FullAddress = request.FullAddress
	}

	if request.Label != "" {
		addressData.Label = request.Label
	}

	updatedAddress, err := s.addressRepository.Update(addressData)

	if err != nil {
		return AddressResponse{}, utils.InternalServerError("Failed to update address")
	}

	return ToAddressResponse(updatedAddress), nil
}

func (s *addressService) DeleteAddress(userID string, addressID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return utils.BadRequest("Invalid user ID")
	}

	parsedAddressID, err := uuid.Parse(addressID)

	if err != nil {
		return utils.BadRequest("Invalid address ID")
	}

	addressData, err := s.addressRepository.FindByID(parsedAddressID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("Address not found")
		}

		return utils.InternalServerError("Failed to get address")
	}

	if addressData.UserID != parsedUserID {
		return utils.Forbidden("You are not allowed to delete this address")
	}

	err = s.addressRepository.Delete(parsedAddressID)

	if err != nil {
		return utils.InternalServerError("Failed to delete address")
	}

	return nil
}