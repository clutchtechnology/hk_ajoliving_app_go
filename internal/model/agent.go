package model

import (
	"time"

	"gorm.io/gorm"
)

// AgentStatus 代理状态常量
type AgentStatus string

const (
	AgentStatusActive    AgentStatus = "active"    // 活跃
	AgentStatusInactive  AgentStatus = "inactive"  // 停用
	AgentStatusSuspended AgentStatus = "suspended" // 暂停
)

// LicenseType 牌照类型常量
type LicenseType string

const (
	LicenseTypeIndividual  LicenseType = "individual"  // 个人牌照
	LicenseTypeSalesperson LicenseType = "salesperson" // 营业员牌照
)

// Agent 地产代理表模型
type Agent struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	AgentName         string         `gorm:"type:varchar(100);not null;index" json:"agent_name"`
	AgentNameEn       *string        `gorm:"type:varchar(100)" json:"agent_name_en,omitempty"`
	LicenseNo         string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"license_no"`
	LicenseType       LicenseType    `gorm:"type:varchar(20);not null;index" json:"license_type"`
	LicenseExpiryDate *time.Time     `gorm:"type:date" json:"license_expiry_date,omitempty"`
	AgencyID          *uint          `gorm:"index" json:"agency_id,omitempty"`
	Phone             string         `gorm:"type:varchar(20);not null" json:"phone"`
	Mobile            *string        `gorm:"type:varchar(20)" json:"mobile,omitempty"`
	Email             string         `gorm:"type:varchar(255);not null;index" json:"email"`
	WechatID          *string        `gorm:"type:varchar(50)" json:"wechat_id,omitempty"`
	Whatsapp          *string        `gorm:"type:varchar(20)" json:"whatsapp,omitempty"`
	OfficeAddress     *string        `gorm:"type:varchar(500)" json:"office_address,omitempty"`
	Specialization    *string        `gorm:"type:varchar(200)" json:"specialization,omitempty"`
	YearsExperience   *int           `json:"years_experience,omitempty"`
	ProfilePhoto      *string        `gorm:"type:varchar(500)" json:"profile_photo,omitempty"`
	Bio               *string        `gorm:"type:text" json:"bio,omitempty"`
	Rating            *float64       `gorm:"type:decimal(3,2);index" json:"rating,omitempty"`
	ReviewCount       int            `gorm:"not null;default:0" json:"review_count"`
	PropertiesSold    int            `gorm:"not null;default:0" json:"properties_sold"`
	PropertiesRented  int            `gorm:"not null;default:0" json:"properties_rented"`
	Status            AgentStatus    `gorm:"type:varchar(20);not null;index;default:'active'" json:"status"`
	IsVerified        bool           `gorm:"not null;default:false;index" json:"is_verified"`
	VerifiedAt        *time.Time     `json:"verified_at,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User         *User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Agency       *User                `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	ServiceAreas []AgentServiceArea   `gorm:"foreignKey:AgentID" json:"service_areas,omitempty"`
}

// TableName 指定表名
func (Agent) TableName() string {
	return "agents"
}

// IsActive 判断是否为活跃状态
func (a *Agent) IsActive() bool {
	return a.Status == AgentStatusActive
}

// IsVerified 判断牌照是否已验证
func (a *Agent) IsVerified() bool {
	return a.IsVerified
}

// HasAgency 判断是否加入代理公司
func (a *Agent) HasAgency() bool {
	return a.AgencyID != nil
}

// BeforeCreate GORM hook - 创建前执行
func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	if a.Status == "" {
		a.Status = AgentStatusActive
	}
	return nil
}

// AgentServiceArea 代理服务区域表模型
type AgentServiceArea struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AgentID    uint      `gorm:"not null;index;uniqueIndex:idx_agent_district" json:"agent_id"`
	DistrictID uint      `gorm:"not null;index;uniqueIndex:idx_agent_district" json:"district_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Agent    *Agent    `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	District *District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}

// TableName 指定表名
func (AgentServiceArea) TableName() string {
	return "agent_service_areas"
}
