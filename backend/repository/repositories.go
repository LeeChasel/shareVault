package repository

import (
	"github.com/LeeChasel/shareVault/backend/configs"
	"gorm.io/gorm"
)

type Repositories struct {
	User *UserRepository
	S3   *S3Repository
}

// 初始化所有 repository
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		S3:   NewS3Repository(configs.Env.AWSS3Bucket),
	}
}
