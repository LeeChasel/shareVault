package service

import (
	"context"

	"github.com/LeeChasel/shareVault/internal/models"
	repoInterface "github.com/LeeChasel/shareVault/internal/repository/interfaces"
	serviceInterface "github.com/LeeChasel/shareVault/internal/service/interfaces"
)

type fileService struct {
	fileRepo repoInterface.FileRepository
}

func NewFileService(f repoInterface.FileRepository) serviceInterface.FileService {
	return &fileService{
		fileRepo: f,
	}
}

func (s *fileService) Create(ctx context.Context, file *models.File) (*models.File, error) {
	return s.fileRepo.Create(ctx, file)
}

func (s *fileService) GetByUserId(ctx context.Context, userId string) ([]*models.File, error) {
	return s.fileRepo.GetByUserId(ctx, userId)
}

func (s *fileService) GetByIds(ctx context.Context, fileIds []string) ([]*models.File, error) {
	return s.fileRepo.GetByIds(ctx, fileIds)
}

func (s *fileService) DeleteByIds(ctx context.Context, fileIds []string) error {
	return s.fileRepo.DeleteByIds(ctx, fileIds)
}
