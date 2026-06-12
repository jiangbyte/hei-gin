package storage

import (
	"context"
	"io"
)

type Engine interface {
	Store(bucket, fileKey string, data []byte) (string, error)
	StoreStream(bucket, fileKey string, reader io.Reader) (string, error)
	GetBytes(bucket, fileKey string) ([]byte, error)
	Delete(bucket, fileKey string) error
	Exists(bucket, fileKey string) (bool, error)
	Copy(srcBucket, srcKey, dstBucket, dstKey string) error
}

type SizedStreamer interface {
	StoreStreamWithSize(bucket, fileKey string, reader io.Reader, size int64) (string, error)
}

type ContextSizedStreamer interface {
	StoreStreamWithContext(ctx context.Context, bucket, fileKey string, reader io.Reader, size int64) (string, error)
}

type ContextDeleter interface {
	DeleteWithContext(ctx context.Context, bucket, fileKey string) error
}
