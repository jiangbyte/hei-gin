// internal/modules/profile/identity/param.go 入参定义（snake_case，对齐 hei-boot Jackson SNAKE_CASE）。
//
// Author: Charlie
package identity

import "hei-gin/internal/framework/core/schema"

// RealNameCaseSubmitParam 提交实名认证工单（人工通道）。
type RealNameCaseSubmitParam struct {
	BusinessType      string   `json:"business_type"`
	DocumentType      string   `json:"document_type" binding:"required"`
	RealName          string   `json:"real_name" binding:"required"`
	DocumentNo        string   `json:"document_no" binding:"required"`
	AttachmentIDs     []string `json:"attachment_ids"`
	ApplicantContact  string   `json:"applicant_contact"`
}

// RealNameCaseInitThirdPartyParam 发起第三方实名认证。
type RealNameCaseInitThirdPartyParam struct {
	BusinessType string `json:"business_type"`
	DocumentType string `json:"document_type" binding:"required"`
	RealName     string `json:"real_name" binding:"required"`
	DocumentNo   string `json:"document_no" binding:"required"`
	Provider     string `json:"provider"`
}

// RealNameCaseCallbackParam 第三方实名认证回调。
type RealNameCaseCallbackParam struct {
	CaseID          string `json:"case_id" binding:"required"`
	ProviderOrderNo string `json:"provider_order_no"`
	Success         *bool  `json:"success" binding:"required"`
	Message         string `json:"message"`
}

// RealNameCaseApproveParam 审核通过。
type RealNameCaseApproveParam struct {
	CaseID string `json:"case_id" binding:"required"`
	Remark string `json:"remark"`
}

// RealNameCaseRejectParam 审核驳回。
type RealNameCaseRejectParam struct {
	CaseID       string `json:"case_id" binding:"required"`
	RejectReason string `json:"reject_reason" binding:"required"`
}

// IdentityRevokeParam 撤销实名快照。
type IdentityRevokeParam struct {
	AccountID string `json:"account_id" binding:"required"`
	Remark    string `json:"remark"`
}

// RealNameCaseMyPageParam 我的实名工单分页。
type RealNameCaseMyPageParam struct {
	schema.PageQuery
	BusinessType string `form:"business_type" json:"business_type"`
	Status       string `form:"status" json:"status"`
}

// RealNameCaseReviewPageParam 审核队列分页。
type RealNameCaseReviewPageParam struct {
	schema.PageQuery
	BusinessType string `form:"business_type" json:"business_type"`
	Status       string `form:"status" json:"status"`
	AccountID    string `form:"account_id" json:"account_id"`
}

// IdentityPageParam 实名快照分页。
type IdentityPageParam struct {
	schema.PageQuery
	Status       string `form:"status" json:"status"`
	AccountID    string `form:"account_id" json:"account_id"`
	DocumentType string `form:"document_type" json:"document_type"`
}
