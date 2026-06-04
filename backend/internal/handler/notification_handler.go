package handler

import (
	"strconv"

	"cowork/internal/service"
	"cowork/pkg/errcode"

	"cowork/internal/dto/response"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles HTTP requests for notifications.
type NotificationHandler struct {
	Svc *service.NotificationService
}

// List handles GET /api/notifications?page=1&page_size=10
func (h *NotificationHandler) List(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	page, pageSize := parsePagination(c)

	list, total, err := h.Svc.List(userID, page, pageSize)
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

// GetUnreadCount handles GET /api/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	count, err := h.Svc.GetUnreadCount(userID)
	if err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, gin.H{"unread_count": count})
}

// MarkRead handles PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid notification id")
		return
	}

	if err := h.Svc.MarkRead(userID, uint(id)); err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// MarkAllRead handles PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	if err := h.Svc.MarkAllRead(userID); err != nil {
		response.Error(c, errcode.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// parsePagination extracts page and pageSize from query parameters, providing sensible defaults.
func parsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	return page, pageSize
}
