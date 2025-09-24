package interfaces

import "context"

type S3Service interface {
	DownloadFile(ctx context.Context, key string) ([]byte, error)
}
