package repository

import "gorm.io/gorm"

type Repositories struct {
	User *UserRepository
}

// 初始化所有 repository
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}
