package resource

import (
	"context"
	"encoding/json"
	"log"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
)

type relRoleResource struct {
	ID         string
	RoleID     string
	ResourceID string
}

type relRolePermission struct {
	ID             string
	RoleID         string
	PermissionCode string
}

func ModulePage(c *gin.Context, param *ModulePageParam) {
	ctx := c.Request.Context()
	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 {
		param.Size = 10
	}
	if param.Size > 100 {
		param.Size = 100
	}

	var total int64
	db.DB.WithContext(ctx).Model(&SysModule{}).Count(&total)

	var records []SysModule
	db.DB.WithContext(ctx).Order("created_at DESC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)

	vos := make([]*ModuleVO, len(records))
	for i, r := range records {
		vos[i] = SysModuleToModuleVO(&r)
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func ModuleDetail(c *gin.Context, id string) *ModuleVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysModule
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询模块详情失败: "+err.Error(), 500))
		return nil
	}
	return SysModuleToModuleVO(&e)
}

func ModuleCreate(c *gin.Context, vo *ModuleVO) {
	ctx := c.Request.Context()

	e := ModuleVOToSysModule(vo)
	if e.Status == "" {
		e.Status = string(enums.StatusEnabled)
	}
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加模块失败: "+err.Error(), 500))
		return
	}
}

func ModuleModify(c *gin.Context, vo *ModuleVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var entity SysModule
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询模块失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{"code": vo.Code, "name": vo.Name, "category": vo.Category}
	if vo.SortCode != 0 {
		up["sort_code"] = vo.SortCode
	} else {
		up["sort_code"] = 0
	}
	if vo.Icon != nil {
		up["icon"] = *vo.Icon
	} else {
		up["icon"] = nil
	}
	if vo.Color != nil {
		up["color"] = *vo.Color
	} else {
		up["color"] = nil
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	} else {
		up["description"] = nil
	}
	if vo.IsVisible != "" {
		up["is_visible"] = vo.IsVisible
	}
	if vo.Status != "" {
		up["status"] = vo.Status
	}
	if err := db.DB.WithContext(ctx).Model(&SysModule{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑模块失败: "+err.Error(), 500))
		return
	}
}

func ModuleRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysModule{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除模块失败: "+err.Error(), 500))
		return
	}
}

