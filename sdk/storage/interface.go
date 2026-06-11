package storage

import (
	"context"
	"io"
)

// Engine defines the contract for file storage backends.
// Each engine type (LOCAL, MINIO, S3) implements this interface.
type Engine interface {
	// Store saves data bytes under the given key in the specified bucket.
	// Returns the stored path or key.
	Store(bucket, fileKey string, data []byte) (string, error)

	// StoreStream saves data from a reader under the given key.
	StoreStream(bucket, fileKey string, reader io.Reader) (string, error)

	// GetBytes retrieves the raw bytes for a given key.
	GetBytes(bucket, fileKey string) ([]byte, error)

	// Delete removes the object at the given key.
	Delete(bucket, fileKey string) error

	// Exists checks if an object exists at the given key.
	Exists(bucket, fileKey string) (bool, error)

	// Copy copies an object from source to destination within the same backend.
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
