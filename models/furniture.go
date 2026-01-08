package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// Furniture 家具模型
type Furniture struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	FurnitureNo        string         `gorm:"size:50;uniqueIndex;not null" json:"furniture_no"`        // 家具编号（系统生成）
	Title              string         `gorm:"size:255;not null;index" json:"title"`                    // 家具名称/标题
	Description        string         `gorm:"type:text" json:"description,omitempty"`                  // 详细描述
	Price              float64        `gorm:"not null;index" json:"price"`                             // 价格（港币）
	CategoryID         uint           `gorm:"not null;index" json:"category_id"`                       // 分类ID
	Brand              string         `gorm:"size:100;index" json:"brand,omitempty"`                   // 品牌
	Condition          string         `gorm:"size:20;not null;index" json:"condition"`                 // new=全新, like_new=近全新, good=良好, fair=一般, poor=较差
	PurchaseDate       *time.Time     `json:"purchase_date,omitempty"`                                 // 购买日期
	DeliveryDistrictID uint           `gorm:"not null;index" json:"delivery_district_id"`              // 交收地区ID
	DeliveryTime       string         `gorm:"size:100" json:"delivery_time,omitempty"`                 // 交收时间
	DeliveryMethod     string         `gorm:"size:50;not null;index" json:"delivery_method"`           // self_pickup=自取, delivery=送货, negotiable=面议
	Status             string         `gorm:"size:20;not null;default:'available';index" json:"status"`// available=可用, reserved=已预订, sold=已售出, expired=已过期, cancelled=已取消
	PublisherID        uint           `gorm:"not null;index" json:"publisher_id"`                      // 发布者ID
	PublisherType      string         `gorm:"size:20;not null" json:"publisher_type"`                  // individual=个人, agency=代理公司
	ViewCount          int            `gorm:"default:0" json:"view_count"`                             // 浏览次数
	FavoriteCount      int            `gorm:"default:0" json:"favorite_count"`                         // 收藏次数
	PublishedAt        time.Time      `gorm:"not null;index" json:"published_at"`                      // 刊登日期
	UpdatedAt          time.Time      `gorm:"index" json:"updated_at"`                                 // 更新日期
	ExpiresAt          time.Time      `gorm:"not null;index" json:"expires_at"`                        // 到期日期
	CreatedAt          time.Time      `json:"created_at"`                                              // 创建时间
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`                                          // 软删除时间

	// 关联
	Category         *FurnitureCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	DeliveryDistrict *District          `gorm:"foreignKey:DeliveryDistrictID" json:"delivery_district,omitempty"`
	Images           []FurnitureImage   `gorm:"foreignKey:FurnitureID" json:"images,omitempty"`
}

func (Furniture) TableName() string {
	return "furniture"
}

// FurnitureCategory 家具分类
type FurnitureCategory struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ParentID   *uint          `gorm:"index" json:"parent_id,omitempty"`           // 父分类ID
	NameZhHant string         `gorm:"size:100;not null" json:"name_zh_hant"`      // 中文繁体名称
	NameZhHans string         `gorm:"size:100" json:"name_zh_hans,omitempty"`     // 中文简体名称
	NameEn     string         `gorm:"size:100" json:"name_en,omitempty"`          // 英文名称
	Icon       string         `gorm:"size:100" json:"icon,omitempty"`             // 图标标识
	SortOrder  int            `gorm:"not null" json:"sort_order"`                 // 排序顺序
	IsActive   bool           `gorm:"default:true" json:"is_active"`              // 是否启用
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Parent       *FurnitureCategory   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	SubCategories []FurnitureCategory `gorm:"foreignKey:ParentID" json:"sub_categories,omitempty"`
}

func (FurnitureCategory) TableName() string {
	return "furniture_categories"
}

// FurnitureImage 家具图片
type FurnitureImage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FurnitureID uint      `gorm:"not null;index" json:"furniture_id"`
	ImageURL    string    `gorm:"size:500;not null" json:"image_url"`
	IsCover     bool      `gorm:"default:false" json:"is_cover"` // 是否为封面图
	SortOrder   int       `gorm:"not null" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

func (FurnitureImage) TableName() string {
	return "furniture_images"
}

// ============ Request DTO ============

