package friend

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/utils"
)

// SendRequestParam 发送好友请求
type SendRequestParam struct {
	ReceiverID   string `json:"receiver_id" binding:"required"`
	ReceiverType string `json:"receiver_type" binding:"required"` // BUSINESS | CONSUMER
	Remark       string `json:"remark"`
}

// HandleRequestParam 处理好友请求
type HandleRequestParam struct {
	RequestID string `json:"request_id" binding:"required"`
}

// RemoveFriendParam 删除好友
type RemoveFriendParam struct {
	FriendID   string `json:"friend_id" binding:"required"`
	FriendType string `json:"friend_type" binding:"required"`
}

// FriendVO 好友视图
type FriendVO struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Remark   string `json:"remark"`
	AddedAt  string `json:"added_at"`
}

// FriendRequestVO 好友请求视图
type FriendRequestVO struct {
	ID           string `json:"id"`
	SenderID     string `json:"sender_id"`
	SenderType   string `json:"sender_type"`
	ReceiverID   string `json:"receiver_id"`
	ReceiverType string `json:"receiver_type"`
	Remark       string `json:"remark"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// BlockVO 黑名单视图
type BlockVO struct {
	BlockedID   string `json:"blocked_id"`
	BlockedType string `json:"blocked_type"`
	CreatedAt   string `json:"created_at"`
}

// RemarkParam 修改备注
type RemarkParam struct {
	FriendID   string `json:"friend_id" binding:"required"`
	FriendType string `json:"friend_type" binding:"required"`
	Remark     string `json:"remark"`
}

// BlockParam 拉黑/取消拉黑
type BlockParam struct {
	BlockedID   string `json:"blocked_id" binding:"required"`
	BlockedType string `json:"blocked_type" binding:"required"`
}

// SearchResult 搜索用户结果
type SearchResult struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func ImFriendRequestToFriendRequestVO(src *imModel.FriendRequest) *FriendRequestVO {
	if src == nil {
		return nil
	}
	dst := &FriendRequestVO{}
	dst.ID = src.ID
	dst.SenderID = src.SenderID
	dst.SenderType = src.SenderType
	dst.ReceiverID = src.ReceiverID
	dst.ReceiverType = src.ReceiverType
	dst.Remark = src.Remark
	dst.Status = src.Status
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	return dst
}

func ImFriendshipToFriendVO(src *imModel.Friendship) *FriendVO {
	if src == nil {
		return nil
	}
	dst := &FriendVO{}
	dst.UserID = src.FriendID
	dst.UserType = src.FriendType
	dst.Remark = src.Remark
	dst.AddedAt = utils.FormatDateTimePtr(src.CreatedAt)
	return dst
}

func ImFriendBlockToBlockVO(src *imModel.FriendBlock) *BlockVO {
	if src == nil {
		return nil
	}
	dst := &BlockVO{}
	dst.BlockedID = src.BlockedID
	dst.BlockedType = src.BlockedType
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	return dst
}
