package api

import (
	"github.com/gin-gonic/gin"
	"github.com/group14000/golang-todo/internal/handlers"
)

func SetupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.Login)
}
