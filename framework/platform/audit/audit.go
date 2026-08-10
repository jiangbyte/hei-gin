// Package audit 提供操作审计入队与异步落库。
package audit

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
)

// Event 为操作审计载荷（入队后落库）。
//
// Author: Charlie
type Event struct {
	Module       string         `json:"module"`
	ResourceType string         `json:"resource_type"`
	ResourceID   string         `json:"resource_id"`
	Action       string         `json:"action"`
	Summary      string         `json:"summary"`
	BeforeData   any            `json:"before_data"`
	AfterData    any            `json:"after_data"`
	AccountID    string         `json:"account_id"`
	AccountType  string         `json:"account_type"`
	RequestID    string         `json:"request_id"`
	IP           string         `json:"ip"`
	UserAgent    string         `json:"user_agent"`
	Success      bool           `json:"success"`
	ErrorMessage string         `json:"error_message"`
	Extra        map[string]any `json:"extra"`
}

// LogRow 对应 sys_operation_audit_log 表。
//
// Author: Charlie
type LogRow struct {
	ID           string          `gorm:"column:id;primaryKey;size:64"`
	Module       string          `gorm:"column:module;size:64"`
	ResourceType *string         `gorm:"column:resource_type;size:128"`
	ResourceID   *string         `gorm:"column:resource_id;size:128"`
	Action       string          `gorm:"column:action;size:64"`
	Summary      *string         `gorm:"column:summary;size:255"`
	BeforeData   json.RawMessage `gorm:"column:before_data;type:jsonb"`
	AfterData    json.RawMessage `gorm:"column:after_data;type:jsonb"`
	AccountID    *string         `gorm:"column:account_id;size:64"`
	AccountType  *string         `gorm:"column:account_type;size:32"`
	RequestID    *string         `gorm:"column:request_id;size:64"`
	IP           *string         `gorm:"column:ip;size:64"`
	UserAgent    *string         `gorm:"column:user_agent;size:512"`
	Success      bool            `gorm:"column:success"`
	ErrorMessage *string         `gorm:"column:error_message"`
	CreatedAt    time.Time       `gorm:"column:created_at;autoCreateTime"`
}

// TableName 返回审计日志表名。
func (LogRow) TableName() string { return "sys_operation_audit_log" }

// Queue 为进程内审计队列（后续可换 Redis/MQ）。
//
// Author: Charlie
type Queue struct {
	ch     chan Event
	db     *gorm.DB
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// NewQueue 创建指定容量的审计队列。
func NewQueue(db *gorm.DB, size int) *Queue {
	if size <= 0 {
		size = 1000
	}
	return &Queue{ch: make(chan Event, size), db: db}
}

// Start 启动异步落库消费者。
func (q *Queue) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	q.cancel = cancel
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-q.ch:
				_ = q.persist(ctx, ev)
			}
		}
	}()
}

// Stop 停止消费者并等待退出。
func (q *Queue) Stop() {
	if q.cancel != nil {
		q.cancel()
	}
	q.wg.Wait()
}

// Publish 非阻塞入队；队列满时丢弃事件。
func (q *Queue) Publish(ev Event) {
	select {
	case q.ch <- ev:
	default:
		// 压力下丢弃，避免阻塞业务请求
	}
}

func (q *Queue) persist(ctx context.Context, ev Event) error {
	before, _ := json.Marshal(ev.BeforeData)
	after, _ := json.Marshal(ev.AfterData)
	row := LogRow{
		ID:        idgen.Next(),
		Module:    ev.Module,
		Action:    ev.Action,
		Success:   ev.Success,
		CreatedAt: time.Now().UTC(),
	}
	if ev.ResourceType != "" {
		row.ResourceType = &ev.ResourceType
	}
	if ev.ResourceID != "" {
		row.ResourceID = &ev.ResourceID
	}
	if ev.Summary != "" {
		row.Summary = &ev.Summary
	}
	if ev.AccountID != "" {
		row.AccountID = &ev.AccountID
	}
	if ev.AccountType != "" {
		row.AccountType = &ev.AccountType
	}
	if ev.RequestID != "" {
		row.RequestID = &ev.RequestID
	}
	if ev.IP != "" {
		row.IP = &ev.IP
	}
	if ev.UserAgent != "" {
		row.UserAgent = &ev.UserAgent
	}
	if ev.ErrorMessage != "" {
		row.ErrorMessage = &ev.ErrorMessage
	}
	if len(before) > 0 && string(before) != "null" {
		row.BeforeData = before
	}
	if len(after) > 0 && string(after) != "null" {
		row.AfterData = after
	}
	return q.db.WithContext(ctx).Create(&row).Error
}
