package request

// GlobalSearchRequest 全局搜索请求
type GlobalSearchRequest struct {
	Keyword  string `form:"keyword" binding:"required,min=1,max=200"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=50"`
}

// SearchPropertiesRequest 搜索房产请求
type SearchPropertiesRequest struct {
	Keyword      string   `form:"keyword" binding:"required,min=1,max=200"`
	ListingType  *string  `form:"listing_type" binding:"omitempty,oneof=rent sale"`
	DistrictID   *uint    `form:"district_id"`
	MinPrice     *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice     *float64 `form:"max_price" binding:"omitempty,min=0"`
	Bedrooms     *int     `form:"bedrooms" binding:"omitempty,min=0"`
	PropertyType *string  `form:"property_type"`
	Page         int      `form:"page,default=1" binding:"min=1"`
	PageSize     int      `form:"page_size,default=20" binding:"min=1,max=50"`
}

// SearchEstatesRequest 搜索屋苑请求
type SearchEstatesRequest struct {
	Keyword    string `form:"keyword" binding:"required,min=1,max=200"`
	DistrictID *uint  `form:"district_id"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=50"`
}

// SearchAgentsRequest 搜索代理人请求
type SearchAgentsRequest struct {
	Keyword    string `form:"keyword" binding:"required,min=1,max=200"`
	DistrictID *uint  `form:"district_id"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=50"`
}

// GetSearchSuggestionsRequest 获取搜索建议请求
type GetSearchSuggestionsRequest struct {
	Keyword string  `form:"keyword" binding:"required,min=1,max=200"`
	Type    *string `form:"type" binding:"omitempty,oneof=property estate agent agency"`
	Limit   int     `form:"limit,default=10" binding:"min=1,max=20"`
}

// GetSearchHistoryRequest 获取搜索历史请求
type GetSearchHistoryRequest struct {
	Type     *string `form:"type" binding:"omitempty,oneof=global property estate agent agency news school furniture"`
	Page     int     `form:"page,default=1" binding:"min=1"`
	PageSize int     `form:"page_size,default=20" binding:"min=1,max=50"`
}
