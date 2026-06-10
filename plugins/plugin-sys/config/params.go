package config

// ConfigVO 配置视图对象
type ConfigVO struct {
	ID          string  `json:"id"`
	ConfigKey   *string `json:"config_key"`
	ConfigValue *string `json:"config_value"`
	Category    *string `json:"category"`
	Remark      *string `json:"remark"`
	SortCode    int     `json:"sort_code"`
	Extra       *string `json:"extra"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   string  `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
}

// ConfigPageParam 配置分页参数
type ConfigPageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Category string `json:"category" form:"category"`
	Keyword  string `json:"keyword" form:"keyword"`
}

// ConfigListParam 配置列表查询参数
type ConfigListParam struct {
	Category string `json:"category" form:"category"`
}

// ConfigBatchEditItem 批量编辑配置项
type ConfigBatchEditItem struct {
	ID          string  `json:"id"`
	ConfigKey   *string `json:"config_key"`
	ConfigValue *string `json:"config_value"`
	Remark      *string `json:"remark"`
	SortCode    int     `json:"sort_code"`
}

// ConfigBatchEditParam 批量编辑配置参数
type ConfigBatchEditParam struct {
	Configs []ConfigBatchEditItem `json:"configs"`
}

// ConfigCategoryEditParam 按分类编辑配置参数
type ConfigCategoryEditParam struct {
	Category    string  `json:"category"`
	ConfigKey   *string `json:"config_key"`
	ConfigValue *string `json:"config_value"`
	Remark      *string `json:"remark"`
}
