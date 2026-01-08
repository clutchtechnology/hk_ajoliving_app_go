package models

import (
	"time"
)

// ============ Response DTO ============

// ValuationResponse 屋苑估价响应
type ValuationResponse struct {
	EstateID               uint      `json:"estate_id"`
	EstateName             string    `json:"estate_name"`
	EstateNameEn           string    `json:"estate_name_en,omitempty"`
	DistrictID             uint      `json:"district_id"`
	District               *District `json:"district,omitempty"`
	Address                string    `json:"address"`
	CompletionYear         int       `json:"completion_year,omitempty"`
	TotalUnits             int       `json:"total_units"`
	AvgPricePerSqft        float64   `json:"avg_price_per_sqft"`        // 平均每平方尺价格
	AvgSalePrice           float64   `json:"avg_sale_price"`            // 平均售价
	AvgRentPrice           float64   `json:"avg_rent_price"`            // 平均租金
	MinPricePerSqft        float64   `json:"min_price_per_sqft"`        // 最低每平方尺价格
	MaxPricePerSqft        float64   `json:"max_price_per_sqft"`        // 最高每平方尺价格
	RecentTransactionCount int       `json:"recent_transaction_count"`  // 近期成交数量
	ForSaleCount           int       `json:"for_sale_count"`            // 当前放盘数量
	ForRentCount           int       `json:"for_rent_count"`            // 当前租盘数量
	PriceChange30d         float64   `json:"price_change_30d"`          // 30天价格变化百分比
	PriceChange90d         float64   `json:"price_change_90d"`          // 90天价格变化百分比
	RentalYield            float64   `json:"rental_yield"`              // 租金回报率 (%)
	LastUpdated            time.Time `json:"last_updated"`              // 最后更新时间
}

// EstateValuationDetail 屋苑估价详情
type EstateValuationDetail struct {
	EstateID               uint                      `json:"estate_id"`
	EstateName             string                    `json:"estate_name"`
	EstateNameEn           string                    `json:"estate_name_en,omitempty"`
	Address                string                    `json:"address"`
	District               *District                 `json:"district,omitempty"`
	CompletionYear         int                       `json:"completion_year,omitempty"`
	Developer              string                    `json:"developer,omitempty"`
	TotalBlocks            int                       `json:"total_blocks"`
	TotalUnits             int                       `json:"total_units"`
	PrimarySchoolNet       string                    `json:"primary_school_net,omitempty"`
	SecondarySchoolNet     string                    `json:"secondary_school_net,omitempty"`
	AvgPricePerSqft        float64                   `json:"avg_price_per_sqft"`
	AvgSalePrice           float64                   `json:"avg_sale_price"`
	AvgRentPrice           float64                   `json:"avg_rent_price"`
	MinPricePerSqft        float64                   `json:"min_price_per_sqft"`
	MaxPricePerSqft        float64                   `json:"max_price_per_sqft"`
	RecentTransactionCount int                       `json:"recent_transaction_count"`
	ForSaleCount           int                       `json:"for_sale_count"`
	ForRentCount           int                       `json:"for_rent_count"`
	PriceChange30d         float64                   `json:"price_change_30d"`
	PriceChange90d         float64                   `json:"price_change_90d"`
	RentalYield            float64                   `json:"rental_yield"`
	PriceHistory           []PriceHistoryPoint       `json:"price_history"`           // 价格历史
	UnitTypePrices         []UnitTypePriceBreakdown  `json:"unit_type_prices"`        // 户型价格分布
	RecentTransactions     []TransactionSummary      `json:"recent_transactions"`     // 近期成交
	LastUpdated            time.Time                 `json:"last_updated"`
}

// PriceHistoryPoint 价格历史数据点
type PriceHistoryPoint struct {
	Date            string  `json:"date"`              // YYYY-MM 格式
	AvgPricePerSqft float64 `json:"avg_price_per_sqft"`
	TransactionCount int    `json:"transaction_count"`
}

// UnitTypePriceBreakdown 户型价格分布
type UnitTypePriceBreakdown struct {
	Bedrooms        int     `json:"bedrooms"`          // 房间数
	AvgArea         float64 `json:"avg_area"`          // 平均面积
	AvgPrice        float64 `json:"avg_price"`         // 平均价格
	AvgPricePerSqft float64 `json:"avg_price_per_sqft"`// 平均每平方尺价格
	MinPrice        float64 `json:"min_price"`         // 最低价格
	MaxPrice        float64 `json:"max_price"`         // 最高价格
	AvailableCount  int     `json:"available_count"`   // 可售数量
}

// TransactionSummary 成交摘要
type TransactionSummary struct {
	TransactionDate time.Time `json:"transaction_date"`
	PropertyType    string    `json:"property_type"`
	Bedrooms        int       `json:"bedrooms"`
	Area            float64   `json:"area"`
	Price           float64   `json:"price"`
	PricePerSqft    float64   `json:"price_per_sqft"`
}

// ListValuationsRequest 获取估价列表请求
type ListValuationsRequest struct {
	DistrictID         *uint    `form:"district_id"`                                    // 地区ID
	PrimarySchoolNet   *string  `form:"primary_school_net"`                             // 小学校网
	SecondarySchoolNet *string  `form:"secondary_school_net"`                           // 中学校网
	MinAvgPrice        *float64 `form:"min_avg_price"`                                  // 最低平均价
	MaxAvgPrice        *float64 `form:"max_avg_price"`                                  // 最高平均价
	MinRentalYield     *float64 `form:"min_rental_yield"`                               // 最低租金回报率
	Keyword            string   `form:"keyword"`                                        // 搜索关键词
	SortBy             string   `form:"sort_by" binding:"omitempty,oneof=price yield transactions name"` // 排序字段
	SortOrder          string   `form:"sort_order" binding:"omitempty,oneof=asc desc"`  // 排序方向
	Page               int      `form:"page" binding:"min=1"`                           // 页码
	PageSize           int      `form:"page_size" binding:"min=1,max=100"`              // 每页数量
}

// SearchValuationsRequest 搜索估价请求
type SearchValuationsRequest struct {
	Keyword  string `form:"keyword" binding:"required"` // 搜索关键词（屋苑名称、地址）
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// PaginatedValuationsResponse 分页估价列表响应
type PaginatedValuationsResponse struct {
	Data       []ValuationResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// DistrictValuationSummary 地区估价汇总
type DistrictValuationSummary struct {
	DistrictID        uint                `json:"district_id"`
	District          *District           `json:"district,omitempty"`
	EstateCount       int                 `json:"estate_count"`        // 屋苑数量
	AvgPricePerSqft   float64             `json:"avg_price_per_sqft"`  // 平均每平方尺价格
	MinPricePerSqft   float64             `json:"min_price_per_sqft"`  // 最低每平方尺价格
	MaxPricePerSqft   float64             `json:"max_price_per_sqft"`  // 最高每平方尺价格
	AvgRentalYield    float64             `json:"avg_rental_yield"`    // 平均租金回报率
	TotalTransactions int                 `json:"total_transactions"`  // 总成交数量
	Estates           []ValuationResponse `json:"estates"`             // 屋苑列表
}
