package s3

import (
	"context"
	"io"
)

//go:generate mockgen -source interface.go -destination=mocks/mock_s3.go -package=mocks

type S3Repo interface {
	Get(ctx context.Context, bucket string, objectName string) (io.ReadCloser, error)
	Put(ctx context.Context, upload Upload) (string, error)
	Remove(ctx context.Context, bucket string, objectName string) error
}
