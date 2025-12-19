package request

// ListSchoolNetsRequest 校网列表请求
type ListSchoolNetsRequest struct {
	DistrictID *uint   `form:"district_id"`                                    // 地区ID筛选
	Level      *string `form:"level" binding:"omitempty,oneof=primary secondary"` // 级别筛选
	Keyword    string  `form:"keyword"`                                        // 关键词搜索
	Page       int     `form:"page,default=1" binding:"min=1"`                 // 页码
	PageSize   int     `form:"page_size,default=20" binding:"min=1,max=100"`   // 每页数量
	SortBy     string  `form:"sort_by,default=net_code"`                       // 排序字段
	SortOrder  string  `form:"sort_order,default=asc" binding:"omitempty,oneof=asc desc"` // 排序方向
}

// ListSchoolsRequest 学校列表请求
type ListSchoolsRequest struct {
	SchoolNetID *uint   `form:"school_net_id"`                                  // 校网ID筛选
	DistrictID  *uint   `form:"district_id"`                                    // 地区ID筛选
	Category    *string `form:"category" binding:"omitempty,oneof=government aided direct_subsidy private international"` // 类别筛选
	Level       *string `form:"level" binding:"omitempty,oneof=kindergarten primary secondary"` // 级别筛选
	Gender      *string `form:"gender" binding:"omitempty,oneof=co-ed boys girls"` // 性别筛选
	Keyword     string  `form:"keyword"`                                        // 关键词搜索
	Page        int     `form:"page,default=1" binding:"min=1"`                 // 页码
	PageSize    int     `form:"page_size,default=20" binding:"min=1,max=100"`   // 每页数量
	SortBy      string  `form:"sort_by,default=name_zh_hant"`                   // 排序字段
	SortOrder   string  `form:"sort_order,default=asc" binding:"omitempty,oneof=asc desc"` // 排序方向
}

// SearchSchoolNetsRequest 搜索校网请求
type SearchSchoolNetsRequest struct {
	Keyword   string `form:"keyword" binding:"required"` // 搜索关键词
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// SearchSchoolsRequest 搜索学校请求
type SearchSchoolsRequest struct {
	Keyword   string `form:"keyword" binding:"required"` // 搜索关键词
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
