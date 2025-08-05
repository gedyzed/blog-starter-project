package routers

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
)

func RegisterBlogRoutes(r *gin.Engine, blogHandler *controllers.BlogHandler, commentHandler *controllers.CommentHandler) {
	blog := r.Group("/blogs")
	{
		blog.POST("/", blogHandler.CreateBlog)
		blog.GET("/", blogHandler.GetAllBlogs)
		blog.GET("/:id", blogHandler.GetBlogById)
		blog.PUT("/:id", blogHandler.UpdateBlog)
		blog.DELETE("/:id", blogHandler.DeleteBlog)
		blog.PUT("/:id/like", blogHandler.LikeBlog)
		blog.PUT("/:id/dislike", blogHandler.DislikeBlog)

	}
	comments := r.Group("/comments")
	{
		comments.POST("/:blogId", commentHandler.CreateComment)
		comments.GET("/:blogId", commentHandler.GetAllComments)
		comments.GET("/:blogId/:id", commentHandler.GetCommentByID)
		comments.PUT("/:blogId/:id", commentHandler.EditComment)
		comments.DELETE("/:blogId/:id", commentHandler.DeleteComment)
	}
}

func RegisterUserRoutes(r *gin.Engine, handler *controllers.UserController) {

	users := r.Group("/users")

	{
		users.POST("/register", handler.RegisterUser)
		users.POST("/login", handler.Login)
	}
}

func RegisterTokenRoutes(r *gin.Engine, handler *controllers.TokenController) {

	tokens := r.Group("/tokens/")

	{
		tokens.POST("/send-vcode", handler.SendVerificationEmail) // send verification email
	}
}
