package service

import (
	"cowork/internal/model"
	"cowork/internal/repository"

	"gorm.io/gorm"
)

// NotificationService handles business logic for notifications.
type NotificationService struct {
	DB *gorm.DB
}

// CreateNotification creates a new notification for the given user.
func (s *NotificationService) CreateNotification(userID uint, content, nType string, relatedID uint) error {
	n := &model.Notification{
		UserID:    userID,
		Content:   content,
		Type:      nType,
		RelatedID: relatedID,
		IsRead:    false,
	}
	return repository.CreateNotification(n)
}

// List returns paginated notifications for a user, ordered by most recent first.
func (s *NotificationService) List(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	return repository.ListNotificationsByUser(userID, page, pageSize)
}

// GetUnreadCount returns the number of unread notifications for a user.
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return repository.CountUnread(userID)
}

// MarkRead marks a single notification as read, ensuring the caller owns it.
func (s *NotificationService) MarkRead(userID, notificationID uint) error {
	if err := repository.MarkRead(notificationID, userID); err != nil {
		return err
	}
	return nil
}

// MarkAllRead marks all notifications as read for the given user.
func (s *NotificationService) MarkAllRead(userID uint) error {
	return repository.MarkAllRead(userID)
}
