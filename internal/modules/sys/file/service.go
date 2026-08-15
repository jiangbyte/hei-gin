// internal/modules/sys/file/service.go 业务服务（对齐 hei-boot FileServiceImpl + FileAccessUrls）。
//
// Author: Charlie

package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
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

// Upload 上传文件并写入元数据（按 STORAGE_UPLOAD_* 运行时配置校验；storageProvider 记录到元数据）。
func (s *Service) Upload(ctx context.Context, fh *multipart.FileHeader, storageProvider, accountID string) (*File, error) {
	if err := s.validateUpload(ctx, fh); err != nil {
		return nil, err
	}
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	originalName := safeOriginalName(fh.Filename)
	objectName := s.buildObjectName(ctx, originalName, "uploads")
	ct := sanitizeContentType(fh.Header.Get("Content-Type"))
	// 按传入 storage_provider 解析引擎（缺省活动引擎；对齐 hei-boot storageSettingsResolver.resolve）。
	prov := s.sto.Provider()
	provider := s.sto.ProviderName()
	if strings.TrimSpace(storageProvider) != "" {
		p := s.sto.ProviderByName(ctx, storageProvider)
		if p != nil {
			prov = p
			provider = strings.ToLower(strings.TrimSpace(storageProvider))
		}
	}
	urlVal, err := prov.Put(ctx, objectName, f, fh.Size, ct)
	if err != nil {
		return nil, err
	}
	bucket := s.sto.Bucket()
	row := File{
		ID: idgen.Next(), ObjectName: objectName, OriginalName: originalName,
		StorageProvider: provider, ContentType: ct, Size: fh.Size, URL: urlVal,
	}
	if bucket != "" {
		row.Bucket = &bucket
	}
	if accountID != "" {
		row.CreatedBy = &accountID
		row.UpdatedBy = &accountID
	}
	if err := s.repo.Create(ctx, &row); err != nil {
		return nil, err
	}
	return s.withResolvedURL(ctx, &row), nil
}

// Delete 批量删除文件及存储对象（object_name 兼容 URL/路径形式；存储删除失败不阻断元数据清理）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	rows, err := s.repo.ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	for i := range rows {
		objectKey := toObjectKey(rows[i].ObjectName, s.publicPath())
		if objectKey == "" {
			continue
		}
		_ = s.providerFor(ctx, &rows[i]).Delete(ctx, objectKey)
	}
	return s.repo.DeleteByIDs(ctx, ids)
}

// DeleteByObjectName 按对象名删除（跨模块清理附件用；外部 URL 忽略）。
func (s *Service) DeleteByObjectName(ctx context.Context, objectName string) error {
	name := strings.TrimSpace(objectName)
	if name == "" || isExternalURL(name) {
		return nil
	}
	row, err := s.repo.FindByObjectName(ctx, toObjectKey(name, s.publicPath()))
	if err != nil {
		return nil
	}
	if key := toObjectKey(row.ObjectName, s.publicPath()); key != "" {
		_ = s.providerFor(ctx, row).Delete(ctx, key)
	}
	return s.repo.DeleteByIDs(ctx, []string{row.ID})
}

// Update 更新文件元数据（original_name 必填；对齐 hei-boot FileServiceImpl.update）。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	if _, err := s.repo.GetByID(ctx, req.ID); err != nil {
		return fmt.Errorf("file not found")
	}
	return s.repo.Update(ctx, req.ID, map[string]any{
		"original_name": safeOriginalName(req.OriginalName),
	})
}

// Detail 详情（重算访问 URL + 回填创建/更新人名）。
func (s *Service) Detail(ctx context.Context, id string) (*File, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.fillNames(ctx, []File{*row})
	return s.withResolvedURL(ctx, row), nil
}

// Page 分页（重算访问 URL + 回填人名；支持 original_name/object_name/storage_provider/content_type 过滤）。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []File, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	if len(rows) > 0 {
		s.fillNames(ctx, rows)
	}
	for i := range rows {
		s.withResolvedURL(ctx, &rows[i])
	}
	return rows, total, current, size, err
}

