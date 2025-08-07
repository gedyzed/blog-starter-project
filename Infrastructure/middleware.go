package infrastructure

import (
	"net/http"
	"strings"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	TokenService domain.ITokenService
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
