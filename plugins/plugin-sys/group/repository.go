package group

import (
	"context"

	"hei-gin/sdk/db"
)

func Page(ctx context.Context, p *GroupPageParam) ([]SysGroup, int64) {
	q := db.DB.WithContext(ctx).Model(&SysGroup{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.OrgID != "" {
		q = q.Where("org_id = ?", p.OrgID)
	}
	var total int64
	q.Count(&total)
	var rows []SysGroup
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func List(ctx context.Context, category, orgID string) []SysGroup {
	q := db.DB.WithContext(ctx).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if orgID != "" {
		q = q.Where("org_id = ?", orgID)
	}
	var all []SysGroup
	q.Find(&all)
	return all
}

func ListAll(ctx context.Context) []SysGroup {
	var records []SysGroup
	db.DB.WithContext(ctx).Order("sort_code ASC").Find(&records)
	return records
}

func FindByID(ctx context.Context, id string) (*SysGroup, error) {
	var e SysGroup
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysGroup) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysGroup{}).Where("id = ?", id).Updates(up).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysGroup{}).Error
}

func ClearUserGroupIDs(ctx context.Context, ids []string) {
	db.DB.WithContext(ctx).Table("sys_user").Where("group_id IN ?", ids).Update("group_id", nil)
}
