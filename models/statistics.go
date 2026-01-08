package models

import "time"

// ============ Request DTO ============

// GetOverviewStatisticsRequest 获取总览统计请求
type GetOverviewStatisticsRequest struct {
	StartDate *string `form:"start_date"` // 开始日期 (YYYY-MM-DD)
	EndDate   *string `form:"end_date"`   // 结束日期 (YYYY-MM-DD)
}

// GetPropertyStatisticsRequest 获取房产统计请求
type GetPropertyStatisticsRequest struct {
	StartDate  *string `form:"start_date"`  // 开始日期 (YYYY-MM-DD)
	EndDate    *string `form:"end_date"`    // 结束日期 (YYYY-MM-DD)
	DistrictID *uint   `form:"district_id"` // 地区ID
}

// GetTransactionStatisticsRequest 获取成交统计请求
type GetTransactionStatisticsRequest struct {
	StartDate  *string `form:"start_date"`  // 开始日期 (YYYY-MM-DD)
	EndDate    *string `form:"end_date"`    // 结束日期 (YYYY-MM-DD)
	DistrictID *uint   `form:"district_id"` // 地区ID
}

// GetUserStatisticsRequest 获取用户统计请求
type GetUserStatisticsRequest struct {
	StartDate *string `form:"start_date"` // 开始日期 (YYYY-MM-DD)
	EndDate   *string `form:"end_date"`   // 结束日期 (YYYY-MM-DD)
}

// ============ Response DTO ============

// OverviewStatisticsResponse 总览统计响应
type OverviewStatisticsResponse struct {
	// 房产统计
	TotalProperties       int64 `json:"total_properties"`        // 总房产数
	SaleProperties        int64 `json:"sale_properties"`         // 出售房产数
	RentProperties        int64 `json:"rent_properties"`         // 出租房产数
	AvailableProperties   int64 `json:"available_properties"`    // 可用房产数
	NewPropertiesThisWeek int64 `json:"new_properties_this_week"` // 本周新增房产

	// 用户统计
	TotalUsers        int64 `json:"total_users"`         // 总用户数
	IndividualUsers   int64 `json:"individual_users"`    // 普通用户数
	AgencyUsers       int64 `json:"agency_users"`        // 代理公司数
	NewUsersThisMonth int64 `json:"new_users_this_month"` // 本月新增用户

	// 代理统计
	TotalAgents     int64 `json:"total_agents"`      // 总代理人数
	VerifiedAgents  int64 `json:"verified_agents"`   // 已验证代理数
	ActiveAgents    int64 `json:"active_agents"`     // 活跃代理数

	// 屋苑统计
	TotalEstates    int64 `json:"total_estates"`     // 总屋苑数

	// 家具统计
	TotalFurniture      int64 `json:"total_furniture"`       // 总家具数
	AvailableFurniture  int64 `json:"available_furniture"`   // 可用家具数

	// 时间范围
	StartDate *string `json:"start_date,omitempty"` // 统计开始日期
	EndDate   *string `json:"end_date,omitempty"`   // 统计结束日期
}

