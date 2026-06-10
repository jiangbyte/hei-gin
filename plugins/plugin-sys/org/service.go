package org

import (
	"context"
	"sort"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	groupModel "hei-gin/plugins/plugin-sys/group"
	posModel "hei-gin/plugins/plugin-sys/position"
	userModel "hei-gin/plugins/plugin-sys/user"


	"github.com/gin-gonic/gin"
)



func orgToVOMap(entity *SysOrg) map[string]interface{} {
	node := map[string]interface{}{
		"id":        entity.ID,
		"code":      entity.Code,
		"name":      entity.Name,
		"category":  entity.Category,
		"status":    entity.Status,
		"sort_code": entity.SortCode,
		"children":  make([]map[string]interface{}, 0),
	}
	if entity.ParentID != nil { node["parent_id"] = *entity.ParentID }
	if entity.Description != nil { node["description"] = *entity.Description }
	if entity.Extra != nil { node["extra"] = *entity.Extra }
	if entity.CreatedAt != nil { node["created_at"] = entity.CreatedAt.Format("2006-01-02 15:04:05") }
	if entity.CreatedBy != nil { node["created_by"] = *entity.CreatedBy }
	if entity.UpdatedAt != nil { node["updated_at"] = entity.UpdatedAt.Format("2006-01-02 15:04:05") }
	if entity.UpdatedBy != nil { node["updated_by"] = *entity.UpdatedBy }
	return node
}

func sortTreeNodes(nodes []map[string]interface{}) {
	sort.Slice(nodes, func(i, j int) bool {
		si, _ := nodes[i]["sort_code"].(int)
		sj, _ := nodes[j]["sort_code"].(int)
		return si < sj
	})
	for _, n := range nodes {
		if children, ok := n["children"].([]map[string]interface{}); ok {
			sortTreeNodes(children)
		}
	}
}

func getParentIDKey(parentID *string) string {
	if parentID == nil || *parentID == "" || *parentID == "0" { return "" }
	return *parentID
}

func Page(c *gin.Context, param *OrgPageParam) {
	crud.Page(c, &SysOrg{}, param, func(q *gorm.DB) *gorm.DB {
		if param.ParentID != "" {
			if param.ParentID == "0" {
				q = q.Where("(parent_id IS NULL OR parent_id = '' OR id = ?)", param.ParentID)
			} else {
				q = q.Where("(id = ? OR parent_id = ?)", param.ParentID, param.ParentID)
			}
		}
		if param.Keyword != "" {
			q = q.Where("name LIKE ?", "%"+param.Keyword+"%")
		}
		return q
	}, "sort_code ASC", func(e *SysOrg) any { return toVO(e) })
}

func Tree(c *gin.Context, param *OrgTreeParam) []map[string]interface{} {
	ctx := c.Request.Context()
	query := db.DB.WithContext(ctx).Model(&SysOrg{}).Order("sort_code ASC")
	if param.Category != "" {
		query = query.Where("category = ?", param.Category)
	}

	var all []SysOrg
	query.Find(&all)

	if len(all) == 0 { return make([]map[string]interface{}, 0) }

	nodeMap := make(map[string]map[string]interface{}, len(all))
	for _, e := range all {
		entry := e
		nodeMap[entry.ID] = orgToVOMap(&entry)
	}

	roots := make([]map[string]interface{}, 0)
	for _, e := range all {
		node := nodeMap[e.ID]
		pid := getParentIDKey(e.ParentID)
		if pid == "" {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[pid]; ok {
				parent["children"] = append(parent["children"].([]map[string]interface{}), node)
			}
		}
	}

	sortTreeNodes(roots)
	return roots
}

