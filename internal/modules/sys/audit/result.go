// internal/modules/sys/audit/result.go API 出参（wire 字符串标量，对齐 hei-boot）。
//
// Author: Charlie

package audit

import (
	"time"

	"hei-gin/internal/framework/core/schema"
)

// OperationLogResult 操作审计日志 API 响应。
type OperationLogResult struct {
	ID           string          `json:"id"`
	Module       string          `json:"module"`
	ResourceType *string         `json:"resource_type"`
	ResourceID   *string         `json:"resource_id"`
	Action       string          `json:"action"`
	Summary      *string         `json:"summary"`
	BeforeData   schema.JSONOrNull `json:"before_data"`
	AfterData    schema.JSONOrNull `json:"after_data"`
	AccountID    *string         `json:"account_id"`
	AccountType  *string         `json:"account_type"`
	RequestID    *string         `json:"request_id"`
	IP           *string         `json:"ip"`
	UserAgent    *string         `json:"user_agent"`
	ActionName   *string         `json:"action_name"`
	ActionType   *string         `json:"action_type"`
	ModuleLabel  *string         `json:"module_label"`
	OperatorName *string         `json:"operator_name"`
	DurationMs   *string         `json:"duration_ms"`
	Success      schema.WireBool `json:"success"`
	ErrorMessage *string         `json:"error_message"`
	CreatedAt    time.Time       `json:"created_at"`
	CreatedBy    *string         `json:"created_by"`
	UpdatedBy    *string         `json:"updated_by"`
}

func toOperationLogResult(row OperationLog) OperationLogResult {
	return OperationLogResult{
		ID:           row.ID,
		Module:       row.Module,
		ResourceType: schema.StringPtr(row.ResourceType),
		ResourceID:   schema.StringPtr(row.ResourceID),
		Action:       row.Action,
		Summary:      schema.StringPtr(row.Summary),
		BeforeData:   schema.JSONOrNullFromBytes(row.BeforeData),
		AfterData:    schema.JSONOrNullFromBytes(row.AfterData),
		AccountID:    schema.StringPtr(row.AccountID),
		AccountType:  schema.StringPtr(row.AccountType),
		RequestID:    schema.StringPtr(row.RequestID),
		IP:           schema.StringPtr(row.IP),
		UserAgent:    schema.StringPtr(row.UserAgent),
		ActionName:   schema.StringPtr(row.ActionName),
		ActionType:   schema.StringPtr(row.ActionType),
		ModuleLabel:  schema.StringPtr(row.ModuleLabel),
		OperatorName: schema.StringPtr(row.OperatorName),
		DurationMs:   schema.IntStringPtr(row.DurationMs),
		Success:      schema.WireBoolValue(row.Success),
		ErrorMessage: row.ErrorMessage,
		CreatedAt:    row.CreatedAt,
		CreatedBy:    row.CreatedBy,
		UpdatedBy:    row.UpdatedBy,
	}
}

func toOperationLogResults(rows []OperationLog) []OperationLogResult {
	out := make([]OperationLogResult, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOperationLogResult(row))
	}
	return out
}
