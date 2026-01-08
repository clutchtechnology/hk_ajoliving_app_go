package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// Estate 屋苑模型
type Estate struct {
	ID                           uint           `gorm:"primaryKey" json:"id"`
	Name                         string         `gorm:"size:200;not null;index" json:"name"`                              // 屋苑名称
	NameEn                       string         `gorm:"size:200" json:"name_en,omitempty"`                                // 英文名称
	Address                      string         `gorm:"size:500;not null" json:"address"`                                 // 详细地址
	DistrictID                   uint           `gorm:"not null;index" json:"district_id"`                                // 所属地区ID
	TotalBlocks                  int            `json:"total_blocks,omitempty"`                                           // 总座数
	TotalUnits                   int            `json:"total_units,omitempty"`                                            // 总单位数
	CompletionYear               int            `json:"completion_year,omitempty"`                                        // 落成年份
	Developer                    string         `gorm:"size:200" json:"developer,omitempty"`                              // 发展商
	ManagementCompany            string         `gorm:"size:200" json:"management_company,omitempty"`                     // 管理公司
	PrimarySchoolNet             string         `gorm:"size:50;index" json:"primary_school_net,omitempty"`                // 小学校网
	SecondarySchoolNet           string         `gorm:"size:50;index" json:"secondary_school_net,omitempty"`              // 中学校网
	RecentTransactionsCount      int            `gorm:"default:0" json:"recent_transactions_count"`                       // 近期成交数量
	ForSaleCount                 int            `gorm:"default:0" json:"for_sale_count"`                                  // 当前放盘数量
	ForRentCount                 int            `gorm:"default:0" json:"for_rent_count"`                                  // 当前租盘数量
	AvgTransactionPrice          float64        `gorm:"index" json:"avg_transaction_price,omitempty"`                     // 平均成交价（港币/平方尺）
	AvgTransactionPriceUpdatedAt *time.Time     `json:"avg_transaction_price_updated_at,omitempty"`                       // 平均成交价更新时间
	Description                  string         `gorm:"type:text" json:"description,omitempty"`                           // 屋苑描述
	ViewCount                    int            `gorm:"default:0" json:"view_count"`                                      // 浏览次数
	FavoriteCount                int            `gorm:"default:0" json:"favorite_count"`                                  // 收藏次数
	IsFeatured                   bool           `gorm:"default:false;index" json:"is_featured"`                           // 是否精选屋苑
	CreatedAt                    time.Time      `gorm:"index" json:"created_at"`                                          // 创建时间
	UpdatedAt                    time.Time      `json:"updated_at"`                                                       // 更新时间
	DeletedAt                    gorm.DeletedAt `gorm:"index" json:"-"`                                                   // 软删除时间

	// 关联
	District   *District       `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Images     []EstateImage   `gorm:"foreignKey:EstateID" json:"images,omitempty"`
	Facilities []Facility      `gorm:"many2many:estate_facilities;" json:"facilities,omitempty"`
	Properties []Property      `gorm:"-" json:"-"` // 逻辑关联，不创建数据库外键约束
}

func (Estate) TableName() string {
	return "estates"
}

// EstateImage 屋苑图片
type EstateImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EstateID  uint      `gorm:"not null;index" json:"estate_id"`
	ImageURL  string    `gorm:"size:500;not null" json:"image_url"`
	ImageType string    `gorm:"size:20;not null" json:"image_type"` // exterior=外观, facilities=设施, environment=环境, aerial=航拍
	Title     string    `gorm:"size:200" json:"title,omitempty"`
	SortOrder int       `gorm:"not null" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (EstateImage) TableName() string {
	return "estate_images"
}

// EstateFacility 屋苑设施关联表
type EstateFacility struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EstateID   uint      `gorm:"not null;index;uniqueIndex:idx_estate_facility" json:"estate_id"`
	FacilityID uint      `gorm:"not null;index;uniqueIndex:idx_estate_facility" json:"facility_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (EstateFacility) TableName() string {
	return "estate_facilities"
}

// Facility 设施字典
type Facility struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	NameZhHant  string    `gorm:"size:100;not null" json:"name_zh_hant"`    // 中文繁体名称
	NameZhHans  string    `gorm:"size:100" json:"name_zh_hans,omitempty"`   // 中文简体名称
	NameEn      string    `gorm:"size:100" json:"name_en,omitempty"`        // 英文名称
	Icon        string    `gorm:"size:100" json:"icon,omitempty"`           // 图标标识
	Category    string    `gorm:"size:50;not null;index" json:"category"`   // building=大厦设施, unit=单位设施
	SortOrder   int       `gorm:"not null" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Facility) TableName() string {
	return "facilities"
}

