// Package audit æä¾›æ“ä½œå®¡è®¡å…¥é˜Ÿä¸Žå¼‚æ­¥è½åº“ï¼ˆoutbox + Redis Stream / è¿›ç¨‹å†…é˜Ÿåˆ—ï¼‰ã€‚
//
// Author: Charlie
package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/config"
	"hei-gin/internal/framework/platform/idgen"
)

// Event ä¸ºæ“ä½œå®¡è®¡è½½è·ï¼ˆå…¥é˜ŸåŽè½åº“ï¼‰ã€‚
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
	OutboxID     string         `json:"outbox_id,omitempty"`
}

// LogRow å¯¹åº” sys_operation_audit_log è¡¨ã€‚
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

// TableName è¿”å›žå®¡è®¡æ—¥å¿—è¡¨åã€‚
func (LogRow) TableName() string { return "sys_operation_audit_log" }

// OutboxRow å¯¹åº” sys_operation_audit_outbox è¡¨ã€‚
//
// Author: Charlie
type OutboxRow struct {
	ID        string     `gorm:"column:id;primaryKey;size:64"`
	Payload   string     `gorm:"column:payload"`
	Status    string     `gorm:"column:status;size:32"`
	Attempts  int        `gorm:"column:attempts"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	ClaimedAt *time.Time `gorm:"column:claimed_at"`
}

// TableName è¿”å›ž outbox è¡¨åã€‚
func (OutboxRow) TableName() string { return "sys_operation_audit_outbox" }

// Queue ä¸ºå®¡è®¡é˜Ÿåˆ—ï¼ˆoutbox + Redis Streamï¼Œæˆ–è¿›ç¨‹å†… channel å›žé€€ï¼‰ã€‚
//
// Author: Charlie
type Queue struct {
	ch           chan Event
	db           *gorm.DB
	rdb          *redis.Client
	cfg          config.AuditConfig
	wg           sync.WaitGroup
	cancel       context.CancelFunc
	consumerName string
}

// NewQueue åˆ›å»ºå®¡è®¡é˜Ÿåˆ—ï¼ˆæ”¯æŒ Redis Stream ä¸Žè¿›ç¨‹å†…å›žé€€ï¼‰ã€‚
//
// Author: Charlie
func NewQueue(db *gorm.DB, rdb *redis.Client, cfg config.AuditConfig) *Queue {
	size := cfg.OperationQueueSize
	if size <= 0 {
		size = 1000
	}
	if cfg.StreamKey == "" {
		cfg.StreamKey = "hei:audit:ops"
	}
	if cfg.ConsumerGroup == "" {
		cfg.ConsumerGroup = "hei-gin-audit"
	}
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	return &Queue{
		ch:           make(chan Event, size),
		db:           db,
		rdb:          rdb,
		cfg:          cfg,
		consumerName: "consumer-" + host,
	}
}

// Start å¯åŠ¨å¼‚æ­¥è½åº“æ¶ˆè´¹è€…ã€‚
//
// Author: Charlie
func (q *Queue) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	q.cancel = cancel
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		if q.useStream() {
			q.consumeStream(ctx)
			return
		}
		q.consumeChannel(ctx)
	}()
}

// Stop åœæ­¢æ¶ˆè´¹è€…å¹¶ç­‰å¾…é€€å‡ºã€‚
//
// Author: Charlie
func (q *Queue) Stop() {
	if q.cancel != nil {
		q.cancel()
	}
	q.wg.Wait()
}

// Publish å†™å…¥ outboxï¼Œå† XADD æˆ–è¿›ç¨‹å†…å…¥é˜Ÿã€‚
//
// Author: Charlie
func (q *Queue) Publish(ev Event) {
	if q == nil || q.db == nil {
		return
	}
	outboxID := idgen.Next()
	ev.OutboxID = outboxID
	payload, err := json.Marshal(ev)
	if err != nil {
		return
	}
	row := OutboxRow{
		ID:        outboxID,
		Payload:   string(payload),
		Status:    "PENDING",
		Attempts:  0,
		CreatedAt: time.Now().UTC(),
	}
	if err := q.db.Create(&row).Error; err != nil {
		return
	}

	if q.useStream() {
		fields := map[string]any{
			"data":      string(payload),
			"outbox_id": outboxID,
			"timestamp": fmt.Sprintf("%d", time.Now().UnixMilli()),
		}
		if err := q.rdb.XAdd(context.Background(), &redis.XAddArgs{
			Stream: q.cfg.StreamKey,
			Values: fields,
		}).Err(); err != nil {
			// Stream å¤±è´¥æ—¶åŒæ­¥è½åº“ï¼Œé¿å…ä»…å…¥ channel å´æ— æ¶ˆè´¹è€…
			if perr := q.persist(context.Background(), ev); perr == nil {
				q.markOutboxDone(context.Background(), outboxID)
			}
		}
		return
	}
	q.enqueue(ev)
}

func (q *Queue) useStream() bool {
	return q.cfg.UseStream && q.rdb != nil
}

func (q *Queue) enqueue(ev Event) {
	select {
	case q.ch <- ev:
	default:
		// åŽ‹åŠ›ä¸‹ä¸¢å¼ƒï¼Œé¿å…é˜»å¡žä¸šåŠ¡è¯·æ±‚
	}
}

func (q *Queue) consumeChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-q.ch:
			if err := q.persist(ctx, ev); err == nil {
				q.markOutboxDone(ctx, ev.OutboxID)
			}
		}
	}
}

func (q *Queue) consumeStream(ctx context.Context) {
	q.ensureGroup(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		streams, err := q.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    q.cfg.ConsumerGroup,
			Consumer: q.consumerName,
			Streams:  []string{q.cfg.StreamKey, ">"},
			Count:    10,
			Block:    2 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil || ctx.Err() != nil {
				continue
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, st := range streams {
			for _, msg := range st.Messages {
				ev, ok := q.eventFromStream(msg.Values)
				if !ok {
					_ = q.rdb.XAck(ctx, q.cfg.StreamKey, q.cfg.ConsumerGroup, msg.ID)
					continue
				}
				if err := q.persist(ctx, ev); err == nil {
					q.markOutboxDone(ctx, ev.OutboxID)
					_ = q.rdb.XAck(ctx, q.cfg.StreamKey, q.cfg.ConsumerGroup, msg.ID)
				}
			}
		}
	}
}

func (q *Queue) ensureGroup(ctx context.Context) {
	err := q.rdb.XGroupCreateMkStream(ctx, q.cfg.StreamKey, q.cfg.ConsumerGroup, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		// ç»„å·²å­˜åœ¨æˆ–å…¶å®ƒçž¬æ—¶é”™è¯¯ï¼šå¯åŠ¨åŽä»å¯ XREADGROUP
		_ = err
	}
}

func (q *Queue) eventFromStream(values map[string]any) (Event, bool) {
	var ev Event
	raw, _ := values["data"].(string)
	if raw == "" {
		return ev, false
	}
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		return ev, false
	}
	if ev.OutboxID == "" {
		if oid, ok := values["outbox_id"].(string); ok {
			ev.OutboxID = oid
		}
	}
	return ev, true
}

func (q *Queue) markOutboxDone(ctx context.Context, id string) {
	if id == "" || q.db == nil {
		return
	}
	_ = q.db.WithContext(ctx).Model(&OutboxRow{}).
		Where("id = ?", id).
		Update("status", "DONE").Error
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
