package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Port  string         `mapstructure:"port" validate:"required,min=1,max=65535"`
	Mongo MongoConfig `mapstructure:"mongo" validate:"required"`
	Auth  AuthConfig  `mapstructure:"auth" validate:"required"`
}

type MongoConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

type AuthConfig struct {
	AccessTokenKey  string `mapstructure:"access_token_key" validate:"required,min=10"`
	RefreshTokenKey string `mapstructure:"refresh_token_key" validate:"required,min=10"`
}

func ValidateConfig(cfg *Config) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	viper.BindEnv("mongo.url", "MONGO_URL")
	viper.BindEnv("auth.access_token_key", "AUTH_ACCESS_TOKEN_KEY")
	viper.BindEnv("auth.refresh_token_key", "AUTH_REFRESH_TOKEN_KEY")

	viper.SetDefault("port", "8080")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := ValidateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DbInit() *mongo.Database {
	// 1) Load the MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println("Mongo URI:", mongoURI) // For debug; remove in prod

	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	// 2) Configure the ServerAPI (required for Atlas)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// 3) Connect to MongoDB Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// 4) Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("✅ Successfully connected to MongoDB Atlas")

	// 5) Return the reference to your database
	return client.Database("BlogDB") // Change if your database name is different
}
