package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	routers "github.com/gedyzed/blog-starter-project/Delivery/Routers"
	infrastructure "github.com/gedyzed/blog-starter-project/Infrastructure"
	"github.com/gedyzed/blog-starter-project/Infrastructure/config"
	repository "github.com/gedyzed/blog-starter-project/Repository"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
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
	ctx, cancel := context.WithCancel(context.Background())
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
	commentRepo := repository.NewCommentRepository(commentCollection, blogCollection)

	//to initialize the indexes
	if err := blogRepo.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	dispatcher := infrastructure.NewBlogQueue()
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
	blogUsecase := usecases.NewBlogUsecase(blogRepo, commentRepo,dispatcher)
	commentUsecase := usecases.NewCommentUsecase(commentRepo,dispatcher)

	// Setup handlers
	userHandler := controllers.NewUserController(userUsecase)
	blogHandler := controllers.NewBlogHandler(blogUsecase)
	commentHandler := controllers.NewCommentHandler(commentUsecase)
	tokenHandler := controllers.NewTokenController(tokenUsecase)

	infrastructure.StartBlogRefreshWorker(ctx, blogUsecase)

	r := gin.Default()

	routers.RegisterBlogRoutes(r, blogHandler, commentHandler)
	routers.RegisterUserRoutes(r, userHandler)
	routers.RegisterTokenRoutes(r, tokenHandler)

	r.Run(":" + conf.Port)
}
