package response

import "time"

// PriceIndexResponse 楼价指数响应
type PriceIndexResponse struct {
	ID               uint              `json:"id"`
	IndexType        string            `json:"index_type"`
	DistrictID       *uint             `json:"district_id,omitempty"`
	District         *DistrictResponse `json:"district,omitempty"`
	EstateID         *uint             `json:"estate_id,omitempty"`
	Estate           *EstateBasicInfo  `json:"estate,omitempty"`
	PropertyType     *string           `json:"property_type,omitempty"`
	IndexValue       float64           `json:"index_value"`
	ChangeValue      float64           `json:"change_value"`
	ChangePercent    float64           `json:"change_percent"`
	AvgPrice         *float64          `json:"avg_price,omitempty"`
	AvgPricePerSqft  *float64          `json:"avg_price_per_sqft,omitempty"`
	TransactionCount int               `json:"transaction_count"`
	Period           string            `json:"period"`
	Year             int               `json:"year"`
	Month            int               `json:"month"`
	Day              *int              `json:"day,omitempty"`
	DataSource       string            `json:"data_source"`
	Notes            *string           `json:"notes,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// PriceIndexListItemResponse 楼价指数列表项响应
type PriceIndexListItemResponse struct {
	ID               uint     `json:"id"`
	IndexType        string   `json:"index_type"`
	IndexValue       float64  `json:"index_value"`
	ChangeValue      float64  `json:"change_value"`
	ChangePercent    float64  `json:"change_percent"`
	AvgPrice         *float64 `json:"avg_price,omitempty"`
	AvgPricePerSqft  *float64 `json:"avg_price_per_sqft,omitempty"`
	TransactionCount int      `json:"transaction_count"`
	Period           string   `json:"period"`
	Year             int      `json:"year"`
	Month            int      `json:"month"`
}

// LatestPriceIndexResponse 最新楼价指数响应
type LatestPriceIndexResponse struct {
	Overall      *PriceIndexResponse   `json:"overall"`       // 整体指数
	ByDistrict   []PriceIndexResponse  `json:"by_district"`   // 各地区指数
	ByPropertyType []PriceIndexResponse `json:"by_property_type"` // 各物业类型指数
	UpdatedAt    time.Time             `json:"updated_at"`    // 更新时间
}

// PriceTrendResponse 价格走势响应
type PriceTrendResponse struct {
	IndexType    string                `json:"index_type"`
	DistrictID   *uint                 `json:"district_id,omitempty"`
	DistrictName *string               `json:"district_name,omitempty"`
	EstateID     *uint                 `json:"estate_id,omitempty"`
	EstateName   *string               `json:"estate_name,omitempty"`
	PropertyType *string               `json:"property_type,omitempty"`
	StartPeriod  string                `json:"start_period"`
	EndPeriod    string                `json:"end_period"`
	DataPoints   []PriceTrendDataPoint `json:"data_points"`
	Statistics   *TrendStatistics      `json:"statistics"`
}

// PriceTrendDataPoint 价格走势数据点
type PriceTrendDataPoint struct {
	Period           string   `json:"period"`
	IndexValue       float64  `json:"index_value"`
	ChangeValue      float64  `json:"change_value"`
	ChangePercent    float64  `json:"change_percent"`
	AvgPrice         *float64 `json:"avg_price,omitempty"`
	AvgPricePerSqft  *float64 `json:"avg_price_per_sqft,omitempty"`
	TransactionCount int      `json:"transaction_count"`
}

// TrendStatistics 走势统计信息
type TrendStatistics struct {
	HighestValue      float64 `json:"highest_value"`
	LowestValue       float64 `json:"lowest_value"`
	AverageValue      float64 `json:"average_value"`
	TotalChange       float64 `json:"total_change"`
	TotalChangePercent float64 `json:"total_change_percent"`
	VolatilityRate    float64 `json:"volatility_rate"`
}

// ComparePriceIndexResponse 对比楼价指数响应
type ComparePriceIndexResponse struct {
	CompareType string                   `json:"compare_type"`
	StartPeriod string                   `json:"start_period"`
	EndPeriod   string                   `json:"end_period"`
	Series      []CompareSeriesData      `json:"series"`
}

// CompareSeriesData 对比数据系列
type CompareSeriesData struct {
	ID           uint                    `json:"id"`
	Name         string                  `json:"name"`
	Type         string                  `json:"type"` // district, estate, property_type
	DataPoints   []PriceTrendDataPoint   `json:"data_points"`
}

// ExportPriceDataResponse 导出价格数据响应
type ExportPriceDataResponse struct {
	FileName    string `json:"file_name"`
	DownloadURL string `json:"download_url"`
	Format      string `json:"format"`
	RecordCount int    `json:"record_count"`
	ExportedAt  time.Time `json:"exported_at"`
}

// PriceIndexHistoryResponse 历史楼价指数响应
type PriceIndexHistoryResponse struct {
	IndexType    string                `json:"index_type"`
	DistrictID   *uint                 `json:"district_id,omitempty"`
	DistrictName *string               `json:"district_name,omitempty"`
	EstateID     *uint                 `json:"estate_id,omitempty"`
	EstateName   *string               `json:"estate_name,omitempty"`
	PropertyType *string               `json:"property_type,omitempty"`
	Years        int                   `json:"years"`
	DataPoints   []PriceTrendDataPoint `json:"data_points"`
	YearlyStats  []YearlyStatistics    `json:"yearly_stats"`
}

// YearlyStatistics 年度统计
type YearlyStatistics struct {
	Year              int     `json:"year"`
	AverageValue      float64 `json:"average_value"`
	YearStartValue    float64 `json:"year_start_value"`
	YearEndValue      float64 `json:"year_end_value"`
	YearChange        float64 `json:"year_change"`
	YearChangePercent float64 `json:"year_change_percent"`
	HighestValue      float64 `json:"highest_value"`
	LowestValue       float64 `json:"lowest_value"`
	TotalTransactions int     `json:"total_transactions"`
}

// EstateBasicInfo 屋苑基本信息（用于楼价指数）
type EstateBasicInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	NameEn *string `json:"name_en,omitempty"`
}

// CreatePriceIndexResponse 创建楼价指数响应
type CreatePriceIndexResponse struct {
	ID      uint   `json:"id"`
	Period  string `json:"period"`
	Message string `json:"message"`
}

// UpdatePriceIndexResponse 更新楼价指数响应
type UpdatePriceIndexResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
