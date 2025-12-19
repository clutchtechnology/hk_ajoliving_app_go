package request

// ListAgentsRequest 代理人列表请求
type ListAgentsRequest struct {
	AgencyID       *uint    `form:"agency_id"`                                      // 代理公司ID筛选
	DistrictID     *uint    `form:"district_id"`                                    // 服务地区筛选
	Status         *string  `form:"status" binding:"omitempty,oneof=active inactive suspended"` // 状态筛选
	IsVerified     *bool    `form:"is_verified"`                                    // 是否已验证
	Specialization string   `form:"specialization"`                                 // 专业领域筛选
	MinRating      *float64 `form:"min_rating" binding:"omitempty,min=0,max=5"`     // 最低评分
	Keyword        string   `form:"keyword"`                                        // 关键词搜索
	Page           int      `form:"page,default=1" binding:"min=1"`                 // 页码
	PageSize       int      `form:"page_size,default=20" binding:"min=1,max=100"`   // 每页数量
	SortBy         string   `form:"sort_by,default=rating"`                         // 排序字段
	SortOrder      string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"` // 排序方向
}

// ContactAgentRequest 联系代理人请求
type ContactAgentRequest struct {
	PropertyID  *uint  `json:"property_id"`                                                // 房源ID（可选）
	Name        string `json:"name" binding:"required,max=100"`                            // 联系人姓名
	Phone       string `json:"phone" binding:"required,max=20"`                            // 联系电话
	Email       string `json:"email" binding:"omitempty,email,max=255"`                    // 联系邮箱
	Message     string `json:"message" binding:"required"`                                 // 留言内容
	ContactType string `json:"contact_type" binding:"omitempty,oneof=inquiry viewing valuation other"` // 联系类型
}
