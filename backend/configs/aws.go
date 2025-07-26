package configs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var awsConfig aws.Config
var awsConfigLoaded bool
var S3Client *s3.Client

// 初始化 AWS 設定，應於專案啟動時呼叫
func InitAWSConfig() {
	if awsConfigLoaded {
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("無法載入 AWS 設定: %v", err)
	}
	awsConfig = cfg
	awsConfigLoaded = true

	S3Client = s3.NewFromConfig(cfg)
}

func GetAWSConfig() aws.Config {
	if !awsConfigLoaded {
		InitAWSConfig()
	}
	return awsConfig
}
