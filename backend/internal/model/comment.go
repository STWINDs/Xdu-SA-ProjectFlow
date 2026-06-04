package model

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectID *uint     `gorm:"index" json:"project_id,omitempty"`
	TaskID    *uint     `gorm:"index" json:"task_id,omitempty"`
	ParentID  *uint     `gorm:"index" json:"parent_id,omitempty"`
	Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
