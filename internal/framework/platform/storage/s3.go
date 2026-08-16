// internal/framework/platform/storage/s3.go S3 兼容对象存储。
//
// Author: Charlie

package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 兼容对象存储（AWS S3 / MinIO / OSS / RustFS / COS）。
//
// Author: Charlie
type S3 struct {
	client       *s3.Client
	bucket       string
	baseURL      string
	fallback     string
	presign      time.Duration
	bucketPublic bool
	pathStyle    bool
	endpointURL  string
}

// NewS3 创建 S3 兼容 Provider；minio/rustfs 使用 path-style。
func NewS3(cfg ResolvedConfig) (*S3, error) {
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

	provider := ResolveProvider(cfg.Provider)
	pathStyle := provider == ProviderMinIO || provider == ProviderRustFS

	clientOpts := []func(*s3.Options){}
	endpoint := strings.TrimSpace(cfg.Endpoint)
	endpointURL := ""
	if endpoint != "" {
		ep := endpoint
		if !strings.HasPrefix(ep, "http://") && !strings.HasPrefix(ep, "https://") {
			if cfg.UseSSL {
				ep = "https://" + ep
			} else {
				ep = "http://" + ep
			}
		}
		endpointURL = strings.TrimRight(ep, "/")
		usePath := pathStyle
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpointURL)
			o.UsePathStyle = usePath
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)
	presign := time.Duration(cfg.PresignExpireSeconds) * time.Second
	if presign <= 0 {
		presign = 3600 * time.Second
	}
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	fallback := ""
	if endpointURL != "" {
		fallback = endpointURL + "/" + cfg.Bucket
	}
	if fallback == "" {
		fallback = (&url.URL{Scheme: "https", Host: cfg.Bucket + ".s3.amazonaws.com", Path: "/"}).String()
		fallback = strings.TrimRight(fallback, "/")
	}
	return &S3{
		client:       client,
		bucket:       cfg.Bucket,
		baseURL:      baseURL,
		fallback:     fallback,
		presign:      presign,
		bucketPublic: cfg.BucketPublic,
		pathStyle:    pathStyle,
		endpointURL:  endpointURL,
	}, nil
}

// Put 上传对象并返回访问 URL。
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
	return s.PublicURL(ctx, key), nil
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

// BucketName 返回对象存储桶名（元数据落库用）。
func (s *S3) BucketName() string { return s.bucket }

// PublicURL 对齐 hei-boot / hei-fastapi：
// 公开桶 → 永久直链（BASE_URL 或 endpoint+bucket）；私有桶 → 预签名 GET。
func (s *S3) PublicURL(ctx context.Context, objectName string) string {
	key := strings.TrimLeft(objectName, "/")
	if s.bucketPublic {
		return s.publicDirectURL(key)
	}
	u, err := s.PresignedURL(ctx, key, s.presign)
	if err == nil && u != "" {
		return u
	}
	return s.fallback + "/" + key
}

func (s *S3) publicDirectURL(key string) string {
	if s.baseURL != "" {
		return s.baseURL + "/" + key
	}
	if s.endpointURL != "" {
		if s.pathStyle {
			return s.endpointURL + "/" + s.bucket + "/" + key
		}
		if u, err := url.Parse(s.endpointURL); err == nil && u.Host != "" {
			scheme := u.Scheme
			if scheme == "" {
				scheme = "https"
			}
			return scheme + "://" + s.bucket + "." + u.Host + "/" + key
		}
		return s.endpointURL + "/" + s.bucket + "/" + key
	}
	return s.fallback + "/" + key
}

// PresignedURL 生成 S3 预签名 GET URL。
func (s *S3) PresignedURL(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	key := strings.TrimLeft(objectName, "/")
	if expire <= 0 {
		expire = s.presign
	}
	pc := s3.NewPresignClient(s.client)
	req, err := pc.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(o *s3.PresignOptions) {
		o.Expires = expire
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}
