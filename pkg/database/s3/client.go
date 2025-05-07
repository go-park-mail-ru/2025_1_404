package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/go-park-mail-ru/2025_1_404/pkg/logger"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3Repo struct {
	client *minio.Client
	logger logger.Logger
}

func New(cfg *config.MinioConfig, logger logger.Logger) (*s3Repo, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &s3Repo{client: client, logger: logger}, nil
}

func (repo *s3Repo) Put(ctx context.Context, upload Upload) (string, error) {
	objectName := repo.generateFileName(upload.Filename)

	options := minio.PutObjectOptions{
		ContentType: upload.ContentType,
	}
	_, err := repo.client.PutObject(ctx, upload.Bucket, objectName, upload.File, upload.Size, options)
	if err != nil {
		return "", fmt.Errorf("failed to put object: %v", err)
	}

	return objectName, nil
}

func (repo *s3Repo) Get(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	object, err := repo.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	return object, nil
}

func (repo *s3Repo) Remove(ctx context.Context, bucket, objectName string) error {
	err := repo.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete objetc: %v", err)
	}

	return nil
}

func (repo *s3Repo) generateFileName(fileName string) string {
	uuid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uuid, fileName)
}
