package home

import (
	"context"
	"time"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	resModel "hei-gin/plugins/plugin-sys/resource"

	"github.com/gin-gonic/gin"
)

func HomeGet(c *gin.Context) *HomeVO {
	userID := auth.GetLoginIDDefaultNull(c)
	vo := &HomeVO{
		QuickActions:       make([]QuickActionVO, 0),
		AvailableResources: make([]QuickActionVO, 0),
		Notices:            make([]HomeNotice, 0),
	}
	if userID != "" {
		vo.QuickActions = findQuickActionsByUserID(c.Request.Context(), userID)
		vo.AvailableResources = getAvailableResources(c.Request.Context(), userID)
	}
	vo.Notices = getNotices(c.Request.Context())
	vo.Stats.TotalUsers = getUserCount(c.Request.Context())
	return vo
}

func HomeAddQuickAction(c *gin.Context, param *AddQuickActionParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()

	var count int64
	db.DB.WithContext(ctx).Model(&SysQuickAction{}).Where("user_id = ? AND resource_id = ?", userID, param.ResourceID).Count(&count)
	if count > 0 {
		return
	}

	var actionCount int64
	db.DB.WithContext(ctx).Model(&SysQuickAction{}).Where("user_id = ?", userID).Count(&actionCount)
	entity := SysQuickAction{
		UserID: userID, ResourceID: param.ResourceID,
		SortCode: int(actionCount+1) * 10,
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加快捷方式失败: "+err.Error(), 500))
		return
	}
}

func HomeRemoveQuickAction(c *gin.Context, param *utils.IdsParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id IN ?", param.IDs).Delete(&SysQuickAction{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("移除快捷方式失败: "+err.Error(), 500))
		return
	}
}

func HomeSortQuickActions(c *gin.Context, param *utils.IdsParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()
	tx := db.DB.WithContext(ctx).Begin()
	for idx, id := range param.IDs {
		if err := tx.Model(&SysQuickAction{}).Where("id = ?", id).Update("sort_code", (idx+1)*10).Error; err != nil {
			tx.Rollback()
			result.WriteError(c, exception.NewBusinessError("排序快捷方式失败: "+err.Error(), 500))
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("排序快捷方式提交事务失败: "+err.Error(), 500))
		return
	}
}

func findQuickActionsByUserID(ctx context.Context, userID string) []QuickActionVO {
	var actions []SysQuickAction
	db.DB.WithContext(ctx).Where("user_id = ?", userID).Order("sort_code ASC, created_at ASC").Find(&actions)
	if len(actions) == 0 {
		return make([]QuickActionVO, 0)
	}

	resourceIDs := make([]string, len(actions))
	for i, a := range actions {
		resourceIDs[i] = a.ResourceID
	}

	var resources []resModel.SysResource
	db.DB.WithContext(ctx).Where("id IN ?", resourceIDs).Find(&resources)
	resourceMap := make(map[string]resModel.SysResource)
	for _, r := range resources {
		resourceMap[r.ID] = r
	}

	vos := make([]QuickActionVO, 0, len(actions))
	for _, a := range actions {
		vo := SysQuickActionToQuickActionVO(&a)
		if r, ok := resourceMap[a.ResourceID]; ok {
			vo.Name = r.Name
			vo.Type = r.Type
			if r.Icon != nil {
				vo.Icon = *r.Icon
			}
			if r.RoutePath != nil {
				vo.RoutePath = *r.RoutePath
			}
			if r.ParentID != nil {
				vo.ParentID = *r.ParentID
			}
		}
		vos = append(vos, *vo)
	}
	return vos
}

func getAvailableResources(ctx context.Context, userID string) []QuickActionVO {
	var actionIDs []string
	db.DB.WithContext(ctx).Model(&SysQuickAction{}).Where("user_id = ?", userID).Select("resource_id").Find(&actionIDs)

	q := db.DB.WithContext(ctx).Model(&resModel.SysResource{}).Where("category IN ? AND status = ?", []string{string(enums.ResourceCategoryBackendMenu), string(enums.ResourceCategoryFrontendMenu)}, string(enums.StatusEnabled))
	if len(actionIDs) > 0 {
		q = q.Where("id NOT IN ?", actionIDs)
	}

	var resources []resModel.SysResource
	q.Order("sort_code ASC").Find(&resources)

	vos := make([]QuickActionVO, len(resources))
	for i, r := range resources {
		vos[i] = QuickActionVO{
			ResourceID: r.ID, Name: r.Name, Type: r.Type,
		}
		if r.Icon != nil {
			vos[i].Icon = *r.Icon
		}
		if r.RoutePath != nil {
			vos[i].RoutePath = *r.RoutePath
		}
		if r.ParentID != nil {
			vos[i].ParentID = *r.ParentID
		}
	}
	return vos
}

func getNotices(ctx context.Context) []HomeNotice {
	type noticeRow struct {
		ID        string
		Title     string
		Level     string
		CreatedAt *time.Time
	}
	var rows []noticeRow
	db.DB.WithContext(ctx).Table("sys_notice").
		Where("status = ?", string(enums.StatusEnabled)).
		Where("category = ?", "PLATFORM").
		Order("sort_code ASC, is_top DESC").
		Select("id, title, level, created_at").
		Limit(5).
		Find(&rows)

	results := make([]HomeNotice, len(rows))
	for i, r := range rows {
		results[i] = HomeNotice{
			ID: r.ID, Title: r.Title, Level: r.Level,
			CreatedAt: utils.FormatDateTimePtr(r.CreatedAt),
		}
	}
	return results
}

func getUserCount(ctx context.Context) int {
	var count int64
	db.DB.WithContext(ctx).Table("sys_user").Count(&count)
	return int(count)
}
