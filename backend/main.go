package main

import (
	"log"

	"github.com/LeeChasel/shareVault/backend/configs"
	"github.com/LeeChasel/shareVault/backend/middleware"
	"github.com/LeeChasel/shareVault/backend/models"
	"github.com/LeeChasel/shareVault/backend/repository"
	"github.com/LeeChasel/shareVault/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	configs.LoadEnv()
	configs.InitAWSConfig()

	configs.InitDB()
	db := configs.GetDB()
	migrationError := db.AutoMigrate(models.AllModels()...)
	if migrationError != nil {
		log.Fatalf("Migration error: %v", migrationError)
	}

	repos := repository.NewRepositories(db)

	r := gin.Default()
	r.Use(middleware.InjectRepos(repos))
	routes.InitRoutes(r)
	r.Run(":8080")
}
