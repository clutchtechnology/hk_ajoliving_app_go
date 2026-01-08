package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// User 用户模型
type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserType        string         `gorm:"size:20;not null;index" json:"user_type"`                        // individual=普通用户, agency=地产代理公司
	Email           string         `gorm:"size:255;uniqueIndex;not null" json:"email"`                     // 邮箱地址
	PasswordHash    string         `gorm:"size:255;not null" json:"-"`                                     // 密码哈希值
	Name            string         `gorm:"size:100;not null" json:"name"`                                  // 用户名称/公司名称
	Phone           string         `gorm:"size:20" json:"phone"`                                           // 联系电话
	Status          string         `gorm:"size:20;not null;default:'active';index" json:"status"`          // active=活跃, inactive=停用, suspended=暂停
	EmailVerified   bool           `gorm:"not null;default:false" json:"email_verified"`                   // 邮箱是否已验证
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`                                    // 邮箱验证时间
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`                                        // 最后登录时间
	CreatedAt       time.Time      `gorm:"index" json:"created_at"`                                        // 创建时间
	UpdatedAt       time.Time      `json:"updated_at"`                                                     // 更新时间
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`                                                 // 软删除时间
}

func (User) TableName() string {
	return "users"
}

// ============ Request DTO ============

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	UserType string `json:"user_type" binding:"required,oneof=individual agency"` // individual=普通用户, agency=地产代理公司
	Email    string `json:"email" binding:"required,email"`                       // 邮箱地址
	Password string `json:"password" binding:"required,min=6,max=50"`             // 密码
	Name     string `json:"name" binding:"required,min=2,max=100"`                // 用户名称/公司名称
	Phone    string `json:"phone" binding:"omitempty,max=20"`                     // 联系电话
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`       // 邮箱地址
	Password string `json:"password" binding:"required,min=6"`    // 密码
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=2,max=100"` // 用户名称/公司名称
	Phone *string `json:"phone" binding:"omitempty,max=20"`       // 联系电话
}

// ============ Response DTO ============

// UserResponse 用户响应
type UserResponse struct {
	ID              uint       `json:"id"`
	UserType        string     `json:"user_type"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Phone           string     `json:"phone,omitempty"`
	Status          string     `json:"status"`
	EmailVerified   bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// ToUserResponse 转换为用户响应
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		UserType:        u.UserType,
		Email:           u.Email,
		Name:            u.Name,
		Phone:           u.Phone,
		Status:          u.Status,
		EmailVerified:   u.EmailVerified,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
