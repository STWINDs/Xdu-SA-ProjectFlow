package handler

import (
	"strconv"

	"cowork/internal/dto/request"
	"cowork/internal/dto/response"
	"cowork/internal/service"
	"cowork/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Svc *service.CommentService
}

func (h *CommentHandler) ListByProject(c *gin.Context) {
	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid project id")
		return
	}

	page, pageSize := parsePagination(c)

	comments, total, err := h.Svc.ListByProject(uint(pid), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, response.PageData{
		List:     comments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *CommentHandler) CreateByProject(c *gin.Context) {
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

	var req request.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	pidUint := uint(pid)
	req.ProjectID = &pidUint

	comment, err := h.Svc.CreateComment(userID, req)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, comment)
}

func (h *CommentHandler) ListByTask(c *gin.Context) {
	tidStr := c.Param("id")
	tid, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	page, pageSize := parsePagination(c)

	comments, total, err := h.Svc.ListByTask(uint(tid), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, response.PageData{
		List:     comments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *CommentHandler) CreateByTask(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	tidStr := c.Param("id")
	tid, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid task id")
		return
	}

	var req request.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	tidUint := uint(tid)
	req.TaskID = &tidUint

	comment, err := h.Svc.CreateComment(userID, req)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, comment)
}

func (h *CommentHandler) CreateReply(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	commentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid comment id")
		return
	}

	var req request.CreateReplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	reply, err := h.Svc.CreateReply(userID, uint(commentID), req)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, reply)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized, "unauthorized")
		return
	}

	idStr := c.Param("id")
	commentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid comment id")
		return
	}

	if err := h.Svc.Delete(userID, uint(commentID)); err != nil {
		response.Error(c, errcode.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}
