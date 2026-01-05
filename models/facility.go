package models

import (
	"time"
)

// FacilityCategory 设施分类常量
type FacilityCategory string

const (
	FacilityCategoryBuilding FacilityCategory = "building" // 大厦设施
	FacilityCategoryUnit     FacilityCategory = "unit"     // 单位设施
)

// Facility 设施字典表模型
type Facility struct {
	ID         uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	NameZhHant string           `gorm:"type:varchar(100);not null" json:"name_zh_hant"`
	NameZhHans *string          `gorm:"type:varchar(100)" json:"name_zh_hans,omitempty"`
	NameEn     *string          `gorm:"type:varchar(100)" json:"name_en,omitempty"`
	Icon       *string          `gorm:"type:varchar(100)" json:"icon,omitempty"`
	Category   FacilityCategory `gorm:"type:varchar(50);not null;index" json:"category"`
	SortOrder  int              `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt  time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Facility) TableName() string {
	return "facilities"
}

// ============ Request DTO ============

// ListFacilitiesRequest 获取设施列表请求
type ListFacilitiesRequest struct {
	Keyword   *string `form:"keyword"`
	Category  *string `form:"category"`
	IsActive  *bool   `form:"is_active"`
	Page      int     `form:"page,default=1" binding:"min=1"`
	PageSize  int     `form:"page_size,default=50" binding:"min=1,max=100"`
	SortBy    string  `form:"sort_by,default=sort_order"`
	SortOrder string  `form:"sort_order,default=asc" binding:"omitempty,oneof=asc desc"`
}

// GetLocalizedName 根据语言获取本地化名称
func (f *Facility) GetLocalizedName(lang string) string {
	switch lang {
	case "zh-Hans", "zh_CN":
		if f.NameZhHans != nil {
			return *f.NameZhHans
		}
		return f.NameZhHant
	case "en":
		if f.NameEn != nil {
			return *f.NameEn
		}
		return f.NameZhHant
	default: // zh-Hant, zh_HK
		return f.NameZhHant
	}
}

// IsBuildingFacility 判断是否为大厦设施
func (f *Facility) IsBuildingFacility() bool {
	return f.Category == FacilityCategoryBuilding
}

// IsUnitFacility 判断是否为单位设施
func (f *Facility) IsUnitFacility() bool {
	return f.Category == FacilityCategoryUnit
}

// PropertyFacility 房产设施关联表模型
type PropertyFacility struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PropertyID uint      `gorm:"not null;index;uniqueIndex:idx_property_facility" json:"property_id"`
	FacilityID uint      `gorm:"not null;index;uniqueIndex:idx_property_facility" json:"facility_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Property *Property `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"property,omitempty"`
	Facility *Facility `gorm:"foreignKey:FacilityID;constraint:OnDelete:CASCADE" json:"facility,omitempty"`
}

// TableName 指定表名
func (PropertyFacility) TableName() string {
	return "property_facilities"
}
