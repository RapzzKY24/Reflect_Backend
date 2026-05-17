package product

type CreateProductRequest struct {
	Name         string          `json:"name" binding:"required"`
	Slug         string          `json:"slug" binding:"required"`
	Description  string          `json:"description"`
	Category     ProductCategory `json:"category" binding:"required"`
	Collection   string          `json:"collection"`
	Image        string          `json:"image"`
	Price        int             `json:"price" binding:"required"`
	Stock        int             `json:"stock"`
	IsNewArrival bool            `json:"is_new_arrival"`
	IsFeatured   bool            `json:"is_featured"`
}

type UpdateProductRequest struct {
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Category     ProductCategory `json:"category"`
	Collection   string          `json:"collection"`
	Image        string          `json:"image"`
	Price        int             `json:"price"`
	Stock        int             `json:"stock"`
	IsNewArrival bool            `json:"is_new_arrival"`
	IsFeatured   bool            `json:"is_featured"`
}