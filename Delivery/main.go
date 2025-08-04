package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	"github.com/gedyzed/blog-starter-project/Delivery/Routers"
	"github.com/gedyzed/blog-starter-project/Infrastructure"
	"github.com/gedyzed/blog-starter-project/Infrastructure/config"
	"github.com/gedyzed/blog-starter-project/Repository"
	"github.com/gedyzed/blog-starter-project/Usecases"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 1) Configure the ServerAPI (required for Atlas)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conf.Mongo.URL).SetServerAPIOptions(serverAPI)

	// 2) Connect to MongoDB Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// 3) Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("✅ Successfully connected to MongoDB Atlas")

	// 4) Return the reference to your database
	db := client.Database("BlogDB") // Change if your database name is different

	// setup collections
	blogCollection := db.Collection("blogs")
	commentCollection := db.Collection("comments")
	userCollection := db.Collection("users")
	tokenCollection := db.Collection("tokens")

	// Setup repo
	tokenRepo := repository.NewMongoTokenRepository(tokenCollection)
	userRepo := repository.NewMongoUserRepo(userCollection)
	blogRepo := repository.NewBlogRepository(blogCollection)
  commentRepo := repository.NewCommentRepository(commentCollection)



	// Setup services
	passService := infrastructure.NewPasswordService()
    tokenService := infrastructure.NewTokenServices()
	jwtService := infrastructure.NewJWTTokenService(
		tokenRepo,
		conf.Auth.AccessTokenKey,
		conf.Auth.RefreshTokenKey,
		30*(24*time.Hour), // 1 month
		60*(24*time.Hour), // 2 month
	)

    // Setup token Usecase
	tokenUsecase := usecases.NewTokenUsecase(tokenRepo, tokenService, jwtService)

	// Setup usecases
	userUsecase := usecases.NewUserUsecase(userRepo, tokenUsecase, passService)
	blogUsecase := usecases.NewBlogUsecase(blogRepo)
  commentUsecase := usecases.NewCommentUsecase(commentRepo)

	// Setup handlers
	userHandler := controllers.NewUserController(userUsecase)
	blogHandler := controllers.NewBlogHandler(blogUsecase)
  commentHandler := controllers.NewCommentHandler(commentUsecase)
  tokenHandler := controllers.NewTokenController(tokenUsecase)


	r := gin.Default()

	routers.RegisterBlogRoutes(r, blogHandler, commentHandler)
	routers.RegisterUserRoutes(r, userHandler)
  routers.RegisterTokenRoutes(r, tokenHandler)

	r.Run(":" + conf.Port)
}
