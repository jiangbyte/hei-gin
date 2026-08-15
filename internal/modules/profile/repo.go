// internal/modules/profile/repo.go 用户中心资料持久化（表名参数化，admin/portal 共用）。
//
// Author: Charlie

package profile

import (
	"context"

	"gorm.io/gorm"
)

// Repo 用户资料仓储（table 区分 profile_user_admin / profile_user_portal）。
//
// Author: Charlie
type Repo struct {
	db    *gorm.DB
	table string
}

// NewRepo 构造按表名绑定的仓储。
func NewRepo(db *gorm.DB, table string) *Repo { return &Repo{db: db, table: table} }

// AdminRepo 管理端资料仓储。
func AdminRepo(db *gorm.DB) *Repo { return NewRepo(db, ProfileTableAdmin) }

// PortalRepo 门户端资料仓储。
func PortalRepo(db *gorm.DB) *Repo { return NewRepo(db, ProfileTablePortal) }

// Table 返回绑定的表名。
func (r *Repo) Table() string { return r.table }

// DB 返回底层 DB（供跨模块清理等）。
func (r *Repo) DB() *gorm.DB { return r.db }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table(r.table)
}

// ListByAccountIDs 批量查资料，返回 account_id → Profile（供列表页避免 N+1）。
func (r *Repo) ListByAccountIDs(ctx context.Context, ids []string) (map[string]*Profile, error) {
	out := make(map[string]*Profile, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var rows []Profile
	if err := r.with(ctx).Where("account_id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		out[rows[i].AccountID] = &rows[i]
	}
	return out, nil
}

// GetProfile 按账号 ID 查资料。
func (r *Repo) GetProfile(ctx context.Context, accountID string) (*Profile, error) {
	var p Profile
	if err := r.with(ctx).First(&p, "account_id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateProfile 创建资料。
func (r *Repo) CreateProfile(ctx context.Context, p *Profile) error {
	return r.with(ctx).Create(p).Error
}

// UpdateProfile 更新资料字段。
func (r *Repo) UpdateProfile(ctx context.Context, accountID string, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}
	return r.with(ctx).Where("account_id = ?", accountID).Updates(updates).Error
}

// UpsertProfile 按 account_id 插入或更新资料（供 iam/account 跨模块调用）。
func (r *Repo) UpsertProfile(ctx context.Context, p *Profile) error {
	var n int64
	if err := r.with(ctx).Where("account_id = ?", p.AccountID).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return r.CreateProfile(ctx, p)
	}
	updates := map[string]any{
		"name": p.Name, "nickname": p.Nickname, "avatar": p.Avatar, "signature": p.Signature,
		"phone": p.Phone, "email": p.Email,
	}
	if r.table == ProfileTableAdmin {
		updates["remark"] = p.Remark
	}
	return r.UpdateProfile(ctx, p.AccountID, updates)
}

// DeleteByAccountIDs 按账号 ID 批量删除资料。
func (r *Repo) DeleteByAccountIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.with(ctx).Where("account_id IN ?", ids).Delete(&Profile{}).Error
}

// GetAccountPassword 查账号密码哈希。
func (r *Repo) GetAccountPassword(ctx context.Context, accountID string) (*AccountPassword, error) {
	var acc AccountPassword
	if err := r.db.WithContext(ctx).First(&acc, "id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

// UpdateAccountPassword 更新密码哈希。
func (r *Repo) UpdateAccountPassword(ctx context.Context, accountID, hash string) error {
	return r.db.WithContext(ctx).Model(&AccountPassword{}).Where("id = ?", accountID).
		Update("password_hash", hash).Error
}

// GetAccountIdentifier 返回账号主登录标识（sys_account_identity ACCOUNT 类型）。
func (r *Repo) GetAccountIdentifier(ctx context.Context, accountID string) (string, error) {
	var identifier string
	err := r.db.WithContext(ctx).Table("sys_account_identity").
		Select("identifier").Where("account_id = ? AND identity_type = ?", accountID, "ACCOUNT").
		Order("is_primary DESC, id ASC").Limit(1).Scan(&identifier).Error
	return identifier, err
}

// LoadIDNames 批量按 ID 加载 {id,name}（保持入参顺序，缺省跳过）。
func (r *Repo) LoadIDNames(ctx context.Context, table string, ids []string) []IDName {
	if len(ids) == 0 {
		return []IDName{}
	}
	var rows []struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	if err := r.db.WithContext(ctx).Table(table).
		Select("id", "name").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return []IDName{}
	}
	byID := make(map[string]string, len(rows))
	for _, row := range rows {
		byID[row.ID] = row.Name
	}
	out := make([]IDName, 0, len(ids))
	for _, id := range ids {
		if name, ok := byID[id]; ok {
			out = append(out, IDName{ID: id, Name: name})
		}
	}
	return out
}
