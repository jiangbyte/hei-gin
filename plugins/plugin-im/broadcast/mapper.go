package broadcast

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/utils"
)

// ImBroadcastToBroadcastVO 将 imModel.Broadcast 映射到 BroadcastVO
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
