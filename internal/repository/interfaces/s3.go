package interfaces

import (
	"context"
	"io"
)

type S3UploadItem struct {
	Key         string
	Body        io.Reader
	ContentType string
}

type S3UploadResult struct {
	Key      string
	Location string
	ETag     string
	Error    error
}

type S3Repository interface {
	UploadFiles(ctx context.Context, items []S3UploadItem) ([]*S3UploadResult, error)
	DeleteFiles(ctx context.Context, keys []string) error
	GetObjectStream(ctx context.Context, key string) (io.ReadCloser, error)
}
