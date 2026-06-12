package message

import (
	"context"
	"sort"
	"strings"

	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"hei-gin/plugins/plugin-im/group"

	"github.com/gin-gonic/gin"
)

// ==================== MessageConversations ====================

func MessageConversations(c *gin.Context, cursor string, size int) ([]ConversationVO, bool) {
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	currentUserID := getLoginID(c)
	userType := getUserType(c)

	ctx := c.Request.Context()
	singleResult := buildFromMessages(ctx, currentUserID, userType)

	groupConvs := group.MyGroupConversationsWithContext(ctx, currentUserID, userType)
	for _, gv := range groupConvs {
		singleResult["group:"+gv.GroupID] = &ConversationVO{
			ConversationID:   "group:" + gv.GroupID,
			ConversationType: ConvTypeGroup,
			GroupID:          gv.GroupID,
			GroupName:        gv.GroupName,
			GroupAvatar:      gv.GroupAvatar,
			MemberCount:      gv.MemberCount,
			LastContent:      gv.LastContent,
			LastTime:         gv.LastTime,
			UnreadCount:      gv.UnreadCount,
		}
	}

	result := make([]ConversationVO, 0, len(singleResult))
	for _, vo := range singleResult {
		result = append(result, *vo)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastTime > result[j].LastTime
	})

	hasMore := false
	if cursor != "" {
		cut := -1
		for i, vo := range result {
			if vo.LastTime <= cursor {
				cut = i
				break
			}
		}
		if cut >= 0 {
			result = result[:cut]
		}
	}
	if len(result) > size {
		result = result[:size]
		hasMore = true
	}

	return result, hasMore
}

func buildFromMessages(ctx context.Context, currentUserID, userType string) map[string]*ConversationVO {
	rows := ListConversationLatest(ctx, currentUserID, userType)

	resultMap := make(map[string]*ConversationVO, len(rows))
	userKeys := make(map[string]bool)

	convIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		convIDs = append(convIDs, r.ConversationID)
	}

	unreads := CountConversationUnread(ctx, convIDs, currentUserID, userType)
	unreadMap := make(map[string]int64, len(unreads))
	for _, u := range unreads {
		unreadMap[u.ConversationID] = u.Count
	}

	for _, r := range rows {
		var otherID, otherType string
		if r.SenderID == currentUserID && r.SenderType == userType {
			otherID, otherType = r.ReceiverID, r.ReceiverType
		} else {
			otherID, otherType = r.SenderID, r.SenderType
		}

		userKeys[otherType+":"+otherID] = true

		resultMap[r.ConversationID] = &ConversationVO{
			ConversationID:   r.ConversationID,
			ConversationType: ConvTypeSingle,
			OtherUserID:      otherID,
			OtherUserType:    otherType,
			LastContent:      r.Content,
			LastTime:         utils.FormatDateTime(r.CreatedAt),
			UnreadCount:      unreadMap[r.ConversationID],
		}
	}

	if len(userKeys) > 0 {
		businessIDs, consumerIDs := []string{}, []string{}
		for key := range userKeys {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) == 2 {
				switch parts[0] {
				case string(enums.LoginTypeBusiness):
					businessIDs = append(businessIDs, parts[1])
				case string(enums.LoginTypeConsumer):
					consumerIDs = append(consumerIDs, parts[1])
				}
			}
		}

		nicknameMap := make(map[string]string)
		avatarMap := make(map[string]string)

		if len(businessIDs) > 0 {
			busUsers := FindBusinessUsers(ctx, businessIDs)
			for _, u := range busUsers {
				if u.Nickname != nil {
					nicknameMap[string(enums.LoginTypeBusiness)+":"+u.ID] = *u.Nickname
				}
				if u.Avatar != nil {
					avatarMap[string(enums.LoginTypeBusiness)+":"+u.ID] = *u.Avatar
				}
			}
		}
		if len(consumerIDs) > 0 {
			conUsers := FindConsumerUsers(ctx, consumerIDs)
			for _, u := range conUsers {
				if u.Nickname != nil {
					nicknameMap[string(enums.LoginTypeConsumer)+":"+u.ID] = *u.Nickname
				}
				if u.Avatar != nil {
					avatarMap[string(enums.LoginTypeConsumer)+":"+u.ID] = *u.Avatar
				}
			}
		}

		for _, vo := range resultMap {
			key := vo.OtherUserType + ":" + vo.OtherUserID
			vo.OtherNickname = nicknameMap[key]
			vo.OtherAvatar = avatarMap[key]
		}
	}

	return resultMap
}

// ==================== MessageConversationMessages ====================

func MessageConversationMessages(c *gin.Context, conversationID, cursor string, size int) ([]ConversationMessageVO, bool) {
	currentUserID := getLoginID(c)
	userType := getUserType(c)

	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	records := ListConversationMessages(c.Request.Context(), conversationID, currentUserID, userType, cursor, size)

	hasMore := len(records) > size
	if hasMore {
		records = records[:size]
	}

	result := make([]ConversationMessageVO, len(records))
	for i, m := range records {
		fileURL := ""
		if m.MsgType == "IMAGE" || m.MsgType == "FILE" {
			fileURL = ResolveFileURL(m.Content, m.Extra)
		}
		result[i] = *ImMessageToConversationMessageVO(&m, fileURL)
	}
	return result, hasMore
}

// ==================== MessageGetOrCreateConversation ====================

func MessageGetOrCreateConversation(c *gin.Context, param *GetOrCreateConversationParam) {
	currentUserID := getLoginID(c)
	userType := getUserType(c)

	if param.UserID == "" || param.UserType == "" {
		result.WriteError(c, exception.NewBusinessError("参数错误", 400))
		return
	}
	cid := imModel.GenerateConversationID(currentUserID, enums.LoginTypeEnum(userType), param.UserID, enums.LoginTypeEnum(param.UserType))

	displayName := param.UserID
	if param.UserType == string(enums.LoginTypeBusiness) {
		if u, err := FindBusinessUser(c.Request.Context(), param.UserID); err == nil && u.Nickname != nil {
			displayName = *u.Nickname
		}
	} else {
		if u, err := FindConsumerUser(c.Request.Context(), param.UserID); err == nil && u.Nickname != nil {
			displayName = *u.Nickname
		}
	}
	result.Success(c, gin.H{"conversation_id": cid, "display_name": displayName})
}
