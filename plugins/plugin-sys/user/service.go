package user

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"hei-gin/sdk/constants"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type cacheOrg struct {
	ID       string
	Name     string
	ParentID *string
}
type cacheGroup struct {
	ID       string
	Name     string
	ParentID *string
}

type nameEntry struct {
	Name     string
	ParentID *string
}

var (
	cachedOrgs   []cacheOrg
	cachedGroups []cacheGroup
	cacheMu      sync.RWMutex
	cacheExpiry  time.Time
)

func batchRoleIDs(uids []string) map[string][]string {
	if len(uids) == 0 {
		return nil
	}
	ctx := context.Background()
	var rr []RelUserRole
	db.DB.WithContext(ctx).Where("user_id IN ?", uids).Find(&rr)
	m := make(map[string][]string)
	for _, r := range rr {
		m[r.UserID] = append(m[r.UserID], r.RoleID)
	}
	return m
}

func loadNameCache(ctx context.Context) {
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
	db.DB.WithContext(ctx).Table("sys_org").Select("id,name,parent_id").Find(&cachedOrgs)
	db.DB.WithContext(ctx).Table("sys_group").Select("id,name,parent_id").Find(&cachedGroups)
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

func enrichNames(vos []*UserVO) {
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
		type positionRow struct{ ID, Name string }
		var ps []positionRow
		db.DB.WithContext(ctx).Table("sys_position").Where("id IN ?", pids).Find(&ps)
		for _, p := range ps {
			pn[p.ID] = p.Name
		}
	}

	loadNameCache(ctx)
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

func UserPage(c *gin.Context, p *UserPageParam) {
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
	q := db.DB.WithContext(ctx).Model(&SysUser{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like, like)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}

	var total int64
	q.Count(&total)

	var rows []SysUser
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*UserVO, len(rows))
	for i, r := range rows {
		vos[i] = SysUserToUserVO(&r)
	}
	enrichNames(vos)
	uids := make([]string, len(rows))
	for i, r := range rows {
		uids[i] = r.ID
	}
	rm := batchRoleIDs(uids)
	for _, v := range vos {
		v.RoleIDs = rm[v.ID]
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func UserCreate(c *gin.Context, v *UserVO) {
	ctx := c.Request.Context()

	if v.Username != nil {
		var cnt int64
		db.DB.WithContext(ctx).Model(&SysUser{}).Where("username = ?", *v.Username).Count(&cnt)
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

	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加用户失败: "+err.Error(), 500))
		return
	}
}

func UserDetail(c *gin.Context, id string) *UserVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysUser
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询用户详情失败: "+err.Error(), 500))
		return nil
	}
	vo := SysUserToUserVO(&e)
	enrichNames([]*UserVO{vo})
	if rm := batchRoleIDs([]string{e.ID}); rm != nil {
		vo.RoleIDs = rm[e.ID]
	}
	return vo
}

