package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// Property 房产模型
type Property struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PropertyNo     string         `gorm:"size:50;uniqueIndex;not null" json:"property_no"`              // 物业编号（系统生成）
	EstateNo       string         `gorm:"size:50;index" json:"estate_no,omitempty"`                     // 楼盘编号（外部编号）
	ListingType    string         `gorm:"size:20;not null;index" json:"listing_type"`                   // sale=出售, rent=出租
	Title          string         `gorm:"size:255;not null" json:"title"`                               // 房产标题
	Description    string         `gorm:"type:text" json:"description,omitempty"`                       // 房产描述
	Area           float64        `gorm:"not null" json:"area"`                                         // 面积（平方尺）
	Price          float64        `gorm:"not null;index" json:"price"`                                  // 价格（港币）
	Address        string         `gorm:"size:500;not null" json:"address"`                             // 详细地址
	DistrictID     uint           `gorm:"not null;index" json:"district_id"`                            // 所属地区ID
	BuildingName   string         `gorm:"size:200;index" json:"building_name,omitempty"`                // 大厦/楼宇名称
	Floor          string         `gorm:"size:20" json:"floor,omitempty"`                               // 楼层
	Orientation    string         `gorm:"size:50" json:"orientation,omitempty"`                         // 座向
	Bedrooms       int            `gorm:"not null;index" json:"bedrooms"`                               // 房间数
	Bathrooms      int            `json:"bathrooms,omitempty"`                                          // 浴室数
	PrimarySchool  string         `gorm:"size:50;index" json:"primary_school_net,omitempty"`            // 小学校网
	SecondarySchool string        `gorm:"size:50;index" json:"secondary_school_net,omitempty"`          // 中学校网
	PropertyType   string         `gorm:"size:50;not null;index" json:"property_type"`                  // 物业类型
	Status         string         `gorm:"size:20;not null;default:'available';index" json:"status"`     // 状态
	PublisherID    uint           `gorm:"not null;index" json:"publisher_id"`                           // 发布者ID
	PublisherType  string         `gorm:"size:20;not null" json:"publisher_type"`                       // individual=个人, agency=代理公司
	AgentID        *uint          `gorm:"index" json:"agent_id,omitempty"`                              // 负责地产代理ID
	ViewCount      int            `gorm:"default:0" json:"view_count"`                                  // 浏览次数
	FavoriteCount  int            `gorm:"default:0" json:"favorite_count"`                              // 收藏次数
	PublishedAt    *time.Time     `gorm:"index" json:"published_at,omitempty"`                          // 发布时间
	ExpiredAt      *time.Time     `gorm:"index" json:"expired_at,omitempty"`                            // 过期时间
	CreatedAt      time.Time      `gorm:"index" json:"created_at"`                                      // 创建时间
	UpdatedAt      time.Time      `json:"updated_at"`                                                   // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                                               // 软删除时间

	// 关联
	District *District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Images   []PropertyImage `gorm:"foreignKey:PropertyID" json:"images,omitempty"`
}

func (Property) TableName() string {
	return "properties"
}

// PropertyImage 房产图片
type PropertyImage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID uint      `gorm:"not null;index" json:"property_id"`
	ImageURL   string    `gorm:"size:500;not null" json:"image_url"`
	ImageType  string    `gorm:"size:20;not null" json:"image_type"` // cover=封面, interior=室内, exterior=外观, floorplan=户型图
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PropertyImage) TableName() string {
	return "property_images"
}

// ============ Request DTO ============

