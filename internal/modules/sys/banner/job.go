// internal/modules/sys/banner/job.go 定时任务。
//
// Author: Charlie

package banner

import (
	"context"
	"fmt"
	"time"
)

// SyncStatusBySchedule 按 start_at/end_at 启用或停用 Banner（SnailJob: bannerStatusJob）。
func (s *Service) SyncStatusBySchedule(ctx context.Context) (expired, activated int64, err error) {
	now := time.Now()
	res := s.repo.DB().WithContext(ctx).Model(&Banner{}).
		Where("status = ? AND end_at IS NOT NULL AND end_at < ?", "ENABLED", now).
		Updates(map[string]any{"status": "DISABLED", "updated_at": now})
	if res.Error != nil {
		return 0, 0, res.Error
	}
	expired = res.RowsAffected

	res = s.repo.DB().WithContext(ctx).Model(&Banner{}).
		Where("status = ? AND start_at IS NOT NULL AND start_at <= ?", "DISABLED", now).
		Where("(end_at IS NULL OR end_at >= ?)", now).
		Updates(map[string]any{"status": "ENABLED", "updated_at": now})
	if res.Error != nil {
		return expired, 0, res.Error
	}
	activated = res.RowsAffected
	return expired, activated, nil
}

// bannerStatusJobHandler SnailJob Handler。
func (s *Service) bannerStatusJobHandler(ctx context.Context, _ string) error {
	expired, activated, err := s.SyncStatusBySchedule(ctx)
	if err != nil {
		return err
	}
	_ = fmt.Sprintf("expired=%d,activated=%d", expired, activated)
	return nil
}
