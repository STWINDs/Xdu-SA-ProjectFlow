package handler

import (
	"net/http"

	"cowork/internal/dto/request"
	"cowork/internal/dto/response"
	"cowork/internal/service"
	"cowork/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	Svc *service.AuthService
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	svcReq := service.RegisterReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.Svc.Register(svcReq)
	if err != nil {
		if appErr, ok := service.IsAppError(err); ok {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, errcode.ErrInternal, "internal server error")
		}
		return
	}

	response.Success(c, user)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	svcReq := service.LoginReq{
		Username:      req.Username,
		Password:      req.Password,
		CaptchaID:     req.CaptchaID,
		CaptchaAnswer: req.CaptchaAnswer,
	}

	resp, err := h.Svc.Login(svcReq, req.CaptchaID, req.CaptchaAnswer)
	if err != nil {
		if appErr, ok := service.IsAppError(err); ok {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, errcode.ErrInternal, "internal server error")
		}
		return
	}

	response.Success(c, resp)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrBadRequest, "invalid request: "+err.Error())
		return
	}

	resp, err := h.Svc.RefreshToken(req.RefreshToken)
	if err != nil {
		if appErr, ok := service.IsAppError(err); ok {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, errcode.ErrInternal, "internal server error")
		}
		return
	}

	response.Success(c, resp)
}

// GetProfile handles GET /api/auth/profile (requires auth middleware)
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Error(c, errcode.ErrUnauthorized, "user not authenticated")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		response.Error(c, errcode.ErrInternal, "invalid user ID type")
		return
	}

	user, err := h.Svc.GetProfile(userID)
	if err != nil {
		if appErr, ok := service.IsAppError(err); ok {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, errcode.ErrInternal, "internal server error")
		}
		return
	}

	response.Success(c, user)
}

// GetCaptcha handles GET /api/auth/captcha
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	id, b64, err := h.Svc.GenerateAndStoreCaptcha()
	if err != nil {
		response.Error(c, errcode.ErrInternal, "failed to generate captcha")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"captcha_id":     id,
			"captcha_image":  b64,
		},
	})
}
