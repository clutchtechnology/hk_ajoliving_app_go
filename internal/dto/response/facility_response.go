package response

// FacilityDetailResponse 设施详细响应
type FacilityDetailResponse struct {
	ID         uint    `json:"id"`
	NameZhHant string  `json:"name_zh_hant"`
	NameZhHans *string `json:"name_zh_hans,omitempty"`
	NameEn     *string `json:"name_en,omitempty"`
	Icon       *string `json:"icon,omitempty"`
	Category   string  `json:"category"`
	SortOrder  int     `json:"sort_order"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// FacilityListResponse 设施列表响应
type FacilityListResponse struct {
	Facilities []*FacilityDetailResponse `json:"facilities"`
	Pagination *Pagination               `json:"pagination"`
}
