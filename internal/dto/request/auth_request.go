package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6,max=100"`
	Phone       string `json:"phone" binding:"omitempty,min=8,max=20"`
	FullName    string `json:"full_name" binding:"omitempty,max=100"`
	UserType    string `json:"user_type" binding:"omitempty,oneof=buyer seller agent admin"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // 可以是 email 或 username
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

// VerifyCodeRequest 验证码验证请求
type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
	Type  string `json:"type" binding:"required,oneof=register reset_password"`
}
