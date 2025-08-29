package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/LeeChasel/shareVault/internal/api/dto"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/google/uuid"
)

type FileService interface {
	Create(ctx context.Context, file *models.File) (*models.File, error)
	UploadFiles(ctx context.Context, userId uuid.UUID, fileHeaders []*multipart.FileHeader) []dto.UploadResult
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*models.File, error)
	GetByIds(ctx context.Context, fileIds []string) ([]*models.File, error)
	DeleteByIds(ctx context.Context, fileIds []string) error
}