// ListPropertiesRequest 获取房产列表请求
type ListPropertiesRequest struct {
	ListingType      *string  `form:"listing_type" binding:"omitempty,oneof=sale rent"`      // 房源类型
	DistrictID       *uint    `form:"district_id"`                                           // 地区ID
	MinPrice         *float64 `form:"min_price" binding:"omitempty,gt=0"`                    // 最低价格
	MaxPrice         *float64 `form:"max_price" binding:"omitempty,gt=0"`                    // 最高价格
	MinArea          *float64 `form:"min_area" binding:"omitempty,gt=0"`                     // 最小面积
	MaxArea          *float64 `form:"max_area" binding:"omitempty,gt=0"`                     // 最大面积
	Bedrooms         *int     `form:"bedrooms" binding:"omitempty,min=0"`                    // 房间数
	PropertyType     *string  `form:"property_type"`                                         // 物业类型
	BuildingName     *string  `form:"building_name"`                                         // 大厦名称
	PrimarySchool    *string  `form:"primary_school_net"`                                    // 小学校网
	SecondarySchool  *string  `form:"secondary_school_net"`                                  // 中学校网
	Status           *string  `form:"status" binding:"omitempty,oneof=available pending sold cancelled"` // 状态
	SortBy           string   `form:"sort_by" binding:"omitempty,oneof=price_asc price_desc area_asc area_desc created_at_desc"` // 排序方式
	Page             int      `form:"page,default=1" binding:"min=1"`                        // 页码
	PageSize         int      `form:"page_size,default=20" binding:"min=1,max=100"`          // 每页数量
}

// CreatePropertyRequest 创建房产请求
type CreatePropertyRequest struct {
	EstateNo        string   `json:"estate_no" binding:"omitempty,max=50"`                           // 楼盘编号
	ListingType     string   `json:"listing_type" binding:"required,oneof=sale rent"`                // sale=出售, rent=出租
	Title           string   `json:"title" binding:"required,max=255"`                               // 房产标题
	Description     string   `json:"description" binding:"omitempty"`                                // 房产描述
	Area            float64  `json:"area" binding:"required,gt=0"`                                   // 面积（平方尺）
	Price           float64  `json:"price" binding:"required,gt=0"`                                  // 价格（港币）
	Address         string   `json:"address" binding:"required,max=500"`                             // 详细地址
	DistrictID      uint     `json:"district_id" binding:"required"`                                 // 所属地区ID
	BuildingName    string   `json:"building_name" binding:"omitempty,max=200"`                      // 大厦/楼宇名称
	Floor           string   `json:"floor" binding:"omitempty,max=20"`                               // 楼层
	Orientation     string   `json:"orientation" binding:"omitempty,max=50"`                         // 座向
	Bedrooms        int      `json:"bedrooms" binding:"required,min=0"`                              // 房间数
	Bathrooms       int      `json:"bathrooms" binding:"omitempty,min=0"`                            // 浴室数
	PrimarySchool   string   `json:"primary_school_net" binding:"omitempty,max=50"`                  // 小学校网
	SecondarySchool string   `json:"secondary_school_net" binding:"omitempty,max=50"`                // 中学校网
	PropertyType    string   `json:"property_type" binding:"required,max=50"`                        // 物业类型
	AgentID         *uint    `json:"agent_id" binding:"omitempty"`                                   // 负责地产代理ID
	ImageURLs       []string `json:"image_urls" binding:"omitempty,max=20,dive,url"`                 // 图片URL列表
}

// UpdatePropertyRequest 更新房产请求
type UpdatePropertyRequest struct {
	Title           *string  `json:"title" binding:"omitempty,max=255"`
	Description     *string  `json:"description"`
	Price           *float64 `json:"price" binding:"omitempty,gt=0"`
	Area            *float64 `json:"area" binding:"omitempty,gt=0"`
	Floor           *string  `json:"floor" binding:"omitempty,max=20"`
	Orientation     *string  `json:"orientation" binding:"omitempty,max=50"`
	Bathrooms       *int     `json:"bathrooms" binding:"omitempty,min=0"`
	Status          *string  `json:"status" binding:"omitempty,oneof=available pending sold cancelled"`
	AgentID         *uint    `json:"agent_id"`
}

// ============ Response DTO ============

