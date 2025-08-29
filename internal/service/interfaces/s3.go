package interfaces

import "context"

type S3Service interface {
	DeleteFiles(ctx context.Context, keys []string) error
	DownloadFile(ctx context.Context, key string) ([]byte, error)
}
