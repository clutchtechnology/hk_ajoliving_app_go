package response

import "time"

// AuthResponse 认证响应
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` // 秒
	User         *UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	UserType  string    `json:"user_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Message string    `json:"message"`
	User    *UserInfo `json:"user"`
}

// ForgotPasswordResponse 忘记密码响应
type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// VerifyCodeResponse 验证码验证响应
type VerifyCodeResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // 用于后续操作的临时令牌
}

// LogoutResponse 登出响应
type LogoutResponse struct {
	Message string `json:"message"`
}
