// internal/modules/sys/audit/operator.go 操作人昵称回显（对齐 hei-boot AuditOperatorSupport）。
//
// Author: Charlie

package audit

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"hei-gin/internal/modules/profile"
)

func enrichOperatorNames(ctx context.Context, db *gorm.DB, rows []OperationLog) {
	if len(rows) == 0 || db == nil {
		return
	}
	adminIDs := map[string]struct{}{}
	portalIDs := map[string]struct{}{}
	for i := range rows {
		if !needsNicknameBackfill(&rows[i]) {
			continue
		}
		if rows[i].AccountID == nil {
			continue
		}
		if isPortalAccountType(rows[i].AccountType) {
			portalIDs[*rows[i].AccountID] = struct{}{}
		} else {
			adminIDs[*rows[i].AccountID] = struct{}{}
		}
	}
	nicknames := map[string]string{}
	if len(adminIDs) > 0 {
		mergeNicknames(ctx, profile.AdminRepo(db), adminIDs, nicknames)
	}
	if len(portalIDs) > 0 {
		mergeNicknames(ctx, profile.PortalRepo(db), portalIDs, nicknames)
	}
	for i := range rows {
		if !needsNicknameBackfill(&rows[i]) || rows[i].AccountID == nil {
			continue
		}
		if nick, ok := nicknames[*rows[i].AccountID]; ok && nick != "" {
			rows[i].OperatorName = &nick
		}
	}
}

func mergeNicknames(ctx context.Context, repo *profile.Repo, ids map[string]struct{}, out map[string]string) {
	list := make([]string, 0, len(ids))
	for id := range ids {
		list = append(list, id)
	}
	profiles, err := repo.ListByAccountIDs(ctx, list)
	if err != nil {
		return
	}
	for id, p := range profiles {
		if p.Nickname != nil && strings.TrimSpace(*p.Nickname) != "" {
			out[id] = strings.TrimSpace(*p.Nickname)
		}
	}
}

func needsNicknameBackfill(row *OperationLog) bool {
	if row == nil || row.AccountID == nil {
		return false
	}
	if row.OperatorName == nil || strings.TrimSpace(*row.OperatorName) == "" {
		return true
	}
	return strings.TrimSpace(*row.OperatorName) == strings.TrimSpace(*row.AccountID)
}

func isPortalAccountType(accountType *string) bool {
	if accountType == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(*accountType), "portal")
}
