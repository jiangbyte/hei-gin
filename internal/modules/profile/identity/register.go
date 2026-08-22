// Package identity 实名认证模块注册。
//
// Author: Charlie
package identity

import (
	"hei-gin/internal/framework/platform/module"
)

// ServiceKey Deps 服务袋键：注册 *Service 单例供 profile 等消费。
const ServiceKey = "profile_identity_service"

// FromDeps 从依赖袋取出实名认证服务；缺失时按 Deps 新建并注册。
func FromDeps(d *module.Deps) *Service {
	if d == nil {
		return nil
	}
	if v, ok := d.Service(ServiceKey); ok {
		if s, ok := v.(*Service); ok && s != nil {
			return s
		}
	}
	s := mustNewService(d)
	d.Provide(ServiceKey, s)
	return s
}

func mustNewService(d *module.Deps) *Service {
	fallbackKey := ""
	if d.Cfg != nil {
		fallbackKey = d.Cfg.Crypto.FernetKey
	}
	crypto, err := NewFieldCrypto(d.Runtime, fallbackKey)
	if err != nil {
		panic("profile.identity: " + err.Error())
	}
	return NewService(d.DB, crypto, d.Storage)
}

// init 自注册 profile.identity 模块。
func init() {
	module.Register("profile.identity", 69, func(d *module.Deps) module.Module {
		s := FromDeps(d)
		return module.Module{
			Name:   "profile.identity",
			Order:  69,
			Models: []any{&ProfileIdentity{}, &RealNameCase{}, &RealNameCaseRecord{}},
			Routes: []module.RouteRegistrar{s.registerRoutes(d)},
		}
	})
}
