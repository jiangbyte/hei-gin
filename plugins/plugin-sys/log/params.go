package log

import "hei-gin/sdk/utils"

type LogVO struct {
	ID         string `json:"id"`
	Category   string `json:"category"`
	Name       string `json:"name"`
	ExeStatus  string `json:"exe_status"`
	ExeMessage string `json:"exe_message"`
	OpIP       string `json:"op_ip"`
	OpAddress  string `json:"op_address"`
	OpBrowser  string `json:"op_browser"`
	OpOs       string `json:"op_os"`
	ClassName  string `json:"class_name"`
	MethodName string `json:"method_name"`
	ReqMethod  string `json:"req_method"`
	ReqURL     string `json:"req_url"`
	ParamJSON  string `json:"param_json"`
	ResultJSON string `json:"result_json"`
	OpTime     string `json:"op_time"`
	TraceID    string `json:"trace_id"`
	OpUser     string `json:"op_user"`
	SignData   string `json:"sign_data"`
	CreatedAt  string `json:"created_at"`
	CreatedBy  string `json:"created_by"`
	UpdatedAt  string `json:"updated_at"`
	UpdatedBy  string `json:"updated_by"`
}

type LogPageParam struct {
	Current   int    `json:"current" form:"current"`
	Size      int    `json:"size" form:"size"`
	Keyword   string `json:"keyword" form:"keyword"`
	Category  string `json:"category" form:"category"`
	ExeStatus string `json:"exe_status" form:"exe_status"`
}

type LogDeleteByCategoryParam struct {
	Category string `json:"category"`
}

type BarChartData struct {
	Days   []string         `json:"days"`
	Series []CategorySeries `json:"series"`
}

type CategorySeries struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type PieChartData struct {
	Data []CategoryTotal `json:"data"`
}

type CategoryTotal struct {
	Category string `json:"category"`
	Total    int    `json:"total"`
}

func SysLogToLogVO(src *SysLog) *LogVO {
	if src == nil {
		return nil
	}

	dst := &LogVO{}
	dst.ID = src.ID
	dst.Category = src.Category
	dst.Name = src.Name
	dst.ExeStatus = src.ExeStatus
	dst.ExeMessage = src.ExeMessage
	dst.OpIP = src.OpIP
	dst.OpAddress = src.OpAddress
	dst.OpBrowser = src.OpBrowser
	dst.OpOs = src.OpOs
	dst.ClassName = src.ClassName
	dst.MethodName = src.MethodName
	dst.ReqMethod = src.ReqMethod
	dst.ReqURL = src.ReqURL
	dst.ParamJSON = src.ParamJSON
	dst.ResultJSON = src.ResultJSON
	dst.TraceID = src.TraceID
	dst.OpUser = src.OpUser
	dst.SignData = src.SignData
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.OpTime = utils.FormatDateTimePtr(src.OpTime)
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func LogVOToSysLog(src *LogVO) *SysLog {
	if src == nil {
		return nil
	}

	dst := &SysLog{}
	dst.ID = src.ID
	dst.Category = src.Category
	dst.Name = src.Name
	dst.ExeStatus = src.ExeStatus
	dst.ExeMessage = src.ExeMessage
	dst.OpIP = src.OpIP
	dst.OpAddress = src.OpAddress
	dst.OpBrowser = src.OpBrowser
	dst.OpOs = src.OpOs
	dst.ClassName = src.ClassName
	dst.MethodName = src.MethodName
	dst.ReqMethod = src.ReqMethod
	dst.ReqURL = src.ReqURL
	dst.ParamJSON = src.ParamJSON
	dst.ResultJSON = src.ResultJSON
	dst.TraceID = src.TraceID
	dst.OpUser = src.OpUser
	dst.SignData = src.SignData
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.OpTime = utils.ParseDateTimePtr(&src.OpTime)
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
