package request

// ListNewDevelopmentsRequest 新楼盘列表请求
type ListNewDevelopmentsRequest struct {
	DistrictID *uint   `form:"district_id"`
	Developer  string  `form:"developer"`
	Status     string  `form:"status" binding:"omitempty,oneof=upcoming presale selling completed"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Bedrooms   *int     `form:"bedrooms"`
	SchoolNet  string   `form:"school_net"`
	IsFeatured *bool    `form:"is_featured"`
	SortBy     string   `form:"sort_by,default=created_at"`
	SortOrder  string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page       int      `form:"page,default=1" binding:"min=1"`
	PageSize   int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// CreateNewDevelopmentRequest 创建新楼盘请求
type CreateNewDevelopmentRequest struct {
	Name               string  `json:"name" binding:"required,max=200"`
	NameEn             string  `json:"name_en" binding:"omitempty,max=200"`
	Address            string  `json:"address" binding:"required,max=500"`
	DistrictID         uint    `json:"district_id" binding:"required"`
	Developer          string  `json:"developer" binding:"required,max=200"`
	ManagementCompany  string  `json:"management_company" binding:"omitempty,max=200"`
	TotalUnits         int     `json:"total_units" binding:"required,min=1"`
	TotalBlocks        int     `json:"total_blocks" binding:"required,min=1"`
	MaxFloors          int     `json:"max_floors" binding:"required,min=1"`
	PrimarySchoolNet   string  `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string  `json:"secondary_school_net" binding:"omitempty,max=50"`
	WebsiteURL         string  `json:"website_url" binding:"omitempty,url,max=500"`
	SalesOfficeAddress string  `json:"sales_office_address" binding:"omitempty,max=500"`
	SalesPhone         string  `json:"sales_phone" binding:"omitempty,max=50"`
	ExpectedCompletion string  `json:"expected_completion"` // 格式: YYYY-MM-DD
	Description        string  `json:"description"`
	Status             string  `json:"status" binding:"required,oneof=upcoming presale selling completed"`
}

// UpdateNewDevelopmentRequest 更新新楼盘请求
type UpdateNewDevelopmentRequest struct {
	Name               string  `json:"name" binding:"omitempty,max=200"`
	NameEn             string  `json:"name_en" binding:"omitempty,max=200"`
	Address            string  `json:"address" binding:"omitempty,max=500"`
	DistrictID         uint    `json:"district_id"`
	Developer          string  `json:"developer" binding:"omitempty,max=200"`
	ManagementCompany  string  `json:"management_company" binding:"omitempty,max=200"`
	TotalUnits         int     `json:"total_units" binding:"omitempty,min=1"`
	TotalBlocks        int     `json:"total_blocks" binding:"omitempty,min=1"`
	MaxFloors          int     `json:"max_floors" binding:"omitempty,min=1"`
	PrimarySchoolNet   string  `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string  `json:"secondary_school_net" binding:"omitempty,max=50"`
	WebsiteURL         string  `json:"website_url" binding:"omitempty,url,max=500"`
	SalesOfficeAddress string  `json:"sales_office_address" binding:"omitempty,max=500"`
	SalesPhone         string  `json:"sales_phone" binding:"omitempty,max=50"`
	ExpectedCompletion string  `json:"expected_completion"`
	Description        string  `json:"description"`
	Status             string  `json:"status" binding:"omitempty,oneof=upcoming presale selling completed"`
	IsFeatured         *bool   `json:"is_featured"`
}
