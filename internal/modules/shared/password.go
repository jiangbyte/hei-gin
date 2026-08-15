// internal/modules/shared/password.go 密码策略（对齐 hei-boot PasswordPolicyApiProvider / PasswordHelper）。
//
// Author: Charlie

package shared

import (
	"context"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/runtimecfg"
)

// PasswordHistory 密码历史实体，对应表 sys_account_password_history。
//
// Author: Charlie
type PasswordHistory struct {
	ID           string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	AccountID    string    `gorm:"column:account_id;size:64;not null" json:"account_id"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null" json:"password_hash"`
	ChangedBy    *string   `gorm:"column:changed_by;size:64" json:"changed_by"`
	ChangeReason *string   `gorm:"column:change_reason;size:64" json:"change_reason"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 返回表名。
func (PasswordHistory) TableName() string { return "sys_account_password_history" }

// PasswordPolicy 密码策略校验与历史/过期辅助。
//
// Author: Charlie
type PasswordPolicy struct {
	db      *gorm.DB
	runtime *runtimecfg.Settings
}

// NewPasswordPolicy 构造密码策略。
func NewPasswordPolicy(db *gorm.DB, rt *runtimecfg.Settings) *PasswordPolicy {
	return &PasswordPolicy{db: db, runtime: rt}
}

func (p *PasswordPolicy) str(ctx context.Context, key, def string) string {
	if p.runtime != nil {
		return p.runtime.GetString(ctx, key, def)
	}
	return def
}

func (p *PasswordPolicy) num(ctx context.Context, key string, def int) int {
	if p.runtime != nil {
		return p.runtime.GetInt(ctx, key, def)
	}
	return def
}

func (p *PasswordPolicy) flag(ctx context.Context, key string, def bool) bool {
	if p.runtime != nil {
		return p.runtime.GetBool(ctx, key, def)
	}
	return def
}

// Validate 断言密码满足策略（复杂度、连续字符、弱口令、用户信息包含、历史复用）。
func (p *PasswordPolicy) Validate(ctx context.Context, raw, accountID, accountName, email, phone string) error {
	if strings.TrimSpace(raw) == "" {
		return ErrPasswordRequired
	}
	minLen := p.num(ctx, "PASSWORD_MIN_LENGTH", 8)
	maxLen := p.num(ctx, "PASSWORD_MAX_LENGTH", 128)
	if len(raw) < minLen {
		return &PolicyError{"密码长度不能少于 " + strconv.Itoa(minLen) + " 位"}
	}
	if len(raw) > maxLen {
		return &PolicyError{"密码长度不能超过 " + strconv.Itoa(maxLen) + " 位"}
	}
	if err := checkComplexity(raw, p.str(ctx, "PASSWORD_COMPLEXITY", "DIGITS_UPPER_LOWER_SPECIAL")); err != nil {
		return err
	}
	if maxConsecutive := p.num(ctx, "PASSWORD_MAX_CONSECUTIVE_CHARS", 3); maxConsecutive > 0 {
		if err := checkMaxConsecutive(raw, maxConsecutive); err != nil {
			return err
		}
	}
	if p.flag(ctx, "PASSWORD_FORBID_WEAK_LIST", true) {
		if isCustomWeak(raw, p.str(ctx, "PASSWORD_CUSTOM_WEAK_WORDS", "")) {
			return &PolicyError{"密码过于常见"}
		}
		if p.isWeakInTable(ctx, raw) {
			return &PolicyError{"密码过于常见"}
		}
	}
	if p.flag(ctx, "PASSWORD_FORBID_USER_INFO", true) && containsUserInfo(raw, accountName, email, phone) {
		return &PolicyError{"密码不能包含账号、邮箱或手机号信息"}
	}
	if accountID != "" && p.flag(ctx, "PASSWORD_FORBID_HISTORICAL", true) {
		limit := p.num(ctx, "PASSWORD_HISTORY_CHECK_COUNT", 5)
		if ok, _ := p.matchesRecent(ctx, accountID, raw, limit); ok {
			return &PolicyError{"不能使用近期使用过的密码"}
		}
	}
	return nil
}

// RecordHistory 记录密码变更历史（reason: register / admin_reset / self_reset / password_expired）。
func (p *PasswordPolicy) RecordHistory(ctx context.Context, accountID, raw, changedBy, reason string) error {
	if accountID == "" || raw == "" {
		return nil
	}
	hash, err := security.HashPassword(raw)
	if err != nil {
		return err
	}
	cb := changedBy
	cr := reason
	row := PasswordHistory{
		ID: idgen.Next(), AccountID: accountID, PasswordHash: hash,
		ChangedBy: &cb, ChangeReason: &cr, CreatedAt: time.Now().UTC(),
	}
	return p.db.WithContext(ctx).Create(&row).Error
}

// PasswordExpired 判断密码是否过期（PASSWORD_VALIDITY_DAYS<=0 不过期）；返回是否过期与剩余提醒天数。
func (p *PasswordPolicy) PasswordExpired(ctx context.Context, accountID string) (bool, int) {
	validityDays := p.num(ctx, "PASSWORD_VALIDITY_DAYS", 0)
	if validityDays <= 0 {
		return false, 0
	}
	ageDays, ok := p.passwordAgeDays(ctx, accountID)
	if !ok {
		return false, 0
	}
	if ageDays >= validityDays {
		return true, 0
	}
	warningDays := p.num(ctx, "PASSWORD_EXPIRY_WARNING_DAYS", 7)
	remain := validityDays - ageDays
	if warningDays > 0 && remain <= warningDays {
		return false, remain
	}
	return false, 0
}

func (p *PasswordPolicy) passwordAgeDays(ctx context.Context, accountID string) (int, bool) {
	var row PasswordHistory
	if err := p.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("created_at DESC").Limit(1).First(&row).Error; err != nil {
		return 0, false
	}
	days := int(time.Since(row.CreatedAt).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return days, true
}

func (p *PasswordPolicy) matchesRecent(ctx context.Context, accountID, raw string, limit int) (bool, error) {
	if limit <= 0 {
		return false, nil
	}
	var rows []PasswordHistory
	if err := p.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("created_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return false, err
	}
	for _, row := range rows {
		if security.CheckPassword(row.PasswordHash, raw) {
			return true, nil
		}
	}
	return false, nil
}

func (p *PasswordPolicy) isWeakInTable(ctx context.Context, raw string) bool {
	var n int64
	if err := p.db.WithContext(ctx).Table("sys_weak_password").
		Where("password = ?", strings.TrimSpace(raw)).Count(&n).Error; err != nil {
		return false
	}
	return n > 0
}

func isCustomWeak(raw, custom string) bool {
	lowered := strings.ToLower(raw)
	for _, part := range strings.Split(custom, ",") {
		word := strings.ToLower(strings.TrimSpace(part))
		if word != "" && word == lowered {
			return true
		}
	}
	return false
}

func containsUserInfo(raw, accountName, email, phone string) bool {
	lowered := strings.ToLower(raw)
	if matchesFragment(lowered, accountName) {
		return true
	}
	if email != "" {
		normalized := strings.ToLower(strings.TrimSpace(email))
		if matchesFragment(lowered, normalized) {
			return true
		}
		if at := strings.Index(normalized, "@"); at > 0 && matchesFragment(lowered, normalized[:at]) {
			return true
		}
	}
	return matchesFragment(lowered, phone)
}

func matchesFragment(loweredPassword, fragment string) bool {
	item := strings.ToLower(strings.TrimSpace(fragment))
	return len(item) >= 3 && strings.Contains(loweredPassword, item)
}

func checkComplexity(password, complexity string) error {
	key := strings.ToUpper(strings.TrimSpace(complexity))
	hasUpper, hasLower, hasDigit, hasSpecial, hasLetter := passwordClasses(password)

	switch key {
	case "", "NO_LIMIT":
		return nil
	case "DIGITS_AND_LETTERS":
		if !hasDigit || !hasLetter {
			return &PolicyError{"密码必须同时包含字母和数字"}
		}
		return nil
	case "DIGITS_AND_UPPERCASE":
		if !hasDigit || !hasUpper {
			return &PolicyError{"密码必须同时包含数字和大写字母"}
		}
		return nil
	case "DIGITS_UPPER_LOWER_SPECIAL":
		return requireClasses(hasUpper, hasLower, hasDigit, hasSpecial, true, true, true, true)
	case "TWO_OF_THREE":
		classes := 0
		if hasLetter {
			classes++
		}
		if hasDigit {
			classes++
		}
		if hasSpecial {
			classes++
		}
		if classes < 2 {
			return &PolicyError{"密码需至少包含字母、数字、特殊字符中的两类"}
		}
		return nil
	case "THREE_OF_FOUR":
		classes := 0
		if hasUpper {
			classes++
		}
		if hasLower {
			classes++
		}
		if hasDigit {
			classes++
		}
		if hasSpecial {
			classes++
		}
		if classes < 3 {
			return &PolicyError{"密码需至少包含大写、小写、数字、特殊字符中的三类"}
		}
		return nil
	default:
		return requireClasses(hasUpper, hasLower, hasDigit, hasSpecial, true, true, true, true)
	}
}

func passwordClasses(password string) (hasUpper, hasLower, hasDigit, hasSpecial, hasLetter bool) {
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper, hasLetter = true, true
		case unicode.IsLower(r):
			hasLower, hasLetter = true, true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsLetter(r):
			hasLetter = true
		default:
			hasSpecial = true
		}
	}
	return
}

func requireClasses(hasUpper, hasLower, hasDigit, hasSpecial, needUpper, needLower, needDigit, needSpecial bool) error {
	if needUpper && !hasUpper {
		return &PolicyError{"密码必须包含至少一个大写字母"}
	}
	if needLower && !hasLower {
		return &PolicyError{"密码必须包含至少一个小写字母"}
	}
	if needDigit && !hasDigit {
		return &PolicyError{"密码必须包含至少一个数字"}
	}
	if needSpecial && !hasSpecial {
		return &PolicyError{"密码必须包含至少一个特殊字符"}
	}
	return nil
}

func checkMaxConsecutive(password string, max int) error {
	if max <= 0 || len(password) <= max {
		return nil
	}
	run := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			run++
			if run > max {
				return &PolicyError{"密码不能包含超过 " + strconv.Itoa(max) + " 个连续相同字符"}
			}
		} else {
			run = 1
		}
	}
	return nil
}

// PolicyError 策略校验错误。
//
// Author: Charlie
type PolicyError struct{ msg string }

func (e *PolicyError) Error() string { return e.msg }

// ErrPasswordRequired 密码为空。
var ErrPasswordRequired = &PolicyError{"密码不能为空"}
