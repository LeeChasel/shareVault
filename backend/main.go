package main

import (
	"github.com/LeeChasel/shareVault/backend/configs"
	"github.com/LeeChasel/shareVault/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	configs.LoadEnv()
	configs.InitDB()
	configs.GetDB().AutoMigrate()
	r := gin.Default()
	routes.InitRoutes(r)
	r.Run(":8080")
}