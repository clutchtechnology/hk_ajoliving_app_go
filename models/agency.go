package models

import (
	"time"
)

// ============ GORM Model ============

// AgencyDetail 代理公司详情
type AgencyDetail struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	UserID                 uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	CompanyName            string         `gorm:"size:200;not null" json:"company_name"`
	CompanyNameEn          string         `gorm:"size:200" json:"company_name_en"`
	LicenseNo              string         `gorm:"size:50;uniqueIndex;not null" json:"license_no"`
	BusinessRegistrationNo string         `gorm:"size:50" json:"business_registration_no"`
	Address                string         `gorm:"size:500;not null" json:"address"`
	Phone                  string         `gorm:"size:20;not null" json:"phone"`
	Fax                    string         `gorm:"size:20" json:"fax"`
	Email                  string         `gorm:"size:255;not null" json:"email"`
	WebsiteURL             string         `gorm:"size:500" json:"website_url"`
	EstablishedYear        int            `json:"established_year"`
	AgentCount             int            `gorm:"default:0" json:"agent_count"`
	Description            string         `gorm:"type:text" json:"description"`
	LogoURL                string         `gorm:"size:500" json:"logo_url"`
	CoverImageURL          string         `gorm:"size:500" json:"cover_image_url"`
	Rating                 float64        `gorm:"type:decimal(3,2)" json:"rating"`
	ReviewCount            int            `gorm:"default:0" json:"review_count"`
	IsVerified             bool           `gorm:"default:false" json:"is_verified"`
	VerifiedAt             *time.Time     `json:"verified_at"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AgencyDetail) TableName() string {
	return "agency_details"
}

// AgencyContact 代理公司联系记录
type AgencyContact struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AgencyID   uint      `gorm:"index;not null" json:"agency_id"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Phone      string    `gorm:"size:20;not null" json:"phone"`
	Email      string    `gorm:"size:255" json:"email"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	PropertyID *uint     `gorm:"index" json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	Agency   *User     `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Property *Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (AgencyContact) TableName() string {
	return "agency_contacts"
}

// ============ Request DTO ============

// ListAgenciesRequest 获取代理公司列表请求
type ListAgenciesRequest struct {
	DistrictID  *uint   `form:"district_id"`
	MinRating   *float64 `form:"min_rating"`
	IsVerified  *bool   `form:"is_verified"`
	Keyword     string  `form:"keyword"`
	SortBy      string  `form:"sort_by" binding:"omitempty,oneof=rating agent_count established_year created_at"`
	SortOrder   string  `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page        int     `form:"page,default=1" binding:"min=1"`
	PageSize    int     `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ContactAgencyRequest 联系代理公司请求
type ContactAgencyRequest struct {
	Name       string `json:"name" binding:"required,max=100"`
	Phone      string `json:"phone" binding:"required,max=20"`
	Email      string `json:"email" binding:"omitempty,email"`
	Message    string `json:"message" binding:"required"`
	PropertyID *uint  `json:"property_id"`
}

// SearchAgenciesRequest 搜索代理公司请求
type SearchAgenciesRequest struct {
	Keyword   string `form:"keyword" binding:"required,min=1"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ============ Response DTO ============

// AgencyResponse 代理公司响应（列表用）
type AgencyResponse struct {
	ID              uint    `json:"id"`
	CompanyName     string  `json:"company_name"`
	CompanyNameEn   string  `json:"company_name_en,omitempty"`
	LogoURL         string  `json:"logo_url,omitempty"`
	Address         string  `json:"address"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	WebsiteURL      string  `json:"website_url,omitempty"`
	EstablishedYear int     `json:"established_year,omitempty"`
	AgentCount      int     `json:"agent_count"`
	Rating          float64 `json:"rating"`
	ReviewCount     int     `json:"review_count"`
	IsVerified      bool    `json:"is_verified"`
}

// AgencyDetailResponse 代理公司详情响应
type AgencyDetailResponse struct {
	ID                     uint       `json:"id"`
	CompanyName            string     `json:"company_name"`
	CompanyNameEn          string     `json:"company_name_en,omitempty"`
	LicenseNo              string     `json:"license_no"`
	BusinessRegistrationNo string     `json:"business_registration_no,omitempty"`
	Address                string     `json:"address"`
	Phone                  string     `json:"phone"`
	Fax                    string     `json:"fax,omitempty"`
	Email                  string     `json:"email"`
	WebsiteURL             string     `json:"website_url,omitempty"`
	EstablishedYear        int        `json:"established_year,omitempty"`
	AgentCount             int        `json:"agent_count"`
	Description            string     `json:"description,omitempty"`
	LogoURL                string     `json:"logo_url,omitempty"`
	CoverImageURL          string     `json:"cover_image_url,omitempty"`
	Rating                 float64    `json:"rating"`
	ReviewCount            int        `json:"review_count"`
	IsVerified             bool       `json:"is_verified"`
	VerifiedAt             *time.Time `json:"verified_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// PaginatedAgenciesResponse 分页代理公司响应
type PaginatedAgenciesResponse struct {
	Agencies   []AgencyResponse `json:"agencies"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// ContactAgencyResponse 联系代理公司响应
type ContactAgencyResponse struct {
	ID        uint      `json:"id"`
	AgencyID  uint      `json:"agency_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
