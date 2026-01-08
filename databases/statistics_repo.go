package databases

import (
	"context"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// StatisticsRepo 统计数据仓储
type StatisticsRepo struct {
	db *gorm.DB
}

// NewStatisticsRepo 创建统计数据仓储实例
func NewStatisticsRepo(db *gorm.DB) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}

// GetOverviewStatistics 获取总览统计
func (r *StatisticsRepo) GetOverviewStatistics(ctx context.Context, startDate, endDate *time.Time) (*models.OverviewStatisticsResponse, error) {
	var stats models.OverviewStatisticsResponse

	// 房产统计
	r.db.WithContext(ctx).Model(&models.Property{}).Count(&stats.TotalProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("listing_type = ?", "sale").Count(&stats.SaleProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("listing_type = ?", "rent").Count(&stats.RentProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("status = ?", "available").Count(&stats.AvailableProperties)

	// 本周新增房产
	weekAgo := time.Now().AddDate(0, 0, -7)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("created_at >= ?", weekAgo).Count(&stats.NewPropertiesThisWeek)

	// 用户统计
	r.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("user_type = ?", "individual").Count(&stats.IndividualUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("user_type = ?", "agency").Count(&stats.AgencyUsers)

	// 本月新增用户
	monthAgo := time.Now().AddDate(0, -1, 0)
	r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", monthAgo).Count(&stats.NewUsersThisMonth)

	// 代理统计
	r.db.WithContext(ctx).Model(&models.Agent{}).Count(&stats.TotalAgents)
	r.db.WithContext(ctx).Model(&models.Agent{}).Where("is_verified = ?", true).Count(&stats.VerifiedAgents)
	r.db.WithContext(ctx).Model(&models.Agent{}).Where("status = ?", "active").Count(&stats.ActiveAgents)

	// 屋苑统计
	r.db.WithContext(ctx).Model(&models.Estate{}).Count(&stats.TotalEstates)

	// 家具统计
	r.db.WithContext(ctx).Model(&models.Furniture{}).Count(&stats.TotalFurniture)
	r.db.WithContext(ctx).Model(&models.Furniture{}).Where("status = ?", "available").Count(&stats.AvailableFurniture)

	return &stats, nil
}

// GetPropertyStatistics 获取房产统计
func (r *StatisticsRepo) GetPropertyStatistics(ctx context.Context, startDate, endDate *time.Time, districtID *uint) (*models.PropertyStatisticsResponse, error) {
	var stats models.PropertyStatisticsResponse

	query := r.db.WithContext(ctx).Model(&models.Property{})
	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}
	if districtID != nil {
		query = query.Where("district_id = ?", *districtID)
	}

	// 房产数量统计
	query.Count(&stats.TotalProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("listing_type = ?", "sale").Count(&stats.SaleProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("listing_type = ?", "rent").Count(&stats.RentProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("status = ?", "available").Count(&stats.AvailableProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("status = ?", "pending").Count(&stats.PendingProperties)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("status IN ?", []string{"sold", "rented"}).Count(&stats.SoldProperties)

	// 价格统计
	var saleStats struct {
		AvgPrice float64
		MaxPrice float64
		MinPrice float64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("listing_type = ?", "sale").
		Select("AVG(price) as avg_price, MAX(price) as max_price, MIN(price) as min_price").
		Scan(&saleStats)
	stats.AvgSalePrice = saleStats.AvgPrice
	stats.MaxSalePrice = saleStats.MaxPrice
	stats.MinSalePrice = saleStats.MinPrice

	var rentStats struct {
		AvgPrice float64
		MaxPrice float64
		MinPrice float64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("listing_type = ?", "rent").
		Select("AVG(price) as avg_price, MAX(price) as max_price, MIN(price) as min_price").
		Scan(&rentStats)
	stats.AvgRentPrice = rentStats.AvgPrice
	stats.MaxRentPrice = rentStats.MaxPrice
	stats.MinRentPrice = rentStats.MinPrice

	// 面积统计
	var areaStats struct {
		AvgArea float64
		MaxArea float64
		MinArea float64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("AVG(area) as avg_area, MAX(area) as max_area, MIN(area) as min_area").
		Scan(&areaStats)
	stats.AvgArea = areaStats.AvgArea
	stats.MaxArea = areaStats.MaxArea
	stats.MinArea = areaStats.MinArea

	// 房间数分布
	var bedroomDist []models.BedroomStat
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("bedrooms, COUNT(*) as count").
		Group("bedrooms").
		Order("bedrooms").
		Scan(&bedroomDist)
	stats.BedroomDistribution = bedroomDist

	// 物业类型分布
	var propertyTypeDist []models.PropertyTypeStat
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("property_type, COUNT(*) as count").
		Group("property_type").
		Order("count DESC").
		Scan(&propertyTypeDist)
	stats.PropertyTypeDistribution = propertyTypeDist

	// 地区分布
	var districtDist []struct {
		DistrictID   uint
		DistrictName string
		Count        int64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("properties.district_id, districts.name_zh as district_name, COUNT(*) as count").
		Joins("LEFT JOIN districts ON properties.district_id = districts.id").
		Group("properties.district_id, districts.name_zh").
		Order("count DESC").
		Limit(10).
		Scan(&districtDist)

	for _, d := range districtDist {
		stats.DistrictDistribution = append(stats.DistrictDistribution, models.DistrictStat{
			DistrictID:   d.DistrictID,
			DistrictName: d.DistrictName,
			Count:        d.Count,
		})
	}

	// 时间趋势
	weekAgo := time.Now().AddDate(0, 0, -7)
	monthAgo := time.Now().AddDate(0, -1, 0)
	yearAgo := time.Now().AddDate(-1, 0, 0)

	r.db.WithContext(ctx).Model(&models.Property{}).Where("created_at >= ?", weekAgo).Count(&stats.NewPropertiesThisWeek)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("created_at >= ?", monthAgo).Count(&stats.NewPropertiesThisMonth)
	r.db.WithContext(ctx).Model(&models.Property{}).Where("created_at >= ?", yearAgo).Count(&stats.NewPropertiesThisYear)

	// 浏览与收藏统计
	var viewStats struct {
		TotalViews     int64
		TotalFavorites int64
		AvgViews       float64
		AvgFavorites   float64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("SUM(view_count) as total_views, SUM(favorite_count) as total_favorites, AVG(view_count) as avg_views, AVG(favorite_count) as avg_favorites").
		Scan(&viewStats)
	stats.TotalViews = viewStats.TotalViews
	stats.TotalFavorites = viewStats.TotalFavorites
	stats.AvgViews = viewStats.AvgViews
	stats.AvgFavorites = viewStats.AvgFavorites

	return &stats, nil
}

// GetTransactionStatistics 获取成交统计
func (r *StatisticsRepo) GetTransactionStatistics(ctx context.Context, startDate, endDate *time.Time, districtID *uint) (*models.TransactionStatisticsResponse, error) {
	var stats models.TransactionStatisticsResponse

	// Note: 由于数据库中没有 transactions 表，这里使用 properties 表的 sold/rented 状态模拟
	// 实际应用中需要创建独立的 transactions 表

	query := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("status IN ?", []string{"sold", "rented"})

	if startDate != nil {
		query = query.Where("updated_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("updated_at <= ?", endDate)
	}
	if districtID != nil {
		query = query.Where("district_id = ?", *districtID)
	}

	// 成交数量统计
	query.Count(&stats.TotalTransactions)

	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("listing_type = ? AND status = ?", "sale", "sold").
		Count(&stats.SaleTransactions)

	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("listing_type = ? AND status = ?", "rent", "rented").
		Count(&stats.RentTransactions)

	// 时间趋势
	weekAgo := time.Now().AddDate(0, 0, -7)
	monthAgo := time.Now().AddDate(0, -1, 0)
	yearAgo := time.Now().AddDate(-1, 0, 0)

	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("status IN ? AND updated_at >= ?", []string{"sold", "rented"}, weekAgo).
		Count(&stats.TransactionsThisWeek)

	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("status IN ? AND updated_at >= ?", []string{"sold", "rented"}, monthAgo).
		Count(&stats.TransactionsThisMonth)

	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("status IN ? AND updated_at >= ?", []string{"sold", "rented"}, yearAgo).
		Count(&stats.TransactionsThisYear)

	// 成交金额统计
	var priceStats struct {
		TotalValue float64
		AvgPrice   float64
		MaxPrice   float64
		MinPrice   float64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("status IN ?", []string{"sold", "rented"}).
		Select("SUM(price) as total_value, AVG(price) as avg_price, MAX(price) as max_price, MIN(price) as min_price").
		Scan(&priceStats)
	stats.TotalTransactionValue = priceStats.TotalValue
	stats.AvgTransactionPrice = priceStats.AvgPrice
	stats.MaxTransactionPrice = priceStats.MaxPrice
	stats.MinTransactionPrice = priceStats.MinPrice

	// 地区成交统计
	var districtTrans []struct {
		DistrictID   uint
		DistrictName string
		Count        int64
		TotalValue   float64
		AvgPrice     float64
	}
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("properties.district_id, districts.name_zh as district_name, COUNT(*) as count, SUM(properties.price) as total_value, AVG(properties.price) as avg_price").
		Joins("LEFT JOIN districts ON properties.district_id = districts.id").
		Where("properties.status IN ?", []string{"sold", "rented"}).
		Group("properties.district_id, districts.name_zh").
		Order("count DESC").
		Limit(10).
		Scan(&districtTrans)

	for _, d := range districtTrans {
		stats.DistrictTransactions = append(stats.DistrictTransactions, models.DistrictTransactionStat{
			DistrictID:       d.DistrictID,
			DistrictName:     d.DistrictName,
			TransactionCount: d.Count,
			TotalValue:       d.TotalValue,
			AvgPrice:         d.AvgPrice,
		})
	}

	// 物业类型成交统计
	var propertyTypeTrans []models.PropertyTypeTransactionStat
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("property_type, COUNT(*) as transaction_count, SUM(price) as total_value, AVG(price) as avg_price").
		Where("status IN ?", []string{"sold", "rented"}).
		Group("property_type").
		Order("transaction_count DESC").
		Scan(&propertyTypeTrans)
	stats.PropertyTypeTransactions = propertyTypeTrans

	// 月度趋势（最近12个月）
	var monthlyTrend []models.MonthlyTransactionStat
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("DATE_FORMAT(updated_at, '%Y-%m') as month, COUNT(*) as transaction_count, SUM(price) as total_value, AVG(price) as avg_price").
		Where("status IN ? AND updated_at >= ?", []string{"sold", "rented"}, time.Now().AddDate(0, -12, 0)).
		Group("DATE_FORMAT(updated_at, '%Y-%m')").
		Order("month DESC").
		Scan(&monthlyTrend)
	stats.MonthlyTrend = monthlyTrend

	return &stats, nil
}

// GetUserStatistics 获取用户统计
func (r *StatisticsRepo) GetUserStatistics(ctx context.Context, startDate, endDate *time.Time) (*models.UserStatisticsResponse, error) {
	var stats models.UserStatisticsResponse

	query := r.db.WithContext(ctx).Model(&models.User{})
	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	// 用户总数统计
	r.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("user_type = ?", "individual").Count(&stats.IndividualUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("user_type = ?", "agency").Count(&stats.AgencyUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("status = ?", "inactive").Count(&stats.InactiveUsers)
	r.db.WithContext(ctx).Model(&models.User{}).Where("status = ?", "suspended").Count(&stats.SuspendedUsers)

	// 新增用户统计
	weekAgo := time.Now().AddDate(0, 0, -7)
	monthAgo := time.Now().AddDate(0, -1, 0)
	yearAgo := time.Now().AddDate(-1, 0, 0)

	r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", weekAgo).Count(&stats.NewUsersThisWeek)
	r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", monthAgo).Count(&stats.NewUsersThisMonth)
	r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", yearAgo).Count(&stats.NewUsersThisYear)

	// 用户活跃度
	r.db.WithContext(ctx).Model(&models.User{}).Where("email_verified = ?", true).Count(&stats.VerifiedEmailUsers)

	// 有发布记录的用户数
	r.db.WithContext(ctx).Model(&models.User{}).
		Joins("INNER JOIN properties ON users.id = properties.publisher_id").
		Distinct("users.id").
		Count(&stats.UsersWithListings)

	// 代理统计
	r.db.WithContext(ctx).Model(&models.Agent{}).Count(&stats.TotalAgents)
	r.db.WithContext(ctx).Model(&models.Agent{}).Where("is_verified = ?", true).Count(&stats.VerifiedAgents)
	r.db.WithContext(ctx).Model(&models.Agent{}).Where("status = ?", "active").Count(&stats.ActiveAgents)

	// 用户登录统计
	today := time.Now().Truncate(24 * time.Hour)
	r.db.WithContext(ctx).Model(&models.User{}).Where("last_login_at >= ?", today).Count(&stats.UsersLoggedInToday)
	r.db.WithContext(ctx).Model(&models.User{}).Where("last_login_at >= ?", weekAgo).Count(&stats.UsersLoggedInThisWeek)
	r.db.WithContext(ctx).Model(&models.User{}).Where("last_login_at >= ?", monthAgo).Count(&stats.UsersLoggedInThisMonth)

	// 月度用户增长（最近12个月）
	var monthlyGrowth []models.MonthlyUserGrowthStat
	r.db.WithContext(ctx).Raw(`
		SELECT 
			DATE_FORMAT(created_at, '%Y-%m') as month,
			COUNT(*) as new_users,
			(SELECT COUNT(*) FROM users u2 WHERE u2.created_at <= LAST_DAY(users.created_at)) as cumulative_users
		FROM users
		WHERE created_at >= ?
		GROUP BY DATE_FORMAT(created_at, '%Y-%m')
		ORDER BY month DESC
	`, time.Now().AddDate(0, -12, 0)).Scan(&monthlyGrowth)
	stats.MonthlyUserGrowth = monthlyGrowth

	return &stats, nil
}
