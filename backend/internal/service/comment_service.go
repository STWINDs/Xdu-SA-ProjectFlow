package service

import (
	"fmt"
	"regexp"

	"cowork/internal/db"
	"cowork/internal/dto/request"
	"cowork/internal/model"
	"cowork/internal/repository"
	"cowork/pkg/errcode"

	"gorm.io/gorm"
)

type CommentService struct {
	DB *gorm.DB
}

var mentionRegex = regexp.MustCompile(`@(\w+)`)

func (s *CommentService) CreateComment(userID uint, req request.CreateCommentReq) (*model.Comment, error) {
	// Validate that at least one of ProjectID/TaskID is set
	if req.ProjectID == nil && req.TaskID == nil {
		return nil, fmt.Errorf("project_id or task_id is required")
	}

	comment := &model.Comment{
		Content:   req.Content,
		UserID:    userID,
		ProjectID: req.ProjectID,
		TaskID:    req.TaskID,
	}

	if err := repository.CreateComment(comment); err != nil {
		return nil, err
	}

	// Parse @mentions and create notifications
	s.processMentions(userID, comment)

	// Reload with User preloaded
	return repository.FindCommentByID(comment.ID)
}

func (s *CommentService) CreateReply(userID uint, parentID uint, req request.CreateReplyReq) (*model.Comment, error) {
	// Find parent comment
	parent, err := repository.FindCommentByID(parentID)
	if err != nil {
		return nil, fmt.Errorf("parent comment not found")
	}

	comment := &model.Comment{
		Content:   req.Content,
		UserID:    userID,
		ProjectID: parent.ProjectID,
		TaskID:    parent.TaskID,
		ParentID:  &parentID,
	}

	if err := repository.CreateComment(comment); err != nil {
		return nil, err
	}

	// Parse @mentions and create notifications
	s.processMentions(userID, comment)

	// Notify parent comment author
	var currentUser model.User
	if err := db.DB.First(&currentUser, userID).Error; err == nil {
		if parent.UserID != userID {
			notif := &model.Notification{
				UserID:    parent.UserID,
				Content:   fmt.Sprintf("%s replied to your comment", currentUser.Username),
				Type:      "CommentReplied",
				RelatedID: comment.ID,
			}
			db.DB.Create(notif)
		}
	}

	return repository.FindCommentByID(comment.ID)
}

func (s *CommentService) processMentions(userID uint, comment *model.Comment) {
	var currentUser model.User
	if err := db.DB.First(&currentUser, userID).Error; err != nil {
		return
	}

	matches := mentionRegex.FindAllStringSubmatch(comment.Content, -1)
	mentionedUsers := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			username := match[1]
			if username == currentUser.Username {
				continue
			}
			if mentionedUsers[username] {
				continue
			}
			mentionedUsers[username] = true

			var user model.User
			if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
				continue
			}

			notif := &model.Notification{
				UserID:    user.ID,
				Content:   fmt.Sprintf("%s mentioned you in a comment", currentUser.Username),
				Type:      "CommentReplied",
				RelatedID: comment.ID,
			}
			db.DB.Create(notif)
		}
	}
}

func (s *CommentService) ListByProject(projectID uint, page, pageSize int) ([]model.Comment, int64, error) {
	return repository.ListCommentsByProject(projectID, page, pageSize)
}

func (s *CommentService) ListByTask(taskID uint, page, pageSize int) ([]model.Comment, int64, error) {
	return repository.ListCommentsByTask(taskID, page, pageSize)
}

func (s *CommentService) Delete(userID, commentID uint) error {
	comment, err := repository.FindCommentByID(commentID)
	if err != nil {
		return fmt.Errorf("comment not found, error code: %d", errcode.ErrCommentNotFound)
	}

	if comment.UserID != userID {
		return fmt.Errorf("you can only delete your own comments, error code: %d", errcode.ErrBadRequest)
	}

	return repository.DeleteComment(commentID)
}
