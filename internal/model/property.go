package model

import (
	"time"

	"gorm.io/gorm"
)

// ListingType 房源类型常量
type ListingType string

const (
	ListingTypeSale ListingType = "sale" // 出售
	ListingTypeRent ListingType = "rent" // 出租
)

// PropertyType 物业类型常量
type PropertyType string

const (
	PropertyTypeApartment  PropertyType = "apartment"  // 公寓
	PropertyTypeVilla      PropertyType = "villa"      // 别墅
	PropertyTypeTownhouse  PropertyType = "townhouse"  // 联排别墅
	PropertyTypeStudio     PropertyType = "studio"     // 开放式单位
	PropertyTypeDuplex     PropertyType = "duplex"     // 复式
	PropertyTypePenthouse  PropertyType = "penthouse"  // 顶层豪宅
	PropertyTypeShophouse  PropertyType = "shophouse"  // 商住两用
)

// PropertyStatus 房产状态常量
type PropertyStatus string

const (
	PropertyStatusAvailable PropertyStatus = "available" // 可用
	PropertyStatusPending   PropertyStatus = "pending"   // 待定
	PropertyStatusSold      PropertyStatus = "sold"      // 已售/已租
	PropertyStatusCancelled PropertyStatus = "cancelled" // 已取消
)

// Property 房产表模型
type Property struct {
	ID                 uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	PropertyNo         string          `gorm:"type:varchar(50);not null;uniqueIndex" json:"property_no"`
	EstateNo           *string         `gorm:"type:varchar(50);index" json:"estate_no,omitempty"`
	ListingType        ListingType     `gorm:"type:varchar(20);not null;index" json:"listing_type"`
	Title              string          `gorm:"type:varchar(255);not null" json:"title"`
	Description        *string         `gorm:"type:text" json:"description,omitempty"`
	Area               float64         `gorm:"type:decimal(10,2);not null" json:"area"` // 平方尺
	Price              float64         `gorm:"type:decimal(15,2);not null;index" json:"price"`
	Address            string          `gorm:"type:varchar(500);not null" json:"address"`
	DistrictID         uint            `gorm:"not null;index" json:"district_id"`
	BuildingName       *string         `gorm:"type:varchar(200);index" json:"building_name,omitempty"`
	Floor              *string         `gorm:"type:varchar(20)" json:"floor,omitempty"`
	Orientation        *string         `gorm:"type:varchar(50)" json:"orientation,omitempty"`
	Bedrooms           int             `gorm:"not null;index" json:"bedrooms"`
	Bathrooms          *int            `json:"bathrooms,omitempty"`
	PrimarySchoolNet   *string         `gorm:"type:varchar(50);index" json:"primary_school_net,omitempty"`
	SecondarySchoolNet *string         `gorm:"type:varchar(50);index" json:"secondary_school_net,omitempty"`
	PropertyType       PropertyType    `gorm:"type:varchar(50);not null;index" json:"property_type"`
	Status             PropertyStatus  `gorm:"type:varchar(20);not null;index;default:'available'" json:"status"`
	PublisherID        uint            `gorm:"not null;index" json:"publisher_id"`
	PublisherType      PublisherType   `gorm:"type:varchar(20);not null" json:"publisher_type"`
	AgentID            *uint           `gorm:"index" json:"agent_id,omitempty"`
	ViewCount          int             `gorm:"not null;default:0" json:"view_count"`
	FavoriteCount      int             `gorm:"not null;default:0" json:"favorite_count"`
	PublishedAt        *time.Time      `gorm:"index" json:"published_at,omitempty"`
	ExpiredAt          *time.Time      `gorm:"index" json:"expired_at,omitempty"`
	CreatedAt          time.Time       `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`

	// 关联
	District   *District          `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Publisher  *User              `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	Agent      *Agent             `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Images     []PropertyImage    `gorm:"foreignKey:PropertyID" json:"images,omitempty"`
	Facilities []Facility         `gorm:"many2many:property_facilities" json:"facilities,omitempty"`
}

// TableName 指定表名
func (Property) TableName() string {
	return "properties"
}

// IsSale 判断是否为出售
func (p *Property) IsSale() bool {
	return p.ListingType == ListingTypeSale
}

// IsRent 判断是否为出租
func (p *Property) IsRent() bool {
	return p.ListingType == ListingTypeRent
}

// IsAvailable 判断是否可用
func (p *Property) IsAvailable() bool {
	return p.Status == PropertyStatusAvailable
}

// IsSold 判断是否已售出/已租出
func (p *Property) IsSold() bool {
	return p.Status == PropertyStatusSold
}

// IsPublished 判断是否已发布
func (p *Property) IsPublished() bool {
	return p.PublishedAt != nil && !p.PublishedAt.After(time.Now())
}

// IsExpired 判断是否已过期
func (p *Property) IsExpired() bool {
	return p.ExpiredAt != nil && p.ExpiredAt.Before(time.Now())
}

// GetPricePerSqft 计算每平方尺价格
func (p *Property) GetPricePerSqft() float64 {
	if p.Area == 0 {
		return 0
	}
	return p.Price / p.Area
}

// BeforeCreate GORM hook - 创建前执行
func (p *Property) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = PropertyStatusAvailable
	}
	// 如果没有设置发布时间，默认为当前时间
	if p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}
