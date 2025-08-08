package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService domain.ITokenService
	oauthService domain.IOAuthServices
}

func NewAuthMiddleware(ts domain.ITokenService, os domain.IOAuthServices) *AuthMiddleware{
	 return &AuthMiddleware{
		tokenService: ts,
		oauthService: os,
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
	uid, err := m.tokenService.VerifyAccessToken(token)
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

