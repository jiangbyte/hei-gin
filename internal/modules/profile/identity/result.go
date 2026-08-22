// internal/modules/profile/identity/result.go 出参定义（snake_case，对齐 hei-boot Jackson SNAKE_CASE）。
//
// Author: Charlie
package identity

import "time"

// IdentityStatusResult 当前账号实名认证快照与进行中工单摘要。
type IdentityStatusResult struct {
	Status           string                     `json:"status"`
	DocumentType     *string                    `json:"document_type"`
	RealNameMasked   string                     `json:"real_name_masked"`
	DocumentNoMasked string                     `json:"document_no_masked"`
	VerifyChannel    *string                    `json:"verify_channel"`
	Provider         *string                    `json:"provider"`
	VerifiedAt       *time.Time                 `json:"verified_at"`
	RevokedAt        *time.Time                 `json:"revoked_at"`
	PendingCase      *RealNameCaseSummaryResult `json:"pending_case"`
}

// RealNameCaseOptionsResult 实名认证可选项。
type RealNameCaseOptionsResult struct {
	BusinessTypes []RealNameBusinessOptionResult `json:"business_types"`
	DocumentTypes []string                       `json:"document_types"`
}

// RealNameBusinessOptionResult 业务类型与可用通道。
type RealNameBusinessOptionResult struct {
	BusinessType string   `json:"business_type"`
	Label        string   `json:"label"`
	Channels     []string `json:"channels"`
}

// RealNameCaseInitResult 第三方实名认证初始化结果。
type RealNameCaseInitResult struct {
	CaseID          string `json:"case_id"`
	Provider        string `json:"provider"`
	ProviderOrderNo string `json:"provider_order_no"`
	RedirectURL     string `json:"redirect_url"`
}

// RealNameCaseSummaryResult 实名工单摘要。
type RealNameCaseSummaryResult struct {
	CaseID           string     `json:"case_id"`
	AccountID        *string    `json:"account_id"`
	BusinessType     string     `json:"business_type"`
	VerifyChannel    string     `json:"verify_channel"`
	Status           string     `json:"status"`
	DocumentType     *string    `json:"document_type"`
	RealNameMasked   *string    `json:"real_name_masked"`
	DocumentNoMasked *string    `json:"document_no_masked"`
	CreatedAt        time.Time  `json:"created_at"`
	ReviewedAt       *time.Time `json:"reviewed_at"`
	RejectReason     *string    `json:"reject_reason"`
}

// RealNameCaseDetailResult 实名工单详情。
type RealNameCaseDetailResult struct {
	RealNameCaseSummaryResult
	Provider        *string                      `json:"provider"`
	ProviderOrderNo *string                      `json:"provider_order_no"`
	SubmitterID     *string                      `json:"submitter_id"`
	ReviewerID      *string                      `json:"reviewer_id"`
	Attachments     []RealNameCaseAttachmentResult `json:"attachments"`
}

// RealNameCaseAttachmentResult 工单附件。
type RealNameCaseAttachmentResult struct {
	ObjectName   string  `json:"object_name"`
	ID           *string `json:"id"`
	OriginalName *string `json:"original_name"`
	ContentType  *string `json:"content_type"`
	Size         *int64  `json:"size"`
	URL          *string `json:"url"`
}

// IdentityPageResult 实名快照分页行。
type IdentityPageResult struct {
	AccountID        string     `json:"account_id"`
	Status           string     `json:"status"`
	DocumentType     *string    `json:"document_type"`
	RealNameMasked   *string    `json:"real_name_masked"`
	DocumentNoMasked *string    `json:"document_no_masked"`
	VerifyChannel    *string    `json:"verify_channel"`
	Provider         *string    `json:"provider"`
	VerifiedAt       *time.Time `json:"verified_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
}
