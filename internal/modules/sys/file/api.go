// Package file 的跨模块门面（对齐 hei-boot FileApi）。
//
// Author: Charlie

package file

import (
	"context"
	"io"
	"mime/multipart"

	"hei-gin/internal/framework/platform/module"
)

// ServiceKey Deps 服务袋键：注册 *Service 单例供 profile 等消费。
const ServiceKey = "sys_file_service"

// API 跨模块文件能力（薄门面；实现为 *Service）。
//
// Author: Charlie
type API interface {
	Upload(ctx context.Context, fh *multipart.FileHeader, storageProvider, accountID string) (*File, error)
	CreateFromStream(ctx context.Context, objectName, originalName, storageProvider, contentType string, size int64, r io.Reader, accountID string) (*File, error)
	DeleteByObjectName(ctx context.Context, objectName string) error
}

// FromDeps 从依赖袋取出文件服务；缺失时按 Deps 新建并注册（保证同进程单例）。
func FromDeps(d *module.Deps) *Service {
	if d == nil {
		return nil
	}
	if v, ok := d.Service(ServiceKey); ok {
		if s, ok := v.(*Service); ok && s != nil {
			return s
		}
	}
	s := NewService(d.DB, d.Storage, d.Runtime)
	d.Provide(ServiceKey, s)
	return s
}
