// internal/modules/sys/file/service.go 业务服务。
//
// Author: Charlie

package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/runtimecfg"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/shared"
)

// Service 文件存储业务服务。
//
// Author: Charlie
type Service struct {
	repo    *Repo
	sto     *storage.Manager
	runtime *runtimecfg.Settings
}

// NewService 构造文件服务。
func NewService(db *gorm.DB, sto *storage.Manager, rt *runtimecfg.Settings) *Service {
	return &Service{repo: NewRepo(db), sto: sto, runtime: rt}
}

// New 构建 sys.file 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Storage, d.Runtime)
	return module.Module{
		Name:   "sys.file",
		Models: []any{&File{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Upload 上传文件并写入元数据（按 STORAGE_UPLOAD_* 运行时配置校验大小/类型/扩展名）。
func (s *Service) Upload(ctx context.Context, fh *multipart.FileHeader) (*File, error) {
	if err := s.validateUpload(ctx, fh); err != nil {
		return nil, err
	}
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	now := time.Now().UTC()
	ext := path.Ext(fh.Filename)
	objectName := fmt.Sprintf("%04d/%02d/%02d/%s%s", now.Year(), now.Month(), now.Day(), strings.ReplaceAll(uuid.NewString(), "-", ""), ext)
	ct := fh.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	url, err := s.sto.Provider().Put(ctx, objectName, f, fh.Size, ct)
	if err != nil {
		return nil, err
	}
	row := File{
		ID: idgen.Next(), ObjectName: objectName, OriginalName: fh.Filename,
		StorageProvider: "local", ContentType: ct, Size: fh.Size, URL: url,
	}
	if err := s.repo.Create(ctx, &row); err != nil {
		return nil, err
	}
	return &row, nil
}

// Delete 批量删除文件及存储对象。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	rows, err := s.repo.ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	for _, r := range rows {
		_ = s.sto.Provider().Delete(ctx, r.ObjectName)
	}
	return s.repo.DeleteByIDs(ctx, ids)
}

// Update 更新文件元数据。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	updates := map[string]any{}
	if req.OriginalName != nil {
		updates["original_name"] = *req.OriginalName
	}
	if len(updates) == 0 {
		return nil
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*File, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []File, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// ListByIDs 批量查询。
func (s *Service) ListByIDs(ctx context.Context, ids []string) ([]File, error) {
	return s.repo.ListByIDs(ctx, ids)
}

// OpenByID 打开文件内容流。
func (s *Service) OpenByID(ctx context.Context, id string) (*File, io.ReadCloser, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	rc, err := s.sto.Provider().Get(ctx, row.ObjectName)
	if err != nil {
		return nil, nil, err
	}
	return row, rc, nil
}

// URL 获取对象访问 URL。
func (s *Service) URL(ctx context.Context, objectName string) (*URLResult, error) {
	row, err := s.repo.FindByObjectName(ctx, objectName)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	return &URLResult{ObjectName: objectName, URL: row.URL}, nil
}

// PresignedURL 获取对象预签名 URL（本地返回公开 URL；S3 按运行时 STORAGE_PRESIGN_EXPIRE_SECONDS 生成）。
func (s *Service) PresignedURL(ctx context.Context, objectName string) (*URLResult, error) {
	if _, err := s.repo.FindByObjectName(ctx, objectName); err != nil {
		return nil, fmt.Errorf("file not found")
	}
	expire := time.Duration(s.sto.PresignExpireSeconds(ctx)) * time.Second
	if p, ok := s.sto.Provider().(storage.Presigner); ok {
		u, err := p.PresignedURL(ctx, objectName, expire)
		if err != nil {
			return nil, err
		}
		return &URLResult{ObjectName: objectName, URL: u}, nil
	}
	pub := s.sto.Provider().PublicURL(objectName)
	return &URLResult{ObjectName: objectName, URL: pub}, nil
}

// validateUpload 按运行时上传限制校验（STORAGE_UPLOAD_*；缺省 20MB、放行全部类型/扩展名）。
func (s *Service) validateUpload(ctx context.Context, fh *multipart.FileHeader) error {
	maxBytes := 20 * 1024 * 1024
	if s.runtime != nil {
		if v := s.runtime.GetInt(ctx, "STORAGE_UPLOAD_MAX_BYTES", 0); v > 0 {
			maxBytes = v
		}
	}
	if fh.Size > int64(maxBytes) {
		return fmt.Errorf("文件大小超过限制（最大 %d 字节）", maxBytes)
	}
	ext := strings.ToLower(path.Ext(fh.Filename))
	ct := strings.ToLower(fh.Header.Get("Content-Type"))
	allowedTypes := strings.Split(s.str(ctx, "STORAGE_UPLOAD_ALLOWED_CONTENT_TYPES"), ",")
	allowedExts := strings.Split(s.str(ctx, "STORAGE_UPLOAD_ALLOWED_EXTENSIONS"), ",")
	deniedExts := strings.Split(s.str(ctx, "STORAGE_UPLOAD_DENIED_EXTENSIONS"), ",")
	for _, d := range deniedExts {
		if d = strings.TrimSpace(strings.ToLower(d)); d != "" && (d == ext || d == "."+strings.TrimPrefix(ext, ".")) {
			return fmt.Errorf("不允许上传该文件类型：%s", ext)
		}
	}
	if len(allowedTypes) > 0 && allowedTypes[0] != "" && !containsFold(allowedTypes, ct) {
		return fmt.Errorf("不允许上传该内容类型：%s", ct)
	}
	if len(allowedExts) > 0 && allowedExts[0] != "" && !containsFold(allowedExts, ext) {
		return fmt.Errorf("不允许上传该扩展名：%s", ext)
	}
	return nil
}

func (s *Service) str(ctx context.Context, key string) string {
	if s.runtime == nil {
		return ""
	}
	return s.runtime.GetString(ctx, key, "")
}

func containsFold(list []string, v string) bool {
	for _, it := range list {
		if strings.EqualFold(strings.TrimSpace(it), v) {
			return true
		}
	}
	return false
}

// OpenByObjectName 按对象名打开文件。
func (s *Service) OpenByObjectName(ctx context.Context, objectName string) (contentType string, rc io.ReadCloser, err error) {
	rc, err = s.sto.Provider().Get(ctx, objectName)
	if err != nil {
		return "", nil, err
	}
	ct := "application/octet-stream"
	if row, err := s.repo.FindByObjectName(ctx, objectName); err == nil {
		ct = row.ContentType
	}
	return ct, rc, nil
}
