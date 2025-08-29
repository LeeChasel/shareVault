package interfaces

import (
	"context"

	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/google/uuid"
)

type FileRepository interface {
	Create(ctx context.Context, file *models.File) (*models.File, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*models.File, error)
	GetByIds(ctx context.Context, fileIds []string) ([]*models.File, error)
	DeleteByIds(ctx context.Context, fileIds []string) error
}
