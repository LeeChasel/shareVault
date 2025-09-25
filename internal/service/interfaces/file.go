package interfaces

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/LeeChasel/shareVault/internal/api/dto"
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/google/uuid"
)

type FileContent struct {
	Name   string
	Stream io.ReadCloser
}

type FileService interface {
	UploadFiles(ctx context.Context, userId uuid.UUID, fileHeaders []*multipart.FileHeader) []dto.UploadResult
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*models.File, error)
	GetByIds(ctx context.Context, fileIds []string) ([]*models.File, error)
	DeleteFiles(ctx context.Context, files []*models.File) error
	GetFileStreams(ctx context.Context, files []*models.File) ([]*FileContent, error)
}
