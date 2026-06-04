package request

// CreateProjectReq is the request body for creating a new project.
type CreateProjectReq struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
}

// UpdateProjectReq is the request body for updating an existing project.
type UpdateProjectReq struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	Version     int    `json:"version" binding:"required"`
}

// AddMemberReq is the request body for adding a member to a project.
type AddMemberReq struct {
	UserID uint `json:"user_id" binding:"required"`
}
