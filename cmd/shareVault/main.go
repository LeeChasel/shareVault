package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LeeChasel/shareVault/configs"
	"github.com/LeeChasel/shareVault/internal/api/routes"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/LeeChasel/shareVault/internal/repository"
	"github.com/LeeChasel/shareVault/internal/service"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

const PORT = 8083

func main() {
	configs.LoadEnv()
	awsCfg, err := configs.InitAWSConfig(context.Background())
	if err != nil {
		log.Fatalf("無法載入 AWS 設定: %v", err)
	}

	db, err := configs.InitDB()
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	// Repository
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)
	s3Repo := repository.NewS3Repository(s3.NewFromConfig(awsCfg), configs.Env.AWSS3Bucket)

	// Service
	services := &service.ApplicationServices{
		UserService: service.NewUserService(userRepo),
		FileService: service.NewFileService(fileRepo),
		S3Service:   service.NewS3Service(s3Repo),
	}

	r := gin.Default()
	routes.InitRoutes(r, services)
	r.Run(fmt.Sprintf(":%d", PORT))
}
