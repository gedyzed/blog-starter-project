package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService domain.ITokenService
}

func (m *AuthMiddleware) IsLogin(c *gin.Context) {
	header := c.Request.Header["Authorization"]
	token := strings.Split(header[0], " ")[1]

	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization header",
		})
		c.Abort()
		return
	}

	user, err := m.tokenService.VerifyAccessToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"errro": "invalid token",
		})
		c.Abort()
		return
	}

	c.Set("userID", user)
	c.Next()
}
