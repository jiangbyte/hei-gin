// internal/modules/dashboard/result.go 出参定义（对齐 hei-boot DashboardOverviewResult 嵌套结构）。
//
// Author: Charlie

package dashboard

// OverviewResult 仪表盘概览统计（web dashboard/index.vue 消费契约）。
//
// Author: Charlie
type OverviewResult struct {
	Summary  SummaryResult  `json:"summary"`
	Accounts AccountsResult `json:"accounts"`
	IAM      IAMResult      `json:"iam"`
	OpsToday OpsTodayResult `json:"ops_today"`
	Trends   TrendsResult   `json:"trends"`
	Files    FilesResult    `json:"files"`
}

// SummaryResult 顶部汇总（对齐 hei-boot DashboardSummaryResult）。
//
// Author: Charlie
type SummaryResult struct {
	AccountTotal   int64 `json:"account_total"`
	OnlineSessions int64 `json:"online_sessions"`
	FileTotal      int64 `json:"file_total"`
	StorageBytes   int64 `json:"storage_bytes"`
}

// AccountsResult 账号分布（对齐 hei-boot DashboardAccountsResult）。
//
// Author: Charlie
type AccountsResult struct {
	Enabled  int64        `json:"enabled"`
	Disabled int64        `json:"disabled"`
	TodayNew int64        `json:"today_new"`
	ByType   []StatusItem `json:"by_type"`
}

// IAMResult IAM 资源计数（对齐 hei-boot DashboardIamResult）。
//
// Author: Charlie
type IAMResult struct {
	RoleCount  int64 `json:"role_count"`
	DeptCount  int64 `json:"dept_count"`
	GroupCount int64 `json:"group_count"`
	MenuCount  int64 `json:"menu_count"`
}

// OpsTodayResult 今日运维指标（对齐 hei-boot DashboardOpsTodayResult）。
//
// Author: Charlie
type OpsTodayResult struct {
	AuditTotal      int64 `json:"audit_total"`
	AuditFailed     int64 `json:"audit_failed"`
	FeedbackPending int64 `json:"feedback_pending"`
}

// TrendsResult 近七日趋势（对齐 hei-boot DashboardTrendsResult）。
//
// Author: Charlie
type TrendsResult struct {
	AccountTrend []TrendPoint `json:"account_trend"`
	AuditTrend   []TrendPoint `json:"audit_trend"`
}

// TrendPoint 趋势点（date/type/value；对齐 hei-boot DashboardTrendPointResult）。
//
// Author: Charlie
type TrendPoint struct {
	Date  string `json:"date"`
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

// FilesResult 文件类型分布（对齐 hei-boot DashboardFilesResult）。
//
// Author: Charlie
type FilesResult struct {
	ByContentType []StatusItem `json:"by_content_type"`
}

// StatusItem 名称/计数项（对齐 hei-boot DashboardStatusItemResult）。
//
// Author: Charlie
type StatusItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}
