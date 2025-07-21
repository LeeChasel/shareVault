package main

import (
	"github.com/LeeChasel/sharevault/backend/configs"
	"github.com/LeeChasel/sharevault/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	configs.LoadEnv()
	r := gin.Default()
	routes.InitRoutes(r)
	r.Run(":8080")
}