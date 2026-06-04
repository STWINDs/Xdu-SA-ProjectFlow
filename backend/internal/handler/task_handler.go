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

type TaskHandler struct {
	Svc *service.TaskService
}

// Create handles POST /api/projects/:pid/tasks
func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	var req request.CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}
	req.ProjectID = uint(pid)

	task, err := h.Svc.Create(userID, req)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, task)
}

// ListByProject handles GET /api/projects/:pid/tasks
func (h *TaskHandler) ListByProject(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}

	status := c.Query("status")
	priority := c.Query("priority")

	var assigneeID *uint
	if assigneeStr := c.Query("assignee_id"); assigneeStr != "" {
		id, err := strconv.ParseUint(assigneeStr, 10, 64)
		if err != nil {
			response.Error(c, errcode.ErrBadRequest, "invalid assignee_id")
			return
		}
		uid := uint(id)
		assigneeID = &uid
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

	tasks, total, err := h.Svc.ListByProject(uint(pid), status, priority, assigneeID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, response.PageData{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetDetail handles GET /api/tasks/:id
func (h *TaskHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	task, err := h.Svc.GetDetail(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, errcode.ErrTaskNotFound, "task not found")
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, task)
}

// Update handles PUT /api/tasks/:id
func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	var req request.UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	task, err := h.Svc.Update(userID, uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, errcode.ErrTaskNotFound, "task not found")
			return
		}
		if service.IsVersionConflict(err) {
			response.Error(c, errcode.ErrTaskVersionConflict, "version conflict: the task has been modified by another user, please refresh and try again")
			return
		}
		if service.IsForbidden(err) {
			response.Error(c, errcode.ErrForbidden, "only the creator or assignee can update this task")
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, task)
}

// Delete handles DELETE /api/tasks/:id
func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	err = h.Svc.Delete(userID, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, errcode.ErrTaskNotFound, "task not found")
			return
		}
		if service.IsForbidden(err) {
			response.Error(c, errcode.ErrForbidden, "only the creator can delete this task")
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// Assign handles PUT /api/tasks/:id/assign
func (h *TaskHandler) Assign(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	var req request.AssignTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	task, err := h.Svc.Assign(uint(id), userID, req.AssigneeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, errcode.ErrTaskNotFound, "task not found")
			return
		}
		if service.IsVersionConflict(err) {
			response.Error(c, errcode.ErrTaskVersionConflict, "version conflict: the task has been modified by another user, please refresh and try again")
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, task)
}

// TransferStatus handles PUT /api/tasks/:id/status
func (h *TaskHandler) TransferStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	var req request.TransferStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	task, err := h.Svc.TransferStatus(userID, uint(id), req.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, errcode.ErrTaskNotFound, "task not found")
			return
		}
		if service.IsForbidden(err) {
			response.Error(c, errcode.ErrForbidden, "only the creator or assignee can transfer the status")
			return
		}
		if service.IsInvalidTransition(err) {
			response.Error(c, errcode.ErrInvalidTaskStatus, "invalid status transition")
			return
		}
		if service.IsVersionConflict(err) {
			response.Error(c, errcode.ErrTaskVersionConflict, "version conflict: the task has been modified by another user, please refresh and try again")
			return
		}
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, task)
}
