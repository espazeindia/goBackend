package middlewares

import (
	"net/http"
	"strings"

	"espazeBackend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Please provide a valid JWT token",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "The provided token is invalid or expired",
			})
			c.Abort()
			return
		}

		// Check if token is expired for all roles
		expired, err := utils.IsTokenExpired(tokenString)
		if err != nil || expired {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token expired",
				"message": "The provided token has expired",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("name", claims.Name)
		c.Set("role", claims.Role)

		c.Next()
	}
}
