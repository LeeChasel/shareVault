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

	configs.InitDB()
	db := configs.GetDB()
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)
	migrationError := db.AutoMigrate(&models.User{})
	if migrationError != nil {
		log.Fatalf("Migration error: %v", migrationError)
	}

	repos := repository.NewRepositories(db)

	r := gin.Default()
	r.Use(middleware.InjectRepos(repos))
	routes.InitRoutes(r)
	r.Run(":8080")
}
