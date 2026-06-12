package org

import (
	"context"

	groupModel "hei-gin/plugins/plugin-sys/group"
	posModel "hei-gin/plugins/plugin-sys/position"
	userModel "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/db"
)

func Page(ctx context.Context, p *OrgPageParam) ([]SysOrg, int64) {
	q := db.DB.WithContext(ctx).Model(&SysOrg{})
	if p.ParentID != "" {
		if p.ParentID == "0" {
			q = q.Where("(parent_id IS NULL OR parent_id = '' OR id = ?)", p.ParentID)
		} else {
			q = q.Where("(id = ? OR parent_id = ?)", p.ParentID, p.ParentID)
		}
	}
	if p.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+p.Keyword+"%")
	}
	var total int64
	q.Count(&total)
	var rows []SysOrg
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func List(ctx context.Context, category string) ([]SysOrg, error) {
	q := db.DB.WithContext(ctx).Model(&SysOrg{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	var all []SysOrg
	if err := q.Find(&all).Error; err != nil {
		return nil, err
	}
	return all, nil
}

func ListAll(ctx context.Context) ([]SysOrg, error) {
	var all []SysOrg
	if err := db.DB.WithContext(ctx).Find(&all).Error; err != nil {
		return nil, err
	}
	return all, nil
}

func FindByID(ctx context.Context, id string) (*SysOrg, error) {
	var e SysOrg
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysOrg) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysOrg{}).Where("id = ?", id).Updates(up).Error
}

func CountUsersByOrgIDs(ctx context.Context, ids []string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&userModel.SysUser{}).Where("org_id IN ?", ids).Count(&count)
	return count
}

func CountGroupsByOrgIDs(ctx context.Context, ids []string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&groupModel.SysGroup{}).Where("org_id IN ?", ids).Count(&count)
	return count
}

func CountPositionsByOrgIDs(ctx context.Context, ids []string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&posModel.SysPosition{}).Where("org_id IN ?", ids).Count(&count)
	return count
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysOrg{}).Error
}
