package repository

import (
	"bytes"
	"context"

	"github.com/LeeChasel/shareVault/backend/configs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	bucket string
}

func NewS3Repository(bucket string) *S3Repository {
	return &S3Repository{bucket: bucket}
}

func (r *S3Repository) UploadFiles(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := configs.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	return err
}

func (r *S3Repository) DeleteFile(ctx context.Context, key string) error {
	_, err := configs.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	
	return err
}