package notice

import (
	"context"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
)

func Page(ctx context.Context, p *NoticePageParam) ([]SysNotice, int64) {
	q := db.DB.WithContext(ctx).Model(&SysNotice{})
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
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func PublicPage(ctx context.Context, p *NoticePageParam) ([]SysNotice, int64) {
	q := db.DB.WithContext(ctx).Model(&SysNotice{}).Where("status = ?", string(enums.StatusEnabled))
	if p.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+p.Keyword+"%")
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	var total int64
	q.Count(&total)
	var rows []SysNotice
	q.Order("is_top DESC, sort_code DESC, created_at DESC").Limit(p.Size).Offset((p.Current-1)*p.Size).Find(&rows)
	return rows, total
}

func FindByID(ctx context.Context, id string) (*SysNotice, error) {
	var e SysNotice
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func FindEnabledByID(ctx context.Context, id string) (*SysNotice, error) {
	var e SysNotice
	if err := db.DB.WithContext(ctx).Where("id = ? AND status = ?", id, string(enums.StatusEnabled)).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func Create(ctx context.Context, entity *SysNotice) error {
	return db.DB.WithContext(ctx).Create(entity).Error
}

func UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&SysNotice{}).Where("id = ?", id).Updates(up).Error
}

func DeleteByIDs(ctx context.Context, ids []string) error {
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysNotice{}).Error
}

func ListAll(ctx context.Context) []SysNotice {
	var rows []SysNotice
	db.DB.WithContext(ctx).Model(&SysNotice{}).Order("sort_code ASC").Find(&rows)
	return rows
}

func Latest(ctx context.Context, size int) []SysNotice {
	var rows []SysNotice
	db.DB.WithContext(ctx).
		Where("status = ?", string(enums.StatusEnabled)).
		Order("is_top DESC, sort_code DESC, created_at DESC").
		Limit(size).
		Find(&rows)
	return rows
}
