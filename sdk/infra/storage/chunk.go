package storage

import (
	"context"
	"io"
)

type ChunkInfo struct {
	UploadID    string
	ChunkIndex  int
	TotalChunks int
	Checksum    string
	Size        int64
	Data        io.Reader
}

type ChunkedUploader interface {
	InitChunkUpload(bucket, fileKey string, totalChunks int) (uploadID string, err error)
	UploadChunk(bucket, fileKey, uploadID string, chunk ChunkInfo) error
	CompleteChunkUpload(bucket, fileKey, uploadID string) (filePath string, err error)
	AbortChunkUpload(bucket, fileKey, uploadID string) error
}

type ContextChunkedUploader interface {
	InitChunkUploadWithContext(ctx context.Context, bucket, fileKey string, totalChunks int) (uploadID string, err error)
	UploadChunkWithContext(ctx context.Context, bucket, fileKey, uploadID string, chunk ChunkInfo) error
	CompleteChunkUploadWithContext(ctx context.Context, bucket, fileKey, uploadID string) (filePath string, err error)
	AbortChunkUploadWithContext(ctx context.Context, bucket, fileKey, uploadID string) error
}
