package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
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