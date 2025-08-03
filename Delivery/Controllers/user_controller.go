package controllers

import (
	"net/http"
	"fmt"
	"regexp"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase *usecases.UserUsecases
}

func NewUserController(uc *usecases.UserUsecases) *UserController {
	return &UserController{userUsecase: uc}
}

func (uc *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()

	// accepting user input
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	// checking for required fields
	if user.Username != "" || user.Password != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "fill all required fields: username, password"})
		c.Abort()
		return
	}

	token, err := uc.userUsecase.Login(ctx, user)
	if err != nil {
		switch err {
		case usecases.ErrInvalidCredential:
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "login successfully",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (uc *UserController) RegisterUser(c *gin.Context) {

	ctx := c.Request.Context()

	// accepting user input
	var user *domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	// checking for required fields
	if user.Email == "" || user.Username == "" || user.Password == "" || user.Firstname == "" {
		c.IndentedJSON(400, gin.H{"error": "fill all required fields: email, username, password, firstname"})
		c.Abort()
		return
	}

	// Basic validation
	if len(user.Password) < 6 {
		c.IndentedJSON(400, gin.H{"error": "password must be at least 6 characters long"})
		c.Abort()
		return
	}

	if len(user.Username) < 3 {
		c.IndentedJSON(400, gin.H{"error": "username must be at least 3 characters long"})
		c.Abort()
		return
	}

	// Basic email validation 
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	match := re.MatchString(user.Email)
	if !match {
		c.IndentedJSON(400, gin.H{"error": "invalid email format"})
		c.Abort()
		return
	}

	if user.VCode == "" {
		c.IndentedJSON(401, gin.H{"error": "insert verification code from your email"})
		c.Abort()
		return
	}

	err := uc.userUsecase.VerifyCode(ctx, user.Email, user.VCode)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err = uc.userUsecase.Register(ctx, user)
	if err != nil {
		switch err.Error() {
		case "username already exists":
			fmt.Println("gb", "username")
			c.IndentedJSON(409, gin.H{"error": err.Error()})
		case "email already exists":
			fmt.Println("email")
			c.IndentedJSON(409, gin.H{"error": err.Error()})
		default:
			c.IndentedJSON(500, gin.H{"error": err.Error()})
		}

		c.Abort()
		return
	}

	err = uc.userUsecase.DeleteVCode(ctx, user.Email)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "user created successfully"})
}


