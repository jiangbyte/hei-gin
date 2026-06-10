package dict

import (
	"context"
	"sort"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/crud"
	"hei-gin/sdk/exception"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)



// Page handles GET /api/v1/sys/dict/page
func Page(c *gin.Context, p *DictPageParam) {
	crud.Page(c, &SysDict{}, p, func(q *gorm.DB) *gorm.DB {
		if p.Keyword != "" {
			like := "%" + p.Keyword + "%"
			q = q.Where("code LIKE ? OR label LIKE ? OR value LIKE ?", like, like, like)
		}
		if p.Category != "" {
			q = q.Where("category = ?", p.Category)
		}
		if p.ParentID != "" {
			q = q.Where("id = ? OR parent_id = ?", p.ParentID, p.ParentID)
		}
		if p.DictGroup == "FRM" {
			q = q.Where("category = ?", "FRM")
		}
		if p.DictGroup == "BIZ" {
			q = q.Where("category = ?", "BIZ")
		}
		return q
	}, "created_at DESC", func(e *SysDict) any { return entToVO(e) })
}

func Tree(c *gin.Context, param *DictTreeParam) []map[string]interface{} {
	ctx := c.Request.Context()
	query := db.DB.WithContext(ctx).Model(&SysDict{}).Order("sort_code ASC")
	if param.Category != "" { query = query.Where("category = ?", param.Category) }
	if param.DictGroup == "FRM" {
		query = query.Where("category = ?", "FRM")
	}
	if param.DictGroup == "BIZ" {
		query = query.Where("category = ?", "BIZ")
	}
	var all []SysDict
	query.Find(&all)
	if len(all) == 0 { return make([]map[string]interface{}, 0) }

	childrenMap := make(map[string][]SysDict)
	for _, r := range all { pid := getParentIDKey(r.ParentID); childrenMap[pid] = append(childrenMap[pid], r) }
	roots := childrenMap[""]
	result := make([]map[string]interface{}, 0, len(roots))
	for _, r := range roots { node := entityToNode(&r); node["children"] = buildTreeChildren(childrenMap, r.ID, 0); result = append(result, node) }
	sortTreeNodes(result)
	return result
}

func buildTreeChildren(childrenMap map[string][]SysDict, parentID string, depth int) []map[string]interface{} {
	if depth > 50 {
		return nil
	}
	children := childrenMap[parentID]
	if len(children) == 0 { return []map[string]interface{}{} }
	result := make([]map[string]interface{}, 0, len(children))
	for _, r := range children { node := entityToNode(&r); node["children"] = buildTreeChildren(childrenMap, r.ID, depth+1); result = append(result, node) }
	sortTreeNodes(result)
	return result
}

func Create(c *gin.Context, vo *DictVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	dictCheckDuplicate(ctx, vo, "")
	dictCheckCircularParent(c.Request.Context(), "", safeStrPtr(vo.ParentID))

	entity := SysDict{ID: utils.GenerateID(), Code: vo.Code, Status: string(enums.StatusEnabled), SortCode: vo.SortCode, CreatedAt: &now, UpdatedAt: &now}
	if vo.Label != nil { entity.Label = vo.Label }
	if vo.Value != nil { entity.Value = vo.Value }
	if vo.Color != nil { entity.Color = vo.Color }
	if vo.Category != nil { entity.Category = vo.Category }
	if vo.ParentID != nil { entity.ParentID = vo.ParentID }
	if userID != "" { entity.CreatedBy = &userID; entity.UpdatedBy = &userID }
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil { panic(exception.NewBusinessError("添加字典失败: "+err.Error(), 500)) }
}

