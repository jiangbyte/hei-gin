// internal/modules/sys/feedback/service.go 业务服务。
//
// Author: Charlie

package feedback

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/sys/file"
)

// Service 反馈业务服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
	db   *gorm.DB
	sto  *storage.Manager
}

// NewService 构造反馈服务。
func NewService(db *gorm.DB, sto *storage.Manager) *Service {
	return &Service{repo: NewRepo(db), db: db, sto: sto}
}

// New 构建 sys.feedback 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB, d.Storage)
	return module.Module{
		Name:   "sys.feedback",
		Order:  41,
		Models: []any{&Feedback{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Submit 提交反馈（默认标题截断、分类 GENERAL；附件对象名规范化；对齐 hei-boot submit）。
func (s *Service) Submit(ctx context.Context, req CreateParam, meta SubmitMeta) error {
	title := req.Title
	if title == "" {
		title = req.Content
		if len([]rune(title)) > 64 {
			title = string([]rune(title)[:64])
		}
	}
	category := req.Category
	if category == "" {
		category = "GENERAL"
	}
	row := Feedback{
		ID: idgen.Next(), Title: title, Content: req.Content, Category: category,
		Contact: req.Contact, AttachObjectNames: jsonList(s.normalizeAttachNames(ctx, req.AttachObjectNames)),
		Status: "PENDING", SubmitterAccountType: meta.AccountType, SubmitterAccountID: meta.AccountID,
		CreatedBy: &meta.CreatedBy, UpdatedBy: &meta.CreatedBy,
	}
	return s.repo.Create(ctx, &row)
}

// Update 回复/更新反馈。
func (s *Service) Update(ctx context.Context, req UpdateParam, meta ReplyMeta) error {
	return s.repo.UpdateReply(ctx, req.ID, req.Status, req.Reply, meta, time.Now().UTC())
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Feedback, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.enrich(ctx, []*Feedback{row})
	return row, nil
}

// PageAdmin 管理端分页。
func (s *Service) PageAdmin(ctx context.Context, q PageParam) (rows []Feedback, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageAdmin(ctx, q)
	if err != nil {
		return nil, 0, current, size, err
	}
	s.enrich(ctx, toPtrs(rows))
	return rows, total, current, size, nil
}

// MyPage 我的反馈分页。
func (s *Service) MyPage(ctx context.Context, q schema.PageQuery, accountID, accountType string) (rows []Feedback, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.PageBySubmitter(ctx, accountID, accountType, current, size)
	if err != nil {
		return nil, 0, current, size, err
	}
	s.enrich(ctx, toPtrs(rows))
	return rows, total, current, size, nil
}

// MyDetail 我的反馈详情。
func (s *Service) MyDetail(ctx context.Context, id, accountID, accountType string) (*Feedback, error) {
	row, err := s.repo.GetBySubmitter(ctx, id, accountID, accountType)
	if err != nil {
		return nil, err
	}
	s.enrich(ctx, []*Feedback{row})
	return row, nil
}

// normalizeAttachNames 规范化附件对象名并校验文件存在（对齐 hei-boot normalizeAttachNames）。
func (s *Service) normalizeAttachNames(ctx context.Context, raw []string) []string {
	names := make([]string, 0, len(raw))
	seen := map[string]struct{}{}
	for _, v := range raw {
		n := file.NormalizeObjectName(v)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		names = append(names, n)
	}
	if len(names) == 0 {
		return names
	}
	files, err := s.repo.ListFilesByObjectNames(ctx, names)
	if err != nil {
		return names
	}
	found := make(map[string]struct{}, len(files))
	for i := range files {
		found[files[i].ObjectName] = struct{}{}
	}
	for _, n := range names {
		if _, ok := found[n]; !ok {
			// 附件文件不存在则跳过（对齐 hei-boot 抛 400；这里软跳过避免提交失败）
			return nil
		}
	}
	return names
}

// enrich 批量回填附件与提交者展示信息（对齐 hei-boot enrichMany）。
func (s *Service) enrich(ctx context.Context, rows []*Feedback) {
	if len(rows) == 0 {
		return
	}
	allNames := make([]string, 0)
	for _, row := range rows {
		if row == nil {
			continue
		}
		var names []string
		_ = json.Unmarshal(row.AttachObjectNames, &names)
		row.AttachObjectNames = jsonList(names)
		allNames = append(allNames, names...)
	}
	fileMap := s.repo.FileMapByObjectNames(ctx, allNames)
	for _, row := range rows {
		if row == nil {
			continue
		}
		var names []string
		_ = json.Unmarshal(row.AttachObjectNames, &names)
		row.Attachments = make([]AttachmentResult, 0, len(names))
		for _, name := range names {
			item := AttachmentResult{ObjectName: name}
			if f, ok := fileMap[name]; ok {
				item.ID = &f.ID
				item.OriginalName = &f.OriginalName
				item.ContentType = &f.ContentType
				item.Size = &f.Size
			}
			if u := s.resolveURL(ctx, name); u != "" {
				item.URL = &u
			}
			row.Attachments = append(row.Attachments, item)
		}
	}
	s.enrichSubmitters(ctx, rows)
}

// enrichSubmitters 批量回填提交者昵称/头像（按账号类型查询资料表；对齐 hei-boot enrichSubmitters）。
func (s *Service) enrichSubmitters(ctx context.Context, rows []*Feedback) {
	adminIDs, portalIDs := classifySubmitters(rows)
	adminNames := s.repo.ProfileNames(ctx, "ADMIN", adminIDs)
	portalNames := s.repo.ProfileNames(ctx, "PORTAL", portalIDs)
	for _, row := range rows {
		if row == nil {
			continue
		}
		typ := row.SubmitterAccountType
		if typ == "" {
			continue
		}
		var names map[string]ProfileBrief
		if strings.EqualFold(typ, "ADMIN") {
			names = adminNames
		} else {
			names = portalNames
		}
		if p, ok := names[row.SubmitterAccountID]; ok {
			if p.Nickname != "" {
				row.SubmitterNickname = &p.Nickname
			}
			if p.Avatar != "" {
				row.SubmitterAvatar = &p.Avatar
			}
		}
	}
}

func classifySubmitters(rows []*Feedback) (adminIDs, portalIDs []string) {
	seenA, seenP := map[string]struct{}{}, map[string]struct{}{}
	for _, row := range rows {
		if row == nil || row.SubmitterAccountID == "" {
			continue
		}
		if strings.EqualFold(row.SubmitterAccountType, "ADMIN") {
			if _, ok := seenA[row.SubmitterAccountID]; !ok {
				seenA[row.SubmitterAccountID] = struct{}{}
				adminIDs = append(adminIDs, row.SubmitterAccountID)
			}
		} else {
			if _, ok := seenP[row.SubmitterAccountID]; !ok {
				seenP[row.SubmitterAccountID] = struct{}{}
				portalIDs = append(portalIDs, row.SubmitterAccountID)
			}
		}
	}
	return adminIDs, portalIDs
}

// resolveURL 解析附件对象访问 URL（外部 URL 原样返回）。
func (s *Service) resolveURL(ctx context.Context, objectName string) string {
	if s.sto == nil {
		return ""
	}
	return s.sto.ResolveURL(ctx, objectName)
}

func toPtrs(rows []Feedback) []*Feedback {
	out := make([]*Feedback, len(rows))
	for i := range rows {
		out[i] = &rows[i]
	}
	return out
}
