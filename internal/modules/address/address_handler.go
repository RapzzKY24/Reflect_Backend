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

// @Summary      Get my addresses
// @Description  Retrieve all saved addresses for the authenticated user
// @Tags         Addresses
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.SwaggerSuccessResponse{data=[]address.AddressResponse}
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /addresses [get]
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

// @Summary      Create address
// @Description  Create a new shipping address for the authenticated user
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateAddressRequest true "Create Address Request"
// @Success      201 {object} utils.SwaggerSuccessResponse{data=address.AddressResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Router       /addresses [post]
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

// @Summary      Update address
// @Description  Update an existing shipping address
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Address ID"
// @Param        request body UpdateAddressRequest true "Update Address Request"
// @Success      200 {object} utils.SwaggerSuccessResponse{data=address.AddressResponse}
// @Failure      400 {object} utils.SwaggerValidationErrorResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /addresses/{id} [put]
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

// @Summary      Delete address
// @Description  Delete a shipping address by ID
// @Tags         Addresses
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Address ID"
// @Success      200 {object} utils.SwaggerSuccessResponse
// @Failure      401 {object} utils.SwaggerErrorResponse
// @Failure      404 {object} utils.SwaggerErrorResponse
// @Router       /addresses/{id} [delete]
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