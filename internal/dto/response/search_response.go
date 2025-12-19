package response

// GlobalSearchResponse 全局搜索响应
type GlobalSearchResponse struct {
	Properties []*PropertySearchResult `json:"properties,omitempty"`
	Estates    []*EstateSearchResult   `json:"estates,omitempty"`
	Agents     []*AgentSearchResult    `json:"agents,omitempty"`
	News       []*NewsSearchResult     `json:"news,omitempty"`
	TotalCount int                     `json:"total_count"`
	Keyword    string                  `json:"keyword"`
}

// PropertySearchResult 房产搜索结果
type PropertySearchResult struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Price       float64  `json:"price"`
	Area        float64  `json:"area"`
	Bedrooms    int      `json:"bedrooms"`
	Bathrooms   int      `json:"bathrooms"`
	ListingType string   `json:"listing_type"`
	Address     string   `json:"address"`
	District    *string  `json:"district,omitempty"`
	EstateName  *string  `json:"estate_name,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

// EstateSearchResult 屋苑搜索结果
type EstateSearchResult struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	NameEn          *string `json:"name_en,omitempty"`
	District        string  `json:"district"`
	Address         string  `json:"address"`
	BuildingCount   *int    `json:"building_count,omitempty"`
	UnitCount       *int    `json:"unit_count,omitempty"`
	PropertyCount   int     `json:"property_count"`
	ImageURL        *string `json:"image_url,omitempty"`
}

// AgentSearchResult 代理人搜索结果
type AgentSearchResult struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	LicenseNo      *string `json:"license_no,omitempty"`
	AgencyName     *string `json:"agency_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
	PropertyCount  int     `json:"property_count"`
	Rating         float64 `json:"rating"`
}

// NewsSearchResult 新闻搜索结果
type NewsSearchResult struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Summary     *string `json:"summary,omitempty"`
	Category    string  `json:"category"`
	ImageURL    *string `json:"image_url,omitempty"`
	ViewCount   int     `json:"view_count"`
	PublishedAt string  `json:"published_at"`
}

// SearchPropertiesResponse 搜索房产响应
type SearchPropertiesResponse struct {
	Properties []*PropertySearchResult `json:"properties"`
	Pagination *Pagination             `json:"pagination"`
	Keyword    string                  `json:"keyword"`
}

// SearchEstatesResponse 搜索屋苑响应
type SearchEstatesResponse struct {
	Estates    []*EstateSearchResult `json:"estates"`
	Pagination *Pagination           `json:"pagination"`
	Keyword    string                `json:"keyword"`
}

// SearchAgentsResponse 搜索代理人响应
type SearchAgentsResponse struct {
	Agents     []*AgentSearchResult `json:"agents"`
	Pagination *Pagination          `json:"pagination"`
	Keyword    string               `json:"keyword"`
}

// SearchSuggestion 搜索建议项
type SearchSuggestion struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Count int    `json:"count,omitempty"`
}

// SearchSuggestionsResponse 搜索建议响应
type SearchSuggestionsResponse struct {
	Suggestions []*SearchSuggestion `json:"suggestions"`
	Keyword     string              `json:"keyword"`
}

// SearchHistoryItem 搜索历史项
type SearchHistoryItem struct {
	ID          uint   `json:"id"`
	Keyword     string `json:"keyword"`
	SearchType  string `json:"search_type"`
	ResultCount int    `json:"result_count"`
	CreatedAt   string `json:"created_at"`
}

// SearchHistoryResponse 搜索历史响应
type SearchHistoryResponse struct {
	Histories  []*SearchHistoryItem `json:"histories"`
	Pagination *Pagination          `json:"pagination"`
}
