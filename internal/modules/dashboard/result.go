// internal/modules/dashboard/result.go 出参定义。
//
// Author: Charlie

package dashboard

// OverviewResult 仪表盘概览统计。
//
// Author: Charlie
type OverviewResult struct {
	AccountTotal    int64 `json:"account_total"`
	NoticeTotal     int64 `json:"notice_total"`
	FeedbackTotal   int64 `json:"feedback_total"`
	FeedbackPending int64 `json:"feedback_pending"`
	FileTotal       int64 `json:"file_total"`
	RoleTotal       int64 `json:"role_total"`
	DeptTotal       int64 `json:"dept_total"`
}
