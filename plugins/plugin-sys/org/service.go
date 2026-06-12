package org

import (
	"context"
	"sort"
	"sync"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type service struct {
	repo *repository
}

const orgTreeCacheTTL = 30 * time.Second

type orgTreeCacheEntry struct {
	expires time.Time
	data    []map[string]interface{}
}

var (
	orgTreeMu    sync.RWMutex
	orgTreeCache = make(map[string]orgTreeCacheEntry)
)

func invalidateOrgTreeCache() {
	orgTreeMu.Lock()
	orgTreeCache = make(map[string]orgTreeCacheEntry)
	orgTreeMu.Unlock()
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

func (s *service) OrgPage(c *gin.Context, p *OrgPageParam) {
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

	vos := make([]*OrgVO, len(rows))
	for i, r := range rows {
		vos[i] = SysOrgToOrgVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *service) OrgTree(c *gin.Context, p *OrgTreeParam) []map[string]interface{} {
	cacheKey := p.Category
	orgTreeMu.RLock()
	if cached, ok := orgTreeCache[cacheKey]; ok && time.Now().Before(cached.expires) {
		orgTreeMu.RUnlock()
		return cached.data
	}
	orgTreeMu.RUnlock()

	ctx := c.Request.Context()
	all, err := s.repo.List(ctx, p.Category)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("查询组织树失败: "+err.Error(), 500))
		return nil
	}

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
	orgTreeMu.Lock()
	orgTreeCache[cacheKey] = orgTreeCacheEntry{expires: time.Now().Add(orgTreeCacheTTL), data: roots}
	orgTreeMu.Unlock()
	return roots
}

func (s *service) OrgCreate(c *gin.Context, vo *OrgVO) {
	ctx := c.Request.Context()

	if vo.Code == "" || vo.Name == "" || vo.Category == "" {
		result.WriteError(c, exception.NewBusinessError("组织编码、名称、类别不能为空", 400))
		return
	}

	e := OrgVOToSysOrg(vo)
	e.Status = string(enums.StatusEnabled)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加组织失败: "+err.Error(), 500))
		return
	}
	invalidateOrgTreeCache()
}

func (s *service) OrgModify(c *gin.Context, vo *OrgVO) {
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
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑组织失败: "+err.Error(), 500))
		return
	}
	invalidateOrgTreeCache()
}

func (s *service) OrgRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()

	allIDs := s.collectDescendantOrgIDs(ctx, ids)

	userCount := s.repo.CountUsersByOrgIDs(ctx, allIDs)
	if userCount > 0 {
		result.WriteError(c, exception.NewBusinessError("组织存在关联用户，无法删除", 400))
		return
	}

	groupCount := s.repo.CountGroupsByOrgIDs(ctx, allIDs)
	if groupCount > 0 {
		result.WriteError(c, exception.NewBusinessError("组织存在关联用户组，无法删除", 400))
		return
	}

	posCount := s.repo.CountPositionsByOrgIDs(ctx, allIDs)
	if posCount > 0 {
		result.WriteError(c, exception.NewBusinessError("组织存在关联职位，无法删除", 400))
		return
	}

	if err := s.repo.DeleteByIDs(ctx, allIDs); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除组织失败: "+err.Error(), 500))
		return
	}
	invalidateOrgTreeCache()
}

func (s *service) OrgDetail(c *gin.Context, id string) *OrgVO {
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
		result.WriteError(c, exception.NewBusinessError("查询组织详情失败: "+err.Error(), 500))
		return nil
	}
	return SysOrgToOrgVO(e)
}

func (s *service) OrgOptions(c *gin.Context) []*OrgVO {
	ctx := c.Request.Context()
	rows, err := s.repo.ListAll(ctx)
	if err != nil {
		return make([]*OrgVO, 0)
	}
	vos := make([]*OrgVO, len(rows))
	for i, r := range rows {
		vos[i] = SysOrgToOrgVO(&r)
	}
	return vos
}

func (s *service) collectDescendantOrgIDs(ctx context.Context, ids []string) []string {
	allIDs := make(map[string]bool)
	for _, id := range ids {
		allIDs[id] = true
	}

	all, err := s.repo.ListAll(ctx)
	if err != nil {
		return ids
	}
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

func OrgPage(c *gin.Context, p *OrgPageParam) { defaultModule.service.OrgPage(c, p) }
func OrgTree(c *gin.Context, p *OrgTreeParam) []map[string]interface{} {
	return defaultModule.service.OrgTree(c, p)
}
func OrgCreate(c *gin.Context, vo *OrgVO) { defaultModule.service.OrgCreate(c, vo) }
func OrgModify(c *gin.Context, vo *OrgVO) { defaultModule.service.OrgModify(c, vo) }
func OrgRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.OrgRemove(c, param)
}
func OrgDetail(c *gin.Context, id string) *OrgVO { return defaultModule.service.OrgDetail(c, id) }
func OrgOptions(c *gin.Context) []*OrgVO         { return defaultModule.service.OrgOptions(c) }
