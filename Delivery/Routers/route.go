package routers

import (
	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	repository "github.com/gedyzed/blog-starter-project/Repository"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func BlogRoutes(r *gin.Engine, blogCollection *mongo.Collection) {
	blogRepo := repository.NewBlogRepository(blogCollection)
	blogUsecase := usecases.NewBlogUsecase(blogRepo)
	BlogHandler := controllers.NewBlogHandler(blogUsecase)

	blogRoutes := r.Group("/blogs")
	{
		blogRoutes.PUT("/:id", BlogHandler.UpdateBlog)
		blogRoutes.DELETE("/:id", BlogHandler.DeleteBlog)
	}
}
