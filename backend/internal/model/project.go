package model

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:30;default:Draft;not null" json:"status"`
	OwnerID     uint      `gorm:"not null;index" json:"owner_id"`
	Owner       *User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Version     int       `gorm:"default:1" json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index" json:"project_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      string    `gorm:"size:20;default:Developer" json:"role"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
}
