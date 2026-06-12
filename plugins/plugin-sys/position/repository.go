package position

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *PositionPageParam) ([]SysPosition, int64) {
	q := r.db.WithContext(ctx).Model(&SysPosition{})
	if p.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	var total int64
	q.Count(&total)
	var rows []SysPosition
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysPosition, error) {
	var e SysPosition
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysPosition) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysPosition{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	r.db.WithContext(ctx).Table("sys_user").Where("position_id IN ?", ids).Update("position_id", nil)
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysPosition{}).Error
}

func (r *repository) ListAll(ctx context.Context) []SysPosition {
	var rows []SysPosition
	r.db.WithContext(ctx).Model(&SysPosition{}).Order("sort_code ASC").Find(&rows)
	return rows
}
