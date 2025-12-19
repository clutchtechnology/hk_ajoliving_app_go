package request

// ListFacilitiesRequest 获取设施列表请求
type ListFacilitiesRequest struct {
	Category  *string `form:"category" binding:"omitempty,oneof=building unit"`
	Keyword   *string `form:"keyword"`
	Page      int     `form:"page,default=1" binding:"min=1"`
	PageSize  int     `form:"page_size,default=50" binding:"min=1,max=100"`
	SortBy    string  `form:"sort_by,default=sort_order"`
	SortOrder string  `form:"sort_order,default=asc" binding:"omitempty,oneof=asc desc"`
}

// CreateFacilityRequest 创建设施请求
type CreateFacilityRequest struct {
	NameZhHant string  `json:"name_zh_hant" binding:"required,max=100"`
	NameZhHans *string `json:"name_zh_hans,omitempty" binding:"omitempty,max=100"`
	NameEn     *string `json:"name_en,omitempty" binding:"omitempty,max=100"`
	Icon       *string `json:"icon,omitempty" binding:"omitempty,max=100"`
	Category   string  `json:"category" binding:"required,oneof=building unit"`
	SortOrder  int     `json:"sort_order" binding:"min=0"`
}

// UpdateFacilityRequest 更新设施请求
type UpdateFacilityRequest struct {
	NameZhHant *string `json:"name_zh_hant,omitempty" binding:"omitempty,max=100"`
	NameZhHans *string `json:"name_zh_hans,omitempty" binding:"omitempty,max=100"`
	NameEn     *string `json:"name_en,omitempty" binding:"omitempty,max=100"`
	Icon       *string `json:"icon,omitempty" binding:"omitempty,max=100"`
	Category   *string `json:"category,omitempty" binding:"omitempty,oneof=building unit"`
	SortOrder  *int    `json:"sort_order,omitempty" binding:"omitempty,min=0"`
}
