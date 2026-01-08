package models

import (
	"time"
)

// ============ GORM Model ============

// District 地区模型
type District struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	NameZhHant string    `gorm:"size:100;not null" json:"name_zh_hant"`           // 中文繁体名称
	NameZhHans string    `gorm:"size:100" json:"name_zh_hans,omitempty"`         // 中文简体名称
	NameEn     string    `gorm:"size:100" json:"name_en,omitempty"`              // 英文名称
	Region     string    `gorm:"size:50;not null;index" json:"region"`           // HK_ISLAND=港岛, KOWLOON=九龙, NEW_TERRITORIES=新界
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`           // 排序顺序
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (District) TableName() string {
	return "districts"
}

// ============ Request DTO ============

// ListDistrictsRequest 获取地区列表请求
type ListDistrictsRequest struct {
	Region string `form:"region" binding:"omitempty,oneof=HK_ISLAND KOWLOON NEW_TERRITORIES"`
}

// GetDistrictPropertiesRequest 获取地区房源请求
type GetDistrictPropertiesRequest struct {
	ListingType  *string `form:"listing_type" binding:"omitempty,oneof=sale rent"`
	PropertyType *string `form:"property_type"`
	MinPrice     *float64 `form:"min_price"`
	MaxPrice     *float64 `form:"max_price"`
	Bedrooms     *int     `form:"bedrooms"`
	Page         int      `form:"page,default=1" binding:"min=1"`
	PageSize     int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// GetDistrictEstatesRequest 获取地区屋苑请求
type GetDistrictEstatesRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ============ Response DTO ============

// DistrictResponse 地区响应
type DistrictResponse struct {
	ID         uint   `json:"id"`
	NameZhHant string `json:"name_zh_hant"`
	NameZhHans string `json:"name_zh_hans,omitempty"`
	NameEn     string `json:"name_en,omitempty"`
	Region     string `json:"region"`
	SortOrder  int    `json:"sort_order"`
}

// ToDistrictResponse 转换为地区响应
func (d *District) ToDistrictResponse() *DistrictResponse {
	return &DistrictResponse{
		ID:         d.ID,
		NameZhHant: d.NameZhHant,
		NameZhHans: d.NameZhHans,
		NameEn:     d.NameEn,
		Region:     d.Region,
		SortOrder:  d.SortOrder,
	}
}

// DistrictDetailResponse 地区详情响应
type DistrictDetailResponse struct {
	ID                uint   `json:"id"`
	NameZhHant        string `json:"name_zh_hant"`
	NameZhHans        string `json:"name_zh_hans,omitempty"`
	NameEn            string `json:"name_en,omitempty"`
	Region            string `json:"region"`
	PropertyCount     int64  `json:"property_count"`
	EstateCount       int64  `json:"estate_count"`
	AvgPropertyPrice  float64 `json:"avg_property_price"`
	SortOrder         int    `json:"sort_order"`
}

// DistrictStatisticsResponse 地区统计数据响应
type DistrictStatisticsResponse struct {
	DistrictID            uint    `json:"district_id"`
	DistrictName          string  `json:"district_name"`
	TotalProperties       int64   `json:"total_properties"`
	PropertiesForSale     int64   `json:"properties_for_sale"`
	PropertiesForRent     int64   `json:"properties_for_rent"`
	TotalEstates          int64   `json:"total_estates"`
	AvgSalePrice          float64 `json:"avg_sale_price"`
	AvgRentPrice          float64 `json:"avg_rent_price"`
	AvgPricePerSqft       float64 `json:"avg_price_per_sqft"`
	NewPropertiesThisWeek int64   `json:"new_properties_this_week"`
	NewPropertiesThisMonth int64  `json:"new_properties_this_month"`
}
