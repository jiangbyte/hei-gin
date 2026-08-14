package relation

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// Repo å…³ç³»æŒä¹…åŒ–ï¼ˆsys_iam_relationï¼‰ã€‚
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo æž„é€  Repoã€‚
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// ListRelations åˆ—å‡ºä¸»ä½“æŒ‡å®šå…³ç³»ç±»åž‹çš„å…³ç³»è¡Œï¼ˆaccountType ä¸ºç©ºä¸è¿‡æ»¤ï¼‰ã€‚
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

// deleteSubjectRelations åˆ é™¤ä¸»ä½“æŒ‡å®šå…³ç³»ç±»åž‹çš„å…³ç³»ï¼ˆaccountType ä¸ºç©ºåˆ å…¨éƒ¨ç±»åž‹ï¼Œä¾›äº‹åŠ¡å†…è°ƒç”¨ï¼‰ã€‚
func (r *Repo) deleteSubjectRelations(db *gorm.DB, subjectType, subjectID, relationType, accountType string) error {
	q := db.Where("subject_type = ? AND subject_id = ? AND relation_type = ?", subjectType, subjectID, relationType)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	return q.Delete(&Relation{}).Error
}

// DeleteBySubjectIDs æŒ‰ä¸»ä½“ id é›†åˆåˆ é™¤æŒ‡å®šå…³ç³»ç±»åž‹çš„å…³ç³»ï¼ˆæ‰¹é‡æ¸…æŒ‰é’®æƒé™ç»‘å®šç”¨ï¼‰ã€‚
func (r *Repo) DeleteBySubjectIDs(ctx context.Context, subjectType string, subjectIDs []string, relationType string) error {
	if len(subjectIDs) == 0 {
		return nil
	}
	return r.with(ctx).Where("subject_type = ? AND subject_id IN ? AND relation_type = ?",
		subjectType, subjectIDs, relationType).Delete(&Relation{}).Error
}

// CreateInBatches æ‰¹é‡æ’å…¥å…³ç³»è¡Œï¼ˆä¾›äº‹åŠ¡å†…è°ƒç”¨ï¼‰ã€‚
func (r *Repo) CreateInBatches(db *gorm.DB, rows []Relation) error {
	return db.CreateInBatches(rows, 200).Error
}
