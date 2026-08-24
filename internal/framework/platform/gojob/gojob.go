// Package gojob 内嵌任务调度：DB 扫描 + Redis 锁 + sys_job/sys_job_log（列名与 hei-boot 一致）。
//
// Author: Charlie
package gojob

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/logger"
	"hei-gin/internal/framework/platform/idgen"
)

const (
	TypeCRON  = "CRON"
	TypeFIXED = "FIXED"

	ExecutorSystem = "system"

	maxScanLimit       = 50
	lockExpireSeconds  = 30 * 60
	lockAcquireTimeout = time.Second
	maxResultLength    = 500
)

// SysJob 任务定义（sys_job，对齐 hei-boot SysJob）。
//
// Author: Charlie
type SysJob struct {
	ID            string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	Name          string         `gorm:"column:name;size:128;not null" json:"name"`
	Handler       string         `gorm:"column:handler;size:255;not null" json:"handler"`
	TriggerType   string         `gorm:"column:trigger_type;size:16;not null" json:"trigger_type"`
	TriggerConfig string         `gorm:"column:trigger_config;size:255;not null" json:"trigger_config"`
	Params        datatypes.JSON `gorm:"column:params;type:json" json:"params"`
	LastRunTime   *time.Time     `gorm:"column:last_run_time" json:"last_run_time"`
	NextRunTime   time.Time      `gorm:"column:next_run_time;not null;index:idx_sys_job_enabled_next" json:"next_run_time"`
	LastResult    *string        `gorm:"column:last_result;size:500" json:"last_result"`
	Enabled       bool           `gorm:"column:enabled;not null;default:true;index:idx_sys_job_enabled_next" json:"enabled"`
	Description   *string        `gorm:"column:description;size:500" json:"description"`
	Sort          int            `gorm:"column:sort;not null;default:0" json:"sort"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy     *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy     *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (SysJob) TableName() string { return "sys_job" }

// SysJobLog 任务执行日志（sys_job_log，对齐 hei-boot SysJobLog）。
//
// Author: Charlie
type SysJobLog struct {
	ID         string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	JobID      string         `gorm:"column:job_id;size:64;not null;index" json:"job_id"`
	Params     datatypes.JSON `gorm:"column:params;type:json" json:"params"`
	StartedAt  time.Time      `gorm:"column:started_at;not null" json:"started_at"`
	DurationMs *int64         `gorm:"column:duration_ms" json:"duration_ms"`
	Success    bool           `gorm:"column:success;not null" json:"success"`
	Result     *string        `gorm:"column:result;type:text" json:"result"`
	Executor   *string        `gorm:"column:executor;size:64" json:"executor"`
	IP         *string        `gorm:"column:ip;size:64" json:"ip"`
	ProcessID  *string        `gorm:"column:process_id;size:32" json:"process_id"`
	AppDir     *string        `gorm:"column:app_dir;size:500" json:"app_dir"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy  *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy  *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (SysJobLog) TableName() string { return "sys_job_log" }

// HandlerFunc 任务处理器：paramJSON 为 params 的 JSON 文本（可空），返回结果摘要。
//
// Author: Charlie
type HandlerFunc func(ctx context.Context, paramJSON string) (string, error)

// HandlerDef 注册的处理器。
//
// Author: Charlie
type HandlerDef struct {
	Key  string
	Name string
	Run  HandlerFunc
}

// Config 调度器配置（对齐 hei-fastapi JobSettings）。
//
// Author: Charlie
type Config struct {
	ScanIntervalMS   int
	PoolSize         int
	LogRetentionDays int
	LogBatchSize     int
}

func (c Config) normalized() Config {
	if c.ScanIntervalMS <= 0 {
		c.ScanIntervalMS = 1000
	}
	if c.PoolSize <= 0 {
		c.PoolSize = 4
	}
	if c.LogRetentionDays <= 0 {
		c.LogRetentionDays = 30
	}
	if c.LogBatchSize <= 0 {
		c.LogBatchSize = 1000
	}
	return c
}

// Manager 内嵌任务调度器。
//
// Author: Charlie
type Manager struct {
	db       *gorm.DB
	rdb      *redis.Client
	cfg      Config
	handlers map[string]HandlerDef
	mu       sync.RWMutex
	sem      chan struct{}
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	started  bool
	ip       string
	pid      string
	appDir   string
}

// NewManager 创建调度器（不启动）。
func NewManager(db *gorm.DB, rdb *redis.Client, cfg Config, handlers []HandlerDef) *Manager {
	cfg = cfg.normalized()
	m := &Manager{
		db:       db,
		rdb:      rdb,
		cfg:      cfg,
		handlers: map[string]HandlerDef{},
		sem:      make(chan struct{}, cfg.PoolSize),
		ip:       "127.0.0.1",
		pid:      strconv.Itoa(os.Getpid()),
		appDir:   mustCwd(),
	}
	m.SetHandlers(handlers)
	return m
}

func mustCwd() string {
	d, err := os.Getwd()
	if err != nil {
		return "."
	}
	return filepath.Clean(d)
}

// SetHandlers 填充/替换处理器表。
func (m *Manager) SetHandlers(handlers []HandlerDef) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers = map[string]HandlerDef{}
	for _, h := range handlers {
		if h.Key == "" || h.Run == nil {
			continue
		}
		m.handlers[h.Key] = h
	}
}

