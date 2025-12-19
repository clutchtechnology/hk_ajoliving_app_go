package model

import (
	"time"

	"gorm.io/gorm"
)

// IndexType 指数类型常量
type IndexType string

const (
	IndexTypeOverall   IndexType = "overall"   // 整体指数
	IndexTypeDistrict  IndexType = "district"  // 地区指数
	IndexTypeEstate    IndexType = "estate"    // 屋苑指数
	IndexTypePropertyType IndexType = "property_type" // 物业类型指数
)

// PriceIndex 楼价指数表模型
type PriceIndex struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	IndexType      IndexType      `gorm:"type:varchar(20);not null;index" json:"index_type"`
	DistrictID     *uint          `gorm:"index" json:"district_id,omitempty"`
	EstateID       *uint          `gorm:"index" json:"estate_id,omitempty"`
	PropertyType   *string        `gorm:"type:varchar(50);index" json:"property_type,omitempty"`
	IndexValue     float64        `gorm:"type:decimal(10,2);not null" json:"index_value"`
	ChangeValue    float64        `gorm:"type:decimal(10,2)" json:"change_value"`
	ChangePercent  float64        `gorm:"type:decimal(5,2)" json:"change_percent"`
	AvgPrice       *float64       `gorm:"type:decimal(15,2)" json:"avg_price,omitempty"`
	AvgPricePerSqft *float64      `gorm:"type:decimal(10,2)" json:"avg_price_per_sqft,omitempty"`
	TransactionCount int           `gorm:"not null;default:0" json:"transaction_count"`
	Period         string         `gorm:"type:varchar(20);not null;index" json:"period"` // YYYY-MM 或 YYYY-MM-DD
	Year           int            `gorm:"not null;index" json:"year"`
	Month          int            `gorm:"not null;index" json:"month"`
	Day            *int           `gorm:"index" json:"day,omitempty"`
	DataSource     string         `gorm:"type:varchar(100)" json:"data_source"`
	Notes          *string        `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	District *District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Estate   *Estate   `gorm:"foreignKey:EstateID" json:"estate,omitempty"`
}

// TableName 指定表名
func (PriceIndex) TableName() string {
	return "price_indices"
}

// IsPositiveChange 判断是否为正增长
func (pi *PriceIndex) IsPositiveChange() bool {
	return pi.ChangeValue > 0
}

// IsNegativeChange 判断是否为负增长
func (pi *PriceIndex) IsNegativeChange() bool {
	return pi.ChangeValue < 0
}

// GetPeriodString 获取周期字符串
func (pi *PriceIndex) GetPeriodString() string {
	return pi.Period
}

// BeforeCreate GORM hook - 创建前执行
func (pi *PriceIndex) BeforeCreate(tx *gorm.DB) error {
	// 自动设置 Year 和 Month
	if pi.Period != "" {
		// 解析 Period 字符串设置 Year 和 Month
		// 这里可以添加解析逻辑
	}
	return nil
}
