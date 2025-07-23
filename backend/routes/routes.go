package routes

import (
	"github.com/LeeChasel/shareVault/backend/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)
}
