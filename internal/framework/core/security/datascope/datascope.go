// Package datascope 按会话数据范围过滤 GORM 查询。
//
// Author: Charlie
package datascope

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// ErrDenied 表示当前会话无权访问目标数据。
//
// Author: Charlie
var ErrDenied = errors.New("datascope: access denied")

// Scope 与 security.DataScope 对齐的别名，便于本包引用。
type Scope = security.DataScope

const (
	ScopeAll          = security.DataScopeAll
	ScopeDeptAndChild = security.DataScopeDeptAndChild
	ScopeDept         = security.DataScopeDept
	ScopeSelf         = security.DataScopeSelf
	ScopeCustom       = security.DataScopeCustom
)

// Effective 从会话权限授权解析有效数据范围与部门 ID 列表。
// *:*:* 或任一 ALL 授权视为全部数据；否则取授权中最宽范围。
//
// Author: Charlie
func Effective(db *gorm.DB, sess *security.Session) (Scope, []string) {
	if sess == nil {
		return ScopeSelf, nil
	}
	for _, k := range sess.PermissionKeys {
		if k == "*:*:*" {
			return ScopeAll, nil
		}
	}
	scope := ScopeSelf
	var custom []string
	for _, g := range sess.PermissionGrants {
		switch g.DataScope {
		case ScopeAll:
			return ScopeAll, nil
		case ScopeDeptAndChild:
			if scope != ScopeAll {
				scope = ScopeDeptAndChild
			}
		case ScopeDept:
			if scope != ScopeAll && scope != ScopeDeptAndChild {
				scope = ScopeDept
			}
		case ScopeCustom:
			if scope == ScopeSelf {
				scope = ScopeCustom
			}
			custom = append(custom, g.CustomScopeDeptIDs...)
		}
	}
	switch scope {
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

// Apply 按会话数据范围过滤查询。
// SELF 使用 created_by = account_id；部门类范围使用 ownerDeptColumn IN (...)。
//
// Author: Charlie
func Apply(db *gorm.DB, sess *security.Session, ownerDeptColumn string) *gorm.DB {
	if db == nil {
		return db
	}
	if sess == nil {
		return db.Where("1 = 0")
	}
	scope, deptIDs := Effective(db, sess)
	switch scope {
	case ScopeAll:
		return db
	case ScopeSelf:
		return db.Where("created_by = ?", sess.AccountID)
	case ScopeDept, ScopeDeptAndChild, ScopeCustom:
		col := sanitizeColumn(ownerDeptColumn)
		if col == "" || len(deptIDs) == 0 {
			return db.Where("1 = 0")
		}
		return db.Where(col+" IN ?", deptIDs)
	default:
		return db.Where("created_by = ?", sess.AccountID)
	}
}

// Assert 写操作数据范围校验：ALL 放行；SELF 比 ownerAccountID；部门范围要求 ownerDeptID 落在可见部门内。
//
// Author: Charlie
func Assert(sess *security.Session, ownerDeptID, ownerAccountID string) error {
	if sess == nil {
		return ErrDenied
	}
	scope, deptIDs := Effective(nil, sess)
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
		// DEPT_AND_CHILD 在 Assert 时 db 为 nil，未展开子部门：回退到会话 DeptIDs 精确匹配已在上面完成；
		// 若需要子部门校验，调用方应在查询路径用 Apply。
		return ErrDenied
	default:
		return ErrDenied
	}
}

func sanitizeColumn(col string) string {
	col = strings.TrimSpace(col)
	if col == "" {
		return ""
	}
	for _, r := range col {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '.' {
			continue
		}
		return ""
	}
	return col
}

func unique(ids []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func expandDeptAndChildren(db *gorm.DB, roots []string) []string {
	if db == nil || len(roots) == 0 {
		return roots
	}
	var ids []string
	// PostgreSQL recursive CTE：本部门及全部子部门。
	q := `
WITH RECURSIVE tree AS (
  SELECT id FROM sys_dept WHERE id IN ?
  UNION ALL
  SELECT d.id FROM sys_dept d INNER JOIN tree t ON d.parent_id = t.id
)
SELECT DISTINCT id FROM tree`
	if err := db.Raw(q, roots).Scan(&ids).Error; err != nil {
		return roots
	}
	if len(ids) == 0 {
		return roots
	}
	return unique(ids)
}
