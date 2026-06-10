package org

import (
	"context"
	"sort"

	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	groupModel "hei-gin/plugins/plugin-sys/group"
	posModel "hei-gin/plugins/plugin-sys/position"
	userModel "hei-gin/plugins/plugin-sys/user"

	"github.com/gin-gonic/gin"
)

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

func OrgPage(c *gin.Context, p *OrgPageParam) {
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

	q := db.DB.WithContext(ctx).Model(&SysOrg{})
	if p.ParentID != "" {
		if p.ParentID == "0" {
			q = q.Where("(parent_id IS NULL OR parent_id = '' OR id = ?)", p.ParentID)
		} else {
			q = q.Where("(id = ? OR parent_id = ?)", p.ParentID, p.ParentID)
		}
	}
	if p.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+p.Keyword+"%")
	}

	var total int64
	q.Count(&total)

	var rows []SysOrg
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*OrgVO, len(rows))
	for i, r := range rows {
		vos[i] = SysOrgToOrgVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func OrgTree(c *gin.Context, p *OrgTreeParam) []map[string]interface{} {
	ctx := c.Request.Context()
	q := db.DB.WithContext(ctx).Model(&SysOrg{}).Order("sort_code ASC")
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}

	var all []SysOrg
	q.Find(&all)

	if len(all) == 0 {
		return make([]map[string]interface{}, 0)
	}

	nodeMap := make(map[string]map[string]interface{}, len(all))
	for _, e := range all {
		entry := e
		node := map[string]interface{}{
			"id":        entry.ID,
			"code":      entry.Code,
			"name":      entry.Name,
			"category":  entry.Category,
			"status":    entry.Status,
			"sort_code": entry.SortCode,
			"children":  make([]map[string]interface{}, 0),
		}
		if entry.ParentID != nil {
			node["parent_id"] = *entry.ParentID
		}
		if entry.Description != nil {
			node["description"] = *entry.Description
		}
		if entry.Extra != nil {
			node["extra"] = *entry.Extra
		}
		if entry.CreatedAt != nil {
			node["created_at"] = utils.FormatDateTimePtr(entry.CreatedAt)
		}
		if entry.CreatedBy != nil {
			node["created_by"] = *entry.CreatedBy
		}
		if entry.UpdatedAt != nil {
			node["updated_at"] = utils.FormatDateTimePtr(entry.UpdatedAt)
		}
		if entry.UpdatedBy != nil {
			node["updated_by"] = *entry.UpdatedBy
		}
		nodeMap[entry.ID] = node
	}

	roots := make([]map[string]interface{}, 0)
	for _, e := range all {
		node := nodeMap[e.ID]
		pid := ""
		if e.ParentID != nil && *e.ParentID != "" && *e.ParentID != "0" {
			pid = *e.ParentID
		}
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

func OrgCreate(c *gin.Context, vo *OrgVO) {
	ctx := c.Request.Context()

	if vo.Code == "" || vo.Name == "" || vo.Category == "" {
		result.WriteError(c, exception.NewBusinessError("组织编码、名称、类别不能为空", 400))
		return
	}

	e := OrgVOToSysOrg(vo)
	e.Status = string(enums.StatusEnabled)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加组织失败: "+err.Error(), 500))
		return
	}
}

func OrgModify(c *gin.Context, vo *OrgVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e SysOrg
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询组织失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"code": vo.Code, "name": vo.Name, "category": vo.Category,
		"sort_code": vo.SortCode,
	}
	if vo.ParentID != nil {
		up["parent_id"] = *vo.ParentID
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	}
	if vo.Extra != nil {
		up["extra"] = *vo.Extra
	}
	if err := db.DB.WithContext(ctx).Model(&SysOrg{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑组织失败: "+err.Error(), 500))
		return
	}
}

func OrgRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()

	allIDs := collectDescendantOrgIDs(ctx, ids)

	var userCount int64
	db.DB.WithContext(ctx).Model(&userModel.SysUser{}).Where("org_id IN ?", allIDs).Count(&userCount)
	if userCount > 0 {
		result.WriteError(c, exception.NewBusinessError("组织存在关联用户，无法删除", 400))
		return
	}

	var groupCount int64
	db.DB.WithContext(ctx).Model(&groupModel.SysGroup{}).Where("org_id IN ?", allIDs).Count(&groupCount)
	if groupCount > 0 {
		result.WriteError(c, exception.NewBusinessError("组织存在关联用户组，无法删除", 400))
		return
	}

	var posCount int64
	db.DB.WithContext(ctx).Model(&posModel.SysPosition{}).Where("org_id IN ?", allIDs).Count(&posCount)
	if posCount > 0 {
		result.WriteError(c, exception.NewBusinessError("组织存在关联职位，无法删除", 400))
		return
	}

	if err := db.DB.WithContext(ctx).Where("id IN ?", allIDs).Delete(&SysOrg{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除组织失败: "+err.Error(), 500))
		return
	}
}

func OrgDetail(c *gin.Context, id string) *OrgVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysOrg
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询组织详情失败: "+err.Error(), 500))
		return nil
	}
	return SysOrgToOrgVO(&e)
}

func OrgOptions(c *gin.Context) []*OrgVO {
	ctx := c.Request.Context()
	var rows []SysOrg
	db.DB.WithContext(ctx).Order("sort_code ASC").Find(&rows)
	vos := make([]*OrgVO, len(rows))
	for i, r := range rows {
		vos[i] = SysOrgToOrgVO(&r)
	}
	return vos
}

func collectDescendantOrgIDs(ctx context.Context, ids []string) []string {
	allIDs := make(map[string]bool)
	for _, id := range ids {
		allIDs[id] = true
	}

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
	for id := range allIDs {
		result = append(result, id)
	}
	return result
}
