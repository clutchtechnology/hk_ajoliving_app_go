package models

import (
	"time"
)

// ============ GORM Model ============

// SearchHistory 搜索历史记录
type SearchHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"index" json:"user_id"`                    // 用户ID（可选，未登录用户为空）
	Keyword    string    `gorm:"size:255;not null;index" json:"keyword"`  // 搜索关键词
	SearchType string    `gorm:"size:50;index" json:"search_type"`        // 搜索类型：global, property, estate, agent
	ResultCount int      `gorm:"default:0" json:"result_count"`           // 搜索结果数量
	IPAddress  string    `gorm:"size:45" json:"ip_address"`               // IP地址
	UserAgent  string    `gorm:"size:500" json:"user_agent"`              // 浏览器信息
	CreatedAt  time.Time `json:"created_at"`
}

func (SearchHistory) TableName() string {
	return "search_histories"
}

// ============ Request DTO ============

// GlobalSearchRequest 全局搜索请求
type GlobalSearchRequest struct {
	Keyword  string `form:"keyword" binding:"required,min=1"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// SearchPropertiesRequest 搜索房产请求
type SearchPropertiesRequest struct {
	Keyword      string   `form:"keyword" binding:"required,min=1"`
	ListingType  *string  `form:"listing_type" binding:"omitempty,oneof=sale rent"`
	DistrictID   *uint    `form:"district_id"`
	MinPrice     *float64 `form:"min_price"`
	MaxPrice     *float64 `form:"max_price"`
	Bedrooms     *int     `form:"bedrooms"`
	PropertyType *string  `form:"property_type"`
	Page         int      `form:"page,default=1" binding:"min=1"`
	PageSize     int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// SearchEstatesRequest 搜索屋苑请求
type SearchEstatesRequest struct {
	Keyword    string `form:"keyword" binding:"required,min=1"`
	DistrictID *uint  `form:"district_id"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// SearchAgentsRequest 搜索代理人请求
type SearchAgentsRequest struct {
	Keyword        string `form:"keyword" binding:"required,min=1"`
	DistrictID     *uint  `form:"district_id"`
	Specialization *string `form:"specialization"`
	Page           int    `form:"page,default=1" binding:"min=1"`
	PageSize       int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// GetSearchSuggestionsRequest 搜索建议请求
type GetSearchSuggestionsRequest struct {
	Keyword string `form:"keyword" binding:"required,min=1,max=50"`
	Limit   int    `form:"limit,default=10" binding:"min=1,max=20"`
}

// GetSearchHistoryRequest 搜索历史请求
type GetSearchHistoryRequest struct {
	SearchType string `form:"search_type" binding:"omitempty,oneof=global property estate agent"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ============ Response DTO ============

// GlobalSearchResponse 全局搜索响应
type GlobalSearchResponse struct {
	Properties      []PropertySearchResult      `json:"properties"`
	Estates         []EstateSearchResult        `json:"estates"`
	Agents          []AgentSearchResult         `json:"agents"`
	Agencies        []AgencySearchResult        `json:"agencies"`
	TotalResults    int                         `json:"total_results"`
	PropertyCount   int                         `json:"property_count"`
	EstateCount     int                         `json:"estate_count"`
	AgentCount      int                         `json:"agent_count"`
	AgencyCount     int                         `json:"agency_count"`
}

// PropertySearchResult 房产搜索结果
type PropertySearchResult struct {
	ID           uint    `json:"id"`
	PropertyNo   string  `json:"property_no"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Area         float64 `json:"area"`
	Bedrooms     int     `json:"bedrooms"`
	PropertyType string  `json:"property_type"`
	ListingType  string  `json:"listing_type"`
	Address      string  `json:"address"`
	DistrictName string  `json:"district_name"`
	Status       string  `json:"status"`
	CoverImage   string  `json:"cover_image,omitempty"`
}

// EstateSearchResult 屋苑搜索结果
type EstateSearchResult struct {
	ID                  uint    `json:"id"`
	Name                string  `json:"name"`
	NameEn              string  `json:"name_en,omitempty"`
	Address             string  `json:"address"`
	DistrictName        string  `json:"district_name"`
	TotalUnits          int     `json:"total_units"`
	AvgTransactionPrice float64 `json:"avg_transaction_price"`
	CompletionYear      int     `json:"completion_year"`
}

// AgentSearchResult 代理人搜索结果
type AgentSearchResult struct {
	ID              uint    `json:"id"`
	AgentName       string  `json:"agent_name"`
	AgentNameEn     string  `json:"agent_name_en,omitempty"`
	LicenseNo       string  `json:"license_no"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	Specialization  string  `json:"specialization,omitempty"`
	Rating          float64 `json:"rating"`
	PropertiesSold  int     `json:"properties_sold"`
	ProfilePhoto    string  `json:"profile_photo,omitempty"`
}

// AgencySearchResult 代理公司搜索结果
type AgencySearchResult struct {
	ID              uint    `json:"id"`
	CompanyName     string  `json:"company_name"`
	CompanyNameEn   string  `json:"company_name_en,omitempty"`
	LicenseNo       string  `json:"license_no"`
	Phone           string  `json:"phone"`
	Address         string  `json:"address"`
	AgentCount      int     `json:"agent_count"`
	Rating          float64 `json:"rating"`
	IsVerified      bool    `json:"is_verified"`
	LogoURL         string  `json:"logo_url,omitempty"`
}

// SearchSuggestion 搜索建议
type SearchSuggestion struct {
	Keyword    string `json:"keyword"`
	Type       string `json:"type"`        // property, estate, agent, district
	Count      int    `json:"count"`       // 搜索次数或相关结果数
	Label      string `json:"label"`       // 显示标签
}

// SearchHistoryResponse 搜索历史响应
type SearchHistoryResponse struct {
	ID          uint      `json:"id"`
	Keyword     string    `json:"keyword"`
	SearchType  string    `json:"search_type"`
	ResultCount int       `json:"result_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// PaginatedSearchHistoryResponse 分页搜索历史响应
type PaginatedSearchHistoryResponse struct {
	Histories  []SearchHistoryResponse `json:"histories"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}
