package model

import "time"



// Broadcast represents a site-wide announcement.
type Broadcast struct {
	ID        string     `gorm:"primaryKey;size:32" json:"id"`
	Title     string     `gorm:"size:255;not null" json:"title"`
	Content   string     `gorm:"type:text" json:"content"`
	Scope     string     `gorm:"size:20;not null;default:ALL" json:"scope"` // ALL | BUSINESS | CONSUMER
	SenderID  string     `gorm:"size:32;not null" json:"sender_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (Broadcast) TableName() string { return "im_broadcast" }

// BroadcastRead tracks which users have read a broadcast.
type BroadcastRead struct {
	ID        string     `gorm:"primaryKey;size:32" json:"id"`
	BroadcastID string     `gorm:"size:32;not null;uniqueIndex:idx_br_broadcast_user" json:"broadcast_id"`
	UserID      string     `gorm:"size:32;not null;uniqueIndex:idx_br_broadcast_user" json:"user_id"`
	UserType    string     `gorm:"size:20;not null;uniqueIndex:idx_br_broadcast_user" json:"user_type"`
	ReadAt      *time.Time `json:"read_at"`
}

func (BroadcastRead) TableName() string { return "im_broadcast_read" }
