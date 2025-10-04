package api

import (
	"github.com/gin-gonic/gin"
	"github.com/group14000/golang-todo/internal/handlers"
	"github.com/group14000/golang-todo/internal/middleware"
)

func SetupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, todoHandler *handlers.TodoHandler, aiHandler *handlers.AIHandler, authMW *middleware.AuthMiddleware) {
	// Public routes
	r.POST("/signup", authHandler.SignUp)
	r.POST("/verify-otp", authHandler.VerifyOTP)
	r.POST("/login", authHandler.Login)
	r.POST("/forgot-password", authHandler.ForgotPassword)
	r.POST("/reset-password", authHandler.ResetPassword)

	// Protected routes
	protected := r.Group("")
	protected.Use(authMW.Handler())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/ai/chat", aiHandler.Chat)
	}

	// Todo routes (protected)
	api := r.Group("/todos")
	api.Use(authMW.Handler())
	{
		api.POST("", todoHandler.Create)
		api.GET("", todoHandler.List)
		api.GET(":id", todoHandler.Get)
		api.PATCH(":id", todoHandler.Update)
		api.DELETE(":id", todoHandler.Delete)
	}
}
