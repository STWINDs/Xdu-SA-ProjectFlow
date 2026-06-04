package request

type CreateCommentReq struct {
	Content   string `json:"content" binding:"required"`
	ProjectID *uint  `json:"project_id"`
	TaskID    *uint  `json:"task_id"`
}

type CreateReplyReq struct {
	Content string `json:"content" binding:"required"`
}
