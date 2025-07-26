package routes

import (
	"github.com/LeeChasel/shareVault/backend/middleware"
	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api/" + apiVersion)

	RegisterAuthRoutes(api.Group("/auth"))

	filesGroup := api.Group("/files")
	filesGroup.Use(middleware.JWTAuthMiddleware())
	RegisterFilesRoutes(filesGroup)
}
