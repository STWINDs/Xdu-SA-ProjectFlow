package handler

import (
	"errors"
	"strconv"

	"cowork/internal/dto/request"
	"cowork/internal/dto/response"
	"cowork/internal/service"
	"cowork/pkg/errcode"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProjectHandler handles HTTP requests for project management.
type ProjectHandler struct {
	Svc *service.ProjectService
}

// Create handles POST /api/projects (PM only via middleware).
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	var req request.CreateProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	project, err := h.Svc.Create(userID, req)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, project)
}

// ListMyProjects handles GET /api/projects?page=&page_size=.
func (h *ProjectHandler) ListMyProjects(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	projects, total, err := h.Svc.ListMyProjects(userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, response.PageData{
		List:     projects,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetDetail handles GET /api/projects/:id. Returns the project together with
// its member list.
func (h *ProjectHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	project, err := h.Svc.GetDetail(projectID)
	if err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	// Fetch members separately (not cached) and return a combined response
	members, _ := h.Svc.ListMembers(projectID)

	response.Success(c, gin.H{
		"project": project,
		"members": members,
	})
}

// Update handles PUT /api/projects/:id. The caller must be the owner or a PM.
func (h *ProjectHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	var req request.UpdateProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	project, err := h.Svc.Update(projectID, userID, req)
	if err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsProjectVersionConflict(err) {
			response.Error(c, errcode.ErrProjectVersion, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, project)
}

// SubmitForApproval handles POST /api/projects/:id/submit.
func (h *ProjectHandler) SubmitForApproval(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	if err := h.Svc.SubmitForApproval(projectID, userID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		if service.IsInvalidStatusTrans(err) {
			response.Error(c, errcode.ErrInvalidStatusTrans, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// Approve handles POST /api/projects/:id/approve (PM only via middleware).
func (h *ProjectHandler) Approve(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	if err := h.Svc.Approve(projectID, userID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsInvalidStatusTrans(err) {
			response.Error(c, errcode.ErrInvalidStatusTrans, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// StartDevelopment handles POST /api/projects/:id/start.
func (h *ProjectHandler) StartDevelopment(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	if err := h.Svc.StartDevelopment(projectID, userID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		if service.IsInvalidStatusTrans(err) {
			response.Error(c, errcode.ErrInvalidStatusTrans, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// Complete handles POST /api/projects/:id/complete.
func (h *ProjectHandler) Complete(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	if err := h.Svc.Complete(projectID, userID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		if service.IsInvalidStatusTrans(err) {
			response.Error(c, errcode.ErrInvalidStatusTrans, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// Archive handles POST /api/projects/:id/archive.
func (h *ProjectHandler) Archive(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	if err := h.Svc.Archive(projectID, userID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		if service.IsInvalidStatusTrans(err) {
			response.Error(c, errcode.ErrInvalidStatusTrans, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// AddMember handles POST /api/projects/:id/members.
func (h *ProjectHandler) AddMember(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	var req request.AddMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	if err := h.Svc.AddMember(projectID, userID, req.UserID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		if errors.Is(err, service.ErrMemberAlreadyExists) {
			response.Error(c, errcode.ErrConflict, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// RemoveMember handles DELETE /api/projects/:id/members/:uid.
func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	uidStr := c.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid user id")
		return
	}
	memberID := uint(uid)

	if err := h.Svc.RemoveMember(projectID, userID, memberID); err != nil {
		if service.IsProjectNotFound(err) {
			response.Error(c, errcode.ErrProjectNotFound, err.Error())
			return
		}
		if service.IsNotProjectOwner(err) {
			response.Error(c, errcode.ErrNotProjectOwner, err.Error())
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListMembers handles GET /api/projects/:id/members.
func (h *ProjectHandler) ListMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	projectID := uint(id)

	members, err := h.Svc.ListMembers(projectID)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, members)
}

// Ensure gorm is used.
var _ = gorm.ErrRecordNotFound
