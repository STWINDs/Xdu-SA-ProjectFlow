package handler

import (
	"strconv"

	"cowork/internal/service"
	"cowork/pkg/errcode"

	"cowork/internal/dto/response"

	"github.com/gin-gonic/gin"
)

// LogHandler handles HTTP requests for operation logs. Only ProjectManager should access these endpoints.
type LogHandler struct {
	Svc *service.LogService
}

// List handles GET /api/logs?user_id=1&action=upload_attachment&page=1&page_size=10
func (h *LogHandler) List(c *gin.Context) {
	// Parse optional user_id filter
	var userID *uint
	if uidStr := c.Query("user_id"); uidStr != "" {
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			response.Error(c, errcode.ErrBadRequest, "invalid user_id")
			return
		}
		uidUint := uint(uid)
		userID = &uidUint
	}

	action := c.Query("action")

	page, pageSize := parsePagination(c)

	list, total, err := h.Svc.List(userID, action, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, response.PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
