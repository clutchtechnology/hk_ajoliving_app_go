package response

import "time"

// AgencyListItemResponse 代理公司列表项响应
type AgencyListItemResponse struct {
	ID                uint     `json:"id"`                               // 代理公司ID
	CompanyName       string   `json:"company_name"`                     // 公司名称（中文）
	CompanyNameEn     *string  `json:"company_name_en,omitempty"`        // 公司名称（英文）
	LogoURL           *string  `json:"logo_url,omitempty"`               // 公司Logo URL
	Address           string   `json:"address"`                          // 地址
	Phone             string   `json:"phone"`                            // 电话
	Email             string   `json:"email"`                            // 电子邮箱
	AgentCount        int      `json:"agent_count"`                      // 代理人数量
	Rating            *float64 `json:"rating,omitempty"`                 // 评分
	ReviewCount       int      `json:"review_count"`                     // 评价数量
	IsVerified        bool     `json:"is_verified"`                      // 是否已验证
	EstablishedYear   *int     `json:"established_year,omitempty"`       // 成立年份
	PropertyCount     int      `json:"property_count"`                   // 房源数量
}

// AgencyResponse 代理公司详情响应
type AgencyResponse struct {
	ID                     uint       `json:"id"`                               // 代理公司ID
	CompanyName            string     `json:"company_name"`                     // 公司名称（中文）
	CompanyNameEn          *string    `json:"company_name_en,omitempty"`        // 公司名称（英文）
	LicenseNo              string     `json:"license_no"`                       // 牌照号码
	BusinessRegistrationNo *string    `json:"business_registration_no,omitempty"` // 商业登记号
	Address                string     `json:"address"`                          // 地址
	Phone                  string     `json:"phone"`                            // 电话
	Fax                    *string    `json:"fax,omitempty"`                    // 传真
	Email                  string     `json:"email"`                            // 电子邮箱
	WebsiteURL             *string    `json:"website_url,omitempty"`            // 官网URL
	EstablishedYear        *int       `json:"established_year,omitempty"`       // 成立年份
	AgentCount             int        `json:"agent_count"`                      // 代理人数量
	Description            *string    `json:"description,omitempty"`            // 公司简介
	LogoURL                *string    `json:"logo_url,omitempty"`               // 公司Logo URL
	CoverImageURL          *string    `json:"cover_image_url,omitempty"`        // 封面图URL
	Rating                 *float64   `json:"rating,omitempty"`                 // 评分
	ReviewCount            int        `json:"review_count"`                     // 评价数量
	IsVerified             bool       `json:"is_verified"`                      // 是否已验证
	VerifiedAt             *time.Time `json:"verified_at,omitempty"`            // 验证时间
	PropertyCount          int        `json:"property_count"`                   // 房源数量
	TopAgents              []AgentBasicInfo `json:"top_agents,omitempty"`       // 优秀代理人
	ServiceDistricts       []DistrictInfo `json:"service_districts,omitempty"`  // 服务地区
	CreatedAt              time.Time  `json:"created_at"`                       // 创建时间
	UpdatedAt              time.Time  `json:"updated_at"`                       // 更新时间
}

// AgentBasicInfo 代理人基本信息（用于代理公司详情）
type AgentBasicInfo struct {
	ID            uint     `json:"id"`                          // 代理人ID
	AgentName     string   `json:"agent_name"`                  // 代理人姓名
	ProfilePhoto  *string  `json:"profile_photo,omitempty"`     // 头像
	Phone         string   `json:"phone"`                       // 电话
	Email         string   `json:"email"`                       // 邮箱
	Rating        *float64 `json:"rating,omitempty"`            // 评分
	ReviewCount   int      `json:"review_count"`                // 评价数
}

// DistrictInfo 地区信息
type DistrictInfo struct {
	ID         uint   `json:"id"`          // 地区ID
	NameZh     string `json:"name_zh"`     // 地区名称（中文）
	NameEn     string `json:"name_en"`     // 地区名称（英文）
}

// ContactAgencyResponse 联系代理公司响应
type ContactAgencyResponse struct {
	Success bool   `json:"success"`         // 是否成功
	Message string `json:"message"`         // 提示信息
}
