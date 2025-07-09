package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URI     string
	Name    string
	Timeout time.Duration
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
	}

	expiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Name:    getEnv("DATABASE_NAME", "user_management_api"),
			Timeout: 10 * time.Second,
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "default_secret_key"),
			ExpiresIn: expiresIn,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
