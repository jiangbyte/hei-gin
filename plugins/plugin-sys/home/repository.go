package home

import (
	"context"
	"time"

	resModel "hei-gin/plugins/plugin-sys/resource"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
)

type homeNoticeRow struct {
	ID        string
	Title     string
	Level     string
	CreatedAt *time.Time
}

func QuickActionExists(ctx context.Context, userID, resourceID string) bool {
	var count int64
	db.DB.WithContext(ctx).Model(&SysQuickAction{}).Where("user_id = ? AND resource_id = ?", userID, resourceID).Count(&count)
	return count > 0
}

func QuickActionCount(ctx context.Context, userID string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&SysQuickAction{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func CreateQuickAction(ctx context.Context, entity *SysQuickAction) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func DeleteQuickActions(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysQuickAction{}).Error
}

func SortQuickActions(ctx context.Context, ids []string) error {
	tx := db.DB.WithContext(ctx).Begin()
	for idx, id := range ids {
		if err := tx.Model(&SysQuickAction{}).Where("id = ?", id).Update("sort_code", (idx+1)*10).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func ListQuickActions(ctx context.Context, userID string) []SysQuickAction {
	var actions []SysQuickAction
	db.DB.WithContext(ctx).Where("user_id = ?", userID).Order("sort_code ASC, created_at ASC").Find(&actions)
	return actions
}

func ListResourcesByIDs(ctx context.Context, ids []string) []resModel.SysResource {
	var resources []resModel.SysResource
	db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&resources)
	return resources
}

func QuickActionResourceIDs(ctx context.Context, userID string) []string {
	var actionIDs []string
	db.DB.WithContext(ctx).Model(&SysQuickAction{}).Where("user_id = ?", userID).Select("resource_id").Find(&actionIDs)
	return actionIDs
}

func AvailableResources(ctx context.Context, actionIDs []string) []resModel.SysResource {
	q := db.DB.WithContext(ctx).Model(&resModel.SysResource{}).Where(
		"category IN ? AND status = ?",
		[]string{string(enums.ResourceCategoryBackendMenu), string(enums.ResourceCategoryFrontendMenu)},
		string(enums.StatusEnabled),
	)
	if len(actionIDs) > 0 {
		q = q.Where("id NOT IN ?", actionIDs)
	}
	var resources []resModel.SysResource
	q.Order("sort_code ASC").Find(&resources)
	return resources
}

func LatestNotices(ctx context.Context) []homeNoticeRow {
	var rows []homeNoticeRow
	db.DB.WithContext(ctx).Table("sys_notice").
		Where("status = ?", string(enums.StatusEnabled)).
		Where("category = ?", "PLATFORM").
		Order("sort_code ASC, is_top DESC").
		Select("id, title, level, created_at").
		Limit(5).
		Find(&rows)
	return rows
}

func UserCount(ctx context.Context) int {
	var count int64
	db.DB.WithContext(ctx).Table("sys_user").Count(&count)
	return int(count)
}
