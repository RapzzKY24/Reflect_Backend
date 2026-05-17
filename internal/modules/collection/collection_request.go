package collection

type CreateCollectionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=120"`
	Slug        string `json:"slug" binding:"required,min=2,max=150"`
	Description string `json:"description"`
	Image       string `json:"image"`
	IsFeatured  bool   `json:"is_featured"`
	IsActive    bool   `json:"is_active"`
}

type UpdateCollectionRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Image       string `json:"image"`
	IsFeatured  *bool  `json:"is_featured"`
	IsActive    *bool  `json:"is_active"`
}