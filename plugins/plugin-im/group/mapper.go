package group

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/utils"
	"time"
)

// ImGroupToGroupVO 将 imModel.Group 映射到 GroupVO
func ImGroupToGroupVO(src *imModel.Group) *GroupVO {
	if src == nil {
		return nil
	}
	dst := &GroupVO{}
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Avatar = src.Avatar
	dst.OwnerID = src.OwnerID
	dst.OwnerType = src.OwnerType
	dst.GroupType = src.GroupType
	dst.Notice = src.Notice
	return dst
}

// ImGroupMemberToMemberVO 将 imModel.GroupMember 映射到 MemberVO
func ImGroupMemberToMemberVO(src *imModel.GroupMember) *MemberVO {
	if src == nil {
		return nil
	}
	dst := &MemberVO{}
	dst.UserID = src.UserID
	dst.UserType = src.UserType
	dst.Role = src.Role
	dst.Nickname = src.Nickname
	dst.IsMuted = src.MutedUntil != nil && src.MutedUntil.After(time.Now())
	dst.JoinedAt = utils.FormatDateTimePtr(src.JoinedAt)
	return dst
}

// ImGroupMessageToMessageVO 将 imModel.GroupMessage 映射到 MessageVO
func ImGroupMessageToMessageVO(src *imModel.GroupMessage, fileURL string) *MessageVO {
	if src == nil {
		return nil
	}
	dst := &MessageVO{}
	dst.ID = src.ID
	dst.SenderID = src.SenderID
	dst.SenderType = src.SenderType
	dst.Content = src.Content
	dst.Extra = src.Extra
	dst.MsgType = src.MsgType
	dst.ReplyTo = src.ReplyTo
	dst.FileURL = fileURL
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	return dst
}
