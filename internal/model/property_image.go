package model

import (
	"time"
)

// PropertyImage 房产图片表模型
type PropertyImage struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PropertyID uint      `gorm:"not null;index" json:"property_id"`
	URL        string    `gorm:"type:varchar(500);not null;column:image_url" json:"url"`
	Caption    *string   `gorm:"type:varchar(255)" json:"caption,omitempty"`
	ImageType  ImageType `gorm:"type:varchar(20);not null" json:"image_type"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	IsCover    bool      `gorm:"not null;default:false" json:"is_cover"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Property *Property `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"property,omitempty"`
}

// TableName 指定表名
func (PropertyImage) TableName() string {
	return "property_images"
}

// IsCover 判断是否为封面图
func (pi *PropertyImage) IsCover() bool {
	return pi.ImageType == ImageTypeCover
}

// IsInterior 判断是否为室内图
func (pi *PropertyImage) IsInterior() bool {
	return pi.ImageType == ImageTypeInterior
}

// IsExterior 判断是否为外观图
func (pi *PropertyImage) IsExterior() bool {
	return pi.ImageType == ImageTypeExterior
}

// IsFloorplan 判断是否为户型图
func (pi *PropertyImage) IsFloorplan() bool {
	return pi.ImageType == ImageTypeFloorplan
}
