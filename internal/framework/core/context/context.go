// Package contextx åœ¨ context.Context ä¸ŠæŒ‚è½½è¯·æ±‚ IDã€ä¼šè¯ä¸Žè´¦å·ç­‰è¯·æ±‚çº§é”®ã€‚
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

// WithRequestID å†™å…¥è¯·æ±‚è¿½è¸ª IDã€‚
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID è¯»å–è¯·æ±‚è¿½è¸ª IDã€‚
func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestID).(string)
	return v
}

// WithSession å†™å…¥ä¼šè¯è½½è·ã€‚
func WithSession(ctx context.Context, s *security.SessionPayload) context.Context {
	return context.WithValue(ctx, keySession, s)
}

// Session è¯»å–ä¼šè¯è½½è·ï¼Œæœªç™»å½•æ—¶è¿”å›ž nilã€‚
func Session(ctx context.Context) *security.SessionPayload {
	v, _ := ctx.Value(keySession).(*security.SessionPayload)
	return v
}

// WithAccount å†™å…¥è´¦å· ID ä¸Žè´¦å·ç±»åž‹ã€‚
func WithAccount(ctx context.Context, id string, t security.AccountType) context.Context {
	ctx = context.WithValue(ctx, keyAccountID, id)
	return context.WithValue(ctx, keyAccountType, t)
}

// AccountID è¯»å–å½“å‰è´¦å· IDã€‚
func AccountID(ctx context.Context) string {
	v, _ := ctx.Value(keyAccountID).(string)
	return v
}

// AccountType è¯»å–å½“å‰è´¦å·ç±»åž‹ã€‚
func AccountType(ctx context.Context) security.AccountType {
	v, _ := ctx.Value(keyAccountType).(security.AccountType)
	return v
}

// WithClientIP å†™å…¥å®¢æˆ·ç«¯ IPã€‚
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, keyClientIP, ip)
}

// ClientIP è¯»å–å®¢æˆ·ç«¯ IPã€‚
func ClientIP(ctx context.Context) string {
	v, _ := ctx.Value(keyClientIP).(string)
	return v
}
