// internal/modules/iam/account/repo.go 持久化仓储。
//
// Author: Charlie

package account

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
	"hei-gin/internal/framework/platform/db/dialect"
	"hei-gin/internal/modules/auth/oauth"
)

// passwordKeyPrefix Redis 密码加密密钥前缀（与 auth 模块一致）。
const passwordKeyPrefix = "password:crypto:"

// Repo 账号持久化（仅 sys_account / sys_account_identity；资料表归 user 模块）。
//
// Author: Charlie
type Repo struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB, rdb *redis.Client) *Repo { return &Repo{db: db, rdb: rdb} }

// DecryptPassword 用 Redis 私钥解密管理端 RSA 加密密码（对齐 hei-boot passwordCryptoApi）。
func (r *Repo) DecryptPassword(ctx context.Context, keyID, encryptedValue string) (string, error) {
	if keyID == "" {
		return "", fmt.Errorf("缺少 password_key_id")
	}
	raw, err := r.rdb.Get(ctx, passwordKeyPrefix+keyID).Result()
	_ = r.rdb.Del(ctx, passwordKeyPrefix+keyID)
	if err == redis.Nil || raw == "" {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	if err != nil {
		return "", err
	}
	if encryptedValue == "" {
		return "", nil
	}
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	rsaKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("无效或过期的密码加密密钥")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("无效的加密密码")
	}
	plain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("无效的加密密码")
	}
	return string(plain), nil
}

// FindIdentitiesByAccountIDs 批量查账号身份。
func (r *Repo) FindIdentitiesByAccountIDs(ctx context.Context, accountIDs []string) (map[string][]Identity, error) {
	out := make(map[string][]Identity, len(accountIDs))
	if len(accountIDs) == 0 {
		return out, nil
	}
	var rows []Identity
	if err := r.with(ctx).Where("account_id IN ?", accountIDs).
		Order("account_id, identity_type, is_primary DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.AccountID] = append(out[row.AccountID], row)
	}
	return out, nil
}

// FindOAuthBindingsByAccountIDs 批量查账号三方绑定。
func (r *Repo) FindOAuthBindingsByAccountIDs(ctx context.Context, accountIDs []string) (map[string][]oauth.AccountOAuthBinding, error) {
	out := make(map[string][]oauth.AccountOAuthBinding, len(accountIDs))
	if len(accountIDs) == 0 {
		return out, nil
	}
	var rows []oauth.AccountOAuthBinding
	if err := r.with(ctx).Where("account_id IN ?", accountIDs).
		Order("account_id, provider").Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.AccountID] = append(out[row.AccountID], row)
	}
	return out, nil
}

// FindIdentities 查账号全部身份。
func (r *Repo) FindIdentities(ctx context.Context, accountID string) ([]Identity, error) {
	var rows []Identity
	err := r.with(ctx).Where("account_id = ?", accountID).Order("identity_type, is_primary DESC").Find(&rows).Error
	return rows, err
}

// FindOAuthBindings 查账号三方绑定。
func (r *Repo) FindOAuthBindings(ctx context.Context, accountID string) ([]oauth.AccountOAuthBinding, error) {
	var rows []oauth.AccountOAuthBinding
	err := r.with(ctx).Where("account_id = ?", accountID).Order("provider").Find(&rows).Error
	return rows, err
}

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

// FindAccountIdentities 批量查账号主登录标识（保持去重，先到先得）。
func (r *Repo) FindAccountIdentities(ctx context.Context, ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}
	var rows []struct {
		AccountID  string `gorm:"column:account_id"`
		Identifier string `gorm:"column:identifier"`
	}
	if err := r.with(ctx).Table("sys_account_identity").
		Select("account_id", "identifier").
		Where("account_id IN ? AND identity_type = ?", ids, IdentityAccount).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		if _, ok := out[row.AccountID]; !ok {
			out[row.AccountID] = row.Identifier
		}
	}
	return out, nil
}

// DeleteIdentitiesExcept 删除账号除指定类型外的全部身份。
func (r *Repo) DeleteIdentitiesExcept(ctx context.Context, accountID, keepType string) error {
	return r.with(ctx).Where("account_id = ? AND identity_type <> ?", accountID, keepType).Delete(&Identity{}).Error
}

// CreateIdentity 创建身份。
func (r *Repo) CreateIdentity(ctx context.Context, row *Identity) error {
	return r.with(ctx).Create(row).Error
}

// ListGroupIDs 查账号已加入的用户组 ID。
func (r *Repo) ListGroupIDs(ctx context.Context, accountID string) ([]string, error) {
	var rows []struct {
		TargetID string `gorm:"column:target_id"`
	}
	if err := r.with(ctx).Table("sys_iam_relation").
		Select("target_id").
		Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND status = ?",
			"ACCOUNT", accountID, "ACCOUNT_GROUP", "ENABLED").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.TargetID)
	}
	return out, nil
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
	TargetKey          string   `gorm:"column:target_key"`
	DataScope          string   `gorm:"column:data_scope"`
	SourceID           string   `gorm:"column:subject_id"`
	SourceType         string   `gorm:"-"`
	CustomScopeDeptIDs []string `gorm:"-"`
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

