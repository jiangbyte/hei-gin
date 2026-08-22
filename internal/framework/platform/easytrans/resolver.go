// internal/framework/platform/easytrans/resolver.go ID→名称批量解析。
//
// Author: Charlie

package easytrans

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type idNameRow struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

// LookupNames 按主键批量查 name 列（表需含 id、name）。
func LookupNames(ctx context.Context, db *gorm.DB, table string, ids []string) map[string]string {
	out := map[string]string{}
	if db == nil || table == "" || len(ids) == 0 {
		return out
	}
	seen := map[string]struct{}{}
	uniq := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	if len(uniq) == 0 {
		return out
	}
	var rows []idNameRow
	_ = db.WithContext(ctx).Table(table).Select("id, name").Where("id IN ?", uniq).Scan(&rows).Error
	for _, r := range rows {
		out[r.ID] = r.Name
	}
	return out
}

// LookupAccountIdentifiers 批量查 ACCOUNT 登录标识。
func LookupAccountIdentifiers(ctx context.Context, db *gorm.DB, accountIDs []string) map[string]string {
	out := map[string]string{}
	if db == nil || len(accountIDs) == 0 {
		return out
	}
	type row struct {
		AccountID  string `gorm:"column:account_id"`
		Identifier string `gorm:"column:identifier"`
	}
	var rows []row
	_ = db.WithContext(ctx).Table("sys_account_identity").
		Select("account_id, identifier").
		Where("account_id IN ? AND identity_type = ? AND bind_status = ?", accountIDs, "ACCOUNT", "BOUND").
		Scan(&rows).Error
	for _, r := range rows {
		if _, ok := out[r.AccountID]; !ok {
			out[r.AccountID] = r.Identifier
		}
	}
	return out
}

// Ptr 返回字符串指针；空串返回 nil（boot 部分字段为 null）。
func Ptr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// Str 空指针转空串（boot 脱敏字段常用 ""）。
func Str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// EmptySlice 保证 JSON 序列化为 [] 而非 null。
func EmptySlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