func ResourcePage(c *gin.Context, param *ResourcePageParam) {
	ctx := c.Request.Context()
	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 {
		param.Size = 10
	}
	if param.Size > 100 {
		param.Size = 100
	}

	var total int64
	db.DB.WithContext(ctx).Model(&SysResource{}).Count(&total)

	var records []SysResource
	db.DB.WithContext(ctx).Order("sort_code ASC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)

	vos := make([]*ResourceVO, len(records))
	for i, r := range records {
		vos[i] = SysResourceToResourceVO(&r)
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func ResourceTree(c *gin.Context, category string) []map[string]interface{} {
	ctx := c.Request.Context()
	q := db.DB.WithContext(ctx).Model(&SysResource{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	var all []SysResource
	q.Find(&all)

	cm := make(map[string][]SysResource)
	for _, r := range all {
		pid := ""
		if r.ParentID != nil && *r.ParentID != "" && *r.ParentID != "0" {
			pid = *r.ParentID
		}
		cm[pid] = append(cm[pid], r)
	}
	return buildRT(cm, "", 0)
}

func buildRT(cm map[string][]SysResource, pid string, depth int) []map[string]interface{} {
	if depth > 50 {
		return nil
	}
	cs := cm[pid]
	r := make([]map[string]interface{}, 0, len(cs))
	for _, c := range cs {
		n := resToNode(&c)
		n["children"] = buildRT(cm, c.ID, depth+1)
		r = append(r, n)
	}
	return r
}

func resToNode(r *SysResource) map[string]interface{} {
	n := map[string]interface{}{
		"id": r.ID, "code": r.Code, "name": r.Name, "category": r.Category, "type": r.Type,
		"route_path": r.RoutePath, "component_path": r.ComponentPath, "redirect_path": r.RedirectPath,
		"icon": r.Icon, "color": r.Color, "is_visible": r.IsVisible, "is_cache": r.IsCache,
		"is_affix": r.IsAffix, "is_breadcrumb": r.IsBreadcrumb, "external_url": r.ExternalURL,
		"sort_code": r.SortCode, "status": r.Status,
	}
	if r.ParentID != nil {
		n["parent_id"] = *r.ParentID
	} else {
		n["parent_id"] = nil
	}
	if r.Description != nil {
		n["description"] = *r.Description
	}
	if r.Extra != nil {
		n["extra"] = *r.Extra
	}
	return n
}

func ResourceMenu(c *gin.Context) []map[string]interface{} {
	ctx := c.Request.Context()
	var all []SysResource
	db.DB.WithContext(ctx).Model(&SysResource{}).Order("sort_code ASC").Find(&all)
	cm := make(map[string][]SysResource)
	for _, r := range all {
		pid := ""
		if r.ParentID != nil && *r.ParentID != "" && *r.ParentID != "0" {
			pid = *r.ParentID
		}
		cm[pid] = append(cm[pid], r)
	}

	roots := cm[""]
	r := make([]map[string]interface{}, 0, len(roots))
	for _, rt := range roots {
		n := menuNode(&rt)
		n["children"] = buildMenuTree(cm, rt.ID, 0)
		r = append(r, n)
	}
	return r
}

func buildMenuTree(cm map[string][]SysResource, pid string, depth int) []map[string]interface{} {
	if depth > 50 {
		return nil
	}
	cs := cm[pid]
	if len(cs) == 0 {
		return []map[string]interface{}{}
	}
	r := make([]map[string]interface{}, 0, len(cs))
	for _, c := range cs {
		n := menuNode(&c)
		n["children"] = buildMenuTree(cm, c.ID, depth+1)
		r = append(r, n)
	}
	return r
}

func menuNode(r *SysResource) map[string]interface{} {
	n := map[string]interface{}{
		"id": r.ID, "code": r.Code, "name": r.Name, "type": r.Type, "category": r.Category,
		"route_path": r.RoutePath, "component_path": r.ComponentPath, "redirect_path": r.RedirectPath,
		"icon": r.Icon, "color": r.Color, "is_visible": r.IsVisible, "is_cache": r.IsCache,
		"is_affix": r.IsAffix, "is_breadcrumb": r.IsBreadcrumb, "external_url": r.ExternalURL,
		"sort_code": r.SortCode, "status": r.Status,
	}
	if r.ParentID != nil {
		n["parent_id"] = *r.ParentID
	} else {
		n["parent_id"] = nil
	}
	if r.Description != nil {
		n["description"] = *r.Description
	}
	return n
}

func ResourceDetail(c *gin.Context, id string) *ResourceVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysResource
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询资源详情失败: "+err.Error(), 500))
		return nil
	}
	return SysResourceToResourceVO(&e)
}

func ResourceCreate(c *gin.Context, vo *ResourceVO) {
	ctx := c.Request.Context()

	e := ResourceVOToSysResource(vo)
	if e.Status == "" {
		e.Status = string(enums.StatusEnabled)
	}
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加资源失败: "+err.Error(), 500))
		return
	}
}

func ResourceModify(c *gin.Context, vo *ResourceVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var old SysResource
	if err := db.DB.WithContext(ctx).First(&old, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询资源失败: "+err.Error(), 500))
		return
	}

	oldExtra := old.Extra
	up := map[string]interface{}{
		"code": vo.Code, "name": vo.Name, "category": vo.Category, "type": vo.Type,
		"sort_code": vo.SortCode,
	}
	if vo.ParentID != nil {
		up["parent_id"] = *vo.ParentID
	} else {
		up["parent_id"] = nil
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	} else {
		up["description"] = nil
	}
	if vo.RoutePath != nil {
		up["route_path"] = *vo.RoutePath
	} else {
		up["route_path"] = nil
	}
	if vo.ComponentPath != nil {
		up["component_path"] = *vo.ComponentPath
	} else {
		up["component_path"] = nil
	}
	if vo.RedirectPath != nil {
		up["redirect_path"] = *vo.RedirectPath
	} else {
		up["redirect_path"] = nil
	}
	if vo.Icon != nil {
		up["icon"] = *vo.Icon
	} else {
		up["icon"] = nil
	}
	if vo.Color != nil {
		up["color"] = *vo.Color
	} else {
		up["color"] = nil
	}
	if vo.IsVisible != "" {
		up["is_visible"] = vo.IsVisible
	}
	if vo.IsCache != "" {
		up["is_cache"] = vo.IsCache
	}
	if vo.IsAffix != "" {
		up["is_affix"] = vo.IsAffix
	}
	if vo.IsBreadcrumb != "" {
		up["is_breadcrumb"] = vo.IsBreadcrumb
	}
	if vo.ExternalURL != nil {
		up["external_url"] = *vo.ExternalURL
	} else {
		up["external_url"] = nil
	}
	if vo.Extra != nil {
		up["extra"] = *vo.Extra
	} else {
		up["extra"] = nil
	}
	if vo.Status != "" {
		up["status"] = vo.Status
	}

	if err := db.DB.WithContext(ctx).Model(&SysResource{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑资源失败: "+err.Error(), 500))
		return
	}

	syncPerm(ctx, vo.ID, oldExtra, vo.Extra)
}

func ResourceRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	all := collectDescendant(ctx, ids)

	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Table("rel_role_resource").Where("resource_id IN ?", all).Delete(nil).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除资源角色关联失败: "+err.Error(), 500))
		return
	}
	subQuery := tx.Table("rel_role_resource").Select("role_id").Where("resource_id IN ?", all)
	if err := tx.Table("rel_role_permission").Where("role_id IN (?)", subQuery).Delete(nil).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除资源权限关联失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("id IN ?", all).Delete(&SysResource{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除资源失败: "+err.Error(), 500))
		return
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}


func collectDescendant(ctx context.Context, ids []string) []string {
	m := make(map[string]bool)
	for _, id := range ids {
		m[id] = true
	}

	var all []SysResource
	db.DB.WithContext(ctx).Find(&all)
	cm := make(map[string][]string)
	for _, r := range all {
		if r.ParentID != nil && *r.ParentID != "" {
			cm[*r.ParentID] = append(cm[*r.ParentID], r.ID)
		}
	}

	q := make([]string, len(ids))
	copy(q, ids)
	for len(q) > 0 {
		pid := q[len(q)-1]
		q = q[:len(q)-1]
		for _, cid := range cm[pid] {
			if !m[cid] {
				m[cid] = true
				q = append(q, cid)
			}
		}
	}
	r := make([]string, 0, len(m))
	for id := range m {
		r = append(r, id)
	}
	return r
}

func syncPerm(ctx context.Context, resourceID string, oldExtra, newExtra *string) {
	oldCode := extractPermCode(oldExtra)
	newCode := extractPermCode(newExtra)
	if oldCode == newCode {
		return
	}

	tx := db.DB.WithContext(ctx).Begin()

	var roleResources []relRoleResource
	if err := tx.Table("rel_role_resource").Where("resource_id = ?", resourceID).Find(&roleResources).Error; err != nil {
		tx.Rollback()
		log.Printf("[RESOURCE] Failed to query role resources: %v", err)
		return
	}
	if len(roleResources) == 0 {
		tx.Rollback()
		return
	}

	roleIDs := make([]string, len(roleResources))
	for i, rr := range roleResources {
		roleIDs[i] = rr.RoleID
	}

	if oldCode != "" {
		if err := tx.Table("rel_role_permission").Where("role_id IN ? AND permission_code = ?", roleIDs, oldCode).Delete(nil).Error; err != nil {
			tx.Rollback()
			log.Printf("[RESOURCE] Failed to delete old permissions: %v", err)
			return
		}
	}

	if newCode != "" {
		var existing []struct{ RoleID string }
		if err := tx.Table("rel_role_permission").Where("role_id IN ? AND permission_code = ?", roleIDs, newCode).Select("role_id").Find(&existing).Error; err != nil {
			tx.Rollback()
			log.Printf("[RESOURCE] Failed to query existing permissions: %v", err)
			return
		}
		existingMap := make(map[string]bool)
		for _, e := range existing {
			existingMap[e.RoleID] = true
		}
		permBatch := make([]struct {
			ID             string
			RoleID         string
			PermissionCode string
		}, 0)
		for _, rid := range roleIDs {
			if existingMap[rid] {
				continue
			}
			permBatch = append(permBatch, struct {
				ID             string
				RoleID         string
				PermissionCode string
			}{utils.GenerateID(), rid, newCode})
		}
		if len(permBatch) > 0 {
			if err := tx.Table("rel_role_permission").CreateInBatches(permBatch, 100).Error; err != nil {
				tx.Rollback()
				log.Printf("[RESOURCE] Failed to batch insert permissions: %v", err)
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[RESOURCE] Failed to commit transaction: %v", err)
	}
}

func extractPermCode(extra *string) string {
	if extra == nil || *extra == "" {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(*extra), &m); err != nil {
		return ""
	}
	code, _ := m["permission_code"].(string)
	return code
}
