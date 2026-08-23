// internal/modules/workspace/service.go 业务服务（对齐 hei-boot WorkspaceServiceImpl / WorkspaceShortcutServiceImpl）。
//
// Author: Charlie

package workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/platform/audit"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

const (
	maxShortcuts = 16
	homeCode     = "workspace"
)

// Service 工作台服务。
//
// Author: Charlie
type Service struct {
	repo *Repo
}

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 workspace 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "workspace",
		Order:  50,
		Models: []any{&WorkspaceShortcut{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Overview 工作台总览（快捷应用 + 本人近期操作/登录日志）。
func (s *Service) Overview(ctx context.Context) (OverviewResult, error) {
	sess, err := requireSession(ctx)
	if err != nil {
		return OverviewResult{}, err
	}
	shortcuts, err := s.listMine(ctx, sess)
	if err != nil {
		return OverviewResult{}, err
	}
	ops, err := s.repo.ListRecentOperations(ctx, sess.AccountID)
	if err != nil {
		return OverviewResult{}, err
	}
	logins, err := s.repo.ListRecentLogins(ctx, sess.AccountID)
	if err != nil {
		return OverviewResult{}, err
	}
	return OverviewResult{
		Shortcuts:        shortcuts,
		RecentOperations: toActivityItems(ops),
		RecentLogins:     toActivityItems(logins),
	}, nil
}

// ListShortcuts 查询当前用户快捷应用。
func (s *Service) ListShortcuts(ctx context.Context) ([]ShortcutResult, error) {
	sess, err := requireSession(ctx)
	if err != nil {
		return nil, err
	}
	return s.listMine(ctx, sess)
}

// ReplaceShortcuts 整体替换当前用户快捷应用。
func (s *Service) ReplaceShortcuts(ctx context.Context, resourceIDs []string) ([]ShortcutResult, error) {
	sess, err := requireSession(ctx)
	if err != nil {
		return nil, err
	}
	normalized := normalizeResourceIDs(resourceIDs)
	if len(normalized) > maxShortcuts {
		return nil, fmt.Errorf("快捷应用最多 %d 个", maxShortcuts)
	}
	granted, fullAccess, err := s.resolveGrantedResourceIDs(ctx, sess)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	entities := make([]WorkspaceShortcut, 0, len(normalized))
	for i, resourceID := range normalized {
		menu, err := s.repo.GetMenuByID(ctx, resourceID)
		if err != nil {
			return nil, fmt.Errorf("存在不可用的菜单资源")
		}
		if menu.ResourceType != "MENU" ||
			menu.Status != security.StatusEnabled ||
			!hasText(menu.Path) ||
			homeCode == menu.Code {
			return nil, fmt.Errorf("存在不可用的菜单资源")
		}
		if !fullAccess {
			if _, ok := granted[resourceID]; !ok {
				return nil, fmt.Errorf("存在未授权的菜单：%s", menu.Name)
			}
		}
		accountID := sess.AccountID
		entities = append(entities, WorkspaceShortcut{
			ID:         idgen.Next(),
			AccountID:  sess.AccountID,
			ResourceID: resourceID,
			Sort:       i + 1,
			CreatedAt:  now,
			CreatedBy:  &accountID,
			UpdatedAt:  now,
			UpdatedBy:  &accountID,
		})
	}
	if err := s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &Repo{db: tx}
		if err := txRepo.DeleteShortcutsByAccount(ctx, sess.AccountID); err != nil {
			return err
		}
		return txRepo.CreateShortcuts(ctx, entities)
	}); err != nil {
		return nil, err
	}
	return s.listMine(ctx, sess)
}

