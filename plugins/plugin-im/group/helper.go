package group

import (
	"context"
	"encoding/json"
	"strings"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/infra/storage"
	"hei-gin/sdk/web/exception"

	imModel "hei-gin/plugins/plugin-im/model"

	"github.com/gin-gonic/gin"
)

// ── Internal helpers ─────────────────────────────────────────────

func getLoginID(c *gin.Context) string {
	if v, ok := c.Get("login_id"); ok {
		if uid, ok := v.(string); ok && uid != "" {
			return uid
		}
	}
	if getUserType(c) == string(enums.LoginTypeConsumer) {
		return auth.Consumer.GetLoginID(c)
	}
	return auth.Business.GetLoginID(c)
}

func getUserType(c *gin.Context) string {
	if v, ok := c.Get("login_type"); ok {
		if loginType, ok := v.(string); ok && loginType != "" {
			return loginType
		}
	}
	path := c.Request.URL.Path
	if isConsumerAPIPath(path) {
		return string(enums.LoginTypeConsumer)
	}
	return string(enums.LoginTypeBusiness)
}

func isConsumerAPIPath(path string) bool {
	if !strings.HasPrefix(path, "/api/v") {
		return false
	}
	afterVersionPrefix := path[len("/api/v"):]
	slash := strings.IndexByte(afterVersionPrefix, '/')
	if slash < 0 {
		return false
	}
	return strings.HasPrefix(afterVersionPrefix[slash+1:], "c/")
}

// ==================== GroupCreate ====================

func validateMemberType(groupType, userType string) error {
	if groupType == GroupTypeConsumerOnly && userType != string(enums.LoginTypeConsumer) {
		return exception.NewBusinessError("该群仅限C端用户", 403)
	}
	return nil
}
func (s *Service) checkOwnerOrAdmin(ctx context.Context, groupID, userID, userType string) (*imModel.Group, *imModel.GroupMember, error) {
	if groupID == "" || userID == "" {
		return nil, nil, exception.NewBusinessError("参数错误", 400)
	}

	group, err := s.repo.FindGroupByID(ctx, groupID)
	if err != nil {
		return nil, nil, exception.NewBusinessError("群不存在", 400)
	}
	member, err := s.repo.FindActiveMember(ctx, groupID, userID, userType)
	if err != nil {
		return nil, nil, exception.NewBusinessError("不在群中", 400)
	}
	if member.Role != RoleOwner && member.Role != RoleAdmin {
		return nil, nil, exception.NewBusinessError("无权限", 403)
	}
	return group, member, nil
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
