package address

type CreateAddressRequest struct {
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`

	Province   string `json:"province" binding:"required"`
	City       string `json:"city" binding:"required"`
	District   string `json:"district" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`

	FullAddress string `json:"full_address" binding:"required"`

	Label string `json:"label"`

	IsPrimary bool `json:"is_primary"`
}

type UpdateAddressRequest struct {
	RecipientName string `json:"recipient_name"`
	PhoneNumber   string `json:"phone_number"`

	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	PostalCode string `json:"postal_code"`

	FullAddress string `json:"full_address"`

	Label string `json:"label"`

	IsPrimary *bool `json:"is_primary"`
}