// PropertyStatisticsResponse 房产统计响应
type PropertyStatisticsResponse struct {
	// 房产数量统计
	TotalProperties     int64 `json:"total_properties"`      // 总房产数
	SaleProperties      int64 `json:"sale_properties"`       // 出售房产数
	RentProperties      int64 `json:"rent_properties"`       // 出租房产数
	AvailableProperties int64 `json:"available_properties"`  // 可用房产数
	PendingProperties   int64 `json:"pending_properties"`    // 待定房产数
	SoldProperties      int64 `json:"sold_properties"`       // 已售/已租房产数

	// 价格统计
	AvgSalePrice       float64 `json:"avg_sale_price"`        // 平均售价
	AvgRentPrice       float64 `json:"avg_rent_price"`        // 平均租金
	MaxSalePrice       float64 `json:"max_sale_price"`        // 最高售价
	MinSalePrice       float64 `json:"min_sale_price"`        // 最低售价
	MaxRentPrice       float64 `json:"max_rent_price"`        // 最高租金
	MinRentPrice       float64 `json:"min_rent_price"`        // 最低租金

	// 面积统计
	AvgArea float64 `json:"avg_area"` // 平均面积
	MaxArea float64 `json:"max_area"` // 最大面积
	MinArea float64 `json:"min_area"` // 最小面积

	// 房型统计
	BedroomDistribution []BedroomStat `json:"bedroom_distribution"` // 房间数分布

	// 物业类型统计
	PropertyTypeDistribution []PropertyTypeStat `json:"property_type_distribution"` // 物业类型分布

	// 地区统计
	DistrictDistribution []DistrictStat `json:"district_distribution"` // 地区分布

	// 时间趋势
	NewPropertiesThisWeek  int64 `json:"new_properties_this_week"`  // 本周新增
	NewPropertiesThisMonth int64 `json:"new_properties_this_month"` // 本月新增
	NewPropertiesThisYear  int64 `json:"new_properties_this_year"`  // 本年新增

	// 浏览与收藏
	TotalViews     int64 `json:"total_views"`     // 总浏览次数
	TotalFavorites int64 `json:"total_favorites"` // 总收藏次数
	AvgViews       float64 `json:"avg_views"`       // 平均浏览次数
	AvgFavorites   float64 `json:"avg_favorites"`   // 平均收藏次数

	// 时间范围
	StartDate  *string `json:"start_date,omitempty"`  // 统计开始日期
	EndDate    *string `json:"end_date,omitempty"`    // 统计结束日期
	DistrictID *uint   `json:"district_id,omitempty"` // 地区ID
}

// BedroomStat 房间数统计
type BedroomStat struct {
	Bedrooms int   `json:"bedrooms"` // 房间数
	Count    int64 `json:"count"`    // 数量
}

// PropertyTypeStat 物业类型统计
type PropertyTypeStat struct {
	PropertyType string `json:"property_type"` // 物业类型
	Count        int64  `json:"count"`         // 数量
}

// DistrictStat 地区统计
type DistrictStat struct {
	DistrictID   uint   `json:"district_id"`   // 地区ID
	DistrictName string `json:"district_name"` // 地区名称
	Count        int64  `json:"count"`         // 数量
}

// TransactionStatisticsResponse 成交统计响应
type TransactionStatisticsResponse struct {
	// 成交数量统计
	TotalTransactions      int64 `json:"total_transactions"`       // 总成交数
	SaleTransactions       int64 `json:"sale_transactions"`        // 买卖成交数
	RentTransactions       int64 `json:"rent_transactions"`        // 租赁成交数
	TransactionsThisWeek   int64 `json:"transactions_this_week"`   // 本周成交数
	TransactionsThisMonth  int64 `json:"transactions_this_month"`  // 本月成交数
	TransactionsThisYear   int64 `json:"transactions_this_year"`   // 本年成交数

	// 成交金额统计
	TotalTransactionValue float64 `json:"total_transaction_value"` // 总成交金额
	AvgTransactionPrice   float64 `json:"avg_transaction_price"`   // 平均成交价
	MaxTransactionPrice   float64 `json:"max_transaction_price"`   // 最高成交价
	MinTransactionPrice   float64 `json:"min_transaction_price"`   // 最低成交价

	// 地区成交统计
	DistrictTransactions []DistrictTransactionStat `json:"district_transactions"` // 地区成交统计

	// 物业类型成交统计
	PropertyTypeTransactions []PropertyTypeTransactionStat `json:"property_type_transactions"` // 物业类型成交统计

	// 时间趋势 (最近12个月)
	MonthlyTrend []MonthlyTransactionStat `json:"monthly_trend"` // 月度趋势

	// 时间范围
	StartDate  *string `json:"start_date,omitempty"`  // 统计开始日期
	EndDate    *string `json:"end_date,omitempty"`    // 统计结束日期
	DistrictID *uint   `json:"district_id,omitempty"` // 地区ID
}

