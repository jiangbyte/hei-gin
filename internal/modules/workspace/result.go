// internal/modules/workspace/result.go 出参定义（对齐 hei-boot Workspace*Result）。
//
// Author: Charlie

package workspace

import "time"

// OverviewResult 工作台总览：快捷应用 + 本人近期操作/登录日志。
//
// Author: Charlie
type OverviewResult struct {
	Shortcuts        []ShortcutResult     `json:"shortcuts"`
	RecentOperations []ActivityItemResult `json:"recent_operations"`
	RecentLogins     []ActivityItemResult `json:"recent_logins"`
}

// ShortcutResult 工作台个人快捷应用项。
//
// Author: Charlie
type ShortcutResult struct {
	ID           string  `json:"id"`
	ResourceID   string  `json:"resource_id"`
	Sort         int     `json:"sort"`
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	Icon         *string `json:"icon"`
	Code         string  `json:"code"`
	ResourceType *string `json:"resource_type,omitempty"`
	Status       *string `json:"status,omitempty"`
}

// ActivityItemResult 工作台个人近期日志摘要项。
//
// Author: Charlie
type ActivityItemResult struct {
	ID           string    `json:"id"`
	Module       string    `json:"module"`
	ModuleLabel  *string   `json:"module_label"`
	Action       string    `json:"action"`
	ActionName   *string   `json:"action_name"`
	ActionType   *string   `json:"action_type"`
	Summary      *string   `json:"summary"`
	Success      bool      `json:"success"`
	IP           *string   `json:"ip"`
	UserAgent    *string   `json:"user_agent"`
	OperatorName *string   `json:"operator_name"`
	DurationMs   *string   `json:"duration_ms"`
	ResourceID   *string   `json:"resource_id"`
	CreatedAt    time.Time `json:"created_at"`
}
