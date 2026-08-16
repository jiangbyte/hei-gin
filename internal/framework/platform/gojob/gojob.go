// Package gojob 内嵌任务调度：基于 go-job（cybergarage/go-job）执行引擎 +
// robfig/cron 触发 + sys_job/sys_job_log 表驱动配置（web/admin 任务管理控制台）。
//
// Author: Charlie
package gojob

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cybergarage/go-job/job"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// 任务状态。
const (
	StatusEnabled  = "ENABLED"
	StatusDisabled = "DISABLED"
)

// 执行日志状态。
const (
	LogSuccess = "SUCCESS"
	LogFail    = "FAIL"
	LogCancel  = "CANCEL"
	LogTimeout = "TIMEOUT"
)

// SysJob 任务定义（sys_job）。
//
// Author: Charlie
type SysJob struct {
	ID          string     `gorm:"column:id;primaryKey;size:64" json:"id"`
	HandlerKey  string     `gorm:"column:handler_key;size:64;uniqueIndex;not null" json:"handler_key"`
	Name        string     `gorm:"column:name;size:128;not null" json:"name"`
	CronExpr    string     `gorm:"column:cron_expr;size:64;not null" json:"cron_expr"`
	Params      string     `gorm:"column:params;type:text" json:"params"`
	Status      string     `gorm:"column:status;size:16;not null;default:ENABLED" json:"status"`
	Description *string    `gorm:"column:description;type:text" json:"description"`
	LastRunAt   *time.Time `gorm:"column:last_run_at" json:"last_run_at"`
	NextRunAt   *time.Time `gorm:"column:next_run_at" json:"next_run_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy   *string    `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy   *string    `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (SysJob) TableName() string { return "sys_job" }

// SysJobLog 任务执行日志（sys_job_log）。
//
// Author: Charlie
type SysJobLog struct {
	ID         string    `gorm:"column:id;primaryKey;size:64" json:"id"`
	JobID      *string   `gorm:"column:job_id;size:64;index" json:"job_id"`
	HandlerKey string    `gorm:"column:handler_key;size:64;index;not null" json:"handler_key"`
	JobName    string    `gorm:"column:job_name;size:128;not null" json:"job_name"`
	Status     string    `gorm:"column:status;size:16;not null" json:"status"`
	Message    *string   `gorm:"column:message;type:text" json:"message"`
	DurationMS int64     `gorm:"column:duration_ms;not null;default:0" json:"duration_ms"`
	StartedAt  time.Time `gorm:"column:started_at;not null" json:"started_at"`
	FinishedAt time.Time `gorm:"column:finished_at;not null" json:"finished_at"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index" json:"created_at"`
}

// TableName 返回表名。
func (SysJobLog) TableName() string { return "sys_job_log" }

// HandlerDef 模块注册的任务处理器（由 module.Job 收集而来）。
//
// Author: Charlie
type HandlerDef struct {
	Key  string
	Name string
	Run  func(ctx context.Context, param string) error
}

// HandlerInfo 处理器轻量描述（供控制台下拉选择；不含 func，可 JSON 序列化）。
//
// Author: Charlie
type HandlerInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Manager 内嵌任务调度器。
//
// Author: Charlie
type Manager struct {
	db       *gorm.DB
	handlers map[string]HandlerDef
	mgr      job.Manager
	cron     *cron.Cron
	mu       sync.RWMutex
	entries  map[string]cron.EntryID // handlerKey -> cron entry
	started  bool
}

// NewManager 从模块注册的 Job Handlers 创建调度器（不启动）。
func NewManager(db *gorm.DB, handlers []HandlerDef) *Manager {
	m := &Manager{
		db:       db,
		handlers: map[string]HandlerDef{},
		entries:  map[string]cron.EntryID{},
	}
	for _, h := range handlers {
		if _, dup := m.handlers[h.Key]; dup {
			continue
		}
		m.handlers[h.Key] = h
	}
	return m
}

// SetHandlers 填充/替换处理器表（模块装配完成后调用，幂等）。
func (m *Manager) SetHandlers(handlers []HandlerDef) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, h := range handlers {
		if _, dup := m.handlers[h.Key]; dup {
			continue
		}
		m.handlers[h.Key] = h
	}
}