func (s *Service) listMine(ctx context.Context, sess *security.SessionPayload) ([]ShortcutResult, error) {
	shortcuts, err := s.repo.ListShortcutsByAccount(ctx, sess.AccountID)
	if err != nil {
		return nil, err
	}
	if len(shortcuts) == 0 {
		return []ShortcutResult{}, nil
	}
	resourceIDs := make([]string, 0, len(shortcuts))
	seen := map[string]struct{}{}
	for _, row := range shortcuts {
		id := strings.TrimSpace(row.ResourceID)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		resourceIDs = append(resourceIDs, id)
	}
	menus, err := s.repo.ListMenusByIDs(ctx, resourceIDs)
	if err != nil {
		return nil, err
	}
	menuMap := make(map[string]MenuResource, len(menus))
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}
	granted, fullAccess, err := s.resolveGrantedResourceIDs(ctx, sess)
	if err != nil {
		return nil, err
	}
	out := make([]ShortcutResult, 0, len(shortcuts))
	for _, shortcut := range shortcuts {
		menu, ok := menuMap[shortcut.ResourceID]
		if !ok || !hasText(menu.Path) {
			continue
		}
		if !fullAccess {
			if _, ok := granted[shortcut.ResourceID]; !ok {
				continue
			}
		}
		out = append(out, toShortcutResult(shortcut, menu))
	}
	return out, nil
}

func (s *Service) resolveGrantedResourceIDs(ctx context.Context, sess *security.SessionPayload) (map[string]struct{}, bool, error) {
	if isFullAccess(sess.PermissionKeys) {
		return nil, true, nil
	}
	if ok, err := s.repo.HasSuperAdminRole(ctx, sess.RoleIDs); err != nil {
		return nil, false, err
	} else if ok {
		return nil, true, nil
	}
	roleIDs, err := s.repo.ListRoleIDs(ctx, sess.AccountID)
	if err != nil {
		return nil, false, err
	}
	groupIDs, err := s.repo.ListGroupIDs(ctx, sess.AccountID)
	if err != nil {
		return nil, false, err
	}
	ids, err := s.repo.ListGrantedResourceIDs(ctx, sess.AccountID, groupIDs, roleIDs, string(security.AccountAdmin))
	if err != nil {
		return nil, false, err
	}
	granted := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		granted[id] = struct{}{}
	}
	return granted, false, nil
}

func requireSession(ctx context.Context) (*security.SessionPayload, error) {
	sess := contextx.Session(ctx)
	if sess == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	return sess, nil
}

func isFullAccess(keys []string) bool {
	for _, k := range keys {
		if k == "*:*:*" {
			return true
		}
	}
	return false
}

func normalizeResourceIDs(resourceIDs []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(resourceIDs))
	for _, id := range resourceIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func hasText(s *string) bool {
	return s != nil && strings.TrimSpace(*s) != ""
}

func toShortcutResult(shortcut WorkspaceShortcut, menu MenuResource) ShortcutResult {
	path := ""
	if menu.Path != nil {
		path = *menu.Path
	}
	return ShortcutResult{
		ID:         shortcut.ID,
		ResourceID: shortcut.ResourceID,
		Sort:       shortcut.Sort,
		Name:       menu.Name,
		Path:       path,
		Icon:       menu.Icon,
		Code:       menu.Code,
	}
}

func toActivityItems(rows []AuditActivity) []ActivityItemResult {
	out := make([]ActivityItemResult, 0, len(rows))
	for _, row := range rows {
		moduleLabel := row.ModuleLabel
		actionName := row.ActionName
		actionType := row.ActionType
		rt := ""
		if row.ResourceType != nil {
			rt = *row.ResourceType
		}
		audit.EnrichActivityLabels(row.Module, rt, row.Action, &actionName, &actionType, &moduleLabel)
		out = append(out, ActivityItemResult{
			ID:           row.ID,
			Module:       row.Module,
			ModuleLabel:  moduleLabel,
			Action:       row.Action,
			ActionName:   actionName,
			ActionType:   actionType,
			Summary:      row.Summary,
			Success:      schema.WireBoolValue(row.Success),
			IP:           row.IP,
			UserAgent:    row.UserAgent,
			OperatorName: row.OperatorName,
			DurationMs:   schema.IntStringPtr(row.DurationMs),
			ResourceID:   schema.StringPtr(row.ResourceID),
			CreatedAt:    row.CreatedAt,
		})
	}
	return out
}
