package routers

import (
	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
<<<<<<< HEAD
	"github.com/gin-gonic/gin"
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
=======
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
>>>>>>> 5326082c22b972240493a44a880a1835f7d591f8
