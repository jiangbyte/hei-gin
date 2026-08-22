// internal/modules/profile/identity/repo.go 数据访问。
//
// Author: Charlie
package identity

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/modules/sys/file"
)

// Repo 实名认证仓储。
type Repo struct {
	db *gorm.DB
}

// NewRepo 构造仓储。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// GetIdentity 按账号 ID 查询实名快照。
func (r *Repo) GetIdentity(ctx context.Context, accountID string) (*ProfileIdentity, error) {
	var row ProfileIdentity
	err := r.with(ctx).Where("account_id = ?", accountID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// SaveIdentity 插入或更新实名快照。
func (r *Repo) SaveIdentity(ctx context.Context, row *ProfileIdentity) error {
	var n int64
	if err := r.with(ctx).Model(&ProfileIdentity{}).Where("account_id = ?", row.AccountID).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return r.with(ctx).Create(row).Error
	}
	return r.with(ctx).Save(row).Error
}

// PageIdentity 实名快照分页。
func (r *Repo) PageIdentity(ctx context.Context, q IdentityPageParam) (rows []ProfileIdentity, total int64, current, size int, err error) {
	current, size = q.Normalize()
	tx := r.with(ctx).Model(&ProfileIdentity{})
	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if q.AccountID != "" {
		tx = tx.Where("account_id = ?", q.AccountID)
	}
	if q.DocumentType != "" {
		tx = tx.Where("document_type = ?", q.DocumentType)
	}
	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, current, size, err
	}
	err = tx.Order("verified_at DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, current, size, err
}

// FindVerifiedByDocumentHash 查询已绑定该证件的实名快照。
func (r *Repo) FindVerifiedByDocumentHash(ctx context.Context, hash, excludeAccountID string) (*ProfileIdentity, error) {
	tx := r.with(ctx).Where("document_no_hash = ? AND status = ?", hash, StatusVerified)
	if excludeAccountID != "" {
		tx = tx.Where("account_id <> ?", excludeAccountID)
	}
	var row ProfileIdentity
	err := tx.First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// CreateCase 创建实名工单。
func (r *Repo) CreateCase(ctx context.Context, row *RealNameCase) error {
	if row.CaseID == "" {
		row.CaseID = idgen.Next()
	}
	return r.with(ctx).Create(row).Error
}

// UpdateCase 更新实名工单。
func (r *Repo) UpdateCase(ctx context.Context, row *RealNameCase) error {
	return r.with(ctx).Save(row).Error
}

// GetCase 按 case_id 查询工单。
func (r *Repo) GetCase(ctx context.Context, caseID string) (*RealNameCase, error) {
	var row RealNameCase
	err := r.with(ctx).Where("case_id = ?", caseID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// CountPendingCases 统计账号进行中工单数。
func (r *Repo) CountPendingCases(ctx context.Context, accountID, businessType string) (int64, error) {
	var n int64
	err := r.with(ctx).Model(&RealNameCase{}).
		Where("account_id = ? AND business_type = ? AND status = ?", accountID, businessType, CaseStatusPending).
		Count(&n).Error
	return n, err
}

// FindPendingCaseByDocumentHash 查询证件进行中的工单。
func (r *Repo) FindPendingCaseByDocumentHash(ctx context.Context, hash, excludeAccountID string) (*RealNameCase, error) {
	tx := r.with(ctx).Where("document_no_hash = ? AND status = ?", hash, CaseStatusPending)
	if excludeAccountID != "" {
		tx = tx.Where("account_id <> ?", excludeAccountID)
	}
	var row RealNameCase
	err := tx.First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// FindLatestPendingCase 查询账号最新进行中工单。
func (r *Repo) FindLatestPendingCase(ctx context.Context, accountID string) (*RealNameCase, error) {
	var row RealNameCase
	err := r.with(ctx).
		Where("account_id = ? AND status = ?", accountID, CaseStatusPending).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// PageCasesByAccount 我的工单分页。
func (r *Repo) PageCasesByAccount(ctx context.Context, accountID string, q RealNameCaseMyPageParam) (rows []RealNameCase, total int64, current, size int, err error) {
	current, size = q.Normalize()
	tx := r.with(ctx).Model(&RealNameCase{}).Where("account_id = ?", accountID)
	if q.BusinessType != "" {
		tx = tx.Where("business_type = ?", q.BusinessType)
	}
	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, current, size, err
	}
	err = tx.Order("created_at DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, current, size, err
}

// PageReviewCases 审核队列分页。
func (r *Repo) PageReviewCases(ctx context.Context, q RealNameCaseReviewPageParam) (rows []RealNameCase, total int64, current, size int, err error) {
	current, size = q.Normalize()
	businessType := q.BusinessType
	if businessType == "" {
		businessType = BusinessAccountVerify
	}
	tx := r.with(ctx).Model(&RealNameCase{}).Where("business_type = ?", businessType)
	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if q.AccountID != "" {
		tx = tx.Where("account_id = ?", q.AccountID)
	}
	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, current, size, err
	}
	err = tx.Order("created_at DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, current, size, err
}

// AppendRecord 写入工单流水。
func (r *Repo) AppendRecord(ctx context.Context, entity *RealNameCase, action string, statusBefore, statusAfter, operatorID, remark string) error {
	before, after := statusBefore, statusAfter
	rec := RealNameCaseRecord{
		RecordID:      idgen.Next(),
		CaseID:        entity.CaseID,
		AccountID:     entity.AccountID,
		BusinessType:  entity.BusinessType,
		Action:        action,
		StatusBefore:  strPtr(before),
		StatusAfter:   strPtr(after),
		VerifyChannel: &entity.VerifyChannel,
		Provider:      entity.Provider,
		OperatorID:    strPtr(operatorID),
		Remark:        strPtr(remark),
	}
	return r.with(ctx).Create(&rec).Error
}

// ListFilesByObjectNames 批量查询文件元数据。
func (r *Repo) ListFilesByObjectNames(ctx context.Context, objectNames []string) ([]file.File, error) {
	if len(objectNames) == 0 {
		return nil, nil
	}
	var rows []file.File
	err := r.with(ctx).Where("object_name IN ?", objectNames).Find(&rows).Error
	return rows, err
}

func strPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
