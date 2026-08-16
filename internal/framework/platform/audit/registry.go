// internal/framework/platform/audit/registry.go 操作审计路由注册表（对齐 hei-boot @OperationAudit 声明式覆盖）。
//
// 各业务模块在 registerRoutes 中与路由并列登记待审计的写接口，
// middleware.Audit 在请求成功后按当前 method+path 命中注册表并发布审计事件。
//
// Author: Charlie

package audit

import (
	"path"
	"strings"
	"sync"
)

// AuditSpec 描述一个待审计接口（镜像 hei-boot @OperationAudit(resourceType, action)）。
//
// Author: Charlie
type AuditSpec struct {
	Method       string // HTTP 方法（大写，如 POST）
	PathPattern  string // 路由路径（含 /api 前缀，路径参数用 *，如 /v1/admin/oauth/*/bind/authorize）
	ResourceType string // 资源类型（对齐 hei-boot 注解值，如 sys_banner、iam_account）
	Action       string // 动作（对齐 hei-boot 注解值，如 create、update、delete）
}

// Registry 审计路由注册表（线程安全）。
//
// Author: Charlie
type Registry struct {
	mu    sync.RWMutex
	specs []AuditSpec
}

// NewRegistry 创建空注册表。
//
// Author: Charlie
func NewRegistry() *Registry {
	return &Registry{}
}

// Register 登记一个待审计接口。
//
// Author: Charlie
func (r *Registry) Register(method, pathPattern, resourceType, action string) {
	if r == nil {
		return
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" || pathPattern == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.specs = append(r.specs, AuditSpec{
		Method:       method,
		PathPattern:  pathPattern,
		ResourceType: resourceType,
		Action:       action,
	})
}

// RegisterSpecs 批量登记审计规格。
//
// Author: Charlie
func (r *Registry) RegisterSpecs(specs ...AuditSpec) {
	for _, s := range specs {
		r.Register(s.Method, s.PathPattern, s.ResourceType, s.Action)
	}
}

// Match 按 method+path 命中审计规格（首个匹配返回）。
//
// Author: Charlie
func (r *Registry) Match(method, requestPath string) (AuditSpec, bool) {
	if r == nil {
		return AuditSpec{}, false
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.specs {
		if s.Method != method {
			continue
		}
		if ok, _ := path.Match(s.PathPattern, requestPath); ok {
			return s, true
		}
	}
	return AuditSpec{}, false
}

// All 返回当前全部审计规格（测试与排查用）。
//
// Author: Charlie
func (r *Registry) All() []AuditSpec {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]AuditSpec, len(r.specs))
	copy(out, r.specs)
	return out
}

// BuildModule 对齐 hei-boot AuditServiceImpl.buildModule：
// resourceType=="resources" → "resource"，其余一律 "iam"。
//
// Author: Charlie
func BuildModule(resourceType string) string {
	if resourceType == "resources" {
		return "resource"
	}
	return "iam"
}
