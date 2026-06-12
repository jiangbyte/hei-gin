package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	core          *minio.Core
	defaultBucket string
	baseURL       string
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

func (m *Minio) client() *minio.Client {
	return m.core.Client
}

func (m *Minio) GetDefaultBucket() string {
	return m.defaultBucket
}

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

func (m *Minio) StoreStream(bucket, fileKey string, reader io.Reader) (string, error) {
	return m.StoreStreamWithSize(bucket, fileKey, reader, -1)
}

func (m *Minio) StoreStreamWithSize(bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	return m.StoreStreamWithContext(context.Background(), bucket, fileKey, reader, size)
}

func (m *Minio) StoreStreamWithContext(ctx context.Context, bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	if err := m._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	_, err := m.client().PutObject(ctx, bucket, fileKey, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

func (m *Minio) GetBytes(bucket, fileKey string) ([]byte, error) {
	ctx := context.Background()
	obj, err := m.client().GetObject(ctx, bucket, fileKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

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

func (m *Minio) GetAuthURL(bucket, fileKey string, timeoutMs int) (string, error) {
	ctx := context.Background()
	expiry := time.Duration(timeoutMs) * time.Millisecond
	u, err := m.client().PresignedGetObject(ctx, bucket, fileKey, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m *Minio) Delete(bucket, fileKey string) error {
	return m.DeleteWithContext(context.Background(), bucket, fileKey)
}

func (m *Minio) DeleteWithContext(ctx context.Context, bucket, fileKey string) error {
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

func (m *Minio) Copy(srcBucket, srcKey, dstBucket, dstKey string) error {
	ctx := context.Background()
	src := minio.CopySrcOptions{Bucket: srcBucket, Object: srcKey}
	dst := minio.CopyDestOptions{Bucket: dstBucket, Object: dstKey}
	_, err := m.client().CopyObject(ctx, dst, src)
	return err
}

func (m *Minio) InitChunkUpload(bucket, fileKey string, totalChunks int) (string, error) {
	return m.InitChunkUploadWithContext(context.Background(), bucket, fileKey, totalChunks)
}

func (m *Minio) InitChunkUploadWithContext(ctx context.Context, bucket, fileKey string, totalChunks int) (string, error) {
	if err := m._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	uploadID, err := m.core.NewMultipartUpload(ctx, bucket, fileKey, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return uploadID, nil
}

func (m *Minio) UploadChunk(bucket, fileKey, uploadID string, chunk ChunkInfo) error {
	return m.UploadChunkWithContext(context.Background(), bucket, fileKey, uploadID, chunk)
}

func (m *Minio) UploadChunkWithContext(ctx context.Context, bucket, fileKey, uploadID string, chunk ChunkInfo) error {
	partNumber := chunk.ChunkIndex + 1
	_, err := m.core.PutObjectPart(ctx, bucket, fileKey, uploadID, partNumber, chunk.Data, chunk.Size, minio.PutObjectPartOptions{})
	return err
}

func (m *Minio) CompleteChunkUpload(bucket, fileKey, uploadID string) (string, error) {
	return m.CompleteChunkUploadWithContext(context.Background(), bucket, fileKey, uploadID)
}

func (m *Minio) CompleteChunkUploadWithContext(ctx context.Context, bucket, fileKey, uploadID string) (string, error) {
	var parts []minio.CompletePart
	partNumberMarker := 0
	for {
		result, err := m.core.ListObjectParts(ctx, bucket, fileKey, uploadID, partNumberMarker, 1000)
		if err != nil {
			return "", err
		}
		for _, p := range result.ObjectParts {
			parts = append(parts, minio.CompletePart{PartNumber: p.PartNumber, ETag: p.ETag})
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

func (m *Minio) AbortChunkUpload(bucket, fileKey, uploadID string) error {
	return m.AbortChunkUploadWithContext(context.Background(), bucket, fileKey, uploadID)
}

func (m *Minio) AbortChunkUploadWithContext(ctx context.Context, bucket, fileKey, uploadID string) error {
	return m.core.AbortMultipartUpload(ctx, bucket, fileKey, uploadID)
}
