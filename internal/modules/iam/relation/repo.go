// internal/modules/iam/relation/repo.go 持久化仓储。
//
// Author: Charlie

package relation

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo 关系持久化（sys_iam_relation）。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// ListRelations 列出主体指定关系类型的关系行（accountType 为空不过滤）。
func (r *Repo) ListRelations(ctx context.Context, subjectType, subjectID, relationType, accountType string) ([]Relation, error) {
	db := r.with(ctx).Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND status = ?",
		subjectType, subjectID, relationType, security.StatusEnabled)
	if accountType != "" {
		db = db.Where("account_type = ?", accountType)
	}
	var rows []Relation
	if err := db.Order("sort asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// deleteSubjectRelations 删除主体指定关系类型的关系（accountType 为空删全部类型，供事务内调用）。
func (r *Repo) deleteSubjectRelations(db *gorm.DB, subjectType, subjectID, relationType, accountType string) error {
	q := db.Where("subject_type = ? AND subject_id = ? AND relation_type = ?", subjectType, subjectID, relationType)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	return q.Delete(&Relation{}).Error
}

// DeleteBySubjectIDs 按主体 id 集合删除指定关系类型的关系（批量清按钮权限绑定用）。
func (r *Repo) DeleteBySubjectIDs(ctx context.Context, subjectType string, subjectIDs []string, relationType string) error {
	if len(subjectIDs) == 0 {
		return nil
	}
	return r.with(ctx).Where("subject_type = ? AND subject_id IN ? AND relation_type = ?",
		subjectType, subjectIDs, relationType).Delete(&Relation{}).Error
}

// CreateInBatches 批量插入关系行（供事务内调用）。
func (r *Repo) CreateInBatches(db *gorm.DB, rows []Relation) error {
	return db.CreateInBatches(rows, 200).Error
}
