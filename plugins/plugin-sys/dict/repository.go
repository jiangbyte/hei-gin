package dict

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *DictPageParam) ([]SysDict, int64) {
	q := r.db.WithContext(ctx).Model(&SysDict{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("code LIKE ? OR label LIKE ? OR value LIKE ?", like, like, like)
	}
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.ParentID != "" {
		q = q.Where("id = ? OR parent_id = ?", p.ParentID, p.ParentID)
	}
	if p.DictGroup == "FRM" {
		q = q.Where("category = ?", "FRM")
	}
	if p.DictGroup == "BIZ" {
		q = q.Where("category = ?", "BIZ")
	}
	var total int64
	q.Count(&total)
	var rows []SysDict
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) ListForTree(ctx context.Context, category, dictGroup string) []SysDict {
	q := r.db.WithContext(ctx).Model(&SysDict{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if dictGroup == "FRM" {
		q = q.Where("category = ?", "FRM")
	}
	if dictGroup == "BIZ" {
		q = q.Where("category = ?", "BIZ")
	}
	var all []SysDict
	q.Find(&all)
	return all
}

func (r *repository) FindByID(ctx context.Context, id string) (*SysDict, error) {
	var e SysDict
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *SysDict) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&SysDict{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysDict{}).Error
}

func (r *repository) ListOptions(ctx context.Context, category, parentID string) []SysDict {
	q := r.db.WithContext(ctx).Model(&SysDict{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if parentID != "" {
		q = q.Where("id = ? OR parent_id = ?", parentID, parentID)
	}
	var records []SysDict
	q.Find(&records)
	return records
}

func (r *repository) ListByCategoryAndKeyword(ctx context.Context, category, keyword string) []SysDict {
	q := r.db.WithContext(ctx).Model(&SysDict{}).Order("sort_code ASC")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("label LIKE ? OR code LIKE ?", kw, kw)
	}
	var records []SysDict
	q.Find(&records)
	return records
}

func (r *repository) FindByTypeCodeAndValue(ctx context.Context, typeCode, value string) (*SysDict, error) {
	var entity SysDict
	if err := r.db.WithContext(ctx).
		Where("parent_id IN (SELECT id FROM sys_dict WHERE code = ?)", typeCode).
		Where("value = ?", value).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*SysDict, error) {
	var parent SysDict
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&parent).Error; err != nil {
		return nil, err
	}
	return &parent, nil
}

func (r *repository) ListChildren(ctx context.Context, parentID string) []SysDict {
	var records []SysDict
	r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort_code ASC").Find(&records)
	return records
}

func (r *repository) CountDuplicateValue(ctx context.Context, parentID *string, value string, excludeID string) int64 {
	var count int64
	q := r.db.WithContext(ctx).Model(&SysDict{}).Where("parent_id = ?", parentID).Where("value = ?", value)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count
}

func (r *repository) ListAll(ctx context.Context) []SysDict {
	var all []SysDict
	r.db.WithContext(ctx).Find(&all)
	return all
}
