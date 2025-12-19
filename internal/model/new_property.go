package model

import (
	"time"

	"gorm.io/gorm"
)

// NewPropertyStatus 新盘状态常量
type NewPropertyStatus string

const (
	NewPropertyStatusUpcoming  NewPropertyStatus = "upcoming"  // 即将推出
	NewPropertyStatusPresale   NewPropertyStatus = "presale"   // 预售中
	NewPropertyStatusSelling   NewPropertyStatus = "selling"   // 销售中
	NewPropertyStatusCompleted NewPropertyStatus = "completed" // 已完成
)

// NewProperty 新盘表模型
type NewProperty struct {
	ID                   uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                 string            `gorm:"type:varchar(200);not null;index" json:"name"`
	NameEn               *string           `gorm:"type:varchar(200)" json:"name_en,omitempty"`
	Address              string            `gorm:"type:varchar(500);not null" json:"address"`
	DistrictID           uint              `gorm:"not null;index" json:"district_id"`
	Status               NewPropertyStatus `gorm:"type:varchar(20);not null;index" json:"status"`
	UnitsForSale         *int              `json:"units_for_sale,omitempty"`
	UnitsSold            *int              `json:"units_sold,omitempty"`
	Developer            string            `gorm:"type:varchar(200);not null;index" json:"developer"`
	ManagementCompany    *string           `gorm:"type:varchar(200)" json:"management_company,omitempty"`
	TotalUnits           int               `gorm:"not null" json:"total_units"`
	TotalBlocks          int               `gorm:"not null" json:"total_blocks"`
	MaxFloors            int               `gorm:"not null" json:"max_floors"`
	PrimarySchoolNet     *string           `gorm:"type:varchar(50);index" json:"primary_school_net,omitempty"`
	SecondarySchoolNet   *string           `gorm:"type:varchar(50);index" json:"secondary_school_net,omitempty"`
	WebsiteURL           *string           `gorm:"type:varchar(500)" json:"website_url,omitempty"`
	SalesOfficeAddress   *string           `gorm:"type:varchar(500)" json:"sales_office_address,omitempty"`
	SalesPhone           *string           `gorm:"type:varchar(50)" json:"sales_phone,omitempty"`
	ExpectedCompletion   *time.Time        `gorm:"type:date" json:"expected_completion,omitempty"`
	OccupationDate       *time.Time        `gorm:"type:date" json:"occupation_date,omitempty"`
	Description          *string           `gorm:"type:text" json:"description,omitempty"`
	ViewCount            int               `gorm:"not null;default:0" json:"view_count"`
	FavoriteCount        int               `gorm:"not null;default:0" json:"favorite_count"`
	SortOrder            int               `gorm:"not null;default:0" json:"sort_order"`
	IsFeatured           bool              `gorm:"not null;default:false;index" json:"is_featured"`
	CreatedAt            time.Time         `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt            time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt    `gorm:"index" json:"-"`

	// 关联
	District *District            `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Images   []NewPropertyImage   `gorm:"foreignKey:NewPropertyID" json:"images,omitempty"`
	Layouts  []NewPropertyLayout  `gorm:"foreignKey:NewPropertyID" json:"layouts,omitempty"`
}

// TableName 指定表名
func (NewProperty) TableName() string {
	return "new_properties"
}

// IsUpcoming 判断是否即将推出
func (np *NewProperty) IsUpcoming() bool {
	return np.Status == NewPropertyStatusUpcoming
}

// IsPresale 判断是否预售中
func (np *NewProperty) IsPresale() bool {
	return np.Status == NewPropertyStatusPresale
}

// IsSelling 判断是否销售中
func (np *NewProperty) IsSelling() bool {
	return np.Status == NewPropertyStatusSelling
}

// IsCompleted 判断是否已完成
func (np *NewProperty) IsCompleted() bool {
	return np.Status == NewPropertyStatusCompleted
}

// IsFeatured 判断是否为精选
func (np *NewProperty) IsFeatured() bool {
	return np.IsFeatured
}

// GetSalesProgress 获取销售进度百分比
func (np *NewProperty) GetSalesProgress() float64 {
	if np.UnitsForSale == nil || np.UnitsSold == nil {
		return 0
	}
	totalUnits := *np.UnitsForSale + *np.UnitsSold
	if totalUnits == 0 {
		return 0
	}
	return float64(*np.UnitsSold) / float64(totalUnits) * 100
}

// GetAvailableUnits 获取可售单位数
func (np *NewProperty) GetAvailableUnits() int {
	if np.UnitsForSale == nil {
		return 0
	}
	return *np.UnitsForSale
}

// NewPropertyImage 新盘图片表模型
type NewPropertyImage struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NewPropertyID uint      `gorm:"not null;index" json:"new_property_id"`
	ImageURL      string    `gorm:"type:varchar(500);not null" json:"image_url"`
	ImageType     ImageType `gorm:"type:varchar(20);not null" json:"image_type"`
	Title         *string   `gorm:"type:varchar(200)" json:"title,omitempty"`
	SortOrder     int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	NewProperty *NewProperty `gorm:"foreignKey:NewPropertyID;constraint:OnDelete:CASCADE" json:"new_property,omitempty"`
}

// TableName 指定表名
func (NewPropertyImage) TableName() string {
	return "new_property_images"
}

// NewPropertyLayout 新盘户型表模型
type NewPropertyLayout struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NewPropertyID  uint      `gorm:"not null;index" json:"new_property_id"`
	UnitType       string    `gorm:"type:varchar(50);not null" json:"unit_type"`
	Bedrooms       int       `gorm:"not null" json:"bedrooms"`
	Bathrooms      *int      `json:"bathrooms,omitempty"`
	SaleableArea   float64   `gorm:"type:decimal(10,2);not null" json:"saleable_area"` // 实用面积
	GrossArea      *float64  `gorm:"type:decimal(10,2)" json:"gross_area,omitempty"`   // 建筑面积
	MinPrice       float64   `gorm:"type:decimal(15,2);not null" json:"min_price"`
	MaxPrice       *float64  `gorm:"type:decimal(15,2)" json:"max_price,omitempty"`
	PricePerSqft   *float64  `gorm:"type:decimal(10,2)" json:"price_per_sqft,omitempty"`
	AvailableUnits int       `gorm:"not null" json:"available_units"`
	FloorplanURL   *string   `gorm:"type:varchar(500)" json:"floorplan_url,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	NewProperty *NewProperty `gorm:"foreignKey:NewPropertyID;constraint:OnDelete:CASCADE" json:"new_property,omitempty"`
}

// TableName 指定表名
func (NewPropertyLayout) TableName() string {
	return "new_property_layouts"
}

// GetPriceRange 获取价格区间字符串
func (npl *NewPropertyLayout) GetPriceRange() string {
	if npl.MaxPrice == nil {
		return ""
	}
	return "" // 可以在这里实现格式化逻辑
}

// GetAvgPrice 获取平均价格
func (npl *NewPropertyLayout) GetAvgPrice() float64 {
	if npl.MaxPrice == nil {
		return npl.MinPrice
	}
	return (npl.MinPrice + *npl.MaxPrice) / 2
}

// IsAvailable 判断是否有可售单位
func (npl *NewPropertyLayout) IsAvailable() bool {
	return npl.AvailableUnits > 0
}

// CalculatePricePerSqft 计算每平方尺价格
func (npl *NewPropertyLayout) CalculatePricePerSqft() float64 {
	if npl.SaleableArea == 0 {
		return 0
	}
	avgPrice := npl.GetAvgPrice()
	return avgPrice / npl.SaleableArea
}
