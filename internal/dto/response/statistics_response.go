package response

// OverviewStatisticsResponse 总览统计响应
type OverviewStatisticsResponse struct {
	Properties  *PropertyOverview  `json:"properties"`
	Users       *UserOverview      `json:"users"`
	Agents      *AgentOverview     `json:"agents"`
	Transactions *TransactionOverview `json:"transactions"`
	PlatformMetrics *PlatformMetrics `json:"platform_metrics"`
}

// PropertyOverview 房产概览
type PropertyOverview struct {
	TotalCount      int     `json:"total_count"`
	ActiveCount     int     `json:"active_count"`
	RentCount       int     `json:"rent_count"`
	SaleCount       int     `json:"sale_count"`
	NewToday        int     `json:"new_today"`
	NewThisWeek     int     `json:"new_this_week"`
	NewThisMonth    int     `json:"new_this_month"`
	AveragePrice    float64 `json:"average_price"`
	TotalValue      float64 `json:"total_value"`
}

// UserOverview 用户概览
type UserOverview struct {
	TotalCount      int `json:"total_count"`
	ActiveCount     int `json:"active_count"`
	NewToday        int `json:"new_today"`
	NewThisWeek     int `json:"new_this_week"`
	NewThisMonth    int `json:"new_this_month"`
	VerifiedCount   int `json:"verified_count"`
}

// AgentOverview 代理人概览
type AgentOverview struct {
	TotalCount      int     `json:"total_count"`
	ActiveCount     int     `json:"active_count"`
	AverageRating   float64 `json:"average_rating"`
	TopAgents       int     `json:"top_agents"`
}

// TransactionOverview 成交概览
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
	TotalViews      int     `json:"total_views"`
	TodayViews      int     `json:"today_views"`
	SearchCount     int     `json:"search_count"`
	ConversionRate  float64 `json:"conversion_rate"`
}

// PropertyStatisticsResponse 房产统计响应
type PropertyStatisticsResponse struct {
	Summary      *PropertyStatisticsSummary `json:"summary"`
	TrendData    []*StatisticsTrendItem     `json:"trend_data"`
	Distribution *PropertyDistribution      `json:"distribution"`
}

// PropertyStatisticsSummary 房产统计汇总
type PropertyStatisticsSummary struct {
	TotalCount      int     `json:"total_count"`
	RentCount       int     `json:"rent_count"`
	SaleCount       int     `json:"sale_count"`
	AveragePrice    float64 `json:"average_price"`
	MedianPrice     float64 `json:"median_price"`
	HighestPrice    float64 `json:"highest_price"`
	LowestPrice     float64 `json:"lowest_price"`
	TotalValue      float64 `json:"total_value"`
}

// StatisticsTrendItem 统计趋势项
type StatisticsTrendItem struct {
	Period string  `json:"period"` // 时间周期
	Count  int     `json:"count"`  // 数量
	Value  float64 `json:"value"`  // 值（价格/金额）
}

// PropertyDistribution 房产分布
type PropertyDistribution struct {
	ByDistrict      []*DistrictStatItem      `json:"by_district"`
	ByPropertyType  []*PropertyTypeStatItem  `json:"by_property_type"`
	ByPriceRange    []*PriceRangeStatItem    `json:"by_price_range"`
	ByBedroomCount  []*BedroomCountStatItem  `json:"by_bedroom_count"`
}

// DistrictStatItem 地区统计项
type DistrictStatItem struct {
	DistrictID   uint    `json:"district_id"`
	DistrictName string  `json:"district_name"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
	AveragePrice float64 `json:"average_price"`
}

// PropertyTypeStatItem 物业类型统计项
type PropertyTypeStatItem struct {
	PropertyType string  `json:"property_type"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
}

// PriceRangeStatItem 价格区间统计项
type PriceRangeStatItem struct {
	Range      string  `json:"range"`       // 例如 "0-1M", "1M-2M"
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// BedroomCountStatItem 房间数统计项
type BedroomCountStatItem struct {
	Bedrooms   int     `json:"bedrooms"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// TransactionStatisticsResponse 成交统计响应
type TransactionStatisticsResponse struct {
	Summary      *TransactionStatisticsSummary `json:"summary"`
	TrendData    []*StatisticsTrendItem        `json:"trend_data"`
	Distribution *TransactionDistribution      `json:"distribution"`
}

// TransactionStatisticsSummary 成交统计汇总
type TransactionStatisticsSummary struct {
	TotalCount       int     `json:"total_count"`
	TotalAmount      float64 `json:"total_amount"`
	AverageAmount    float64 `json:"average_amount"`
	MedianAmount     float64 `json:"median_amount"`
	HighestAmount    float64 `json:"highest_amount"`
	LowestAmount     float64 `json:"lowest_amount"`
	AveragePricePerSqft float64 `json:"average_price_per_sqft"`
}

// TransactionDistribution 成交分布
type TransactionDistribution struct {
	ByDistrict []*DistrictTransactionItem `json:"by_district"`
	ByEstate   []*EstateTransactionItem   `json:"by_estate"`
	ByMonth    []*MonthTransactionItem    `json:"by_month"`
}

// DistrictTransactionItem 地区成交统计项
type DistrictTransactionItem struct {
	DistrictID    uint    `json:"district_id"`
	DistrictName  string  `json:"district_name"`
	Count         int     `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// EstateTransactionItem 屋苑成交统计项
type EstateTransactionItem struct {
	EstateID      uint    `json:"estate_id"`
	EstateName    string  `json:"estate_name"`
	Count         int     `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// MonthTransactionItem 月度成交统计项
type MonthTransactionItem struct {
	Month         string  `json:"month"` // YYYY-MM
	Count         int     `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

// UserStatisticsResponse 用户统计响应
type UserStatisticsResponse struct {
	Summary      *UserStatisticsSummary `json:"summary"`
	TrendData    []*StatisticsTrendItem `json:"trend_data"`
	Distribution *UserDistribution      `json:"distribution"`
}

// UserStatisticsSummary 用户统计汇总
type UserStatisticsSummary struct {
	TotalCount       int     `json:"total_count"`
	ActiveCount      int     `json:"active_count"`
	NewUsersToday    int     `json:"new_users_today"`
	NewUsersThisWeek int     `json:"new_users_this_week"`
	NewUsersThisMonth int    `json:"new_users_this_month"`
	VerifiedCount    int     `json:"verified_count"`
	RetentionRate    float64 `json:"retention_rate"`
}

// UserDistribution 用户分布
type UserDistribution struct {
	ByRole         []*UserRoleStatItem    `json:"by_role"`
	ByStatus       []*UserStatusStatItem  `json:"by_status"`
	ByRegistrationSource []*RegistrationSourceStatItem `json:"by_registration_source"`
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
	Source     string  `json:"source"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}