// ListByIDs 批量查询（重算访问 URL + 回填人名）。
func (s *Service) ListByIDs(ctx context.Context, ids []string) ([]File, error) {
	rows, err := s.repo.ListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(rows) > 0 {
		s.fillNames(ctx, rows)
	}
	for i := range rows {
		s.withResolvedURL(ctx, &rows[i])
	}
	return rows, nil
}

// OpenByID 打开文件内容流。
func (s *Service) OpenByID(ctx context.Context, id string) (*File, io.ReadCloser, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	rc, err := s.providerFor(ctx, row).Get(ctx, toObjectKey(row.ObjectName, s.publicPath()))
	if err != nil {
		return nil, nil, err
	}
	return row, rc, nil
}

// URL 获取对象访问 URL（始终重算，避免返回已过期的预签名；对齐 hei-boot FileServiceImpl.url）。
func (s *Service) URL(ctx context.Context, objectName string) (*URLResult, error) {
	key := toObjectKey(objectName, s.publicPath())
	if key == "" {
		return nil, fmt.Errorf("file not found")
	}
	row, err := s.repo.FindByObjectName(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	return &URLResult{ObjectName: key, URL: s.providerFor(ctx, row).PublicURL(key)}, nil
}

// PresignedURL 获取对象预签名 URL（本地返回公开 URL；S3 按运行时 STORAGE_PRESIGN_EXPIRE_SECONDS 生成）。
func (s *Service) PresignedURL(ctx context.Context, objectName string) (*URLResult, error) {
	key := toObjectKey(objectName, s.publicPath())
	if key == "" {
		return nil, fmt.Errorf("file not found")
	}
	row, err := s.repo.FindByObjectName(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	p := s.providerFor(ctx, row)
	expire := time.Duration(s.sto.PresignExpireSeconds(ctx)) * time.Second
	if pr, ok := p.(storage.Presigner); ok {
		u, err := pr.PresignedURL(ctx, key, expire)
		if err != nil {
			return nil, err
		}
		return &URLResult{ObjectName: key, URL: u}, nil
	}
	return &URLResult{ObjectName: key, URL: p.PublicURL(key)}, nil
}

// OpenByObjectName 按对象名打开文件（公开访问：先校验元数据存在，避免越权读取存储对象）。
func (s *Service) OpenByObjectName(ctx context.Context, objectName string) (contentType string, rc io.ReadCloser, err error) {
	if !validObjectName(objectName) {
		return "", nil, fmt.Errorf("file not found")
	}
	key := toObjectKey(objectName, s.publicPath())
	if key == "" {
		return "", nil, fmt.Errorf("file not found")
	}
	row, err := s.repo.FindByObjectName(ctx, key)
	if err != nil {
		return "", nil, fmt.Errorf("file not found")
	}
	rc, err = s.providerFor(ctx, row).Get(ctx, key)
	if err != nil {
		return "", nil, err
	}
	return row.ContentType, rc, nil
}

// AssertOwnedByCurrent 校验文件归属（门户端仅本人可访问；对齐 hei-boot assertOwnedByCurrent）。
func (s *Service) AssertOwnedByCurrent(row *File, accountID string) error {
	if row == nil {
		return fmt.Errorf("file not found")
	}
	if row.CreatedBy == nil || *row.CreatedBy == "" || *row.CreatedBy != accountID {
		return fmt.Errorf("无权访问该文件")
	}
	return nil
}

// withResolvedURL 按行 storage_provider 重算访问 URL（库中 url 可能为已过期预签名）。
func (s *Service) withResolvedURL(ctx context.Context, row *File) *File {
	if row == nil || row.ObjectName == "" {
		return row
	}
	key := toObjectKey(row.ObjectName, s.publicPath())
	if key == "" {
		return row
	}
	if u := s.providerFor(ctx, row).PublicURL(key); u != "" {
		row.URL = u
	}
	return row
}

// providerFor 按行存储提供商解析引擎（缺省回退活动引擎；对齐 hei-boot storageFor）。
func (s *Service) providerFor(ctx context.Context, row *File) storage.Provider {
	if row != nil && strings.TrimSpace(row.StorageProvider) != "" {
		return s.sto.ProviderByName(ctx, row.StorageProvider)
	}
	return s.sto.Provider()
}

// fillNames 批量回填 created_name / updated_name（ACCOUNT 登录标识，对齐 hei-boot easy-trans）。
func (s *Service) fillNames(ctx context.Context, rows []File) {
	ids := make([]string, 0, len(rows)*2)
	seen := map[string]struct{}{}
	add := func(v *string) {
		if v == nil || *v == "" {
			return
		}
		if _, ok := seen[*v]; ok {
			return
		}
		seen[*v] = struct{}{}
		ids = append(ids, *v)
	}
	for i := range rows {
		add(rows[i].CreatedBy)
		add(rows[i].UpdatedBy)
	}
	names := s.repo.LoadAccountNames(ctx, ids)
	for i := range rows {
		if rows[i].CreatedBy != nil {
			if n, ok := names[*rows[i].CreatedBy]; ok {
				rows[i].CreatedName = &n
			}
		}
		if rows[i].UpdatedBy != nil {
			if n, ok := names[*rows[i].UpdatedBy]; ok {
				rows[i].UpdatedName = &n
			}
		}
	}
}

// publicPath 公开路径前缀（默认 /api/v1/files）。
func (s *Service) publicPath() string {
	if s.sto != nil {
		if p := s.sto.PublicPath(); p != "" {
			return p
		}
	}
	return "/api/v1/files"
}

// validateUpload 按运行时上传限制校验（STORAGE_UPLOAD_*：JSON 数组或逗号分隔；缺省 10MB）。
func (s *Service) validateUpload(ctx context.Context, fh *multipart.FileHeader) error {
	maxBytes := s.int(ctx, "STORAGE_UPLOAD_MAX_BYTES", 10*1024*1024)
	if fh.Size > int64(maxBytes) {
		return fmt.Errorf("文件大小超过限制（最大 %d 字节）", maxBytes)
	}
	originalName := safeOriginalName(fh.Filename)
	ext := strings.ToLower(path.Ext(originalName))
	ct := strings.ToLower(strings.TrimSpace(fh.Header.Get("Content-Type")))

	denied := s.strList(ctx, "STORAGE_UPLOAD_DENIED_EXTENSIONS")
	for _, d := range denied {
		if normalizeExt(d) == ext {
			return fmt.Errorf("不允许上传该文件类型：%s", ext)
		}
	}
	allowedExts := s.strList(ctx, "STORAGE_UPLOAD_ALLOWED_EXTENSIONS")
	if len(allowedExts) > 0 && !containsFold(allowedExts, ext) {
		return fmt.Errorf("不允许上传该扩展名：%s", ext)
	}
	allowedTypes := s.strList(ctx, "STORAGE_UPLOAD_ALLOWED_CONTENT_TYPES")
	if ct != "" && len(allowedTypes) > 0 && !containsFold(allowedTypes, ct) {
		return fmt.Errorf("不允许上传该内容类型：%s", ct)
	}
	return nil
}

// buildObjectName 形如 uploads/YYYY/MM/DD/{uuid}{ext}（category 长度受 STORAGE_UPLOAD_CATEGORY_MAX_LENGTH 约束）。
func (s *Service) buildObjectName(ctx context.Context, filename, category string) string {
	maxCategoryLen := s.int(ctx, "STORAGE_UPLOAD_CATEGORY_MAX_LENGTH", 64)
	if maxCategoryLen < 1 {
		maxCategoryLen = 1
	}
	category = strings.TrimSpace(category)
	if category == "" {
		category = "uploads"
	}
	if len(category) > maxCategoryLen {
		category = category[:maxCategoryLen]
	}
	ext := strings.ToLower(path.Ext(filename))
	now := time.Now().UTC()
	return fmt.Sprintf("%s/%04d/%02d/%02d/%s%s",
		category, now.Year(), now.Month(), now.Day(),
		strings.ReplaceAll(uuid.NewString(), "-", ""), ext)
}

func (s *Service) str(ctx context.Context, key string) string {
	if s.runtime == nil {
		return ""
	}
	return s.runtime.GetString(ctx, key, "")
}

func (s *Service) int(ctx context.Context, key string, def int) int {
	if s.runtime == nil {
		return def
	}
	return s.runtime.GetInt(ctx, key, def)
}

// strList 读取列表配置：兼容 JSON 数组与逗号分隔两种形态。
func (s *Service) strList(ctx context.Context, key string) []string {
	raw := strings.TrimSpace(s.str(ctx, key))
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "[") {
		var list []string
		if err := json.Unmarshal([]byte(raw), &list); err == nil {
			return list
		}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func normalizeExt(ext string) string {
	v := strings.ToLower(strings.TrimSpace(ext))
	if v != "" && !strings.HasPrefix(v, ".") {
		return "." + v
	}
	return v
}

func containsFold(list []string, v string) bool {
	for _, it := range list {
		if strings.EqualFold(strings.TrimSpace(it), v) {
			return true
		}
	}
	return false
}

func safeOriginalName(filename string) string {
	safe := strings.ReplaceAll(filename, "\\", "/")
	if idx := strings.LastIndex(safe, "/"); idx >= 0 {
		safe = safe[idx+1:]
	}
	if strings.TrimSpace(safe) == "" {
		return "file"
	}
	return safe
}

func sanitizeContentType(ct string) string {
	ct = strings.TrimSpace(ct)
	if ct == "" {
		return "application/octet-stream"
	}
	return ct
}

// isExternalURL 是否为 http(s)/data/blob 外部引用。
func isExternalURL(value string) bool {
	u, err := url.Parse(strings.TrimSpace(value))
	if err != nil || u.Scheme == "" {
		return false
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https", "data", "blob":
		return true
	}
	return false
}

// NormalizeObjectName 规范化对象名（对齐 hei-boot FileAccessUrls.normalizeObjectName）：外部 URL 原样返回，
// 否则去掉公开路径前缀后返回纯 object key（供 banner/feedback 等落库前规范化图片/附件引用）。
func NormalizeObjectName(value, publicPath string) string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return ""
	}
	if isExternalURL(raw) {
		return raw
	}
	return toObjectKey(raw, publicPath)
}

// toObjectKey 把任意对象引用（纯 key / /api/v1/files/... 路径 / 完整 URL）转成存储纯 key。
func toObjectKey(raw, publicPath string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	pathOnly := raw
	if strings.Contains(raw, "://") {
		if u, err := url.Parse(raw); err == nil {
			pathOnly = u.Path
		}
	}
	pathOnly = strings.ReplaceAll(pathOnly, "\\", "/")
	prefix := strings.TrimRight(publicPath, "/") + "/"
	if strings.HasPrefix(pathOnly, prefix) {
		pathOnly = pathOnly[len(prefix):]
	} else if strings.TrimRight(pathOnly, "/") == strings.TrimRight(publicPath, "/") {
		return ""
	}
	key := strings.TrimLeft(pathOnly, "/")
	if key == "" {
		return ""
	}
	if strings.Contains(key, "..") {
		return ""
	}
	return key
}

// CleanupManaged 跨模块清理托管文件（换头像/删附件等）：跳过外部 URL，删存储对象并清元数据。
func CleanupManaged(ctx context.Context, db *gorm.DB, sto *storage.Manager, objectName string) error {
	name := strings.TrimSpace(objectName)
	if name == "" || isExternalURL(name) {
		return nil
	}
	key := toObjectKey(name, sto.PublicPath())
	if key == "" {
		return nil
	}
	var row File
	if err := db.WithContext(ctx).Where("object_name = ?", key).First(&row).Error; err != nil {
		return nil
	}
	p := sto.ProviderByName(ctx, row.StorageProvider)
	_ = p.Delete(ctx, key)
	return db.WithContext(ctx).Where("id = ?", row.ID).Delete(&File{}).Error
}

// validObjectName 公开访问路径安全校验（对齐 hei-boot publicDownload）。
func validObjectName(objectName string) bool {
	if objectName == "" {
		return false
	}
	if strings.Contains(objectName, "..") || strings.HasPrefix(objectName, "/") || strings.Contains(objectName, "\\") {
		return false
	}
	return true
}
