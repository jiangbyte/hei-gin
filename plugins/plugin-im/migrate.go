package plugin_im

import (
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/infra/db"
)

func init() {
	db.RegisterModel(&imModel.Message{})
	db.RegisterModel(&imModel.Conversation{})
	db.RegisterModel(&imModel.ConversationUnread{})
	db.RegisterModel(&imModel.ImFile{})
	db.RegisterModel(&imModel.Broadcast{})
	db.RegisterModel(&imModel.BroadcastRead{})
	db.RegisterModel(&imModel.FriendRequest{})
	db.RegisterModel(&imModel.Friendship{})
	db.RegisterModel(&imModel.FriendBlock{})
	db.RegisterModel(&imModel.Group{})
	db.RegisterModel(&imModel.GroupMember{})
	db.RegisterModel(&imModel.GroupJoinRequest{})
	db.RegisterModel(&imModel.GroupMessage{})
	db.RegisterModel(&imModel.GroupMessageRead{})
}
