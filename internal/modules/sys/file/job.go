// internal/modules/sys/file/job.go 定时任务。
//
// Author: Charlie

package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hei-gin/internal/framework/platform/storage"
)

// CleanupLocalOrphans 清理本地存储中早于 minAge 且无 sys_file 元数据的孤立文件
// （对齐 hei-fastapi sysFileCleanupLocalOrphans）。返回 scanned/deleted/skipped。
func (s *Service) CleanupLocalOrphans(ctx context.Context, minAgeSeconds, limit int64) (scanned, deleted, skipped int64, err error) {
	prov := s.sto.ProviderByName(ctx, "local")
	local, ok := prov.(*storage.Local)
	if !ok {
		return 0, 0, 0, nil // 非本地存储引擎时不处理
	}
	root := local.Root()
	if root == "" {
		return 0, 0, 0, nil
	}
	if _, statErr := os.Stat(root); statErr != nil {
		return 0, 0, 0, nil
	}
	// 最小保护窗：避免误删刚上传、元数据尚未落库的文件。
	if minAgeSeconds < 300 {
		minAgeSeconds = 300
	}
	if limit <= 0 {
		limit = 200
	}
	cutoff := time.Now().Add(-time.Duration(minAgeSeconds) * time.Second)

	// 收集早于 cutoff 的候选对象（相对 root 的斜杠路径），上限 limit*5。
	candidates := make([]string, 0, limit*5)
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return nil
		}
		if info.ModTime().After(cutoff) {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		candidates = append(candidates, filepath.ToSlash(rel))
		if int64(len(candidates)) >= limit*5 {
			return filepath.SkipAll
		}
		return nil
	})
	if len(candidates) == 0 {
		return 0, 0, 0, nil
	}

	scanned = int64(len(candidates))
	existing, listErr := s.repo.ListByObjectNames(ctx, candidates)
	if listErr != nil {
		return scanned, 0, 0, listErr
	}
	existSet := make(map[string]bool, len(existing))
	for i := range existing {
		existSet[existing[i].ObjectName] = true
	}

	for _, name := range candidates {
		if existSet[name] {
			skipped++
			continue
		}
		if delErr := prov.Delete(ctx, name); delErr != nil {
			continue // 删除失败不阻断
		}
		deleted++
	}
	return scanned, deleted, skipped, nil
}

// sysFileCleanupLocalOrphansHandler 任务 Handler（param 为保留分钟数，缺省 60 分钟）。
func (s *Service) sysFileCleanupLocalOrphansHandler(ctx context.Context, param string) error {
	minAgeSeconds := int64(3600)
	if p := strings.TrimSpace(param); p != "" {
		if n, convErr := strconv.ParseInt(p, 10, 64); convErr == nil && n > 0 {
			minAgeSeconds = n * 60
		}
	}
	scanned, deleted, skipped, err := s.CleanupLocalOrphans(ctx, minAgeSeconds, 200)
	if err != nil {
		return err
	}
	_ = fmt.Sprintf("scanned=%d,deleted=%d,skipped=%d", scanned, deleted, skipped)
	return nil
}
