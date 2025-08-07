package routers

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	infrastructure "github.com/gedyzed/blog-starter-project/Infrastructure"
)

func RegisterBlogRoutes(r *gin.Engine, blogHandler *controllers.BlogHandler, commentHandler *controllers.CommentHandler, authMiddleware *infrastructure.AuthMiddleware) {
	blog := r.Group("/blogs")
	{
		// Public routes
		blog.GET("/", blogHandler.GetAllBlogs)
		blog.GET("/filter", blogHandler.FilterBlogs)
		blog.GET("/search", blogHandler.SearchBlogs)
		blog.GET("/:id", blogHandler.GetBlogById)

		// Protected routes
		blog.POST("/", authMiddleware.IsLogin, blogHandler.CreateBlog)
		blog.PUT("/:id", authMiddleware.IsLogin, blogHandler.UpdateBlog)
		blog.DELETE("/:id", authMiddleware.IsLogin, blogHandler.DeleteBlog)
		blog.POST("/:id/like", authMiddleware.IsLogin, blogHandler.LikeBlog)
		blog.POST("/:id/dislike", authMiddleware.IsLogin, blogHandler.DislikeBlog)
	}
	comments := r.Group("/comments")
	{
		comments.POST("/:blogId", authMiddleware.IsLogin, commentHandler.CreateComment)
		comments.GET("/:blogId", commentHandler.GetAllComments)
		comments.GET("/:blogId/:id", commentHandler.GetCommentByID)
		comments.PUT("/:blogId/:id", authMiddleware.IsLogin, commentHandler.EditComment)
		comments.DELETE("/:blogId/:id", authMiddleware.IsLogin, commentHandler.DeleteComment)
	}
}

func RegisterUserRoutes(r *gin.Engine, handler *controllers.UserController) {

	users := r.Group("/users")

	{
		users.POST("/register", handler.RegisterUser)
		users.POST("/login", handler.Login)
		users.POST("/forgot-password", handler.ForgotPassword)
		users.POST("/reset-password", handler.ResetPassword)
		users.POST("/update-profile", handler.ProfileUpdate)
	}

	admins := r.Group("/admins")

	{
		admins.POST("/promote-demote-user", handler.PromoteDemoteUser)
	}

}

func RegisterTokenRoutes(r *gin.Engine, handler *controllers.TokenController) {

	tokens := r.Group("/tokens/")

	{
		tokens.POST("/send-vcode", handler.SendVerificationEmail)
	}
}
