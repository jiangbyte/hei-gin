package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Local implements Engine on the local filesystem.
type Local struct {
	uploadFolder string
	baseURL      string
}

// NewLocal creates a new local filesystem storage backend.
func NewLocal(uploadFolder, baseURL string) *Local {
	return &Local{uploadFolder: uploadFolder, baseURL: baseURL}
}

func (s *Local) safePath(bucket, fileKey string) (string, error) {
	if strings.Contains(bucket, "..") || strings.Contains(fileKey, "..") ||
		strings.Contains(bucket, "/") || strings.Contains(bucket, "\\") ||
		strings.Contains(fileKey, "/") || strings.Contains(fileKey, "\\") {
		return "", fmt.Errorf("invalid path: contains directory traversal characters")
	}
	return filepath.Join(s.uploadFolder, bucket, fileKey), nil
}

func (s *Local) ensurePath(bucket, fileKey string) (string, error) {
	fullPath, err := s.safePath(bucket, fileKey)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return fullPath, nil
}

func (s *Local) getPath(bucket, fileKey string) (string, error) {
	return s.safePath(bucket, fileKey)
}

// Store writes raw bytes to the local filesystem and returns the full path.
func (s *Local) Store(bucket, fileKey string, data []byte) (string, error) {
	path, err := s.ensurePath(bucket, fileKey)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

// StoreStream copies data from a reader to the local filesystem.
func (s *Local) StoreStream(bucket, fileKey string, reader io.Reader) (string, error) {
	return s.StoreStreamWithSize(bucket, fileKey, reader, -1)
}

func (s *Local) StoreStreamWithSize(bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	path, err := s.ensurePath(bucket, fileKey)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}
	return path, nil
}

// GetBytes reads and returns the raw bytes of a stored file.
func (s *Local) GetBytes(bucket, fileKey string) ([]byte, error) {
	path, err := s.getPath(bucket, fileKey)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// GetURL returns the HTTP access URL for the file.
// If baseURL is configured (e.g. "http://localhost:18886/"), returns baseURL + bucket + "/" + fileKey.
// Otherwise returns a relative path "/uploads/" + bucket + "/" + fileKey for the static file server.
func (s *Local) GetURL(bucket, fileKey string) string {
	if s.baseURL != "" {
		return s.baseURL + bucket + "/" + fileKey
	}
	return "/uploads/" + bucket + "/" + fileKey
}

// Delete removes a file from the local filesystem.
func (s *Local) Delete(bucket, fileKey string) error {
	path, err := s.getPath(bucket, fileKey)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Exists checks whether a file exists on the local filesystem.
func (s *Local) Exists(bucket, fileKey string) (bool, error) {
	path, err := s.getPath(bucket, fileKey)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// Copy copies a file from source to destination on the local filesystem.
func (s *Local) Copy(srcBucket, srcKey, dstBucket, dstKey string) error {
	src, err := s.getPath(srcBucket, srcKey)
	if err != nil {
		return err
	}
	dst, err := s.ensurePath(dstBucket, dstKey)
	if err != nil {
		return err
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