// ============ Request DTO ============

// ListEstatesRequest 获取屋苑列表请求
type ListEstatesRequest struct {
	DistrictID         *uint   `form:"district_id"`                                                       // 地区ID
	PrimarySchoolNet   *string `form:"primary_school_net"`                                                // 小学校网
	SecondarySchoolNet *string `form:"secondary_school_net"`                                              // 中学校网
	Developer          *string `form:"developer"`                                                         // 发展商
	MinAvgPrice        *float64 `form:"min_avg_price"`                                                    // 最低平均价
	MaxAvgPrice        *float64 `form:"max_avg_price"`                                                    // 最高平均价
	IsFeatured         *bool   `form:"is_featured"`                                                       // 是否精选
	Keyword            string  `form:"keyword"`                                                           // 搜索关键词
	SortBy             string  `form:"sort_by" binding:"omitempty,oneof=name recent_transactions avg_price view_count"` // 排序字段
	SortOrder          string  `form:"sort_order" binding:"omitempty,oneof=asc desc"`                     // 排序方向
	Page               int     `form:"page" binding:"min=1"`                                              // 页码
	PageSize           int     `form:"page_size" binding:"min=1,max=100"`                                 // 每页数量
}

// GetEstatePropertiesRequest 获取屋苑内房源列表请求
type GetEstatePropertiesRequest struct {
	ListingType  *string  `form:"listing_type" binding:"omitempty,oneof=sale rent"` // sale=出售, rent=出租
	MinPrice     *float64 `form:"min_price"`                                        // 最低价格
	MaxPrice     *float64 `form:"max_price"`                                        // 最高价格
	Bedrooms     *int     `form:"bedrooms"`                                         // 房间数
	PropertyType *string  `form:"property_type"`                                    // 物业类型
	Status       string   `form:"status" binding:"omitempty,oneof=available pending sold cancelled"` // 状态
	Page         int      `form:"page" binding:"min=1"`                            // 页码
	PageSize     int      `form:"page_size" binding:"min=1,max=100"`               // 每页数量
}

// CreateEstateRequest 创建屋苑请求
type CreateEstateRequest struct {
	Name               string  `json:"name" binding:"required,max=200"`
	NameEn             string  `json:"name_en" binding:"omitempty,max=200"`
	Address            string  `json:"address" binding:"required,max=500"`
	DistrictID         uint    `json:"district_id" binding:"required"`
	TotalBlocks        int     `json:"total_blocks" binding:"omitempty,min=1"`
	TotalUnits         int     `json:"total_units" binding:"omitempty,min=1"`
	CompletionYear     int     `json:"completion_year" binding:"omitempty,min=1900,max=2100"`
	Developer          string  `json:"developer" binding:"omitempty,max=200"`
	ManagementCompany  string  `json:"management_company" binding:"omitempty,max=200"`
	PrimarySchoolNet   string  `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string  `json:"secondary_school_net" binding:"omitempty,max=50"`
	Description        string  `json:"description"`
	IsFeatured         bool    `json:"is_featured"`
	FacilityIDs        []uint  `json:"facility_ids"` // 设施ID列表
}

// UpdateEstateRequest 更新屋苑请求
type UpdateEstateRequest struct {
	Name               *string `json:"name" binding:"omitempty,max=200"`
	NameEn             *string `json:"name_en" binding:"omitempty,max=200"`
	Address            *string `json:"address" binding:"omitempty,max=500"`
	DistrictID         *uint   `json:"district_id"`
	TotalBlocks        *int    `json:"total_blocks" binding:"omitempty,min=1"`
	TotalUnits         *int    `json:"total_units" binding:"omitempty,min=1"`
	CompletionYear     *int    `json:"completion_year" binding:"omitempty,min=1900,max=2100"`
	Developer          *string `json:"developer" binding:"omitempty,max=200"`
	ManagementCompany  *string `json:"management_company" binding:"omitempty,max=200"`
	PrimarySchoolNet   *string `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet *string `json:"secondary_school_net" binding:"omitempty,max=50"`
	Description        *string `json:"description"`
	IsFeatured         *bool   `json:"is_featured"`
	FacilityIDs        []uint  `json:"facility_ids"` // 设施ID列表（覆盖更新）
}