func Modify(c *gin.Context, vo *DictVO, userID string) {
	ctx := c.Request.Context()
	var entity SysDict
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound { panic(exception.NewBusinessError("数据不存在", 400)) }
		panic(exception.NewBusinessError("查询字典失败: "+err.Error(), 500))
	}
	dictCheckDuplicate(ctx, vo, vo.ID)
	if vo.ParentID != nil && *vo.ParentID != "" && *vo.ParentID != getParentIDKey(entity.ParentID) {
		dictCheckCircularParent(c.Request.Context(), vo.ID, safeStrPtr(vo.ParentID))
	}
	up := map[string]interface{}{"code": vo.Code, "sort_code": vo.SortCode, "updated_at": time.Now()}
	if vo.Label != nil { up["label"] = *vo.Label }
	if vo.Value != nil { up["value"] = *vo.Value }
	if vo.Color != nil { up["color"] = *vo.Color }
	if vo.Category != nil { up["category"] = *vo.Category }
	if vo.ParentID != nil { up["parent_id"] = *vo.ParentID } else { up["parent_id"] = nil }
	if userID != "" { up["updated_by"] = userID }
	if err := db.DB.WithContext(ctx).Model(&SysDict{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil { panic(exception.NewBusinessError("编辑字典失败: "+err.Error(), 500)) }
}

func Remove(c *gin.Context, ids []string) {
	if len(ids) == 0 { return }
	ctx := c.Request.Context()
	allIDs := dictCollectDescendantIDs(c.Request.Context(), ids)
	if err := db.DB.WithContext(ctx).Where("id IN ?", allIDs).Delete(&SysDict{}).Error; err != nil { panic(exception.NewBusinessError("删除字典失败: "+err.Error(), 500)) }
}

func Detail(c *gin.Context, id string) *DictVO {
	if id == "" { return nil }
	ctx := c.Request.Context()
	var entity SysDict
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return nil }
		panic(exception.NewBusinessError("查询字典详情失败: "+err.Error(), 500))
	}
	return entToVO(&entity)
}

func Options(c *gin.Context, param *DictOptionsParam) []*DictVO {
	ctx := c.Request.Context()
	query := db.DB.WithContext(ctx).Model(&SysDict{}).Order("sort_code ASC")
	if param.Category != "" { query = query.Where("category = ?", param.Category) }
	if param.ParentID != "" { query = query.Where("id = ? OR parent_id = ?", param.ParentID, param.ParentID) }
	var records []SysDict
	query.Find(&records)
	vos := make([]*DictVO, len(records))
	for i, r := range records { vos[i] = entToVO(&r) }
	return vos
}

func safeStrPtr(s *string) string { if s == nil { return "" }; return *s }
func dictCheckDuplicate(ctx context.Context, vo *DictVO, excludeID string) {
	if vo.Value != nil && *vo.Value != "" {
		var cnt int64
		q := db.DB.WithContext(ctx).Model(&SysDict{}).Where("parent_id = ?", vo.ParentID).Where("value = ?", *vo.Value)
		if excludeID != "" { q = q.Where("id != ?", excludeID) }
		q.Count(&cnt)
		if cnt > 0 { panic(exception.NewBusinessError("同一父字典下已存在相同值"+*vo.Value, 400)) }
	}
}

func dictCheckCircularParent(ctx context.Context, entityID, newParentID string) {
	if newParentID == "" || newParentID == "0" || entityID == "" { return }
	
	var all []SysDict
	db.DB.WithContext(ctx).Find(&all)
	parentMap := make(map[string]string)
	for _, e := range all { if e.ParentID != nil { parentMap[e.ID] = *e.ParentID } }
	current := newParentID
	for current != "" {
		if current == entityID { panic(exception.NewBusinessError("父级不能选择自身或子节点", 400)) }
		current = parentMap[current]
	}
}

func dictCollectDescendantIDs(ctx context.Context, ids []string) []string {
	
	var all []SysDict
	db.DB.WithContext(ctx).Find(&all)
	childrenMap := make(map[string][]string)
	for _, r := range all { pid := getParentIDKey(r.ParentID); childrenMap[pid] = append(childrenMap[pid], r.ID) }
	allIDs := make(map[string]bool)
	for _, id := range ids { allIDs[id] = true }
	stack := make([]string, len(ids))
	copy(stack, ids)
	for len(stack) > 0 {
		parentID := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, childID := range childrenMap[parentID] {
			if !allIDs[childID] { allIDs[childID] = true; stack = append(stack, childID) }
		}
	}
	result := make([]string, 0, len(allIDs))
	for id := range allIDs { result = append(result, id) }
	return result
}

