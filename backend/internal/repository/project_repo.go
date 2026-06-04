package repository

import (
	"cowork/internal/db"
	"cowork/internal/model"
	"time"
)

// CreateProject inserts a new project into the database.
func CreateProject(p *model.Project) error {
	return db.DB.Create(p).Error
}

// FindProjectByID finds a project by primary key ID, preloading the Owner association.
func FindProjectByID(id uint) (*model.Project, error) {
	var p model.Project
	err := db.DB.Preload("Owner").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProject saves all fields of the project to the database.
func UpdateProject(p *model.Project) error {
	return db.DB.Save(p).Error
}

// UpdateProjectWithVersion performs an optimistic-lock update. It uses
// WHERE id=? AND version=? and increments the version. Returns true if a
// row was actually updated, false if the version no longer matches.
func UpdateProjectWithVersion(p *model.Project) (bool, error) {
	currentVersion := p.Version
	p.Version = currentVersion + 1

	now := time.Now()
	result := db.DB.Exec(
		"UPDATE projects SET name=?, description=?, status=?, owner_id=?, version=?, updated_at=? WHERE id=? AND version=?",
		p.Name, p.Description, p.Status, p.OwnerID, p.Version, now, p.ID, currentVersion,
	)
	return result.RowsAffected > 0, result.Error
}

// ListProjectsByUser returns paginated projects where the given user is
// either the owner or a member. Results are ordered by creation time
// descending and include the Owner association.
func ListProjectsByUser(userID uint, page, pageSize int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	// Subquery: project IDs where the user is a member
	memberSubQuery := db.DB.Model(&model.ProjectMember{}).Select("project_id").Where("user_id = ?", userID)

	query := db.DB.Model(&model.Project{}).Where("owner_id = ? OR id IN (?)", userID, memberSubQuery)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Owner").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// AddMember inserts a new project member record.
func AddMember(pm *model.ProjectMember) error {
	return db.DB.Create(pm).Error
}

// RemoveMember deletes a member record by project ID and user ID.
func RemoveMember(projectID, userID uint) error {
	return db.DB.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&model.ProjectMember{}).Error
}

// ListMembers returns all members of a project with the User association preloaded.
func ListMembers(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	err := db.DB.Preload("User").Where("project_id = ?", projectID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
