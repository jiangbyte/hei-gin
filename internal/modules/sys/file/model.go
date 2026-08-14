// Package file 提供系统文件上传、下载与元数据管理。
//
// Author: Charlie
package file

import "time"

// File 文件元数据实体，对应表 sys_file。
//
// Author: Charlie
type File struct {
	ID              string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	ObjectName      string    `gorm:"column:object_name;size:255;uniqueIndex;not null" json:"object_name"`
	OriginalName    string    `gorm:"column:original_name;size:255;not null" json:"original_name"`
	StorageProvider string    `gorm:"column:storage_provider;size:32;not null" json:"storage_provider"`
	Bucket          *string   `gorm:"column:bucket;size:255" json:"bucket"`
	ContentType     string    `gorm:"column:content_type;size:128;not null" json:"content_type"`
	Size            int64     `gorm:"column:size;not null" json:"size"`
	URL             string    `gorm:"column:url;size:1024;not null" json:"url"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy       *string   `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy       *string   `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回 File 对应的数据库表名。
func (File) TableName() string { return "sys_file" }
