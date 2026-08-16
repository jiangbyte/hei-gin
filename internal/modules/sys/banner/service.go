// internal/modules/sys/banner/service.go 业务服务。
//
// Author: Charlie

package banner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/sys/file"
)

// interactionDeltaKey Redis 中 Banner 互动计数增量哈希键
// （对齐 hei-fastapi banner_interaction_delta_key）。
const interactionDeltaKey = "hei:banner:interaction:deltas"

// Service Banner 业务服务。
//
// Author: Charlie
type Service struct {
	repo  *Repo
	sto   *storage.Manager
	redis *redis.Client
}

// NewService 构造 Banner 服务。
func NewService(db *gorm.DB, sto *storage.Manager, rdb *redis.Client) *Service {
	return &Service{repo: NewRepo(db), sto: sto, redis: rdb}
}

// New 构建 sys.banner 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB, d.Storage, d.Redis)
	return module.Module{
		Name:   "sys.banner",
		Models: []any{&Banner{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
		Jobs: []module.Job{
			{
				Name: "sys_banner_status_sync",
				Run:  s.bannerStatusJobHandler,
			},
			{
				Name: "sys_banner_flush_interactions",
				Run:  s.bannerFlushInteractionsHandler,
			},
		},
	}
}

// Create 创建 Banner。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	lt := req.LinkType
	if lt == "" {
		lt = "URL"
	}
	targets := req.TargetAccountTypes
	if len(targets) == 0 {
		targets = datatypes.JSON([]byte("[]"))
	}
	row := Banner{
		ID: idgen.Next(), Title: req.Title, Image: s.normalizeImage(req.Image), URL: req.URL, LinkType: lt,
		Summary: req.Summary, Description: req.Description, Category: req.Category, Type: req.Type,
		Position: req.Position, TargetAccountTypes: targets, Sort: req.Sort, Status: statusOr(req.Status),
		StartAt: req.StartAt, EndAt: req.EndAt,
	}
	return s.repo.Create(ctx, &row)
}

// normalizeImage 规范化图片对象名（对齐 hei-boot create/update 的 fileApi.normalizeObjectName）。
func (s *Service) normalizeImage(image string) string {
	if image == "" {
		return ""
	}
	return file.NormalizeObjectName(image)
}

// Update 更新 Banner。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	lt := req.LinkType
	if lt == "" {
		lt = "URL"
	}
	updates := map[string]any{
		"title": req.Title, "image": s.normalizeImage(req.Image), "url": req.URL, "link_type": lt,
		"summary": req.Summary, "description": req.Description, "category": req.Category,
		"type": req.Type, "position": req.Position, "sort": req.Sort, "status": statusOr(req.Status),
		"start_at": req.StartAt, "end_at": req.EndAt,
	}
	if len(req.TargetAccountTypes) > 0 {
		updates["target_account_types"] = req.TargetAccountTypes
	}
	return s.repo.Update(ctx, req.ID, updates)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Banner, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.withResolvedImageURL(row)
	return row, nil
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Banner, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	for i := range rows {
		s.withResolvedImageURL(&rows[i])
	}
	return rows, total, current, size, err
}

// List 管理端可见 Banner 列表。
func (s *Service) List(ctx context.Context, q ListParam) ([]Banner, error) {
	rows, err := s.repo.List(ctx, q.Position, q.Category, q.Type, "ADMIN", security.StatusEnabled)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		s.withResolvedImageURL(&rows[i])
	}
	return rows, nil
}

// Interaction 互动上报：Banner 须存在、启用且目标含 PORTAL；计数先入 Redis 增量，
// 由 bannerFlushInteractions 周期任务刷入 DB（对齐 hei-fastapi record_interaction + flush）。
func (s *Service) Interaction(ctx context.Context, id string) error {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("banner not found")
	}
	if row.Status != security.StatusEnabled {
		return errors.New("banner not enabled")
	}
	portal := false
	var targets []string
	_ = json.Unmarshal(row.TargetAccountTypes, &targets)
	for _, t := range targets {
		if t == "PORTAL" {
			portal = true
			break
		}
	}
	if !portal {
		return errors.New("banner not targeted to portal")
	}
	if s.redis == nil {
		return nil
	}
	return s.redis.HIncrBy(ctx, interactionDeltaKey, id, 1).Err()
}

// FlushInteractionDeltas 将 Redis 中的 Banner 互动增量刷入数据库并清空哈希，
// 返回处理条数（对齐 hei-fastapi flush_interaction_deltas）。
func (s *Service) FlushInteractionDeltas(ctx context.Context) (int64, error) {
	if s.redis == nil {
		return 0, nil
	}
	raw, err := s.redis.HGetAll(ctx, interactionDeltaKey).Result()
	if err != nil {
		return 0, err
	}
	if len(raw) == 0 {
		return 0, nil
	}
	type delta struct {
		id    string
		value int64
	}
	var deltas []delta
	for id, text := range raw {
		n, convErr := strconv.ParseInt(text, 10, 64)
		if convErr != nil || n <= 0 {
			continue
		}
		deltas = append(deltas, delta{id: id, value: n})
	}
	if len(deltas) == 0 {
		// 无有效增量（如非法值），清理残留键
		_ = s.redis.Del(ctx, interactionDeltaKey).Err()
		return 0, nil
	}
	ids := make([]string, 0, len(deltas))
	for _, d := range deltas {
		ids = append(ids, d.id)
		if _, incErr := s.repo.IncrementInteractionBy(ctx, d.id, d.value); incErr != nil {
			return 0, incErr
		}
	}
	_ = s.redis.HDel(ctx, interactionDeltaKey, ids...).Err()
	return int64(len(deltas)), nil
}

// bannerFlushInteractionsHandler 任务 Handler（对齐 hei-fastapi bannerFlushInteractions）。
func (s *Service) bannerFlushInteractionsHandler(ctx context.Context, _ string) (string, error) {
	n, err := s.FlushInteractionDeltas(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("flushed=%d", n), nil
}

// PortalList 门户端有效 Banner 列表。
func (s *Service) PortalList(ctx context.Context, q PortalListParam) ([]Banner, error) {
	rows, err := s.repo.ListPortal(ctx, q, security.StatusEnabled)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		s.withResolvedImageURL(&rows[i])
	}
	return rows, nil
}

// withResolvedImageURL 回填 image_url（对齐 hei-boot withResolvedImageUrl）。
func (s *Service) withResolvedImageURL(row *Banner) {
	if row == nil || row.Image == "" {
		return
	}
	if u := s.sto.ResolveURL(context.Background(), row.Image); u != "" {
		row.ImageURL = &u
	}
}

func statusOr(st string) string {
	if st == "" {
		return security.StatusEnabled
	}
	return st
}