// ListFurnitureRequest 获取家具列表请求
type ListFurnitureRequest struct {
	CategoryID         *uint    `form:"category_id"`                                                 // 分类ID
	MinPrice           *float64 `form:"min_price"`                                                   // 最低价格
	MaxPrice           *float64 `form:"max_price"`                                                   // 最高价格
	Brand              *string  `form:"brand"`                                                       // 品牌
	Condition          *string  `form:"condition" binding:"omitempty,oneof=new like_new good fair poor"` // 新旧程度
	DeliveryDistrictID *uint    `form:"delivery_district_id"`                                        // 交收地区ID
	DeliveryMethod     *string  `form:"delivery_method" binding:"omitempty,oneof=self_pickup delivery negotiable"` // 交收方法
	Status             string   `form:"status" binding:"omitempty,oneof=available reserved sold expired cancelled"` // 状态
	Keyword            string   `form:"keyword"`                                                     // 搜索关键词
	SortBy             string   `form:"sort_by" binding:"omitempty,oneof=price published_at view_count"` // 排序字段
	SortOrder          string   `form:"sort_order" binding:"omitempty,oneof=asc desc"`               // 排序方向
	Page               int      `form:"page" binding:"min=1"`                                        // 页码
	PageSize           int      `form:"page_size" binding:"min=1,max=100"`                           // 每页数量
}

// CreateFurnitureRequest 发布家具请求
type CreateFurnitureRequest struct {
	Title              string     `json:"title" binding:"required,max=255"`
	Description        string     `json:"description"`
	Price              float64    `json:"price" binding:"required,gt=0"`
	CategoryID         uint       `json:"category_id" binding:"required"`
	Brand              string     `json:"brand" binding:"omitempty,max=100"`
	Condition          string     `json:"condition" binding:"required,oneof=new like_new good fair poor"`
	PurchaseDate       *time.Time `json:"purchase_date"`
	DeliveryDistrictID uint       `json:"delivery_district_id" binding:"required"`
	DeliveryTime       string     `json:"delivery_time" binding:"omitempty,max=100"`
	DeliveryMethod     string     `json:"delivery_method" binding:"required,oneof=self_pickup delivery negotiable"`
	ImageURLs          []string   `json:"image_urls" binding:"required,min=1"` // 至少1张图片
}

// UpdateFurnitureRequest 更新家具请求
type UpdateFurnitureRequest struct {
	Title              *string    `json:"title" binding:"omitempty,max=255"`
	Description        *string    `json:"description"`
	Price              *float64   `json:"price" binding:"omitempty,gt=0"`
	CategoryID         *uint      `json:"category_id"`
	Brand              *string    `json:"brand" binding:"omitempty,max=100"`
	Condition          *string    `json:"condition" binding:"omitempty,oneof=new like_new good fair poor"`
	PurchaseDate       *time.Time `json:"purchase_date"`
	DeliveryDistrictID *uint      `json:"delivery_district_id"`
	DeliveryTime       *string    `json:"delivery_time" binding:"omitempty,max=100"`
	DeliveryMethod     *string    `json:"delivery_method" binding:"omitempty,oneof=self_pickup delivery negotiable"`
	ImageURLs          []string   `json:"image_urls"` // 覆盖更新图片
}

// UpdateFurnitureStatusRequest 更新家具状态请求
type UpdateFurnitureStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=available reserved sold expired cancelled"`
}

// ============ Response DTO ============

// FurnitureResponse 家具响应
type FurnitureResponse struct {
	ID                 uint               `json:"id"`
	FurnitureNo        string             `json:"furniture_no"`
	Title              string             `json:"title"`
	Description        string             `json:"description,omitempty"`
	Price              float64            `json:"price"`
	Category           *FurnitureCategory `json:"category,omitempty"`
	Brand              string             `json:"brand,omitempty"`
	Condition          string             `json:"condition"`
	PurchaseDate       *time.Time         `json:"purchase_date,omitempty"`
	DeliveryDistrictID uint               `json:"delivery_district_id"`
	DeliveryDistrict   *District          `json:"delivery_district,omitempty"`
	DeliveryTime       string             `json:"delivery_time,omitempty"`
	DeliveryMethod     string             `json:"delivery_method"`
	Status             string             `json:"status"`
	PublisherID        uint               `json:"publisher_id"`
	PublisherType      string             `json:"publisher_type"`
	ViewCount          int                `json:"view_count"`
	FavoriteCount      int                `json:"favorite_count"`
	Images             []FurnitureImage   `json:"images,omitempty"`
	PublishedAt        time.Time          `json:"published_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	ExpiresAt          time.Time          `json:"expires_at"`
}

// PaginatedFurnitureResponse 分页家具列表响应
type PaginatedFurnitureResponse struct {
	Data       []FurnitureResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// FurnitureCategoryResponse 家具分类响应
type FurnitureCategoryResponse struct {
	ID            uint                        `json:"id"`
	ParentID      *uint                       `json:"parent_id,omitempty"`
	NameZhHant    string                      `json:"name_zh_hant"`
	NameZhHans    string                      `json:"name_zh_hans,omitempty"`
	NameEn        string                      `json:"name_en,omitempty"`
	Icon          string                      `json:"icon,omitempty"`
	SortOrder     int                         `json:"sort_order"`
	IsActive      bool                        `json:"is_active"`
	SubCategories []FurnitureCategoryResponse `json:"sub_categories,omitempty"`
	FurnitureCount int                        `json:"furniture_count"` // 该分类下的家具数量
}
