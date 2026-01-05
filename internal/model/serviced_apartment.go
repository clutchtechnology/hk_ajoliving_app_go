package model

import (
	"time"

	"gorm.io/gorm"
)

// ServicedApartmentStatus 服务式住宅状态常量
type ServicedApartmentStatus string

const (
	ServicedApartmentStatusActive   ServicedApartmentStatus = "active"   // 营业中
	ServicedApartmentStatusInactive ServicedApartmentStatus = "inactive" // 暂停营业
	ServicedApartmentStatusClosed   ServicedApartmentStatus = "closed"   // 已关闭
)

// ServicedApartmentImageType 服务式住宅图片类型常量
type ServicedApartmentImageType string

const (
	ServicedApartmentImageTypeExterior  ServicedApartmentImageType = "exterior"  // 外观
	ServicedApartmentImageTypeLobby     ServicedApartmentImageType = "lobby"     // 大堂
	ServicedApartmentImageTypeRoom      ServicedApartmentImageType = "room"      // 房间
	ServicedApartmentImageTypeBathroom  ServicedApartmentImageType = "bathroom"  // 浴室
	ServicedApartmentImageTypeFacility  ServicedApartmentImageType = "facility"  // 设施
)

// ServicedApartment 服务式住宅表模型
type ServicedApartment struct {
	ID            uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string                  `gorm:"type:varchar(200);not null;index" json:"name"`
	NameEn        *string                 `gorm:"type:varchar(200)" json:"name_en,omitempty"`
	Address       string                  `gorm:"type:varchar(500);not null" json:"address"`
	DistrictID    uint                    `gorm:"not null;index" json:"district_id"`
	Description   *string                 `gorm:"type:text" json:"description,omitempty"`
	Phone         string                  `gorm:"type:varchar(50);not null" json:"phone"`
	WebsiteURL    *string                 `gorm:"type:varchar(500)" json:"website_url,omitempty"`
	Email         *string                 `gorm:"type:varchar(255)" json:"email,omitempty"`
	CompanyID     uint                    `gorm:"not null;index" json:"company_id"`
	CheckInTime   *string                 `gorm:"type:varchar(50)" json:"check_in_time,omitempty"`
	CheckOutTime  *string                 `gorm:"type:varchar(50)" json:"check_out_time,omitempty"`
	MinStayDays   *int                    `json:"min_stay_days,omitempty"`
	Status        ServicedApartmentStatus `gorm:"type:varchar(20);not null;index;default:'active'" json:"status"`
	Rating        *float64                `gorm:"type:decimal(3,2)" json:"rating,omitempty"`
	ReviewCount   int                     `gorm:"not null;default:0" json:"review_count"`
	ViewCount     int                     `gorm:"not null;default:0" json:"view_count"`
	FavoriteCount int                     `gorm:"not null;default:0" json:"favorite_count"`
	IsFeatured    bool                    `gorm:"not null;default:false;index" json:"is_featured"`
	CreatedAt     time.Time               `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt     time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt          `gorm:"index" json:"-"`

	// 关联
	District   *District                    `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Company    *User                        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Units      []ServicedApartmentUnit      `gorm:"foreignKey:ServicedApartmentID" json:"units,omitempty"`
	Images     []ServicedApartmentImage     `gorm:"foreignKey:ServicedApartmentID" json:"images,omitempty"`
	Facilities []Facility                   `gorm:"many2many:serviced_apartment_facilities" json:"facilities,omitempty"`
}

// TableName 指定表名
func (ServicedApartment) TableName() string {
	return "serviced_apartments"
}

// IsActive 判断是否营业中
func (sa *ServicedApartment) IsActive() bool {
	return sa.Status == ServicedApartmentStatusActive
}

// IsClosed 判断是否已关闭
func (sa *ServicedApartment) IsClosed() bool {
	return sa.Status == ServicedApartmentStatusClosed
}

// CheckIsFeatured 判断是否为精选推荐
func (sa *ServicedApartment) CheckIsFeatured() bool {
	return sa.IsFeatured
}

// HasWebsite 判断是否有官网
func (sa *ServicedApartment) HasWebsite() bool {
	return sa.WebsiteURL != nil && *sa.WebsiteURL != ""
}

// GetMinPrice 获取最低价格（从房型中）
func (sa *ServicedApartment) GetMinPrice() float64 {
	if len(sa.Units) == 0 {
		return 0
	}
	minPrice := sa.Units[0].MonthlyPrice
	for _, unit := range sa.Units {
		if unit.MonthlyPrice < minPrice {
			minPrice = unit.MonthlyPrice
		}
	}
	return minPrice
}

// BeforeCreate GORM hook - 创建前执行
func (sa *ServicedApartment) BeforeCreate(tx *gorm.DB) error {
	if sa.Status == "" {
		sa.Status = ServicedApartmentStatusActive
	}
	return nil
}

