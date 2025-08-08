package routers

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	infrastructure "github.com/gedyzed/blog-starter-project/Infrastructure"
)

func RegisterBlogRoutes(r *gin.Engine, blogHandler *controllers.BlogHandler, commentHandler *controllers.CommentHandler) {

	blog := r.Group("/blogs")
	{
		blog.POST("/", blogHandler.CreateBlog)
		blog.GET("/", blogHandler.GetAllBlogs)
		blog.GET("/:id", blogHandler.GetBlogById)
		blog.PUT("/:id", blogHandler.UpdateBlog)
		blog.DELETE("/:id", blogHandler.DeleteBlog)
		blog.POST("/:id/like", blogHandler.LikeBlog)
		blog.POST("/:id/dislike", blogHandler.DislikeBlog)
		blog.GET("/filter", blogHandler.FilterBlogs)
		blog.GET("/search", blogHandler.SearchBlogs)
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

func RegisterUserRoutes(r *gin.Engine, handler *controllers.UserController, authMiddleware *infrastructure.AuthMiddleware) {

	users := r.Group("/users")

	{
		users.POST("/register", handler.RegisterUser)
		users.POST("/login", handler.Login)
		users.POST("/forgot-password", handler.ForgotPassword)
		users.POST("/reset-password", handler.ResetPassword)
	}

	protectedUser := r.Group("/users")
	protectedUser.Use(authMiddleware.IsLogin)
	{
		protectedUser.POST("/update-profile", handler.ProfileUpdate)
	}

	protectedAdmins := r.Group("/admins")
	protectedAdmins.Use(authMiddleware.IsLogin)
	{
		protectedAdmins.POST("/promote-demote-user", handler.PromoteDemoteUser)
	}

}

func RegisterTokenRoutes(r *gin.Engine, handler *controllers.TokenController) {

	tokens := r.Group("/tokens/")

	{
		tokens.POST("/send-vcode", handler.SendVerificationEmail)
	}
}

func RegisterOAuthRoutes(r *gin.Engine, handler *controllers.OAuthController) {

	oauth := r.Group("/oauth")

	{
		oauth.GET("/auth/login", handler.OAuthHandler)
		oauth.GET("/callback", handler.OAuthCallBack)
		oauth.POST("/refresh-token", handler.RefreshToken)
	}
}

func RegisterGenerativeAIRoutes(r *gin.Engine, handler *controllers.GenerativeAIController, authMiddleware *infrastructure.AuthMiddleware) {

	protectedAI:= r.Group("/ai")
	protectedAI.Use(authMiddleware.IsLogin)
	{
		protectedAI.POST("/generate", handler.GenerativeAI)
	}
}
