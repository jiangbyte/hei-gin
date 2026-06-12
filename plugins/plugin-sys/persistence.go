package plugin_sys

import (
	"time"

	syslog "hei-gin/plugins/plugin-sys/log"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/shared/contracts"
	"hei-gin/sdk/utils"
)

// logPersister implements contracts.LogPersistenceAPI by saving to the sys_log table.
type logPersister struct{}

func (p *logPersister) SaveLog(entry contracts.LogEntry) error {
	now := time.Now()
	sysEntry := syslog.SysLog{
		ID:         utils.GenerateID(),
		Name:       entry.Name,
		Category:   entry.Category,
		ExeStatus:  entry.ExeStatus,
		ExeMessage: entry.ExeMessage,
		OpIP:       entry.OpIP,
		OpAddress:  entry.OpAddress,
		OpBrowser:  entry.OpBrowser,
		OpOs:       entry.OpOS,
		OpUser:     entry.OpUser,
		TraceID:    entry.TraceID,
		SignData:   entry.SignData,
		ReqMethod:  entry.ReqMethod,
		ReqURL:     entry.ReqURL,
		ParamJSON:  entry.ParamJSON,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	if entry.OpTime != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", entry.OpTime); err == nil {
			sysEntry.OpTime = &parsed
		}
	}
	return db.DB.Create(&sysEntry).Error
}

// Ensure compile-time check.
var _ contracts.LogPersistenceAPI = (*logPersister)(nil)
