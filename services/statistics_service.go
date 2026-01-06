package services

import (
	"context"
	"time"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"go.uber.org/zap"
)

// StatisticsService 统计服务接口
type StatisticsService interface {
	GetOverviewStatistics(ctx context.Context, req *models.GetOverviewStatisticsRequest) (*models.OverviewStatisticsResponse, error)
	GetPropertyStatistics(ctx context.Context, req *models.GetPropertyStatisticsRequest) (*models.PropertyStatisticsResponse, error)
	GetTransactionStatistics(ctx context.Context, req *models.GetTransactionStatisticsRequest) (*models.TransactionStatisticsResponse, error)
	GetUserStatistics(ctx context.Context, req *models.GetUserStatisticsRequest) (*models.UserStatisticsResponse, error)
}

type statisticsService struct {
	repo   databases.StatisticsRepository
	logger *zap.Logger
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService(repo databases.StatisticsRepository, logger *zap.Logger) StatisticsService {
	return &statisticsService{
		repo:   repo,
		logger: logger,
	}
}

// GetOverviewStatistics 获取总览统计
func (s *statisticsService) GetOverviewStatistics(ctx context.Context, req *models.GetOverviewStatisticsRequest) (*models.OverviewStatisticsResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)
	monthStart := todayStart.AddDate(0, -1, 0)

	// 房产统计
	todayStartStr := todayStart.Format("2006-01-02")
	weekStartStr := weekStart.Format("2006-01-02")
	monthStartStr := monthStart.Format("2006-01-02")
	
	propertyActive, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{})
	propertyNewToday, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		StartDate: &todayStartStr,
	})
	propertyNewWeek, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		StartDate: &weekStartStr,
	})
	propertyNewMonth, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		StartDate: &monthStartStr,
	})

	rentType := "rent"
	saleType := "sale"
	propertyRent, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		ListingType: &rentType,
	})
	propertySale, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		ListingType: &saleType,
	})

	propertyPriceStats, _ := s.repo.GetPropertyPriceStats(ctx, &models.GetPropertyStatisticsRequest{})

	// 用户统计
	userTotal, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{})
	userStatusMap, _ := s.repo.GetUserCountByStatus(ctx, &models.GetUserStatisticsRequest{})
	userNewToday, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{
		StartDate: &todayStartStr,
	})
	userNewWeek, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{
		StartDate: &weekStartStr,
	})
	userNewMonth, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{
		StartDate: &monthStartStr,
	})

	// 代理人统计
	agentTotal, _ := s.repo.GetAgentCount(ctx)
	agentAvgRating, _ := s.repo.GetAgentAverageRating(ctx)

	// 成交统计
	transactionTotal, _ := s.repo.GetTransactionCount(ctx, &models.GetTransactionStatisticsRequest{})
	transactionToday, _ := s.repo.GetTransactionCount(ctx, &models.GetTransactionStatisticsRequest{
		StartDate: &todayStartStr,
	})
	transactionWeek, _ := s.repo.GetTransactionCount(ctx, &models.GetTransactionStatisticsRequest{
		StartDate: &weekStartStr,
	})
	transactionMonth, _ := s.repo.GetTransactionCount(ctx, &models.GetTransactionStatisticsRequest{
		StartDate: &monthStartStr,
	})
	transactionAmountStats, _ := s.repo.GetTransactionAmountStats(ctx, &models.GetTransactionStatisticsRequest{})

	resp := &models.OverviewStatisticsResponse{
		Properties: models.PropertyOverview{
			TotalCount:    int(propertyActive),
			ActiveCount:   int(propertyActive),
			RentCount:     int(propertyRent),
			SaleCount:     int(propertySale),
			NewToday:      int(propertyNewToday),
			NewThisWeek:   int(propertyNewWeek),
			NewThisMonth:  int(propertyNewMonth),
			AveragePrice:  propertyPriceStats["average"],
			TotalValue:    propertyPriceStats["total"],
		},
		Users: models.UserOverview{
			TotalCount:    int(userTotal),
			ActiveCount:   int(userStatusMap["active"]),
			NewToday:      int(userNewToday),
			NewThisWeek:   int(userNewWeek),
			NewThisMonth:  int(userNewMonth),
			VerifiedCount: int(userStatusMap["verified"]),
		},
		Agents: models.AgentOverview{
			TotalCount:    int(agentTotal),
			ActiveCount:   int(agentTotal),
			AverageRating: agentAvgRating,
			TopAgents:     0,
		},
		Transactions: models.TransactionOverview{
			TotalCount:      int(transactionTotal),
			TodayCount:      int(transactionToday),
			ThisWeekCount:   int(transactionWeek),
			ThisMonthCount:  int(transactionMonth),
			TotalAmount:     transactionAmountStats["total"],
			AverageAmount:   transactionAmountStats["average"],
		},
		PlatformMetrics: models.PlatformMetrics{
			TotalViews:     0,
			TodayViews:     0,
			SearchCount:    0,
			ConversionRate: 0.0,
		},
	}

	return resp, nil
}

