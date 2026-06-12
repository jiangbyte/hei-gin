package position

import (
	"context"

	"hei-gin/sdk/db"
)

func Page(ctx context.Context, p *PositionPageParam) ([]SysPosition, int64) {
	q := db.DB.WithContext(ctx).Model(&SysPosition{})
	if p.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	var total int64
	q.Count(&total)
	var rows []SysPosition
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysPosition, error) {
	var e SysPosition
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysPosition) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysPosition{}).Where("id = ?", id).Updates(up).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	db.DB.WithContext(ctx).Table("sys_user").Where("position_id IN ?", ids).Update("position_id", nil)
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysPosition{}).Error
}

func ListAll(ctx context.Context) []SysPosition {
	var rows []SysPosition
	db.DB.WithContext(ctx).Model(&SysPosition{}).Order("sort_code ASC").Find(&rows)
	return rows
}
