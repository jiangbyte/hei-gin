package account

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/core/security"
)

// Repo 账号持久化。
//
// Author: Charlie
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

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

// GetAdminProfile 查管理端资料。
func (r *Repo) GetAdminProfile(ctx context.Context, accountID string) (*AdminUserProfile, error) {
	var p AdminUserProfile
	if err := r.with(ctx).First(&p, "account_id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPortalProfile 查门户资料。
func (r *Repo) GetPortalProfile(ctx context.Context, accountID string) (*PortalUserProfile, error) {
	var p PortalUserProfile
	if err := r.with(ctx).First(&p, "account_id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &p, nil
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

// CreateBundle 事务创建账号、身份与资料。
func (r *Repo) CreateBundle(ctx context.Context, acc Account, ident Identity, admin *AdminUserProfile, portal *PortalUserProfile) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&acc).Error; err != nil {
			return err
		}
		if err := tx.Create(&ident).Error; err != nil {
			return err
		}
		if admin != nil {
			return tx.Create(admin).Error
		}
		if portal != nil {
			return tx.Create(portal).Error
		}
		return nil
	})
}

// UpdateBundle 事务更新账号、身份与资料。
func (r *Repo) UpdateBundle(ctx context.Context, id string, updates map[string]any, accountIdent string, accountType string, profile map[string]any) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Account{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		_ = tx.Model(&Identity{}).Where("account_id = ? AND identity_type = ?", id, IdentityAccount).
			Update("identifier", accountIdent).Error
		if accountType == string(security.AccountAdmin) {
			return tx.Model(&AdminUserProfile{}).Where("account_id = ?", id).Updates(profile).Error
		}
		return tx.Model(&PortalUserProfile{}).Where("account_id = ?", id).Updates(profile).Error
	})
}

// DeleteByIDs 事务删除账号及关联。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id IN ?", ids).Delete(&Identity{}).Error; err != nil {
			return err
		}
		if err := tx.Where("account_id IN ?", ids).Delete(&AdminUserProfile{}).Error; err != nil {
			return err
		}
		if err := tx.Where("account_id IN ?", ids).Delete(&PortalUserProfile{}).Error; err != nil {
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
