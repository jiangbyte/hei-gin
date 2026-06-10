package message

// MessageVO is the view object for a message.
type MessageVO struct {
	ConversationID string  `json:"conversation_id"`
	ID             string  `json:"id"`
	Content        string  `json:"content,omitempty"`
	MsgType        string  `json:"msg_type"`
	Extra          string  `json:"extra,omitempty"`
	SenderID       string  `json:"sender_id"`
	SenderType     string  `json:"sender_type"`
	ReceiverID     string  `json:"receiver_id"`
	ReceiverType   string  `json:"receiver_type"`
	Status         string  `json:"status"`
	ReadAt         *string `json:"read_at"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type MessagePageParam struct {
	Current int    `json:"current" form:"current"`
	Size    int    `json:"size" form:"size"`
	Status  string `json:"status" form:"status"`
}

type MessageSendParam struct {
	Content      string   `json:"content"`
	MsgType      string   `json:"msg_type"`
	Extra        string   `json:"extra"`
	ReceiverIDs  []string `json:"receiver_ids"`
	ReceiverType string   `json:"receiver_type"`
}

type RecallParam struct {
	MessageID string `json:"message_id"`
}

type ForwardParam struct {
	MessageID  string   `json:"message_id"`
	TargetIDs  []string `json:"target_ids"`
	TargetType string   `json:"target_type"`
}

type SearchParam struct {
	Keyword string `json:"keyword" form:"keyword"`
	Cursor  string `json:"cursor" form:"cursor"`
	Size    int    `json:"size" form:"size"`
}

type UnreadCountVO struct {
	Count int64 `json:"count"`
}

// Conversation types
const (
	ConvTypeSingle = "single"
	ConvTypeGroup  = "group"
)

type ConversationVO struct {
	ConversationID   string `json:"conversation_id"`
	ConversationType string `json:"conversation_type"`

	// Single-chat fields
	OtherUserID   string `json:"other_user_id,omitempty"`
	OtherUserType string `json:"other_user_type,omitempty"`
	OtherNickname string `json:"other_nickname,omitempty"`
	OtherAvatar   string `json:"other_avatar,omitempty"`

	// Group fields
	GroupID     string `json:"group_id,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	GroupAvatar string `json:"group_avatar,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`

	LastContent string `json:"last_content"`
	LastTime    string `json:"last_time"`
	UnreadCount int64  `json:"unread_count"`
}

type ConversationMessageVO struct {
	FileURL      string `json:"file_url"`
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	SenderType string `json:"sender_type"`
	Content    string `json:"content"`
	MsgType    string `json:"msg_type"`
	Extra      string `json:"extra,omitempty"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// DeleteParam 删除消息参数
type DeleteParam struct {
	IDs []string `json:"ids"`
}

// ConversationReadParam 会话已读参数
type ConversationReadParam struct {
	ConversationID string `json:"conversation_id"`
}

type GetOrCreateConversationParam struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
}
