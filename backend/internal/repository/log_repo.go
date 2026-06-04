package repository

import (
	"cowork/internal/db"
	"cowork/internal/model"
)

// CreateLog inserts a new operation log record into the database.
func CreateLog(log *model.OperationLog) error {
	return db.DB.Create(log).Error
}

// ListLogs returns paginated operation logs with optional user and action filters.
// Results are ordered by most recent first and include the associated User.
func ListLogs(page, pageSize int, userID *uint, action string) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	var total int64

	query := db.DB.Model(&model.OperationLog{}).Preload("User")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}
