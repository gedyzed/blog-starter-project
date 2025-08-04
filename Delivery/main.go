package main

import (
	"log"
	"os"
	"github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	"github.com/gedyzed/blog-starter-project/Delivery/Routers"
	"github.com/gedyzed/blog-starter-project/Infrastructure/config"
	"github.com/gedyzed/blog-starter-project/Repository"
	"github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load it")
	}

	db := config.DbInit()
	blogCollection := db.Collection("blogs")
	commentCollection := db.Collection("comments")

	blogRepo := repository.NewBlogRepository(blogCollection)
	blogUsecase := usecases.NewBlogUsecase(blogRepo)
	blogHandler := controllers.NewBlogHandler(blogUsecase)
	commentRepo := repository.NewCommentRepository(commentCollection)
	commentUsecase := usecases.NewCommentUsecase(commentRepo)
	commentHandler := controllers.NewCommentHandler(commentUsecase)

	r := gin.Default()

	routers.RegisterBlogRoutes(r, blogHandler, commentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
