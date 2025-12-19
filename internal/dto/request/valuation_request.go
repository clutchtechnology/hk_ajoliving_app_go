package request

// ListValuationsRequest 获取屋苑估价列表请求
type ListValuationsRequest struct {
	DistrictID       *uint    `form:"district_id"`
	MinPrice         *float64 `form:"min_price"`
	MaxPrice         *float64 `form:"max_price"`
	MinArea          *float64 `form:"min_area"`
	MaxArea          *float64 `form:"max_area"`
	SchoolNet        string   `form:"school_net"`
	SortBy           string   `form:"sort_by"`           // avg_price, name, completion_year
	SortOrder        string   `form:"sort_order"`        // asc, desc
	Page             int      `form:"page,default=1" binding:"min=1"`
	PageSize         int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// SearchValuationsRequest 搜索屋苑估价请求
type SearchValuationsRequest struct {
	Keyword          string   `form:"keyword" binding:"required"`
	DistrictID       *uint    `form:"district_id"`
	MinPrice         *float64 `form:"min_price"`
	MaxPrice         *float64 `form:"max_price"`
	Page             int      `form:"page,default=1" binding:"min=1"`
	PageSize         int      `form:"page_size,default=20" binding:"min=1,max=100"`
}
