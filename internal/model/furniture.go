package model

import (
	"time"

	"gorm.io/gorm"
)

// FurnitureCondition 家具新旧程度常量
type FurnitureCondition string

const (
	FurnitureConditionNew      FurnitureCondition = "new"      // 全新
	FurnitureConditionLikeNew  FurnitureCondition = "like_new" // 近全新
	FurnitureConditionGood     FurnitureCondition = "good"     // 良好
	FurnitureConditionFair     FurnitureCondition = "fair"     // 一般
	FurnitureConditionPoor     FurnitureCondition = "poor"     // 较差
)

// DeliveryMethod 交收方法常量
type DeliveryMethod string

const (
	DeliveryMethodSelfPickup  DeliveryMethod = "self_pickup" // 自取
	DeliveryMethodDelivery    DeliveryMethod = "delivery"    // 送货
	DeliveryMethodNegotiable  DeliveryMethod = "negotiable"  // 面议
)

// FurnitureStatus 家具状态常量
type FurnitureStatus string

const (
	FurnitureStatusAvailable FurnitureStatus = "available" // 可用
	FurnitureStatusReserved  FurnitureStatus = "reserved"  // 已预订
	FurnitureStatusSold      FurnitureStatus = "sold"      // 已售出
	FurnitureStatusExpired   FurnitureStatus = "expired"   // 已过期
	FurnitureStatusCancelled FurnitureStatus = "cancelled" // 已取消
)

