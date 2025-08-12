package infrastructure

import (
	"fmt"
	"net/http"
	"strings"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	TokenService domain.ITokenService
	oauthService domain.IOAuthServices
	userUsecase  *usecases.UserUsecases
}

func NewAuthMiddleware(ts domain.ITokenService, os domain.IOAuthServices, uc *usecases.UserUsecases) *AuthMiddleware{
	 return &AuthMiddleware{
		TokenService: ts,
		oauthService: os,
		userUsecase: uc,
	}
}


func (m *AuthMiddleware) IsLogin(c *gin.Context) {

	ctx := c.Request.Context()
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing or invalid authorization header"})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	var userID string

	// First try local JWT verification
	uid, err := m.TokenService.VerifyAccessToken(token)
	if err == nil {
		userID = uid
	} else {
		
		// Fallback to Google OAuth2 verification
		resolvedID, err := m.oauthService.VerifyGoogleIDToken(ctx, token)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		userID = resolvedID
	}

	c.Set("userID", userID)
	c.Next()
}

// Add this function after your existing IsLogin function
func (m *AuthMiddleware) IsLoginWithRole() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
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

		var UserID string
		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := m.TokenService.VerifyAccessToken(token)
		if err == nil {
			UserID  = userID
		} else {

			// Fallback to Google OAuth2 verification
			resolvedID, err := m.oauthService.VerifyGoogleIDToken(ctx, token)
			fmt.Println("userID : ", resolvedID)
			if err != nil {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
				c.Abort()
				return
		}
		 UserID = resolvedID
		}

		// Get user role from database
		user, err := m.userUsecase.FindByUserID(ctx, UserID)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return	
		}


		// Set both userID and role in context
		c.Set("userID", UserID)
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
