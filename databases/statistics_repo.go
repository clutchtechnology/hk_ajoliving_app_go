package databases

import (
	"context"
	"time"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// StatisticsRepository 统计数据仓库接口
type StatisticsRepository interface {
	// 房产统计
	GetPropertyCount(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (int64, error)
	GetPropertyCountByStatus(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (map[string]int64, error)
	GetPropertyCountByListingType(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (map[string]int64, error)
	GetPropertyPriceStats(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (map[string]float64, error)
	GetPropertyCountByDistrict(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error)
	GetPropertyCountByType(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error)
	GetPropertyCountByBedrooms(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error)
	GetPropertyTrendData(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error)

	// 成交统计
	GetTransactionCount(ctx context.Context, filter *models.GetTransactionStatisticsRequest) (int64, error)
	GetTransactionAmountStats(ctx context.Context, filter *models.GetTransactionStatisticsRequest) (map[string]float64, error)
	GetTransactionCountByDistrict(ctx context.Context, filter *models.GetTransactionStatisticsRequest) ([]map[string]interface{}, error)
	GetTransactionCountByEstate(ctx context.Context, filter *models.GetTransactionStatisticsRequest) ([]map[string]interface{}, error)
	GetTransactionTrendData(ctx context.Context, filter *models.GetTransactionStatisticsRequest) ([]map[string]interface{}, error)

	// 用户统计
	GetUserCount(ctx context.Context, filter *models.GetUserStatisticsRequest) (int64, error)
	GetUserCountByStatus(ctx context.Context, filter *models.GetUserStatisticsRequest) (map[string]int64, error)
	GetUserCountByRole(ctx context.Context, filter *models.GetUserStatisticsRequest) (map[string]int64, error)
	GetUserTrendData(ctx context.Context, filter *models.GetUserStatisticsRequest) ([]map[string]interface{}, error)

	// 代理人统计
	GetAgentCount(ctx context.Context) (int64, error)
	GetAgentAverageRating(ctx context.Context) (float64, error)
}

type statisticsRepository struct {
	db *gorm.DB
}

// NewStatisticsRepository 创建统计数据仓库实例
func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

// ========== 房产统计 ==========

func (r *statisticsRepository) GetPropertyCount(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *statisticsRepository) GetPropertyCountByStatus(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (map[string]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)
	if err := query.Select("status, COUNT(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, r := range results {
		statusMap[r.Status] = r.Count
	}
	return statusMap, nil
}

func (r *statisticsRepository) GetPropertyCountByListingType(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (map[string]int64, error) {
	type Result struct {
		ListingType string
		Count       int64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)
	if err := query.Select("listing_type, COUNT(*) as count").Group("listing_type").Scan(&results).Error; err != nil {
		return nil, err
	}

	typeMap := make(map[string]int64)
	for _, r := range results {
		typeMap[r.ListingType] = r.Count
	}
	return typeMap, nil
}

func (r *statisticsRepository) GetPropertyPriceStats(ctx context.Context, filter *models.GetPropertyStatisticsRequest) (map[string]float64, error) {
	type Result struct {
		Total   float64
		Average float64
		Max     float64
		Min     float64
	}
	var result Result
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)
	if err := query.Select("SUM(price) as total, AVG(price) as average, MAX(price) as max, MIN(price) as min").Scan(&result).Error; err != nil {
		return nil, err
	}

	return map[string]float64{
		"total":   result.Total,
		"average": result.Average,
		"max":     result.Max,
		"min":     result.Min,
	}, nil
}

func (r *statisticsRepository) GetPropertyCountByDistrict(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error) {
	type Result struct {
		DistrictID   uint
		DistrictName string
		Count        int64
		AveragePrice float64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.Property{}).
		Select("properties.district_id, districts.name_zh_hant as district_name, COUNT(*) as count, AVG(properties.price) as average_price").
		Joins("LEFT JOIN districts ON properties.district_id = districts.id")
	query = r.applyPropertyFilters(query, filter)
	if err := query.Group("properties.district_id, districts.name_zh_hant").Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range results {
		output = append(output, map[string]interface{}{
			"district_id":    r.DistrictID,
			"district_name":  r.DistrictName,
			"count":          r.Count,
			"average_price":  r.AveragePrice,
		})
	}
	return output, nil
}

func (r *statisticsRepository) GetPropertyCountByType(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error) {
	type Result struct {
		PropertyType string
		Count        int64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)
	if err := query.Select("property_type, COUNT(*) as count").Group("property_type").Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range results {
		output = append(output, map[string]interface{}{
			"property_type": r.PropertyType,
			"count":         r.Count,
		})
	}
	return output, nil
}

func (r *statisticsRepository) GetPropertyCountByBedrooms(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error) {
	type Result struct {
		Bedrooms int
		Count    int64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)
	if err := query.Select("bedrooms, COUNT(*) as count").Group("bedrooms").Order("bedrooms ASC").Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range results {
		output = append(output, map[string]interface{}{
			"bedrooms": r.Bedrooms,
			"count":    r.Count,
		})
	}
	return output, nil
}

func (r *statisticsRepository) GetPropertyTrendData(ctx context.Context, filter *models.GetPropertyStatisticsRequest) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := r.db.WithContext(ctx).Model(&models.Property{})
	query = r.applyPropertyFilters(query, filter)

	// 根据 period 分组
	var period string
	if filter.Period != nil {
		period = *filter.Period
	}
	switch period {
	case "day":
		query = query.Select("DATE(created_at) as period, COUNT(*) as count, AVG(price) as value").
			Group("DATE(created_at)").
			Order("DATE(created_at) ASC")
	case "week":
		query = query.Select("DATE_TRUNC('week', created_at) as period, COUNT(*) as count, AVG(price) as value").
			Group("DATE_TRUNC('week', created_at)").
			Order("DATE_TRUNC('week', created_at) ASC")
	case "month":
		query = query.Select("DATE_TRUNC('month', created_at) as period, COUNT(*) as count, AVG(price) as value").
			Group("DATE_TRUNC('month', created_at)").
			Order("DATE_TRUNC('month', created_at) ASC")
	case "year":
		query = query.Select("DATE_TRUNC('year', created_at) as period, COUNT(*) as count, AVG(price) as value").
			Group("DATE_TRUNC('year', created_at)").
			Order("DATE_TRUNC('year', created_at) ASC")
	default: // 默认按月
		query = query.Select("DATE_TRUNC('month', created_at) as period, COUNT(*) as count, AVG(price) as value").
			Group("DATE_TRUNC('month', created_at)").
			Order("DATE_TRUNC('month', created_at) ASC")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *statisticsRepository) applyPropertyFilters(query *gorm.DB, filter *models.GetPropertyStatisticsRequest) *gorm.DB {
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.EstateID != nil {
		query = query.Where("estate_id = ?", *filter.EstateID)
	}
	if filter.ListingType != nil {
		query = query.Where("listing_type = ?", *filter.ListingType)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	return query
}

// ========== 成交统计 ==========

func (r *statisticsRepository) GetTransactionCount(ctx context.Context, filter *models.GetTransactionStatisticsRequest) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Table("properties")
	query = r.applyTransactionFilters(query, filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *statisticsRepository) GetTransactionAmountStats(ctx context.Context, filter *models.GetTransactionStatisticsRequest) (map[string]float64, error) {
	type Result struct {
		Total   float64
		Average float64
		Max     float64
		Min     float64
	}
	var result Result
	query := r.db.WithContext(ctx).Table("properties")
	query = r.applyTransactionFilters(query, filter)
	if err := query.Select("SUM(price) as total, AVG(price) as average, MAX(price) as max, MIN(price) as min").Scan(&result).Error; err != nil {
		return nil, err
	}

	return map[string]float64{
		"total":   result.Total,
		"average": result.Average,
		"max":     result.Max,
		"min":     result.Min,
	}, nil
}

func (r *statisticsRepository) GetTransactionCountByDistrict(ctx context.Context, filter *models.GetTransactionStatisticsRequest) ([]map[string]interface{}, error) {
	type Result struct {
		DistrictID    uint
		DistrictName  string
		Count         int64
		TotalAmount   float64
		AverageAmount float64
	}
	var results []Result
	query := r.db.WithContext(ctx).Table("properties").
		Select("properties.district_id, districts.name_zh_hant as district_name, COUNT(*) as count, SUM(properties.price) as total_amount, AVG(properties.price) as average_amount").
		Joins("LEFT JOIN districts ON properties.district_id = districts.id")
	query = r.applyTransactionFilters(query, filter)
	if err := query.Group("properties.district_id, districts.name_zh_hant").Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range results {
		output = append(output, map[string]interface{}{
			"district_id":    r.DistrictID,
			"district_name":  r.DistrictName,
			"count":          r.Count,
			"total_amount":   r.TotalAmount,
			"average_amount": r.AverageAmount,
		})
	}
	return output, nil
}

func (r *statisticsRepository) GetTransactionCountByEstate(ctx context.Context, filter *models.GetTransactionStatisticsRequest) ([]map[string]interface{}, error) {
	type Result struct {
		EstateID      uint
		EstateName    string
		Count         int64
		TotalAmount   float64
		AverageAmount float64
	}
	var results []Result
	query := r.db.WithContext(ctx).Table("properties").
		Select("properties.estate_id, estates.name_zh_hant as estate_name, COUNT(*) as count, SUM(properties.price) as total_amount, AVG(properties.price) as average_amount").
		Joins("LEFT JOIN estates ON properties.estate_id = estates.id")
	query = r.applyTransactionFilters(query, filter)
	if err := query.Where("properties.estate_id IS NOT NULL").Group("properties.estate_id, estates.name_zh_hant").Limit(10).Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range results {
		output = append(output, map[string]interface{}{
			"estate_id":      r.EstateID,
			"estate_name":    r.EstateName,
			"count":          r.Count,
			"total_amount":   r.TotalAmount,
			"average_amount": r.AverageAmount,
		})
	}
	return output, nil
}

func (r *statisticsRepository) GetTransactionTrendData(ctx context.Context, filter *models.GetTransactionStatisticsRequest) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := r.db.WithContext(ctx).Table("properties")
	query = r.applyTransactionFilters(query, filter)

	// 根据 period 分组
	var period string
	if filter.Period != nil {
		period = *filter.Period
	}
	switch period {
	case "day":
		query = query.Select("DATE(created_at) as period, COUNT(*) as count, SUM(price) as total_amount, AVG(price) as average_amount").
			Group("DATE(created_at)").
			Order("DATE(created_at) ASC")
	case "week":
		query = query.Select("DATE_TRUNC('week', created_at) as period, COUNT(*) as count, SUM(price) as total_amount, AVG(price) as average_amount").
			Group("DATE_TRUNC('week', created_at)").
			Order("DATE_TRUNC('week', created_at) ASC")
	case "month":
		query = query.Select("TO_CHAR(created_at, 'YYYY-MM') as period, COUNT(*) as count, SUM(price) as total_amount, AVG(price) as average_amount").
			Group("TO_CHAR(created_at, 'YYYY-MM')").
			Order("TO_CHAR(created_at, 'YYYY-MM') ASC")
	case "year":
		query = query.Select("DATE_TRUNC('year', created_at) as period, COUNT(*) as count, SUM(price) as total_amount, AVG(price) as average_amount").
			Group("DATE_TRUNC('year', created_at)").
			Order("DATE_TRUNC('year', created_at) ASC")
	default: // 默认按月
		query = query.Select("TO_CHAR(created_at, 'YYYY-MM') as period, COUNT(*) as count, SUM(price) as total_amount, AVG(price) as average_amount").
			Group("TO_CHAR(created_at, 'YYYY-MM')").
			Order("TO_CHAR(created_at, 'YYYY-MM') ASC")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *statisticsRepository) applyTransactionFilters(query *gorm.DB, filter *models.GetTransactionStatisticsRequest) *gorm.DB {
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.EstateID != nil {
		query = query.Where("estate_id = ?", *filter.EstateID)
	}
	if filter.ListingType != nil {
		query = query.Where("listing_type = ?", *filter.ListingType)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	return query
}

// ========== 用户统计 ==========

func (r *statisticsRepository) GetUserCount(ctx context.Context, filter *models.GetUserStatisticsRequest) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.User{})
	query = r.applyUserFilters(query, filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *statisticsRepository) GetUserCountByStatus(ctx context.Context, filter *models.GetUserStatisticsRequest) (map[string]int64, error) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.User{})
	query = r.applyUserFilters(query, filter)
	if err := query.Select("status, COUNT(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, r := range results {
		statusMap[r.Status] = r.Count
	}
	return statusMap, nil
}

func (r *statisticsRepository) GetUserCountByRole(ctx context.Context, filter *models.GetUserStatisticsRequest) (map[string]int64, error) {
	type Result struct {
		Role  string
		Count int64
	}
	var results []Result
	query := r.db.WithContext(ctx).Model(&models.User{})
	query = r.applyUserFilters(query, filter)
	if err := query.Select("role, COUNT(*) as count").Group("role").Scan(&results).Error; err != nil {
		return nil, err
	}

	roleMap := make(map[string]int64)
	for _, r := range results {
		roleMap[r.Role] = r.Count
	}
	return roleMap, nil
}

func (r *statisticsRepository) GetUserTrendData(ctx context.Context, filter *models.GetUserStatisticsRequest) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := r.db.WithContext(ctx).Model(&models.User{})
	query = r.applyUserFilters(query, filter)

	// 根据 period 分组
	var period string
	if filter.Period != nil {
		period = *filter.Period
	}
	switch period {
	case "day":
		query = query.Select("DATE(created_at) as period, COUNT(*) as count").
			Group("DATE(created_at)").
			Order("DATE(created_at) ASC")
	case "week":
		query = query.Select("DATE_TRUNC('week', created_at) as period, COUNT(*) as count").
			Group("DATE_TRUNC('week', created_at)").
			Order("DATE_TRUNC('week', created_at) ASC")
	case "month":
		query = query.Select("DATE_TRUNC('month', created_at) as period, COUNT(*) as count").
			Group("DATE_TRUNC('month', created_at)").
			Order("DATE_TRUNC('month', created_at) ASC")
	case "year":
		query = query.Select("DATE_TRUNC('year', created_at) as period, COUNT(*) as count").
			Group("DATE_TRUNC('year', created_at)").
			Order("DATE_TRUNC('year', created_at) ASC")
	default: // 默认按月
		query = query.Select("DATE_TRUNC('month', created_at) as period, COUNT(*) as count").
			Group("DATE_TRUNC('month', created_at)").
			Order("DATE_TRUNC('month', created_at) ASC")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *statisticsRepository) applyUserFilters(query *gorm.DB, filter *models.GetUserStatisticsRequest) *gorm.DB {
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	return query
}

// ========== 代理人统计 ==========

func (r *statisticsRepository) GetAgentCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Agent{}).Where("status = ?", "active").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *statisticsRepository) GetAgentAverageRating(ctx context.Context) (float64, error) {
	type Result struct {
		AverageRating float64
	}
	var result Result
	if err := r.db.WithContext(ctx).Model(&models.Agent{}).
		Select("AVG(rating) as average_rating").
		Where("status = ?", "active").
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.AverageRating, nil
}

// ========== 辅助函数 ==========

// GetTimeRangeForPeriod 根据 period 计算时间范围
func GetTimeRangeForPeriod(period string) (time.Time, time.Time) {
	now := time.Now()
	var start time.Time

	switch period {
	case "day":
		start = now.AddDate(0, 0, -30) // 最近30天
	case "week":
		start = now.AddDate(0, 0, -90) // 最近90天（约13周）
	case "month":
		start = now.AddDate(-1, 0, 0) // 最近1年
	case "year":
		start = now.AddDate(-5, 0, 0) // 最近5年
	default:
		start = now.AddDate(-1, 0, 0) // 默认最近1年
	}

	return start, now
}
