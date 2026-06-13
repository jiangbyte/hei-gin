package user

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"hei-gin/sdk/constants"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type nameEntry struct {
	Name     string
	ParentID *string
}

type Service struct {
	repo *repository
}

var (
	cachedOrgs   []cacheOrg
	cachedGroups []cacheGroup
	cacheMu      sync.RWMutex
	cacheExpiry  time.Time
)

func (s *Service) batchRoleIDs(uids []string) map[string][]string {
	if len(uids) == 0 {
		return nil
	}
	ctx := context.Background()
	rr := s.repo.FindRoleRelsByUserIDs(ctx, uids)
	m := make(map[string][]string)
	for _, r := range rr {
		m[r.UserID] = append(m[r.UserID], r.RoleID)
	}
	return m
}

func (s *Service) loadNameCache(ctx context.Context) {
	cacheMu.RLock()
	if time.Since(cacheExpiry) <= 5*time.Minute {
		cacheMu.RUnlock()
		return
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()
	if time.Since(cacheExpiry) <= 5*time.Minute {
		return
	}
	cachedOrgs = s.repo.FindOrgCacheRows(ctx)
	cachedGroups = s.repo.FindGroupCacheRows(ctx)
	cacheExpiry = time.Now()
}

func buildNameMaps() (map[string]nameEntry, map[string]nameEntry) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	orgMap := make(map[string]nameEntry, len(cachedOrgs))
	for _, o := range cachedOrgs {
		orgMap[o.ID] = nameEntry{o.Name, o.ParentID}
	}
	grpMap := make(map[string]nameEntry, len(cachedGroups))
	for _, g := range cachedGroups {
		grpMap[g.ID] = nameEntry{g.Name, g.ParentID}
	}
	return orgMap, grpMap
}

