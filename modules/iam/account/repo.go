package account

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/core/security"
)

// Repo 账号持久化（仅 sys_account / sys_account_identity；资料表归 user 模块）。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// DB 返回底层 DB（供同事务扩展；一般业务勿用）。
func (r *Repo) DB() *gorm.DB { return r.db }

// FindIdentity 按类型与标识查身份。
func (r *Repo) FindIdentity(ctx context.Context, identityType, identifier string) (*Identity, error) {
	var ident Identity
	if err := r.with(ctx).Where("identity_type = ? AND identifier = ?", identityType, identifier).First(&ident).Error; err != nil {
		return nil, err
	}
	return &ident, nil
}

// GetByID 按主键查账号。
func (r *Repo) GetByID(ctx context.Context, id string) (*Account, error) {
	var acc Account
	if err := r.with(ctx).First(&acc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

// FindAccountIdentity 查账号主登录身份。
func (r *Repo) FindAccountIdentity(ctx context.Context, accountID string) (*Identity, error) {
	var ident Identity
	if err := r.with(ctx).Where("account_id = ? AND identity_type = ?", accountID, IdentityAccount).First(&ident).Error; err != nil {
		return nil, err
	}
	return &ident, nil
}

// ListRoleIDs 查账号已启用角色 ID。
func (r *Repo) ListRoleIDs(ctx context.Context, accountID string) ([]string, error) {
	var roleRels []struct {
		TargetID string `gorm:"column:target_id"`
	}
	if err := r.with(ctx).Table("sys_iam_relation").
		Select("target_id").
		Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND status = ?", "ACCOUNT", accountID, "ACCOUNT_ROLE", "ENABLED").
		Find(&roleRels).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(roleRels))
	for _, row := range roleRels {
		out = append(out, row.TargetID)
	}
	return out, nil
}

type permRow struct {
	TargetKey string `gorm:"column:target_key"`
	DataScope string `gorm:"column:data_scope"`
	SourceID  string `gorm:"column:subject_id"`
}

// ListRolePermissions 按角色列出权限键。
func (r *Repo) ListRolePermissions(ctx context.Context, roleIDs []string) ([]permRow, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var rows []permRow
	err := r.with(ctx).Table("sys_iam_relation").
		Select("target_key, data_scope, subject_id").
		Where("relation_type = ? AND target_type = ? AND status = ?", "ROLE_PERMISSION", "PERMISSION", "ENABLED").
		Where("subject_type = ? AND subject_id IN ?", "ROLE", roleIDs).
		Find(&rows).Error
	return rows, err
}

// CreateAccount 事务创建账号与主身份。
func (r *Repo) CreateAccount(ctx context.Context, acc Account, ident Identity) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&acc).Error; err != nil {
			return err
		}
		return tx.Create(&ident).Error
	})
}

// UpdateAccount 更新账号字段与主登录标识。
func (r *Repo) UpdateAccount(ctx context.Context, id string, updates map[string]any, accountIdent string) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Account{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Model(&Identity{}).Where("account_id = ? AND identity_type = ?", id, IdentityAccount).
			Update("identifier", accountIdent).Error
	})
}

// DeleteByIDs 事务删除身份与账号（资料由 user 模块先删）。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id IN ?", ids).Delete(&Identity{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&Account{}).Error
	})
}

// PageAccounts 分页查询账号。
func (r *Repo) PageAccounts(ctx context.Context, p PageParam) (rows []Account, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Account{})
	if p.AccountType != "" {
		db = db.Where("account_type = ?", p.AccountType)
	}
	if p.AccountStatus != "" {
		db = db.Where("account_status = ?", p.AccountStatus)
	}
	if p.Account != "" {
		db = db.Where("id IN (SELECT account_id FROM sys_account_identity WHERE identity_type = ? AND identifier ILIKE ?)",
			IdentityAccount, "%"+p.Account+"%")
	}
	if p.Name != "" {
		db = db.Where(
			`(account_type = ? AND id IN (SELECT account_id FROM admin_user_profile WHERE name ILIKE ?))
			 OR (account_type = ? AND id IN (SELECT account_id FROM portal_user_profile WHERE name ILIKE ?))`,
			string(security.AccountAdmin), "%"+p.Name+"%", string(security.AccountPortal), "%"+p.Name+"%",
		)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
