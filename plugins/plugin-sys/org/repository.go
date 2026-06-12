package org

import (
	"context"

	"gorm.io/gorm"
	groupModel "hei-gin/plugins/plugin-sys/group"
	posModel "hei-gin/plugins/plugin-sys/position"
	userModel "hei-gin/plugins/plugin-sys/user"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *OrgPageParam) ([]SysOrg, int64) {
	q := r.db.WithContext(ctx).Model(&SysOrg{})
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
	q.Order("sort_code ASC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) List(ctx context.Context, category string) ([]SysOrg, error) {
	q := r.db.WithContext(ctx).Model(&SysOrg{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	var all []SysOrg
	if err := q.Find(&all).Error; err != nil {
		return nil, err
	}
	return all, nil
}

func (r *repository) ListAll(ctx context.Context) ([]SysOrg, error) {
	var all []SysOrg
	if err := r.db.WithContext(ctx).Find(&all).Error; err != nil {
		return nil, err
	}
	return all, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysOrg, error) {
	var e SysOrg
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysOrg) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysOrg{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) CountUsersByOrgIDs(ctx context.Context, ids []string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&userModel.SysUser{}).Where("org_id IN ?", ids).Count(&count)
	return count
}

func (r *repository) CountGroupsByOrgIDs(ctx context.Context, ids []string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&groupModel.SysGroup{}).Where("org_id IN ?", ids).Count(&count)
	return count
}

func (r *repository) CountPositionsByOrgIDs(ctx context.Context, ids []string) int64 {
	var count int64
	r.db.WithContext(ctx).Model(&posModel.SysPosition{}).Where("org_id IN ?", ids).Count(&count)
	return count
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysOrg{}).Error
}
