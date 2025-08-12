package config

import (
	"errors"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Config struct {
	App   AppConfig   `mapstructure:"app" validate:"required"`
	Port  string      `mapstructure:"port" validate:"required,min=1,max=65535"`
	Mongo MongoConfig `mapstructure:"mongo" validate:"required"`
	Auth  AuthConfig  `mapstructure:"auth" validate:"required"`
	OAuth OAuthConfig `mapstructure:"oauth" validate:"required"`
	Email EmailConfig `mapstructure:"email" validate:"required"`
	AI    AIConfig    `mapstructure:"ai" validate:"required"`
}

type MongoConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

type AppConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

type AIConfig struct {
	ApiKey string `mapstructure:"api_key" validate:"required"`
}

type EmailConfig struct {
	AppPassword string `mapstructure:"app_password" validate:"required,min=10"`
	SenderEmail string `mapstructure:"sender_email" validate:"required,email"`
	SMTPHost    string `mapstructure:"smtp_host" validate:"required,hostname"`
	SMTPPort    string `mapstructure:"smtp_port" validate:"required,numeric"`
}

type OAuthConfig struct {
	ClientID     string         `mapstructure:"client_id" validate:"required"`
	ClientSecret string         `mapstructure:"client_secret" validate:"required"`
	Endpoint     oauth2.Endpoint `mapstructure:"endpoint" validate:"required"`
	RedirectURL  string         `mapstructure:"redirect_url" validate:"required"`
	Scopes       []string       `mapstructure:"scopes" validate:"required"`
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
	// Load .env only if it exists (local dev)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: error loading .env file: %v", err)
		}
	}

	// Load config.yaml if it exists (local dev)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config.yaml file found or error reading it: %v", err)
	}

	// Enable automatic env override
	viper.AutomaticEnv()

	// Bind env vars explicitly for nested fields (helps in Render)
	viper.BindEnv("mongo.url", "MONGO_URL")
	viper.BindEnv("auth.access_token_key", "AUTH_ACCESS_TOKEN_KEY")
	viper.BindEnv("auth.refresh_token_key", "AUTH_REFRESH_TOKEN_KEY")
	viper.BindEnv("app.url", "APP_URL")
	viper.BindEnv("email.app_password", "EMAIL_APP_PASSWORD")
	viper.BindEnv("email.sender_email", "EMAIL_SENDER_EMAIL")
	viper.BindEnv("email.smtp_host", "EMAIL_SMTP_HOST")
	viper.BindEnv("email.smtp_port", "EMAIL_SMTP_PORT")
	viper.BindEnv("oauth.client_id", "OAUTH_CLIENT_ID")
	viper.BindEnv("oauth.client_secret", "OAUTH_CLIENT_SECRET")
	viper.BindEnv("oauth.redirect_url", "OAUTH_REDIRECT_URL")
	viper.BindEnv("ai.api_key", "GEMINI_API_KEY")

	// Set defaults (including PORT)
	viper.SetDefault("port", "8080")
	viper.SetDefault("oauth.scopes", []string{"email", "profile"})

	// Unmarshal into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Validate
	if err := ValidateConfig(&cfg); err != nil {
		return nil, err
	}

	// Additional sanity check for PORT numeric range
	// (optional)
	portNum := cfg.Port
	if portNum == "" {
		return nil, errors.New("PORT environment variable is required")
	}

	return &cfg, nil
}
