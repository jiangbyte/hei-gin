package dict

import (
	"context"
	"sort"

	"gorm.io/gorm"

	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

// ===== Page =====

func DictPage(c *gin.Context, p *DictPageParam) {
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

	rows, total := Page(ctx, p)

	vos := make([]*DictVO, len(rows))
	for i, r := range rows {
		vos[i] = SysDictToDictVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

// ===== Tree =====

func DictTree(c *gin.Context, param *DictTreeParam) []map[string]interface{} {
	ctx := c.Request.Context()
	all := ListForTree(ctx, param.Category, param.DictGroup)
	if len(all) == 0 {
		return make([]map[string]interface{}, 0)
	}

	childrenMap := make(map[string][]SysDict)
	for _, r := range all {
		pid := getParentIDKey(r.ParentID)
		childrenMap[pid] = append(childrenMap[pid], r)
	}
	roots := childrenMap[""]
	result := make([]map[string]interface{}, 0, len(roots))
	for _, r := range roots {
		node := entityToNode(&r)
		node["children"] = buildTreeChildren(childrenMap, r.ID, 0)
		result = append(result, node)
	}
	sortTreeNodes(result)
	return result
}

func buildTreeChildren(childrenMap map[string][]SysDict, parentID string, depth int) []map[string]interface{} {
	if depth > 50 {
		return nil
	}
	children := childrenMap[parentID]
	if len(children) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(children))
	for _, r := range children {
		node := entityToNode(&r)
		node["children"] = buildTreeChildren(childrenMap, r.ID, depth+1)
		result = append(result, node)
	}
	sortTreeNodes(result)
	return result
}

// ===== Create =====

func DictCreate(c *gin.Context, vo *DictVO) {
	ctx := c.Request.Context()
	if err := dictCheckDuplicate(ctx, vo, ""); err != nil {
		result.WriteError(c, err)
		return
	}
	if err := dictCheckCircularParent(ctx, "", utils.SafeStrPtr(vo.ParentID)); err != nil {
		result.WriteError(c, err)
		return
	}

	e := DictVOToSysDict(vo)
	e.Status = string(enums.StatusEnabled)
	if err := Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加字典失败: "+err.Error(), 500))
		return
	}
}

// ===== Modify =====

func DictModify(c *gin.Context, vo *DictVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	e, err := FindByID(ctx, vo.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询字典失败: "+err.Error(), 500))
		return
	}

	if err := dictCheckDuplicate(ctx, vo, vo.ID); err != nil {
		result.WriteError(c, err)
		return
	}
	if vo.ParentID != nil && *vo.ParentID != "" && *vo.ParentID != getParentIDKey(e.ParentID) {
		if err := dictCheckCircularParent(ctx, vo.ID, utils.SafeStrPtr(vo.ParentID)); err != nil {
			result.WriteError(c, err)
			return
		}
	}

	up := map[string]interface{}{
		"code": vo.Code, "sort_code": vo.SortCode,
	}
	if vo.Label != nil {
		up["label"] = *vo.Label
	}
	if vo.Value != nil {
		up["value"] = *vo.Value
	}
	if vo.Color != nil {
		up["color"] = *vo.Color
	}
	if vo.Category != nil {
		up["category"] = *vo.Category
	}
	if vo.ParentID != nil {
		up["parent_id"] = *vo.ParentID
	} else {
		up["parent_id"] = nil
	}
	if err := UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑字典失败: "+err.Error(), 500))
		return
	}
}

// ===== Remove =====

func DictRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	allIDs := dictCollectDescendantIDs(ctx, ids)
	if err := DeleteByIDs(ctx, allIDs); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除字典失败: "+err.Error(), 500))
		return
	}
}

// ===== Detail =====

func DictDetail(c *gin.Context, id string) *DictVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询字典详情失败: "+err.Error(), 500))
		return nil
	}
	return SysDictToDictVO(e)
}

// ===== Options =====