func Create(c *gin.Context, vo *OrgVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()

	entity := SysOrg{
		ID:        utils.GenerateID(),
		Code:      vo.Code,
		Name:      vo.Name,
		Category:  vo.Category,
		Status: string(enums.StatusEnabled),
		SortCode:  vo.SortCode,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if vo.ParentID != nil { entity.ParentID = vo.ParentID }
	if vo.Description != nil { entity.Description = vo.Description }
	if vo.Extra != nil { entity.Extra = vo.Extra }
	if userID != "" {
		entity.CreatedBy = &userID
		entity.UpdatedBy = &userID
	}

	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加组织失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *OrgVO, userID string) {
	ctx := c.Request.Context()
	var entity SysOrg
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound { panic(exception.NewBusinessError("数据不存在", 400)) }
		panic(exception.NewBusinessError("查询组织失败: "+err.Error(), 500))
	}

	updates := map[string]interface{}{
		"code":       vo.Code,
		"name":       vo.Name,
		"category":   vo.Category,
		"sort_code":  vo.SortCode,
		"updated_at": time.Now(),
	}
	if vo.ParentID != nil { updates["parent_id"] = *vo.ParentID } else { updates["parent_id"] = nil }
	if vo.Description != nil { updates["description"] = *vo.Description } else { updates["description"] = nil }
	if vo.Extra != nil { updates["extra"] = *vo.Extra } else { updates["extra"] = nil }
	if userID != "" { updates["updated_by"] = userID }

	if err := db.DB.WithContext(ctx).Model(&SysOrg{}).Where("id = ?", vo.ID).Updates(updates).Error; err != nil {
		panic(exception.NewBusinessError("编辑组织失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	if len(ids) == 0 { return }
	ctx := c.Request.Context()

	allIDs := collectDescendantOrgIDs(c.Request.Context(), ids)

	var userCount int64
	db.DB.WithContext(ctx).Model(&userModel.SysUser{}).Where("org_id IN ?", allIDs).Count(&userCount)
	if userCount > 0 {
		panic(exception.NewBusinessError("组织存在关联用户，无法删除", 400))
	}

	var groupCount int64
	db.DB.WithContext(ctx).Model(&groupModel.SysGroup{}).Where("org_id IN ?", allIDs).Count(&groupCount)
	if groupCount > 0 {
		panic(exception.NewBusinessError("组织存在关联用户组，无法删除", 400))
	}

	var posCount int64
	db.DB.WithContext(ctx).Model(&posModel.SysPosition{}).Where("org_id IN ?", allIDs).Count(&posCount)
	if posCount > 0 {
		panic(exception.NewBusinessError("组织存在关联职位，无法删除", 400))
	}

	if err := db.DB.WithContext(ctx).Where("id IN ?", allIDs).Delete(&SysOrg{}).Error; err != nil {
		panic(exception.NewBusinessError("删除组织失败: "+err.Error(), 500))
	}
}

func Detail(c *gin.Context, id string) *OrgVO {
	if id == "" { return nil }
	ctx := c.Request.Context()
	var entity SysOrg
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return nil }
		panic(exception.NewBusinessError("查询组织详情失败: "+err.Error(), 500))
	}
	return toVO(&entity)
}

func Options(c *gin.Context) []*OrgVO {
	ctx := c.Request.Context()
	var records []SysOrg
	db.DB.WithContext(ctx).Order("sort_code ASC").Find(&records)
	vos := make([]*OrgVO, len(records))
	for i, r := range records { vos[i] = toVO(&r) }
	return vos
}

func collectDescendantOrgIDs(ctx context.Context, ids []string) []string {
	
	allIDs := make(map[string]bool)
	for _, id := range ids { allIDs[id] = true }

	var all []SysOrg
	db.DB.WithContext(ctx).Find(&all)
	cm := make(map[string][]string)
	for _, o := range all {
		if o.ParentID != nil && *o.ParentID != "" {
			cm[*o.ParentID] = append(cm[*o.ParentID], o.ID)
		}
	}

	queue := make([]string, len(ids))
	copy(queue, ids)

	for len(queue) > 0 {
		parentID := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		for _, childID := range cm[parentID] {
			if !allIDs[childID] {
				allIDs[childID] = true
				queue = append(queue, childID)
			}
		}
	}

	result := make([]string, 0, len(allIDs))
	for id := range allIDs { result = append(result, id) }
	return result
}
