package request

// RegisterReq is the request body for user registration.
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginReq is the request body for user login.
type LoginReq struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaAnswer string `json:"captcha_answer" binding:"required"`
}

// RefreshReq is the request body for token refresh.
type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
