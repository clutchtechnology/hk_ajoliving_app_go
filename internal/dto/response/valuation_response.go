package response

import "time"

// ValuationResponse 屋苑估价响应
type ValuationResponse struct {
	EstateID                uint       `json:"estate_id"`
	EstateName              string     `json:"estate_name"`
	EstateNameEn            string     `json:"estate_name_en,omitempty"`
	Address                 string     `json:"address"`
	DistrictID              uint       `json:"district_id"`
	DistrictName            string     `json:"district_name"`
	TotalBlocks             int        `json:"total_blocks,omitempty"`
	TotalUnits              int        `json:"total_units,omitempty"`
	CompletionYear          int        `json:"completion_year,omitempty"`
	PrimarySchoolNet        string     `json:"primary_school_net,omitempty"`
	SecondarySchoolNet      string     `json:"secondary_school_net,omitempty"`
	AvgPricePerSqft         float64    `json:"avg_price_per_sqft"`         // 平均每平方尺价格
	AvgTransactionPrice     float64    `json:"avg_transaction_price"`      // 平均成交价
	RecentTransactionsCount int        `json:"recent_transactions_count"`  // 最近成交数量
	PriceUpdatedAt          *time.Time `json:"price_updated_at,omitempty"` // 价格更新时间
	PriceTrend              string     `json:"price_trend"`                // up, down, stable
	PriceTrendPercentage    float64    `json:"price_trend_percentage"`     // 涨跌百分比
}

// ValuationListItemResponse 屋苑估价列表项响应
type ValuationListItemResponse struct {
	EstateID                uint    `json:"estate_id"`
	EstateName              string  `json:"estate_name"`
	Address                 string  `json:"address"`
	DistrictName            string  `json:"district_name"`
	CompletionYear          int     `json:"completion_year,omitempty"`
	AvgPricePerSqft         float64 `json:"avg_price_per_sqft"`
	AvgTransactionPrice     float64 `json:"avg_transaction_price"`
	RecentTransactionsCount int     `json:"recent_transactions_count"`
	PriceTrend              string  `json:"price_trend"`
	PriceTrendPercentage    float64 `json:"price_trend_percentage"`
}

// DistrictValuationResponse 地区估价响应
type DistrictValuationResponse struct {
	DistrictID              uint                        `json:"district_id"`
	DistrictName            string                      `json:"district_name"`
	TotalEstates            int                         `json:"total_estates"`
	AvgPricePerSqft         float64                     `json:"avg_price_per_sqft"`
	MedianPricePerSqft      float64                     `json:"median_price_per_sqft"`
	MinPricePerSqft         float64                     `json:"min_price_per_sqft"`
	MaxPricePerSqft         float64                     `json:"max_price_per_sqft"`
	TotalTransactions       int                         `json:"total_transactions"`
	PriceTrend              string                      `json:"price_trend"`
	PriceTrendPercentage    float64                     `json:"price_trend_percentage"`
	Estates                 []ValuationListItemResponse `json:"estates"`
}
