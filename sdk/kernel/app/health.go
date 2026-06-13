package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/sdk/config"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/kernel/plugin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": config.C.App.Name + " is running",
		"version": config.C.App.Version,
	})
}

func LiveHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": config.C.App.Name,
		"version": config.C.App.Version,
	})
}

func ReadyHandler(c *gin.Context) {
	report := readinessReport(c.Request.Context())
	status := http.StatusOK
	if !report.Ready {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, report)
}

type readinessState struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
}

type pluginStatusBrief struct {
	Name      string `json:"name"`
	InitOK    bool   `json:"init_ok"`
	StartOK   bool   `json:"start_ok"`
	LastError string `json:"last_error,omitempty"`
}

type readinessPayload struct {
	Ready      bool                `json:"ready"`
	Service    string              `json:"service"`
	Version    string              `json:"version"`
	CheckedAt  string              `json:"checked_at"`
	Components []readinessState    `json:"components"`
	Plugins    []pluginStatusBrief `json:"plugins"`
}

func readinessReport(parent context.Context) readinessPayload {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	components := []readinessState{
		checkMySQL(ctx),
		checkRedis(ctx),
	}
	pluginsReady, snapshot := plugin.Ready()
	pluginItems := make([]pluginStatusBrief, 0, len(snapshot))
	for _, item := range snapshot {
		pluginItems = append(pluginItems, pluginStatusBrief{
			Name:      item.Name,
			InitOK:    item.InitOK,
			StartOK:   item.StartOK,
			LastError: item.LastError,
		})
	}

	ready := pluginsReady
	for _, item := range components {
		if !item.OK {
			ready = false
			break
		}
	}

	return readinessPayload{
		Ready:      ready,
		Service:    config.C.App.Name,
		Version:    config.C.App.Version,
		CheckedAt:  time.Now().Format(time.RFC3339),
		Components: components,
		Plugins:    pluginItems,
	}
}

func checkMySQL(ctx context.Context) readinessState {
	if db.DB == nil {
		return readinessState{Name: "mysql", OK: false, Detail: "not initialized"}
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		return readinessState{Name: "mysql", OK: false, Detail: err.Error()}
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return readinessState{Name: "mysql", OK: false, Detail: err.Error()}
	}
	return readinessState{Name: "mysql", OK: true}
}

func checkRedis(ctx context.Context) readinessState {
	if db.Redis == nil {
		return readinessState{Name: "redis", OK: false, Detail: "not initialized"}
	}
	if err := db.Redis.Ping(ctx).Err(); err != nil {
		return readinessState{Name: "redis", OK: false, Detail: err.Error()}
	}
	return readinessState{Name: "redis", OK: true}
}
