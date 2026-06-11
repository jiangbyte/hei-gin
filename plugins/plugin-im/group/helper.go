package group

import (
	"context"
	"encoding/json"
	"strings"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/storage"

	imModel "hei-gin/plugins/plugin-im/model"

	"github.com/gin-gonic/gin"
)

// ── Internal helpers ─────────────────────────────────────────────

func getLoginID(c *gin.Context) string {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return auth.Consumer.GetLoginID(c)
	}
	return auth.GetLoginID(c)
}

func getUserType(c *gin.Context) string {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return string(enums.LoginTypeConsumer)
	}
	return string(enums.LoginTypeBusiness)
}

// ==================== GroupCreate ====================

func validateMemberType(groupType, userType string) error {
	if groupType == GroupTypeConsumerOnly && userType != string(enums.LoginTypeConsumer) {
		return exception.NewBusinessError("该群仅限C端用户", 403)
	}
	return nil
}
func checkOwnerOrAdmin(ctx context.Context, groupID, userID, userType string) (*imModel.Group, *imModel.GroupMember, error) {
	if groupID == "" || userID == "" {
		return nil, nil, exception.NewBusinessError("参数错误", 400)
	}

	var group imModel.Group
	if err := db.DB.WithContext(ctx).First(&group, "id = ?", groupID).Error; err != nil {
		return nil, nil, exception.NewBusinessError("群不存在", 400)
	}
	var member imModel.GroupMember
	if err := db.DB.WithContext(ctx).
		Where("group_id = ? AND user_id = ? AND user_type = ? AND status = ?",
			groupID, userID, userType, MemberActive).
		First(&member).Error; err != nil {
		return nil, nil, exception.NewBusinessError("不在群中", 400)
	}
	if member.Role != RoleOwner && member.Role != RoleAdmin {
		return nil, nil, exception.NewBusinessError("无权限", 403)
	}
	return &group, &member, nil
}
func resolveFileURL(content, extra string) string {
	if strings.HasPrefix(content, "http") {
		return content
	}
	if content == "" {
		return ""
	}
	engine := "LOCAL"
	bucket := "DEFAULT"
	if extra != "" {
		var meta struct {
			Engine string `json:"engine"`
			Bucket string `json:"bucket"`
		}
		if err := json.Unmarshal([]byte(extra), &meta); err == nil {
			if meta.Engine != "" {
				engine = meta.Engine
			}
			if meta.Bucket != "" {
				bucket = meta.Bucket
			}
		}
	}
	return storage.GetURL(engine, bucket, content)
}
