package response

// Pagination 分页信息
type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// DistrictBasicResponse 地区基本信息
type DistrictBasicResponse struct {
	ID         uint    `json:"id"`
	NameZhHant string  `json:"name_zh_hant"`
	NameZhHans *string `json:"name_zh_hans,omitempty"`
	NameEn     *string `json:"name_en,omitempty"`
}

// FacilityResponse 设施响应
type FacilityResponse struct {
	ID         uint    `json:"id"`
	NameZhHant string  `json:"name_zh_hant"`
	NameZhHans *string `json:"name_zh_hans,omitempty"`
	NameEn     *string `json:"name_en,omitempty"`
	Icon       *string `json:"icon,omitempty"`
	Category   string  `json:"category"`
}
