package service

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"cowork/internal/db"
	"cowork/internal/model"

	"gorm.io/gorm"
)

// MaxFileSize is the maximum allowed upload size: 20 MB.
const MaxFileSize = 20 * 1024 * 1024 // 20971520 bytes

// UploadDir is the local directory where uploaded files are stored.
const UploadDir = "./uploads"

// allowedContentTypes defines the MIME types permitted for upload.
// Patterns ending with "/*" match any subtype within that MIME type.
var allowedContentTypes = []string{
	"image/*",
	"application/pdf",
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/zip",
	"application/x-rar-compressed",
	"text/plain",
}

// AttachmentService handles file upload, listing, and deletion.
type AttachmentService struct {
	DB *gorm.DB
}

// Upload validates and saves a file, then creates an Attachment record.
// taskID and projectID are optional pointers; at least one should be non-nil in practice.
func (s *AttachmentService) Upload(userID uint, file *multipart.FileHeader, taskID *uint, projectID *uint) (*model.Attachment, error) {
	// 1. Check file size
	if file.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds 20MB limit")
	}

	// 2. Check content type
	contentType := file.Header.Get("Content-Type")
	if !isContentTypeAllowed(contentType) {
		return nil, fmt.Errorf("file type not allowed: %s", contentType)
	}

	// 3. Generate unique filename
	uniqueID, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique filename: %w", err)
	}
	_ = filepath.Ext(file.Filename)
	savedName := fmt.Sprintf("%s_%s", uniqueID, file.Filename)
	savedPath := filepath.Join(UploadDir, savedName)

	// 4. Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 5. Save file to disk
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(savedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 6. Create Attachment record
	att := &model.Attachment{
		FileName:    file.Filename,
		FileURL:     "/uploads/" + savedName,
		FileSize:    file.Size,
		ContentType: contentType,
		UserID:      userID,
		TaskID:      taskID,
		ProjectID:   projectID,
	}

	if err := s.DB.Create(att).Error; err != nil {
		// Clean up the saved file on DB failure
		os.Remove(savedPath)
		return nil, fmt.Errorf("failed to create attachment record: %w", err)
	}

	// 7. Log operation (fire and forget)
	_ = db.DB.Create(&model.OperationLog{
		UserID: userID,
		Action: "upload_attachment",
		Detail: fmt.Sprintf("Uploaded file: %s (%d bytes)", file.Filename, file.Size),
	})

	return att, nil
}

// ListByTask returns all attachments linked to the given task.
func (s *AttachmentService) ListByTask(taskID uint) ([]model.Attachment, error) {
	var list []model.Attachment
	err := s.DB.Where("task_id = ?", taskID).Order("created_at DESC").Find(&list).Error
	return list, err
}

// ListByProject returns all attachments linked to the given project.
func (s *AttachmentService) ListByProject(projectID uint) ([]model.Attachment, error) {
	var list []model.Attachment
	err := s.DB.Where("project_id = ?", projectID).Order("created_at DESC").Find(&list).Error
	return list, err
}

// Delete removes an attachment after verifying the requesting user owns it.
// It deletes both the database record and the file on disk.
func (s *AttachmentService) Delete(requesterID uint, attachmentID uint) error {
	var att model.Attachment
	if err := s.DB.First(&att, attachmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("attachment not found")
		}
		return err
	}

	if att.UserID != requesterID {
		return fmt.Errorf("permission denied: you do not own this attachment")
	}

	// Delete file from disk
	filePath := filepath.Join(".", strings.TrimPrefix(att.FileURL, "/"))
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: failed to delete file %s from disk: %v", filePath, err)
	}

	// Delete DB record
	if err := s.DB.Delete(&att).Error; err != nil {
		return err
	}

	// Log operation (fire and forget)
	_ = db.DB.Create(&model.OperationLog{
		UserID: requesterID,
		Action: "delete_attachment",
		Detail: fmt.Sprintf("Deleted attachment: %s", att.FileName),
	})

	return nil
}

// isContentTypeAllowed checks whether the given MIME type matches any entry in the allowed list.
// Patterns ending with "/*" match any subtype within that MIME type (e.g., "image/*").
func isContentTypeAllowed(contentType string) bool {
	if contentType == "" {
		return false
	}

	for _, allowed := range allowedContentTypes {
		if strings.HasSuffix(allowed, "/*") {
			prefix := strings.TrimSuffix(allowed, "*")
			if strings.HasPrefix(contentType, prefix) {
				return true
			}
		} else if contentType == allowed {
			return true
		}
	}
	return false
}

// generateUUID creates a random UUID-like string in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func generateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Set version 4 bits (RFC 4122 variant)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
