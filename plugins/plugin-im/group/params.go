package group

type CreateParam struct {
	Name       string   `json:"name" binding:"required"`
	Avatar     string   `json:"avatar"`
	MemberIDs  []string `json:"member_ids"`
	MemberType string   `json:"member_type"` // BUSINESS | CONSUMER
}

type UpdateParam struct {
	GroupID string  `json:"group_id" binding:"required"`
	Name    *string `json:"name"`
	Avatar  *string `json:"avatar"`
	Notice  *string `json:"notice"`
}

type InviteParam struct {
	GroupID  string   `json:"group_id" binding:"required"`
	UserIDs  []string `json:"user_ids" binding:"required"`
	UserType string   `json:"user_type" binding:"required"`
}

type KickParam struct {
	GroupID  string `json:"group_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
}

type SetRoleParam struct {
	GroupID  string `json:"group_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Role     string `json:"role" binding:"required"` // admin | owner
}

type SendMessageParam struct {
	GroupID string `json:"group_id" binding:"required"`
	Content string `json:"content"`
	MsgType string `json:"msg_type"`
	Extra   string `json:"extra"` // JSON TEXT
	ReplyTo string `json:"reply_to"`
}

type GroupVO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	OwnerID     string `json:"owner_id"`
	OwnerType   string `json:"owner_type"`
	GroupType   string `json:"group_type"`
	Notice      string `json:"notice"`
	MemberCount int    `json:"member_count"`
	LastContent string `json:"last_content"`
	LastTime    string `json:"last_time"`
	UnreadCount int64  `json:"unread_count"`
}

type MemberVO struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Role     string `json:"role"`
	Nickname string `json:"nickname"`
	JoinedAt string `json:"joined_at"`
	IsMuted  bool   `json:"is_muted"`
}

type MessageVO struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	SenderType string `json:"sender_type"`
	Content    string `json:"content"`
	Extra      string `json:"extra"`
	MsgType    string `json:"msg_type"`
	ReplyTo    string `json:"reply_to"`
	FileURL    string `json:"file_url"`
	CreatedAt  string `json:"created_at"`
}

// HandleJoinRequestParam 处理入群申请
type HandleJoinRequestParam struct {
	RequestID string `json:"request_id" binding:"required"`
	Action    string `json:"action" binding:"required"` // approved | rejected
}

// TransferOwnerParam 转让群
type TransferOwnerParam struct {
	GroupID      string `json:"group_id" binding:"required"`
	NewOwnerID   string `json:"new_owner_id" binding:"required"`
	NewOwnerType string `json:"new_owner_type" binding:"required"`
}

// SetNicknameParam 设置群昵称
type SetNicknameParam struct {
	GroupID  string `json:"group_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Nickname string `json:"nickname"`
}

// DissolveParam 解散群
type DissolveParam struct {
	GroupID string `json:"group_id" binding:"required"`
}

// JoinOrLeaveParam 加入/退出群
type JoinOrLeaveParam struct {
	GroupID string `json:"group_id" binding:"required"`
}

// RecallMessageParam 撤回消息
type RecallMessageParam struct {
	GroupID   string `json:"group_id" binding:"required"`
	MessageID string `json:"message_id" binding:"required"`
}

// MarkReadParam 标记已读
type MarkReadParam struct {
	GroupID   string `json:"group_id" binding:"required"`
	MessageID string `json:"message_id" binding:"required"`
}

// MuteParam 禁言
type MuteParam struct {
	GroupID  string `json:"group_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Duration int    `json:"duration"`
}

// UnmuteParam 解禁
type UnmuteParam struct {
	GroupID  string `json:"group_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
}

// ConversationVO is used by MyGroupConversations for the unified conversation list.
type ConversationVO struct {
	ConversationID   string `json:"conversation_id,omitempty"`
	ConversationType string `json:"conversation_type,omitempty"`
	OtherUserID      string `json:"other_user_id,omitempty"`
	OtherUserType    string `json:"other_user_type,omitempty"`
	OtherNickname    string `json:"other_nickname,omitempty"`
	OtherAvatar      string `json:"other_avatar,omitempty"`
	GroupID          string `json:"group_id,omitempty"`
	GroupName        string `json:"group_name,omitempty"`
	GroupAvatar      string `json:"group_avatar,omitempty"`
	MemberCount      int    `json:"member_count,omitempty"`
	LastContent      string `json:"last_content"`
	LastTime         string `json:"last_time"`
	UnreadCount      int64  `json:"unread_count"`
}