// GetPropertyStatistics 获取房产统计
func (s *statisticsService) GetPropertyStatistics(ctx context.Context, req *models.GetPropertyStatisticsRequest) (*models.PropertyStatisticsResponse, error) {
	// 设置默认时间范围
	if req.StartDate == nil && req.EndDate == nil {
		period := "month"
		if req.Period != nil {
			period = *req.Period
		}
		start, end := databases.GetTimeRangeForPeriod(period)
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")
		req.StartDate = &startStr
		req.EndDate = &endStr
	}

	// 汇总统计
	totalCount, _ := s.repo.GetPropertyCount(ctx, req)
	
	rentType := "rent"
	saleType := "sale"
	rentCount, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		DistrictID:  req.DistrictID,
		EstateID:    req.EstateID,
		ListingType: &rentType,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	saleCount, _ := s.repo.GetPropertyCount(ctx, &models.GetPropertyStatisticsRequest{
		DistrictID:  req.DistrictID,
		EstateID:    req.EstateID,
		ListingType: &saleType,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})

	priceStats, _ := s.repo.GetPropertyPriceStats(ctx, req)

	summary := &models.PropertyStatisticsSummary{
		TotalCount:   int(totalCount),
		RentCount:    int(rentCount),
		SaleCount:    int(saleCount),
		AveragePrice: priceStats["average"],
		MedianPrice:  priceStats["average"], // TODO: 实现中位数计算
		HighestPrice: priceStats["max"],
		LowestPrice:  priceStats["min"],
		TotalValue:   priceStats["total"],
	}

	// 趋势数据
	trendDataRaw, _ := s.repo.GetPropertyTrendData(ctx, req)
	trendData := convertToTrendItems(trendDataRaw)

	// 分布统计
	districtData, _ := s.repo.GetPropertyCountByDistrict(ctx, req)
	typeData, _ := s.repo.GetPropertyCountByType(ctx, req)
	bedroomData, _ := s.repo.GetPropertyCountByBedrooms(ctx, req)

	distribution := &models.PropertyDistribution{
		ByDistrict:     convertToDistrictStatItems(districtData, int(totalCount)),
		ByPropertyType: convertToPropertyTypeStatItems(typeData, int(totalCount)),
		ByPriceRange:   calculatePriceRangeDistribution(priceStats, int(totalCount)),
		ByBedroomCount: convertToBedroomCountStatItems(bedroomData, int(totalCount)),
	}

	return &models.PropertyStatisticsResponse{
		Summary:      summary,
		TrendData:    trendData,
		Distribution: distribution,
	}, nil
}

// GetTransactionStatistics 获取成交统计
func (s *statisticsService) GetTransactionStatistics(ctx context.Context, req *models.GetTransactionStatisticsRequest) (*models.TransactionStatisticsResponse, error) {
	// 设置默认时间范围
	if req.StartDate == nil && req.EndDate == nil {
		period := "month"
		if req.Period != nil {
			period = *req.Period
		}
		start, end := databases.GetTimeRangeForPeriod(period)
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")
		req.StartDate = &startStr
		req.EndDate = &endStr
	}

	// 汇总统计
	totalCount, _ := s.repo.GetTransactionCount(ctx, req)
	amountStats, _ := s.repo.GetTransactionAmountStats(ctx, req)

	summary := &models.TransactionStatisticsSummary{
		TotalCount:          int(totalCount),
		TotalAmount:         amountStats["total"],
		AverageAmount:       amountStats["average"],
		MedianAmount:        amountStats["average"], // TODO: 实现中位数计算
		HighestAmount:       amountStats["max"],
		LowestAmount:        amountStats["min"],
		AveragePricePerSqft: 0, // TODO: 计算每平方尺价格
	}

	// 趋势数据
	trendDataRaw, _ := s.repo.GetTransactionTrendData(ctx, req)
	trendData := convertToTransactionTrendItems(trendDataRaw)

	// 分布统计
	districtData, _ := s.repo.GetTransactionCountByDistrict(ctx, req)
	estateData, _ := s.repo.GetTransactionCountByEstate(ctx, req)

	distribution := &models.TransactionDistribution{
		ByDistrict: convertToDistrictTransactionItems(districtData),
		ByEstate:   convertToEstateTransactionItems(estateData),
		ByMonth:    convertToMonthTransactionItems(trendDataRaw),
	}

	return &models.TransactionStatisticsResponse{
		Summary:      summary,
		TrendData:    trendData,
		Distribution: distribution,
	}, nil
}

