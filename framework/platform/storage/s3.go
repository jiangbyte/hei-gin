package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"hei-gin/framework/core/config"
)

// S3 兼容对象存储（AWS S3 / MinIO / OSS）。
//
// Author: Charlie
type S3 struct {
	client  *s3.Client
	bucket  string
	baseURL string
}

// NewS3 创建 S3 兼容 Provider；endpoint 非空时使用 path-style（MinIO）。
func NewS3(cfg config.StorageConfig) (*S3, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("storage: bucket is required for provider %q", cfg.Provider)
	}
	region := cfg.Region
	if region == "" {
		region = "us-east-1"
	}
	loadOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}
	if cfg.AccessKey != "" || cfg.SecretKey != "" {
		loadOpts = append(loadOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}

	clientOpts := []func(*s3.Options){}
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint != "" {
		ep := endpoint
		if !strings.HasPrefix(ep, "http://") && !strings.HasPrefix(ep, "https://") {
			if cfg.UseSSL {
				ep = "https://" + ep
			} else {
				ep = "http://" + ep
			}
		}
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(ep)
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" && endpoint != "" {
		ep := endpoint
		if !strings.HasPrefix(ep, "http://") && !strings.HasPrefix(ep, "https://") {
			if cfg.UseSSL {
				ep = "https://" + ep
			} else {
				ep = "http://" + ep
			}
		}
		base = strings.TrimRight(ep, "/") + "/" + cfg.Bucket
	}
	return &S3{client: client, bucket: cfg.Bucket, baseURL: base}, nil
}

// Put 上传对象并返回公开 URL。
func (s *S3) Put(ctx context.Context, objectName string, r io.Reader, size int64, contentType string) (string, error) {
	key := strings.TrimLeft(objectName, "/")
	in := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   r,
	}
	if size >= 0 {
		in.ContentLength = aws.Int64(size)
	}
	if contentType != "" {
		in.ContentType = aws.String(contentType)
	}
	if _, err := s.client.PutObject(ctx, in); err != nil {
		return "", err
	}
	return s.PublicURL(key), nil
}

// Get 下载对象流。
func (s *S3) Get(ctx context.Context, objectName string) (io.ReadCloser, error) {
	key := strings.TrimLeft(objectName, "/")
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

// Delete 删除对象。
func (s *S3) Delete(ctx context.Context, objectName string) error {
	key := strings.TrimLeft(objectName, "/")
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// PublicURL 拼出对象公开访问 URL。
func (s *S3) PublicURL(objectName string) string {
	key := strings.TrimLeft(objectName, "/")
	if s.baseURL != "" {
		return s.baseURL + "/" + key
	}
	u := url.URL{Scheme: "https", Host: s.bucket + ".s3.amazonaws.com", Path: "/" + key}
	return u.String()
}
