package controllers

import (
	"errors"
	"log"
	"net/http"
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

func (uc *UserController) RefreshToken(c *gin.Context) {
	var requestBody struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	user, err := uc.userUsecase.GetByEmail(c.Request.Context(), requestBody.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
	}

	if user.Email != requestBody.Email {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	token, err := uc.userUsecase.RefreshToken(c.Request.Context(), user.ID.Hex(), requestBody.RefreshToken)
	if err != nil {
		switch err {
		case usecases.ErrExpiredRefreshToken:
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "expired refresh token"})
			c.Abort()
			return
		case domain.ErrTokenNotFound:
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "token not found"})
			c.Abort()
			return
		case domain.ErrInvalidToken:
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			c.Abort()
			return
		default:
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "login successfully",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
func (uc *UserController) Logout(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "include the username as a path parameter"})
		c.Abort()
		return
	}

	err := uc.userUsecase.Logout(c.Request.Context(), username)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		case usecases.ErrInvalidCredential:
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		case domain.ErrTokenNotFound:
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "session does not exist"})
			c.Abort()
			return
		default:
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "logout successfully"})
}

func (uc *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	// checking for required fields
	if requestBody.Username == "" || requestBody.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "fill all required fields: username, password"})
		c.Abort()
		return
	}

	var user = domain.User{
		Username: requestBody.Username,
		Password: requestBody.Password,
	}

	token, err := uc.userUsecase.Login(ctx, user)
	if err != nil {
		switch err {
		case usecases.ErrInvalidCredential:
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		case domain.ErrUserNotFound:
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": "invalid credential",
			})
		default:
			log.Printf("error %s\n", err.Error())
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

	_, err := uc.userUsecase.VerifyCode(ctx, user.VCode)
	if err != nil {
		log.Println(err.Error())
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	_, err = uc.userUsecase.Register(ctx, user)
	if err != nil {
		switch err.Error() {
		case "username already exists":
			c.IndentedJSON(409, gin.H{"error": err.Error()})
		case "email already exists":
			c.IndentedJSON(409, gin.H{"error": err.Error()})
		default:
			log.Println(err.Error())
			c.IndentedJSON(500, gin.H{"error": err.Error()})
		}

		c.Abort()
		return
	}

	err = uc.userUsecase.DeleteVCode(ctx, user.Email)
	if err != nil {
		log.Println(err.Error())
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "user created successfully"})
}

func (uc *UserController) ForgotPassword(c *gin.Context) {

	ctx := c.Request.Context()

	var user *domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	if user.Email == "" {
		c.IndentedJSON(400, gin.H{"error": "fill all required fields: email, password(new), vcode."})
		c.Abort()
		return
	}

	err := uc.userUsecase.ForgotPassword(ctx, user.Email)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "reset link has been sent to your email. please check your email and reset your password"})
}

func (uc *UserController) ResetPassword(c *gin.Context) {

	ctx := c.Request.Context()
	token := c.Query("token")

	var user *domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	if user.Password == "" {
		c.IndentedJSON(400, gin.H{"error": "please insert new password"})
		c.Abort()
		return
	}

	email, err := uc.userUsecase.VerifyCode(ctx, token)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err = uc.userUsecase.ResetPassword(ctx, email, user.Password)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Password Reset Successful"})

}

func (uc *UserController) PromoteDemoteUser(c *gin.Context) {

	ctx := c.Request.Context()

	var promteDemote *domain.PromoteDemoteStruct
	if err := c.ShouldBindJSON(&promteDemote); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	if promteDemote.UserID == "" {
		c.IndentedJSON(400, gin.H{"error": "user Id required"})
		c.Abort()
		return
	}

	err := uc.userUsecase.PromoteDemote(ctx, promteDemote.UserID)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Role has been updated successfully"})
}

func (uc *UserController) ProfileUpdate(c *gin.Context) {

	ctx := c.Request.Context()
	var profileUpdate domain.ProfileUpdateInput
	if err := c.ShouldBindJSON(&profileUpdate); err != nil {
		c.IndentedJSON(400, gin.H{"error": "invalid input format"})
		c.Abort()
		return
	}

	if profileUpdate.UserID == "" {
		c.IndentedJSON(400, gin.H{"error": "user Id required"})
		c.Abort()
		return
	}

	pdi := domain.ProfileUpdateInput{}
	if profileUpdate == pdi {
		c.IndentedJSON(400, gin.H{"error": "No profile field to be updated"})
		c.Abort()
		return
	}

	err := uc.userUsecase.ProfileUpdate(ctx, &profileUpdate)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Profile has been updated successfully"})
}
