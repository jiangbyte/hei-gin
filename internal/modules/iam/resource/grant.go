// internal/modules/iam/resource/grant.go 授权逻辑。
//
// Author: Charlie

package resource

import (
	"context"
	"errors"
	"sort"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/modules/iam/relation"
)

// ListPermissions 返回已注册权限清单（对齐 hei-boot SysRegisteredPermissionResult：permission_key + 三段拆解）。
func (s *Service) ListPermissions() []PermissionItem {
	all := s.perms.All()
	out := make([]PermissionItem, 0, len(all))
	for k, v := range all {
		module, res, action := splitPermissionKey(k)
		out = append(out, PermissionItem{
			PermissionKey: k, Name: v,
			ModuleCode: module, ResourceCode: res, Action: action,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PermissionKey < out[j].PermissionKey })
	return out
}

// splitPermissionKey 将 a:b:c 拆成 module_code/resource_code/action。
func splitPermissionKey(key string) (module, res, action string) {
	parts := strings.Split(key, ":")
	if len(parts) > 0 {
		module = parts[0]
	}
	if len(parts) > 1 {
		res = parts[1]
	}
	if len(parts) > 2 {
		action = parts[2]
	}
	return
}

// BindResourcePermissions 为管理端资源绑定权限（先删后插，事务；对齐 hei-boot bindPermission）。
func (s *Service) BindResourcePermissions(ctx context.Context, req ResourcePermissionBindParam) error {
	if _, err := s.repo.GetResourceByID(ctx, req.ResourceID); err != nil {
		return err
	}
	keys := []string{}
	if req.PermissionKey != "" {
		keys = append(keys, req.PermissionKey)
	}
	return s.rel.BindResourcePermissionDetail(ctx, relation.SubjectResource, req.ResourceID, relation.RelResourcePermission,
		orAdmin(req.AccountType), keys, req.DataScope, req.CustomScopeDeptIDs, req.Sort, req.Description)
}

// BindClientResourcePermissions 为客户端资源绑定权限（先删后插，事务）。
func (s *Service) BindClientResourcePermissions(ctx context.Context, req ResourcePermissionBindParam) error {
	keys := []string{}
	if req.PermissionKey != "" {
		keys = append(keys, req.PermissionKey)
	}
	return s.rel.BindResourcePermissionDetail(ctx, relation.SubjectClientResource, req.ResourceID, relation.RelClientResourcePermission,
		orAdmin(req.AccountType), keys, req.DataScope, req.CustomScopeDeptIDs, req.Sort, req.Description)
}

// CreateButton 创建按钮资源并绑定权限（对齐 hei-boot createButton：建 BUTTON 后写 RESOURCE_PERMISSION）。
func (s *Service) CreateButton(ctx context.Context, req ButtonAddParam) error {
	parent, err := s.repo.GetResourceByID(ctx, req.ParentID)
	if err != nil {
		return err
	}
	row := Resource{
		ID: idgen.Next(), ParentID: &parent.ID, Code: req.Code, Name: req.Name,
		ResourceType: ResourceTypeButton, ModuleID: parent.ModuleID,
		Sort: sortOr(req.Sort), IsVisible: false, IsCache: false, IsAffix: false,
		Status: statusOr(req.Status), Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	if err := s.repo.CreateResource(ctx, &row); err != nil {
		return err
	}
	return s.rel.BindResourcePermissionDetail(ctx, relation.SubjectResource, row.ID, relation.RelResourcePermission,
		string(security.AccountAdmin), strSlice(req.PermissionKey), req.DataScope, req.CustomScopeDeptIDs, req.Sort, req.Description)
}

// UpdateButton 更新按钮资源并重建权限绑定。
func (s *Service) UpdateButton(ctx context.Context, req ButtonEditParam) error {
	row, err := s.repo.GetResourceByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if row.ResourceType != ResourceTypeButton {
		return errors.New("resource is not a button")
	}
	parent, err := s.repo.GetResourceByID(ctx, req.ParentID)
	if err != nil {
		return err
	}
	updates := map[string]any{
		"parent_id": parent.ID, "code": req.Code, "name": req.Name,
		"module_id": parent.ModuleID, "sort": sortOr(req.Sort),
		"status": statusOr(req.Status), "description": req.Description,
	}
	if err := s.repo.UpdateResource(ctx, req.ID, updates); err != nil {
		return err
	}
	return s.rel.BindResourcePermissionDetail(ctx, relation.SubjectResource, req.ID, relation.RelResourcePermission,
		string(security.AccountAdmin), strSlice(req.PermissionKey), req.DataScope, req.CustomScopeDeptIDs, req.Sort, req.Description)
}

func strSlice(v string) []string {
	if v == "" {
		return nil
	}
	return []string{v}
}

// DeleteButtons 批量删除按钮资源（先清按钮权限绑定，再删资源）。
func (s *Service) DeleteButtons(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	rows, err := s.repo.GetResourcesByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(rows) != len(unique(ids)) {
		return gorm.ErrRecordNotFound
	}
	for _, row := range rows {
		if row.ResourceType != ResourceTypeButton {
			return errors.New("resource is not a button")
		}
	}
	if err := s.rel.DeleteBySubjectIDs(ctx, relation.SubjectResource, ids, relation.RelResourcePermission); err != nil {
		return err
	}
	return s.repo.DeleteResources(ctx, ids)
}

// PageButtons 按钮资源分页（对齐 hei-boot toButtonResults：合并权限绑定与父级/模块名称）。
func (s *Service) PageButtons(ctx context.Context, p ButtonPageParam) (rows []ButtonResult, total int64, current, size int, err error) {
	current, size = p.Normalize()
	resRows, total, err := s.repo.PageButtons(ctx, p)
	if err != nil {
		return nil, 0, current, size, err
	}
	if len(resRows) == 0 {
		return []ButtonResult{}, total, current, size, nil
	}
	ids := make([]string, 0, len(resRows))
	for i := range resRows {
		ids = append(ids, resRows[i].ID)
	}
	// 批量加载权限绑定关系 + 父级名称 + 模块名称
	permMap := s.repo.ButtonPermissions(ctx, ids)
	parentRows, _ := s.repo.GetResourcesByIDs(ctx, collectParentIDs(resRows))
	parentName := map[string]string{}
	for i := range parentRows {
		parentName[parentRows[i].ID] = parentRows[i].Name
	}
	moduleRows, _ := s.repo.ModulesByIDs(ctx, collectModuleIDs(resRows))
	moduleName := map[string]string{}
	for i := range moduleRows {
		moduleName[moduleRows[i].ID] = moduleRows[i].Name
	}
	rows = make([]ButtonResult, 0, len(resRows))
	for i := range resRows {
		rr := &resRows[i]
		perm := permMap[rr.ID]
		br := ButtonResult{
			ID: rr.ID, ParentID: rr.ParentID, Code: rr.Code, Name: rr.Name,
			PermissionKey: perm.PermissionKey, PermissionDescription: perm.Description,
			DataScope: perm.DataScope, CustomScopeDeptIDs: perm.CustomScopeDeptIDs,
			ModuleID: rr.ModuleID, Sort: rr.Sort, Status: rr.Status, Description: rr.Description,
			CreatedAt: rr.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: rr.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if rr.ParentID != nil {
			if n, ok := parentName[*rr.ParentID]; ok {
				br.ParentIDName = &n
			}
		}
		if rr.ModuleID != nil {
			if n, ok := moduleName[*rr.ModuleID]; ok {
				br.ModuleIDName = &n
			}
		}
		rows = append(rows, br)
	}
	return rows, total, current, size, nil
}

type buttonPerm struct {
	PermissionKey      string   `gorm:"column:target_key"`
	Description        *string  `gorm:"column:description"`
	DataScope          string   `gorm:"column:data_scope"`
	CustomScopeDeptIDs []string `gorm:"column:custom_scope_dept_ids"`
}

func collectParentIDs(rows []Resource) []string {
	seen := map[string]struct{}{}
	var out []string
	for i := range rows {
		if rows[i].ParentID != nil {
			if _, ok := seen[*rows[i].ParentID]; !ok {
				seen[*rows[i].ParentID] = struct{}{}
				out = append(out, *rows[i].ParentID)
			}
		}
	}
	return out
}

func collectModuleIDs(rows []Resource) []string {
	seen := map[string]struct{}{}
	var out []string
	for i := range rows {
		if rows[i].ModuleID != nil {
			if _, ok := seen[*rows[i].ModuleID]; !ok {
				seen[*rows[i].ModuleID] = struct{}{}
				out = append(out, *rows[i].ModuleID)
			}
		}
	}
	return out
}

func unique(ids []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
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

func orAdmin(t string) string {
	if t == "" {
		return string(security.AccountAdmin)
	}
	return t
}
