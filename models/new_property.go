package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// NewProperty 新盘模型
type NewProperty struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"size:200;not null;index" json:"name"`                          // 新盘名称
	NameEn              string         `gorm:"size:200" json:"name_en,omitempty"`                            // 英文名称
	Address             string         `gorm:"size:500;not null" json:"address"`                             // 详细地址
	DistrictID          uint           `gorm:"not null;index" json:"district_id"`                            // 所属地区ID
	Status              string         `gorm:"size:20;not null;index" json:"status"`                         // upcoming=即将推出, presale=预售中, selling=销售中, completed=已完成
	UnitsForSale        int            `json:"units_for_sale,omitempty"`                                     // 在售单位数
	UnitsSold           int            `json:"units_sold,omitempty"`                                         // 已售单位数
	Developer           string         `gorm:"size:200;not null;index" json:"developer"`                     // 开发商名称
	ManagementCompany   string         `gorm:"size:200" json:"management_company,omitempty"`                 // 管理公司名称
	TotalUnits          int            `gorm:"not null" json:"total_units"`                                  // 物业总伙数
	TotalBlocks         int            `gorm:"not null" json:"total_blocks"`                                 // 座数
	MaxFloors           int            `gorm:"not null" json:"max_floors"`                                   // 最高层数
	PrimarySchoolNet    string         `gorm:"size:50;index" json:"primary_school_net,omitempty"`            // 小学校网
	SecondarySchoolNet  string         `gorm:"size:50;index" json:"secondary_school_net,omitempty"`          // 中学校网
	WebsiteURL          string         `gorm:"size:500" json:"website_url,omitempty"`                        // 官方网页地址
	SalesOfficeAddress  string         `gorm:"size:500" json:"sales_office_address,omitempty"`               // 销售处地址
	SalesPhone          string         `gorm:"size:50" json:"sales_phone,omitempty"`                         // 销售电话
	ExpectedCompletion  *time.Time     `json:"expected_completion,omitempty"`                                // 预计落成日期
	OccupationDate      *time.Time     `json:"occupation_date,omitempty"`                                    // 入伙日期
	Description         string         `gorm:"type:text" json:"description,omitempty"`                       // 项目描述
	ViewCount           int            `gorm:"default:0" json:"view_count"`                                  // 浏览次数
	FavoriteCount       int            `gorm:"default:0" json:"favorite_count"`                              // 收藏次数
	SortOrder           int            `gorm:"default:0" json:"sort_order"`                                  // 排序顺序
	IsFeatured          bool           `gorm:"default:false;index" json:"is_featured"`                       // 是否精选推荐
	CreatedAt           time.Time      `gorm:"index" json:"created_at"`                                      // 创建时间
	UpdatedAt           time.Time      `json:"updated_at"`                                                   // 更新时间
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`                                               // 软删除时间

	// 关联
	District *District                 `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Images   []NewPropertyImage        `gorm:"foreignKey:NewPropertyID" json:"images,omitempty"`
	Layouts  []NewPropertyLayout       `gorm:"foreignKey:NewPropertyID" json:"layouts,omitempty"`
}

func (NewProperty) TableName() string {
	return "new_properties"
}

// NewPropertyImage 新盘图片
type NewPropertyImage struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NewPropertyID uint      `gorm:"not null;index" json:"new_property_id"`
	ImageURL      string    `gorm:"size:500;not null" json:"image_url"`
	ImageType     string    `gorm:"size:20;not null" json:"image_type"` // exterior=外观, interior=室内示范单位, facilities=设施, floorplan=户型图, location=位置图
	Title         string    `gorm:"size:200" json:"title,omitempty"`
	SortOrder     int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
}

func (NewPropertyImage) TableName() string {
	return "new_property_images"
}

// NewPropertyLayout 新盘户型
type NewPropertyLayout struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NewPropertyID uint      `gorm:"not null;index" json:"new_property_id"`
	UnitType      string    `gorm:"size:50;not null" json:"unit_type"`          // 户型类型（如：1房、2房、3房）
	Bedrooms      int       `gorm:"not null" json:"bedrooms"`                   // 房间数
	Bathrooms     int       `json:"bathrooms,omitempty"`                        // 浴室数
	SaleableArea  float64   `gorm:"not null" json:"saleable_area"`              // 实用面积（平方尺）
	GrossArea     float64   `json:"gross_area,omitempty"`                       // 建筑面积（平方尺）
	MinPrice      float64   `gorm:"not null" json:"min_price"`                  // 最低售价（港币）
	MaxPrice      float64   `json:"max_price,omitempty"`                        // 最高售价（港币）
	PricePerSqft  float64   `json:"price_per_sqft,omitempty"`                   // 每平方尺价格
	AvailableUnits int      `gorm:"not null" json:"available_units"`            // 可售单位数
	FloorplanURL  string    `gorm:"size:500" json:"floorplan_url,omitempty"`    // 户型图URL
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (NewPropertyLayout) TableName() string {
	return "new_property_layouts"
}

// ============ Request DTO ============

// ListNewPropertiesRequest 获取新盘列表请求
type ListNewPropertiesRequest struct {
	DistrictID         *uint    `form:"district_id"`                                                            // 地区ID
	Status             *string  `form:"status" binding:"omitempty,oneof=upcoming presale selling completed"`   // 状态
	Developer          *string  `form:"developer"`                                                              // 开发商
	MinPrice           *float64 `form:"min_price" binding:"omitempty,gt=0"`                                     // 最低价格
	MaxPrice           *float64 `form:"max_price" binding:"omitempty,gt=0"`                                     // 最高价格
	PrimarySchoolNet   *string  `form:"primary_school_net"`                                                     // 小学校网
	SecondarySchoolNet *string  `form:"secondary_school_net"`                                                   // 中学校网
	IsFeatured         *bool    `form:"is_featured"`                                                            // 是否精选
	Page               int      `form:"page,default=1" binding:"min=1"`                                         // 页码
	PageSize           int      `form:"page_size,default=20" binding:"min=1,max=100"`                           // 每页数量
}

// ============ Response DTO ============

// NewPropertyResponse 新盘响应（列表用）
type NewPropertyResponse struct {
	ID                 uint       `json:"id"`
	Name               string     `json:"name"`
	NameEn             string     `json:"name_en,omitempty"`
	Address            string     `json:"address"`
	Status             string     `json:"status"`
	UnitsForSale       int        `json:"units_for_sale,omitempty"`
	UnitsSold          int        `json:"units_sold,omitempty"`
	Developer          string     `json:"developer"`
	TotalUnits         int        `json:"total_units"`
	ExpectedCompletion *time.Time `json:"expected_completion,omitempty"`
	ViewCount          int        `json:"view_count"`
	FavoriteCount      int        `json:"favorite_count"`
	IsFeatured         bool       `json:"is_featured"`
	CoverImage         string     `json:"cover_image,omitempty"`
	District           *District  `json:"district,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// NewPropertyDetailResponse 新盘详情响应
