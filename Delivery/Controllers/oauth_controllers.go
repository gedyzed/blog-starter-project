package controllers

import (
	"net/http"


	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OAuthController struct {
	config      *oauth2.Config
	services 	domain.IOAuthServices
}

func NewOAuthController(cfg *oauth2.Config, svs domain.IOAuthServices) *OAuthController {
	return &OAuthController{
		config: 	  cfg,
		services: 	 svs,
	}
}

func (oa *OAuthController) OAuthHandler(c *gin.Context) {
	url := oa.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oa *OAuthController) OAuthCallBack(c *gin.Context) {

	ctx := c.Request.Context()
	code := c.Query("code")

	tokens, err := oa.services.OAuthCallBack(ctx, code)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"token":        tokens,
		"redirect_url": "/dashboard",
	})

}

func (oa *OAuthController) RefreshToken(c *gin.Context) {

	ctx := c.Request.Context()
	var token *domain.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	token, err := oa.services.RefreshToken(ctx, token)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": "token refresh successfully",
		"token":   token,
	})
}

     