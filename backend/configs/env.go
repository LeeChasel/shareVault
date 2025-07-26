package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	JwtSecret   string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	AWSRegion   string
	AWSS3Bucket string
}

var Env *EnvConfig

func LoadEnv() {
	_ = godotenv.Load()

	Env = &EnvConfig{
		JwtSecret:   mustGetEnv("JWT_SECRET"),
		DBHost:      mustGetEnv("DB_HOST"),
		DBPort:      mustGetEnv("DB_PORT"),
		DBUser:      mustGetEnv("DB_USER"),
		DBPassword:  mustGetEnv("DB_PASSWORD"),
		DBName:      mustGetEnv("DB_NAME"),
		AWSRegion:   mustGetEnv("AWS_REGION"),
		AWSS3Bucket: mustGetEnv("AWS_S3_BUCKET"),
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustGetEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Fatalf("Environment variable %s is required but not set", key)
	return ""
}
