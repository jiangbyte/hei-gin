// internal/modules/iam/resource/grant.go 授权逻辑。
//
// Author: Charlie

package resource

import (
	"context"
	"errors"
	"sort"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/modules/iam/relation"
)

// ListPermissions 返回已注册权限清单（按权限键排序）。
func (s *Service) ListPermissions() []PermissionItem {
	all := s.perms.All()
	out := make([]PermissionItem, 0, len(all))
	for k, v := range all {
		out = append(out, PermissionItem{Key: k, Name: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// BindResourcePermissions 为管理端资源绑定权限（先删后插，事务）。
func (s *Service) BindResourcePermissions(ctx context.Context, req ResourcePermissionBindParam) error {
	if _, err := s.repo.GetResourceByID(ctx, req.ResourceID); err != nil {
		return err
	}
	return s.rel.BindResourcePermissions(ctx, relation.SubjectResource, req.ResourceID, relation.RelResourcePermission, orAdmin(req.AccountType), req.PermissionKeys)
}

// BindClientResourcePermissions 为客户端资源绑定权限（先删后插，事务）。
func (s *Service) BindClientResourcePermissions(ctx context.Context, req ResourcePermissionBindParam) error {
	return s.rel.BindResourcePermissions(ctx, relation.SubjectClientResource, req.ResourceID, relation.RelClientResourcePermission, orAdmin(req.AccountType), req.PermissionKeys)
}

// CreateButton 创建按钮资源（resource_type=BUTTON，挂在菜单下并继承模块）。
func (s *Service) CreateButton(ctx context.Context, req ButtonAddParam) error {
	parent, err := s.repo.GetResourceByID(ctx, req.ResourceID)
	if err != nil {
		return err
	}
	row := Resource{
		ID: idgen.Next(), ParentID: &parent.ID, Code: req.Code, Name: req.Name,
		ResourceType: ResourceTypeButton, ModuleID: parent.ModuleID,
		Sort: sortOr(req.Sort), IsVisible: false, IsCache: false, IsAffix: false,
		Status: statusOr(req.Status), Description: req.Description, Extra: datatypes.JSON([]byte("{}")),
	}
	return s.repo.CreateResource(ctx, &row)
}

// UpdateButton 更新按钮资源。
func (s *Service) UpdateButton(ctx context.Context, req ButtonEditParam) error {
	row, err := s.repo.GetResourceByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if row.ResourceType != ResourceTypeButton {
		return errors.New("resource is not a button")
	}
	parent, err := s.repo.GetResourceByID(ctx, req.ResourceID)
	if err != nil {
		return err
	}
	updates := map[string]any{
		"parent_id": parent.ID, "code": req.Code, "name": req.Name,
		"module_id": parent.ModuleID, "sort": sortOr(req.Sort),
		"status": statusOr(req.Status), "description": req.Description,
	}
	return s.repo.UpdateResource(ctx, req.ID, updates)
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

// PageButtons 按钮资源分页。
func (s *Service) PageButtons(ctx context.Context, p ButtonPageParam) (rows []Resource, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.PageButtons(ctx, p)
	return rows, total, current, size, err
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
