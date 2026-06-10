package role

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	userModel "hei-gin/plugins/plugin-sys/user"


	"github.com/gin-gonic/gin"
)


func RolePage(c *gin.Context, p *RolePageParam) {
	crud.Page(c, &SysRole{}, p, func(q *gorm.DB) *gorm.DB {
		if p.Keyword != "" {
			like := "%" + p.Keyword + "%"
			q = q.Where("name LIKE ? OR code LIKE ?", like, like)
		}
		if p.Category != "" {
			q = q.Where("category = ?", p.Category)
		}
		return q
	}, "created_at DESC", func(e *SysRole) any { return toVO(e) })
}

func RoleCreate(c *gin.Context, vo *RoleVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	e := SysRole{ID: utils.GenerateID(), Code: vo.Code, Name: vo.Name, Category: vo.Category, SortCode: vo.SortCode, Status: string(enums.StatusEnabled), CreatedAt: &now, UpdatedAt: &now}
	if vo.Description != nil { e.Description = vo.Description }
	if vo.Status != "" { e.Status = vo.Status }
	if vo.Extra != nil { e.Extra = vo.Extra }
	if userID != "" { e.CreatedBy = &userID; e.UpdatedBy = &userID }
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil { panic(exception.NewBusinessError("添加角色失败: "+err.Error(), 500)) }
}

func RoleModify(c *gin.Context, vo *RoleVO, userID string) {
	ctx := c.Request.Context()
	if vo.ID == "" { panic(exception.NewBusinessError("ID不能为空", 400)) }

	var e SysRole
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound { panic(exception.NewBusinessError("数据不存在", 400)) }
		panic(exception.NewBusinessError("查询角色失败: "+err.Error(), 500))
	}

	up := map[string]interface{}{"code": vo.Code, "name": vo.Name, "category": vo.Category, "sort_code": vo.SortCode, "updated_at": time.Now()}
	if vo.Description != nil { up["description"] = *vo.Description }
	if vo.Status != "" { up["status"] = vo.Status }
	if vo.Extra != nil { up["extra"] = *vo.Extra }
	if userID != "" { up["updated_by"] = userID }
	if err := db.DB.WithContext(ctx).Model(&SysRole{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil { panic(exception.NewBusinessError("编辑角色失败: "+err.Error(), 500)) }
}

func RoleRemove(c *gin.Context, ids []string) {
	if len(ids) == 0 { return }
	ctx := c.Request.Context()

	var cnt int64
	db.DB.WithContext(ctx).Model(&userModel.RelUserRole{}).Where("role_id IN ?", ids).Count(&cnt)
	if cnt > 0 { panic(exception.NewBusinessError("角色存在关联用户，无法删除", 400)) }

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelRolePermission{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("删除角色权限关联失败: "+err.Error(), 500))
	}
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelRoleResource{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("删除角色资源关联失败: "+err.Error(), 500))
	}
	if err := tx.Where("role_id IN ?", ids).Delete(&userModel.RelUserRole{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("删除角色用户关联失败: "+err.Error(), 500))
	}
	if err := tx.Where("id IN ?", ids).Delete(&SysRole{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("删除角色失败: "+err.Error(), 500))
	}
	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
	}
}

func RoleDetail(c *gin.Context, id string) *RoleVO {
	if id == "" { return nil }
	ctx := c.Request.Context()
	var e SysRole
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return nil }
		panic(exception.NewBusinessError("查询角色详情失败: "+err.Error(), 500))
	}
	return toVO(&e)
}

