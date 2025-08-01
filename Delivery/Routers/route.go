package routers

import (
	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
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