func resolvePath(id *string, m map[string]nameEntry) []string {
	if id == nil || *id == "" {
		return nil
	}
	var path []string
	cur := *id
	for {
		n, ok := m[cur]
		if !ok {
			break
		}
		path = append(path, n.Name)
		if n.ParentID == nil || *n.ParentID == "" {
			break
		}
		cur = *n.ParentID
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func (s *Service) enrichNames(vos []*UserVO) {
	if len(vos) == 0 {
		return
	}
	ctx := context.Background()

	var pids []string
	for _, v := range vos {
		if v.PositionID != nil && *v.PositionID != "" {
			pids = append(pids, *v.PositionID)
		}
	}
	pn := make(map[string]string)
	if len(pids) > 0 {
		ps := s.repo.FindPositionRows(ctx, pids)
		for _, p := range ps {
			pn[p.ID] = p.Name
		}
	}

	s.loadNameCache(ctx)
	orgMap, grpMap := buildNameMaps()

	for _, v := range vos {
		if v.PositionID != nil {
			if n, ok := pn[*v.PositionID]; ok {
				v.PositionName = &n
			}
		}
		if v.OrgID != nil {
			v.OrgNames = resolvePath(v.OrgID, orgMap)
		}
		if v.GroupID != nil {
			v.GroupNames = resolvePath(v.GroupID, grpMap)
		}
	}
}

func (s *Service) Page(c *gin.Context, p *UserPageParam) {
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

	vos := make([]*UserVO, len(rows))
	for i, r := range rows {
		vos[i] = SysUserToUserVO(&r)
	}
	s.enrichNames(vos)
	uids := make([]string, len(rows))
	for i, r := range rows {
		uids[i] = r.ID
	}
	rm := s.batchRoleIDs(uids)
	for _, v := range vos {
		v.RoleIDs = rm[v.ID]
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) Create(c *gin.Context, v *UserVO) {
	ctx := c.Request.Context()

	if v.Username != nil {
		cnt := s.repo.CountByUsername(ctx, *v.Username, "")
		if cnt > 0 {
			result.WriteError(c, exception.NewBusinessError("账号已存在", 400))
			return
		}
	}

	e := UserVOToSysUser(v)
	e.Status = string(enums.UserStatusActive)

	if v.Password != nil {
		h, _ := bcrypt.GenerateFromPassword([]byte(*v.Password), bcrypt.DefaultCost)
		s := string(h)
		e.Password = &s
	}

	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加用户失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Detail(c *gin.Context, id string) *UserVO {
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
		result.WriteError(c, exception.NewBusinessError("查询用户详情失败: "+err.Error(), 500))
		return nil
	}
	vo := SysUserToUserVO(e)
	s.enrichNames([]*UserVO{vo})
	if rm := s.batchRoleIDs([]string{e.ID}); rm != nil {
		vo.RoleIDs = rm[e.ID]
	}
	return vo
}

func (s *Service) Modify(c *gin.Context, v *UserVO) {
	ctx := c.Request.Context()
	if v.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}
	old, err := s.repo.FindByID(ctx, v.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	up := map[string]interface{}{}
	if v.Username != nil {
		cnt := s.repo.CountByUsername(ctx, *v.Username, v.ID)
		if cnt > 0 {
			result.WriteError(c, exception.NewBusinessError("账号已存在", 400))
			return
		}
		up["username"] = *v.Username
	}
	if v.Nickname != nil {
		up["nickname"] = *v.Nickname
	}
	if v.Avatar != nil {
		up["avatar"] = *v.Avatar
	}
	if v.Motto != nil {
		up["motto"] = *v.Motto
	}
	if v.Gender != nil {
		up["gender"] = *v.Gender
	}
	if v.Birthday != "" {
		up["birthday"] = utils.ParseDateTimePtr(&v.Birthday)
	}
	if v.Email != nil {
		up["email"] = *v.Email
	}
	if v.Github != nil {
		up["github"] = *v.Github
	}
	if v.Phone != nil {
		up["phone"] = *v.Phone
	}
	if v.OrgID != nil {
		up["org_id"] = *v.OrgID
	} else if old.OrgID != nil {
		up["org_id"] = nil
	}
	if v.PositionID != nil {
		up["position_id"] = *v.PositionID
	} else if old.PositionID != nil {
		up["position_id"] = nil
	}
	if v.GroupID != nil {
		up["group_id"] = *v.GroupID
	} else if old.GroupID != nil {
		up["group_id"] = nil
	}
	if v.Status != "" {
		up["status"] = v.Status
	}
	up["updated_at"] = time.Now()
	if uid := auth.Business.GetLoginIDDefaultNull(c); uid != "" {
		up["updated_by"] = uid
	}
	if err := s.repo.UpdateByID(ctx, v.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑用户失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Remove(c *gin.Context, p *utils.IdsParam) {
	ids := p.IDs
	if len(ids) == 0 {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.DeleteUsersCascade(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除用户失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) GrantRole(c *gin.Context, p *GrantRoleParam) {
	if p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("用户ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()
	seen := make(map[string]bool)
	batch := make([]RelUserRole, 0)
	for _, id := range p.RoleIDs {
		if !seen[id] {
			seen[id] = true
			batch = append(batch, RelUserRole{UserID: p.UserID, RoleID: id})
		}
	}
	if err := s.repo.ReplaceUserRoles(ctx, p.UserID, batch); err != nil {
		result.WriteError(c, exception.NewBusinessError("分配角色失败: "+err.Error(), 500))
		return
	}
	if err := auth.Business.RefreshUserSessionsACL(ctx, p.UserID); err != nil {
		result.WriteError(c, exception.NewBusinessError("刷新用户会话权限失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) GrantPermission(c *gin.Context, p *GrantUserPermissionParam) {
	if p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("用户ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()
	batch := make([]RelUserPermission, len(p.Permissions))
	for i, pi := range p.Permissions {
		r := RelUserPermission{UserID: p.UserID, PermissionCode: pi.PermissionCode, Scope: pi.Scope}
		if pi.CustomScopeGroupIds != nil {
			r.CustomScopeGroupIds = pi.CustomScopeGroupIds
		}
		if pi.CustomScopeOrgIds != nil {
			r.CustomScopeOrgIds = pi.CustomScopeOrgIds
		}
		batch[i] = r
	}
	if err := s.repo.ReplaceUserPermissions(ctx, p.UserID, batch); err != nil {
		result.WriteError(c, exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
		return
	}
	if err := auth.Business.RefreshUserSessionsACL(ctx, p.UserID); err != nil {
		result.WriteError(c, exception.NewBusinessError("刷新用户会话权限失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) RefreshSessionACL(c *gin.Context, p *RefreshSessionACLParam) {
	if p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("用户ID不能为空", 400))
		return
	}
	if err := auth.Business.RefreshUserSessionsACL(c.Request.Context(), p.UserID); err != nil {
		result.WriteError(c, exception.NewBusinessError("刷新用户会话权限失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) BatchRefreshSessionACL(c *gin.Context, p *BatchRefreshSessionACLParam) {
	if len(p.UserIDs) == 0 {
		result.WriteError(c, exception.NewBusinessError("用户ID不能为空", 400))
		return
	}
	seen := make(map[string]struct{}, len(p.UserIDs))
	for _, userID := range p.UserIDs {
		if userID == "" {
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		if err := auth.Business.RefreshUserSessionsACL(c.Request.Context(), userID); err != nil {
			result.WriteError(c, exception.NewBusinessError("刷新用户会话权限失败: "+err.Error(), 500))
			return
		}
	}
}

func (s *Service) UpdateStatus(c *gin.Context, p *UpdateStatusParam) {
	if len(p.IDs) == 0 {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}
	if err := s.repo.UpdateStatusByIDs(c.Request.Context(), p.IDs, p.Status); err != nil {
		result.WriteError(c, exception.NewBusinessError("更新用户状态失败: "+err.Error(), 500))
		return
	}
	if p.Status != string(enums.UserStatusActive) {
		for _, userID := range p.IDs {
			auth.Business.Sessions().KickoutUser(c.Request.Context(), userID)
		}
	}
}

func (s *Service) OwnRoleIDs(c *gin.Context, uid string) []string {
	rr := s.repo.FindRoleRelsByUserID(c.Request.Context(), uid)
	ids := make([]string, len(rr))
	for i, r := range rr {
		ids[i] = r.RoleID
	}
	return ids
}

func (s *Service) OwnPermissionDetails(c *gin.Context, uid string) []map[string]interface{} {
	pp := s.repo.FindPermissionRelsByUserID(c.Request.Context(), uid)
	r := make([]map[string]interface{}, len(pp))
	for i, p := range pp {
		r[i] = map[string]interface{}{
			"permission_code":        p.PermissionCode,
			"scope":                  p.Scope,
			"custom_scope_group_ids": p.CustomScopeGroupIds,
			"custom_scope_org_ids":   p.CustomScopeOrgIds,
		}
	}
	return r
}

func (s *Service) UpdateProfile(c *gin.Context, p *UpdateProfileParam) {
	uid := auth.Business.GetLoginIDDefaultNull(c)
	if uid == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	up := map[string]interface{}{}
	if p.Username != nil {
		up["username"] = *p.Username
	}
	if p.Nickname != nil {
		up["nickname"] = *p.Nickname
	}
	if p.Motto != nil {
		up["motto"] = *p.Motto
	}
	if p.Gender != nil {
		up["gender"] = *p.Gender
	}
	if p.Birthday != "" {
		up["birthday"] = utils.ParseDateTimePtr(&p.Birthday)
	}
	if p.Email != nil {
		up["email"] = *p.Email
	}
	if p.Github != nil {
		up["github"] = *p.Github
	}
	if p.Phone != nil {
		up["phone"] = *p.Phone
	}
	up["updated_at"] = time.Now()
	if err := s.repo.UpdateByID(c.Request.Context(), uid, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("更新个人信息失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) UpdateAvatar(c *gin.Context, p *UpdateAvatarParam) {
	uid := auth.Business.GetLoginIDDefaultNull(c)
	avatar := p.Avatar
	if uid == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	if avatar == "" {
		result.WriteError(c, exception.NewBusinessError("头像不能为空", 400))
		return
	}
	avatar = utils.CompressBase64Image(avatar, 512, 512, 80)
	ctx := c.Request.Context()
	entity, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户不存在", 404))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	if err := s.repo.UpdateAvatar(ctx, entity, avatar); err != nil {
		result.WriteError(c, exception.NewBusinessError("保存头像失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) UpdatePassword(c *gin.Context, p *UpdatePasswordParam) {
	uid := auth.Business.GetLoginIDDefaultNull(c)
	if uid == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户不存在", 404))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	if e.Password == nil || *e.Password == "" {
		result.WriteError(c, exception.NewBusinessError("未设置密码，无法修改", 400))
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(*e.Password), []byte(utils.Decrypt(p.CurrentPassword))) != nil {
		result.WriteError(c, exception.NewBusinessError("当前密码不正确", 400))
		return
	}
	newPwd := utils.Decrypt(p.NewPassword)
	if newPwd == "" {
		result.WriteError(c, exception.NewBusinessError("新密码解密失败", 400))
		return
	}
	h, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("密码加密失败", 500))
		return
	}
	if err := s.repo.UpdatePassword(ctx, uid, string(h)); err != nil {
		result.WriteError(c, exception.NewBusinessError("修改密码失败: "+err.Error(), 500))
		return
	}
}

func buildResourceTree(resources []rawResource) []map[string]interface{} {
	cm := make(map[string][]rawResource)
	for _, r := range resources {
		pid := ""
		if r.ParentID != nil && *r.ParentID != "" {
			pid = *r.ParentID
		}
		cm[pid] = append(cm[pid], r)
	}
	return buildUserMenuTree(cm, "")
}

func (s *Service) OwnRoles(c *gin.Context, uid string) gin.H {
	roleIDs := s.OwnRoleIDs(c, uid)
	return gin.H{"code": 200, "message": "请求成功", "success": true, "data": roleIDs}
}

func (s *Service) Current(c *gin.Context) *UserVO {
	userID := auth.Business.GetLoginIDDefaultNull(c)
	if userID == "" {
		return nil
	}
	return s.Detail(c, userID)
}

func (s *Service) Menus(c *gin.Context) []map[string]interface{} {
	userID := auth.Business.GetLoginIDDefaultNull(c)
	if userID == "" {
		return make([]map[string]interface{}, 0)
	}

	roleIDs := s.OwnRoleIDs(c, userID)
	isSuperAdmin := false
	if len(roleIDs) > 0 {
		roleRows := s.repo.FindRoleCodesByIDs(c.Request.Context(), roleIDs)
		for _, role := range roleRows {
			if role.Code == constants.SUPER_ADMIN_CODE {
				isSuperAdmin = true
				break
			}
		}
	}
	if isSuperAdmin {
		resources := s.repo.FindEnabledResources(c.Request.Context())
		return buildResourceTree(resources)
	}

	if len(roleIDs) == 0 {
		return make([]map[string]interface{}, 0)
	}

	rr := s.repo.FindRoleResourcesByRoleIDs(c.Request.Context(), roleIDs)
	if len(rr) == 0 {
		return make([]map[string]interface{}, 0)
	}

	resourceIDs := make([]string, len(rr))
	for i, r := range rr {
		resourceIDs[i] = r.ResourceID
	}

	resources := s.repo.FindEnabledResourcesByIDs(c.Request.Context(), resourceIDs)
	return buildResourceTree(resources)
}

func buildUserMenuTree(cm map[string][]rawResource, pid string) []map[string]interface{} {
	cs := cm[pid]
	r := make([]map[string]interface{}, 0, len(cs))
	for _, c := range cs {
		n := map[string]interface{}{
			"id": c.ID, "code": c.Code, "name": c.Name, "category": c.Category, "type": c.Type,
			"route_path": c.RoutePath, "component_path": c.ComponentPath, "redirect_path": c.RedirectPath,
			"icon": c.Icon, "color": c.Color, "is_visible": c.IsVisible, "is_cache": c.IsCache,
			"is_affix": c.IsAffix, "is_breadcrumb": c.IsBreadcrumb, "external_url": c.ExternalURL,
			"sort_code": c.SortCode, "status": c.Status,
		}
		if c.ParentID != nil {
			n["parent_id"] = *c.ParentID
		} else {
			n["parent_id"] = nil
		}
		if c.Description != nil {
			n["description"] = *c.Description
		}
		n["children"] = buildUserMenuTree(cm, c.ID)
		r = append(r, n)
	}
	return r
}

func (s *Service) Permissions(c *gin.Context) []string {
	userID := auth.Business.GetLoginIDDefaultNull(c)
	if userID == "" {
		return make([]string, 0)
	}

	roleIDs := s.OwnRoleIDs(c, userID)
	var permCodes []string

	if len(roleIDs) > 0 {
		rp := s.repo.FindDistinctRolePermissionCodes(c.Request.Context(), roleIDs)
		for _, p := range rp {
			permCodes = append(permCodes, p.PermissionCode)
		}
	}

	up := s.repo.FindDistinctUserPermissionCodes(c.Request.Context(), userID)
	for _, p := range up {
		permCodes = append(permCodes, p.PermissionCode)
	}

	return permCodes
}
