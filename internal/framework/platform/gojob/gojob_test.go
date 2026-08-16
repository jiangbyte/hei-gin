package gojob

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	glebarezsqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(glebarezsqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return db
}

// TestManager_TriggerRunsAndLogs 验证：注册 handler → SyncJob → Trigger 立即执行
// → go-job 执行 executor → 写 sys_job_log + 更新 last_run_at。
func TestManager_TriggerRunsAndLogs(t *testing.T) {
	db := newTestDB(t)
	var mu sync.Mutex
	var ran bool
	m := NewManager(db, []HandlerDef{
		{
			Key:  "testJob",
			Name: "Test Job",
			Run: func(ctx context.Context, param string) error {
				mu.Lock()
				ran = true
				mu.Unlock()
				if param == "boom" {
					return errors.New("boom")
				}
				return nil
			},
		},
	})
	ctx := context.Background()
	if err := m.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer func() { _ = m.Stop(ctx) }()

	row := SysJob{
		ID:         "j1",
		HandlerKey: "testJob",
		Name:       "Test Job",
		CronExpr:   "0 * * * * *",
		Params:     "",
		Status:     StatusEnabled,
	}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create job: %v", err)
	}
	if err := m.SyncJob(ctx, &row); err != nil {
		t.Fatalf("sync job: %v", err)
	}
	if err := m.Trigger(ctx, "testJob", ""); err != nil {
		t.Fatalf("trigger: %v", err)
	}

	// 等待 worker 执行 + 日志写入
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		var logs []SysJobLog
		if err := db.Where("handler_key = ?", "testJob").Find(&logs).Error; err == nil && len(logs) > 0 {
			if logs[0].Status == LogSuccess && logs[0].JobName == "Test Job" && logs[0].DurationMS >= 0 {
				var updated SysJob
				if err := db.First(&updated, "id = ?", "j1").Error; err == nil && updated.LastRunAt != nil {
					mu.Lock()
					ok := ran
					mu.Unlock()
					if ok {
						return // success
					}
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("expected job to run and write success log + last_run_at")
}

// TestManager_Handlers 验证处理器表填充与排序无关的集合查询。
func TestManager_Handlers(t *testing.T) {
	db := newTestDB(t)
	m := NewManager(db, []HandlerDef{{Key: "a", Name: "A", Run: func(ctx context.Context, p string) error { return nil }}})
	m.SetHandlers([]HandlerDef{{Key: "b", Name: "B", Run: func(ctx context.Context, p string) error { return nil }}})
	found := map[string]bool{}
	for _, h := range m.Handlers() {
		found[h.Key] = true
	}
	if !found["a"] || !found["b"] {
		t.Fatalf("handlers missing: %v", found)
	}
}

// TestManager_InvalidCron 验证非法 cron 在 SyncJob 时报错（不会 panic）。
func TestManager_InvalidCron(t *testing.T) {
	db := newTestDB(t)
	m := NewManager(db, []HandlerDef{{Key: "x", Name: "X", Run: func(ctx context.Context, p string) error { return nil }}})
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer func() { _ = m.Stop(context.Background()) }()
	row := SysJob{ID: "x1", HandlerKey: "x", Name: "X", CronExpr: "not-a-cron", Status: StatusEnabled}
	if err := m.SyncJob(context.Background(), &row); err == nil {
		t.Fatalf("expected invalid cron error")
	} else if !strings.Contains(err.Error(), "invalid cron") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestManager_StartWithExistingJobs å›žå½’ï¼šStart å‰å·²æœ‰ ENABLED ä»»åŠ¡æ—¶ï¼Œ
// ä¸å¾—å› é”åµŒå¥—æ­»é”ï¼Œä¸”åº”å®Œæˆ cron æ³¨å†Œå¹¶æ›´æ–° next_run_atã€‚
func TestManager_StartWithExistingJobs(t *testing.T) {
	db := newTestDB(t)
	m := NewManager(db, []HandlerDef{{
		Key:  "seedJob",
		Name: "Seed",
		Run:  func(ctx context.Context, p string) error { return nil },
	}})
	ctx := context.Background()
	if err := db.AutoMigrate(&SysJob{}, &SysJobLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	row := SysJob{ID: "seed1", HandlerKey: "seedJob", Name: "Seed", CronExpr: "0 0 * * * *", Status: StatusEnabled}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- m.Start(ctx) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("start: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Start deadlocked (lock nesting)")
	}
	defer func() { _ = m.Stop(ctx) }()

	m.mu.RLock()
	_, scheduled := m.entries["seedJob"]
	m.mu.RUnlock()
	if !scheduled {
		t.Fatalf("seedJob not scheduled after Start")
	}
	var updated SysJob
	if err := db.First(&updated, "id = ?", "seed1").Error; err != nil {
		t.Fatalf("load updated: %v", err)
	}
	if updated.NextRunAt == nil {
		t.Fatalf("next_run_at not set after Start")
	}
}


// TestHandlerInfoJSONMarshal 回归：处理器清单接口必须返回可 JSON 序列化的
// 轻量 DTO（HandlerDef 含 func 字段，直接序列化会报 unsupported type → handlers 接口 500）。
func TestHandlerInfoJSONMarshal(t *testing.T) {
	m := NewManager(nil, []HandlerDef{{
		Key:  "a",
		Name: "A",
		Run:  func(ctx context.Context, p string) error { return nil },
	}})
	infos := m.HandlerInfos()
	b, err := json.Marshal(infos)
	if err != nil {
		t.Fatalf("marshal err: %v", err)
	}
	if len(infos) != 1 || infos[0].Key != "a" {
		t.Fatalf("unexpected infos: %+v", infos)
	}
	t.Logf("bytes=%s", string(b))
}