// Start 启动 go-job 引擎并按 sys_job 表加载启用的任务调度。
func (m *Manager) Start(ctx context.Context) error {
	// 注意：此处只在短临界内检查/标记 started，不能在持锁时调用
	// scheduleJob/removeEntry/nextRunAt（它们内部会再次加锁，Mutex 不可重入会死锁）。
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	// 兜底建表：db.sql 提供完整结构，此处仅防首次部署缺表导致启动失败。
	_ = m.db.WithContext(ctx).AutoMigrate(&SysJob{}, &SysJobLog{})

	mgr, err := job.NewManager()
	if err != nil {
		return fmt.Errorf("gojob: create manager: %w", err)
	}
	if err := mgr.Start(); err != nil {
		return fmt.Errorf("gojob: start manager: %w", err)
	}
	m.mu.Lock()
	m.mgr = mgr
	m.cron = cron.New(cron.WithSeconds()) // 6 段 cron（含秒）
	m.cron.Start()
	m.started = true
	m.mu.Unlock()

	var jobs []SysJob
	if err := m.db.WithContext(ctx).Where("status = ?", StatusEnabled).Find(&jobs).Error; err != nil {
		return fmt.Errorf("gojob: load jobs: %w", err)
	}
	for i := range jobs {
		if err := m.scheduleJob(ctx, &jobs[i]); err != nil {
			continue // 单任务失败不影响整体启动
		}
	}
	return nil
}

// Stop 停止调度器与执行引擎。
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_ = ctx
	if !m.started {
		return nil
	}
	if m.cron != nil {
		m.cron.Stop()
	}
	if m.mgr != nil {
		_ = m.mgr.Stop()
	}
	m.started = false
	return nil
}

// Handlers 返回全部可用处理器（含执行函数，仅供调度使用）。
func (m *Manager) Handlers() []HandlerDef {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]HandlerDef, 0, len(m.handlers))
	for _, h := range m.handlers {
		out = append(out, h)
	}
	return out
}

// HandlerInfos 返回处理器轻量清单（不含 func，可安全 JSON 序列化给控制台）。
func (m *Manager) HandlerInfos() []HandlerInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]HandlerInfo, 0, len(m.handlers))
	for _, h := range m.handlers {
		out = append(out, HandlerInfo{Key: h.Key, Name: h.Name})
	}
	return out
}

// buildJobDef 由 handler 与任务行构造 go-job Job 定义；executor 内记录执行日志。
func (m *Manager) buildJobDef(row *SysJob) (job.Job, error) {
	h, ok := m.handlers[row.HandlerKey]
	if !ok {
		return nil, fmt.Errorf("handler not registered: %s", row.HandlerKey)
	}
	exec := func(ctx context.Context, p string) error {
		started := time.Now()
		err := h.Run(ctx, p)
		finished := time.Now()
		status := LogSuccess
		message := "ok"
		if err != nil {
			status = LogFail
			message = err.Error()
		}
		m.writeLog(ctx, row, status, message, &started, &finished)
		return err
	}
	// 注意：go-job worker 直接调用 CompleteProcessor/TerminateProcessor（无 nil 检查），
	// 必须显式提供非 nil 处理器，否则执行完成时 panic。
	return job.NewJob(
		job.WithKind(row.HandlerKey),
		job.WithDescription(row.Name),
		job.WithExecutor(exec),
		job.WithCompleteProcessor(func(_ job.Instance, _ []any) {}),
		job.WithTerminateProcessor(func(_ job.Instance, err error) error { return err }),
	)
}

// trigger 触发一次执行（立即入 go-job 队列）。
func (m *Manager) trigger(ctx context.Context, row *SysJob, param string) error {
	if m.mgr == nil {
		return fmt.Errorf("gojob manager not started")
	}
	jobDef, err := m.buildJobDef(row)
	if err != nil {
		return err
	}
	if param == "" {
		param = row.Params
	}
	_, err = m.mgr.ScheduleJob(jobDef, job.WithArguments(param))
	return err
}

// Trigger 手动触发任务（业务模块调用）。
func (m *Manager) Trigger(ctx context.Context, handlerKey, param string) error {
	row, err := m.jobByHandlerKey(ctx, handlerKey)
	if err != nil {
		return err
	}
	if row.Status != StatusEnabled {
		return fmt.Errorf("job %s is disabled", handlerKey)
	}
	return m.trigger(ctx, row, param)
}

