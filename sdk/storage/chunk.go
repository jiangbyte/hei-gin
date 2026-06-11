package storage

import "io"

// ChunkInfo describes a single file chunk in a chunked upload session.
type ChunkInfo struct {
	UploadID    string // Upload session identifier
	ChunkIndex  int    // 0-based chunk index
	TotalChunks int    // Total number of chunks in this upload
	Checksum    string // Optional chunk-level checksum (SHA256 hex)
	Size        int64  // Chunk size in bytes when known
	Data        io.Reader
}

// ChunkedUploader is an optional interface for storage backends that support
// chunked/resumable uploads of large files.
type ChunkedUploader interface {
	// InitChunkUpload initializes a chunked upload session.
	// Returns a unique uploadID for subsequent chunk operations.
	InitChunkUpload(bucket, fileKey string, totalChunks int) (uploadID string, err error)

	// UploadChunk stores a single chunk.
	// chunks can be sent in any order; the backend reassembles by ChunkIndex.
	UploadChunk(bucket, fileKey, uploadID string, chunk ChunkInfo) error

	// CompleteChunkUpload finalizes the upload, merges all chunks into the final file,
	// and returns the final file path. Cleans up temporary data on success.
	CompleteChunkUpload(bucket, fileKey, uploadID string) (filePath string, err error)

	// AbortChunkUpload cancels the upload and cleans up all temporary data.
	AbortChunkUpload(bucket, fileKey, uploadID string) error
}
