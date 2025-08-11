package main

import (
	"context"
	"fmt"
	"log"
	"time"

	controllers "github.com/gedyzed/blog-starter-project/Delivery/Controllers"
	routers "github.com/gedyzed/blog-starter-project/Delivery/Routers"
	infrastructure "github.com/gedyzed/blog-starter-project/Infrastructure"
	config "github.com/gedyzed/blog-starter-project/Infrastructure/config"
	"github.com/gedyzed/blog-starter-project/Infrastructure/oauth"
	repository "github.com/gedyzed/blog-starter-project/Repository"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"github.com/gin-gonic/gin"
)

func main() {

	conf, err := config.LoadConfig()
	googleOauthConfig := oauth.NewGoogleOauthConfig(&conf.OAuth)
	
	if err != nil {
		log.Fatal("Cannot configure the amazon ses")
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	db := infrastructure.DbInit(conf.Mongo.URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	cacheSize := 100
	lruCache, err := infrastructure.NewLRUCache(cacheSize)
	if err != nil{
        log.Fatalf("Failed to initialize LRU cache: %v", err)
    }
	
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

	blogRepo := repository.NewBlogRepository(blogCollection, userRepo, lruCache.BlogCache(), lruCache. SortedBlogsCache())
	commentRepo := repository.NewCommentRepository(commentCollection, blogCollection, userRepo,  lruCache.CommentCache())

	//to initialize the indexes
	if err := blogRepo.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	dispatcher := infrastructure.NewBlogQueue()
	// Setup services
	passService := infrastructure.NewPasswordService()
	vtokenService := infrastructure.NewTokenService(conf.Email, conf.App.URL)
	tokenService := infrastructure.NewJWTTokenService(
		tokenRepo,
		conf.Auth.AccessTokenKey,
		conf.Auth.RefreshTokenKey,
		30*(24*time.Hour), // 1 month
		60*(24*time.Hour), // 2 month
	)
	

	// Setup usecases
	tokenUsecase := usecases.NewTokenUsecase(tokenRepo, vtokenRepo, vtokenService, tokenService)
	userUsecase := usecases.NewUserUsecase(userRepo, tokenUsecase, passService)

	blogUsecase := usecases.NewBlogUsecase(blogRepo, commentRepo, dispatcher)
	commentUsecase := usecases.NewCommentUsecase(commentRepo, dispatcher)

	// oauth servcive
	oauthService := oauth.NewOAuthServices(googleOauthConfig, userUsecase)

	// Setup handlers
	userHandler := controllers.NewUserController(userUsecase)
	blogHandler := controllers.NewBlogHandler(blogUsecase)
	commentHandler := controllers.NewCommentHandler(commentUsecase)
	tokenHandler := controllers.NewTokenController(tokenUsecase)
	oAuthHandler := controllers.NewOAuthController(googleOauthConfig, oauthService)
	genAIHandler := controllers.NewGenerativeAIController(&conf.AI)

	// middlewares 
	authMiddleware := infrastructure.NewAuthMiddleware(tokenService, oauthService, userUsecase)

	infrastructure.StartBlogRefreshWorker(ctx, blogUsecase)

	r := gin.Default()

	routers.RegisterUserRoutes(r, userHandler, authMiddleware)
	routers.RegisterTokenRoutes(r, tokenHandler, )
	routers.RegisterOAuthRoutes(r,  oAuthHandler)
	routers.RegisterGenerativeAIRoutes(r, genAIHandler, authMiddleware)
	routers.RegisterBlogRoutes(r, blogHandler, commentHandler, authMiddleware)

	r.Run(":" + conf.Port)
}
