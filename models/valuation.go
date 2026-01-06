package models

// valuation.go - 物业估价相关的响应类型定义

// ============ 估价响应 DTO ============

// ValuationListItemResponse 估价列表项响应
type ValuationListItemResponse struct {
	EstateID                uint     `json:"estate_id"`
	EstateName              string   `json:"estate_name"`
	EstateNameEn            string   `json:"estate_name_en,omitempty"`
	Address                 string   `json:"address"`
	DistrictName            string   `json:"district_name"`
	CompletionYear          int      `json:"completion_year,omitempty"`
	AvgTransactionPrice     float64  `json:"avg_transaction_price"`
	AvgPricePerSqft         float64  `json:"avg_price_per_sqft"`
	RecentTransactionsCount int      `json:"recent_transactions_count"`
	PriceTrend              string   `json:"price_trend"`              // up, down, stable
	PriceTrendPercentage    float64  `json:"price_trend_percentage"`   // 价格趋势百分比
}

// ValuationDetailResponse 估价详情响应
type ValuationDetailResponse struct {
	EstateID                uint     `json:"estate_id"`
	EstateName              string   `json:"estate_name"`
	EstateNameEn            string   `json:"estate_name_en,omitempty"`
	Address                 string   `json:"address"`
	DistrictID              uint     `json:"district_id"`
	DistrictName            string   `json:"district_name"`
	TotalBlocks             int      `json:"total_blocks,omitempty"`
	TotalUnits              int      `json:"total_units,omitempty"`
	CompletionYear          int      `json:"completion_year,omitempty"`
	PrimarySchoolNet        string   `json:"primary_school_net,omitempty"`
	SecondarySchoolNet      string   `json:"secondary_school_net,omitempty"`
	AvgTransactionPrice     float64  `json:"avg_transaction_price"`
	AvgPricePerSqft         float64  `json:"avg_price_per_sqft"`
	RecentTransactionsCount int      `json:"recent_transactions_count"`
	PriceUpdatedAt          *string  `json:"price_updated_at,omitempty"`
	PriceTrend              string   `json:"price_trend"`              // up, down, stable
	PriceTrendPercentage    float64  `json:"price_trend_percentage"`   // 价格趋势百分比
}

// DistrictValuationResponse 地区估价响应
type DistrictValuationResponse struct {
	DistrictID            uint                         `json:"district_id"`
	DistrictName          string                       `json:"district_name"`
	TotalEstates          int                          `json:"total_estates"`
	AvgPricePerSqft       float64                      `json:"avg_price_per_sqft"`
	MedianPricePerSqft    float64                      `json:"median_price_per_sqft"`
	MinPricePerSqft       float64                      `json:"min_price_per_sqft"`
	MaxPricePerSqft       float64                      `json:"max_price_per_sqft"`
	TotalTransactions     int                          `json:"total_transactions"`
	PriceTrend            string                       `json:"price_trend"`
	PriceTrendPercentage  float64                      `json:"price_trend_percentage"`
	Estates               []ValuationListItemResponse  `json:"estates"`
}
