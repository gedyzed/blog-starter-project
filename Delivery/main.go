package main

import (
	"fmt"
	"time"

	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	routers "github.com/gedyzed/blog-starter-project/Delivery/Routers"
	infrastructure "github.com/gedyzed/blog-starter-project/Infrastructure"
	"github.com/gedyzed/blog-starter-project/Infrastructure/config"
	repository "github.com/gedyzed/blog-starter-project/Repository"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
)

func main() {

	conf, err := config.LoadConfig()
	
	if err != nil {
		fmt.Println(err)
		return
	}

	db := infrastructure.DbInit(conf.Mongo.URL)


  ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// setup collections
	blogCollection := db.Collection("blogs")
	commentCollection := db.Collection("comments")
	userCollection := db.Collection("users")
	tokenCollection := db.Collection("tokens")
	vtokenCollection := db.Collection("vtokens")

	// Setup repo
	tokenRepo := repository.NewMongoTokenRepository(tokenCollection)
	vtokenRepo := repository.NewMongoVTokenRepository(vtokenCollection)
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
	tokenService := infrastructure.NewTokenService(conf.Email, conf.App.URL)
	jwtService := infrastructure.NewJWTTokenService(
		tokenRepo,
		conf.Auth.AccessTokenKey,
		conf.Auth.RefreshTokenKey,
		30*(24*time.Hour), // 1 month
		60*(24*time.Hour), // 2 month
	)

	// Setup usecases
	tokenUsecase := usecases.NewTokenUsecase(tokenRepo, vtokenRepo, tokenService, jwtService)
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
