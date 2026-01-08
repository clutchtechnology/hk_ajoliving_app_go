package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// ServicedApartment 服务式住宅模型
type ServicedApartment struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:200;not null;index" json:"name"`                    // 住宅名称
	NameEn       string         `gorm:"size:200" json:"name_en,omitempty"`                      // 英文名称
	Address      string         `gorm:"size:500;not null" json:"address"`                       // 详细地址
	DistrictID   uint           `gorm:"not null;index" json:"district_id"`                      // 所属地区ID
	Description  string         `gorm:"type:text" json:"description,omitempty"`                 // 详细描述
	Phone        string         `gorm:"size:50;not null" json:"phone"`                          // 联系电话
	WebsiteURL   string         `gorm:"size:500" json:"website_url,omitempty"`                  // 官方网站
	Email        string         `gorm:"size:255" json:"email,omitempty"`                        // 联系邮箱
	CompanyID    uint           `gorm:"not null;index" json:"company_id"`                       // 所属公司ID（关联users表，user_type='agency'）
	CheckInTime  string         `gorm:"size:50" json:"check_in_time,omitempty"`                 // 入住时间
	CheckOutTime string         `gorm:"size:50" json:"check_out_time,omitempty"`                // 退房时间
	MinStayDays  int            `json:"min_stay_days,omitempty"`                                // 最少入住天数
	Status       string         `gorm:"size:20;not null;default:'active';index" json:"status"`  // active=营业中, inactive=暂停营业, closed=已关闭
	Rating       float64        `gorm:"type:decimal(3,2)" json:"rating,omitempty"`              // 评分（0-5）
	ReviewCount  int            `gorm:"default:0" json:"review_count"`                          // 评价数量
	ViewCount    int            `gorm:"default:0" json:"view_count"`                            // 浏览次数
	FavoriteCount int           `gorm:"default:0" json:"favorite_count"`                        // 收藏次数
	IsFeatured   bool           `gorm:"default:false;index" json:"is_featured"`                 // 是否精选推荐
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`                                // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                                             // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                                         // 软删除时间

	// 关联
	District *District                    `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Company  *User                        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Images   []ServicedApartmentImage     `gorm:"foreignKey:ServicedApartmentID" json:"images,omitempty"`
	Units    []ServicedApartmentUnit      `gorm:"foreignKey:ServicedApartmentID" json:"units,omitempty"`
}

func (ServicedApartment) TableName() string {
	return "serviced_apartments"
}

// ServicedApartmentUnit 服务式住宅房型
type ServicedApartmentUnit struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	ServicedApartmentID  uint      `gorm:"not null;index" json:"serviced_apartment_id"`
	UnitType             string    `gorm:"size:50;not null" json:"unit_type"`          // 房型名称
	Bedrooms             int       `gorm:"not null" json:"bedrooms"`                   // 房间数
	Bathrooms            int       `json:"bathrooms,omitempty"`                        // 浴室数
	Area                 float64   `gorm:"not null" json:"area"`                       // 面积（平方尺）
	MaxOccupancy         int       `gorm:"not null" json:"max_occupancy"`              // 最多入住人数
	DailyPrice           float64   `json:"daily_price,omitempty"`                      // 日租价格（港币）
	WeeklyPrice          float64   `json:"weekly_price,omitempty"`                     // 周租价格（港币）
	MonthlyPrice         float64   `gorm:"not null;index" json:"monthly_price"`        // 月租价格（港币）
	AvailableUnits       int       `gorm:"not null" json:"available_units"`            // 可用单位数
	Description          string    `gorm:"type:text" json:"description,omitempty"`     // 房型描述
	SortOrder            int       `gorm:"default:0" json:"sort_order"`                // 排序顺序
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (ServicedApartmentUnit) TableName() string {
	return "serviced_apartment_units"
}

