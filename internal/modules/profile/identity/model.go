// Package identity 实名认证（对齐 hei-boot profile.modules.identity）。
//
// Author: Charlie
package identity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const (
	StatusUnverified = "UNVERIFIED"
	StatusVerified   = "VERIFIED"
	StatusRevoked    = "REVOKED"

	CaseStatusPending  = "PENDING"
	CaseStatusApproved = "APPROVED"
	CaseStatusRejected = "REJECTED"

	BusinessAccountVerify = "ACCOUNT_VERIFY"

	ChannelManual      = "MANUAL"
	ChannelThirdParty  = "THIRD_PARTY"
	ProviderMock       = "MOCK"
)

// ProfileIdentity 账号实名快照，对应表 profile_identity。
type ProfileIdentity struct {
	AccountID        string     `gorm:"column:account_id;primaryKey;size:64" json:"account_id"`
	Status           string     `gorm:"column:status;size:32;not null;default:UNVERIFIED" json:"status"`
	DocumentType     *string    `gorm:"column:document_type;size:32" json:"document_type"`
	RealNameCipher   *string    `gorm:"column:real_name_cipher;type:text" json:"-"`
	DocumentNoCipher *string    `gorm:"column:document_no_cipher;type:text" json:"-"`
	DocumentNoHash   *string    `gorm:"column:document_no_hash;size:128" json:"-"`
	VerifyChannel    *string    `gorm:"column:verify_channel;size:32" json:"verify_channel"`
	Provider         *string    `gorm:"column:provider;size:32" json:"provider"`
	ProviderOrderNo  *string    `gorm:"column:provider_order_no;size:128" json:"provider_order_no"`
	VerifiedAt       *time.Time `gorm:"column:verified_at" json:"verified_at"`
	SourceCaseID     *string    `gorm:"column:source_case_id;size:64" json:"source_case_id"`
	RevokedAt        *time.Time `gorm:"column:revoked_at" json:"revoked_at"`
	RevokedBy        *string    `gorm:"column:revoked_by;size:64" json:"revoked_by"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy        *string    `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy        *string    `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (ProfileIdentity) TableName() string { return "profile_identity" }

// RealNameCase 实名业务工单，对应表 real_name_case。
type RealNameCase struct {
	CaseID                  string     `gorm:"column:case_id;primaryKey;size:64" json:"case_id"`
	BusinessType            string     `gorm:"column:business_type;size:64;not null" json:"business_type"`
	VerifyChannel           string     `gorm:"column:verify_channel;size:32;not null" json:"verify_channel"`
	Status                  string     `gorm:"column:status;size:32;not null" json:"status"`
	AccountID               *string    `gorm:"column:account_id;size:64" json:"account_id"`
	TargetAccountHintCipher *string    `gorm:"column:target_account_hint_cipher;type:text" json:"-"`
	ApplicantContactCipher  *string    `gorm:"column:applicant_contact_cipher;type:text" json:"-"`
	DocumentType            *string    `gorm:"column:document_type;size:32" json:"document_type"`
	RealNameCipher          *string    `gorm:"column:real_name_cipher;type:text" json:"-"`
	DocumentNoCipher        *string    `gorm:"column:document_no_cipher;type:text" json:"-"`
	DocumentNoHash          *string    `gorm:"column:document_no_hash;size:128" json:"-"`
	AttachmentIDs           StringList `gorm:"column:attachment_ids;type:text" json:"attachment_ids"`
	PayloadCipher           *string    `gorm:"column:payload_cipher;type:text" json:"-"`
	HandlerDeptID           *string    `gorm:"column:handler_dept_id;size:64" json:"handler_dept_id"`
	Provider                *string    `gorm:"column:provider;size:32" json:"provider"`
	ProviderOrderNo         *string    `gorm:"column:provider_order_no;size:128" json:"provider_order_no"`
	SubmitterID             *string    `gorm:"column:submitter_id;size:64" json:"submitter_id"`
	ReviewerID              *string    `gorm:"column:reviewer_id;size:64" json:"reviewer_id"`
	ReviewedAt              *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	RejectReason            *string    `gorm:"column:reject_reason;size:512" json:"reject_reason"`
	ExpireAt                *time.Time `gorm:"column:expire_at" json:"expire_at"`
	CreatedAt               time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy               *string    `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt               time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy               *string    `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (RealNameCase) TableName() string { return "real_name_case" }

// RealNameCaseRecord 实名工单流水，对应表 real_name_case_record。
type RealNameCaseRecord struct {
	RecordID     string    `gorm:"column:record_id;primaryKey;size:64" json:"record_id"`
	CaseID       string    `gorm:"column:case_id;size:64;not null" json:"case_id"`
	AccountID    *string   `gorm:"column:account_id;size:64" json:"account_id"`
	BusinessType string    `gorm:"column:business_type;size:64;not null" json:"business_type"`
	Action       string    `gorm:"column:action;size:32;not null" json:"action"`
	StatusBefore *string   `gorm:"column:status_before;size:32" json:"status_before"`
	StatusAfter  *string   `gorm:"column:status_after;size:32" json:"status_after"`
	VerifyChannel *string  `gorm:"column:verify_channel;size:32" json:"verify_channel"`
	Provider     *string   `gorm:"column:provider;size:32" json:"provider"`
	OperatorID   *string   `gorm:"column:operator_id;size:64" json:"operator_id"`
	DeptID       *string   `gorm:"column:dept_id;size:64" json:"dept_id"`
	Remark       *string   `gorm:"column:remark;size:512" json:"remark"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 返回表名。
func (RealNameCaseRecord) TableName() string { return "real_name_case_record" }

// StringList JSON 字符串数组列（TEXT/JSON）。
type StringList []string

// Value 实现 driver.Valuer。
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

// Scan 实现 sql.Scanner。
func (s *StringList) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		*s = nil
		return nil
	}
	if len(raw) == 0 {
		*s = StringList{}
		return nil
	}
	return json.Unmarshal(raw, s)
}
