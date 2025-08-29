package repository

import (
	"context"

	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/LeeChasel/shareVault/internal/repository/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) interfaces.FileRepository {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) Create(ctx context.Context, file *models.File) (*models.File, error) {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *fileRepository) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*models.File, error) {
	var files []*models.File
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId.String()).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) GetByIds(ctx context.Context, fileIds []string) ([]*models.File, error) {
	if len(fileIds) == 0 {
		return nil, nil
	}

	var files []*models.File
	if err := r.db.WithContext(ctx).Where("id IN ?", fileIds).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}

func (r *fileRepository) DeleteByIds(ctx context.Context, fileIds []string) error {
	if len(fileIds) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Where("id IN ?", fileIds).Delete(&models.File{}).Error; err != nil {
		return err
	}

	return nil
}
