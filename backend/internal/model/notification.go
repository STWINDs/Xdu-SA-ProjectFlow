package model

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Type      string    `gorm:"size:30;not null" json:"type"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	RelatedID uint      `json:"related_id"`
	CreatedAt time.Time `json:"created_at"`
}
