package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURL    string
	JWTSecret     string
	EmailHost     string
	EmailPort     int
	EmailUser     string
	EmailPassword string
	EmailUseTLS   bool
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

	emailHost := os.Getenv("EMAIL_HOST")
	if emailHost == "" {
		log.Fatal("EMAIL_HOST environment variable is required")
	}

	emailPortStr := os.Getenv("EMAIL_PORT")
	if emailPortStr == "" {
		log.Fatal("EMAIL_PORT environment variable is required")
	}
	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		log.Fatal("EMAIL_PORT must be a valid integer")
	}

	emailUser := os.Getenv("EMAIL_HOST_USER")
	if emailUser == "" {
		log.Fatal("EMAIL_HOST_USER environment variable is required")
	}

	emailPassword := os.Getenv("EMAIL_HOST_PASSWORD")
	if emailPassword == "" {
		log.Fatal("EMAIL_HOST_PASSWORD environment variable is required")
	}

	emailUseTLS := os.Getenv("EMAIL_USE_TLS") == "True"

	return &Config{
		MongoDBURL:    mongoURL,
		JWTSecret:     jwtSecret,
		EmailHost:     emailHost,
		EmailPort:     emailPort,
		EmailUser:     emailUser,
		EmailPassword: emailPassword,
		EmailUseTLS:   emailUseTLS,
	}
}
