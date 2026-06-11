package storage

import (
	"bytes"
	"context"
	"strings"
	"errors"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage implements FileStorage using MinIO (S3-compatible object storage).
// It uses minio.Core for both standard operations (via embedded *Client) and
// low-level multipart upload APIs.
type Minio struct {
	core          *minio.Core
	defaultBucket string
	baseURL      string
	endpoint      string
}
func NewMinio(endpoint, accessKey, secretKey, defaultBucket string, secure bool, region string, baseURL string) *Minio {
	core, err := minio.NewCore(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region,
	})
	if err != nil {
		log.Printf("[storage/minio] failed to create MinIO client: %v", err)
		return nil
	}
	return &Minio{
		core:          core,
		defaultBucket: defaultBucket,
		baseURL:       baseURL,
		endpoint:      endpoint,
}
	}

// client returns the embedded *minio.Client for standard operations.
func (m *Minio) client() *minio.Client {
	return m.core.Client
}

// GetDefaultBucket returns the default bucket name.
func (m *Minio) GetDefaultBucket() string {
	return m.defaultBucket
}

// _ensureBucket creates the bucket if it does not already exist.
func (m *Minio) _ensureBucket(ctx context.Context, bucket string) error {
	exists, err := m.client().BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return m.client().MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}

// Store uploads raw bytes to MinIO and returns the object key.
func (m *Minio) Store(bucket, fileKey string, data []byte) (string, error) {
	ctx := context.Background()
	if err := m._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	_, err := m.client().PutObject(ctx, bucket, fileKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

// StoreStream uploads data from a reader to MinIO.
func (m *Minio) StoreStream(bucket, fileKey string, reader io.Reader) (string, error) {
	ctx := context.Background()
	if err := m._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	_, err = m.client().PutObject(ctx, bucket, fileKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

// GetBytes downloads and returns the raw bytes of an object from MinIO.
func (m *Minio) GetBytes(bucket, fileKey string) ([]byte, error) {
	ctx := context.Background()
	obj, err := m.client().GetObject(ctx, bucket, fileKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

// GetURL returns the public endpoint-based URL for the object.
func (m *Minio) GetURL(bucket, fileKey string) string {
	if m.baseURL != "" {
		return m.baseURL + bucket + "/" + fileKey
	}
	scheme := "http"
	endpoint := m.endpoint
	if strings.HasPrefix(endpoint, "https://") {
		scheme = "https"
		endpoint = strings.TrimPrefix(endpoint, "https://")
	} else if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	}
	return scheme + "://" + endpoint + "/" + bucket + "/" + fileKey
}

// GetAuthURL returns a time-limited presigned URL for the object.
func (m *Minio) GetAuthURL(bucket, fileKey string, timeoutMs int) (string, error) {
	ctx := context.Background()
	expiry := time.Duration(timeoutMs) * time.Millisecond
	u, err := m.client().PresignedGetObject(ctx, bucket, fileKey, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Delete removes an object from MinIO.
func (m *Minio) Delete(bucket, fileKey string) error {
	ctx := context.Background()
	err := m.client().RemoveObject(ctx, bucket, fileKey, minio.RemoveObjectOptions{})
	if err != nil {
		var errResp minio.ErrorResponse
		if errors.As(err, &errResp) && errResp.Code == "NoSuchKey" {
			return nil
		}
		return err
	}
	return nil
}

// Exists checks whether an object exists in MinIO.
func (m *Minio) Exists(bucket, fileKey string) (bool, error) {
	ctx := context.Background()
	_, err := m.client().StatObject(ctx, bucket, fileKey, minio.StatObjectOptions{})
	if err != nil {
		var errResp minio.ErrorResponse
		if errors.As(err, &errResp) && errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Copy copies an object from source to destination within MinIO.
func (m *Minio) Copy(srcBucket, srcKey, dstBucket, dstKey string) error {
	ctx := context.Background()
	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcKey,
	}
	dst := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: dstKey,
	}
	_, err := m.client().CopyObject(ctx, dst, src)
	return err
}

// ===== ChunkedUploader implementation (MinIO native multipart upload) =====

// InitChunkUpload initializes a MinIO multipart upload session.
func (m *Minio) InitChunkUpload(bucket, fileKey string, totalChunks int) (string, error) {
	ctx := context.Background()
	if err := m._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	uploadID, err := m.core.NewMultipartUpload(ctx, bucket, fileKey, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return uploadID, nil
}

// UploadChunk uploads a single chunk as a part in the MinIO multipart upload.
func (m *Minio) UploadChunk(bucket, fileKey, uploadID string, chunk ChunkInfo) error {
	ctx := context.Background()
	data, err := io.ReadAll(chunk.Data)
	if err != nil {
		return err
	}
	partNumber := chunk.ChunkIndex + 1 // MinIO parts are 1-based
	_, err = m.core.PutObjectPart(ctx, bucket, fileKey, uploadID, partNumber,
		bytes.NewReader(data), int64(len(data)), minio.PutObjectPartOptions{})
	return err
}

// CompleteChunkUpload lists all uploaded parts and completes the multipart upload.
func (m *Minio) CompleteChunkUpload(bucket, fileKey, uploadID string) (string, error) {
	ctx := context.Background()

	var parts []minio.CompletePart
	partNumberMarker := 0
	for {
		result, err := m.core.ListObjectParts(ctx, bucket, fileKey, uploadID, partNumberMarker, 1000)
		if err != nil {
			return "", err
		}
		for _, p := range result.ObjectParts {
			parts = append(parts, minio.CompletePart{
				PartNumber: p.PartNumber,
				ETag:       p.ETag,
			})
		}
		if !result.IsTruncated {
			break
		}
		partNumberMarker = result.NextPartNumberMarker
	}

	if len(parts) == 0 {
		return "", errors.New("no parts to complete")
	}

	_, err := m.core.CompleteMultipartUpload(ctx, bucket, fileKey, uploadID, parts, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

// AbortChunkUpload aborts the MinIO multipart upload and cleans up partial data.
func (m *Minio) AbortChunkUpload(bucket, fileKey, uploadID string) error {
	ctx := context.Background()
	return m.core.AbortMultipartUpload(ctx, bucket, fileKey, uploadID)
}