// ServicedApartmentUnit 服务式住宅房型表模型
type ServicedApartmentUnit struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServicedApartmentID  uint      `gorm:"not null;index" json:"serviced_apartment_id"`
	UnitType             string    `gorm:"type:varchar(50);not null" json:"unit_type"`
	Bedrooms             int       `gorm:"not null" json:"bedrooms"`
	Bathrooms            *int      `json:"bathrooms,omitempty"`
	Area                 float64   `gorm:"type:decimal(10,2);not null" json:"area"` // 平方尺
	MaxOccupancy         int       `gorm:"not null" json:"max_occupancy"`
	DailyPrice           *float64  `gorm:"type:decimal(10,2)" json:"daily_price,omitempty"`
	WeeklyPrice          *float64  `gorm:"type:decimal(10,2)" json:"weekly_price,omitempty"`
	MonthlyPrice         float64   `gorm:"type:decimal(10,2);not null;index" json:"monthly_price"`
	AvailableUnits       int       `gorm:"not null" json:"available_units"`
	Description          *string   `gorm:"type:text" json:"description,omitempty"`
	SortOrder            int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	ServicedApartment *ServicedApartment         `gorm:"foreignKey:ServicedApartmentID;constraint:OnDelete:CASCADE" json:"serviced_apartment,omitempty"`
	Images            []ServicedApartmentImage   `gorm:"foreignKey:UnitID" json:"images,omitempty"`
}

// TableName 指定表名
func (ServicedApartmentUnit) TableName() string {
	return "serviced_apartment_units"
}

// IsAvailable 判断是否有可用单位
func (sau *ServicedApartmentUnit) IsAvailable() bool {
	return sau.AvailableUnits > 0
}

// GetPricePerSqft 计算每平方尺月租
func (sau *ServicedApartmentUnit) GetPricePerSqft() float64 {
	if sau.Area == 0 {
		return 0
	}
	return sau.MonthlyPrice / sau.Area
}

// HasDailyPrice 判断是否提供日租
func (sau *ServicedApartmentUnit) HasDailyPrice() bool {
	return sau.DailyPrice != nil && *sau.DailyPrice > 0
}

// HasWeeklyPrice 判断是否提供周租
func (sau *ServicedApartmentUnit) HasWeeklyPrice() bool {
	return sau.WeeklyPrice != nil && *sau.WeeklyPrice > 0
}

// ServicedApartmentImage 服务式住宅图片表模型
type ServicedApartmentImage struct {
	ID                  uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	ServicedApartmentID *uint                      `gorm:"index" json:"serviced_apartment_id,omitempty"`
	UnitID              *uint                      `gorm:"index" json:"unit_id,omitempty"`
	ImageURL            string                     `gorm:"type:varchar(500);not null" json:"image_url"`
	ImageType           ServicedApartmentImageType `gorm:"type:varchar(20);not null" json:"image_type"`
	Title               *string                    `gorm:"type:varchar(200)" json:"title,omitempty"`
	SortOrder           int                        `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt           time.Time                  `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	ServicedApartment *ServicedApartment      `gorm:"foreignKey:ServicedApartmentID;constraint:OnDelete:CASCADE" json:"serviced_apartment,omitempty"`
	Unit              *ServicedApartmentUnit  `gorm:"foreignKey:UnitID;constraint:OnDelete:CASCADE" json:"unit,omitempty"`
}

// TableName 指定表名
func (ServicedApartmentImage) TableName() string {
	return "serviced_apartment_images"
}

// IsApartmentImage 判断是否为公寓整体图片
func (sai *ServicedApartmentImage) IsApartmentImage() bool {
	return sai.ServicedApartmentID != nil
}

// IsUnitImage 判断是否为房型图片
func (sai *ServicedApartmentImage) IsUnitImage() bool {
	return sai.UnitID != nil
}

// ServicedApartmentFacility 服务式住宅设施关联表模型
type ServicedApartmentFacility struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServicedApartmentID uint      `gorm:"not null;index;uniqueIndex:idx_serviced_apt_facility" json:"serviced_apartment_id"`
	FacilityID          uint      `gorm:"not null;index;uniqueIndex:idx_serviced_apt_facility" json:"facility_id"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	ServicedApartment *ServicedApartment `gorm:"foreignKey:ServicedApartmentID;constraint:OnDelete:CASCADE" json:"serviced_apartment,omitempty"`
	Facility          *Facility          `gorm:"foreignKey:FacilityID;constraint:OnDelete:CASCADE" json:"facility,omitempty"`
}

// TableName 指定表名
func (ServicedApartmentFacility) TableName() string {
	return "serviced_apartment_facilities"
}