// GetUserStatistics 获取用户统计
func (s *statisticsService) GetUserStatistics(ctx context.Context, req *models.GetUserStatisticsRequest) (*models.UserStatisticsResponse, error) {
	// 设置默认时间范围
	if req.StartDate == nil && req.EndDate == nil {
		period := "month"
		if req.Period != nil {
			period = *req.Period
		}
		start, end := databases.GetTimeRangeForPeriod(period)
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")
		req.StartDate = &startStr
		req.EndDate = &endStr
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)
	monthStart := todayStart.AddDate(0, -1, 0)
	
	todayStartStr := todayStart.Format("2006-01-02")
	weekStartStr := weekStart.Format("2006-01-02")
	monthStartStr := monthStart.Format("2006-01-02")

	// 汇总统计
	totalCount, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{})
	statusMap, _ := s.repo.GetUserCountByStatus(ctx, req)
	newToday, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{
		StartDate: &todayStartStr,
	})
	newWeek, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{
		StartDate: &weekStartStr,
	})
	newMonth, _ := s.repo.GetUserCount(ctx, &models.GetUserStatisticsRequest{
		StartDate: &monthStartStr,
	})

	summary := &models.UserStatisticsSummary{
		TotalCount:        int(totalCount),
		ActiveCount:       int(statusMap["active"]),
		NewUsersToday:     int(newToday),
		NewUsersThisWeek:  int(newWeek),
		NewUsersThisMonth: int(newMonth),
		VerifiedCount:     int(statusMap["verified"]),
		RetentionRate:     0.0, // TODO: 实现留存率计算
	}

	// 趋势数据
	trendDataRaw, _ := s.repo.GetUserTrendData(ctx, req)
	trendData := convertToTrendItems(trendDataRaw)

	// 分布统计
	roleMap, _ := s.repo.GetUserCountByRole(ctx, req)

	distribution := &models.UserDistribution{
		ByRole:               convertToUserRoleStatItems(roleMap, int(totalCount)),
		ByStatus:             convertToUserStatusStatItems(statusMap, int(totalCount)),
		ByRegistrationSource: []*models.RegistrationSourceStatItem{}, // TODO: 实现注册来源统计
	}

	return &models.UserStatisticsResponse{
		Summary:      summary,
		TrendData:    trendData,
		Distribution: distribution,
	}, nil
}

// ========== 转换辅助函数 ==========

func convertToTrendItems(data []map[string]interface{}) []*models.TrendItem {
	var items []*models.TrendItem
	for _, d := range data {
		item := &models.TrendItem{}
		if period, ok := d["period"].(time.Time); ok {
			item.Period = period.Format("2006-01-02")
		} else if period, ok := d["period"].(string); ok {
			item.Period = period
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
		}
		if value, ok := d["value"].(float64); ok {
			item.Value = value
		}
		items = append(items, item)
	}
	return items
}

func convertToTransactionTrendItems(data []map[string]interface{}) []*models.TransactionTrendItem {
	var items []*models.TransactionTrendItem
	for _, d := range data {
		item := &models.TransactionTrendItem{}
		if period, ok := d["period"].(time.Time); ok {
			item.Period = period.Format("2006-01-02")
		} else if period, ok := d["period"].(string); ok {
			item.Period = period
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
		}
		if avgAmount, ok := d["average_amount"].(float64); ok {
			item.AverageAmount = avgAmount
		}
		items = append(items, item)
	}
	return items
}

func convertToDistrictStatItems(data []map[string]interface{}, total int) []*models.DistrictStatItem {
	var items []*models.DistrictStatItem
	for _, d := range data {
		item := &models.DistrictStatItem{}
		if id, ok := d["district_id"].(uint); ok {
			item.DistrictID = id
		}
		if name, ok := d["district_name"].(string); ok {
			item.DistrictName = name
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
			if total > 0 {
				item.Percentage = float64(count) / float64(total) * 100
			}
		}
		if avgPrice, ok := d["average_price"].(float64); ok {
			item.AveragePrice = avgPrice
		}
		items = append(items, item)
	}
	return items
}

