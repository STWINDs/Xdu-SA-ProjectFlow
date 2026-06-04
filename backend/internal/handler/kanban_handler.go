package handler

import (
	"net/http"
	"strconv"
	"time"

	"cowork/internal/db"
	"cowork/internal/dto/response"
	"cowork/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KanbanHandler struct {
	DB *gorm.DB
}

type KanbanResponse struct {
	Columns []KanbanColumn `json:"columns"`
}

type KanbanColumn struct {
	Status string       `json:"status"`
	Title  string       `json:"title"`
	Tasks  []KanbanTask `json:"tasks"`
}

type KanbanTask struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Priority     string     `json:"priority"`
	Status       string     `json:"status"`
	AssigneeID   *uint      `json:"assignee_id"`
	AssigneeName string     `json:"assignee_name,omitempty"`
	Deadline     *time.Time `json:"deadline"`
	Version      int        `json:"version"`
}

var kanbanStatuses = []struct {
	Status string
	Title  string
}{
	{"Todo", "Todo"},
	{"InProgress", "In Progress"},
	{"Review", "Review"},
	{"Testing", "Testing"},
	{"Done", "Done"},
}

func (h *KanbanHandler) GetKanban(c *gin.Context) {
	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		response.Error(c, 90002, "invalid project id")
		return
	}

	// Verify project exists
	var project model.Project
	if err := db.DB.First(&project, pid).Error; err != nil {
		response.Error(c, 90005, "project not found")
		return
	}

	// Build query for tasks in this project
	query := db.DB.Where("project_id = ?", uint(pid)).Preload("Assignee")

	// Optional filters
	if assigneeIDStr := c.Query("assignee_id"); assigneeIDStr != "" {
		if aid, err := strconv.ParseUint(assigneeIDStr, 10, 64); err == nil {
			query = query.Where("assignee_id = ?", uint(aid))
		}
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	var tasks []model.Task
	if err := query.Find(&tasks).Error; err != nil {
		response.Error(c, 90001, "failed to query tasks")
		return
	}

	// Group tasks by status
	taskMap := make(map[string][]KanbanTask)
	for _, t := range tasks {
		kt := KanbanTask{
			ID:         t.ID,
			Title:      t.Title,
			Priority:   t.Priority,
			Status:     t.Status,
			AssigneeID: t.AssigneeID,
			Deadline:   t.Deadline,
			Version:    t.Version,
		}
		if t.Assignee != nil {
			kt.AssigneeName = t.Assignee.Username
		}
		taskMap[t.Status] = append(taskMap[t.Status], kt)
	}

	// Build columns (always return all 5)
	columns := make([]KanbanColumn, 0, 5)
	for _, ks := range kanbanStatuses {
		colTasks := taskMap[ks.Status]
		if colTasks == nil {
			colTasks = []KanbanTask{}
		}
		columns = append(columns, KanbanColumn{
			Status: ks.Status,
			Title:  ks.Title,
			Tasks:  colTasks,
		})
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "ok",
		Data:    KanbanResponse{Columns: columns},
	})
}
