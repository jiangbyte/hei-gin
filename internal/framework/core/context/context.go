// Package contextx 在 context.Context 上挂载请求 ID、会话与账号等请求级键。
//
// Author: Charlie
package contextx

import (
	"context"

	"hei-gin/internal/framework/core/security"
)

type ctxKey int

const (
	keyRequestID ctxKey = iota
	keySession
	keyAccountID
	keyAccountType
	keyClientIP
)

// WithRequestID 写入请求追踪 ID。
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID 读取请求追踪 ID。
func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestID).(string)
	return v
}

// WithSession 写入会话载荷。
func WithSession(ctx context.Context, s *security.SessionPayload) context.Context {
	return context.WithValue(ctx, keySession, s)
}

// Session 读取会话载荷，未登录时返回 nil。
func Session(ctx context.Context) *security.SessionPayload {
	v, _ := ctx.Value(keySession).(*security.SessionPayload)
	return v
}

// WithAccount 写入账号 ID 与账号类型。
func WithAccount(ctx context.Context, id string, t security.AccountType) context.Context {
	ctx = context.WithValue(ctx, keyAccountID, id)
	return context.WithValue(ctx, keyAccountType, t)
}

// AccountID 读取当前账号 ID。
func AccountID(ctx context.Context) string {
	v, _ := ctx.Value(keyAccountID).(string)
	return v
}

// AccountType 读取当前账号类型。
func AccountType(ctx context.Context) security.AccountType {
	v, _ := ctx.Value(keyAccountType).(security.AccountType)
	return v
}

// WithClientIP 写入客户端 IP。
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, keyClientIP, ip)
}

// ClientIP 读取客户端 IP。
func ClientIP(ctx context.Context) string {
	v, _ := ctx.Value(keyClientIP).(string)
	return v
}