// SyncJob 按任务行同步调度（创建/更新/启停后调用）。
func (m *Manager) SyncJob(ctx context.Context, row *SysJob) error {
	return m.scheduleJob(ctx, row)
}

// RemoveJob 移除任务调度（删除时调用）。
func (m *Manager) RemoveJob(handlerKey string) {
	m.removeEntry(handlerKey)
}

// scheduleJob 注册/更新 cron 调度（ENABLED 且 cron 合法时）。
func (m *Manager) scheduleJob(ctx context.Context, row *SysJob) error {
	if _, ok := m.handlers[row.HandlerKey]; !ok {
		return fmt.Errorf("handler not registered: %s", row.HandlerKey)
	}
	m.removeEntry(row.HandlerKey)
	if row.Status != StatusEnabled || row.CronExpr == "" {
		return nil
	}
	entryID, err := m.cron.AddFunc(row.CronExpr, func() {
		row := row
		if err := m.trigger(context.Background(), row, row.Params); err != nil {
			m.writeLog(context.Background(), row, LogFail, "schedule trigger failed: "+err.Error(), nil, nil)
		}
	})
	if err != nil {
		return fmt.Errorf("invalid cron %q for %s: %w", row.CronExpr, row.HandlerKey, err)
	}
	m.mu.Lock()
	m.entries[row.HandlerKey] = entryID
	m.mu.Unlock()

	// 更新下次执行时间
	if nxt := m.nextRunAt(row.HandlerKey); nxt != nil {
		_ = m.db.WithContext(ctx).Model(&SysJob{}).Where("id = ?", row.ID).
			Update("next_run_at", *nxt).Error
	}
	return nil
}

// removeEntry 移除指定任务的 cron 条目。
func (m *Manager) removeEntry(handlerKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if id, ok := m.entries[handlerKey]; ok {
		m.cron.Remove(id)
		delete(m.entries, handlerKey)
	}
}

// nextRunAt 取指定任务 cron 的下次执行时间。
func (m *Manager) nextRunAt(handlerKey string) *time.Time {
	m.mu.RLock()
	entryID, ok := m.entries[handlerKey]
	m.mu.RUnlock()
	if !ok || m.cron == nil {
		return nil
	}
	entry := m.cron.Entry(entryID)
	if entry.Schedule == nil {
		return nil
	}
	nxt := entry.Schedule.Next(time.Now())
	if nxt.IsZero() {
		return nil
	}
	return &nxt
}

// writeLog 写执行日志并更新任务 last_run_at / next_run_at。
func (m *Manager) writeLog(ctx context.Context, row *SysJob, status, message string, startedAt, finishedAt *time.Time) {
	if row == nil {
		return
	}
	now := time.Now()
	st := now
	if startedAt != nil {
		st = *startedAt
	}
	fi := now
	if finishedAt != nil {
		fi = *finishedAt
	}
	jobID := row.ID
	log := SysJobLog{
		ID:         newID(),
		JobID:      &jobID,
		HandlerKey: row.HandlerKey,
		JobName:    row.Name,
		Status:     status,
		Message:    &message,
		DurationMS: fi.Sub(st).Milliseconds(),
		StartedAt:  st,
		FinishedAt: fi,
	}
	_ = m.db.WithContext(ctx).Create(&log).Error
	updates := map[string]any{"last_run_at": fi}
	if nxt := m.nextRunAt(row.HandlerKey); nxt != nil {
		updates["next_run_at"] = *nxt
	}
	_ = m.db.WithContext(ctx).Model(&SysJob{}).Where("id = ?", row.ID).Updates(updates).Error
}

// newID 生成日志主键（时间戳 + 纳秒后缀）。
func newID() string {
	return fmt.Sprintf("%d%06d", time.Now().UnixMilli(), time.Now().Nanosecond()/1000%1000000)
}

// jobByHandlerKey 按 handler_key 查询任务。
func (m *Manager) jobByHandlerKey(ctx context.Context, handlerKey string) (*SysJob, error) {
	var row SysJob
	if err := m.db.WithContext(ctx).Where("handler_key = ?", handlerKey).First(&row).Error; err != nil {
		return nil, fmt.Errorf("job not found: %s", handlerKey)
	}
	return &row, nil
}