// Furniture 家具表模型
type Furniture struct {
	ID                 uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	FurnitureNo        string             `gorm:"type:varchar(50);not null;uniqueIndex" json:"furniture_no"`
	Title              string             `gorm:"type:varchar(255);not null;index" json:"title"`
	Description        *string            `gorm:"type:text" json:"description,omitempty"`
	Price              float64            `gorm:"type:decimal(10,2);not null;index" json:"price"`
	CategoryID         uint               `gorm:"not null;index" json:"category_id"`
	Brand              *string            `gorm:"type:varchar(100);index" json:"brand,omitempty"`
	Condition          FurnitureCondition `gorm:"type:varchar(20);not null;index" json:"condition"`
	PurchaseDate       *time.Time         `gorm:"type:date" json:"purchase_date,omitempty"`
	DeliveryDistrictID uint               `gorm:"not null;index" json:"delivery_district_id"`
	DeliveryTime       *string            `gorm:"type:varchar(100)" json:"delivery_time,omitempty"`
	DeliveryMethod     DeliveryMethod     `gorm:"type:varchar(50);not null;index" json:"delivery_method"`
	Status             FurnitureStatus    `gorm:"type:varchar(20);not null;index;default:'available'" json:"status"`
	PublisherID        uint               `gorm:"not null;index" json:"publisher_id"`
	PublisherType      PublisherType      `gorm:"type:varchar(20);not null" json:"publisher_type"`
	ViewCount          int                `gorm:"not null;default:0" json:"view_count"`
	FavoriteCount      int                `gorm:"not null;default:0" json:"favorite_count"`
	PublishedAt        time.Time          `gorm:"not null;index" json:"published_at"`
	UpdatedAt          time.Time          `gorm:"autoUpdateTime;index" json:"updated_at"`
	ExpiresAt          time.Time          `gorm:"not null;index" json:"expires_at"`
	CreatedAt          time.Time          `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`

	// 关联
	Category         *FurnitureCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Publisher        *User              `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	DeliveryDistrict *District          `gorm:"foreignKey:DeliveryDistrictID" json:"delivery_district,omitempty"`
	Images           []FurnitureImage   `gorm:"foreignKey:FurnitureID" json:"images,omitempty"`
}

// TableName 指定表名
func (Furniture) TableName() string {
	return "furniture"
}

// IsAvailable 判断是否可用
func (f *Furniture) IsAvailable() bool {
	return f.Status == FurnitureStatusAvailable && !f.IsExpired()
}

// IsReserved 判断是否已预订
func (f *Furniture) IsReserved() bool {
	return f.Status == FurnitureStatusReserved
}

// IsSold 判断是否已售出
func (f *Furniture) IsSold() bool {
	return f.Status == FurnitureStatusSold
}

// IsExpired 判断是否已过期
func (f *Furniture) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// IsNew 判断是否全新
func (f *Furniture) IsNew() bool {
	return f.Condition == FurnitureConditionNew
}

// GetAge 获取家具使用年限（从购买日期计算）
func (f *Furniture) GetAge() int {
	if f.PurchaseDate == nil {
		return 0
	}
	years := time.Since(*f.PurchaseDate).Hours() / 24 / 365
	return int(years)
}

// SupportsDelivery 判断是否支持送货
func (f *Furniture) SupportsDelivery() bool {
	return f.DeliveryMethod == DeliveryMethodDelivery || f.DeliveryMethod == DeliveryMethodNegotiable
}

// SupportsSelfPickup 判断是否支持自取
func (f *Furniture) SupportsSelfPickup() bool {
	return f.DeliveryMethod == DeliveryMethodSelfPickup || f.DeliveryMethod == DeliveryMethodNegotiable
}

// GetDaysUntilExpiry 获取距离过期的天数
func (f *Furniture) GetDaysUntilExpiry() int {
	if f.IsExpired() {
		return 0
	}
	duration := time.Until(f.ExpiresAt)
	return int(duration.Hours() / 24)
}

// BeforeCreate GORM hook - 创建前执行
func (f *Furniture) BeforeCreate(tx *gorm.DB) error {
	if f.Status == "" {
		f.Status = FurnitureStatusAvailable
	}
	// 设置默认过期时间为90天后
	if f.ExpiresAt.IsZero() {
		f.ExpiresAt = time.Now().AddDate(0, 0, 90)
	}
	return nil
}

// FurnitureCategory 家具分类表模型
type FurnitureCategory struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID   *uint     `gorm:"index" json:"parent_id,omitempty"`
	NameZhHant string    `gorm:"type:varchar(100);not null" json:"name_zh_hant"`
	NameZhHans *string   `gorm:"type:varchar(100)" json:"name_zh_hans,omitempty"`
	NameEn     *string   `gorm:"type:varchar(100)" json:"name_en,omitempty"`
	Icon       *string   `gorm:"type:varchar(100)" json:"icon,omitempty"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Parent       *FurnitureCategory   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Subcategories []FurnitureCategory `gorm:"foreignKey:ParentID" json:"subcategories,omitempty"`
}

// TableName 指定表名
func (FurnitureCategory) TableName() string {
	return "furniture_categories"
}

// GetLocalizedName 根据语言获取本地化名称
func (fc *FurnitureCategory) GetLocalizedName(lang string) string {
	switch lang {
	case "zh-Hans", "zh_CN":
		if fc.NameZhHans != nil {
			return *fc.NameZhHans
		}
		return fc.NameZhHant
	case "en":
		if fc.NameEn != nil {
			return *fc.NameEn
		}
		return fc.NameZhHant
	default: // zh-Hant, zh_HK
		return fc.NameZhHant
	}
}

// IsTopLevel 判断是否为顶级分类
func (fc *FurnitureCategory) IsTopLevel() bool {
	return fc.ParentID == nil
}

// HasSubcategories 判断是否有子分类
func (fc *FurnitureCategory) HasSubcategories() bool {
	return len(fc.Subcategories) > 0
}

// IsActive 判断是否启用
func (fc *FurnitureCategory) IsActive() bool {
	return fc.IsActive
}

