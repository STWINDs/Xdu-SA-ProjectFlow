package request

import "time"

type CreateTaskReq struct {
	Title       string     `json:"title" binding:"required,max=200"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" binding:"required,oneof=Low Medium High Critical"`
	ProjectID   uint       `json:"project_id" binding:"required"`
	AssigneeID  *uint      `json:"assignee_id"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTaskReq struct {
	Title       string     `json:"title" binding:"required,max=200"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" binding:"required,oneof=Low Medium High Critical"`
	AssigneeID  *uint      `json:"assignee_id"`
	Deadline    *time.Time `json:"deadline"`
	Version     int        `json:"version" binding:"required"`
}

type AssignTaskReq struct {
	AssigneeID uint `json:"assignee_id" binding:"required"`
}

type TransferStatusReq struct {
	Status string `json:"status" binding:"required,oneof=Todo InProgress Review Testing Done"`
}
