package repository

import (
	"context"

	"github.com/LeeChasel/shareVault/backend/models"

	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(ctx context.Context, file *models.File) (*models.File, error) {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *FileRepository) GetByUserId(ctx context.Context, userId string) ([]*models.File, error) {
	var files []*models.File
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
