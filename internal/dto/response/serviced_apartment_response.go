package response

import "time"

// ServicedApartmentResponse 服务式公寓详情响应
type ServicedApartmentResponse struct {
	ID            uint                          `json:"id"`
	Name          string                        `json:"name"`
	NameEn        string                        `json:"name_en,omitempty"`
	Address       string                        `json:"address"`
	DistrictID    uint                          `json:"district_id"`
	District      *DistrictResponse             `json:"district,omitempty"`
	Description   string                        `json:"description,omitempty"`
	Phone         string                        `json:"phone"`
	WebsiteURL    string                        `json:"website_url,omitempty"`
	Email         string                        `json:"email,omitempty"`
	CheckInTime   string                        `json:"check_in_time,omitempty"`
	CheckOutTime  string                        `json:"check_out_time,omitempty"`
	MinStayDays   int                           `json:"min_stay_days,omitempty"`
	Status        string                        `json:"status"`
	Rating        float64                       `json:"rating,omitempty"`
	ReviewCount   int                           `json:"review_count"`
	ViewCount     int                           `json:"view_count"`
	FavoriteCount int                           `json:"favorite_count"`
	IsFeatured    bool                          `json:"is_featured"`
	MinPrice      float64                       `json:"min_price"`
	Units         []ServicedApartmentUnitResponse `json:"units,omitempty"`
	Images        []ServicedApartmentImageResponse `json:"images,omitempty"`
	Facilities    []FacilityResponse            `json:"facilities,omitempty"`
	CreatedAt     time.Time                     `json:"created_at"`
	UpdatedAt     time.Time                     `json:"updated_at"`
}

// ServicedApartmentListItemResponse 服务式公寓列表项响应
type ServicedApartmentListItemResponse struct {
	ID            uint              `json:"id"`
	Name          string            `json:"name"`
	NameEn        string            `json:"name_en,omitempty"`
	Address       string            `json:"address"`
	DistrictID    uint              `json:"district_id"`
	District      *DistrictResponse `json:"district,omitempty"`
	Phone         string            `json:"phone"`
	MinStayDays   int               `json:"min_stay_days,omitempty"`
	Status        string            `json:"status"`
	Rating        float64           `json:"rating,omitempty"`
	ReviewCount   int               `json:"review_count"`
	ViewCount     int               `json:"view_count"`
	FavoriteCount int               `json:"favorite_count"`
	IsFeatured    bool              `json:"is_featured"`
	MinPrice      float64           `json:"min_price"`
	CoverImage    string            `json:"cover_image,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

// ServicedApartmentUnitResponse 服务式公寓房型响应
type ServicedApartmentUnitResponse struct {
	ID           uint    `json:"id"`
	UnitType     string  `json:"unit_type"`
	Bedrooms     int     `json:"bedrooms"`
	Bathrooms    int     `json:"bathrooms,omitempty"`
	Area         float64 `json:"area"`
	DailyPrice   float64 `json:"daily_price"`
	WeeklyPrice  float64 `json:"weekly_price,omitempty"`
	MonthlyPrice float64 `json:"monthly_price"`
	MaxGuests    int     `json:"max_guests"`
	Description  string  `json:"description,omitempty"`
}

// ServicedApartmentImageResponse 服务式公寓图片响应
type ServicedApartmentImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	ImageType string `json:"image_type"`
	Title     string `json:"title,omitempty"`
	SortOrder int    `json:"sort_order"`
}

// CreateServicedApartmentResponse 创建服务式公寓响应
type CreateServicedApartmentResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

// UpdateServicedApartmentResponse 更新服务式公寓响应
type UpdateServicedApartmentResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
