package routes

import (
	"github.com/LeeChasel/shareVault/backend/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
}