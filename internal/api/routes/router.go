package routes

import (
	"github.com/LeeChasel/shareVault/internal/api/handlers"
	"github.com/LeeChasel/shareVault/internal/api/middleware"
	"github.com/LeeChasel/shareVault/internal/service"
	"github.com/gin-gonic/gin"
)

type groupRoute struct {
	route    *gin.RouterGroup
	services *service.ApplicationServices
}

const apiVersion = "v1"

func (r *groupRoute) initAuthRoute() {
	services := r.services
	authRoutes := r.route.Group("/auth")

	authRoutes.POST("/register", handlers.Register(services))
	authRoutes.POST("/login", handlers.Login(services))
}

func (r *groupRoute) initFileRoute() {
	services := r.services
	filesGroup := r.route.Group("/files")
	filesGroup.Use(middleware.JWTAuthMiddleware())

	filesGroup.GET("", handlers.ListUserFiles(services))
	filesGroup.DELETE("", handlers.DeleteFileByIds(services))
	filesGroup.POST("/upload", handlers.UploadFiles(services))
	filesGroup.POST("/download", handlers.DownloadFiles(services))
}

func InitRoutes(r *gin.Engine, services *service.ApplicationServices) {
	api := r.Group("/api/" + apiVersion)

	groupRoute := groupRoute{
		route:    api,
		services: services,
	}

	groupRoute.initAuthRoute()
	groupRoute.initFileRoute()
}