// ============ Response DTO ============

// EstateResponse 屋苑响应
type EstateResponse struct {
	ID                      uint       `json:"id"`
	Name                    string     `json:"name"`
	NameEn                  string     `json:"name_en,omitempty"`
	Address                 string     `json:"address"`
	DistrictID              uint       `json:"district_id"`
	District                *District  `json:"district,omitempty"`
	TotalBlocks             int        `json:"total_blocks,omitempty"`
	TotalUnits              int        `json:"total_units,omitempty"`
	CompletionYear          int        `json:"completion_year,omitempty"`
	Developer               string     `json:"developer,omitempty"`
	ManagementCompany       string     `json:"management_company,omitempty"`
	PrimarySchoolNet        string     `json:"primary_school_net,omitempty"`
	SecondarySchoolNet      string     `json:"secondary_school_net,omitempty"`
	RecentTransactionsCount int        `json:"recent_transactions_count"`
	ForSaleCount            int        `json:"for_sale_count"`
	ForRentCount            int        `json:"for_rent_count"`
	AvgTransactionPrice     float64    `json:"avg_transaction_price,omitempty"`
	Description             string     `json:"description,omitempty"`
	ViewCount               int        `json:"view_count"`
	FavoriteCount           int        `json:"favorite_count"`
	IsFeatured              bool       `json:"is_featured"`
	Images                  []EstateImage `json:"images,omitempty"`
	Facilities              []Facility    `json:"facilities,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// EstateStatisticsResponse 屋苑统计数据响应
type EstateStatisticsResponse struct {
	EstateID                uint    `json:"estate_id"`
	EstateName              string  `json:"estate_name"`
	TotalUnits              int     `json:"total_units"`
	ForSaleCount            int     `json:"for_sale_count"`
	ForRentCount            int     `json:"for_rent_count"`
	RecentTransactionsCount int     `json:"recent_transactions_count"` // 最近3个月
	AvgTransactionPrice     float64 `json:"avg_transaction_price"`     // 港币/平方尺
	MinPrice                float64 `json:"min_price"`                 // 最低价格
	MaxPrice                float64 `json:"max_price"`                 // 最高价格
	AvgArea                 float64 `json:"avg_area"`                  // 平均面积
	ViewCount               int     `json:"view_count"`
	FavoriteCount           int     `json:"favorite_count"`
}

// PaginatedEstatesResponse 分页屋苑列表响应
type PaginatedEstatesResponse struct {
	Data       []EstateResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// ============ Facility Request DTO ============

// ListFacilitiesRequest 获取设施列表请求
type ListFacilitiesRequest struct {
	Category string `form:"category" binding:"omitempty,oneof=building unit"`
}

// CreateFacilityRequest 创建设施请求
type CreateFacilityRequest struct {
	NameZhHant string `json:"name_zh_hant" binding:"required,max=100"`
	NameZhHans string `json:"name_zh_hans" binding:"omitempty,max=100"`
	NameEn     string `json:"name_en" binding:"omitempty,max=100"`
	Icon       string `json:"icon" binding:"omitempty,max=100"`
	Category   string `json:"category" binding:"required,oneof=building unit"`
	SortOrder  int    `json:"sort_order"`
}

// UpdateFacilityRequest 更新设施请求
type UpdateFacilityRequest struct {
	NameZhHant *string `json:"name_zh_hant" binding:"omitempty,max=100"`
	NameZhHans *string `json:"name_zh_hans" binding:"omitempty,max=100"`
	NameEn     *string `json:"name_en" binding:"omitempty,max=100"`
	Icon       *string `json:"icon" binding:"omitempty,max=100"`
	Category   *string `json:"category" binding:"omitempty,oneof=building unit"`
	SortOrder  *int    `json:"sort_order"`
}

// ============ Facility Response DTO ============

// FacilityResponse 设施响应
type FacilityResponse struct {
	ID         uint      `json:"id"`
	NameZhHant string    `json:"name_zh_hant"`
	NameZhHans string    `json:"name_zh_hans,omitempty"`
	NameEn     string    `json:"name_en,omitempty"`
	Icon       string    `json:"icon,omitempty"`
	Category   string    `json:"category"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
