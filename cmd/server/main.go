package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/group14000/golang-todo/api"
	"github.com/group14000/golang-todo/internal/config"
	"github.com/group14000/golang-todo/internal/database"
	"github.com/group14000/golang-todo/internal/handlers"
	"github.com/group14000/golang-todo/internal/middleware"
	"github.com/group14000/golang-todo/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	client, err := database.Connect(cfg.MongoDBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	userRepo := database.NewUserRepository(client)
	otpRepo := database.NewOTPRepository(client)
	emailService := services.NewEmailService(cfg)
	authService := services.NewAuthService(userRepo, otpRepo, emailService, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// Todo dependencies
	todoRepo := database.NewTodoRepository(client)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret)

	r := gin.Default()
	api.SetupRoutes(r, authHandler, todoHandler, authMW)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
