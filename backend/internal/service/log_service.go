package service

import (
	"log"

	"cowork/internal/db"
	"cowork/internal/model"
	"cowork/internal/repository"

	"gorm.io/gorm"
)

// LogService handles creation and querying of operation logs.
type LogService struct {
	DB *gorm.DB
}

// Log creates an operation log entry. This is fire-and-forget: errors are logged
// to stderr but the method always returns nil so callers are not blocked.
func (s *LogService) Log(userID uint, action, ip, detail string) {
	logEntry := &model.OperationLog{
		UserID: userID,
		Action: action,
		IP:     ip,
		Detail: detail,
	}
	if err := db.DB.Create(logEntry).Error; err != nil {
		log.Printf("Warning: failed to create operation log (user=%d action=%s): %v", userID, action, err)
	}
}

// List returns paginated operation logs with optional filters. Intended for ProjectManager use.
func (s *LogService) List(userID *uint, action string, page, pageSize int) ([]model.OperationLog, int64, error) {
	return repository.ListLogs(page, pageSize, userID, action)
}
