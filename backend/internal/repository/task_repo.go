package repository

import (
	"cowork/internal/db"
	"cowork/internal/model"
	"time"
)

func CreateTask(t *model.Task) error {
	return db.DB.Create(t).Error
}

func FindTaskByID(id uint) (*model.Task, error) {
	var t model.Task
	err := db.DB.Preload("Assignee").Preload("Creator").Preload("Project").First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(t *model.Task) error {
	return db.DB.Save(t).Error
}

func UpdateTaskWithVersion(t *model.Task) (bool, error) {
	currentVersion := t.Version
	t.Version = currentVersion + 1

	now := time.Now()
	result := db.DB.Exec(
		"UPDATE tasks SET title=?, description=?, priority=?, status=?, assignee_id=?, deadline=?, version=?, updated_at=? WHERE id=? AND version=?",
		t.Title, t.Description, t.Priority, t.Status, t.AssigneeID, t.Deadline, t.Version, now, t.ID, currentVersion,
	)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func DeleteTask(id uint) error {
	return db.DB.Delete(&model.Task{}, id).Error
}

func ListTasksByProject(projectID uint, status, priority string, assigneeID *uint, page, pageSize int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := db.DB.Model(&model.Task{}).Where("project_id = ?", projectID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Assignee").Preload("Creator").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func ListTasksByAssignee(assigneeID uint, page, pageSize int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := db.DB.Model(&model.Task{}).Where("assignee_id = ?", assigneeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Project").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
