package message

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/utils"
)

// ImMessageToMessageVO 将 imModel.Message 映射到 MessageVO
func ImMessageToMessageVO(src *imModel.Message) *MessageVO {
	if src == nil {
		return nil
	}
	v := &MessageVO{
		ID: src.ID, ConversationID: src.ConversationID,
		Content: src.Content, MsgType: src.MsgType, Extra: src.Extra,
		SenderID: src.SenderID, SenderType: src.SenderType,
		ReceiverID: src.ReceiverID, ReceiverType: src.ReceiverType,
		Status: src.Status,
		CreatedAt: utils.FormatDateTimePtr(src.CreatedAt),
		UpdatedAt: utils.FormatDateTimePtr(src.UpdatedAt),
	}
	if src.ReadAt != nil {
		s := utils.FormatDateTime(*src.ReadAt)
		v.ReadAt = &s
	}
	return v
}

// ImMessageToConversationMessageVO 将 imModel.Message 映射到 ConversationMessageVO
func ImMessageToConversationMessageVO(src *imModel.Message, fileURL string) *ConversationMessageVO {
	if src == nil {
		return nil
	}
	return &ConversationMessageVO{
		ID: src.ID, SenderID: src.SenderID, SenderType: src.SenderType,
		Content: src.Content, MsgType: src.MsgType, Extra: src.Extra,
		Status: src.Status, FileURL: fileURL,
		CreatedAt: utils.FormatDateTimePtr(src.CreatedAt),
	}
}
