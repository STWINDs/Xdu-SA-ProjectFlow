package errcode

const (
	ErrInternal  = 90001
	ErrBadRequest = 90002
	ErrUnauthorized = 90003
	ErrForbidden    = 90004
	ErrNotFound     = 90005
	ErrConflict     = 90006

	ErrUserExists     = 10001
	ErrUserNotFound   = 10002
	ErrWrongPassword  = 10003
	ErrTokenInvalid   = 10004
	ErrTokenExpired   = 10005
	ErrCaptchaInvalid = 10006
	ErrLoginLocked    = 10007
	ErrWeakPassword   = 10008

	ErrProjectNotFound    = 20001
	ErrInvalidStatusTrans = 20002
	ErrNotProjectOwner    = 20003
	ErrProjectVersion     = 20004

	ErrTaskNotFound       = 30001
	ErrTaskVersionConflict = 30002
	ErrInvalidTaskStatus  = 30003

	ErrCommentNotFound = 40001

	ErrFileTooLarge   = 50001
	ErrFileTypeDenied = 50002
	ErrUploadFailed   = 50003

	ErrNotificationNotFound = 60001
)
