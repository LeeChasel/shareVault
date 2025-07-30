package routes

import (
	"github.com/LeeChasel/shareVault/backend/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterFilesRoutes(r *gin.RouterGroup) {
	r.GET("", handlers.ListUserFiles)
	r.POST("/upload", handlers.UploadFiles)
} 