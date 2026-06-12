package log

import "time"

type SysLog struct {
	ID         string     `gorm:"primaryKey;size:32" json:"id"`
	Category   string     `gorm:"size:255;index" json:"category"`
	Name       string     `gorm:"size:255" json:"name"`
	ExeStatus  string     `gorm:"size:255" json:"exe_status"`
	ExeMessage string     `gorm:"type:text" json:"exe_message"`
	OpIP       string     `gorm:"size:255" json:"op_ip"`
	OpAddress  string     `gorm:"size:255" json:"op_address"`
	OpBrowser  string     `gorm:"size:255" json:"op_browser"`
	OpOs       string     `gorm:"size:255" json:"op_os"`
	ClassName  string     `gorm:"size:255" json:"class_name"`
	MethodName string     `gorm:"size:255" json:"method_name"`
	ReqMethod  string     `gorm:"size:255" json:"req_method"`
	ReqURL     string     `gorm:"type:text" json:"req_url"`
	ParamJSON  string     `gorm:"type:longtext" json:"param_json"`
	ResultJSON string     `gorm:"type:longtext" json:"result_json"`
	OpTime     *time.Time `json:"op_time"`
	TraceID    string     `gorm:"size:64" json:"trace_id"`
	OpUser     string     `gorm:"size:255" json:"op_user"`
	SignData   string     `gorm:"type:longtext" json:"sign_data"`
	CreatedAt  *time.Time `json:"created_at"`
	CreatedBy  string     `gorm:"size:32" json:"created_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  string     `gorm:"size:32" json:"updated_by"`
}

func (SysLog) TableName() string { return "sys_log" }
