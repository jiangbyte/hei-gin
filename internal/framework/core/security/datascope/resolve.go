// Package datascope 按权限键解析数据范围（对齐 hei-boot DataScopeSupport）。
//
// Author: Charlie

package datascope

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// FindGrant 在会话授权中定位指定权限键（取最后一条）。
func FindGrant(sess *security.Session, permissionKey string) *security.PermissionGrant {
	if sess == nil || permissionKey == "" {
		return nil
	}
	var matched *security.PermissionGrant
	for i := range sess.PermissionGrants {
		g := &sess.PermissionGrants[i]
		if g.PermissionKey == permissionKey {
			matched = g
		}
	}
	return matched
}

// HasUnrestricted 是否超级权限或指定键为 ALL 范围。
func HasUnrestricted(sess *security.Session, permissionKey string) bool {
	if sess == nil {
		return false
	}
	for _, k := range sess.PermissionKeys {
		if k == "*:*:*" {
			return true
		}
	}
	if g := FindGrant(sess, permissionKey); g != nil && g.DataScope == ScopeAll {
		return true
	}
	return false
}

// EffectiveForKey 按权限键解析有效数据范围与部门 ID 列表。
func EffectiveForKey(db *gorm.DB, sess *security.Session, permissionKey string) (Scope, []string) {
	if sess == nil {
		return ScopeSelf, nil
	}
	if HasUnrestricted(sess, permissionKey) {
		return ScopeAll, nil
	}
	g := FindGrant(sess, permissionKey)
	scope := ScopeSelf
	var custom []string
	if g != nil {
		scope = g.DataScope
		custom = append(custom, g.CustomScopeDeptIDs...)
	}
	switch scope {
	case ScopeAll:
		return ScopeAll, nil
	case ScopeDept:
		return ScopeDept, unique(sess.DeptIDs)
	case ScopeDeptAndChild:
		return ScopeDeptAndChild, expandDeptAndChildren(db, unique(sess.DeptIDs))
	case ScopeCustom:
		return ScopeCustom, unique(custom)
	default:
		return ScopeSelf, nil
	}
}

// ApplyKey 按权限键过滤查询（owner 列用于 SELF；dept 列用于部门范围）。
func ApplyKey(db *gorm.DB, sess *security.Session, permissionKey, ownerDeptColumn, ownerAccountColumn string) *gorm.DB {
	if db == nil {
		return db
	}
	if sess == nil {
		return db.Where("1 = 0")
	}
	scope, deptIDs := EffectiveForKey(db, sess, permissionKey)
	switch scope {
	case ScopeAll:
		return db
	case ScopeSelf:
		col := sanitizeColumn(ownerAccountColumn)
		if col == "" {
			col = "created_by"
		}
		return db.Where(col+" = ?", sess.AccountID)
	case ScopeDept, ScopeDeptAndChild, ScopeCustom:
		col := sanitizeColumn(ownerDeptColumn)
		if col == "" || len(deptIDs) == 0 {
			if col := sanitizeColumn(ownerAccountColumn); col != "" {
				return db.Where(col+" = ?", sess.AccountID)
			}
			return db.Where("1 = 0")
		}
		return db.Where(col+" IN ?", deptIDs)
	default:
		col := sanitizeColumn(ownerAccountColumn)
		if col == "" {
			col = "created_by"
		}
		return db.Where(col+" = ?", sess.AccountID)
	}
}

// ApplyAccountScope 账号列表/分页：SELF 按 id；部门范围经 ACCOUNT_DEPT 子查询。
func ApplyAccountScope(db *gorm.DB, sess *security.Session, permissionKey string) *gorm.DB {
	if db == nil {
		return db
	}
	if sess == nil {
		return db.Where("1 = 0")
	}
	scope, deptIDs := EffectiveForKey(db, sess, permissionKey)
	switch scope {
	case ScopeAll:
		return db
	case ScopeSelf:
		return db.Where("id = ?", sess.AccountID)
	case ScopeDept, ScopeDeptAndChild, ScopeCustom:
		if len(deptIDs) == 0 {
			return db.Where("1 = 0")
		}
		return db.Where(`id IN (
			SELECT subject_id FROM sys_iam_relation
			WHERE subject_type = 'ACCOUNT' AND relation_type = 'ACCOUNT_DEPT'
			  AND target_type = 'DEPT' AND target_id IN ?
		)`, deptIDs)
	default:
		return db.Where("id = ?", sess.AccountID)
	}
}

// AssertKey 写操作数据范围校验（按权限键）。
func AssertKey(sess *security.Session, permissionKey, ownerDeptID, ownerAccountID string) error {
	if sess == nil {
		return ErrDenied
	}
	scope, deptIDs := EffectiveForKey(nil, sess, permissionKey)
	switch scope {
	case ScopeAll:
		return nil
	case ScopeSelf:
		if ownerAccountID != "" && ownerAccountID == sess.AccountID {
			return nil
		}
		return ErrDenied
	case ScopeDept, ScopeDeptAndChild, ScopeCustom:
		if ownerDeptID == "" {
			return ErrDenied
		}
		for _, id := range deptIDs {
			if id == ownerDeptID {
				return nil
			}
		}
		return ErrDenied
	default:
		return ErrDenied
	}
}

// AssertAccountAccessible 断言目标账号在权限键对应数据范围内可访问。
func AssertAccountAccessible(ctx context.Context, db *gorm.DB, sess *security.Session, accountID, permissionKey string) error {
	if sess == nil || accountID == "" {
		return ErrDenied
	}
	scope, deptIDs := EffectiveForKey(db, sess, permissionKey)
	switch scope {
	case ScopeAll:
		return nil
	case ScopeSelf:
		if accountID == sess.AccountID {
			return nil
		}
		return ErrDenied
	case ScopeDept, ScopeDeptAndChild, ScopeCustom:
		if len(deptIDs) == 0 {
			return ErrDenied
		}
		var cnt int64
		err := db.WithContext(ctx).Table("sys_iam_relation").
			Where("subject_type = ? AND subject_id = ? AND relation_type = ? AND target_type = ? AND target_id IN ?",
				"ACCOUNT", accountID, "ACCOUNT_DEPT", "DEPT", deptIDs).
			Count(&cnt).Error
		if err != nil || cnt == 0 {
			return ErrDenied
		}
		return nil
	default:
		return ErrDenied
	}
}

// AssertAccountAccessibleMsg 带业务文案的账号范围校验。
func AssertAccountAccessibleMsg(ctx context.Context, db *gorm.DB, sess *security.Session, accountID, permissionKey string) error {
	if err := AssertAccountAccessible(ctx, db, sess, accountID, permissionKey); err != nil {
		return fmt.Errorf("无权访问该数据")
	}
	return nil
}
