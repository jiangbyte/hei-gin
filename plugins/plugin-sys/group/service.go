package group

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/crud"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func groupToVOMap(entity *SysGroup) map[string]interface{} {
	n := map[string]interface{}{
		"id": entity.ID, "code": entity.Code, "name": entity.Name, "category": entity.Category,
		"org_id": entity.OrgID, "status": entity.Status, "sort_code": entity.SortCode,
		"children": make([]map[string]interface{}, 0),
	}
	if entity.ParentID != nil { n["parent_id"] = *entity.ParentID }
	if entity.Description != nil { n["description"] = *entity.Description }
	return n
}

func Page(c *gin.Context, param *GroupPageParam) {
	crud.Page(c, &SysGroup{}, param, func(q *gorm.DB) *gorm.DB {
		if param.Keyword != "" { q = q.Where("name LIKE ?", "%"+param.Keyword+"%") }
		if param.Category != "" { q = q.Where("category = ?", param.Category) }
	
		if param.OrgID != "" { q = q.Where("org_id = ?", param.OrgID) }
		return q
	}, "sort_code ASC", func(e *SysGroup) any { return toVO(e) })
}

func Tree(c *gin.Context, param *GroupTreeParam) []map[string]interface{} {
	ctx := c.Request.Context()
	var all []SysGroup
	db.DB.WithContext(ctx).Order("sort_code ASC").Find(&all)
	if len(all) == 0 { return make([]map[string]interface{}, 0) }

	childrenMap := make(map[string][]SysGroup)
	for _, r := range all {
		pid := ""
		if r.ParentID != nil { pid = *r.ParentID }
		childrenMap[pid] = append(childrenMap[pid], r)
	}
	roots := childrenMap[""]
	result := make([]map[string]interface{}, 0, len(roots))
	for _, r := range roots {
		node := groupToVOMap(&r)
		buildTree(node, &r, childrenMap)
		result = append(result, node)
	}
	return result
}

func buildTree(node map[string]interface{}, parent *SysGroup, childrenMap map[string][]SysGroup) {
	children := childrenMap[parent.ID]
	if len(children) == 0 { return }
	childNodes := make([]map[string]interface{}, len(children))
	for i, c := range children {
		childNode := groupToVOMap(&c)
		buildTree(childNode, &c, childrenMap)
		childNodes[i] = childNode
	}
	node["children"] = childNodes
}

func Detail(c *gin.Context, id string) *GroupVO {
	ctx := c.Request.Context()
	var entity SysGroup
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		panic(exception.NewBusinessError("群组不存在: "+err.Error(), 500))
	}
	return toVO(&entity)
}

func Create(c *gin.Context, vo *GroupVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	entity := SysGroup{
		ID:        utils.GenerateID(),
		Code:      vo.Code,
		Name:      vo.Name,
		Category:  vo.Category,
		OrgID:     vo.OrgID,
		ParentID:  vo.ParentID,
		Status:    vo.Status,
		SortCode:  vo.SortCode,
		CreatedAt: &now,
		CreatedBy: &userID,
		UpdatedAt: &now,
		UpdatedBy: &userID,
	}
	if vo.Description != nil { entity.Description = vo.Description }
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加群组失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *GroupVO, userID string) {
	ctx := c.Request.Context()
	var entity SysGroup
	if err := db.DB.WithContext(ctx).Where("id = ?", vo.ID).First(&entity).Error; err != nil {
		panic(exception.NewBusinessError("群组不存在: "+err.Error(), 500))
	}
	now := time.Now()
	up := map[string]any{
		"code":       vo.Code,
		"name":       vo.Name,
		"category":   vo.Category,
		"org_id":     vo.OrgID,
		"parent_id":  vo.ParentID,
		"status":     vo.Status,
		"sort_code":  vo.SortCode,
		"updated_at": now,
		"updated_by": userID,
	}
	if vo.Description != nil { up["description"] = *vo.Description }
	if err := db.DB.WithContext(ctx).Model(&entity).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("编辑群组失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	if len(ids) == 0 { return }
	ctx := c.Request.Context()
	allIDs := getAllDescendantIDs(ctx, ids)
	for _, id := range allIDs {
		db.DB.WithContext(ctx).Table("sys_user").Where("group_id = ?", id).Update("group_id", nil)
	}
	db.DB.WithContext(ctx).Where("id IN ?", allIDs).Delete(&SysGroup{})
}

func Options(c *gin.Context) []*GroupVO {
	ctx := c.Request.Context()
	var records []SysGroup
	db.DB.WithContext(ctx).Order("sort_code ASC").Find(&records)
	vos := make([]*GroupVO, len(records))
	for i, r := range records { vos[i] = toVO(&r) }
	return vos
}

func GetAll(c *gin.Context) []*GroupVO {
	ctx := c.Request.Context()
	var records []SysGroup
	db.DB.WithContext(ctx).Order("sort_code ASC").Find(&records)
	vos := make([]*GroupVO, len(records))
	for i, r := range records { vos[i] = toVO(&r) }
	return vos
}

func getAllDescendantIDs(ctx context.Context, ids []string) []string {
	allIDs := make(map[string]bool)
	for _, id := range ids { allIDs[id] = true }

	var all []SysGroup
	db.DB.WithContext(ctx).Find(&all)
	cm := make(map[string][]string)
	for _, g := range all {
		if g.ParentID != nil && *g.ParentID != "" {
			cm[*g.ParentID] = append(cm[*g.ParentID], g.ID)
		}
	}

	q := make([]string, len(ids)); copy(q, ids)
	for len(q) > 0 {
		pid := q[len(q)-1]; q = q[:len(q)-1]
		for _, cid := range cm[pid] { if !allIDs[cid] { allIDs[cid] = true; q = append(q, cid) } }
	}
	r := make([]string, 0, len(allIDs)); for id := range allIDs { r = append(r, id) }; return r
}
