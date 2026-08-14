// Package dict 提供系统数据字典管理（CRUD、分页与树查询）。
package dict

import "time"

// Dict 数据字典实体，对应表 sys_dict。
//
// Author: Charlie
type Dict struct {
	ID        string    `gorm:"column:id;primaryKey;size:32" json:"id"`
	Code      string    `gorm:"column:code;size:50;uniqueIndex;not null" json:"code"`
	Label     *string   `gorm:"column:label;size:255" json:"label"`
	Value     *string   `gorm:"column:value;size:255" json:"value"`
	Color     *string   `gorm:"column:color;size:32" json:"color"`
	Category  *string   `gorm:"column:category;size:64" json:"category"`
	ParentID  *string   `gorm:"column:parent_id;size:32" json:"parent_id"`
	Status    string    `gorm:"column:status;size:16;not null" json:"status"`
	Sort      int       `gorm:"column:sort;not null;default:0" json:"sort"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 Dict 对应的数据库表名。
func (Dict) TableName() string { return "sys_dict" }
