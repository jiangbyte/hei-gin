// internal/modules/sys/job/service.go 任务管理业务服务（对齐 hei-boot JobService）。
//
// Author: Charlie

package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/db/dialect"
	"hei-gin/internal/framework/platform/gojob"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// Service 任务管理业务服务。
//
// Author: Charlie
type Service struct {
	db   *gorm.DB
	jobs *gojob.Manager
}

// NewService 构造服务。
func NewService(db *gorm.DB, jobs *gojob.Manager) *Service {
	return &Service{db: db, jobs: jobs}
}

// New 构建 sys.job 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB, d.Jobs)
	return module.Module{
		Name:   "sys.job",
		Models: []any{&gojob.SysJob{}, &gojob.SysJobLog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
		Jobs: []module.Job{
			{Name: "sys_job_sample", Run: sampleHandler},
			{Name: "sys_job_log_cleanup", Run: s.logCleanupHandler},
		},
	}
}

// Create 创建任务。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	if err := gojob.ValidateTrigger(req.ExecuteType, req.TriggerConfig); err != nil {
		return err
	}
	if s.jobs != nil && !s.jobs.HasHandler(strings.TrimSpace(req.ExecuteClass)) {
		return fmt.Errorf("未找到任务处理器: %s", req.ExecuteClass)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	now := time.Now().UTC()
	next, err := gojob.ComputeNextRunTime(req.ExecuteType, req.TriggerConfig, now)
	if err != nil {
		return err
	}
	row := gojob.SysJob{
		ID:            idgen.Next(),
		JobName:       strings.TrimSpace(req.JobName),
		ExecuteClass:  strings.TrimSpace(req.ExecuteClass),
		ExecuteType:   strings.ToUpper(strings.TrimSpace(req.ExecuteType)),
		TriggerConfig: strings.TrimSpace(req.TriggerConfig),
		ExecuteParam:  gojob.ParamJSON(req.ExecuteParam),
		NextRunTime:   next,
		Enabled:       enabled,
		Description:   req.Description,
		Sort:          req.Sort,
	}
	return s.db.WithContext(ctx).Create(&row).Error
}

// Update 更新任务。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", req.ID).Error; err != nil {
		return fmt.Errorf("job not found")
	}
	if err := gojob.ValidateTrigger(req.ExecuteType, req.TriggerConfig); err != nil {
		return err
	}
	if s.jobs != nil && !s.jobs.HasHandler(strings.TrimSpace(req.ExecuteClass)) {
		return fmt.Errorf("未找到任务处理器: %s", req.ExecuteClass)
	}
	enabled := row.Enabled
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	typeChanged := !strings.EqualFold(row.ExecuteType, req.ExecuteType) ||
		row.TriggerConfig != strings.TrimSpace(req.TriggerConfig)
	updates := map[string]any{
		"job_name":       strings.TrimSpace(req.JobName),
		"execute_class":  strings.TrimSpace(req.ExecuteClass),
		"execute_type":   strings.ToUpper(strings.TrimSpace(req.ExecuteType)),
		"trigger_config": strings.TrimSpace(req.TriggerConfig),
		"execute_param":  gojob.ParamJSON(req.ExecuteParam),
		"enabled":        enabled,
		"description":    req.Description,
		"sort":           req.Sort,
	}
	if typeChanged {
		next, err := gojob.ComputeNextRunTime(req.ExecuteType, req.TriggerConfig, time.Now().UTC())
		if err != nil {
			return err
		}
		updates["next_run_time"] = next
	}
	return s.db.WithContext(ctx).Model(&gojob.SysJob{}).Where("id = ?", row.ID).Updates(updates).Error
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&gojob.SysJob{}).Error
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*gojob.SysJob, error) {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job not found")
	}
	return &row, nil
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) ([]gojob.SysJob, int64, int, int, error) {
	cur, size := q.Normalize()
	tx := s.db.WithContext(ctx).Model(&gojob.SysJob{})
	if n := strings.TrimSpace(q.JobName); n != "" {
		tx = tx.Where(dialect.ILike(tx, "job_name"), "%"+n+"%")
	}
	if t := strings.TrimSpace(q.ExecuteType); t != "" {
		tx = tx.Where("execute_type = ?", strings.ToUpper(t))
	}
	if q.Enabled != nil {
		tx = tx.Where("enabled = ?", *q.Enabled)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, cur, size, err
	}
	var rows []gojob.SysJob
	err := tx.Order("sort ASC, created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, cur, size, err
}

// SetEnabled 启停。
func (s *Service) SetEnabled(ctx context.Context, req EnabledParam) error {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", req.ID).Error; err != nil {
		return fmt.Errorf("job not found")
	}
	updates := map[string]any{"enabled": req.Enabled}
	if req.Enabled {
		next, err := gojob.ComputeNextRunTime(row.ExecuteType, row.TriggerConfig, time.Now().UTC())
		if err != nil {
			return err
		}
		updates["next_run_time"] = next
	}
	return s.db.WithContext(ctx).Model(&gojob.SysJob{}).Where("id = ?", row.ID).Updates(updates).Error
}

// RunNow 立即执行。
func (s *Service) RunNow(ctx context.Context, id, executor string) error {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return fmt.Errorf("job not found")
	}
	if !row.Enabled {
		return fmt.Errorf("任务未启用")
	}
	if s.jobs == nil {
		return fmt.Errorf("调度器未就绪")
	}
	if executor == "" {
		executor = gojob.ExecutorSystem
	}
	s.jobs.SubmitRun(ctx, id, true, executor)
	return nil
}

// Logs 执行日志分页。
func (s *Service) Logs(ctx context.Context, q LogParam) ([]gojob.SysJobLog, int64, int, int, error) {
	cur, size := q.Normalize()
	tx := s.db.WithContext(ctx).Model(&gojob.SysJobLog{})
	if id := strings.TrimSpace(q.JobID); id != "" {
		tx = tx.Where("job_id = ?", id)
	}
	if q.Success != nil {
		tx = tx.Where("success = ?", *q.Success)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, cur, size, err
	}
	var rows []gojob.SysJobLog
	err := tx.Order("execute_time DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, cur, size, err
}

func sampleHandler(_ context.Context, paramJSON string) (string, error) {
	if strings.TrimSpace(paramJSON) == "" || paramJSON == "{}" || paramJSON == "null" {
		return "echo: (无参数)", nil
	}
	return "echo: " + paramJSON, nil
}

func (s *Service) logCleanupHandler(ctx context.Context, paramJSON string) (string, error) {
	retention := 30
	batch := 1000
	if s.jobs != nil {
		cfg := s.jobs.ConfigValues()
		retention = cfg.LogRetentionDays
		batch = cfg.LogBatchSize
	}
	if strings.TrimSpace(paramJSON) != "" && paramJSON != "null" {
		var m map[string]any
		if err := json.Unmarshal([]byte(paramJSON), &m); err == nil {
			if v, ok := asInt(m["retentionDays"]); ok && v > 0 {
				retention = v
			}
			if v, ok := asInt(m["batchSize"]); ok && v > 0 {
				batch = v
			}
		}
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -retention)
	res := s.db.WithContext(ctx).
		Where("execute_time < ?", cutoff).
		Limit(batch).
		Delete(&gojob.SysJobLog{})
	if res.Error != nil {
		return "", res.Error
	}
	return fmt.Sprintf("deleted=%d,retentionDays=%d,batchSize=%d", res.RowsAffected, retention, batch), nil
}

func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	default:
		return 0, false
	}
}
