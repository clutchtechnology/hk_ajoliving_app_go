package model

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
