package friend

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