func UserModify(c *gin.Context, v *UserVO) {
	ctx := c.Request.Context()
	if v.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}
	var old SysUser
	if err := db.DB.WithContext(ctx).First(&old, "id = ?", v.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	up := map[string]interface{}{}
	if v.Username != nil {
		var cnt int64
		db.DB.WithContext(ctx).Model(&SysUser{}).Where("username = ? AND id != ?", *v.Username, v.ID).Count(&cnt)
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
	if uid := auth.GetLoginIDDefaultNull(c); uid != "" {
		up["updated_by"] = uid
	}
	if err := db.DB.WithContext(ctx).Model(&SysUser{}).Where("id = ?", v.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑用户失败: "+err.Error(), 500))
		return
	}
}

func UserRemove(c *gin.Context, p *utils.IdsParam) {
	ids := p.IDs
	if len(ids) == 0 {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("user_id IN ?", ids).Delete(&RelUserRole{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除用户角色关联失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("user_id IN ?", ids).Delete(&RelUserPermission{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除用户权限关联失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("user_id IN ?", ids).Delete(&SysQuickAction{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除快捷操作失败: "+err.Error(), 500))
		return
	}
	if err := tx.Where("id IN ?", ids).Delete(&SysUser{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除用户失败: "+err.Error(), 500))
		return
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}

func UserGrantRole(c *gin.Context, p *GrantRoleParam) {
	if p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("用户ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("user_id = ?", p.UserID).Delete(&RelUserRole{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除已有角色失败: "+err.Error(), 500))
		return
	}
	seen := make(map[string]bool)
	batch := make([]RelUserRole, 0)
	for _, id := range p.RoleIDs {
		if !seen[id] {
			seen[id] = true
			batch = append(batch, RelUserRole{UserID: p.UserID, RoleID: id})
		}
	}
	if len(batch) > 0 {
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("分配角色失败: "+err.Error(), 500))
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}

func UserGrantPermission(c *gin.Context, p *GrantUserPermissionParam) {
	if p.UserID == "" {
		result.WriteError(c, exception.NewBusinessError("用户ID不能为空", 400))
		return
	}
	ctx := c.Request.Context()
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Where("user_id = ?", p.UserID).Delete(&RelUserPermission{}).Error; err != nil {
		tx.Rollback()
		result.WriteError(c, exception.NewBusinessError("删除已有权限失败: "+err.Error(), 500))
		return
	}
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
	if len(batch) > 0 {
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("分配权限失败: "+err.Error(), 500))
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("提交事务失败: "+err.Error(), 500))
		return
	}
}

func UserUpdateStatus(c *gin.Context, p *UpdateStatusParam) {
	if len(p.IDs) == 0 {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}
	if err := db.DB.Model(&SysUser{}).Where("id IN ?", p.IDs).Updates(
		map[string]interface{}{"status": p.Status}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("更新用户状态失败: "+err.Error(), 500))
		return
	}
}

func UserOwnRoleIDs(c *gin.Context, uid string) []string {
	var rr []RelUserRole
	db.DB.Where("user_id = ?", uid).Find(&rr)
	ids := make([]string, len(rr))
	for i, r := range rr {
		ids[i] = r.RoleID
	}
	return ids
}

func UserOwnPermissionDetails(c *gin.Context, uid string) []map[string]interface{} {
	var pp []RelUserPermission
	db.DB.Where("user_id = ?", uid).Find(&pp)
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

func UserUpdateProfile(c *gin.Context, p *UpdateProfileParam) {
	uid := auth.GetLoginIDDefaultNull(c)
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
	if err := db.DB.Model(&SysUser{}).Where("id = ?", uid).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("更新个人信息失败: "+err.Error(), 500))
		return
	}
}

func UserUpdateAvatar(c *gin.Context, p *UpdateAvatarParam) {
	uid := auth.GetLoginIDDefaultNull(c)
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
	var entity SysUser
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户不存在", 404))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	if err := db.DB.WithContext(ctx).Model(&entity).Update("avatar", avatar).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("保存头像失败: "+err.Error(), 500))
		return
	}
}

func UserUpdatePassword(c *gin.Context, p *UpdatePasswordParam) {
	uid := auth.GetLoginIDDefaultNull(c)
	if uid == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	ctx := c.Request.Context()
	var e SysUser
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", uid).Error; err != nil {
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
	if err := db.DB.WithContext(ctx).Model(&SysUser{}).Where("id = ?", uid).Update("password", string(h)).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("修改密码失败: "+err.Error(), 500))
		return
	}
}

type rawResource struct {
	ID            string
	ParentID      *string
	Code          string
	Name          string
	Category      string
	Type          string
	RoutePath     *string
	ComponentPath *string
	RedirectPath  *string
	Icon          *string
	Color         *string
	IsVisible     string
	IsCache       string
	IsAffix       string
	IsBreadcrumb  string
	ExternalURL   *string
	Description   *string
	SortCode      int
	Status        string
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

func UserOwnRoles(c *gin.Context, uid string) gin.H {
	roleIDs := UserOwnRoleIDs(c, uid)
	return gin.H{"code": 200, "message": "请求成功", "success": true, "data": roleIDs}
}

func UserCurrent(c *gin.Context) *UserVO {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		return nil
	}
	return UserDetail(c, userID)
}

func UserMenus(c *gin.Context) []map[string]interface{} {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		return make([]map[string]interface{}, 0)
	}

	roleIDs := UserOwnRoleIDs(c, userID)
	isSuperAdmin := false
	if len(roleIDs) > 0 {
		var roles []struct{ Code string }
		db.DB.Table("sys_role").Where("id IN ?", roleIDs).Find(&roles)
		for _, role := range roles {
			if role.Code == constants.SUPER_ADMIN_CODE {
				isSuperAdmin = true
				break
			}
		}
	}
	if isSuperAdmin {
		var resources []rawResource
		db.DB.Table("sys_resource").Where("status = ?", string(enums.StatusEnabled)).Order("sort_code ASC").Find(&resources)
		return buildResourceTree(resources)
	}

	if len(roleIDs) == 0 {
		return make([]map[string]interface{}, 0)
	}

	var rr []RelRoleResource
	db.DB.Where("role_id IN ?", roleIDs).Find(&rr)
	if len(rr) == 0 {
		return make([]map[string]interface{}, 0)
	}

	resourceIDs := make([]string, len(rr))
	for i, r := range rr {
		resourceIDs[i] = r.ResourceID
	}

	var resources []rawResource
	db.DB.Table("sys_resource").Where("id IN ? AND status = ?", resourceIDs, string(enums.StatusEnabled)).Order("sort_code ASC").Find(&resources)
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

func UserPermissions(c *gin.Context) []string {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		return make([]string, 0)
	}

	roleIDs := UserOwnRoleIDs(c, userID)
	var permCodes []string

	if len(roleIDs) > 0 {
		var rp []RelRolePermission
		db.DB.Where("role_id IN ?", roleIDs).Select("DISTINCT permission_code").Find(&rp)
		for _, p := range rp {
			permCodes = append(permCodes, p.PermissionCode)
		}
	}

	var up []RelUserPermission
	db.DB.Where("user_id = ?", userID).Select("DISTINCT permission_code").Find(&up)
	for _, p := range up {
		permCodes = append(permCodes, p.PermissionCode)
	}

	return permCodes
}
