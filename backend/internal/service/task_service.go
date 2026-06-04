package service

import (
	"errors"
	"strconv"
	"time"

	"cowork/internal/cache"
	"cowork/internal/dto/request"
	"cowork/internal/model"
	"cowork/internal/repository"

	"gorm.io/gorm"
)

type TaskService struct {
	DB *gorm.DB
}

var validStatusTransitions = map[string]string{
	"Todo":       "InProgress",
	"InProgress": "Review",
	"Review":     "Testing",
	"Testing":    "Done",
}

func (s *TaskService) taskDetailCacheKey(taskID uint) string {
	return "task:detail:" + strconv.FormatUint(uint64(taskID), 10)
}

func (s *TaskService) invalidateTaskCache(taskID uint) {
	_ = cache.Del(s.taskDetailCacheKey(taskID))
	_ = cache.DeletePattern("task:list:*")
}

// Create creates a new task. creatorID is the user creating the task.
func (s *TaskService) Create(creatorID uint, req request.CreateTaskReq) (*model.Task, error) {
	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      "Todo",
		ProjectID:   req.ProjectID,
		AssigneeID:  req.AssigneeID,
		CreatorID:   creatorID,
		Deadline:    req.Deadline,
		Version:     1,
	}
	if err := repository.CreateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

// Update updates a task with optimistic locking. userID is checked to be the creator or assignee.
func (s *TaskService) Update(userID, taskID uint, req request.UpdateTaskReq) (*model.Task, error) {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// Verify user is creator or assignee
	if task.CreatorID != userID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		return nil, errors.New("forbidden")
	}

	// Apply optimistic lock: the client must send the current version
	task.Title = req.Title
	task.Description = req.Description
	task.Priority = req.Priority
	task.AssigneeID = req.AssigneeID
	task.Deadline = req.Deadline
	task.Version = req.Version

	ok, err := repository.UpdateTaskWithVersion(task)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("version conflict")
	}

	s.invalidateTaskCache(taskID)

	// Re-fetch to get the updated record
	return repository.FindTaskByID(taskID)
}

// Delete deletes a task. Only the creator can delete.
func (s *TaskService) Delete(userID, taskID uint) error {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return err
	}
	if task.CreatorID != userID {
		return errors.New("forbidden")
	}

	if err := repository.DeleteTask(taskID); err != nil {
		return err
	}

	s.invalidateTaskCache(taskID)
	return nil
}

// Assign assigns a task to a user. assignerID is the user performing the assignment.
func (s *TaskService) Assign(taskID, assignerID, assigneeID uint) (*model.Task, error) {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	task.AssigneeID = &assigneeID

	ok, err := repository.UpdateTaskWithVersion(task)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("version conflict")
	}

	s.invalidateTaskCache(taskID)

	return repository.FindTaskByID(taskID)
}

// TransferStatus transitions a task to a new status. userID is checked to be the assignee or creator.
func (s *TaskService) TransferStatus(userID, taskID uint, newStatus string) (*model.Task, error) {
	task, err := repository.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// Verify user is assignee or creator
	if task.CreatorID != userID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		return nil, errors.New("forbidden")
	}

	// Validate transition
	expectedNext, ok := validStatusTransitions[task.Status]
	if !ok {
		return nil, errors.New("invalid current status")
	}
	if newStatus != expectedNext {
		return nil, errors.New("invalid status transition")
	}

	task.Status = newStatus

	ok, err = repository.UpdateTaskWithVersion(task)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("version conflict")
	}

	s.invalidateTaskCache(taskID)

	return repository.FindTaskByID(taskID)
}

// GetDetail retrieves a task by ID with Cache-Aside pattern.
func (s *TaskService) GetDetail(taskID uint) (*model.Task, error) {
	cacheKey := s.taskDetailCacheKey(taskID)

	var task model.Task
	if err := cache.Get(cacheKey, &task); err == nil {
		return &task, nil
	}

	t, err := repository.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	_ = cache.Set(cacheKey, t, 10*time.Minute)

	return t, nil
}

// ListByProject lists tasks for a project with optional filters.
func (s *TaskService) ListByProject(projectID uint, status, priority string, assigneeID *uint, page, pageSize int) ([]model.Task, int64, error) {
	return repository.ListTasksByProject(projectID, status, priority, assigneeID, page, pageSize)
}

// IsVersionConflict checks if an error is a version conflict.
func IsVersionConflict(err error) bool {
	return err != nil && err.Error() == "version conflict"
}

// IsForbidden checks if an error is a forbidden error.
func IsForbidden(err error) bool {
	return err != nil && err.Error() == "forbidden"
}

// IsInvalidTransition checks if an error is an invalid status transition.
func IsInvalidTransition(err error) bool {
	return err != nil && err.Error() == "invalid status transition"
}
