package routers

import (
	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	repository "github.com/gedyzed/blog-starter-project/Repository"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUp(db *mongo.Database, route *gin.Engine){

	public := route.Group("")
	RegisterUser(db, public)

}


func RegisterUser(db *mongo.Database, route *gin.RouterGroup){

	userRepo := repository.NewUserMongoRepo(db)
	passwordService := usecases.NewPasswordServices()
	tokenService := usecases.NewTokenServices()
	userUsecase:= usecases.NewUserUsecase(userRepo, passwordService, tokenService)
	userController := controllers.NewUserController(userUsecase)

	route.POST("/register", userController.RegisterUser)
}