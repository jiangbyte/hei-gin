package account

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/module"
)

// PurgeExpiredCancelled 物理清理超期已注销账号（对齐 boot accountPurgeCancelledJob）。
func (s *Service) PurgeExpiredCancelled(ctx context.Context, retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		retentionDays = 15
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	var ids []string
	if err := s.repo.DB().WithContext(ctx).Model(&Account{}).
		Where("account_status = ? AND cancelled_at IS NOT NULL AND cancelled_at < ?", security.AccountStatusCancelled, cutoff).
		Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}
	if err := s.Delete(ctx, ids); err != nil {
		return 0, err
	}
	return int64(len(ids)), nil
}

func (s *Service) accountPurgeCancelledJobHandler(ctx context.Context, param string) error {
	days := 15
	if p := strings.TrimSpace(param); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			days = n
		}
	}
	n, err := s.PurgeExpiredCancelled(ctx, days)
	if err != nil {
		return err
	}
	_ = fmt.Sprintf("purged=%d", n)
	return nil
}

// withJobs 附加 XXL-JOB handlers。
func (s *Service) withJobs(m module.Module) module.Module {
	m.Jobs = append(m.Jobs, module.Job{
		Name: "accountPurgeCancelledJob",
		Run:  s.accountPurgeCancelledJobHandler,
	})
	return m
}
