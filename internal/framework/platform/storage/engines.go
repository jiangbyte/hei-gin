// internal/framework/platform/storage/engines.go 文件引擎与 provider 映射（对齐 hei-boot FileEngines / hei-fastapi engines.py）。
//
// Author: Charlie

package storage

import "strings"

// Provider 标识（sys_file.storage_provider / 前端 STORAGE_PROVIDER_OPTIONS）。
const (
	ProviderMinIO  = "minio"
	ProviderRustFS = "rustfs"
	ProviderOSS    = "oss"
	ProviderS3     = "s3"
)

// Engine 标识（DEFAULT_FILE_ENGINE / 配置页）。
const (
	EngineMinIO   = "MINIO"
	EngineRustFS  = "RUSTFS"
	EngineAliyun  = "ALIYUN"
	EngineTencent = "TENCENT"
)

var engineToProvider = map[string]string{
	EngineMinIO:   ProviderMinIO,
	EngineRustFS:  ProviderRustFS,
	EngineAliyun:  ProviderOSS,
	EngineTencent: ProviderS3,
}

var providerToEngine = map[string]string{
	ProviderMinIO:  EngineMinIO,
	ProviderRustFS: EngineRustFS,
	ProviderOSS:    EngineAliyun,
	ProviderS3:     EngineTencent,
}

var providerToKeyPrefix = map[string]string{
	ProviderMinIO:  "STORAGE_MINIO",
	ProviderRustFS: "STORAGE_RUSTFS",
	ProviderOSS:    "STORAGE_ALIYUN",
	ProviderS3:     "STORAGE_TENCENT",
}

// EngineToProvider DEFAULT_FILE_ENGINE → provider。
func EngineToProvider(engine string) string {
	return engineToProvider[strings.ToUpper(strings.TrimSpace(engine))]
}

// ProviderToEngine provider → DEFAULT_FILE_ENGINE。
func ProviderToEngine(provider string) string {
	return providerToEngine[strings.ToLower(strings.TrimSpace(provider))]
}

// ResolveProvider 接受引擎名或 provider 名，返回规范化 provider；未知返回空。
func ResolveProvider(engineOrProvider string) string {
	v := strings.TrimSpace(engineOrProvider)
	if v == "" {
		return ""
	}
	lower := strings.ToLower(v)
	if _, ok := providerToEngine[lower]; ok {
		return lower
	}
	return EngineToProvider(v)
}

// ProviderConfigKeyPrefix provider → sys_config 键前缀。
func ProviderConfigKeyPrefix(provider string) string {
	return providerToKeyPrefix[ResolveProvider(provider)]
}

// IsS3Compatible 是否为支持的对象存储 provider。
func IsS3Compatible(provider string) bool {
	_, ok := providerToEngine[ResolveProvider(provider)]
	return ok
}
