package role

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	userModel "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

func (s *Service) Page(c *gin.Context, p *RolePageParam) {
	ctx := c.Request.Context()
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}

	rows, total := s.repo.Page(ctx, p)

	vos := make([]*RoleVO, len(rows))
	for i, r := range rows {
		vos[i] = SysRoleToRoleVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) Create(c *gin.Context, vo *RoleVO) {
	ctx := c.Request.Context()

	e := RoleVOToSysRole(vo)
	e.Status = string(enums.StatusEnabled)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加角色失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Modify(c *gin.Context, vo *RoleVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	if _, err := s.repo.FindByID(ctx, vo.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询角色失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"code": vo.Code, "name": vo.Name, "category": vo.Category,
		"sort_code": vo.SortCode, "updated_at": time.Now(),
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	}
	if vo.Status != "" {
		up["status"] = vo.Status
	}
	if vo.Extra != nil {
		up["extra"] = *vo.Extra
	}
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑角色失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()

	cnt := s.repo.CountUserRoleRefs(ctx, ids)
	if cnt > 0 {
		result.WriteError(c, exception.NewBusinessError("角色存在关联用户，无法删除", 400))
		return
	}

	if err := s.repo.DeleteCascade(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除角色失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Detail(c *gin.Context, id string) *RoleVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询角色详情失败: "+err.Error(), 500))
		return nil
	}
	return SysRoleToRoleVO(e)
}

func (s *Service) OwnPermissionCodes(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	perms := s.repo.FindRolePermissionCodes(ctx, roleID)
	codes := make([]string, len(perms))
	for i, p := range perms {
		codes[i] = p.PermissionCode
	}
	return codes
}

func (s *Service) OwnPermissionDetails(c *gin.Context, roleID string) []map[string]interface{} {
	ctx := c.Request.Context()
	perms := s.repo.FindRolePermissions(ctx, roleID)
	r := make([]map[string]interface{}, len(perms))
	for i, p := range perms {
		r[i] = map[string]interface{}{
			"permission_code":        p.PermissionCode,
			"scope":                  p.Scope,
			"custom_scope_group_ids": p.CustomScopeGroupIds,
			"custom_scope_org_ids":   p.CustomScopeOrgIds,
		}
	}
	return r
}

func (s *Service) OwnResourceIDs(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	resources := s.repo.FindRoleResources(ctx, roleID)
	ids := make([]string, len(resources))
	for i, r := range resources {
		ids[i] = r.ResourceID
	}
	return ids
}

func (s *Service) GrantPermissions(c *gin.Context, param *GrantPermissionParam) {
	roleID := param.RoleID
	if roleID == "" {
		result.WriteError(c, exception.NewBusinessError("角色ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()

	permBatch := make([]userModel.RelRolePermission, len(param.Permissions))
	for i, p := range param.Permissions {
		r := userModel.RelRolePermission{
			RoleID:         roleID,
			PermissionCode: p.PermissionCode,
			Scope:          p.Scope,
		}
		if p.CustomScopeGroupIds != nil {
			r.CustomScopeGroupIds = p.CustomScopeGroupIds
		}
		if p.CustomScopeOrgIds != nil {
			r.CustomScopeOrgIds = p.CustomScopeOrgIds
		}
		permBatch[i] = r
	}
	if err := s.repo.ReplaceRolePermissions(ctx, roleID, permBatch); err != nil {
		result.WriteError(c, exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
		return
	}
	s.refreshRoleUserSessionsACL(ctx, roleID)
}

func (s *Service) GrantResources(c *gin.Context, param *GrantResourceParam) {
	roleID := param.RoleID
	if roleID == "" {
		result.WriteError(c, exception.NewBusinessError("角色ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()

	uIDs := make([]string, 0)
	seen := make(map[string]bool)
	for _, id := range param.ResourceIDs {
		if !seen[id] {
			seen[id] = true
			uIDs = append(uIDs, id)
		}
	}

	rrBatch := make([]userModel.RelRoleResource, len(uIDs))
	for i, id := range uIDs {
		rrBatch[i] = userModel.RelRoleResource{
			RoleID: roleID, ResourceID: id,
		}
	}
	res := s.repo.FindResourceExtras(ctx, uIDs)
	existingPerms := s.repo.FindRolePermissionCodes(ctx, roleID)
	epm := make(map[string]bool)
	for _, p := range existingPerms {
		epm[p.PermissionCode] = true
	}

	permBatch := make([]userModel.RelRolePermission, 0)
	for _, r := range res {
		if r.Extra == nil || *r.Extra == "" {
			continue
		}
		var em map[string]interface{}
		if err := json.Unmarshal([]byte(*r.Extra), &em); err != nil {
			continue
		}
		pc, ok := em["permission_code"].(string)
		if !ok || pc == "" || epm[pc] {
			continue
		}
		permBatch = append(permBatch, userModel.RelRolePermission{
			ID: utils.GenerateID(), RoleID: roleID, PermissionCode: pc, Scope: "ALL",
		})
	}
	if err := s.repo.ReplaceRoleResourcesAndAppendPerms(ctx, roleID, rrBatch, uIDs, permBatch); err != nil {
		result.WriteError(c, exception.NewBusinessError("分配资源权限失败: "+err.Error(), 500))
		return
	}
	s.refreshRoleUserSessionsACL(ctx, roleID)
}

func (s *Service) RefreshSessionACL(c *gin.Context, param *RefreshRoleSessionACLParam) {
	if param.RoleID == "" {
		result.WriteError(c, exception.NewBusinessError("角色ID不能为空", 400))
		return
	}
	s.refreshRoleUserSessionsACL(c.Request.Context(), param.RoleID)
}

func (s *Service) refreshRoleUserSessionsACL(ctx context.Context, roleID string) {
	userIDs := s.repo.FindUserIDsByRoleID(ctx, roleID)
	seen := make(map[string]struct{}, len(userIDs))
	for _, userID := range userIDs {
		if userID == "" {
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		_ = auth.Business.RefreshUserSessionsACL(ctx, userID)
	}
}