func sortTreeNodes(nodes []map[string]interface{}) {
	sort.Slice(nodes, func(i, j int) bool { si, _ := nodes[i]["sort_code"].(int); sj, _ := nodes[j]["sort_code"].(int); return si < sj })
	for _, n := range nodes { if children, ok := n["children"].([]map[string]interface{}); ok { sortTreeNodes(children) } }
}

func entToVO(entity *SysDict) *DictVO {
	vo := &DictVO{ID: entity.ID, Code: entity.Code, Status: entity.Status, SortCode: entity.SortCode}
	if entity.Label != nil { vo.Label = entity.Label }
	if entity.Value != nil { vo.Value = entity.Value }
	if entity.Color != nil { vo.Color = entity.Color }
	if entity.Category != nil { vo.Category = entity.Category }
	if entity.ParentID != nil { vo.ParentID = entity.ParentID }
	if entity.CreatedAt != nil { vo.CreatedAt = entity.CreatedAt.Format("2006-01-02 15:04:05") }
	if entity.CreatedBy != nil { vo.CreatedBy = entity.CreatedBy }
	if entity.UpdatedAt != nil { vo.UpdatedAt = entity.UpdatedAt.Format("2006-01-02 15:04:05") }
	if entity.UpdatedBy != nil { vo.UpdatedBy = entity.UpdatedBy }
	return vo
}

func entityToNode(e *SysDict) map[string]interface{} {
	node := map[string]interface{}{"id": e.ID, "code": e.Code, "status": e.Status, "sort_code": e.SortCode}
	if e.Label != nil { node["label"] = *e.Label }
	if e.Value != nil { node["value"] = *e.Value }
	if e.Color != nil { node["color"] = *e.Color }
	if e.Category != nil { node["category"] = *e.Category }
	if e.ParentID != nil { node["parent_id"] = *e.ParentID }
	return node
}

func getParentIDKey(parentID *string) string {
	if parentID == nil || *parentID == "" { return "" }
	return *parentID
}

func DictList(c *gin.Context, param *DictListParam) []*DictVO {
	ctx := c.Request.Context()
	query := db.DB.WithContext(ctx).Model(&SysDict{}).Order("sort_code ASC")
	if param.Category != "" { query = query.Where("category = ?", param.Category) }
	if param.Keyword != "" {
		kw := "%" + param.Keyword + "%"
		query = query.Where("label LIKE ? OR code LIKE ?", kw, kw)
	}
	var records []SysDict
	query.Find(&records)
	vos := make([]*DictVO, len(records))
	for i, r := range records { vos[i] = entToVO(&r) }
	return vos
}

func DictGetLabel(c *gin.Context, typeCode, value string) gin.H {
	ctx := c.Request.Context()
	var entity SysDict
	if err := db.DB.WithContext(ctx).Where("parent_id IN (SELECT id FROM sys_dict WHERE code = ?)", typeCode).Where("value = ?", value).First(&entity).Error; err != nil {
		return gin.H{"code": 200, "message": "璇锋眰鎴愬姛", "success": true, "data": nil}
	}
	label := ""
	if entity.Label != nil { label = *entity.Label }
	return gin.H{"code": 200, "message": "璇锋眰鎴愬姛", "success": true, "data": label}
}

func DictGetChildren(c *gin.Context, typeCode string) []*DictVO {
	ctx := c.Request.Context()
	var parent SysDict
	if err := db.DB.WithContext(ctx).Where("code = ?", typeCode).First(&parent).Error; err != nil {
		return make([]*DictVO, 0)
	}
	var records []SysDict
	db.DB.WithContext(ctx).Where("parent_id = ?", parent.ID).Order("sort_code ASC").Find(&records)
	vos := make([]*DictVO, len(records))
	for i, r := range records { vos[i] = entToVO(&r) }
	return vos
}

