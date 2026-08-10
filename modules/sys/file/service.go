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

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/framework/platform/storage"
	"hei-gin/modules/shared"
)

// Service 文件存储业务服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
	sto  *storage.Manager
}

// NewService 构造文件服务。
func NewService(db *gorm.DB, sto *storage.Manager) *Service {
	return &Service{repo: NewRepo(db), sto: sto}
}

// New 构建 sys.file 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Storage)
	return module.Module{
		Name:   "sys.file",
		Models: []any{&File{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Upload 上传文件并写入元数据。
func (s *Service) Upload(ctx context.Context, fh *multipart.FileHeader) (*File, error) {
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