func DictOptions(c *gin.Context, param *DictOptionsParam) []*DictVO {
	ctx := c.Request.Context()
	records := ListOptions(ctx, param.Category, param.ParentID)
	vos := make([]*DictVO, len(records))
	for i, r := range records {
		vos[i] = SysDictToDictVO(&r)
	}
	return vos
}

// ===== List =====

func DictList(c *gin.Context, param *DictListParam) []*DictVO {
	ctx := c.Request.Context()
	records := ListByCategoryAndKeyword(ctx, param.Category, param.Keyword)
	vos := make([]*DictVO, len(records))
	for i, r := range records {
		vos[i] = SysDictToDictVO(&r)
	}
	return vos
}

// ===== GetLabel =====

func DictGetLabel(c *gin.Context, typeCode, value string) *string {
	ctx := c.Request.Context()
	entity, err := FindByTypeCodeAndValue(ctx, typeCode, value)
	if err != nil {
		return nil
	}
	return entity.Label
}

// ===== GetChildren =====

func DictGetChildren(c *gin.Context, typeCode string) []*DictVO {
	ctx := c.Request.Context()
	parent, err := FindByCode(ctx, typeCode)
	if err != nil {
		return make([]*DictVO, 0)
	}
	records := ListChildren(ctx, parent.ID)
	vos := make([]*DictVO, len(records))
	for i, r := range records {
		vos[i] = SysDictToDictVO(&r)
	}
	return vos
}

// ===== Internal helpers =====


func dictCheckDuplicate(ctx context.Context, vo *DictVO, excludeID string) error {
	if vo.Value != nil && *vo.Value != "" {
		cnt := CountDuplicateValue(ctx, vo.ParentID, *vo.Value, excludeID)
		if cnt > 0 {
			return exception.NewBusinessError("同一父字典下已存在相同值"+*vo.Value, 400)
		}
	}
	return nil
}

func dictCheckCircularParent(ctx context.Context, entityID, newParentID string) error {
	if newParentID == "" || newParentID == "0" || entityID == "" {
		return nil
	}

	all := ListAll(ctx)
	parentMap := make(map[string]string)
	for _, e := range all {
		if e.ParentID != nil {
			parentMap[e.ID] = *e.ParentID
		}
	}
	current := newParentID
	for current != "" {
		if current == entityID {
			return exception.NewBusinessError("父级不能选择自身或子节点", 400)
		}
		current = parentMap[current]
	}
	return nil
}

func dictCollectDescendantIDs(ctx context.Context, ids []string) []string {
	all := ListAll(ctx)
	childrenMap := make(map[string][]string)
	for _, r := range all {
		pid := getParentIDKey(r.ParentID)
		childrenMap[pid] = append(childrenMap[pid], r.ID)
	}
	allIDs := make(map[string]bool)
	for _, id := range ids {
		allIDs[id] = true
	}
	stack := make([]string, len(ids))
	copy(stack, ids)
	for len(stack) > 0 {
		parentID := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, childID := range childrenMap[parentID] {
			if !allIDs[childID] {
				allIDs[childID] = true
				stack = append(stack, childID)
			}
		}
	}
	result := make([]string, 0, len(allIDs))
	for id := range allIDs {
		result = append(result, id)
	}
	return result
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

func entityToNode(e *SysDict) map[string]interface{} {
	node := map[string]interface{}{
		"id": e.ID, "code": e.Code, "status": e.Status, "sort_code": e.SortCode,
	}
	if e.Label != nil {
		node["label"] = *e.Label
	}
	if e.Value != nil {
		node["value"] = *e.Value
	}
	if e.Color != nil {
		node["color"] = *e.Color
	}
	if e.Category != nil {
		node["category"] = *e.Category
	}
	if e.ParentID != nil {
		node["parent_id"] = *e.ParentID
	}
	return node
}

func getParentIDKey(parentID *string) string {
	if parentID == nil || *parentID == "" {
		return ""
	}
	return *parentID
}
