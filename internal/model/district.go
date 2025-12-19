package model

import (
	"time"
)

// Region 香港区域常量
type Region string

const (
	RegionHKIsland        Region = "HK_ISLAND"        // 港岛
	RegionKowloon         Region = "KOWLOON"          // 九龙
	RegionNewTerritories  Region = "NEW_TERRITORIES"  // 新界
)

// District 地区表模型
type District struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NameZhHant  string    `gorm:"type:varchar(100);not null" json:"name_zh_hant"` // 中文繁体
	NameZhHans  *string   `gorm:"type:varchar(100)" json:"name_zh_hans,omitempty"` // 中文简体
	NameEn      *string   `gorm:"type:varchar(100)" json:"name_en,omitempty"`      // 英文
	Region      Region    `gorm:"type:varchar(50);not null;index" json:"region"`
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (District) TableName() string {
	return "districts"
}

// GetLocalizedName 根据语言获取本地化名称
func (d *District) GetLocalizedName(lang string) string {
	switch lang {
	case "zh-Hans", "zh_CN":
		if d.NameZhHans != nil {
			return *d.NameZhHans
		}
		return d.NameZhHant
	case "en":
		if d.NameEn != nil {
			return *d.NameEn
		}
		return d.NameZhHant
	default: // zh-Hant, zh_HK
		return d.NameZhHant
	}
}

// IsHKIsland 判断是否为港岛区
func (d *District) IsHKIsland() bool {
	return d.Region == RegionHKIsland
}

// IsKowloon 判断是否为九龙区
func (d *District) IsKowloon() bool {
	return d.Region == RegionKowloon
}

// IsNewTerritories 判断是否为新界区
func (d *District) IsNewTerritories() bool {
	return d.Region == RegionNewTerritories
}
