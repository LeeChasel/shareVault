package service

import (
	"context"

	repoInterface "github.com/LeeChasel/shareVault/internal/repository/interfaces"
	serviceInterface "github.com/LeeChasel/shareVault/internal/service/interfaces"
)

type s3Service struct {
	s3Repo repoInterface.S3Repository
}

func NewS3Service(s3 repoInterface.S3Repository) serviceInterface.S3Service {
	return &s3Service{
		s3Repo: s3,
	}
}

func (s *s3Service) UploadFiles(ctx context.Context, key string, data []byte, contentType string) error {
	return s.s3Repo.UploadFiles(ctx, key, data, contentType)
}

func (s *s3Service) DeleteFiles(ctx context.Context, keys []string) error {
	return s.s3Repo.DeleteFiles(ctx, keys)
}

func (s *s3Service) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	return s.s3Repo.DownloadFile(ctx, key)
}
