package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"cowork/internal/cache"
	"cowork/internal/db"
	"cowork/internal/dto/request"
	"cowork/internal/model"
	"cowork/internal/repository"

	"gorm.io/gorm"
)

// Sentinel errors for project operations. Handlers check these with IsXxx helpers.
var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrInvalidStatusTrans  = errors.New("invalid status transition")
	ErrNotProjectOwner     = errors.New("not project owner")
	ErrProjectVersion      = errors.New("version conflict")
	ErrMemberAlreadyExists = errors.New("member already exists")
)

// ProjectService handles project business logic.
type ProjectService struct {
	DB *gorm.DB
}

// --- Error helpers so handlers can branch on specific failures ---

func IsProjectNotFound(err error) bool       { return errors.Is(err, ErrProjectNotFound) }
func IsInvalidStatusTrans(err error) bool    { return errors.Is(err, ErrInvalidStatusTrans) }
func IsNotProjectOwner(err error) bool       { return errors.Is(err, ErrNotProjectOwner) }
func IsProjectVersionConflict(err error) bool { return errors.Is(err, ErrProjectVersion) }

// --- Cache helpers ---

func (s *ProjectService) projectDetailCacheKey(projectID uint) string {
	return "project:detail:" + strconv.FormatUint(uint64(projectID), 10)
}

func (s *ProjectService) invalidateProjectCache(projectID uint) {
	_ = cache.Del(s.projectDetailCacheKey(projectID))
	_ = cache.DeletePattern("project:stats:*")
}

// Create creates a new project owned by userID with status "Draft" and version 1.
func (s *ProjectService) Create(userID uint, req request.CreateProjectReq) (*model.Project, error) {
	project := &model.Project{
		Name:        req.Name,
		Description: req.Description,
		Status:      "Draft",
		OwnerID:     userID,
		Version:     1,
	}

	if err := repository.CreateProject(project); err != nil {
		return nil, err
	}
	return project, nil
}

// Update modifies a project's name and description. The caller must be the owner
// or have the "ProjectManager" role. Optimistic locking is enforced via req.Version.
func (s *ProjectService) Update(projectID, userID uint, req request.UpdateProjectReq) (*model.Project, error) {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	// Check owner or PM role
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, ErrProjectNotFound // user not found, generic
	}
	if p.OwnerID != userID && user.Role != "ProjectManager" {
		return nil, ErrNotProjectOwner
	}

	// Apply optimistic lock
	p.Name = req.Name
	p.Description = req.Description
	p.Version = req.Version

	updated, err := repository.UpdateProjectWithVersion(p)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, ErrProjectVersion
	}

	s.invalidateProjectCache(projectID)

	return repository.FindProjectByID(projectID)
}

// SubmitForApproval transitions the project from Draft to PendingApproval.
// Only the project owner may submit.
func (s *ProjectService) SubmitForApproval(projectID, userID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.OwnerID != userID {
		return ErrNotProjectOwner
	}
	if p.Status != "Draft" {
		return ErrInvalidStatusTrans
	}

	p.Status = "PendingApproval"
	if err := repository.UpdateProject(p); err != nil {
		return err
	}

	s.invalidateProjectCache(projectID)
	return nil
}

// Approve transitions the project from PendingApproval to Approved.
// The caller must have the ProjectManager role (enforced by middleware).
// An operation log entry is written.
func (s *ProjectService) Approve(projectID, userID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.Status != "PendingApproval" {
		return ErrInvalidStatusTrans
	}

	p.Status = "Approved"
	if err := repository.UpdateProject(p); err != nil {
		return err
	}

	// Record operation log
	log := &model.OperationLog{
		UserID: userID,
		Action: "approve_project",
		Detail: fmt.Sprintf("Approved project: %s (ID: %d)", p.Name, p.ID),
	}
	_ = repository.CreateLog(log)

	s.invalidateProjectCache(projectID)
	return nil
}

// StartDevelopment transitions the project from Approved to Developing.
// Only the project owner may start development.
func (s *ProjectService) StartDevelopment(projectID, userID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.OwnerID != userID {
		return ErrNotProjectOwner
	}
	if p.Status != "Approved" {
		return ErrInvalidStatusTrans
	}

	p.Status = "Developing"
	if err := repository.UpdateProject(p); err != nil {
		return err
	}

	s.invalidateProjectCache(projectID)
	return nil
}

// Complete transitions the project from Developing to Completed.
// Only the project owner may complete.
func (s *ProjectService) Complete(projectID, userID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.OwnerID != userID {
		return ErrNotProjectOwner
	}
	if p.Status != "Developing" {
		return ErrInvalidStatusTrans
	}

	p.Status = "Completed"
	if err := repository.UpdateProject(p); err != nil {
		return err
	}

	s.invalidateProjectCache(projectID)
	return nil
}

// Archive transitions the project from Completed to Archived.
// Only the project owner may archive.
func (s *ProjectService) Archive(projectID, userID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.OwnerID != userID {
		return ErrNotProjectOwner
	}
	if p.Status != "Completed" {
		return ErrInvalidStatusTrans
	}

	p.Status = "Archived"
	if err := repository.UpdateProject(p); err != nil {
		return err
	}

	s.invalidateProjectCache(projectID)
	return nil
}

// GetDetail retrieves a project by ID using the Cache-Aside pattern.
// Cache TTL is 10 minutes. On cache miss the project is loaded from the
// database with the Owner association preloaded. Callers should also call
// ListMembers separately to obtain member data.
func (s *ProjectService) GetDetail(projectID uint) (*model.Project, error) {
	cacheKey := s.projectDetailCacheKey(projectID)

	var p model.Project
	if err := cache.Get(cacheKey, &p); err == nil {
		return &p, nil
	}

	project, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	_ = cache.Set(cacheKey, project, 10*time.Minute)

	return project, nil
}

// ListMyProjects returns paginated projects where the user is owner or member.
// This list is not cached.
func (s *ProjectService) ListMyProjects(userID uint, page, pageSize int) ([]model.Project, int64, error) {
	return repository.ListProjectsByUser(userID, page, pageSize)
}

// AddMember adds a user as a member of a project. Only the project owner can
// add members. Duplicate membership is rejected.
func (s *ProjectService) AddMember(projectID, ownerID, memberID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.OwnerID != ownerID {
		return ErrNotProjectOwner
	}

	// Check for duplicate membership
	var existing model.ProjectMember
	if err := db.DB.Where("project_id = ? AND user_id = ?", projectID, memberID).First(&existing).Error; err == nil {
		return ErrMemberAlreadyExists
	}

	pm := &model.ProjectMember{
		ProjectID: projectID,
		UserID:    memberID,
		Role:      "Developer",
	}

	return repository.AddMember(pm)
}

// RemoveMember removes a user from a project. Only the project owner can
// remove members.
func (s *ProjectService) RemoveMember(projectID, ownerID, memberID uint) error {
	p, err := repository.FindProjectByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if p.OwnerID != ownerID {
		return ErrNotProjectOwner
	}

	return repository.RemoveMember(projectID, memberID)
}

// ListMembers returns all members of a project with the User association preloaded.
func (s *ProjectService) ListMembers(projectID uint) ([]model.ProjectMember, error) {
	return repository.ListMembers(projectID)
}

// Ensure fmt is used (for future use).
var _ = fmt.Sprintf
