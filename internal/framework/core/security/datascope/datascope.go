// Package datascope æŒ‰ä¼šè¯æ•°æ®èŒƒå›´è¿‡æ»¤ GORM æŸ¥è¯¢ã€‚
//
// Author: Charlie
package datascope

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
)

// ErrDenied è¡¨ç¤ºå½“å‰ä¼šè¯æ— æƒè®¿é—®ç›®æ ‡æ•°æ®ã€‚
//
// Author: Charlie
var ErrDenied = errors.New("datascope: access denied")

// Scope ä¸Ž security.DataScope å¯¹é½çš„åˆ«åï¼Œä¾¿äºŽæœ¬åŒ…å¼•ç”¨ã€‚
type Scope = security.DataScope

const (
	ScopeAll          = security.DataScopeAll
	ScopeDeptAndChild = security.DataScopeDeptAndChild
	ScopeDept         = security.DataScopeDept
	ScopeSelf         = security.DataScopeSelf
	ScopeCustom       = security.DataScopeCustom
)

// Effective ä»Žä¼šè¯æƒé™æŽˆæƒè§£æžæœ‰æ•ˆæ•°æ®èŒƒå›´ä¸Žéƒ¨é—¨ ID åˆ—è¡¨ã€‚
// *:*:* æˆ–ä»»ä¸€ ALL æŽˆæƒè§†ä¸ºå…¨éƒ¨æ•°æ®ï¼›å¦åˆ™å–æŽˆæƒä¸­æœ€å®½èŒƒå›´ã€‚
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

// Apply æŒ‰ä¼šè¯æ•°æ®èŒƒå›´è¿‡æ»¤æŸ¥è¯¢ã€‚
// SELF ä½¿ç”¨ created_by = account_idï¼›éƒ¨é—¨ç±»èŒƒå›´ä½¿ç”¨ ownerDeptColumn IN (...)ã€‚
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

// Assert å†™æ“ä½œæ•°æ®èŒƒå›´æ ¡éªŒï¼šALL æ”¾è¡Œï¼›SELF æ¯” ownerAccountIDï¼›éƒ¨é—¨èŒƒå›´è¦æ±‚ ownerDeptID è½åœ¨å¯è§éƒ¨é—¨å†…ã€‚
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
		// DEPT_AND_CHILD åœ¨ Assert æ—¶ db ä¸º nilï¼Œæœªå±•å¼€å­éƒ¨é—¨ï¼šå›žé€€åˆ°ä¼šè¯ DeptIDs ç²¾ç¡®åŒ¹é…å·²åœ¨ä¸Šé¢å®Œæˆï¼›
		// è‹¥éœ€è¦å­éƒ¨é—¨æ ¡éªŒï¼Œè°ƒç”¨æ–¹åº”åœ¨æŸ¥è¯¢è·¯å¾„ç”¨ Applyã€‚
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
	// PostgreSQL recursive CTEï¼šæœ¬éƒ¨é—¨åŠå…¨éƒ¨å­éƒ¨é—¨ã€‚
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

// FormatDeny è¿”å›žå¸¦ä¸Šä¸‹æ–‡çš„æ‹’ç»é”™è¯¯ï¼ˆå¯é€‰ï¼‰ã€‚
func FormatDeny(detail string) error {
	if detail == "" {
		return ErrDenied
	}
	return fmt.Errorf("%w: %s", ErrDenied, detail)
}
