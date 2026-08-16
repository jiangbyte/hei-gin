// internal/modules/iam/account/job.go 定时任务。
//
// Author: Charlie

package account

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/module"
)

// PurgeExpiredCancelled 物理清理超期已注销账号（JobHandler: accountPurgeCancelledJob）。
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

func (s *Service) accountPurgeCancelledJobHandler(ctx context.Context, param string) (string, error) {
	days := 15
	if p := strings.TrimSpace(param); p != "" {
		// 支持纯数字或 {"retentionDays":15}
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			days = n
		} else if strings.Contains(p, "retentionDays") {
			p = strings.TrimSpace(p)
			for _, part := range strings.FieldsFunc(p, func(r rune) bool {
				return r == '{' || r == '}' || r == ',' || r == '"' || r == ':'
			}) {
				if n, err := strconv.Atoi(strings.TrimSpace(part)); err == nil && n > 0 {
					days = n
					break
				}
			}
		}
	}
	n, err := s.PurgeExpiredCancelled(ctx, days)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("purged=%d", n), nil
}

// withJobs 附加任务处理器（gojob 调度器收集）。
func (s *Service) withJobs(m module.Module) module.Module {
	m.Jobs = append(m.Jobs, module.Job{
		Name: "iam_account_purge_cancelled",
		Run:  s.accountPurgeCancelledJobHandler,
	})
	return m
}
