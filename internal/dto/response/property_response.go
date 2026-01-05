package response

import "time"

// PropertyResponse 房产详情响应
type PropertyResponse struct {
	ID                 uint                 `json:"id"`
	PropertyNo         string               `json:"property_no"`
	EstateNo           string               `json:"estate_no,omitempty"`
	ListingType        string               `json:"listing_type"`
	Title              string               `json:"title"`
	Description        string               `json:"description,omitempty"`
	Area               float64              `json:"area"`
	Price              float64              `json:"price"`
	Address            string               `json:"address"`
	DistrictID         uint                 `json:"district_id"`
	District           *DistrictResponse    `json:"district,omitempty"`
	BuildingName       string               `json:"building_name,omitempty"`
	Floor              string               `json:"floor,omitempty"`
	Orientation        string               `json:"orientation,omitempty"`
	Bedrooms           int                  `json:"bedrooms"`
	Bathrooms          int                  `json:"bathrooms,omitempty"`
	PrimarySchoolNet   string               `json:"primary_school_net,omitempty"`
	SecondarySchoolNet string               `json:"secondary_school_net,omitempty"`
	PropertyType       string               `json:"property_type"`
	Status             string               `json:"status"`
	PublisherID        uint                 `json:"publisher_id"`
	PublisherType      string               `json:"publisher_type"`
	AgentID            uint                 `json:"agent_id,omitempty"`
	Agent              *AgentBriefResponse  `json:"agent,omitempty"`
	ViewCount          int                  `json:"view_count"`
	FavoriteCount      int                  `json:"favorite_count"`
	Images             []PropertyImageResponse `json:"images,omitempty"`
	Facilities         []FacilityResponse   `json:"facilities,omitempty"`
	PublishedAt        *time.Time           `json:"published_at,omitempty"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

// PropertyListItemResponse 房产列表项响应
type PropertyListItemResponse struct {
	ID            uint              `json:"id"`
	PropertyNo    string            `json:"property_no"`
	ListingType   string            `json:"listing_type"`
	Title         string            `json:"title"`
	Area          float64           `json:"area"`
	Price         float64           `json:"price"`
	Address       string            `json:"address"`
	DistrictID    uint              `json:"district_id"`
	District      *DistrictResponse `json:"district,omitempty"`
	BuildingName  string            `json:"building_name,omitempty"`
	Bedrooms      int               `json:"bedrooms"`
	Bathrooms     int               `json:"bathrooms,omitempty"`
	PropertyType  string            `json:"property_type"`
	Status        string            `json:"status"`
	CoverImage    string            `json:"cover_image,omitempty"`
	ViewCount     int               `json:"view_count"`
	FavoriteCount int               `json:"favorite_count"`
	CreatedAt     time.Time         `json:"created_at"`
}

// PropertyImageResponse 房产图片响应
type PropertyImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	Caption   string `json:"caption,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsCover   bool   `json:"is_cover"`
}

// DistrictResponse 地区响应
type DistrictResponse struct {
	ID         uint   `json:"id"`
	NameZhHant string `json:"name_zh_hant"`
	NameZhHans string `json:"name_zh_hans,omitempty"`
	NameEn     string `json:"name_en,omitempty"`
	Region     string `json:"region"`
}

// AgentBriefResponse 代理人简要信息响应
type AgentBriefResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	LicenseNo string `json:"license_no,omitempty"`
}

// CreatePropertyResponse 创建房产响应
type CreatePropertyResponse struct {
	ID         uint   `json:"id"`
	PropertyNo string `json:"property_no"`
	Message    string `json:"message"`
}

// UpdatePropertyResponse 更新房产响应
type UpdatePropertyResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

// DeletePropertyResponse 删除房产响应
type DeletePropertyResponse struct {
	Message string `json:"message"`
}
