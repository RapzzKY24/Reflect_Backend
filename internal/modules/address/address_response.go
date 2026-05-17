package address

import (
	"time"

	"github.com/google/uuid"
)

type AddressResponse struct {
	ID            uuid.UUID `json:"id"`
	RecipientName string    `json:"recipient_name"`
	PhoneNumber   string    `json:"phone_number"`

	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	PostalCode string `json:"postal_code"`

	FullAddress string `json:"full_address"`

	Label string `json:"label"`

	IsPrimary bool `json:"is_primary"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAddressResponse(address Address) AddressResponse {
	return AddressResponse{
		ID:            address.ID,
		RecipientName: address.RecipientName,
		PhoneNumber:   address.PhoneNumber,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		PostalCode:    address.PostalCode,
		FullAddress:   address.FullAddress,
		Label:         address.Label,
		IsPrimary:     address.IsPrimary,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}
}

func ToAddressResponses(addresses []Address) []AddressResponse {
	responses := make([]AddressResponse, 0)

	for _, address := range addresses {
		responses = append(responses, ToAddressResponse(address))
	}

	return responses
}