package controllers

import (
	
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)




type TokenController struct {
	tokenUsecase usecases.ITokenUsecase
}

func NewTokenController(ts usecases.ITokenUsecase) *TokenController {
	return &TokenController{tokenUsecase: ts}
}

// send email for email verification
func (ec *TokenController) SendVerificationEmail(c *gin.Context) {

	ctx := c.Request.Context()

	// accepting user input
	var emailRequest *domain.EmailRequest
	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	// call the function to create and send verification code
	tokenType := usecases.Email_Verification
	err := ec.tokenUsecase.CreateSendVCode(ctx, emailRequest.Email, tokenType)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Internal server error. Try again later"})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "we have send verifcation code to your email. Please check your email!"})
}

// send email for email verification
func (ec *TokenController) SendPasswordRestCode (c *gin.Context) {

	ctx := c.Request.Context()

	// accepting user input
	var emailRequest *domain.EmailRequest
	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	// call the function to create and send verification code
	tokenType := "password_reset"
	err := ec.tokenUsecase.CreateSendVCode(ctx, emailRequest.Email, tokenType)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Internal server error. Try again later"})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "we have send verifcation code to your email. Please check your email!"})
}