func RoleAssignResource(c *gin.Context, roleID string, resourceIDs []string) {
	if roleID == "" { panic(exception.NewBusinessError("角色ID不能为空", 400)) }
	ctx := c.Request.Context()

	uIDs := make([]string, 0)
	seen := make(map[string]bool)
	for _, id := range resourceIDs { if !seen[id] { seen[id] = true; uIDs = append(uIDs, id) } }

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&userModel.RelRoleResource{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("删除已有资源权限失败: "+err.Error(), 500))
	}
	rrBatch := make([]userModel.RelRoleResource, len(uIDs))
	for i, id := range uIDs {
		rrBatch[i] = userModel.RelRoleResource{ID: utils.GenerateID(), RoleID: roleID, ResourceID: id}
	}
	if err := tx.Create(&rrBatch).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("分配资源权限失败: "+err.Error(), 500))
	}

	type rr struct{ ID string; Extra *string }
	var res []rr
	tx.Table("sys_resource").Where("id IN ? AND extra IS NOT NULL AND extra != ''", uIDs).Find(&res)

	var existingPerms []userModel.RelRolePermission
	tx.Where("role_id = ?", roleID).Select("permission_code").Find(&existingPerms)
	epm := make(map[string]bool)
	for _, p := range existingPerms { epm[p.PermissionCode] = true }

	permBatch := make([]userModel.RelRolePermission, 0)
	for _, r := range res {
		if r.Extra == nil || *r.Extra == "" { continue }
		var em map[string]interface{}
		if err := json.Unmarshal([]byte(*r.Extra), &em); err != nil { continue }
		pc, ok := em["permission_code"].(string)
		if !ok || pc == "" || epm[pc] { continue }
		permBatch = append(permBatch, userModel.RelRolePermission{ID: utils.GenerateID(), RoleID: roleID, PermissionCode: pc, Scope: "ALL"})
	}
	if len(permBatch) > 0 {
		if err := tx.Create(&permBatch).Error; err != nil {
			tx.Rollback()
			panic(exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
		}
	}
	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
	}
}

func RoleAssignPermission(c *gin.Context, roleID string, permissions []userModel.PermissionItem) {
	if roleID == "" { panic(exception.NewBusinessError("角色ID不能为空", 400)) }
	ctx := c.Request.Context()

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&userModel.RelRolePermission{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("删除已有权限失败: "+err.Error(), 500))
	}
	permBatch := make([]userModel.RelRolePermission, len(permissions))
	for i, p := range permissions {
		r := userModel.RelRolePermission{ID: utils.GenerateID(), RoleID: roleID, PermissionCode: p.PermissionCode, Scope: p.Scope}
		if p.CustomScopeGroupIds != nil { r.CustomScopeGroupIds = p.CustomScopeGroupIds }
		if p.CustomScopeOrgIds != nil { r.CustomScopeOrgIds = p.CustomScopeOrgIds }
		permBatch[i] = r
	}
	if err := tx.Create(&permBatch).Error; err != nil {
		tx.Rollback()
		panic(exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
	}
	if err := tx.Commit().Error; err != nil {
		panic(exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
	}
}

func RoleOwnPermissionCodes(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	var perms []userModel.RelRolePermission
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Select("permission_code").Find(&perms)
	codes := make([]string, len(perms))
	for i, p := range perms { codes[i] = p.PermissionCode }
	return codes
}

func RoleOwnPermissionDetails(c *gin.Context, roleID string) []map[string]interface{} {
	ctx := c.Request.Context()
	var perms []userModel.RelRolePermission
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&perms)
	r := make([]map[string]interface{}, len(perms))
	for i, p := range perms {
		r[i] = map[string]interface{}{"permission_code": p.PermissionCode, "scope": p.Scope, "custom_scope_group_ids": p.CustomScopeGroupIds, "custom_scope_org_ids": p.CustomScopeOrgIds}
	}
	return r
}

func RoleOwnResourceIDs(c *gin.Context, roleID string) []string {
	ctx := c.Request.Context()
	var resources []userModel.RelRoleResource
	db.DB.WithContext(ctx).Where("role_id = ?", roleID).Select("resource_id").Find(&resources)
	ids := make([]string, len(resources))
	for i, r := range resources { ids[i] = r.ResourceID }
	return ids
}

func RoleGrantPermissions(c *gin.Context, roleID string, permissions []PermissionItem, userID string) {
	// Convert PermissionItem to userModel.PermissionItem
	items := make([]userModel.PermissionItem, len(permissions))
	for i, p := range permissions {
		items[i] = userModel.PermissionItem{
			PermissionCode:      p.PermissionCode,
			Scope:               p.Scope,
			CustomScopeGroupIds: p.CustomScopeGroupIds,
			CustomScopeOrgIds:   p.CustomScopeOrgIds,
		}
	}
	RoleAssignPermission(c, roleID, items)
}

func RoleGrantResources(c *gin.Context, roleID string, resourceIDs []string, permissions []ButtonPermissionScope) {
	RoleAssignResource(c, roleID, resourceIDs)
}
