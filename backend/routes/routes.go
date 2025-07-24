package routes

import (
	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api/" + apiVersion)

	RegisterAuthRoutes(api.Group("/auth"))
}
