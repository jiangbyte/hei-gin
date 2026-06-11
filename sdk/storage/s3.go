package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements FileStorage using AWS S3 or any S3-compatible object store.
type S3 struct {
	client        *s3.Client
	defaultBucket string
	baseURL       string
	endpoint      string
}

// NewS3Storage creates a new S3-compatible storage backend.
func NewS3(endpoint, accessKey, secretKey, defaultBucket, region string, pathStyle bool, baseURL string) *S3 {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Printf("[storage/s3] failed to load AWS config: %v", err)
		return nil
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = pathStyle
	})

	return &S3{
		client:        client,
		defaultBucket: defaultBucket,
		baseURL:       baseURL,
		endpoint:      endpoint,
	}
}

// GetDefaultBucket returns the default bucket name.
func (s *S3) GetDefaultBucket() string {
	return s.defaultBucket
}

// _ensureBucket creates the bucket if it does not already exist.
func (s *S3) _ensureBucket(ctx context.Context, bucket string) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		return err
	}
	return nil
}

// Store uploads raw bytes to S3 and returns the object key.
func (s *S3) Store(bucket, fileKey string, data []byte) (string, error) {
	ctx := context.Background()
	if err := s._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

// StoreStream reads all data from the reader and uploads it to S3.
func (s *S3) StoreStream(bucket, fileKey string, reader io.Reader) (string, error) {
	return s.StoreStreamWithSize(bucket, fileKey, reader, -1)
}

func (s *S3) StoreStreamWithSize(bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	ctx := context.Background()
	if err := s._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   reader,
	}
	if size >= 0 {
		input.ContentLength = aws.Int64(size)
	}
	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

// GetBytes downloads and returns the raw bytes of an object from S3.
func (s *S3) GetBytes(bucket, fileKey string) ([]byte, error) {
	ctx := context.Background()
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()
	return io.ReadAll(output.Body)
}

// GetURL returns the endpoint-based URL for the object.
func (s *S3) GetURL(bucket, fileKey string) string {
	if s.baseURL != "" {
		return s.baseURL + bucket + "/" + fileKey
	}
	return s.endpoint + "/" + bucket + "/" + fileKey
}

// GetAuthURL returns a time-limited presigned URL for the object.
func (s *S3) GetAuthURL(bucket, fileKey string, timeoutMs int) (string, error) {
	ctx := context.Background()
	presignClient := s3.NewPresignClient(s.client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(timeoutMs) * time.Millisecond
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// Delete removes an object from S3.
func (s *S3) Delete(bucket, fileKey string) error {
	ctx := context.Background()
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	return err
}

// Exists checks whether an object exists in S3.
func (s *S3) Exists(bucket, fileKey string) (bool, error) {
	ctx := context.Background()
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Copy copies an object from source to destination within S3.
func (s *S3) Copy(srcBucket, srcKey, dstBucket, dstKey string) error {
	ctx := context.Background()
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(dstBucket),
		CopySource: aws.String(srcBucket + "/" + srcKey),
		Key:        aws.String(dstKey),
	})
	return err
}

// ===== ChunkedUploader implementation (S3 native multipart upload) =====

// InitChunkUpload initializes an S3 multipart upload session.
func (s *S3) InitChunkUpload(bucket, fileKey string, totalChunks int) (string, error) {
	ctx := context.Background()
	if err := s._ensureBucket(ctx, bucket); err != nil {
		return "", err
	}
	output, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", err
	}
	return *output.UploadId, nil
}

// UploadChunk uploads a single chunk as a part in the S3 multipart upload.
func (s *S3) UploadChunk(bucket, fileKey, uploadID string, chunk ChunkInfo) error {
	ctx := context.Background()
	partNumber := int32(chunk.ChunkIndex + 1) // S3 parts are 1-based
	_, err := s.client.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(fileKey),
		UploadId:      aws.String(uploadID),
		PartNumber:    aws.Int32(partNumber),
		Body:          chunk.Data,
		ContentLength: aws.Int64(chunk.Size),
	})
	return err
}

// CompleteChunkUpload lists all uploaded parts and completes the S3 multipart upload.
func (s *S3) CompleteChunkUpload(bucket, fileKey, uploadID string) (string, error) {
	ctx := context.Background()

	var completedParts []types.CompletedPart
	var partNumberMarker *string

	for {
		output, err := s.client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:           aws.String(bucket),
			Key:              aws.String(fileKey),
			UploadId:         aws.String(uploadID),
			PartNumberMarker: partNumberMarker,
		})
		if err != nil {
			return "", err
		}
		for _, p := range output.Parts {
			completedParts = append(completedParts, types.CompletedPart{
				PartNumber: p.PartNumber,
				ETag:       p.ETag,
			})
		}
		if output.IsTruncated == nil || !*output.IsTruncated {
			break
		}
		partNumberMarker = output.NextPartNumberMarker
	}

	if len(completedParts) == 0 {
		return "", errors.New("no parts to complete")
	}

	_, err := s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(fileKey),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return "", err
	}
	return bucket + "/" + fileKey, nil
}

// AbortChunkUpload aborts the S3 multipart upload and cleans up partial data.
func (s *S3) AbortChunkUpload(bucket, fileKey, uploadID string) error {
	ctx := context.Background()
	_, err := s.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(fileKey),
		UploadId: aws.String(uploadID),
	})
	return err
}
