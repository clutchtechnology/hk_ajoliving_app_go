package models

// statistics.go - 统计相关的响应类型定义

// ============ 请求类型 ============

// GetOverviewStatisticsRequest 获取总览统计请求
type GetOverviewStatisticsRequest struct {
	Period string `form:"period,default=month" binding:"omitempty,oneof=day week month year"`
}

// ============ 总览统计响应 ============

// PropertyOverview 房产总览统计
type PropertyOverview struct {
	TotalCount    int     `json:"total_count"`
	ActiveCount   int     `json:"active_count"`
	RentCount     int     `json:"rent_count"`
	SaleCount     int     `json:"sale_count"`
	NewToday      int     `json:"new_today"`
	NewThisWeek   int     `json:"new_this_week"`
	NewThisMonth  int     `json:"new_this_month"`
	AveragePrice  float64 `json:"average_price"`
	TotalValue    float64 `json:"total_value"`
}

// UserOverview 用户总览统计
type UserOverview struct {
	TotalCount    int `json:"total_count"`
	ActiveCount   int `json:"active_count"`
	NewToday      int `json:"new_today"`
	NewThisWeek   int `json:"new_this_week"`
	NewThisMonth  int `json:"new_this_month"`
	VerifiedCount int `json:"verified_count"`
}

// AgentOverview 代理人总览统计
type AgentOverview struct {
	TotalCount    int     `json:"total_count"`
	ActiveCount   int     `json:"active_count"`
	AverageRating float64 `json:"average_rating"`
	TopAgents     int     `json:"top_agents"`
}

// TransactionOverview 成交总览统计
type TransactionOverview struct {
	TotalCount      int     `json:"total_count"`
	TodayCount      int     `json:"today_count"`
	ThisWeekCount   int     `json:"this_week_count"`
	ThisMonthCount  int     `json:"this_month_count"`
	TotalAmount     float64 `json:"total_amount"`
	AverageAmount   float64 `json:"average_amount"`
}

// PlatformMetrics 平台指标
type PlatformMetrics struct {
	TotalViews     int     `json:"total_views"`
	TodayViews     int     `json:"today_views"`
	SearchCount    int     `json:"search_count"`
	ConversionRate float64 `json:"conversion_rate"`
}

// OverviewStatisticsResponse 总览统计响应
type OverviewStatisticsResponse struct {
	Properties      PropertyOverview    `json:"properties"`
	Users           UserOverview        `json:"users"`
	Agents          AgentOverview       `json:"agents"`
	Transactions    TransactionOverview `json:"transactions"`
	PlatformMetrics PlatformMetrics     `json:"platform_metrics"`
}

// ============ 房产统计响应 ============

// PropertyStatisticsSummary 房产统计汇总
type PropertyStatisticsSummary struct {
	TotalCount   int     `json:"total_count"`
	RentCount    int     `json:"rent_count"`
	SaleCount    int     `json:"sale_count"`
	AveragePrice float64 `json:"average_price"`
	MedianPrice  float64 `json:"median_price"`
	HighestPrice float64 `json:"highest_price"`
	LowestPrice  float64 `json:"lowest_price"`
	TotalValue   float64 `json:"total_value"`
}

// TrendItem 趋势数据项（通用）
type TrendItem struct {
	Period string  `json:"period"` // 日期格式：2006-01-02 或 2006-01
	Count  int     `json:"count"`
	Value  float64 `json:"value,omitempty"` // 可选的数值（如平均价格）
}

// PropertyDistribution 房产分布统计
type PropertyDistribution struct {
	ByDistrict     []*DistrictStatItem      `json:"by_district"`
	ByPropertyType []*PropertyTypeStatItem  `json:"by_property_type"`
	ByPriceRange   []*PriceRangeStatItem    `json:"by_price_range"`
	ByBedroomCount []*BedroomCountStatItem  `json:"by_bedroom_count"`
}

// DistrictStatItem 地区统计项
type DistrictStatItem struct {
	DistrictID   uint    `json:"district_id"`
	DistrictName string  `json:"district_name"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
	AveragePrice float64 `json:"average_price,omitempty"`
}

// PropertyTypeStatItem 房产类型统计项
type PropertyTypeStatItem struct {
	PropertyType string  `json:"property_type"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
}

