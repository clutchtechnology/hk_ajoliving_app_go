package response

// ConfigResponse 系统配置响应
type ConfigResponse struct {
	System     *SystemConfig     `json:"system"`
	App        *AppConfig        `json:"app"`
	Features   *FeaturesConfig   `json:"features"`
	API        *APIConfig        `json:"api"`
	UI         *UIConfig         `json:"ui"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	Version     string `json:"version"`
	Environment string `json:"environment"` // dev, staging, production
	APIBaseURL  string `json:"api_base_url"`
	WebURL      string `json:"web_url"`
	Timezone    string `json:"timezone"`
	Language    string `json:"language"`
	Currency    string `json:"currency"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Favicon     string `json:"favicon"`
	Copyright   string `json:"copyright"`
	SupportEmail string `json:"support_email"`
	SupportPhone string `json:"support_phone"`
}

// FeaturesConfig 功能配置
type FeaturesConfig struct {
	EnableRegistration     bool `json:"enable_registration"`
	EnableSocialLogin      bool `json:"enable_social_login"`
	EnablePropertyListing  bool `json:"enable_property_listing"`
	EnableFurnitureStore   bool `json:"enable_furniture_store"`
	EnableMortgage         bool `json:"enable_mortgage"`
	EnableValuation        bool `json:"enable_valuation"`
	EnableNews             bool `json:"enable_news"`
	EnableSchoolNet        bool `json:"enable_school_net"`
	EnablePriceIndex       bool `json:"enable_price_index"`
}

// APIConfig API配置
type APIConfig struct {
	Version      string   `json:"version"`
	RateLimit    int      `json:"rate_limit"`     // 每分钟请求限制
	Timeout      int      `json:"timeout"`        // 请求超时时间（秒）
	MaxPageSize  int      `json:"max_page_size"`  // 最大分页大小
	AllowedOrigins []string `json:"allowed_origins"`
}

// UIConfig UI配置
type UIConfig struct {
	Theme        string            `json:"theme"`         // light, dark
	PrimaryColor string            `json:"primary_color"`
	MapProvider  string            `json:"map_provider"`  // google, mapbox
	MapAPIKey    string            `json:"map_api_key"`
	Languages    []LanguageOption  `json:"languages"`
}

// LanguageOption 语言选项
type LanguageOption struct {
	Code  string `json:"code"`  // zh-HK, zh-CN, en
	Name  string `json:"name"`
	Flag  string `json:"flag"`
}

// RegionsResponse 区域配置响应
type RegionsResponse struct {
	Regions []*RegionConfig `json:"regions"`
}

// RegionConfig 区域配置
type RegionConfig struct {
	ID          uint              `json:"id"`
	Code        string            `json:"code"`        // HK, KLN, NT
	NameZhHant  string            `json:"name_zh_hant"`
	NameZhHans  string            `json:"name_zh_hans"`
	NameEn      string            `json:"name_en"`
	Type        string            `json:"type"`        // island, peninsula, territories
	DisplayOrder int              `json:"display_order"`
	Districts   []*DistrictConfig `json:"districts"`
}

// DistrictConfig 地区配置
type DistrictConfig struct {
	ID          uint   `json:"id"`
	RegionID    uint   `json:"region_id"`
	Code        string `json:"code"`
	NameZhHant  string `json:"name_zh_hant"`
	NameZhHans  string `json:"name_zh_hans"`
	NameEn      string `json:"name_en"`
	DisplayOrder int   `json:"display_order"`
	PropertyCount int  `json:"property_count"`
	EstateCount   int  `json:"estate_count"`
}

// PropertyTypesResponse 房产类型配置响应
type PropertyTypesResponse struct {
	PropertyTypes []*PropertyTypeConfig `json:"property_types"`
	ListingTypes  []*ListingTypeConfig  `json:"listing_types"`
	Statuses      []*StatusConfig       `json:"statuses"`
}

// PropertyTypeConfig 房产类型配置
type PropertyTypeConfig struct {
	Code        string `json:"code"`        // apartment, villa, townhouse, etc.
	NameZhHant  string `json:"name_zh_hant"`
	NameZhHans  string `json:"name_zh_hans"`
	NameEn      string `json:"name_en"`
	Icon        string `json:"icon"`
	DisplayOrder int   `json:"display_order"`
	Description string `json:"description"`
}

// ListingTypeConfig 房源类型配置
type ListingTypeConfig struct {
	Code        string `json:"code"`        // rent, sale
	NameZhHant  string `json:"name_zh_hant"`
	NameZhHans  string `json:"name_zh_hans"`
	NameEn      string `json:"name_en"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

// StatusConfig 状态配置
type StatusConfig struct {
	Code        string `json:"code"`        // active, pending, sold, rented, etc.
	NameZhHant  string `json:"name_zh_hant"`
	NameZhHans  string `json:"name_zh_hans"`
	NameEn      string `json:"name_en"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
