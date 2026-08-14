// internal/modules/health/result.go 出参定义。
//
// Author: Charlie

package health

// CheckItem 单项探针结果。
//
// Author: Charlie
type CheckItem struct {
	Enabled bool    `json:"enabled"`
	OK      bool    `json:"ok"`
	Detail  *string `json:"detail"`
}

// LiveResult 存活探针结果。
//
// Author: Charlie
type LiveResult struct {
	Status string `json:"status"`
}

// ReadyResult 就绪探针结果。
//
// Author: Charlie
type ReadyResult struct {
	Status string `json:"status"`
	Checks struct {
		Database CheckItem `json:"database"`
		Redis    CheckItem `json:"redis"`
	} `json:"checks"`
}