// UpdatePasswordHash 仅更新密码哈希。
func (r *Repo) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	return r.with(ctx).Model(&Account{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
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

// PageAccounts 分页查询账号；sess 非空时按数据范围过滤（SELF→created_by）。
func (r *Repo) PageAccounts(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []Account, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&Account{})
	if sess != nil {
		db = datascope.ApplyAccountScope(db, sess, "iam:account:page")
	}
	if p.AccountType != "" {
		db = db.Where("account_type = ?", p.AccountType)
	}
	if p.AccountStatus != "" {
		db = db.Where("account_status = ?", p.AccountStatus)
	}
	if p.Account != "" {
		db = db.Where("id IN (SELECT account_id FROM sys_account_identity WHERE identity_type = ? AND "+dialect.ILike(db, "identifier")+")",
			IdentityAccount, "%"+p.Account+"%")
	}
	if p.Name != "" {
		db = db.Where(
			`(account_type = ? AND id IN (SELECT account_id FROM profile_user_admin WHERE `+dialect.ILike(db, "nickname")+`))
			 OR (account_type = ? AND id IN (SELECT account_id FROM profile_user_portal WHERE `+dialect.ILike(db, "nickname")+`))`,
			string(security.AccountAdmin), "%"+p.Name+"%", string(security.AccountPortal), "%"+p.Name+"%",
		)
	}
	if p.Phone != "" {
		db = db.Where(
			`id IN (SELECT account_id FROM profile_user_admin WHERE `+dialect.ILike(db, "phone")+`)
			 OR id IN (SELECT account_id FROM profile_user_portal WHERE `+dialect.ILike(db, "phone")+`)`,
			"%"+p.Phone+"%", "%"+p.Phone+"%",
		)
	}
	if p.Email != "" {
		db = db.Where(
			`id IN (SELECT account_id FROM profile_user_admin WHERE `+dialect.ILike(db, "email")+`)
			 OR id IN (SELECT account_id FROM profile_user_portal WHERE `+dialect.ILike(db, "email")+`)`,
			"%"+p.Email+"%", "%"+p.Email+"%",
		)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// ListDeptIDs 查账号已关联部门 ID。
func (r *Repo) ListDeptIDs(ctx context.Context, accountID, accountType string) ([]string, error) {
	var ids []string
	q := r.with(ctx).Table("sys_iam_relation").
		Select("target_id").
		Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND target_type = ? AND status = ?",
			"ACCOUNT", accountID, "ACCOUNT_DEPT", "DEPT", "ENABLED")
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	if err := q.Scan(&ids).Error; err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// ListRoleCodes 按角色 ID 查编码。
func (r *Repo) ListRoleCodes(ctx context.Context, roleIDs []string) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}
	var codes []string
	if err := r.with(ctx).Table("sys_role").
		Select("code").Where("id IN ?", roleIDs).Scan(&codes).Error; err != nil {
		return nil, err
	}
	if codes == nil {
		codes = []string{}
	}
	return codes, nil
}

// ListRoleIDsByGroups 查用户组拥有的角色 ID。
func (r *Repo) ListRoleIDsByGroups(ctx context.Context, groupIDs []string, accountType string) ([]string, error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}
	var ids []string
	q := r.with(ctx).Table("sys_iam_relation").
		Select("target_id").
		Where("subject_type = ? AND subject_id IN ? AND relation_type = ? AND target_type = ? AND status = ?",
			"GROUP", groupIDs, "GROUP_ROLE", "ROLE", "ENABLED")
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	if err := q.Scan(&ids).Error; err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}

// ListClientPermissionKeys 查客户端权限键。
func (r *Repo) ListClientPermissionKeys(ctx context.Context, accountID string, roleIDs, groupIDs []string, accountType string) ([]string, error) {
	cond := "(subject_type = ? AND subject_id = ?"
	args := []any{"ACCOUNT", accountID}
	if len(groupIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, "GROUP", groupIDs)
	}
	if len(roleIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, "ROLE", roleIDs)
	}
	cond += ")"
	fullArgs := append([]any{[]string{
		"SUBJECT_CLIENT_RESOURCE_GRANT", "ACCOUNT_CLIENT_RESOURCE", "GROUP_CLIENT_RESOURCE", "ROLE_CLIENT_RESOURCE",
	}, "CLIENT_RESOURCE", "ENABLED"}, args...)
	var keys []string
	err := r.with(ctx).Table("sys_iam_relation").
		Select("DISTINCT target_key").
		Where("relation_type IN ? AND target_type = ? AND status = ? AND target_key <> '' AND "+cond, fullArgs...).
		Scan(&keys).Error
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, nil
}

// ListGrantedClientResourceIDs 查已授权客户端资源 ID。
func (r *Repo) ListGrantedClientResourceIDs(ctx context.Context, accountID string, groupIDs, roleIDs []string, accountType string) ([]string, error) {
	cond := "(subject_type = ? AND subject_id = ?"
	args := []any{"ACCOUNT", accountID}
	if len(groupIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, "GROUP", groupIDs)
	}
	if len(roleIDs) > 0 {
		cond += " OR (subject_type = ? AND subject_id IN ?)"
		args = append(args, "ROLE", roleIDs)
	}
	cond += ")"
	fullArgs := append([]any{[]string{
		"SUBJECT_CLIENT_RESOURCE_GRANT", "ACCOUNT_CLIENT_RESOURCE", "GROUP_CLIENT_RESOURCE", "ROLE_CLIENT_RESOURCE",
	}, "CLIENT_RESOURCE", "ENABLED"}, args...)
	var ids []string
	q := r.with(ctx).Table("sys_iam_relation").
		Select("DISTINCT target_id").
		Where("relation_type IN ? AND target_type = ? AND status = ? AND "+cond, fullArgs...)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	if err := q.Scan(&ids).Error; err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}
