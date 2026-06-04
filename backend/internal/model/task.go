package model

import "time"

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Priority    string     `gorm:"size:20;default:Medium" json:"priority"`
	Status      string     `gorm:"size:20;default:Todo" json:"status"`
	ProjectID   uint       `gorm:"not null;index" json:"project_id"`
	Project     *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	AssigneeID  *uint      `gorm:"index" json:"assignee_id"`
	Assignee    *User      `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	CreatorID   uint       `gorm:"not null;index" json:"creator_id"`
	Creator     *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Deadline    *time.Time `json:"deadline"`
	Version     int        `gorm:"default:1" json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
