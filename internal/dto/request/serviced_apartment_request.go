package request

// ListServicedApartmentsRequest 服务式公寓列表请求
type ListServicedApartmentsRequest struct {
	DistrictID  *uint    `form:"district_id"`
	MinPrice    *float64 `form:"min_price"`
	MaxPrice    *float64 `form:"max_price"`
	MinStayDays *int     `form:"min_stay_days"`
	Status      string   `form:"status" binding:"omitempty,oneof=active inactive closed"`
	IsFeatured  *bool    `form:"is_featured"`
	SortBy      string   `form:"sort_by,default=created_at"`
	SortOrder   string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page        int      `form:"page,default=1" binding:"min=1"`
	PageSize    int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// CreateServicedApartmentRequest 创建服务式公寓请求
type CreateServicedApartmentRequest struct {
	Name          string   `json:"name" binding:"required,max=200"`
	NameEn        string   `json:"name_en" binding:"omitempty,max=200"`
	Address       string   `json:"address" binding:"required,max=500"`
	DistrictID    uint     `json:"district_id" binding:"required"`
	Description   string   `json:"description"`
	Phone         string   `json:"phone" binding:"required,max=50"`
	WebsiteURL    string   `json:"website_url" binding:"omitempty,url,max=500"`
	Email         string   `json:"email" binding:"omitempty,email,max=255"`
	CheckInTime   string   `json:"check_in_time" binding:"omitempty,max=50"`
	CheckOutTime  string   `json:"check_out_time" binding:"omitempty,max=50"`
	MinStayDays   int      `json:"min_stay_days" binding:"omitempty,min=1"`
	Status        string   `json:"status" binding:"required,oneof=active inactive closed"`
	IsFeatured    bool     `json:"is_featured"`
	FacilityIDs   []uint   `json:"facility_ids"`
	ImageURLs     []string `json:"image_urls"`
}

// UpdateServicedApartmentRequest 更新服务式公寓请求
type UpdateServicedApartmentRequest struct {
	Name         string   `json:"name" binding:"omitempty,max=200"`
	NameEn       string   `json:"name_en" binding:"omitempty,max=200"`
	Address      string   `json:"address" binding:"omitempty,max=500"`
	DistrictID   uint     `json:"district_id"`
	Description  string   `json:"description"`
	Phone        string   `json:"phone" binding:"omitempty,max=50"`
	WebsiteURL   string   `json:"website_url" binding:"omitempty,url,max=500"`
	Email        string   `json:"email" binding:"omitempty,email,max=255"`
	CheckInTime  string   `json:"check_in_time" binding:"omitempty,max=50"`
	CheckOutTime string   `json:"check_out_time" binding:"omitempty,max=50"`
	MinStayDays  int      `json:"min_stay_days" binding:"omitempty,min=1"`
	Status       string   `json:"status" binding:"omitempty,oneof=active inactive closed"`
	IsFeatured   *bool    `json:"is_featured"`
	FacilityIDs  []uint   `json:"facility_ids"`
	ImageURLs    []string `json:"image_urls"`
}
