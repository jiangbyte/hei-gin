package model

import (
	"crypto/sha1"
	"fmt"
	"time"

	"hei-gin/sdk/enums"
)

// ─── New Unified Models ───────────────────────────────────────────────

// Message is the unified single-chat message table, replacing both SysMessage and ClientMessage.
type Message struct {
	ID             string     `gorm:"primaryKey;size:32" json:"id"`
	ConversationID string     `gorm:"size:32;not null;index;index:idx_msg_conv_created,priority:1" json:"conversation_id"`
	Content        string     `gorm:"type:text" json:"content"`
	Extra          string     `gorm:"type:text" json:"extra"`
	MsgType        string     `gorm:"size:20;default:TEXT" json:"msg_type"`
	SenderID       string     `gorm:"size:32;index;index:idx_msg_sender_type_created,priority:1" json:"sender_id"`
	SenderType     string     `gorm:"size:20;index:idx_msg_sender_type_created,priority:2" json:"sender_type"`
	ReceiverID     string     `gorm:"size:32;index;index:idx_msg_receiver_type_status,priority:1" json:"receiver_id"`
	ReceiverType   string     `gorm:"size:20;index:idx_msg_receiver_type_status,priority:2" json:"receiver_type"`
	Status         string     `gorm:"size:10;not null;default:unread;index:idx_msg_receiver_type_status,priority:3" json:"status"`
	DeletedBy      string     `gorm:"size:32" json:"deleted_by,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      *time.Time `gorm:"index:idx_msg_sender_type_created,priority:3;index:idx_msg_conv_created,priority:2" json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

func (Message) TableName() string { return "im_message" }

// Conversation tracks conversation metadata (cache table).
type Conversation struct {
	ID        string     `gorm:"primaryKey;size:32" json:"id"`
	FromID    string     `gorm:"column:from_id;size:32;not null" json:"from_id"`
	FromType  string     `gorm:"column:from_type;size:20;not null" json:"from_type"`
	ToID      string     `gorm:"column:to_id;size:32;not null" json:"to_id"`
	ToType    string     `gorm:"column:to_type;size:20;not null" json:"to_type"`
	LastMsg   string     `gorm:"type:text" json:"last_msg"`
	LastTime  *time.Time `json:"last_time"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (Conversation) TableName() string { return "im_conversation" }

// ConversationUnread tracks per-user unread count per conversation.
type ConversationUnread struct {
	ID             string `gorm:"primaryKey;size:32" json:"id"`
	ConversationID string `gorm:"size:32;uniqueIndex:idx_conv_unread" json:"conversation_id"`
	UserID         string `gorm:"size:32;uniqueIndex:idx_conv_unread" json:"user_id"`
	UserType       string `gorm:"size:20;uniqueIndex:idx_conv_unread" json:"user_type"`
	UnreadCount    int    `gorm:"default:0" json:"unread_count"`
}

func (ConversationUnread) TableName() string { return "im_conversation_unread" }

// ─── Legacy Models (keep for backward compat during migration) ──────

// Shared message type constants
const (
	MsgTypeText   = "TEXT"
	MsgTypeImage  = "IMAGE"
	MsgTypeFile   = "FILE"
	MsgTypeSystem = "SYSTEM"
)

// MsgExtraImage holds extra metadata for IMAGE messages.
type MsgExtraImage struct {
	Width     int    `json:"w,omitempty"`
	Height    int    `json:"h,omitempty"`
	Format    string `json:"format,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// MsgExtraFile holds extra metadata for FILE messages.
type MsgExtraFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	MIME string `json:"mime"`
}

// MsgExtraSystem holds extra metadata for SYSTEM messages.
type MsgExtraSystem struct {
	Action     string `json:"action"`
	OperatorID string `json:"operator_id,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	UserType   string `json:"user_type,omitempty"`
}

// GenerateConversationID generates a deterministic conversation ID from two user identifiers.
func GenerateConversationID(u1ID string, u1Type enums.LoginTypeEnum, u2ID string, u2Type enums.LoginTypeEnum) string {
	key1 := string(u1Type) + ":" + u1ID
	key2 := string(u2Type) + ":" + u2ID
	if key1 > key2 {
		key1, key2 = key2, key1
	}
	h := sha1.Sum([]byte(key1 + "|" + key2))
	return fmt.Sprintf("%x", h[:8])
}
