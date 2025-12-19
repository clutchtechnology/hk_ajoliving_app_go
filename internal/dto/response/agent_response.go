package response

import "time"

// AgentResponse 代理人响应
type AgentResponse struct {
	ID                uint                      `json:"id"`
	UserID            uint                      `json:"user_id"`
	AgentName         string                    `json:"agent_name"`
	AgentNameEn       *string                   `json:"agent_name_en,omitempty"`
	LicenseNo         string                    `json:"license_no"`
	LicenseType       string                    `json:"license_type"`
	LicenseExpiryDate *time.Time                `json:"license_expiry_date,omitempty"`
	AgencyID          *uint                     `json:"agency_id,omitempty"`
	AgencyName        *string                   `json:"agency_name,omitempty"`       // 代理公司名称
	Phone             string                    `json:"phone"`
	Mobile            *string                   `json:"mobile,omitempty"`
	Email             string                    `json:"email"`
	WechatID          *string                   `json:"wechat_id,omitempty"`
	Whatsapp          *string                   `json:"whatsapp,omitempty"`
	OfficeAddress     *string                   `json:"office_address,omitempty"`
	Specialization    *string                   `json:"specialization,omitempty"`
	YearsExperience   *int                      `json:"years_experience,omitempty"`
	ProfilePhoto      *string                   `json:"profile_photo,omitempty"`
	Bio               *string                   `json:"bio,omitempty"`
	Rating            *float64                  `json:"rating,omitempty"`
	ReviewCount       int                       `json:"review_count"`
	PropertiesSold    int                       `json:"properties_sold"`
	PropertiesRented  int                       `json:"properties_rented"`
	Status            string                    `json:"status"`
	IsVerified        bool                      `json:"is_verified"`
	VerifiedAt        *time.Time                `json:"verified_at,omitempty"`
	ServiceAreas      []AgentServiceAreaResponse `json:"service_areas,omitempty"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

// AgentListItemResponse 代理人列表项响应
type AgentListItemResponse struct {
	ID               uint       `json:"id"`
	AgentName        string     `json:"agent_name"`
	AgentNameEn      *string    `json:"agent_name_en,omitempty"`
	LicenseNo        string     `json:"license_no"`
	LicenseType      string     `json:"license_type"`
	AgencyID         *uint      `json:"agency_id,omitempty"`
	AgencyName       *string    `json:"agency_name,omitempty"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	ProfilePhoto     *string    `json:"profile_photo,omitempty"`
	Specialization   *string    `json:"specialization,omitempty"`
	YearsExperience  *int       `json:"years_experience,omitempty"`
	Rating           *float64   `json:"rating,omitempty"`
	ReviewCount      int        `json:"review_count"`
	PropertiesSold   int        `json:"properties_sold"`
	PropertiesRented int        `json:"properties_rented"`
	Status           string     `json:"status"`
	IsVerified       bool       `json:"is_verified"`
	CreatedAt        time.Time  `json:"created_at"`
}

// AgentServiceAreaResponse 代理服务区域响应
type AgentServiceAreaResponse struct {
	ID           uint   `json:"id"`
	DistrictID   uint   `json:"district_id"`
	DistrictName string `json:"district_name"`
}

// AgentContactResponse 联系请求响应
type AgentContactResponse struct {
	ID          uint       `json:"id"`
	AgentID     uint       `json:"agent_id"`
	AgentName   string     `json:"agent_name"`
	PropertyID  *uint      `json:"property_id,omitempty"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Message     string     `json:"message"`
	ContactType string     `json:"contact_type"`
	Status      string     `json:"status"`
	ContactedAt *time.Time `json:"contacted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
