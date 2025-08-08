package infrastructure

import (
	"context"
	"net/http"
	"strings"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	TokenService   domain.ITokenService
	UserRepository domain.IUserRepository // Add this field
}

func (m *AuthMiddleware) IsLogin(c *gin.Context) {
	header := c.Request.Header["Authorization"]
	if len(header) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization header",
		})
		c.Abort()
		return
	}

	authHeader := header[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid authorization header format",
		})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization header",
		})
		c.Abort()
		return
	}

	user, err := m.TokenService.VerifyAccessToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		c.Abort()
		return
	}

	c.Set("userID", user)
	c.Next()
}

// Add this function after your existing IsLogin function
func (m *AuthMiddleware) IsLoginWithRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := m.TokenService.VerifyAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get user role from database
		user, err := m.UserRepository.Get(context.Background(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set both userID and role in context
		c.Set("userID", userID)
		c.Set("role", user.Role)
		c.Next()
	}
}

// Add this function for admin-only endpoints
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
