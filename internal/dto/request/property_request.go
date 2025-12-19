package request

// ListPropertiesRequest 房产列表请求
type ListPropertiesRequest struct {
	DistrictID     *uint    `form:"district_id"`
	BuildingName   string   `form:"building_name"`
	MinPrice       *float64 `form:"min_price"`
	MaxPrice       *float64 `form:"max_price"`
	MinArea        *float64 `form:"min_area"`
	MaxArea        *float64 `form:"max_area"`
	Bedrooms       *int     `form:"bedrooms"`
	PropertyType   string   `form:"property_type"`
	ListingType    string   `form:"listing_type" binding:"omitempty,oneof=rent sale"`
	Status         string   `form:"status" binding:"omitempty,oneof=available pending sold cancelled"`
	SchoolNet      string   `form:"school_net"`
	SortBy         string   `form:"sort_by,default=created_at"`
	SortOrder      string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page           int      `form:"page,default=1" binding:"min=1"`
	PageSize       int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// CreatePropertyRequest 创建房产请求
type CreatePropertyRequest struct {
	Title              string   `json:"title" binding:"required,max=255"`
	Description        string   `json:"description"`
	Area               float64  `json:"area" binding:"required,gt=0"`
	Price              float64  `json:"price" binding:"required,gt=0"`
	Address            string   `json:"address" binding:"required,max=500"`
	DistrictID         uint     `json:"district_id" binding:"required"`
	BuildingName       string   `json:"building_name" binding:"omitempty,max=200"`
	Floor              string   `json:"floor" binding:"omitempty,max=20"`
	Orientation        string   `json:"orientation" binding:"omitempty,max=50"`
	Bedrooms           int      `json:"bedrooms" binding:"min=0"`
	Bathrooms          int      `json:"bathrooms" binding:"min=0"`
	PrimarySchoolNet   string   `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string   `json:"secondary_school_net" binding:"omitempty,max=50"`
	PropertyType       string   `json:"property_type" binding:"required,oneof=apartment villa townhouse studio duplex penthouse shophouse"`
	ListingType        string   `json:"listing_type" binding:"required,oneof=rent sale"`
	FacilityIDs        []uint   `json:"facility_ids"`
	ImageURLs          []string `json:"image_urls"`
}

// UpdatePropertyRequest 更新房产请求
type UpdatePropertyRequest struct {
	Title              string   `json:"title" binding:"omitempty,max=255"`
	Description        string   `json:"description"`
	Area               float64  `json:"area" binding:"omitempty,gt=0"`
	Price              float64  `json:"price" binding:"omitempty,gt=0"`
	Address            string   `json:"address" binding:"omitempty,max=500"`
	DistrictID         uint     `json:"district_id"`
	BuildingName       string   `json:"building_name" binding:"omitempty,max=200"`
	Floor              string   `json:"floor" binding:"omitempty,max=20"`
	Orientation        string   `json:"orientation" binding:"omitempty,max=50"`
	Bedrooms           int      `json:"bedrooms" binding:"min=0"`
	Bathrooms          int      `json:"bathrooms" binding:"min=0"`
	PrimarySchoolNet   string   `json:"primary_school_net" binding:"omitempty,max=50"`
	SecondarySchoolNet string   `json:"secondary_school_net" binding:"omitempty,max=50"`
	PropertyType       string   `json:"property_type" binding:"omitempty,oneof=apartment villa townhouse studio duplex penthouse shophouse"`
	Status             string   `json:"status" binding:"omitempty,oneof=available pending sold cancelled"`
	FacilityIDs        []uint   `json:"facility_ids"`
	ImageURLs          []string `json:"image_urls"`
}

// FeaturedPropertiesRequest 精选房源请求
type FeaturedPropertiesRequest struct {
	ListingType string `form:"listing_type" binding:"omitempty,oneof=rent sale"`
	Limit       int    `form:"limit,default=10" binding:"min=1,max=50"`
}

// HotPropertiesRequest 热门房源请求
type HotPropertiesRequest struct {
	ListingType string `form:"listing_type" binding:"omitempty,oneof=rent sale"`
	Limit       int    `form:"limit,default=10" binding:"min=1,max=50"`
}

// SimilarPropertiesRequest 相似房源请求
type SimilarPropertiesRequest struct {
	Limit int `form:"limit,default=6" binding:"min=1,max=20"`
}

// ListBuyPropertiesRequest 买房房源列表请求
type ListBuyPropertiesRequest struct {
	DistrictID   *uint    `form:"district_id"`
	BuildingName string   `form:"building_name"`
	MinPrice     *float64 `form:"min_price"`
	MaxPrice     *float64 `form:"max_price"`
	MinArea      *float64 `form:"min_area"`
	MaxArea      *float64 `form:"max_area"`
	Bedrooms     *int     `form:"bedrooms"`
	PropertyType string   `form:"property_type"`
	SchoolNet    string   `form:"school_net"`
	IsNew        *bool    `form:"is_new"` // 是否新房
	SortBy       string   `form:"sort_by,default=created_at"`
	SortOrder    string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page         int      `form:"page,default=1" binding:"min=1"`
	PageSize     int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListRentPropertiesRequest 租房房源列表请求
type ListRentPropertiesRequest struct {
	DistrictID   *uint    `form:"district_id"`
	BuildingName string   `form:"building_name"`
	MinPrice     *float64 `form:"min_price"`  // 月租最低价
	MaxPrice     *float64 `form:"max_price"`  // 月租最高价
	MinArea      *float64 `form:"min_area"`
	MaxArea      *float64 `form:"max_area"`
	Bedrooms     *int     `form:"bedrooms"`
	PropertyType string   `form:"property_type"`
	SchoolNet    string   `form:"school_net"`
	RentType     string   `form:"rent_type" binding:"omitempty,oneof=short long"` // 短租/长租
	SortBy       string   `form:"sort_by,default=created_at"`
	SortOrder    string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page         int      `form:"page,default=1" binding:"min=1"`
	PageSize     int      `form:"page_size,default=20" binding:"min=1,max=100"`
}
