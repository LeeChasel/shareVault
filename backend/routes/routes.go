package routes

import (
	"github.com/LeeChasel/sharevault/backend/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)
} 