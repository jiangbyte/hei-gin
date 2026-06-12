package banner

import (
	"context"

	"hei-gin/sdk/db"
)

func Page(ctx context.Context, p *BannerPageParam) ([]SysBanner, int64) {
	q := db.DB.WithContext(ctx).Model(&SysBanner{})
	var total int64
	q.Count(&total)
	var rows []SysBanner
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysBanner, error) {
	var e SysBanner
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysBanner) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysBanner{}).Where("id = ?", id).Updates(up).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysBanner{}).Error
}

func ListAll(ctx context.Context) []SysBanner {
	var rows []SysBanner
	db.DB.WithContext(ctx).Model(&SysBanner{}).Order("sort_code ASC").Find(&rows)
	return rows
}
