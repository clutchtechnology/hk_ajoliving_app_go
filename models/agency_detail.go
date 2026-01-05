package models

import (
	"time"
)

// AgencyDetail 代理公司详情表模型
type AgencyDetail struct {
	ID                     uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                 uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	CompanyName            string         `gorm:"type:varchar(200);not null;index" json:"company_name"`
	CompanyNameEn          *string        `gorm:"type:varchar(200)" json:"company_name_en,omitempty"`
	LicenseNo              string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"license_no"`
	BusinessRegistrationNo *string        `gorm:"type:varchar(50)" json:"business_registration_no,omitempty"`
	Address                string         `gorm:"type:varchar(500);not null" json:"address"`
	Phone                  string         `gorm:"type:varchar(20);not null" json:"phone"`
	Fax                    *string        `gorm:"type:varchar(20)" json:"fax,omitempty"`
	Email                  string         `gorm:"type:varchar(255);not null" json:"email"`
	WebsiteURL             *string        `gorm:"type:varchar(500)" json:"website_url,omitempty"`
	EstablishedYear        *int           `json:"established_year,omitempty"`
	AgentCount             int            `gorm:"not null;default:0" json:"agent_count"`
	Description            *string        `gorm:"type:text" json:"description,omitempty"`
	LogoURL                *string        `gorm:"type:varchar(500)" json:"logo_url,omitempty"`
	CoverImageURL          *string        `gorm:"type:varchar(500)" json:"cover_image_url,omitempty"`
	Rating                 *float64       `gorm:"type:decimal(3,2);index" json:"rating,omitempty"`
	ReviewCount            int            `gorm:"not null;default:0" json:"review_count"`
	IsVerified             bool           `gorm:"not null;default:false;index" json:"is_verified"`
	VerifiedAt             *time.Time     `json:"verified_at,omitempty"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (AgencyDetail) TableName() string {
	return "agency_details"
}

// CheckIsVerified 判断是否已验证
func (ad *AgencyDetail) CheckIsVerified() bool {
	return ad.IsVerified
}

// HasWebsite 判断是否有官网
func (ad *AgencyDetail) HasWebsite() bool {
	return ad.WebsiteURL != nil && *ad.WebsiteURL != ""
}

// ============ Request DTO ============

// ListAgenciesRequest 获取代理公司列表请求
type ListAgenciesRequest struct {
	Keyword     *string  `form:"keyword"`
	CompanyName *string  `form:"company_name"`
	IsVerified  *bool    `form:"is_verified"`
	MinRating   *float64 `form:"min_rating" binding:"omitempty,gte=0,lte=5"`
	Page        int      `form:"page,default=1" binding:"min=1"`
	PageSize    int      `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy      string   `form:"sort_by,default=rating"`
	SortOrder   string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// SearchAgenciesRequest 搜索代理公司请求
type SearchAgenciesRequest struct {
	Keyword    string   `form:"keyword" binding:"required"`
	IsVerified *bool    `form:"is_verified"`
	MinRating  *float64 `form:"min_rating"`
	Page       int      `form:"page,default=1" binding:"min=1"`
	PageSize   int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ============ Additional Request DTO ============

// ContactAgencyRequest 联系代理公司请求
type ContactAgencyRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}
