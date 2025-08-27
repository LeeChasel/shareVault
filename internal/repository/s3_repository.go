package repository

import (
	"bytes"
	"context"
	"io"

	"github.com/LeeChasel/shareVault/internal/repository/interfaces"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Repository struct {
	client     *s3.Client
	bucketName string
}

func NewS3Repository(client *s3.Client, bucketName string) interfaces.S3Repository {
	return &s3Repository{
		client:     client,
		bucketName: bucketName,
	}
}

func (r *s3Repository) UploadFiles(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	return err
}

func (r *s3Repository) DeleteFiles(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(key),
		}
	}
	_, err := r.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(r.bucketName),
		Delete: &types.Delete{
			Objects: objects,
		},
	})

	return err
}

func (r *s3Repository) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
