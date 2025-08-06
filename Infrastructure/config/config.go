package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {

	App  AppConfig 	  `mapstructure:"app" validate:"required"`
	Port  string      `mapstructure:"port" validate:"required,min=1,max=65535"`
	Mongo MongoConfig `mapstructure:"mongo" validate:"required"`
	Auth  AuthConfig  `mapstructure:"auth" validate:"required"`
	Email EmailConfig `mapstructure:"email" validate:"required"`
}

type MongoConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

type AppConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

type EmailConfig struct {
	AppPassword string `mapstructure:"app_password" validate:"required,min=10"`
	SenderEmail string `mapstructure:"sender_email" validate:"required,email"`
	SMTPHost    string `mapstructure:"smtp_host" validate:"required,hostname"`
	SMTPPort    string `mapstructure:"smtp_port" validate:"required,numeric"`
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
	viper.BindEnv("app.url", "APP_URL")
	viper.BindEnv("email.app_password", "EMAIL_APP_PASSWORD")
	viper.BindEnv("email.sender_email", "EMAIL_SENDER_EMAIL")
	viper.BindEnv("email.smtp_host", "EMAIL_SMTP_HOST")
	viper.BindEnv("email.smtp_port", "EMAIL_SMTP_PORT")

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


