package repository

import (
	"context"
	"io"

	"github.com/LeeChasel/shareVault/internal/constants"
	"github.com/LeeChasel/shareVault/internal/repository/interfaces"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/errgroup"
)

type s3Repository struct {
	uploader   *manager.Uploader
	downloader *manager.Downloader
	client     *s3.Client
	bucketName string
}

func NewS3Repository(cfg aws.Config, bucketName string) interfaces.S3Repository {
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.Concurrency = 10
		u.PartSize = 10 * constants.MB
	})
	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.Concurrency = 5
		d.PartSize = 10 * constants.MB
	})

	return &s3Repository{
		uploader:   uploader,
		downloader: downloader,
		client:     client,
		bucketName: bucketName,
	}
}

func (r *s3Repository) UploadFiles(ctx context.Context, items []interfaces.S3UploadItem) ([]*interfaces.S3UploadResult, error) {
	results := make([]*interfaces.S3UploadResult, len(items))

	if len(items) == 0 {
		return results, nil
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(3)

	for i, item := range items {
		index, uploadItem := i, item

		g.Go(func() error {
			result := r.uploadSingleFile(ctx, uploadItem)
			results[index] = result

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return results, err
	}

	return results, nil
}

func (r *s3Repository) uploadSingleFile(ctx context.Context, item interfaces.S3UploadItem) *interfaces.S3UploadResult {
	result := &interfaces.S3UploadResult{
		Key: item.Key,
	}

	uploadResult, err := r.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(item.Key),
		Body:        item.Body,
		ContentType: aws.String(item.ContentType),
	})

	if err != nil {
		result.Error = err
		return result
	}

	result.Location = uploadResult.Location
	result.ETag = aws.ToString(uploadResult.ETag)

	return result
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

func (r *s3Repository) GetObjectStream(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}
