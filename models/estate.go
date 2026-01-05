package models

import (
	"time"

	"gorm.io/gorm"
)

// EstateImageType 屋苑图片类型常量
type EstateImageType string

const (
	EstateImageTypeExterior    EstateImageType = "exterior"    // 外观
	EstateImageTypeFacilities  EstateImageType = "facilities"  // 设施
	EstateImageTypeEnvironment EstateImageType = "environment" // 环境
	EstateImageTypeAerial      EstateImageType = "aerial"      // 航拍
)

// Estate 屋苑表模型
type Estate struct {
	ID                           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                         string         `gorm:"type:varchar(200);not null;index" json:"name"`
	NameEn                       *string        `gorm:"type:varchar(200)" json:"name_en,omitempty"`
	Address                      string         `gorm:"type:varchar(500);not null" json:"address"`
	DistrictID                   uint           `gorm:"not null;index" json:"district_id"`
	TotalBlocks                  *int           `json:"total_blocks,omitempty"`
	TotalUnits                   *int           `json:"total_units,omitempty"`
	CompletionYear               *int           `json:"completion_year,omitempty"`
	Developer                    *string        `gorm:"type:varchar(200)" json:"developer,omitempty"`
	ManagementCompany            *string        `gorm:"type:varchar(200)" json:"management_company,omitempty"`
	PrimarySchoolNet             *string        `gorm:"type:varchar(50);index" json:"primary_school_net,omitempty"`
	SecondarySchoolNet           *string        `gorm:"type:varchar(50);index" json:"secondary_school_net,omitempty"`
	RecentTransactionsCount      int            `gorm:"not null;default:0" json:"recent_transactions_count"`
	ForSaleCount                 int            `gorm:"not null;default:0" json:"for_sale_count"`
	ForRentCount                 int            `gorm:"not null;default:0" json:"for_rent_count"`
	AvgTransactionPrice          *float64       `gorm:"type:decimal(15,2);index" json:"avg_transaction_price,omitempty"`
	AvgTransactionPriceUpdatedAt *time.Time     `json:"avg_transaction_price_updated_at,omitempty"`
	Description                  *string        `gorm:"type:text" json:"description,omitempty"`
	ViewCount                    int            `gorm:"not null;default:0" json:"view_count"`
	FavoriteCount                int            `gorm:"not null;default:0" json:"favorite_count"`
	IsFeatured                   bool           `gorm:"not null;default:false;index" json:"is_featured"`
	CreatedAt                    time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt                    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	District   *District        `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Images     []EstateImage    `gorm:"foreignKey:EstateID" json:"images,omitempty"`
	Facilities []Facility       `gorm:"many2many:estate_facilities" json:"facilities,omitempty"`
}

// TableName 指定表名
func (Estate) TableName() string {
	return "estates"
}

// ============ Request DTO ============

// ListEstatesRequest 获取屋苑列表请求
type ListEstatesRequest struct {
	DistrictID        *uint    `form:"district_id"`
	Name              *string  `form:"name"`
	SchoolNet         *string  `form:"school_net"`
	MinYearBuilt      *int     `form:"min_year_built"`
	MaxYearBuilt      *int     `form:"max_year_built"`
	MinCompletionYear *int     `form:"min_completion_year"`
	MaxCompletionYear *int     `form:"max_completion_year"`
	MinPrice          *float64 `form:"min_price"`
	MaxPrice          *float64 `form:"max_price"`
	MinAvgPrice       *float64 `form:"min_avg_price"`
	MaxAvgPrice       *float64 `form:"max_avg_price"`
	HasListings       *bool    `form:"has_listings"`
	HasTransactions   *bool    `form:"has_transactions"`
	IsFeatured        *bool    `form:"is_featured"`
	Page              int      `form:"page,default=1" binding:"min=1"`
	PageSize          int      `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy            string   `form:"sort_by,default=name"`
	SortOrder         string   `form:"sort_order,default=asc" binding:"omitempty,oneof=asc desc"`
}