// Resolve 按 handler key 解析处理器。
func (m *Manager) Resolve(key string) (HandlerDef, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.handlers[key]
	return h, ok
}

// HasHandler 是否已注册。
func (m *Manager) HasHandler(key string) bool {
	_, ok := m.Resolve(key)
	return ok
}

// HandlerInfos 轻量列表（可选控制台用）。
func (m *Manager) HandlerInfos() []HandlerInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]HandlerInfo, 0, len(m.handlers))
	for _, h := range m.handlers {
		name := h.Name
		if name == "" {
			name = h.Key
		}
		out = append(out, HandlerInfo{Key: h.Key, Name: name})
	}
	return out
}

// HandlerInfo 处理器描述。
//
// Author: Charlie
type HandlerInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Start 启动扫描循环。
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return nil
	}
	runCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.started = true
	m.mu.Unlock()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.scanLoop(runCtx)
	}()
	return nil
}

// Stop 停止扫描循环。
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return nil
	}
	cancel := m.cancel
	m.cancel = nil
	m.started = false
	m.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Manager) scanLoop(ctx context.Context) {
	interval := time.Duration(m.cfg.ScanIntervalMS) * time.Millisecond
	if interval < 100*time.Millisecond {
		interval = 100 * time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.scanOnce(ctx)
		}
	}
}

func (m *Manager) scanOnce(ctx context.Context) {
	var jobs []SysJob
	now := time.Now().UTC()
	err := m.db.WithContext(ctx).
		Where("enabled = ? AND next_run_time <= ?", true, now).
		Order("sort ASC, next_run_time ASC").
		Limit(maxScanLimit).
		Find(&jobs).Error
	if err != nil {
		logger.L.Warn("job scan failed", zap.Error(err))
		return
	}
	for _, j := range jobs {
		m.SubmitRun(ctx, j.ID, false, ExecutorSystem)
	}
}

// SubmitRun 有界并发提交执行（立即返回；执行与写库使用独立 context，避免 HTTP 请求结束导致 cancel）。
func (m *Manager) SubmitRun(_ context.Context, jobID string, force bool, executor string) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.sem <- struct{}{}
		defer func() { <-m.sem }()
		_ = m.RunJob(context.Background(), jobID, force, executor)
	}()
}

// RunJob 同步执行单个任务（带 Redis 锁）。
func (m *Manager) RunJob(ctx context.Context, jobID string, force bool, executor string) error {
	lockKey := "sys:job:run:" + jobID
	if m.rdb != nil {
		ok, err := m.acquireLock(ctx, lockKey)
		if err != nil || !ok {
			return err
		}
		defer func() { _ = m.rdb.Del(context.Background(), lockKey).Err() }()
	}
	return m.runLocked(ctx, jobID, force, executor)
}

