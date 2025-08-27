package configs

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Env.DBHost, Env.DBUser, Env.DBPassword, Env.DBName, Env.DBPort)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	database.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)
	return database, nil
}