// ListValuationsRequest 获取估价列表请求
type ListValuationsRequest struct {
	DistrictID *uint    `form:"district_id"`
	EstateNo   *string  `form:"estate_no"`
	SchoolNet  *string  `form:"school_net"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Page       int      `form:"page,default=1" binding:"min=1"`
	PageSize   int      `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy     string   `form:"sort_by,default=created_at"`
	SortOrder  string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// SearchValuationsRequest 搜索估价请求
type SearchValuationsRequest struct {
	Keyword    string   `form:"keyword" binding:"required"`
	DistrictID *uint    `form:"district_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Page       int      `form:"page,default=1" binding:"min=1"`
	PageSize   int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// CheckIsFeatured 判断是否为精选屋苑
func (e *Estate) CheckIsFeatured() bool {
	return e.IsFeatured
}

// HasTransactions 判断是否有近期成交
func (e *Estate) HasTransactions() bool {
	return e.RecentTransactionsCount > 0
}

// HasListings 判断是否有放盘或租盘
func (e *Estate) HasListings() bool {
	return e.ForSaleCount > 0 || e.ForRentCount > 0
}

// GetTotalListings 获取总盘数（放盘+租盘）
func (e *Estate) GetTotalListings() int {
	return e.ForSaleCount + e.ForRentCount
}

// HasAvgPrice 判断是否有平均成交价数据
func (e *Estate) HasAvgPrice() bool {
	return e.AvgTransactionPrice != nil && *e.AvgTransactionPrice > 0
}

// IsAvgPriceRecent 判断平均成交价是否为近期数据（7天内）
func (e *Estate) IsAvgPriceRecent() bool {
	if e.AvgTransactionPriceUpdatedAt == nil {
		return false
	}
	return time.Since(*e.AvgTransactionPriceUpdatedAt) <= 7*24*time.Hour
}

// GetAge 获取屋苑楼龄
func (e *Estate) GetAge() int {
	if e.CompletionYear == nil {
		return 0
	}
	currentYear := time.Now().Year()
	age := currentYear - *e.CompletionYear
	if age < 0 {
		return 0
	}
	return age
}

// EstateImage 屋苑图片表模型
type EstateImage struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	EstateID  uint            `gorm:"not null;index" json:"estate_id"`
	ImageURL  string          `gorm:"type:varchar(500);not null" json:"image_url"`
	ImageType EstateImageType `gorm:"type:varchar(20);not null" json:"image_type"`
	Title     *string         `gorm:"type:varchar(200)" json:"title,omitempty"`
	SortOrder int             `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Estate *Estate `gorm:"foreignKey:EstateID;constraint:OnDelete:CASCADE" json:"estate,omitempty"`
}

// TableName 指定表名
func (EstateImage) TableName() string {
	return "estate_images"
}

// IsExterior 判断是否为外观图
func (ei *EstateImage) IsExterior() bool {
	return ei.ImageType == EstateImageTypeExterior
}

// IsFacility 判断是否为设施图
func (ei *EstateImage) IsFacility() bool {
	return ei.ImageType == EstateImageTypeFacilities
}

// IsEnvironment 判断是否为环境图
func (ei *EstateImage) IsEnvironment() bool {
	return ei.ImageType == EstateImageTypeEnvironment
}

// IsAerial 判断是否为航拍图
func (ei *EstateImage) IsAerial() bool {
	return ei.ImageType == EstateImageTypeAerial
}

// EstateFacility 屋苑设施关联表模型
type EstateFacility struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EstateID   uint      `gorm:"not null;index;uniqueIndex:idx_estate_facility" json:"estate_id"`
	FacilityID uint      `gorm:"not null;index;uniqueIndex:idx_estate_facility" json:"facility_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Estate   *Estate   `gorm:"foreignKey:EstateID;constraint:OnDelete:CASCADE" json:"estate,omitempty"`
	Facility *Facility `gorm:"foreignKey:FacilityID;constraint:OnDelete:CASCADE" json:"facility,omitempty"`
}

// TableName 指定表名
func (EstateFacility) TableName() string {
	return "estate_facilities"
}
