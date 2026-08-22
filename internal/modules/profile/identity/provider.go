// internal/modules/profile/identity/provider.go 第三方实名认证 Provider。
//
// Author: Charlie
package identity

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// VerifyProvider 第三方实名认证插件。
type VerifyProvider interface {
	ProviderCode() string
	Supports(verifyChannel, documentType string) bool
	InitVerify(caseEntity *RealNameCase, param RealNameCaseInitThirdPartyParam) RealNameCaseInitResult
	HandleCallback(caseEntity *RealNameCase, param RealNameCaseCallbackParam)
}

// ProviderRegistry 按 provider code 或通道路由 Provider。
type ProviderRegistry struct {
	providers []VerifyProvider
}

// NewProviderRegistry 构造注册表。
func NewProviderRegistry(providers ...VerifyProvider) *ProviderRegistry {
	return &ProviderRegistry{providers: providers}
}

// Resolve 解析 Provider。
func (r *ProviderRegistry) Resolve(verifyChannel, documentType, preferredProvider string) (VerifyProvider, error) {
	if strings.TrimSpace(preferredProvider) != "" {
		for _, p := range r.providers {
			if strings.EqualFold(preferredProvider, p.ProviderCode()) {
				return p, nil
			}
		}
		return nil, bizErr(400, 400, "Unsupported identity provider: "+preferredProvider)
	}
	for _, p := range r.providers {
		if p.ProviderCode() == ProviderMock {
			continue
		}
		if p.Supports(verifyChannel, documentType) {
			return p, nil
		}
	}
	for _, p := range r.providers {
		if p.Supports(verifyChannel, documentType) {
			return p, nil
		}
	}
	return nil, bizErr(400, 400, fmt.Sprintf("No identity provider for channel=%s", verifyChannel))
}

// MockVerifyProvider 开发/测试用 Mock Provider。
type MockVerifyProvider struct{}

// ProviderCode 返回 MOCK。
func (MockVerifyProvider) ProviderCode() string { return ProviderMock }

// Supports 是否支持第三方通道。
func (MockVerifyProvider) Supports(verifyChannel, _ string) bool {
	return strings.EqualFold(verifyChannel, ChannelThirdParty)
}

// InitVerify 初始化 Mock 认证。
func (MockVerifyProvider) InitVerify(caseEntity *RealNameCase, _ RealNameCaseInitThirdPartyParam) RealNameCaseInitResult {
	return RealNameCaseInitResult{
		CaseID:          caseEntity.CaseID,
		Provider:        ProviderMock,
		ProviderOrderNo: "MOCK-" + strings.ReplaceAll(uuid.NewString(), "-", ""),
		RedirectURL:     "/mock/identity-verify?case_id=" + caseEntity.CaseID,
	}
}

// HandleCallback Mock 回调无额外逻辑。
func (MockVerifyProvider) HandleCallback(_ *RealNameCase, _ RealNameCaseCallbackParam) {}
