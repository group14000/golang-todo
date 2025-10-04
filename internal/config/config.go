package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURL string
	JWTSecret  string
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		log.Fatal("MONGODB_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &Config{
		MongoDBURL: mongoURL,
		JWTSecret:  jwtSecret,
	}
}
