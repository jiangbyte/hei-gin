package home

import (
	"context"

	resModel "hei-gin/plugins/plugin-sys/resource"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type service struct {
	repo *repository
}

func (s *service) HomeGet(c *gin.Context) *HomeVO {
	userID := auth.GetLoginIDDefaultNull(c)
	vo := &HomeVO{
		QuickActions:       make([]QuickActionVO, 0),
		AvailableResources: make([]QuickActionVO, 0),
		Notices:            make([]HomeNotice, 0),
	}
	if userID != "" {
		vo.QuickActions = s.findQuickActionsByUserID(c.Request.Context(), userID)
		vo.AvailableResources = s.getAvailableResources(c.Request.Context(), userID)
	}
	vo.Notices = s.getNotices(c.Request.Context())
	vo.Stats.TotalUsers = s.repo.UserCount(c.Request.Context())
	return vo
}

func (s *service) HomeAddQuickAction(c *gin.Context, param *AddQuickActionParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()

	if s.repo.QuickActionExists(ctx, userID, param.ResourceID) {
		return
	}

	actionCount := s.repo.QuickActionCount(ctx, userID)
	entity := SysQuickAction{
		UserID: userID, ResourceID: param.ResourceID,
		SortCode: int(actionCount+1) * 10,
	}
	if err := s.repo.CreateQuickAction(ctx, &entity); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加快捷方式失败: "+err.Error(), 500))
		return
	}
}

func (s *service) HomeRemoveQuickAction(c *gin.Context, param *utils.IdsParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.DeleteQuickActions(ctx, param.IDs); err != nil {
		result.WriteError(c, exception.NewBusinessError("移除快捷方式失败: "+err.Error(), 500))
		return
	}
}

func (s *service) HomeSortQuickActions(c *gin.Context, param *utils.IdsParam) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("登录用户不存在", 500))
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.SortQuickActions(ctx, param.IDs); err != nil {
		result.WriteError(c, exception.NewBusinessError("排序快捷方式提交事务失败: "+err.Error(), 500))
		return
	}
}

func (s *service) findQuickActionsByUserID(ctx context.Context, userID string) []QuickActionVO {
	actions := s.repo.ListQuickActions(ctx, userID)
	if len(actions) == 0 {
		return make([]QuickActionVO, 0)
	}

	resourceIDs := make([]string, len(actions))
	for i, a := range actions {
		resourceIDs[i] = a.ResourceID
	}

	resources := s.repo.ListResourcesByIDs(ctx, resourceIDs)
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

func (s *service) getAvailableResources(ctx context.Context, userID string) []QuickActionVO {
	actionIDs := s.repo.QuickActionResourceIDs(ctx, userID)
	resources := s.repo.AvailableResources(ctx, actionIDs)

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

func (s *service) getNotices(ctx context.Context) []HomeNotice {
	rows := s.repo.LatestNotices(ctx)

	results := make([]HomeNotice, len(rows))
	for i, r := range rows {
		results[i] = HomeNotice{
			ID: r.ID, Title: r.Title, Level: r.Level,
			CreatedAt: utils.FormatDateTimePtr(r.CreatedAt),
		}
	}
	return results
}

func HomeGet(c *gin.Context) *HomeVO {
	return defaultModule.service.HomeGet(c)
}

func HomeAddQuickAction(c *gin.Context, param *AddQuickActionParam) {
	defaultModule.service.HomeAddQuickAction(c, param)
}

func HomeRemoveQuickAction(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.HomeRemoveQuickAction(c, param)
}

func HomeSortQuickActions(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.HomeSortQuickActions(c, param)
}
