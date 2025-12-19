package request

// ListAgenciesRequest 获取代理公司列表请求
type ListAgenciesRequest struct {
	DistrictID  *uint   `form:"district_id"`                                                 // 服务地区ID
	IsVerified  *bool   `form:"is_verified"`                                                 // 是否已验证
	MinRating   *float64 `form:"min_rating" binding:"omitempty,gte=0,lte=5"`                 // 最低评分
	Keyword     string  `form:"keyword"`                                                     // 关键词搜索（公司名称、简介）
	Page        int     `form:"page,default=1" binding:"min=1"`                              // 页码
	PageSize    int     `form:"page_size,default=20" binding:"min=1,max=100"`                // 每页数量
	SortBy      string  `form:"sort_by,default=rating" binding:"omitempty,oneof=rating agent_count created_at"` // 排序字段
	SortOrder   string  `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`  // 排序方向
}

// SearchAgenciesRequest 搜索代理公司请求
type SearchAgenciesRequest struct {
	Keyword   string `form:"keyword" binding:"required"`                                    // 搜索关键词
	Page      int    `form:"page,default=1" binding:"min=1"`                                // 页码
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`                  // 每页数量
}

// ContactAgencyRequest 联系代理公司请求
type ContactAgencyRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`                            // 联系人姓名
	Phone       string  `json:"phone" binding:"required,max=20"`                            // 联系电话
	Email       *string `json:"email" binding:"omitempty,email,max=255"`                    // 电子邮箱（可选）
	PropertyID  *uint   `json:"property_id"`                                                 // 相关房产ID（可选）
	Subject     string  `json:"subject" binding:"required,max=200"`                         // 咨询主题
	Message     string  `json:"message" binding:"required,max=2000"`                        // 咨询内容
	ContactTime *string `json:"contact_time" binding:"omitempty,max=100"`                   // 期望联系时间（可选）
}