// PriceRangeStatItem 价格区间统计项
type PriceRangeStatItem struct {
	Range      string  `json:"range"`       // 如 "0-1M", "1M-2M"
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// BedroomCountStatItem 睡房数量统计项
type BedroomCountStatItem struct {
	Bedrooms   int     `json:"bedrooms"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// PropertyStatisticsResponse 房产统计响应
type PropertyStatisticsResponse struct {
	Summary      *PropertyStatisticsSummary `json:"summary"`
	TrendData    []*TrendItem               `json:"trend_data"`
	Distribution *PropertyDistribution      `json:"distribution"`
}

// ============ 交易统计响应 ============

// TransactionStatisticsSummary 交易统计汇总
type TransactionStatisticsSummary struct {
	TotalCount          int     `json:"total_count"`
	TotalAmount         float64 `json:"total_amount"`
	AverageAmount       float64 `json:"average_amount"`
	MedianAmount        float64 `json:"median_amount"`
	HighestAmount       float64 `json:"highest_amount"`
	LowestAmount        float64 `json:"lowest_amount"`
	AveragePricePerSqft float64 `json:"average_price_per_sqft"`
}

// TransactionTrendItem 交易趋势数据项
type TransactionTrendItem struct {
	Period        string  `json:"period"` // 日期格式：2006-01-02
	Count         int     `json:"count"`
	AverageAmount float64 `json:"average_amount"`
}

// TransactionDistribution 交易分布统计
type TransactionDistribution struct {
	ByDistrict []*DistrictTransactionItem `json:"by_district"`
	ByEstate   []*EstateTransactionItem   `json:"by_estate"`
	ByMonth    []*MonthTransactionItem    `json:"by_month"`
}

// DistrictTransactionItem 地区交易统计项
type DistrictTransactionItem struct {
	DistrictID    uint    `json:"district_id"`
	DistrictName  string  `json:"district_name"`
	Count         int     `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// EstateTransactionItem 屋苑交易统计项
type EstateTransactionItem struct {
	EstateID      uint    `json:"estate_id"`
	EstateName    string  `json:"estate_name"`
	Count         int     `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// MonthTransactionItem 月度交易统计项
type MonthTransactionItem struct {
	Month         string  `json:"month"` // 格式：2006-01
	Count         int     `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// TransactionStatisticsResponse 交易统计响应
type TransactionStatisticsResponse struct {
	Summary      *TransactionStatisticsSummary `json:"summary"`
	TrendData    []*TransactionTrendItem       `json:"trend_data"`
	Distribution *TransactionDistribution      `json:"distribution"`
}

// ============ 用户统计响应 ============

// UserStatisticsSummary 用户统计汇总
type UserStatisticsSummary struct {
	TotalCount        int     `json:"total_count"`
	ActiveCount       int     `json:"active_count"`
	NewUsersToday     int     `json:"new_users_today"`
	NewUsersThisWeek  int     `json:"new_users_this_week"`
	NewUsersThisMonth int     `json:"new_users_this_month"`
	VerifiedCount     int     `json:"verified_count"`
	RetentionRate     float64 `json:"retention_rate"` // 留存率
}

// UserDistribution 用户分布统计
type UserDistribution struct {
	ByRole               []*UserRoleStatItem               `json:"by_role"`
	ByStatus             []*UserStatusStatItem             `json:"by_status"`
	ByRegistrationSource []*RegistrationSourceStatItem     `json:"by_registration_source"`
}

// UserRoleStatItem 用户角色统计项
type UserRoleStatItem struct {
	Role       string  `json:"role"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// UserStatusStatItem 用户状态统计项
type UserStatusStatItem struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// RegistrationSourceStatItem 注册来源统计项
type RegistrationSourceStatItem struct {
	Source     string  `json:"source"` // google, facebook, email, phone
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// UserStatisticsResponse 用户统计响应
type UserStatisticsResponse struct {
	Summary      *UserStatisticsSummary `json:"summary"`
	TrendData    []*TrendItem           `json:"trend_data"`
	Distribution *UserDistribution      `json:"distribution"`
}
