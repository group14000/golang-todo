// @title           Golang Todo API
// @version         1.0
// @description     A Clean Architecture Todo API with OTP-based authentication, JWT authorization, and AI chat integration.
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Provide your JWT access token as: Bearer <token>
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

	// Swagger docs & handlers
	_ "github.com/group14000/golang-todo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// AI dependencies
	aiService := services.NewAIService(cfg.AIAPIKey)
	aiHandler := handlers.NewAIHandler(aiService)

	r := gin.Default()
	api.SetupRoutes(r, authHandler, todoHandler, aiHandler, authMW)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080 (swagger at /swagger/index.html)")
	r.Run(":8080")
}