// FurnitureImage 家具图片表模型
type FurnitureImage struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FurnitureID uint      `gorm:"not null;index" json:"furniture_id"`
	ImageURL    string    `gorm:"type:varchar(500);not null" json:"image_url"`
	IsCover     bool      `gorm:"not null;default:false" json:"is_cover"`
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Furniture *Furniture `gorm:"foreignKey:FurnitureID;constraint:OnDelete:CASCADE" json:"furniture,omitempty"`
}

// TableName 指定表名
func (FurnitureImage) TableName() string {
	return "furniture_images"
}

// IsCover 判断是否为封面图
func (fi *FurnitureImage) IsCover() bool {
	return fi.IsCover
}

// FurnitureOrderStatus 家具订单状态常量
type FurnitureOrderStatus string

const (
	FurnitureOrderStatusPending   FurnitureOrderStatus = "pending"   // 待确认
	FurnitureOrderStatusConfirmed FurnitureOrderStatus = "confirmed" // 已确认
	FurnitureOrderStatusCompleted FurnitureOrderStatus = "completed" // 已完成
	FurnitureOrderStatusCancelled FurnitureOrderStatus = "cancelled" // 已取消
)

// FurnitureOrder 家具订单表模型
type FurnitureOrder struct {
	ID              uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo         string               `gorm:"type:varchar(50);not null;uniqueIndex" json:"order_no"`
	BuyerID         uint                 `gorm:"not null;index" json:"buyer_id"`
	SellerID        uint                 `gorm:"not null;index" json:"seller_id"`
	FurnitureID     uint                 `gorm:"not null;index" json:"furniture_id"`
	Price           float64              `gorm:"type:decimal(10,2);not null" json:"price"`
	Status          FurnitureOrderStatus `gorm:"type:varchar(20);not null;index" json:"status"`
	DeliveryMethod  string               `gorm:"type:varchar(50);not null" json:"delivery_method"`
	DeliveryAddress *string              `gorm:"type:varchar(500)" json:"delivery_address,omitempty"`
	DeliveryDate    *time.Time           `gorm:"type:date" json:"delivery_date,omitempty"`
	BuyerNote       *string              `gorm:"type:text" json:"buyer_note,omitempty"`
	SellerNote      *string              `gorm:"type:text" json:"seller_note,omitempty"`
	CreatedAt       time.Time            `gorm:"autoCreateTime;index" json:"created_at"`
	ConfirmedAt     *time.Time           `json:"confirmed_at,omitempty"`
	CompletedAt     *time.Time           `json:"completed_at,omitempty"`
	CancelledAt     *time.Time           `json:"cancelled_at,omitempty"`

	// 关联
	Buyer     *User      `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller    *User      `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Furniture *Furniture `gorm:"foreignKey:FurnitureID" json:"furniture,omitempty"`
}

// TableName 指定表名
func (FurnitureOrder) TableName() string {
	return "furniture_orders"
}

// IsPending 判断是否待确认
func (fo *FurnitureOrder) IsPending() bool {
	return fo.Status == FurnitureOrderStatusPending
}

// IsConfirmed 判断是否已确认
func (fo *FurnitureOrder) IsConfirmed() bool {
	return fo.Status == FurnitureOrderStatusConfirmed
}

// IsCompleted 判断是否已完成
func (fo *FurnitureOrder) IsCompleted() bool {
	return fo.Status == FurnitureOrderStatusCompleted
}

// IsCancelled 判断是否已取消
func (fo *FurnitureOrder) IsCancelled() bool {
	return fo.Status == FurnitureOrderStatusCancelled
}

// CanConfirm 判断是否可以确认（当前为待确认状态）
func (fo *FurnitureOrder) CanConfirm() bool {
	return fo.IsPending()
}

// CanComplete 判断是否可以完成（当前为已确认状态）
func (fo *FurnitureOrder) CanComplete() bool {
	return fo.IsConfirmed()
}

// CanCancel 判断是否可以取消（待确认或已确认状态）
func (fo *FurnitureOrder) CanCancel() bool {
	return fo.IsPending() || fo.IsConfirmed()
}
