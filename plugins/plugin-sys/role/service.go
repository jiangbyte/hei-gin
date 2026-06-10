package role

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	userModel "hei-gin/plugins/plugin-sys/user"

	"github.com/gin-gonic/gin"
)

func RolePage(c *gin.Context, p *RolePageParam) {
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

	q := db.DB.WithContext(ctx).Model(&SysRole{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}

	var total int64
	q.Count(&total)

	var rows []SysRole
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*RoleVO, len(rows))
	for i, r := range rows {
		vos[i] = SysRoleToRoleVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func RoleCreate(c *gin.Context, vo *RoleVO) {
	ctx := c.Request.Context()

	e := RoleVOToSysRole(vo)
	e.Status = string(enums.StatusEnabled)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加角色失败: "+err.Error(), 500))
		return
	}
}

func RoleModify(c *gin.Context, vo *RoleVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e SysRole
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
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
	if err := db.DB.WithContext(ctx).Model(&SysRole{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑角色失败: "+err.Error(), 500))
		return
	}
}

func RoleRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()

	var cnt int64
	db.DB.WithContext(ctx).Model(&userModel.RelUserRole{}).Where("role_id IN ?", ids).Count(&cnt)
	if cnt > 0 {
		result.WriteError(c, exception.NewBusinessError("角色存在关联用户，无法删除", 400))
		return
	}

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelRolePermission{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除角色权限关联失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelRoleResource{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除角色资源关联失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelUserRole{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除角色用户关联失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("id IN ?", ids).Delete(&SysRole{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除角色失败: "+err.Error(), 500))
		return
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}

func RoleDetail(c *gin.Context, id string) *RoleVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysRole
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询角色详情失败: "+err.Error(), 500))
		return nil
	}
	return SysRoleToRoleVO(&e)
}

func RoleOwnPermissionCodes(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	var perms []userModel.RelRolePermission
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Select("permission_code").Find(&perms)
	codes := make([]string, len(perms))
	for i, p := range perms {
		codes[i] = p.PermissionCode
	}
	return codes
}

func RoleOwnPermissionDetails(c *gin.Context, roleID string) []map[string]interface{} {
	ctx := c.Request.Context()
	var perms []userModel.RelRolePermission
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&perms)
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

func RoleOwnResourceIDs(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	var resources []userModel.RelRoleResource
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Select("resource_id").Find(&resources)
	ids := make([]string, len(resources))
	for i, r := range resources {
		ids[i] = r.ResourceID
	}
	return ids
}

func RoleGrantPermissions(c *gin.Context, param *GrantPermissionParam) {
	roleID := param.RoleID
	if roleID == "" {
		result.WriteError(c, exception.NewBusinessError("角色ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&userModel.RelRolePermission{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除已有权限失败: "+err.Error(), 500))
		return
	}

	permBatch := make([]userModel.RelRolePermission, len(param.Permissions))
	for i, p := range param.Permissions {
		r := userModel.RelRolePermission{
			RoleID: roleID,
			PermissionCode: p.PermissionCode, 
			Scope: p.Scope,
		}
		if p.CustomScopeGroupIds != nil {
			r.CustomScopeGroupIds = p.CustomScopeGroupIds
		}
		if p.CustomScopeOrgIds != nil {
			r.CustomScopeOrgIds = p.CustomScopeOrgIds
		}
		permBatch[i] = r
	}
	if err := tx.Create(&permBatch).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
		return
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}

func RoleGrantResources(c *gin.Context, param *GrantResourceParam) {
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

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&userModel.RelRoleResource{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除已有资源权限失败: "+err.Error(), 500))
		return
	}

	rrBatch := make([]userModel.RelRoleResource, len(uIDs))
	for i, id := range uIDs {
		rrBatch[i] = userModel.RelRoleResource{
			RoleID: roleID, ResourceID: id,
		}
	}
	if err := tx.Create(&rrBatch).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("分配资源权限失败: "+err.Error(), 500))
		return
	}

	type row struct{ ID string; Extra *string }
	var res []row
	tx.Table("sys_resource").Where("id IN ? AND extra IS NOT NULL AND extra != ''", uIDs).Find(&res)

	var existingPerms []userModel.RelRolePermission
	tx.Where("role_id = ?", roleID).Select("permission_code").Find(&existingPerms)
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
	if len(permBatch) > 0 {
		if err := tx.Create(&permBatch).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}
