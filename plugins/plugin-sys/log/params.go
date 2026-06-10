package log

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
