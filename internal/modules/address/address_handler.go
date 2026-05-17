package address

import (
	"net/http"

	"reflect-backend/internal/middleware"
	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressService AddressService
}

func NewAddressHandler(addressService AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

func (h *AddressHandler) GetMyAddresses(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	addresses, err := h.addressService.GetMyAddresses(userID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Addresses retrieved successfully",
		addresses,
	)
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var request CreateAddressRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	address, err := h.addressService.CreateAddress(userID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"Address created successfully",
		address,
	)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	addressID := c.Param("id")

	var request UpdateAddressRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	address, err := h.addressService.UpdateAddress(userID, addressID, request)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Address updated successfully",
		address,
	)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	addressID := c.Param("id")

	err := h.addressService.DeleteAddress(userID, addressID)

	if err != nil {
		c.Error(err)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Address deleted successfully",
		nil,
	)
}