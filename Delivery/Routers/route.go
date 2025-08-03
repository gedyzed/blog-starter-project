package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/gedyzed/blog-starter-project/Delivery/Controllers"
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

func RegisterUserRoutes(r *gin.Engine, handler *controllers.UserController) {

	users := r.Group("/users")

	{
		users.POST("/register", handler.RegisterUser)
		users.POST("/login", handler.Login)
	}
}
