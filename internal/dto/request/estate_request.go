package request

// ListEstatesRequest 屋苑列表请求
type ListEstatesRequest struct {
	DistrictID         *uint    `form:"district_id"`
	SchoolNet          string   `form:"school_net"`
	MinCompletionYear  *int     `form:"min_completion_year"`
	MaxCompletionYear  *int     `form:"max_completion_year"`
	MinAvgPrice        *float64 `form:"min_avg_price"`
	MaxAvgPrice        *float64 `form:"max_avg_price"`
	HasListings        *bool    `form:"has_listings"`        // 是否有放盘/租盘
	HasTransactions    *bool    `form:"has_transactions"`    // 是否有近期成交
	IsFeatured         *bool    `form:"is_featured"`
	SortBy             string   `form:"sort_by,default=created_at"`
	SortOrder          string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page               int      `form:"page,default=1" binding:"min=1"`
	PageSize           int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// CreateEstateRequest 创建屋苑请求
type CreateEstateRequest struct {
	Name               string `json:"name" binding:"required,max=200"`
	NameEn             string `json:"name_en" binding:"omitempty,max=200"`
	Address            string `json:"address" binding:"required,max=500"`
	DistrictID         uint   `json:"district_id" binding:"required"`
	TotalBlocks        int    `json:"total_blocks" binding:"omitempty,min=1"`
	TotalUnits         int    `json:"total_units" binding:"omitempty,min=1"`
	CompletionYear     int    `json:"completion_year" binding:"omitempty,min=1900,max=2100"`
	Developer          string `json:"developer" binding:"omitempty,max=200"`
	ManagementCompany  string `json:"management_company" binding:"omitempty,max=200"`
	PrimarySchoolNet   string `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string `json:"secondary_school_net" binding:"omitempty,max=50"`
	Description        string `json:"description"`
	IsFeatured         bool   `json:"is_featured"`
	FacilityIDs        []uint `json:"facility_ids"`
	ImageURLs          []string `json:"image_urls"`
}

// UpdateEstateRequest 更新屋苑请求
type UpdateEstateRequest struct {
	Name               string `json:"name" binding:"omitempty,max=200"`
	NameEn             string `json:"name_en" binding:"omitempty,max=200"`
	Address            string `json:"address" binding:"omitempty,max=500"`
	DistrictID         uint   `json:"district_id"`
	TotalBlocks        int    `json:"total_blocks" binding:"omitempty,min=1"`
	TotalUnits         int    `json:"total_units" binding:"omitempty,min=1"`
	CompletionYear     int    `json:"completion_year" binding:"omitempty,min=1900,max=2100"`
	Developer          string `json:"developer" binding:"omitempty,max=200"`
	ManagementCompany  string `json:"management_company" binding:"omitempty,max=200"`
	PrimarySchoolNet   string `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string `json:"secondary_school_net" binding:"omitempty,max=50"`
	Description        string `json:"description"`
	IsFeatured         *bool  `json:"is_featured"`
	FacilityIDs        []uint `json:"facility_ids"`
	ImageURLs          []string `json:"image_urls"`
}
