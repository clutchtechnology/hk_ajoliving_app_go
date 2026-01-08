package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// Agent 地产代理模型
type Agent struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"uniqueIndex;not null" json:"user_id"`                   // 关联用户ID
	AgentName         string         `gorm:"size:100;not null;index" json:"agent_name"`             // 代理人姓名
	AgentNameEn       string         `gorm:"size:100" json:"agent_name_en,omitempty"`               // 英文姓名
	LicenseNo         string         `gorm:"size:50;uniqueIndex;not null" json:"license_no"`        // 地产代理牌照号码
	LicenseType       string         `gorm:"size:20;not null;index" json:"license_type"`            // individual=个人牌照, salesperson=营业员牌照
	LicenseExpiryDate *time.Time     `json:"license_expiry_date,omitempty"`                         // 牌照到期日期
	AgencyID          *uint          `gorm:"index" json:"agency_id,omitempty"`                      // 所属代理公司ID
	Phone             string         `gorm:"size:20;not null" json:"phone"`                         // 联系电话
	Mobile            string         `gorm:"size:20" json:"mobile,omitempty"`                       // 手机号码
	Email             string         `gorm:"size:255;not null;index" json:"email"`                  // 电子邮箱
	WechatID          string         `gorm:"size:50" json:"wechat_id,omitempty"`                    // 微信号
	Whatsapp          string         `gorm:"size:20" json:"whatsapp,omitempty"`                     // WhatsApp号码
	OfficeAddress     string         `gorm:"size:500" json:"office_address,omitempty"`              // 办公地址
	Specialization    string         `gorm:"size:200" json:"specialization,omitempty"`              // 专长领域
	YearsExperience   int            `gorm:"default:0" json:"years_experience"`                     // 从业年限
	ProfilePhoto      string         `gorm:"size:500" json:"profile_photo,omitempty"`               // 个人照片URL
	Bio               string         `gorm:"type:text" json:"bio,omitempty"`                        // 个人简介
	Rating            float64        `gorm:"type:decimal(3,2);index" json:"rating"`                 // 评分（0-5）
	ReviewCount       int            `gorm:"default:0" json:"review_count"`                         // 评价数量
	PropertiesSold    int            `gorm:"default:0" json:"properties_sold"`                      // 已售物业数量
	PropertiesRented  int            `gorm:"default:0" json:"properties_rented"`                    // 已租物业数量
	Status            string         `gorm:"size:20;not null;default:'active';index" json:"status"` // active=活跃, inactive=停用, suspended=暂停
	IsVerified        bool           `gorm:"default:false;index" json:"is_verified"`                // 是否已验证牌照
	VerifiedAt        *time.Time     `json:"verified_at,omitempty"`                                 // 验证时间
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Agent) TableName() string {
	return "agents"
}

// AgentServiceArea 代理服务区域
type AgentServiceArea struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AgentID    uint      `gorm:"not null;index" json:"agent_id"`
	DistrictID uint      `gorm:"not null;index" json:"district_id"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	Agent    *Agent    `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	District *District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}

func (AgentServiceArea) TableName() string {
	return "agent_service_areas"
}

// AgentContact 联系代理人记录
type AgentContact struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"not null;index" json:"agent_id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"`      // 联系人用户ID（可选）
	Name      string    `gorm:"size:100;not null" json:"name"`       // 联系人姓名
	Phone     string    `gorm:"size:20;not null" json:"phone"`       // 联系电话
	Email     string    `gorm:"size:255" json:"email,omitempty"`     // 邮箱
	Message   string    `gorm:"type:text;not null" json:"message"`   // 留言内容
	CreatedAt time.Time `json:"created_at"`
}

func (AgentContact) TableName() string {
	return "agent_contacts"
}

// ============ Request DTO ============

// ListAgentsRequest 代理人列表请求
type ListAgentsRequest struct {
	LicenseType    *string `form:"license_type" binding:"omitempty,oneof=individual salesperson"`
	AgencyID       *uint   `form:"agency_id"`
	DistrictID     *uint   `form:"district_id"`
	Status         *string `form:"status" binding:"omitempty,oneof=active inactive suspended"`
	IsVerified     *bool   `form:"is_verified"`
	Specialization *string `form:"specialization"`
	Keyword        string  `form:"keyword"`
	Page           int     `form:"page,default=1" binding:"min=1"`
	PageSize       int     `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ContactAgentRequest 联系代理人请求
type ContactAgentRequest struct {
	Name    string `json:"name" binding:"required,max=100"`
	Phone   string `json:"phone" binding:"required,max=20"`
	Email   string `json:"email" binding:"omitempty,email,max=255"`
	Message string `json:"message" binding:"required,max=1000"`
}

// ============ Response DTO ============

// AgentResponse 代理人响应
type AgentResponse struct {
	ID               uint    `json:"id"`
	AgentName        string  `json:"agent_name"`
	AgentNameEn      string  `json:"agent_name_en,omitempty"`
	LicenseNo        string  `json:"license_no"`
	LicenseType      string  `json:"license_type"`
	Phone            string  `json:"phone"`
	Email            string  `json:"email"`
	ProfilePhoto     string  `json:"profile_photo,omitempty"`
	Specialization   string  `json:"specialization,omitempty"`
	YearsExperience  int     `json:"years_experience"`
	Rating           float64 `json:"rating"`
	ReviewCount      int     `json:"review_count"`
	PropertiesSold   int     `json:"properties_sold"`
	PropertiesRented int     `json:"properties_rented"`
	Status           string  `json:"status"`
	IsVerified       bool    `json:"is_verified"`
}

// AgentDetailResponse 代理人详情响应
type AgentDetailResponse struct {
	ID                uint              `json:"id"`
	AgentName         string            `json:"agent_name"`
	AgentNameEn       string            `json:"agent_name_en,omitempty"`
	LicenseNo         string            `json:"license_no"`
	LicenseType       string            `json:"license_type"`
	LicenseExpiryDate *time.Time        `json:"license_expiry_date,omitempty"`
	AgencyID          *uint             `json:"agency_id,omitempty"`
	Phone             string            `json:"phone"`
	Mobile            string            `json:"mobile,omitempty"`
	Email             string            `json:"email"`
	WechatID          string            `json:"wechat_id,omitempty"`
	Whatsapp          string            `json:"whatsapp,omitempty"`
	OfficeAddress     string            `json:"office_address,omitempty"`
	Specialization    string            `json:"specialization,omitempty"`
	YearsExperience   int               `json:"years_experience"`
	ProfilePhoto      string            `json:"profile_photo,omitempty"`
	Bio               string            `json:"bio,omitempty"`
	Rating            float64           `json:"rating"`
	ReviewCount       int               `json:"review_count"`
	PropertiesSold    int               `json:"properties_sold"`
	PropertiesRented  int               `json:"properties_rented"`
	Status            string            `json:"status"`
	IsVerified        bool              `json:"is_verified"`
	VerifiedAt        *time.Time        `json:"verified_at,omitempty"`
	ServiceAreas      []*DistrictBrief  `json:"service_areas,omitempty"`
}

// DistrictBrief 地区简要信息
type DistrictBrief struct {
	ID         uint   `json:"id"`
	NameZhHant string `json:"name_zh_hant"`
	NameEn     string `json:"name_en,omitempty"`
	Region     string `json:"region"`
}

// PaginatedAgentsResponse 分页代理人响应
type PaginatedAgentsResponse struct {
	Items      []*AgentResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
