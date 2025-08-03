package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gedyzed/blog-starter-project/Usecases"
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	// checking for required fields
	if user.Email != "" || user.Username != "" || user.Password != "" || user.Firstname != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "fill all required fields: email, username, password, firstname"})
		c.Abort()
		return
	}

	err := uc.userUsecase.Register(ctx, user)

	if err != nil {
		switch err.Error() {
		case "username or email already exists":
			c.IndentedJSON(409, gin.H{"error": err.Error()})
		default:
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}
