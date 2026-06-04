package model

import "time"

type Attachment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FileName    string    `gorm:"size:255;not null" json:"file_name"`
	FileURL     string    `gorm:"size:500;not null" json:"file_url"`
	FileSize    int64     `gorm:"not null" json:"file_size"`
	ContentType string    `gorm:"size:100" json:"content_type"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TaskID      *uint     `gorm:"index" json:"task_id,omitempty"`
	ProjectID   *uint     `gorm:"index" json:"project_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
