// internal/modules/iam/result/loader.go 账号视图加载。
//
// Author: Charlie

package result

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	adminuser "hei-gin/internal/modules/user/admin"
	portaluser "hei-gin/internal/modules/user/portal"
)

// accountRow sys_account è¡ŒæŠ•å½±ï¼ˆä»…å›žæ˜¾æ‰€éœ€å­—æ®µï¼‰ã€‚
type accountRow struct {
	ID                 string     `gorm:"column:id"`
	AccountType        string     `gorm:"column:account_type"`
	AccountStatus      string     `gorm:"column:account_status"`
	CancelledAt        *time.Time `gorm:"column:cancelled_at"`
	CancelledBy        *string    `gorm:"column:cancelled_by"`
	CancelReason       *string    `gorm:"column:cancel_reason"`
	LastLoginIP        *string    `gorm:"column:last_login_ip"`
	LastLoginAddress   *string    `gorm:"column:last_login_address"`
	LastLoginTime      *time.Time `gorm:"column:last_login_time"`
	LastLoginDevice    *string    `gorm:"column:last_login_device"`
	LatestLoginIP      *string    `gorm:"column:latest_login_ip"`
	LatestLoginAddress *string    `gorm:"column:latest_login_address"`
	LatestLoginTime    *time.Time `gorm:"column:latest_login_time"`
	LatestLoginDevice  *string    `gorm:"column:latest_login_device"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	CreatedBy          *string    `gorm:"column:created_by"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
	UpdatedBy          *string    `gorm:"column:updated_by"`
}

// TableName è¿”å›žè¡¨åã€‚
func (accountRow) TableName() string { return "sys_account" }

// identityRow sys_account_identity è¡ŒæŠ•å½±ã€‚
type identityRow struct {
	AccountID  string `gorm:"column:account_id"`
	Identifier string `gorm:"column:identifier"`
}

// TableName è¿”å›žè¡¨åã€‚
func (identityRow) TableName() string { return "sys_account_identity" }

// LoadAccountViews æŒ‰ ID åˆ—è¡¨åŠ è½½è´¦å·ç»“æžœè¡Œï¼ˆä¿æŒå…¥å‚é¡ºåºï¼Œç¼ºçœè·³è¿‡ï¼‰ã€‚
func LoadAccountViews(ctx context.Context, db *gorm.DB, ids []string) ([]AccountView, error) {
	out := make([]AccountView, 0, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var rows []accountRow
	if err := db.WithContext(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	byID := make(map[string]accountRow, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
	}
	// ä¸»ç™»å½•æ ‡è¯†
	var idents []identityRow
	if err := db.WithContext(ctx).Where("account_id IN ? AND identity_type = ?", ids, "ACCOUNT").Find(&idents).Error; err != nil {
		return nil, err
	}
	identByAccount := map[string]string{}
	for _, it := range idents {
		if _, ok := identByAccount[it.AccountID]; !ok {
			identByAccount[it.AccountID] = it.Identifier
		}
	}
	adminRepo := adminuser.NewRepo(db)
	portalRepo := portaluser.NewRepo(db)
	for _, id := range ids {
		row, ok := byID[id]
		if !ok {
			continue
		}
		v := AccountView{
			ID: row.ID, AccountType: row.AccountType, AccountStatus: row.AccountStatus,
			CancelledAt: row.CancelledAt, CancelledBy: row.CancelledBy, CancelReason: row.CancelReason,
			LastLoginIP: row.LastLoginIP, LastLoginAddress: row.LastLoginAddress, LastLoginTime: row.LastLoginTime,
			LastLoginDevice: row.LastLoginDevice, LatestLoginIP: row.LatestLoginIP,
			LatestLoginAddress: row.LatestLoginAddress, LatestLoginTime: row.LatestLoginTime,
			LatestLoginDevice: row.LatestLoginDevice,
			CreatedAt:         row.CreatedAt, CreatedBy: row.CreatedBy, UpdatedAt: row.UpdatedAt, UpdatedBy: row.UpdatedBy,
		}
		v.Account = identByAccount[id]
		if row.AccountType == string(security.AccountAdmin) {
			if p, err := adminRepo.GetProfile(ctx, id); err == nil {
				v.Name, v.Nickname, v.Avatar, v.Signature, v.Phone, v.Email, v.Remark =
					p.Name, p.Nickname, p.Avatar, p.Signature, p.Phone, p.Email, p.Remark
			}
		} else if p, err := portalRepo.GetProfile(ctx, id); err == nil {
			v.Name, v.Nickname, v.Avatar, v.Signature, v.Phone, v.Email =
				p.Name, p.Nickname, p.Avatar, p.Signature, p.Phone, p.Email
		}
		out = append(out, v)
	}
	return out, nil
}

// LoadAccountTypes æŒ‰ ID åˆ—è¡¨æŸ¥è´¦å·ç±»åž‹ï¼ˆä¾›æˆå‘˜æŽˆæƒæ—¶å›žå¡« account_typeï¼‰ã€‚
func LoadAccountTypes(ctx context.Context, db *gorm.DB, ids []string) (map[string]string, error) {
	out := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var rows []struct {
		ID          string `gorm:"column:id"`
		AccountType string `gorm:"column:account_type"`
	}
	if err := db.WithContext(ctx).Table("sys_account").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = r.AccountType
	}
	return out, nil
}
