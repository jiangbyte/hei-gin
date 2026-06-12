package banner

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *BannerPageParam) ([]SysBanner, int64) {
	q := r.db.WithContext(ctx).Model(&SysBanner{})
	var total int64
	q.Count(&total)
	var rows []SysBanner
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysBanner, error) {
	var e SysBanner
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysBanner) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysBanner{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysBanner{}).Error
}

func (r *repository) ListAll(ctx context.Context) []SysBanner {
	var rows []SysBanner
	r.db.WithContext(ctx).Model(&SysBanner{}).Order("sort_code ASC").Find(&rows)
	return rows
}
