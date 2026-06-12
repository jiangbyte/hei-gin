package broadcast

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/utils"
)

// SendBroadcastParam 发送全站通知
type SendBroadcastParam struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Scope   string `json:"scope"` // ALL | BUSINESS | CONSUMER
}

// BroadcastVO 通知视图
type BroadcastVO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content,omitempty"`
	Scope     string `json:"scope"`
	SenderID  string `json:"sender_id"`
	Read      bool   `json:"read"`
	ReadAt    string `json:"read_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

// ReadParam 标记已读参数
type ReadParam struct {
	BroadcastID string `json:"broadcast_id" binding:"required"`
}

func ImBroadcastToBroadcastVO(src *imModel.Broadcast) *BroadcastVO {
	if src == nil {
		return nil
	}
	dst := &BroadcastVO{}
	dst.ID = src.ID
	dst.Title = src.Title
	dst.Content = src.Content
	dst.Scope = src.Scope
	dst.SenderID = src.SenderID
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	return dst
}
