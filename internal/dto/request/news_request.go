package request

// ListNewsRequest 新闻列表请求
type ListNewsRequest struct {
	CategoryID *uint   `form:"category_id"`                                    // 分类ID筛选
	Status     *string `form:"status" binding:"omitempty,oneof=draft published archived"` // 状态筛选
	IsFeatured *bool   `form:"is_featured"`                                    // 是否精选
	IsHot      *bool   `form:"is_hot"`                                         // 是否热门
	IsTop      *bool   `form:"is_top"`                                         // 是否置顶
	Keyword    string  `form:"keyword"`                                        // 关键词搜索
	Tag        string  `form:"tag"`                                            // 标签筛选
	Page       int     `form:"page,default=1" binding:"min=1"`                 // 页码
	PageSize   int     `form:"page_size,default=20" binding:"min=1,max=100"`   // 每页数量
	SortBy     string  `form:"sort_by,default=published_at"`                   // 排序字段
	SortOrder  string  `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"` // 排序方向
}
