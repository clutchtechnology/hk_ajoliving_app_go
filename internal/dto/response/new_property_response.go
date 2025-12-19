package response

import "time"

// NewDevelopmentResponse 新楼盘详情响应
type NewDevelopmentResponse struct {
	ID                 uint                       `json:"id"`
	Name               string                     `json:"name"`
	NameEn             string                     `json:"name_en,omitempty"`
	Address            string                     `json:"address"`
	DistrictID         uint                       `json:"district_id"`
	District           *DistrictResponse          `json:"district,omitempty"`
	Status             string                     `json:"status"`
	UnitsForSale       int                        `json:"units_for_sale"`
	UnitsSold          int                        `json:"units_sold"`
	Developer          string                     `json:"developer"`
	ManagementCompany  string                     `json:"management_company,omitempty"`
	TotalUnits         int                        `json:"total_units"`
	TotalBlocks        int                        `json:"total_blocks"`
	MaxFloors          int                        `json:"max_floors"`
	PrimarySchoolNet   string                     `json:"primary_school_net,omitempty"`
	SecondarySchoolNet string                     `json:"secondary_school_net,omitempty"`
	WebsiteURL         string                     `json:"website_url,omitempty"`
	SalesOfficeAddress string                     `json:"sales_office_address,omitempty"`
	SalesPhone         string                     `json:"sales_phone,omitempty"`
	ExpectedCompletion *time.Time                 `json:"expected_completion,omitempty"`
	OccupationDate     *time.Time                 `json:"occupation_date,omitempty"`
	Description        string                     `json:"description,omitempty"`
	ViewCount          int                        `json:"view_count"`
	FavoriteCount      int                        `json:"favorite_count"`
	IsFeatured         bool                       `json:"is_featured"`
	SalesProgress      float64                    `json:"sales_progress"`
	Images             []NewDevelopmentImageResponse  `json:"images,omitempty"`
	Layouts            []NewDevelopmentLayoutResponse `json:"layouts,omitempty"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
}

// NewDevelopmentListItemResponse 新楼盘列表项响应
type NewDevelopmentListItemResponse struct {
	ID                 uint              `json:"id"`
	Name               string            `json:"name"`
	NameEn             string            `json:"name_en,omitempty"`
	Address            string            `json:"address"`
	DistrictID         uint              `json:"district_id"`
	District           *DistrictResponse `json:"district,omitempty"`
	Status             string            `json:"status"`
	Developer          string            `json:"developer"`
	TotalUnits         int               `json:"total_units"`
	UnitsForSale       int               `json:"units_for_sale"`
	MinPrice           float64           `json:"min_price"`
	MaxPrice           float64           `json:"max_price"`
	ExpectedCompletion *time.Time        `json:"expected_completion,omitempty"`
	CoverImage         string            `json:"cover_image,omitempty"`
	ViewCount          int               `json:"view_count"`
	FavoriteCount      int               `json:"favorite_count"`
	IsFeatured         bool              `json:"is_featured"`
	SalesProgress      float64           `json:"sales_progress"`
	CreatedAt          time.Time         `json:"created_at"`
}

// NewDevelopmentImageResponse 新楼盘图片响应
type NewDevelopmentImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	ImageType string `json:"image_type"`
	Title     string `json:"title,omitempty"`
	SortOrder int    `json:"sort_order"`
}

// NewDevelopmentLayoutResponse 新楼盘户型响应
type NewDevelopmentLayoutResponse struct {
	ID             uint    `json:"id"`
	UnitType       string  `json:"unit_type"`
	Bedrooms       int     `json:"bedrooms"`
	Bathrooms      int     `json:"bathrooms,omitempty"`
	SaleableArea   float64 `json:"saleable_area"`
	GrossArea      float64 `json:"gross_area,omitempty"`
	MinPrice       float64 `json:"min_price"`
	MaxPrice       float64 `json:"max_price,omitempty"`
	PricePerSqft   float64 `json:"price_per_sqft,omitempty"`
	AvailableUnits int     `json:"available_units"`
	FloorplanURL   string  `json:"floorplan_url,omitempty"`
}

// CreateNewDevelopmentResponse 创建新楼盘响应
type CreateNewDevelopmentResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

// UpdateNewDevelopmentResponse 更新新楼盘响应
type UpdateNewDevelopmentResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
