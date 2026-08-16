// internal/modules/sys/job/service.go 任务管理业务服务。
//
// Author: Charlie

package job

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/gojob"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
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
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB, d.Jobs)
	return module.Module{
		Name:   "sys.job",
		Models: []any{&gojob.SysJob{}, &gojob.SysJobLog{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Handlers 可用处理器列表（轻量 DTO，供控制台下拉选择；不含 func，可 JSON 序列化）。
func (s *Service) Handlers() []gojob.HandlerInfo {
	if s.jobs == nil {
		return nil
	}
	return s.jobs.HandlerInfos()
}

// Create 创建任务并同步调度。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = gojob.StatusEnabled
	}
	row := gojob.SysJob{
		ID:         idgen.Next(),
		HandlerKey: strings.TrimSpace(req.HandlerKey),
		Name:       strings.TrimSpace(req.Name),
		CronExpr:   strings.TrimSpace(req.CronExpr),
		Params:     req.Params,
		Status:     status,
		Description: req.Description,
	}
	if err := s.validateJob(&row); err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	if s.jobs != nil {
		_ = s.jobs.SyncJob(ctx, &row)
	}
	return nil
}

// Update 更新任务并同步调度。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", req.ID).Error; err != nil {
		return fmt.Errorf("job not found")
	}
	row.HandlerKey = strings.TrimSpace(req.HandlerKey)
	row.Name = strings.TrimSpace(req.Name)
	row.CronExpr = strings.TrimSpace(req.CronExpr)
	row.Params = req.Params
	row.Status = strings.TrimSpace(req.Status)
	row.Description = req.Description
	if err := s.validateJob(&row); err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Model(&gojob.SysJob{}).Where("id = ?", row.ID).
		Updates(map[string]any{
			"handler_key": row.HandlerKey, "name": row.Name, "cron_expr": row.CronExpr,
			"params": row.Params, "status": row.Status, "description": row.Description,
		}).Error; err != nil {
		return err
	}
	if s.jobs != nil {
		_ = s.jobs.SyncJob(ctx, &row)
	}
	return nil
}

// Delete 批量删除任务并移除调度（保留执行日志）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var rows []gojob.SysJob
	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return err
	}
	if s.jobs != nil {
		for _, r := range rows {
			s.jobs.RemoveJob(r.HandlerKey)
		}
	}
	return s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&gojob.SysJob{}).Error
}

// SetStatus 启停任务。
func (s *Service) SetStatus(ctx context.Context, req StatusParam) error {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", req.ID).Error; err != nil {
		return fmt.Errorf("job not found")
	}
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status != gojob.StatusEnabled && status != gojob.StatusDisabled {
		return fmt.Errorf("invalid status: %s", req.Status)
	}
	row.Status = status
	if err := s.db.WithContext(ctx).Model(&gojob.SysJob{}).Where("id = ?", row.ID).
		Update("status", status).Error; err != nil {
		return err
	}
	if s.jobs != nil {
		_ = s.jobs.SyncJob(ctx, &row)
	}
	return nil
}

// Trigger 立即触发一次执行。
func (s *Service) Trigger(ctx context.Context, req TriggerParam) error {
	if s.jobs == nil {
		return errors.New("job scheduler not started")
	}
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", req.ID).Error; err != nil {
		return fmt.Errorf("job not found")
	}
	return s.jobs.Trigger(ctx, row.HandlerKey, req.Params)
}

// Detail 任务详情。
func (s *Service) Detail(ctx context.Context, id string) (*gojob.SysJob, error) {
	var row gojob.SysJob
	if err := s.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job not found")
	}
	return &row, nil
}

// Page 任务分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []gojob.SysJob, total int64, current, size int, err error) {
	current, size = q.Normalize()
	db := s.db.WithContext(ctx).Model(&gojob.SysJob{})
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if k := strings.TrimSpace(q.Keyword); k != "" {
		db = db.Where("(name ILIKE ? OR handler_key ILIKE ?)", "%"+k+"%", "%"+k+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	err = db.Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, current, size, err
}

// Logs 执行日志分页。
func (s *Service) Logs(ctx context.Context, q LogParam) (rows []gojob.SysJobLog, total int64, current, size int, err error) {
	current, size = q.Normalize()
	db := s.db.WithContext(ctx).Model(&gojob.SysJobLog{})
	if q.JobID != "" {
		db = db.Where("job_id = ?", q.JobID)
	}
	if q.HandlerKey != "" {
		db = db.Where("handler_key = ?", q.HandlerKey)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	err = db.Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, current, size, err
}

// validateJob 校验处理器存在且 cron 表达式合法。
func (s *Service) validateJob(row *gojob.SysJob) error {
	if s.jobs == nil {
		return errors.New("job scheduler not started")
	}
	found := false
	for _, h := range s.jobs.Handlers() {
		if h.Key == row.HandlerKey {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("handler not registered: %s", row.HandlerKey)
	}
	if row.Status == gojob.StatusEnabled && row.CronExpr != "" {
		if _, err := cronParse(row.CronExpr); err != nil {
			return fmt.Errorf("invalid cron expr: %w", err)
		}
	}
	return nil
}