// ServicedApartmentImage 服务式住宅图片
type ServicedApartmentImage struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	ServicedApartmentID *uint     `gorm:"index" json:"serviced_apartment_id,omitempty"` // 关联的服务式住宅ID（整体照片）
	UnitID              *uint     `gorm:"index" json:"unit_id,omitempty"`               // 关联的房型ID（房型照片）
	ImageURL            string    `gorm:"size:500;not null" json:"image_url"`
	ImageType           string    `gorm:"size:20;not null" json:"image_type"` // exterior=外观, lobby=大堂, room=房间, bathroom=浴室, facilities=设施
	Title               string    `gorm:"size:200" json:"title,omitempty"`
	SortOrder           int       `gorm:"default:0" json:"sort_order"`
	CreatedAt           time.Time `json:"created_at"`
}

func (ServicedApartmentImage) TableName() string {
	return "serviced_apartment_images"
}

// ============ Request DTO ============

// ListServicedApartmentsRequest 获取服务式住宅列表请求
type ListServicedApartmentsRequest struct {
	DistrictID    *uint    `form:"district_id"`                                                     // 地区ID
	MinPrice      *float64 `form:"min_price" binding:"omitempty,gt=0"`                              // 最低月租
	MaxPrice      *float64 `form:"max_price" binding:"omitempty,gt=0"`                              // 最高月租
	MinRating     *float64 `form:"min_rating" binding:"omitempty,gte=0,lte=5"`                      // 最低评分
	Status        *string  `form:"status" binding:"omitempty,oneof=active inactive closed"`         // 状态
	IsFeatured    *bool    `form:"is_featured"`                                                     // 是否精选
	Page          int      `form:"page,default=1" binding:"min=1"`                                  // 页码
	PageSize      int      `form:"page_size,default=20" binding:"min=1,max=100"`                    // 每页数量
}

// CreateServicedApartmentRequest 创建服务式住宅请求
type CreateServicedApartmentRequest struct {
	Name         string `json:"name" binding:"required,max=200"`        // 住宅名称
	NameEn       string `json:"name_en" binding:"omitempty,max=200"`    // 英文名称
	Address      string `json:"address" binding:"required,max=500"`     // 详细地址
	DistrictID   uint   `json:"district_id" binding:"required"`         // 所属地区ID
	Description  string `json:"description" binding:"omitempty"`        // 详细描述
	Phone        string `json:"phone" binding:"required,max=50"`        // 联系电话
	WebsiteURL   string `json:"website_url" binding:"omitempty,url"`    // 官方网站
	Email        string `json:"email" binding:"omitempty,email"`        // 联系邮箱
	CheckInTime  string `json:"check_in_time" binding:"omitempty"`      // 入住时间
	CheckOutTime string `json:"check_out_time" binding:"omitempty"`     // 退房时间
	MinStayDays  int    `json:"min_stay_days" binding:"omitempty,min=1"`// 最少入住天数
}

// UpdateServicedApartmentRequest 更新服务式住宅请求
type UpdateServicedApartmentRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=200"`
	Description  *string `json:"description"`
	Phone        *string `json:"phone" binding:"omitempty,max=50"`
	WebsiteURL   *string `json:"website_url" binding:"omitempty,url"`
	Email        *string `json:"email" binding:"omitempty,email"`
	CheckInTime  *string `json:"check_in_time"`
	CheckOutTime *string `json:"check_out_time"`
	MinStayDays  *int    `json:"min_stay_days" binding:"omitempty,min=1"`
	Status       *string `json:"status" binding:"omitempty,oneof=active inactive closed"`
}

// ============ Response DTO ============

// ServicedApartmentResponse 服务式住宅响应（列表用）
type ServicedApartmentResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	NameEn        string    `json:"name_en,omitempty"`
	Address       string    `json:"address"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	Rating        float64   `json:"rating,omitempty"`
	ReviewCount   int       `json:"review_count"`
	ViewCount     int       `json:"view_count"`
	FavoriteCount int       `json:"favorite_count"`
	IsFeatured    bool      `json:"is_featured"`
	MinMonthlyPrice float64 `json:"min_monthly_price,omitempty"` // 最低月租
	CoverImage    string    `json:"cover_image,omitempty"`
	District      *District `json:"district,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ServicedApartmentDetailResponse 服务式住宅详情响应