func convertToPropertyTypeStatItems(data []map[string]interface{}, total int) []*models.PropertyTypeStatItem {
	var items []*models.PropertyTypeStatItem
	for _, d := range data {
		item := &models.PropertyTypeStatItem{}
		if propType, ok := d["property_type"].(string); ok {
			item.PropertyType = propType
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
			if total > 0 {
				item.Percentage = float64(count) / float64(total) * 100
			}
		}
		items = append(items, item)
	}
	return items
}

func convertToBedroomCountStatItems(data []map[string]interface{}, total int) []*models.BedroomCountStatItem {
	var items []*models.BedroomCountStatItem
	for _, d := range data {
		item := &models.BedroomCountStatItem{}
		if bedrooms, ok := d["bedrooms"].(int); ok {
			item.Bedrooms = bedrooms
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
			if total > 0 {
				item.Percentage = float64(count) / float64(total) * 100
			}
		}
		items = append(items, item)
	}
	return items
}

func calculatePriceRangeDistribution(priceStats map[string]float64, total int) []*models.PriceRangeStatItem {
	// 简化版价格区间分布
	// TODO: 根据实际价格数据动态计算区间
	return []*models.PriceRangeStatItem{
		{Range: "0-1M", Count: 0, Percentage: 0},
		{Range: "1M-2M", Count: 0, Percentage: 0},
		{Range: "2M-5M", Count: 0, Percentage: 0},
		{Range: "5M+", Count: 0, Percentage: 0},
	}
}

func convertToDistrictTransactionItems(data []map[string]interface{}) []*models.DistrictTransactionItem {
	var items []*models.DistrictTransactionItem
	for _, d := range data {
		item := &models.DistrictTransactionItem{}
		if id, ok := d["district_id"].(uint); ok {
			item.DistrictID = id
		}
		if name, ok := d["district_name"].(string); ok {
			item.DistrictName = name
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
		}
		if total, ok := d["total_amount"].(float64); ok {
			item.TotalAmount = total
		}
		if avg, ok := d["average_amount"].(float64); ok {
			item.AverageAmount = avg
		}
		items = append(items, item)
	}
	return items
}

func convertToEstateTransactionItems(data []map[string]interface{}) []*models.EstateTransactionItem {
	var items []*models.EstateTransactionItem
	for _, d := range data {
		item := &models.EstateTransactionItem{}
		if id, ok := d["estate_id"].(uint); ok {
			item.EstateID = id
		}
		if name, ok := d["estate_name"].(string); ok {
			item.EstateName = name
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
		}
		if total, ok := d["total_amount"].(float64); ok {
			item.TotalAmount = total
		}
		if avg, ok := d["average_amount"].(float64); ok {
			item.AverageAmount = avg
		}
		items = append(items, item)
	}
	return items
}

func convertToMonthTransactionItems(data []map[string]interface{}) []*models.MonthTransactionItem {
	var items []*models.MonthTransactionItem
	for _, d := range data {
		item := &models.MonthTransactionItem{}
		if period, ok := d["period"].(time.Time); ok {
			item.Month = period.Format("2006-01")
		} else if period, ok := d["period"].(string); ok {
			item.Month = period
		}
		if count, ok := d["count"].(int64); ok {
			item.Count = int(count)
		}
		if total, ok := d["total_amount"].(float64); ok {
			item.TotalAmount = total
		}
		if avg, ok := d["average_amount"].(float64); ok {
			item.AverageAmount = avg
		}
		items = append(items, item)
	}
	return items
}

func convertToUserRoleStatItems(roleMap map[string]int64, total int) []*models.UserRoleStatItem {
	var items []*models.UserRoleStatItem
	for role, count := range roleMap {
		item := &models.UserRoleStatItem{
			Role:  role,
			Count: int(count),
		}
		if total > 0 {
			item.Percentage = float64(count) / float64(total) * 100
		}
		items = append(items, item)
	}
	return items
}

func convertToUserStatusStatItems(statusMap map[string]int64, total int) []*models.UserStatusStatItem {
	var items []*models.UserStatusStatItem
	for status, count := range statusMap {
		item := &models.UserStatusStatItem{
			Status: status,
			Count:  int(count),
		}
		if total > 0 {
			item.Percentage = float64(count) / float64(total) * 100
		}
		items = append(items, item)
	}
	return items
}
