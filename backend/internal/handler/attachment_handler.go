package handler

import (
	"strconv"

	"cowork/internal/service"
	"cowork/pkg/errcode"

	"cowork/internal/dto/response"

	"github.com/gin-gonic/gin"
)

// AttachmentHandler handles HTTP requests for file attachments.
type AttachmentHandler struct {
	Svc *service.AttachmentService
}

// Upload handles POST /api/attachments/upload
// Accepts multipart form: file (required), task_id (optional), project_id (optional).
func (h *AttachmentHandler) Upload(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "file is required")
		return
	}

	// Parse optional task_id
	var taskID *uint
	if tidStr := c.PostForm("task_id"); tidStr != "" {
		tid, parseErr := strconv.ParseUint(tidStr, 10, 64)
		if parseErr != nil {
			response.Error(c, errcode.ErrBadRequest, "invalid task_id")
			return
		}
		tidUint := uint(tid)
		taskID = &tidUint
	}

	// Parse optional project_id
	var projectID *uint
	if pidStr := c.PostForm("project_id"); pidStr != "" {
		pid, parseErr := strconv.ParseUint(pidStr, 10, 64)
		if parseErr != nil {
			response.Error(c, errcode.ErrBadRequest, "invalid project_id")
			return
		}
		pidUint := uint(pid)
		projectID = &pidUint
	}

	att, err := h.Svc.Upload(userID, file, taskID, projectID)
	if err != nil {
		msg := err.Error()
		switch {
		case msg == "file size exceeds 20MB limit":
			response.Error(c, errcode.ErrFileTooLarge, msg)
		case msg[:len("file type not allowed:")] == "file type not allowed:":
			response.Error(c, errcode.ErrFileTypeDenied, msg)
		default:
			response.Error(c, errcode.ErrUploadFailed, msg)
		}
		return
	}

	response.Success(c, att)
}

// ListByTask handles GET /api/tasks/:tid/attachments
func (h *AttachmentHandler) ListByTask(c *gin.Context) {
	tidStr := c.Param("id")
	tid, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	list, err := h.Svc.ListByTask(uint(tid))
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, list)
}

// ListByProject handles GET /api/projects/:pid/attachments
func (h *AttachmentHandler) ListByProject(c *gin.Context) {
	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}

	list, err := h.Svc.ListByProject(uint(pid))
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, list)
}

// Delete handles DELETE /api/attachments/:id
func (h *AttachmentHandler) Delete(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid attachment id")
		return
	}

	if err := h.Svc.Delete(userID, uint(id)); err != nil {
		msg := err.Error()
		switch {
		case msg == "attachment not found":
			response.Error(c, errcode.ErrNotFound, msg)
		case msg == "permission denied: you do not own this attachment":
			response.Error(c, errcode.ErrForbidden, msg)
		default:
			response.Error(c, errcode.ErrInternal, msg)
		}
		return
	}

	response.Success(c, nil)
}
