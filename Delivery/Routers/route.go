package routers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	"github.com/gedyzed/blog-starter-project/Infrastructure"
	"github.com/gedyzed/blog-starter-project/Repository"
	"github.com/gedyzed/blog-starter-project/Usecases"
)

func RegisterBlogRoutes(r *gin.Engine, handler *controllers.BlogHandler) {
	blog := r.Group("/blogs")
	{
		blog.POST("/", handler.CreateBlog)
		blog.GET("/", handler.GetAllBlogs)
		blog.GET("/:id", handler.GetBlogById)
		blog.PUT("/:id", handler.UpdateBlog)
		blog.DELETE("/:id", handler.DeleteBlog)
	}
}

func SetUp(db *mongo.Database, route *gin.Engine) {

	public := route.Group("")
	RegisterUser(db, public)

}

func RegisterUser(db *mongo.Database, route *gin.RouterGroup) {

	userRepo := repository.NewUserMongoRepo(db)
	passwordService := infrastructure.NewPasswordService()
	tokenService := infrastructure.NewJWTTokenService("TODO")
	userUsecase := usecases.NewUserUsecase(userRepo, tokenService, passwordService)
	userController := controllers.NewUserController(userUsecase)

	route.POST("/register", userController.RegisterUser)
	route.POST("/login", userController.Login)
}