// PropertyResponse 房产响应（列表用）
type PropertyResponse struct {
	ID            uint       `json:"id"`
	PropertyNo    string     `json:"property_no"`
	ListingType   string     `json:"listing_type"`
	Title         string     `json:"title"`
	Price         float64    `json:"price"`
	Area          float64    `json:"area"`
	Address       string     `json:"address"`
	BuildingName  string     `json:"building_name,omitempty"`
	Bedrooms      int        `json:"bedrooms"`
	Bathrooms     int        `json:"bathrooms,omitempty"`
	PropertyType  string     `json:"property_type"`
	Status        string     `json:"status"`
	ViewCount     int        `json:"view_count"`
	FavoriteCount int        `json:"favorite_count"`
	CoverImage    string     `json:"cover_image,omitempty"`
	District      *District  `json:"district,omitempty"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// PropertyDetailResponse 房产详情响应
type PropertyDetailResponse struct {
	ID              uint            `json:"id"`
	PropertyNo      string          `json:"property_no"`
	EstateNo        string          `json:"estate_no,omitempty"`
	ListingType     string          `json:"listing_type"`
	Title           string          `json:"title"`
	Description     string          `json:"description,omitempty"`
	Price           float64         `json:"price"`
	Area            float64         `json:"area"`
	Address         string          `json:"address"`
	BuildingName    string          `json:"building_name,omitempty"`
	Floor           string          `json:"floor,omitempty"`
	Orientation     string          `json:"orientation,omitempty"`
	Bedrooms        int             `json:"bedrooms"`
	Bathrooms       int             `json:"bathrooms,omitempty"`
	PrimarySchool   string          `json:"primary_school_net,omitempty"`
	SecondarySchool string          `json:"secondary_school_net,omitempty"`
	PropertyType    string          `json:"property_type"`
	Status          string          `json:"status"`
	PublisherID     uint            `json:"publisher_id"`
	PublisherType   string          `json:"publisher_type"`
	AgentID         *uint           `json:"agent_id,omitempty"`
	ViewCount       int             `json:"view_count"`
	FavoriteCount   int             `json:"favorite_count"`
	District        *District       `json:"district,omitempty"`
	Images          []PropertyImage `json:"images,omitempty"`
	PublishedAt     *time.Time      `json:"published_at,omitempty"`
	ExpiredAt       *time.Time      `json:"expired_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// PaginatedPropertiesResponse 分页房产响应
type PaginatedPropertiesResponse struct {
	Data       []PropertyResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// ToPropertyResponse 转换为房产响应
func (p *Property) ToPropertyResponse() *PropertyResponse {
	resp := &PropertyResponse{
		ID:            p.ID,
		PropertyNo:    p.PropertyNo,
		ListingType:   p.ListingType,
		Title:         p.Title,
		Price:         p.Price,
		Area:          p.Area,
		Address:       p.Address,
		BuildingName:  p.BuildingName,
		Bedrooms:      p.Bedrooms,
		Bathrooms:     p.Bathrooms,
		PropertyType:  p.PropertyType,
		Status:        p.Status,
		ViewCount:     p.ViewCount,
		FavoriteCount: p.FavoriteCount,
		District:      p.District,
		PublishedAt:   p.PublishedAt,
		CreatedAt:     p.CreatedAt,
	}

	// 获取封面图
	if len(p.Images) > 0 {
		for _, img := range p.Images {
			if img.ImageType == "cover" {
				resp.CoverImage = img.ImageURL
				break
			}
		}
		// 如果没有封面图，使用第一张
		if resp.CoverImage == "" {
			resp.CoverImage = p.Images[0].ImageURL
		}
	}

	return resp
}

// ToPropertyDetailResponse 转换为房产详情响应
func (p *Property) ToPropertyDetailResponse() *PropertyDetailResponse {
	return &PropertyDetailResponse{
		ID:              p.ID,
		PropertyNo:      p.PropertyNo,
		EstateNo:        p.EstateNo,
		ListingType:     p.ListingType,
		Title:           p.Title,
		Description:     p.Description,
		Price:           p.Price,
		Area:            p.Area,
		Address:         p.Address,
		BuildingName:    p.BuildingName,
		Floor:           p.Floor,
		Orientation:     p.Orientation,
		Bedrooms:        p.Bedrooms,
		Bathrooms:       p.Bathrooms,
		PrimarySchool:   p.PrimarySchool,
		SecondarySchool: p.SecondarySchool,
		PropertyType:    p.PropertyType,
		Status:          p.Status,
		PublisherID:     p.PublisherID,
		PublisherType:   p.PublisherType,
		AgentID:         p.AgentID,
		ViewCount:       p.ViewCount,
		FavoriteCount:   p.FavoriteCount,
		District:        p.District,
		Images:          p.Images,
		PublishedAt:     p.PublishedAt,
		ExpiredAt:       p.ExpiredAt,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}
