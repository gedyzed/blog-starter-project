package routers

import (
	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(r *gin.Engine, blogHandler *controllers.BlogHandler, commentHandler *controllers.CommentHandler){
	blog := r.Group("/blogs")
	{
		blog.POST("/", blogHandler.CreateBlog)         
		blog.GET("/", blogHandler.GetAllBlogs)        
		blog.GET("/:id", blogHandler.GetBlogById)      
		blog.PUT("/:id", blogHandler.UpdateBlog)   
		blog.DELETE("/:id", blogHandler.DeleteBlog)    
	}
	// Comment Routes
	r.POST("/comments/:blogId", commentHandler.CreateComment)
	r.GET("/comments/:blogId", commentHandler.GetAllComments)
	r.GET("/comments/:blogId/:id", commentHandler.GetCommentByID)
	r.PUT("/comments/:blogId/:id", commentHandler.EditComment)
	r.DELETE("/comments/:blogId/:id", commentHandler.DeleteComment)
}