type NewPropertyDetailResponse struct {
	ID                  uint                  `json:"id"`
	Name                string                `json:"name"`
	NameEn              string                `json:"name_en,omitempty"`
	Address             string                `json:"address"`
	Status              string                `json:"status"`
	UnitsForSale        int                   `json:"units_for_sale,omitempty"`
	UnitsSold           int                   `json:"units_sold,omitempty"`
	Developer           string                `json:"developer"`
	ManagementCompany   string                `json:"management_company,omitempty"`
	TotalUnits          int                   `json:"total_units"`
	TotalBlocks         int                   `json:"total_blocks"`
	MaxFloors           int                   `json:"max_floors"`
	PrimarySchoolNet    string                `json:"primary_school_net,omitempty"`
	SecondarySchoolNet  string                `json:"secondary_school_net,omitempty"`
	WebsiteURL          string                `json:"website_url,omitempty"`
	SalesOfficeAddress  string                `json:"sales_office_address,omitempty"`
	SalesPhone          string                `json:"sales_phone,omitempty"`
	ExpectedCompletion  *time.Time            `json:"expected_completion,omitempty"`
	OccupationDate      *time.Time            `json:"occupation_date,omitempty"`
	Description         string                `json:"description,omitempty"`
	ViewCount           int                   `json:"view_count"`
	FavoriteCount       int                   `json:"favorite_count"`
	IsFeatured          bool                  `json:"is_featured"`
	District            *District             `json:"district,omitempty"`
	Images              []NewPropertyImage    `json:"images,omitempty"`
	Layouts             []NewPropertyLayout   `json:"layouts,omitempty"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

// PaginatedNewPropertiesResponse 分页新盘响应
type PaginatedNewPropertiesResponse struct {
	Data       []NewPropertyResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// ToNewPropertyResponse 转换为新盘响应
func (np *NewProperty) ToNewPropertyResponse() *NewPropertyResponse {
	resp := &NewPropertyResponse{
		ID:                 np.ID,
		Name:               np.Name,
		NameEn:             np.NameEn,
		Address:            np.Address,
		Status:             np.Status,
		UnitsForSale:       np.UnitsForSale,
		UnitsSold:          np.UnitsSold,
		Developer:          np.Developer,
		TotalUnits:         np.TotalUnits,
		ExpectedCompletion: np.ExpectedCompletion,
		ViewCount:          np.ViewCount,
		FavoriteCount:      np.FavoriteCount,
		IsFeatured:         np.IsFeatured,
		District:           np.District,
		CreatedAt:          np.CreatedAt,
	}

	// 获取封面图
	if len(np.Images) > 0 {
		for _, img := range np.Images {
			if img.ImageType == "exterior" {
				resp.CoverImage = img.ImageURL
				break
			}
		}
		if resp.CoverImage == "" {
			resp.CoverImage = np.Images[0].ImageURL
		}
	}

	return resp
}

// ToNewPropertyDetailResponse 转换为新盘详情响应
func (np *NewProperty) ToNewPropertyDetailResponse() *NewPropertyDetailResponse {
	return &NewPropertyDetailResponse{
		ID:                  np.ID,
		Name:                np.Name,
		NameEn:              np.NameEn,
		Address:             np.Address,
		Status:              np.Status,
		UnitsForSale:        np.UnitsForSale,
		UnitsSold:           np.UnitsSold,
		Developer:           np.Developer,
		ManagementCompany:   np.ManagementCompany,
		TotalUnits:          np.TotalUnits,
		TotalBlocks:         np.TotalBlocks,
		MaxFloors:           np.MaxFloors,
		PrimarySchoolNet:    np.PrimarySchoolNet,
		SecondarySchoolNet:  np.SecondarySchoolNet,
		WebsiteURL:          np.WebsiteURL,
		SalesOfficeAddress:  np.SalesOfficeAddress,
		SalesPhone:          np.SalesPhone,
		ExpectedCompletion:  np.ExpectedCompletion,
		OccupationDate:      np.OccupationDate,
		Description:         np.Description,
		ViewCount:           np.ViewCount,
		FavoriteCount:       np.FavoriteCount,
		IsFeatured:          np.IsFeatured,
		District:            np.District,
		Images:              np.Images,
		Layouts:             np.Layouts,
		CreatedAt:           np.CreatedAt,
		UpdatedAt:           np.UpdatedAt,
	}
}
