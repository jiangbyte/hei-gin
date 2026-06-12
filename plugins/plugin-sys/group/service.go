package group

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

type service struct {
	repo *repository
}

func groupToVOMap(entity *SysGroup) map[string]interface{} {
	n := map[string]interface{}{
		"id": entity.ID, "code": entity.Code, "name": entity.Name, "category": entity.Category,
		"org_id": entity.OrgID, "status": entity.Status, "sort_code": entity.SortCode,
		"children": make([]map[string]interface{}, 0),
	}
	if entity.ParentID != nil {
		n["parent_id"] = *entity.ParentID
	}
	if entity.Description != nil {
		n["description"] = *entity.Description
	}
	return n
}

func (s *service) GroupPage(c *gin.Context, p *GroupPageParam) {
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

	vos := make([]*GroupVO, len(rows))
	for i, r := range rows {
		vos[i] = SysGroupToGroupVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *service) GroupTree(c *gin.Context, param *GroupTreeParam) []map[string]interface{} {
	ctx := c.Request.Context()
	all := s.repo.List(ctx, param.Category, param.OrgID)

	if len(all) == 0 {
		return make([]map[string]interface{}, 0)
	}

	childrenMap := make(map[string][]SysGroup)
	for _, r := range all {
		pid := ""
		if r.ParentID != nil {
			pid = *r.ParentID
		}
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
	if len(children) == 0 {
		return
	}
	childNodes := make([]map[string]interface{}, len(children))
	for i, c := range children {
		childNode := groupToVOMap(&c)
		buildTree(childNode, &c, childrenMap)
		childNodes[i] = childNode
	}
	node["children"] = childNodes
}

func (s *service) GroupDetail(c *gin.Context, id string) *GroupVO {
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
		result.WriteError(c, exception.NewBusinessError("查询群组详情失败: "+err.Error(), 500))
		return nil
	}
	return SysGroupToGroupVO(e)
}

func (s *service) GroupCreate(c *gin.Context, vo *GroupVO) {
	ctx := c.Request.Context()

	e := GroupVOToSysGroup(vo)
	if e.Status == "" {
		e.Status = string(enums.StatusEnabled)
	}
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加群组失败: "+err.Error(), 500))
		return
	}
}

func (s *service) GroupModify(c *gin.Context, vo *GroupVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	_, err := s.repo.FindByID(ctx, vo.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询群组失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"code": vo.Code, "name": vo.Name, "category": vo.Category,
		"org_id": vo.OrgID, "parent_id": vo.ParentID,
		"status": vo.Status, "sort_code": vo.SortCode,
	}
	if vo.Description != nil {
		up["description"] = *vo.Description
	}
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑群组失败: "+err.Error(), 500))
		return
	}
}

func (s *service) GroupRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	allIDs := s.getAllDescendantIDs(ctx, ids)
	s.repo.ClearUserGroupIDs(ctx, allIDs)
	_ = s.repo.DeleteByIDs(ctx, allIDs)
}

func (s *service) GroupOptions(c *gin.Context) []*GroupVO {
	ctx := c.Request.Context()
	records := s.repo.ListAll(ctx)
	vos := make([]*GroupVO, len(records))
	for i, r := range records {
		vos[i] = SysGroupToGroupVO(&r)
	}
	return vos
}

func (s *service) GroupGetAll(c *gin.Context) []*GroupVO {
	ctx := c.Request.Context()
	records := s.repo.ListAll(ctx)
	vos := make([]*GroupVO, len(records))
	for i, r := range records {
		vos[i] = SysGroupToGroupVO(&r)
	}
	return vos
}

func (s *service) getAllDescendantIDs(ctx context.Context, ids []string) []string {
	allIDs := make(map[string]bool)
	for _, id := range ids {
		allIDs[id] = true
	}

	all := s.repo.ListAll(ctx)
	cm := make(map[string][]string)
	for _, g := range all {
		if g.ParentID != nil && *g.ParentID != "" {
			cm[*g.ParentID] = append(cm[*g.ParentID], g.ID)
		}
	}

	q := make([]string, len(ids))
	copy(q, ids)
	for len(q) > 0 {
		pid := q[len(q)-1]
		q = q[:len(q)-1]
		for _, cid := range cm[pid] {
			if !allIDs[cid] {
				allIDs[cid] = true
				q = append(q, cid)
			}
		}
	}
	r := make([]string, 0, len(allIDs))
	for id := range allIDs {
		r = append(r, id)
	}
	return r
}

func GroupPage(c *gin.Context, p *GroupPageParam) { defaultModule.service.GroupPage(c, p) }
func GroupTree(c *gin.Context, param *GroupTreeParam) []map[string]interface{} {
	return defaultModule.service.GroupTree(c, param)
}
func GroupDetail(c *gin.Context, id string) *GroupVO { return defaultModule.service.GroupDetail(c, id) }
func GroupCreate(c *gin.Context, vo *GroupVO)        { defaultModule.service.GroupCreate(c, vo) }
func GroupModify(c *gin.Context, vo *GroupVO)        { defaultModule.service.GroupModify(c, vo) }
func GroupRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.GroupRemove(c, param)
}
func GroupOptions(c *gin.Context) []*GroupVO { return defaultModule.service.GroupOptions(c) }
func GroupGetAll(c *gin.Context) []*GroupVO  { return defaultModule.service.GroupGetAll(c) }
