// Package logger 提供基于 zap 的进程级全局日志。
//
// Author: Charlie
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L 为全局 Logger，Setup 前为 Nop。
var L *zap.Logger = zap.NewNop()

// Setup 按 debug 开关初始化开发 / 生产 zap 配置。
func Setup(debug bool) error {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	lg, err := cfg.Build()
	if err != nil {
		return err
	}
	L = lg
	return nil
}

// Sync 刷新日志缓冲。
func Sync() { _ = L.Sync() }
