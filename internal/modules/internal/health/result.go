// internal/modules/internal/health/result.go 出参定义。
//
// Author: Charlie

package health

// CheckItem å•é¡¹æŽ¢é’ˆç»“æžœã€‚
//
// Author: Charlie
type CheckItem struct {
	Enabled bool    `json:"enabled"`
	OK      bool    `json:"ok"`
	Detail  *string `json:"detail"`
}

// LiveResult å­˜æ´»æŽ¢é’ˆç»“æžœã€‚
//
// Author: Charlie
type LiveResult struct {
	Status string `json:"status"`
}

// ReadyResult å°±ç»ªæŽ¢é’ˆç»“æžœã€‚
//
// Author: Charlie
type ReadyResult struct {
	Status string `json:"status"`
	Checks struct {
		Database CheckItem `json:"database"`
		Redis    CheckItem `json:"redis"`
	} `json:"checks"`
}
