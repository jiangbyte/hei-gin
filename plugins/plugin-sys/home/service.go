package home

import (
	"context"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/exception"
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
	vo.Stats.TotalUsers = UserCount(c.Request.Context())
	return vo
}

func HomeAddQuickAction(c *gin.Context, param *AddQuickActionParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()

	if QuickActionExists(ctx, userID, param.ResourceID) {
		return
	}

	actionCount := QuickActionCount(ctx, userID)
	entity := SysQuickAction{
		UserID: userID, ResourceID: param.ResourceID,
		SortCode: int(actionCount+1) * 10,
	}
	if err := CreateQuickAction(ctx, &entity); err != nil {
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
	if err := DeleteQuickActions(ctx, param.IDs); err != nil {
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
	if err := SortQuickActions(ctx, param.IDs); err != nil {
		result.WriteError(c, exception.NewBusinessError("排序快捷方式提交事务失败: "+err.Error(), 500))
		return
	}
}

func findQuickActionsByUserID(ctx context.Context, userID string) []QuickActionVO {
	actions := ListQuickActions(ctx, userID)
	if len(actions) == 0 {
		return make([]QuickActionVO, 0)
	}

	resourceIDs := make([]string, len(actions))
	for i, a := range actions {
		resourceIDs[i] = a.ResourceID
	}

	resources := ListResourcesByIDs(ctx, resourceIDs)
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
	actionIDs := QuickActionResourceIDs(ctx, userID)
	resources := AvailableResources(ctx, actionIDs)

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
	rows := LatestNotices(ctx)

	results := make([]HomeNotice, len(rows))
	for i, r := range rows {
		results[i] = HomeNotice{
			ID: r.ID, Title: r.Title, Level: r.Level,
			CreatedAt: utils.FormatDateTimePtr(r.CreatedAt),
		}
	}
	return results
}
