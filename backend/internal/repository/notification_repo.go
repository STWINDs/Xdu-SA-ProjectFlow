package repository

import (
	"cowork/internal/db"
	"cowork/internal/model"
)

// CreateNotification inserts a new notification record into the database.
func CreateNotification(n *model.Notification) error {
	return db.DB.Create(n).Error
}

// ListNotificationsByUser returns paginated notifications for a user, ordered by most recent first.
func ListNotificationsByUser(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var list []model.Notification
	var total int64

	query := db.DB.Model(&model.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// CountUnread returns the number of unread notifications for a user.
func CountUnread(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// MarkRead marks a single notification as read, scoped to the owner user.
func MarkRead(id, userID uint) error {
	return db.DB.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

// MarkAllRead marks all notifications as read for a given user.
func MarkAllRead(userID uint) error {
	return db.DB.Model(&model.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}
