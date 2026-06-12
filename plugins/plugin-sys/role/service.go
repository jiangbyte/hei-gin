package role

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	userModel "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type service struct {
	repo *repository
}

func (s *service) RolePage(c *gin.Context, p *RolePageParam) {
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

func (s *service) RoleCreate(c *gin.Context, vo *RoleVO) {
	ctx := c.Request.Context()

	e := RoleVOToSysRole(vo)
	e.Status = string(enums.StatusEnabled)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加角色失败: "+err.Error(), 500))
		return
	}
}

func (s *service) RoleModify(c *gin.Context, vo *RoleVO) {
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

func (s *service) RoleRemove(c *gin.Context, param *utils.IdsParam) {
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

func (s *service) RoleDetail(c *gin.Context, id string) *RoleVO {
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

func (s *service) RoleOwnPermissionCodes(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	perms := s.repo.FindRolePermissionCodes(ctx, roleID)
	codes := make([]string, len(perms))
	for i, p := range perms {
		codes[i] = p.PermissionCode
	}
	return codes
}

func (s *service) RoleOwnPermissionDetails(c *gin.Context, roleID string) []map[string]interface{} {
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

func (s *service) RoleOwnResourceIDs(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	resources := s.repo.FindRoleResources(ctx, roleID)
	ids := make([]string, len(resources))
	for i, r := range resources {
		ids[i] = r.ResourceID
	}
	return ids
}

func (s *service) RoleGrantPermissions(c *gin.Context, param *GrantPermissionParam) {
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
}

func (s *service) RoleGrantResources(c *gin.Context, param *GrantResourceParam) {
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
}

func RolePage(c *gin.Context, p *RolePageParam) { defaultModule.service.RolePage(c, p) }
func RoleCreate(c *gin.Context, vo *RoleVO)     { defaultModule.service.RoleCreate(c, vo) }
func RoleModify(c *gin.Context, vo *RoleVO)     { defaultModule.service.RoleModify(c, vo) }
func RoleRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.RoleRemove(c, param)
}
func RoleDetail(c *gin.Context, id string) *RoleVO { return defaultModule.service.RoleDetail(c, id) }
func RoleOwnPermissionCodes(c *gin.Context, roleID string) []string {
	return defaultModule.service.RoleOwnPermissionCodes(c, roleID)
}
func RoleOwnPermissionDetails(c *gin.Context, roleID string) []map[string]interface{} {
	return defaultModule.service.RoleOwnPermissionDetails(c, roleID)
}
func RoleOwnResourceIDs(c *gin.Context, roleID string) []string {
	return defaultModule.service.RoleOwnResourceIDs(c, roleID)
}
func RoleGrantPermissions(c *gin.Context, param *GrantPermissionParam) {
	defaultModule.service.RoleGrantPermissions(c, param)
}
func RoleGrantResources(c *gin.Context, param *GrantResourceParam) {
	defaultModule.service.RoleGrantResources(c, param)
}
