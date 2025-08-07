package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)


type OAuthController struct {
	config *oauth2.Config
	userUsecase *usecases.UserUsecases
}

func NewOAuthController(cfg *oauth2.Config, userUsecase *usecases.UserUsecases) *OAuthController{
	return &OAuthController{
		config: cfg,
		userUsecase: userUsecase,
	}
}

func (oa *OAuthController) OAuthHandler(c *gin.Context) {
	url := oa.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oa *OAuthController) OAuthCallBack(c *gin.Context) {

	ctx := c.Request.Context()
	code := c.Query("code")
	got, err := oa.config.Exchange(ctx, code)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format or content"})
		c.Abort()
		return
		
	}

	refreshTTL := 30 * 24 * time.Hour // 30 days of refresh token expiry
	now := time.Now()

	tokens := &domain.Token{
		AccessToken: got.AccessToken,
		RefreshToken: got.RefreshToken,
		AccessExpiry: got.Expiry,
		RefreshExpiry: now.Add(refreshTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}

	client := oa.config.Client(ctx, got)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		c.Abort()
		return
	}

	var userInfo *domain.UserInfo
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error while decoding response"})
		c.Abort()
		return
	}

	profile := domain.Profile{
		ProfilePic: userInfo.Picture,
		CreatedAt: now,
		UpdatedAt: now,
	}
 
	user := &domain.User{
		Firstname: userInfo.GivenName,
		Lastname: userInfo.FamilyName,
		Email: userInfo.Email,
		Role: "user",
		Provider: "google",
		CreatedAt: now,
		UpdatedAt: now,	

		Profile: profile,
	}

	log.Print("user email: ", user.Email)

	userID, err := oa.userUsecase.Register(ctx, user)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	tokens.UserID = userID
	err = oa.userUsecase.SaveToken(ctx, tokens)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return	
	}

	
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token": tokens,
	})

	log.Printf("OAuth login success: %s (user ID: %s)", user.Email, userID)
	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	
}

func (oa *OAuthController) RefreshToken (c *gin.Context){

	ctx := c.Request.Context()
	var token *domain.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	expiredToken := &oauth2.Token{
		RefreshToken: token.RefreshToken,
		Expiry: time.Now().Add(-time.Hour),
	}

	tokenSource := oa.config.TokenSource(ctx, expiredToken)

	newToken, err := tokenSource.Token()
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "failed to refresh access token"})
		c.Abort()
		return
	}

	token.AccessToken = newToken.AccessToken
	token.AccessExpiry = newToken.Expiry
	token.UpdatedAt = time.Now()

	if newToken.RefreshToken != "" {
		token.RefreshToken = newToken.RefreshToken
	}

	err = oa.userUsecase.SaveToken(ctx, token)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	
	c.IndentedJSON(200, gin.H{
		"message": "token refresh successfully",
		"token" : token,
	})

	

}