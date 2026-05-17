package middleware

import (
	"net/http"

	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Message, nil)
			return
		}

		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Internal server error",
			err.Error(),
		)
	}
}