func (m *Manager) acquireLock(ctx context.Context, key string) (bool, error) {
	deadline := time.Now().Add(lockAcquireTimeout)
	for {
		ok, err := m.rdb.SetNX(ctx, key, "1", time.Duration(lockExpireSeconds)*time.Second).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (m *Manager) runLocked(ctx context.Context, jobID string, force bool, executor string) error {
	var job SysJob
	if err := m.db.WithContext(ctx).First(&job, "id = ?", jobID).Error; err != nil {
		return err
	}
	if !job.Enabled {
		return nil
	}
	now := time.Now().UTC()
	if !force && job.NextRunTime.After(now) {
		return nil
	}

	executeTime := time.Now().UTC()
	handler, ok := m.Resolve(job.Handler)
	var (
		result  string
		success bool
	)
	if !ok {
		result = "执行失败: 未找到任务处理器: " + job.Handler
		success = false
		logger.L.Error("Job execute failed, handler not found",
			zap.String("id", job.ID),
			zap.String("name", job.Name),
			zap.String("handler", job.Handler),
		)
	} else {
		paramJSON := ""
		if len(job.Params) > 0 && string(job.Params) != "null" {
			paramJSON = string(job.Params)
		}
		out, err := handler.Run(ctx, paramJSON)
		if err != nil {
			result = "执行失败: " + err.Error()
			success = false
			logger.L.Error("Job execute failed",
				zap.String("id", job.ID),
				zap.String("name", job.Name),
				zap.Error(err),
			)
		} else {
			result = out
			success = true
		}
	}
	dur := time.Since(executeTime).Milliseconds()
	if dur < 0 {
		dur = 0
	}
	if err := m.recordRun(ctx, &job, executor, success, result, executeTime, dur); err != nil {
		logger.L.Error("Job record run failed",
			zap.String("id", job.ID),
			zap.String("name", job.Name),
			zap.Error(err),
		)
		return err
	}
	logger.L.Info("Job done",
		zap.String("id", job.ID),
		zap.String("name", job.Name),
		zap.Bool("success", success),
		zap.String("result", truncateResult(result)),
		zap.Int64("duration_ms", dur),
		zap.String("executor", executor),
	)
	return nil
}

func (m *Manager) recordRun(ctx context.Context, job *SysJob, executor string, success bool, result string, executeTime time.Time, durMS int64) error {
	next, err := ComputeNextRunTime(job.TriggerType, job.TriggerConfig, executeTime)
	if err != nil {
		next = executeTime.Add(time.Minute)
	}
	truncated := truncateResult(result)
	lastRun := executeTime
	updates := map[string]any{
		"last_run_time": lastRun,
		"next_run_time": next,
		"last_result":   truncated,
	}
	if err := m.db.WithContext(ctx).Model(&SysJob{}).Where("id = ?", job.ID).Updates(updates).Error; err != nil {
		return err
	}

	ip, pid, appDir := m.ip, m.pid, m.appDir
	exec := executor
	logRow := SysJobLog{
		ID:         idgen.Next(),
		JobID:      job.ID,
		Params:     job.Params,
		StartedAt:  executeTime,
		DurationMs: &durMS,
		Success:    success,
		Result:     &truncated,
		Executor:   &exec,
		IP:         &ip,
		ProcessID:  &pid,
		AppDir:     &appDir,
	}
	return m.db.WithContext(ctx).Create(&logRow).Error
}

func truncateResult(s string) string {
	if len(s) <= maxResultLength {
		return s
	}
	return s[:maxResultLength]
}

// --- cron / FIXED ---

var cronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

// ValidateTrigger 校验触发配置。
func ValidateTrigger(triggerType, triggerConfig string) error {
	triggerType = strings.ToUpper(strings.TrimSpace(triggerType))
	triggerConfig = strings.TrimSpace(triggerConfig)
	switch triggerType {
	case TypeFIXED:
		n, err := strconv.Atoi(triggerConfig)
		if err != nil || n < 1 {
			return errors.New("FIXED 触发配置必须为正整数秒数")
		}
		return nil
	case TypeCRON:
		if _, err := cronParser.Parse(triggerConfig); err != nil {
			return fmt.Errorf("CRON 表达式无效: %s", triggerConfig)
		}
		return nil
	default:
		return fmt.Errorf("不支持的触发类型: %s", triggerType)
	}
}

// ComputeNextRunTime 计算下次执行时间（CRON 秒在首位，对齐 Spring/hei-boot）。
func ComputeNextRunTime(triggerType, triggerConfig string, from time.Time) (time.Time, error) {
	triggerType = strings.ToUpper(strings.TrimSpace(triggerType))
	triggerConfig = strings.TrimSpace(triggerConfig)
	if err := ValidateTrigger(triggerType, triggerConfig); err != nil {
		return time.Time{}, err
	}
	if triggerType == TypeFIXED {
		n, _ := strconv.Atoi(triggerConfig)
		return from.Add(time.Duration(n) * time.Second), nil
	}
	sched, err := cronParser.Parse(triggerConfig)
	if err != nil {
		return time.Time{}, err
	}
	return sched.Next(from), nil
}

// ParamJSON 将 map 编码为 JSON（nil → null/empty）。
func ParamJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(b)
}

// ConfigValues 暴露日志清理默认参数。
func (m *Manager) ConfigValues() Config { return m.cfg }