// DistrictTransactionStat 地区成交统计
type DistrictTransactionStat struct {
	DistrictID          uint    `json:"district_id"`           // 地区ID
	DistrictName        string  `json:"district_name"`         // 地区名称
	TransactionCount    int64   `json:"transaction_count"`     // 成交数量
	TotalValue          float64 `json:"total_value"`           // 总成交金额
	AvgPrice            float64 `json:"avg_price"`             // 平均价格
}

// PropertyTypeTransactionStat 物业类型成交统计
type PropertyTypeTransactionStat struct {
	PropertyType     string  `json:"property_type"`      // 物业类型
	TransactionCount int64   `json:"transaction_count"`  // 成交数量
	TotalValue       float64 `json:"total_value"`        // 总成交金额
	AvgPrice         float64 `json:"avg_price"`          // 平均价格
}

// MonthlyTransactionStat 月度成交统计
type MonthlyTransactionStat struct {
	Month            string  `json:"month"`             // 月份 (YYYY-MM)
	TransactionCount int64   `json:"transaction_count"` // 成交数量
	TotalValue       float64 `json:"total_value"`       // 总成交金额
	AvgPrice         float64 `json:"avg_price"`         // 平均价格
}

// UserStatisticsResponse 用户统计响应
type UserStatisticsResponse struct {
	// 用户总数统计
	TotalUsers      int64 `json:"total_users"`       // 总用户数
	IndividualUsers int64 `json:"individual_users"`  // 普通用户数
	AgencyUsers     int64 `json:"agency_users"`      // 代理公司数
	ActiveUsers     int64 `json:"active_users"`      // 活跃用户数
	InactiveUsers   int64 `json:"inactive_users"`    // 停用用户数
	SuspendedUsers  int64 `json:"suspended_users"`   // 暂停用户数

	// 新增用户统计
	NewUsersThisWeek  int64 `json:"new_users_this_week"`  // 本周新增
	NewUsersThisMonth int64 `json:"new_users_this_month"` // 本月新增
	NewUsersThisYear  int64 `json:"new_users_this_year"`  // 本年新增

	// 用户活跃度
	VerifiedEmailUsers int64 `json:"verified_email_users"` // 已验证邮箱用户数
	UsersWithListings  int64 `json:"users_with_listings"`  // 有发布记录的用户数

	// 代理统计
	TotalAgents    int64 `json:"total_agents"`     // 总代理人数
	VerifiedAgents int64 `json:"verified_agents"`  // 已验证代理数
	ActiveAgents   int64 `json:"active_agents"`    // 活跃代理数

	// 用户登录统计
	UsersLoggedInToday     int64 `json:"users_logged_in_today"`      // 今日登录用户数
	UsersLoggedInThisWeek  int64 `json:"users_logged_in_this_week"`  // 本周登录用户数
	UsersLoggedInThisMonth int64 `json:"users_logged_in_this_month"` // 本月登录用户数

	// 时间趋势 (最近12个月)
	MonthlyUserGrowth []MonthlyUserGrowthStat `json:"monthly_user_growth"` // 月度用户增长

	// 时间范围
	StartDate *string `json:"start_date,omitempty"` // 统计开始日期
	EndDate   *string `json:"end_date,omitempty"`   // 统计结束日期
}

// MonthlyUserGrowthStat 月度用户增长统计
type MonthlyUserGrowthStat struct {
	Month          string `json:"month"`           // 月份 (YYYY-MM)
	NewUsers       int64  `json:"new_users"`       // 新增用户数
	CumulativeUsers int64 `json:"cumulative_users"` // 累计用户数
}

// TimeRange 时间范围辅助结构
type TimeRange struct {
	StartDate time.Time
	EndDate   time.Time
}