type ServicedApartmentDetailResponse struct {
	ID            uint                        `json:"id"`
	Name          string                      `json:"name"`
	NameEn        string                      `json:"name_en,omitempty"`
	Address       string                      `json:"address"`
	Description   string                      `json:"description,omitempty"`
	Phone         string                      `json:"phone"`
	WebsiteURL    string                      `json:"website_url,omitempty"`
	Email         string                      `json:"email,omitempty"`
	CompanyID     uint                        `json:"company_id"`
	CheckInTime   string                      `json:"check_in_time,omitempty"`
	CheckOutTime  string                      `json:"check_out_time,omitempty"`
	MinStayDays   int                         `json:"min_stay_days,omitempty"`
	Status        string                      `json:"status"`
	Rating        float64                     `json:"rating,omitempty"`
	ReviewCount   int                         `json:"review_count"`
	ViewCount     int                         `json:"view_count"`
	FavoriteCount int                         `json:"favorite_count"`
	IsFeatured    bool                        `json:"is_featured"`
	District      *District                   `json:"district,omitempty"`
	Images        []ServicedApartmentImage    `json:"images,omitempty"`
	Units         []ServicedApartmentUnit     `json:"units,omitempty"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
}

// PaginatedServicedApartmentsResponse 分页服务式住宅响应
type PaginatedServicedApartmentsResponse struct {
	Data       []ServicedApartmentResponse `json:"data"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}

// ToServicedApartmentResponse 转换为服务式住宅响应
func (sa *ServicedApartment) ToServicedApartmentResponse() *ServicedApartmentResponse {
	resp := &ServicedApartmentResponse{
		ID:            sa.ID,
		Name:          sa.Name,
		NameEn:        sa.NameEn,
		Address:       sa.Address,
		Phone:         sa.Phone,
		Status:        sa.Status,
		Rating:        sa.Rating,
		ReviewCount:   sa.ReviewCount,
		ViewCount:     sa.ViewCount,
		FavoriteCount: sa.FavoriteCount,
		IsFeatured:    sa.IsFeatured,
		District:      sa.District,
		CreatedAt:     sa.CreatedAt,
	}

	// 获取最低月租
	if len(sa.Units) > 0 {
		minPrice := sa.Units[0].MonthlyPrice
		for _, unit := range sa.Units {
			if unit.MonthlyPrice < minPrice && unit.MonthlyPrice > 0 {
				minPrice = unit.MonthlyPrice
			}
		}
		resp.MinMonthlyPrice = minPrice
	}

	// 获取封面图
	if len(sa.Images) > 0 {
		for _, img := range sa.Images {
			if img.ImageType == "exterior" {
				resp.CoverImage = img.ImageURL
				break
			}
		}
		if resp.CoverImage == "" {
			resp.CoverImage = sa.Images[0].ImageURL
		}
	}

	return resp
}

// ToServicedApartmentDetailResponse 转换为服务式住宅详情响应
func (sa *ServicedApartment) ToServicedApartmentDetailResponse() *ServicedApartmentDetailResponse {
	return &ServicedApartmentDetailResponse{
		ID:            sa.ID,
		Name:          sa.Name,
		NameEn:        sa.NameEn,
		Address:       sa.Address,
		Description:   sa.Description,
		Phone:         sa.Phone,
		WebsiteURL:    sa.WebsiteURL,
		Email:         sa.Email,
		CompanyID:     sa.CompanyID,
		CheckInTime:   sa.CheckInTime,
		CheckOutTime:  sa.CheckOutTime,
		MinStayDays:   sa.MinStayDays,
		Status:        sa.Status,
		Rating:        sa.Rating,
		ReviewCount:   sa.ReviewCount,
		ViewCount:     sa.ViewCount,
		FavoriteCount: sa.FavoriteCount,
		IsFeatured:    sa.IsFeatured,
		District:      sa.District,
		Images:        sa.Images,
		Units:         sa.Units,
		CreatedAt:     sa.CreatedAt,
		UpdatedAt:     sa.UpdatedAt,
	}
}
