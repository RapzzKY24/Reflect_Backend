package middleware

import (
	"net/http"
	"strings"

	"reflect-backend/internal/config"
	"reflect-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(utils.Unauthorized("Authorization header is required"))
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		if tokenString == authHeader {
			c.Error(utils.Unauthorized("Bearer token is required"))
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.Error(utils.Unauthorized("Invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.Error(utils.Unauthorized("Invalid token claims"))
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")

		if !exists {
			c.Error(utils.Unauthorized("User role not found"))
			c.Abort()
			return
		}

		role, ok := roleValue.(string)

		if !ok {
			c.Error(utils.Forbidden("Invalid user role"))
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.Error(utils.Forbidden("You are not allowed to access this resource"))
		c.Abort()
	}
}

func GetCurrentUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")

	if !exists {
		return ""
	}

	id, ok := userID.(string)

	if !ok {
		return ""
	}

	return id
}

func GetCurrentUserRole(c *gin.Context) string {
	roleValue, exists := c.Get("role")

	if !exists {
		return ""
	}

	role, ok := roleValue.(string)

	if !ok {
		return ""
	}

	return role
}

func AbortUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  "error",
		"message": message,
	})
	c.Abort()
}