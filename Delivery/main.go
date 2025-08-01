package main

import (
	"log"
	"os"

	routers "github.com/gedyzed/blog-starter-project/Delivery/Routers"
	"github.com/gedyzed/blog-starter-project/Infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found or failed to load it")
	}

	// Connect to DB
	db := config.DbInit()
	blogCollection := db.Collection("blogs")

	// Set up router
	r := gin.Default()
	routers.BlogRoutes(r, blogCollection)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
