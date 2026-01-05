package models

import (
	"time"

	"gorm.io/gorm"
)

// UserType 用户类型常量
type UserType string

const (
	UserTypeIndividual UserType = "individual" // 普通用户
	UserTypeAgency     UserType = "agency"     // 地产代理公司
)

// UserStatus 用户状态常量
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"    // 活跃
	UserStatusInactive  UserStatus = "inactive"  // 停用
	UserStatusSuspended UserStatus = "suspended" // 暂停
)

// User 用户表模型
type User struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Email           string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password        string         `gorm:"type:varchar(255);not null" json:"-"` // 不在JSON中显示
	FullName        string         `gorm:"type:varchar(100);column:full_name" json:"full_name"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone"`
	Avatar          string         `gorm:"type:varchar(500)" json:"avatar"`
	UserType        string         `gorm:"type:varchar(20);not null;index" json:"user_type"` // buyer, seller, agent, admin
	Status          string         `gorm:"type:varchar(20);not null;index;default:'active'" json:"status"` // active, inactive, suspended
	EmailVerified   bool           `gorm:"not null;default:false" json:"email_verified"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	AgencyDetail *AgencyDetail `gorm:"foreignKey:UserID" json:"agency_detail,omitempty"`
	Agent        *Agent        `gorm:"foreignKey:UserID" json:"agent,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 判断是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsEmailVerified 判断邮箱是否已验证
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}

// BeforeCreate GORM hook - 创建前执行
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 设置默认状态
	if u.Status == "" {
		u.Status = "active"
	}
	if u.UserType == "" {
		u.UserType = "buyer"
	}
	return nil
}

// ============ Request DTO ============

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"`
	UserType string `json:"user_type" binding:"required,oneof=buyer seller agent"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
}

// ============ Response DTO ============

// RegisterResponse 注册响应
type RegisterResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// AuthResponse 认证响应（登录/刷新Token）
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ForgotPasswordResponse 忘记密码响应
type ForgotPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// VerifyCodeRequest 验证码请求
type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// VerifyCodeResponse 验证码响应
type VerifyCodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
