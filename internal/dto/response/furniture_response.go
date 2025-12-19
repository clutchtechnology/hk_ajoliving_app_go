package response

import "time"

// FurnitureListItemResponse 家具列表项响应
type FurnitureListItemResponse struct {
	ID                 uint      `json:"id"`
	FurnitureNo        string    `json:"furniture_no"`
	Title              string    `json:"title"`
	Price              float64   `json:"price"`
	CategoryID         uint      `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	Brand              *string   `json:"brand,omitempty"`
	Condition          string    `json:"condition"`
	DeliveryDistrictID uint      `json:"delivery_district_id"`
	DeliveryDistrict   string    `json:"delivery_district"`
	DeliveryMethod     string    `json:"delivery_method"`
	Status             string    `json:"status"`
	ViewCount          int       `json:"view_count"`
	FavoriteCount      int       `json:"favorite_count"`
	CoverImage         *string   `json:"cover_image,omitempty"`
	PublishedAt        time.Time `json:"published_at"`
	ExpiresAt          time.Time `json:"expires_at"`
	DaysUntilExpiry    int       `json:"days_until_expiry"`
}

// FurnitureResponse 家具详情响应
type FurnitureResponse struct {
	ID                 uint                   `json:"id"`
	FurnitureNo        string                 `json:"furniture_no"`
	Title              string                 `json:"title"`
	Description        *string                `json:"description,omitempty"`
	Price              float64                `json:"price"`
	CategoryID         uint                   `json:"category_id"`
	Category           *FurnitureCategoryResponse `json:"category,omitempty"`
	Brand              *string                `json:"brand,omitempty"`
	Condition          string                 `json:"condition"`
	PurchaseDate       *time.Time             `json:"purchase_date,omitempty"`
	Age                int                    `json:"age"`
	DeliveryDistrictID uint                   `json:"delivery_district_id"`
	DeliveryDistrict   *DistrictBasicResponse `json:"delivery_district,omitempty"`
	DeliveryTime       *string                `json:"delivery_time,omitempty"`
	DeliveryMethod     string                 `json:"delivery_method"`
	SupportsDelivery   bool                   `json:"supports_delivery"`
	SupportsSelfPickup bool                   `json:"supports_self_pickup"`
	Status             string                 `json:"status"`
	PublisherID        uint                   `json:"publisher_id"`
	Publisher          *PublisherBasicResponse `json:"publisher,omitempty"`
	ViewCount          int                    `json:"view_count"`
	FavoriteCount      int                    `json:"favorite_count"`
	Images             []FurnitureImageResponse `json:"images"`
	PublishedAt        time.Time              `json:"published_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	ExpiresAt          time.Time              `json:"expires_at"`
	DaysUntilExpiry    int                    `json:"days_until_expiry"`
	IsAvailable        bool                   `json:"is_available"`
	IsExpired          bool                   `json:"is_expired"`
	CreatedAt          time.Time              `json:"created_at"`
}

// FurnitureCategoryResponse 家具分类响应
type FurnitureCategoryResponse struct {
	ID            uint                         `json:"id"`
	ParentID      *uint                        `json:"parent_id,omitempty"`
	NameZhHant    string                       `json:"name_zh_hant"`
	NameZhHans    *string                      `json:"name_zh_hans,omitempty"`
	NameEn        *string                      `json:"name_en,omitempty"`
	Icon          *string                      `json:"icon,omitempty"`
	SortOrder     int                          `json:"sort_order"`
	IsActive      bool                         `json:"is_active"`
	IsTopLevel    bool                         `json:"is_top_level"`
	Subcategories []FurnitureCategoryResponse  `json:"subcategories,omitempty"`
	FurnitureCount int                         `json:"furniture_count,omitempty"`
}

// FurnitureImageResponse 家具图片响应
type FurnitureImageResponse struct {
	ID          uint   `json:"id"`
	ImageURL    string `json:"image_url"`
	SortOrder   int    `json:"sort_order"`
	IsCover     bool   `json:"is_cover"`
}

// PublisherBasicResponse 发布者基本信息响应
type PublisherBasicResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Avatar       *string `json:"avatar,omitempty"`
	PublisherType string `json:"publisher_type"`
}

// DistrictBasicResponse 地区基本信息响应
type DistrictBasicResponse struct {
	ID         uint   `json:"id"`
	NameZhHant string `json:"name_zh_hant"`
	NameZhHans *string `json:"name_zh_hans,omitempty"`
	NameEn     *string `json:"name_en,omitempty"`
}

// CreateFurnitureResponse 创建家具响应
type CreateFurnitureResponse struct {
	ID          uint      `json:"id"`
	FurnitureNo string    `json:"furniture_no"`
	Title       string    `json:"title"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	PublishedAt time.Time `json:"published_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Message     string    `json:"message"`
}

// UpdateFurnitureResponse 更新家具响应
type UpdateFurnitureResponse struct {
	ID          uint      `json:"id"`
	FurnitureNo string    `json:"furniture_no"`
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updated_at"`
	Message     string    `json:"message"`
}

// UpdateFurnitureStatusResponse 更新家具状态响应
type UpdateFurnitureStatusResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	Message   string    `json:"message"`
}
