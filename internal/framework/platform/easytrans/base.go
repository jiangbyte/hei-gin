// Package easytrans 对齐 hei-boot easy-trans 出参（trans_map 与 ID 名称回填）。
//
// Author: Charlie
package easytrans

// Map easy-trans 翻译映射（boot 默认返回空对象 {}）。
type Map map[string]any

// EmptyMap 返回空 trans_map。
func EmptyMap() Map {
	return Map{}
}

// Base 可嵌入出参结构体。
type Base struct {
	TransMap Map `json:"trans_map"`
}

// NewBase 构造带空 trans_map 的基础字段。
func NewBase() Base {
	return Base{TransMap: EmptyMap()}
}
