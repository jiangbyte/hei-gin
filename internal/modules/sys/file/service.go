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
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/shared"
)

// Service æ–‡ä»¶å­˜å‚¨ä¸šåŠ¡æœåŠ¡ã€‚
//
// Author: Charlie
type Service struct {
	repo *Repo
	sto  *storage.Manager
}

// NewService æž„é€ æ–‡ä»¶æœåŠ¡ã€‚
func NewService(db *gorm.DB, sto *storage.Manager) *Service {
	return &Service{repo: NewRepo(db), sto: sto}
}

// New æž„å»º sys.file æ¨¡å—ã€‚
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Storage)
	return module.Module{
		Name:   "sys.file",
		Models: []any{&File{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Upload ä¸Šä¼ æ–‡ä»¶å¹¶å†™å…¥å…ƒæ•°æ®ã€‚
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

// Delete æ‰¹é‡åˆ é™¤æ–‡ä»¶åŠå­˜å‚¨å¯¹è±¡ã€‚
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

// Update æ›´æ–°æ–‡ä»¶å…ƒæ•°æ®ã€‚
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

// Detail è¯¦æƒ…ã€‚
func (s *Service) Detail(ctx context.Context, id string) (*File, error) {
	return s.repo.GetByID(ctx, id)
}

// Page åˆ†é¡µã€‚
func (s *Service) Page(ctx context.Context, q PageParam) (rows []File, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// ListByIDs æ‰¹é‡æŸ¥è¯¢ã€‚
func (s *Service) ListByIDs(ctx context.Context, ids []string) ([]File, error) {
	return s.repo.ListByIDs(ctx, ids)
}

// OpenByID æ‰“å¼€æ–‡ä»¶å†…å®¹æµã€‚
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

// URL èŽ·å–å¯¹è±¡è®¿é—® URLã€‚
func (s *Service) URL(ctx context.Context, objectName string) (*URLResult, error) {
	row, err := s.repo.FindByObjectName(ctx, objectName)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	return &URLResult{ObjectName: objectName, URL: row.URL}, nil
}

// PresignedURL èŽ·å–å¯¹è±¡é¢„ç­¾å URLï¼ˆæœ¬åœ°å­˜å‚¨è¿”å›žå…¬å¼€ URLï¼‰ã€‚
func (s *Service) PresignedURL(ctx context.Context, objectName string) (*URLResult, error) {
	if _, err := s.repo.FindByObjectName(ctx, objectName); err != nil {
		return nil, fmt.Errorf("file not found")
	}
	pub := s.sto.Provider().PublicURL(objectName)
	return &URLResult{ObjectName: objectName, URL: pub}, nil
}

// OpenByObjectName æŒ‰å¯¹è±¡åæ‰“å¼€æ–‡ä»¶ã€‚
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
