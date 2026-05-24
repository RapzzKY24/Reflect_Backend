package utils

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"total_items"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

type SwaggerSuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SwaggerErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type SwaggerValidationErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessPaginatedResponse(c *gin.Context, statusCode int, message string, items interface{}, totalItems int64, page int, limit int) {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data: PaginatedData{
			Items:      items,
			TotalItems: totalItems,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Error:   err,
	})
}

func ValidationErrorResponse(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Status:  "error",
		Message: "Validation error",
		Error:   err,
	})
}