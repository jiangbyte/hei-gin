package friend

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/utils"
)

// ImFriendRequestToFriendRequestVO 将 imModel.FriendRequest 映射到 FriendRequestVO
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

// ImFriendshipToFriendVO 将 imModel.Friendship 映射到 FriendVO
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

// ImFriendBlockToBlockVO 将 imModel.FriendBlock 映射到 BlockVO
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
