package interfaces

import "context"

type S3Repository interface {
	UploadFiles(ctx context.Context, key string, data []byte, contentType string) error
	DeleteFiles(ctx context.Context, keys []string) error
	DownloadFile(ctx context.Context, key string) ([]byte, error)
}
