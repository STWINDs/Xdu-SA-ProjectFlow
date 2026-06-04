package repository

import (
	"cowork/internal/db"
	"cowork/internal/model"

	"gorm.io/gorm"
)

func CreateComment(c *model.Comment) error {
	return db.DB.Create(c).Error
}

func FindCommentByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := db.DB.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func ListCommentsByProject(projectID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := db.DB.Where("project_id = ? AND parent_id IS NULL", projectID)

	if err := query.Model(&model.Comment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func ListCommentsByTask(taskID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := db.DB.Where("task_id = ? AND parent_id IS NULL", taskID)

	if err := query.Model(&model.Comment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func DeleteComment(id uint) error {
	// Delete all replies first
	if err := db.DB.Where("parent_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
		return err
	}
	// Delete the comment itself
	return db.DB.Delete(&model.Comment{}, id).Error
}
