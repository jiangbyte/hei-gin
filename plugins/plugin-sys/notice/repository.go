package notice

import (
	"context"

	"gorm.io/gorm"
	"hei-gin/plugins/plugin-sys/shared"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *NoticePageParam) ([]SysNotice, int64) {
	q := r.db.WithContext(ctx).Model(&SysNotice{})
	if p.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	var total int64
	q.Count(&total)
	var rows []SysNotice
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) PublicPage(ctx context.Context, p *NoticePageParam) ([]SysNotice, int64) {
	q := r.db.WithContext(ctx).Model(&SysNotice{}).Where("status = ?", shared.StatusEnabled)
	if p.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	var total int64
	q.Count(&total)
	var rows []SysNotice
	q.Order("is_top DESC, sort_code DESC, created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysNotice, error) {
	var e SysNotice
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) FindEnabledByID(ctx context.Context, id string) (*SysNotice, error) {
	var e SysNotice
	if err := r.db.WithContext(ctx).Where("id = ? AND status = ?", id, shared.StatusEnabled).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysNotice) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysNotice{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysNotice{}).Error
}

func (r *repository) ListAll(ctx context.Context) []SysNotice {
	var rows []SysNotice
	r.db.WithContext(ctx).Model(&SysNotice{}).Order("sort_code ASC").Find(&rows)
	return rows
}

func (r *repository) Latest(ctx context.Context, size int) []SysNotice {
	var rows []SysNotice
	r.db.WithContext(ctx).
		Where("status = ?", shared.StatusEnabled).
		Order("is_top DESC, sort_code DESC, created_at DESC").
		Limit(size).
		Find(&rows)
	return rows
}
