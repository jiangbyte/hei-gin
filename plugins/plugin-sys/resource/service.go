package resource

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"
)

type service struct {
	repo *repository
}

const resourceTreeCacheTTL = 30 * time.Second

type resourceTreeCacheEntry struct {
	expires time.Time
	data    []map[string]interface{}
}

var (
	resourceTreeMu    sync.RWMutex
	resourceTreeCache = make(map[string]resourceTreeCacheEntry)
)

func invalidateResourceTreeCache() {
	resourceTreeMu.Lock()
	resourceTreeCache = make(map[string]resourceTreeCacheEntry)
	resourceTreeMu.Unlock()
}

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

func (s *service) ModulePage(c *gin.Context, param *ModulePageParam) {
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

	records, total := s.repo.ModulePageQuery(ctx, param.Current, param.Size)

	vos := make([]*ModuleVO, len(records))
	for i, r := range records {
		vos[i] = SysModuleToModuleVO(&r)
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func (s *service) ModuleDetail(c *gin.Context, id string) *ModuleVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询模块详情失败: "+err.Error(), 500))
		return nil
	}
	return SysModuleToModuleVO(e)
}

func (s *service) ModuleCreate(c *gin.Context, vo *ModuleVO) {
	ctx := c.Request.Context()

	e := ModuleVOToSysModule(vo)
	if e.Status == "" {
		e.Status = string(enums.StatusEnabled)
	}
	if err := s.repo.CreateModule(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加模块失败: "+err.Error(), 500))
		return
	}
	invalidateResourceTreeCache()
}

func (s *service) ModuleModify(c *gin.Context, vo *ModuleVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	if _, err := s.repo.FindModuleByID(ctx, vo.ID); err != nil {
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
	if err := s.repo.UpdateModuleByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑模块失败: "+err.Error(), 500))
		return
	}
	invalidateResourceTreeCache()
}

func (s *service) ModuleRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.DeleteModules(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除模块失败: "+err.Error(), 500))
		return
	}
	invalidateResourceTreeCache()
}

func (s *service) ResourcePage(c *gin.Context, param *ResourcePageParam) {
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

	records, total := s.repo.ResourcePageQuery(ctx, param.Current, param.Size)

	vos := make([]*ResourceVO, len(records))
	for i, r := range records {
		vos[i] = SysResourceToResourceVO(&r)
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func (s *service) ResourceTree(c *gin.Context, category string) []map[string]interface{} {
	cacheKey := "tree:" + category
	resourceTreeMu.RLock()
	if cached, ok := resourceTreeCache[cacheKey]; ok && time.Now().Before(cached.expires) {
		resourceTreeMu.RUnlock()
		return cached.data
	}
	resourceTreeMu.RUnlock()

	ctx := c.Request.Context()
	all, err := s.repo.ListResources(ctx, category)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("查询资源树失败: "+err.Error(), 500))
		return nil
	}

	cm := make(map[string][]SysResource)
	for _, r := range all {
		pid := ""
		if r.ParentID != nil && *r.ParentID != "" && *r.ParentID != "0" {
			pid = *r.ParentID
		}
		cm[pid] = append(cm[pid], r)
	}
	data := buildRT(cm, "", 0)
	resourceTreeMu.Lock()
	resourceTreeCache[cacheKey] = resourceTreeCacheEntry{expires: time.Now().Add(resourceTreeCacheTTL), data: data}
	resourceTreeMu.Unlock()
	return data
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

func (s *service) ResourceMenu(c *gin.Context) []map[string]interface{} {
	cacheKey := "menu"
	resourceTreeMu.RLock()
	if cached, ok := resourceTreeCache[cacheKey]; ok && time.Now().Before(cached.expires) {
		resourceTreeMu.RUnlock()
		return cached.data
	}
	resourceTreeMu.RUnlock()

	ctx := c.Request.Context()
	all, err := s.repo.ListAllResources(ctx)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("查询资源菜单失败: "+err.Error(), 500))
		return nil
	}
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
	resourceTreeMu.Lock()
	resourceTreeCache[cacheKey] = resourceTreeCacheEntry{expires: time.Now().Add(resourceTreeCacheTTL), data: r}
	resourceTreeMu.Unlock()
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

func (s *service) ResourceDetail(c *gin.Context, id string) *ResourceVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询资源详情失败: "+err.Error(), 500))
		return nil
	}
	return SysResourceToResourceVO(e)
}

func (s *service) ResourceCreate(c *gin.Context, vo *ResourceVO) {
	ctx := c.Request.Context()

	e := ResourceVOToSysResource(vo)
	if e.Status == "" {
		e.Status = string(enums.StatusEnabled)
	}
	if err := s.repo.CreateResource(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加资源失败: "+err.Error(), 500))
		return
	}
	invalidateResourceTreeCache()
}

func (s *service) ResourceModify(c *gin.Context, vo *ResourceVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	old, err := s.repo.FindResourceByID(ctx, vo.ID)
	if err != nil {
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

	if err := s.repo.UpdateResourceByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑资源失败: "+err.Error(), 500))
		return
	}
	invalidateResourceTreeCache()

	if err := s.repo.SyncPermissions(ctx, vo.ID, extractPermCode(oldExtra), extractPermCode(vo.Extra)); err != nil {
		log.Printf("[RESOURCE] Failed to sync permissions: %v", err)
	}
}

func (s *service) ResourceRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	all := s.repo.CollectResourceDescendants(ctx, ids)
	if err := s.repo.DeleteResourcesCascade(ctx, all); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除资源失败: "+err.Error(), 500))
		return
	}
	invalidateResourceTreeCache()
}

func (s *service) collectDescendant(ctx context.Context, ids []string) []string {
	return s.repo.CollectResourceDescendants(ctx, ids)
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

func ModulePage(c *gin.Context, param *ModulePageParam) { defaultModule.service.ModulePage(c, param) }
func ModuleDetail(c *gin.Context, id string) *ModuleVO {
	return defaultModule.service.ModuleDetail(c, id)
}
func ModuleCreate(c *gin.Context, vo *ModuleVO) { defaultModule.service.ModuleCreate(c, vo) }
func ModuleModify(c *gin.Context, vo *ModuleVO) { defaultModule.service.ModuleModify(c, vo) }
func ModuleRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.ModuleRemove(c, param)
}
func ResourcePage(c *gin.Context, param *ResourcePageParam) {
	defaultModule.service.ResourcePage(c, param)
}
func ResourceTree(c *gin.Context, category string) []map[string]interface{} {
	return defaultModule.service.ResourceTree(c, category)
}
func ResourceMenu(c *gin.Context) []map[string]interface{} {
	return defaultModule.service.ResourceMenu(c)
}
func ResourceDetail(c *gin.Context, id string) *ResourceVO {
	return defaultModule.service.ResourceDetail(c, id)
}
func ResourceCreate(c *gin.Context, vo *ResourceVO) { defaultModule.service.ResourceCreate(c, vo) }
func ResourceModify(c *gin.Context, vo *ResourceVO) { defaultModule.service.ResourceModify(c, vo) }
func ResourceRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.ResourceRemove(c, param)
}
