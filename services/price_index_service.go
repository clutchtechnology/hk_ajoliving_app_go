package services

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

// PriceIndexService 楼价指数服务接口
//
// PriceIndexService Methods:
// 0. NewPriceIndexService(repo databases.PriceIndexRepository, logger *zap.Logger) -> 注入依赖
// 1. GetPriceIndex(ctx context.Context, filter *models.GetPriceIndexRequest) -> 获取楼价指数列表
// 2. GetLatestPriceIndex(ctx context.Context) -> 获取最新楼价指数
// 3. GetDistrictPriceIndex(ctx context.Context, districtID uint, filter *map[string]interface{}) -> 获取地区楼价指数
// 4. GetEstatePriceIndex(ctx context.Context, estateID uint, filter *map[string]interface{}) -> 获取屋苑楼价指数
// 5. GetPriceTrends(ctx context.Context, filter *models.GetPriceTrendsRequest) -> 获取价格走势
// 6. ComparePriceIndex(ctx context.Context, filter *map[string]interface{}) -> 对比楼价指数
// 7. ExportPriceData(ctx context.Context, filter *map[string]interface{}) -> 导出价格数据
// 8. GetPriceIndexHistory(ctx context.Context, filter *map[string]interface{}) -> 获取历史楼价指数
// 9. CreatePriceIndex(ctx context.Context, req *models.PriceIndex) -> 创建楼价指数
// 10. UpdatePriceIndex(ctx context.Context, id uint, req *models.PriceIndex) -> 更新楼价指数
type PriceIndexService interface {
	GetPriceIndex(ctx context.Context, filter *models.GetPriceIndexRequest) ([]*models.PriceIndex, int64, error)
	GetLatestPriceIndex(ctx context.Context) (*models.PriceIndex, error)
	GetDistrictPriceIndex(ctx context.Context, districtID uint, filter *map[string]interface{}) ([]*models.PriceIndex, error)
	GetEstatePriceIndex(ctx context.Context, estateID uint, filter *map[string]interface{}) ([]*models.PriceIndex, error)
	GetPriceTrends(ctx context.Context, filter *models.GetPriceTrendsRequest) (*map[string]interface{}, error)
	ComparePriceIndex(ctx context.Context, filter *map[string]interface{}) (*map[string]interface{}, error)
	ExportPriceData(ctx context.Context, filter *map[string]interface{}) (*[]byte, error)
	GetPriceIndexHistory(ctx context.Context, filter *map[string]interface{}) (*[]models.PriceIndex, error)
	CreatePriceIndex(ctx context.Context, req *models.PriceIndex) (*models.PriceIndex, error)
	UpdatePriceIndex(ctx context.Context, id uint, req *models.PriceIndex) (*models.PriceIndex, error)
}

type priceIndexService struct {
	repo   databases.PriceIndexRepository
	logger *zap.Logger
}

// 0. NewPriceIndexService 创建楼价指数服务
func NewPriceIndexService(repo databases.PriceIndexRepository, logger *zap.Logger) PriceIndexService {
	return &priceIndexService{
		repo:   repo,
		logger: logger,
	}
}

// 1. GetPriceIndex 获取楼价指数列表
func (s *priceIndexService) GetPriceIndex(ctx context.Context, filter *models.GetPriceIndexRequest) ([]*models.PriceIndex, int64, error) {
	indices, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list price indices", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.PriceIndex, 0, len(indices))
	for _, index := range indices {
		result = append(result, convertToPriceIndexListItemResponse(index))
	}
	
	return result, total, nil
}

// 2. GetLatestPriceIndex 获取最新楼价指数
func (s *priceIndexService) GetLatestPriceIndex(ctx context.Context) (*models.PriceIndex, error) {
	// 获取整体最新指数
	overall, err := s.repo.GetLatest(ctx, string(models.IndexTypeOverall))
	if err != nil {
		s.logger.Error("failed to get latest overall price index", zap.Error(err))
		return nil, err
	}
	
	return overall, nil
}

// 3. GetDistrictPriceIndex 获取地区楼价指数
func (s *priceIndexService) GetDistrictPriceIndex(ctx context.Context, districtID uint, filter *map[string]interface{}) ([]*models.PriceIndex, error) {
	var startPeriod, endPeriod *string
	var limit int = 12
	
	if filter != nil {
		if sp, ok := (*filter)["start_period"].(string); ok {
			startPeriod = &sp
		}
		if ep, ok := (*filter)["end_period"].(string); ok {
			endPeriod = &ep
		}
		if l, ok := (*filter)["limit"].(int); ok {
			limit = l
		}
	}
	
	indices, err := s.repo.GetDistrictPriceIndex(ctx, districtID, startPeriod, endPeriod, limit)
	if err != nil {
		s.logger.Error("failed to get district price index", zap.Error(err), zap.Uint("district_id", districtID))
		return nil, err
	}
	
	result := make([]*models.PriceIndex, 0, len(indices))
	for _, index := range indices {
		result = append(result, convertToPriceIndexResponse(index))
	}
	
	return result, nil
}

// 4. GetEstatePriceIndex 获取屋苑楼价指数
func (s *priceIndexService) GetEstatePriceIndex(ctx context.Context, estateID uint, filter *map[string]interface{}) ([]*models.PriceIndex, error) {
	var startPeriod, endPeriod *string
	var limit int = 12
	
	if filter != nil {
		if sp, ok := (*filter)["start_period"].(string); ok {
			startPeriod = &sp
		}
		if ep, ok := (*filter)["end_period"].(string); ok {
			endPeriod = &ep
		}
		if l, ok := (*filter)["limit"].(int); ok {
			limit = l
		}
	}
	
	indices, err := s.repo.GetEstatePriceIndex(ctx, estateID, startPeriod, endPeriod, limit)
	if err != nil {
		s.logger.Error("failed to get estate price index", zap.Error(err), zap.Uint("estate_id", estateID))
		return nil, err
	}
	
	result := make([]*models.PriceIndex, 0, len(indices))
	for _, index := range indices {
		result = append(result, convertToPriceIndexResponse(index))
	}
	
	return result, nil
}

// 5. GetPriceTrends 获取价格走势
func (s *priceIndexService) GetPriceTrends(ctx context.Context, filter *models.GetPriceTrendsRequest) (*map[string]interface{}, error) {
	indices, err := s.repo.GetTrends(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get price trends", zap.Error(err))
		return nil, err
	}
	
	if len(indices) == 0 {
		return nil, tools.ErrNotFound
	}
	
	// 构建数据点
	dataPoints := make([]map[string]interface{}, 0, len(indices))
	for _, index := range indices {
		dataPoints = append(dataPoints, map[string]interface{}{
			"period":            index.Period,
			"index_value":       index.IndexValue,
			"change_value":      index.ChangeValue,
			"change_percent":    index.ChangePercent,
			"avg_price":         index.AvgPrice,
			"avg_price_per_sqft": index.AvgPricePerSqft,
			"transaction_count": index.TransactionCount,
		})
	}
	
	// 构建响应
	resp := map[string]interface{}{
		"index_type":   filter.IndexType,
		"district_id":  filter.DistrictID,
		"estate_id":    filter.EstateID,
		"property_type": filter.PropertyType,
		"start_period": filter.StartPeriod,
		"end_period":   filter.EndPeriod,
		"data_points":  dataPoints,
	}
	
	// 设置名称
	if len(indices) > 0 {
		if indices[0].District != nil {
			resp["district_name"] = indices[0].District.NameZhHant
		}
		if indices[0].Estate != nil {
			resp["estate_name"] = indices[0].Estate.Name
		}
	}
	
	// 计算统计信息
	resp["statistics"] = calculateTrendStatistics(indices)
	
	return &resp, nil
}

// 6. ComparePriceIndex 对比楼价指数
func (s *priceIndexService) ComparePriceIndex(ctx context.Context, filter *map[string]interface{}) (*map[string]interface{}, error) {
	compareType, _ := (*filter)["compare_type"].(string)
	startPeriod, _ := (*filter)["start_period"].(string)
	endPeriod, _ := (*filter)["end_period"].(string)
	
	resp := map[string]interface{}{
		"compare_type": compareType,
		"start_period": startPeriod,
		"end_period":   endPeriod,
		"series":       []map[string]interface{}{},
	}
	
	switch compareType {
	case "districts":
		districtIDs, _ := (*filter)["district_ids"].([]uint)
		indices, err := s.repo.GetForComparison(ctx, string(models.IndexTypeDistrict), districtIDs, startPeriod, endPeriod)
		if err != nil {
			s.logger.Error("failed to get district comparison data", zap.Error(err))
			return nil, err
		}
		resp["series"] = groupIndicesByID(indices, "district")
		
	case "estates":
		estateIDs, _ := (*filter)["estate_ids"].([]uint)
		indices, err := s.repo.GetForComparison(ctx, string(models.IndexTypeEstate), estateIDs, startPeriod, endPeriod)
		if err != nil {
			s.logger.Error("failed to get estate comparison data", zap.Error(err))
			return nil, err
		}
		resp["series"] = groupIndicesByID(indices, "estate")
		
	case "property_types":
		// 对于物业类型，需要特殊处理
		// TODO: 实现物业类型对比逻辑
		resp["series"] = []map[string]interface{}{}
		
	default:
		return nil, fmt.Errorf("invalid compare type: %s", compareType)
	}
	
	return &resp, nil
}

// 7. ExportPriceData 导出价格数据
func (s *priceIndexService) ExportPriceData(ctx context.Context, filter *map[string]interface{}) (*[]byte, error) {
	indexType, _ := (*filter)["index_type"].(*models.IndexType)
	districtID, _ := (*filter)["district_id"].(*uint)
	estateID, _ := (*filter)["estate_id"].(*uint)
	propertyType, _ := (*filter)["property_type"].(*string)
	startPeriod, _ := (*filter)["start_period"].(string)
	endPeriod, _ := (*filter)["end_period"].(string)
	format, _ := (*filter)["format"].(string)
	
	startPeriodPtr := &startPeriod
	endPeriodPtr := &endPeriod
	
	// 构建查询请求
	listFilter := &models.GetPriceIndexRequest{
		IndexType:    indexType,
		DistrictID:   districtID,
		EstateID:     estateID,
		PropertyType: propertyType,
		StartPeriod:  startPeriodPtr,
		EndPeriod:    endPeriodPtr,
		Page:         1,
		PageSize:     10000, // 导出所有数据
	}
	
	_, total, err := s.repo.List(ctx, listFilter)
	if err != nil {
		s.logger.Error("failed to get data for export", zap.Error(err))
		return nil, err
	}
	
	// TODO: 实现实际的文件生成和上传逻辑
	// 这里只是返回模拟数据
	fileName := fmt.Sprintf("price_index_%s_%s.%s", startPeriod, endPeriod, format)
	
	s.logger.Info("price data exported", zap.String("file_name", fileName), zap.Int64("record_count", total))
	
	// 返回空字节数组，实际应该返回文件内容
	result := make([]byte, 0)
	return &result, nil
}

// 8. GetPriceIndexHistory 获取历史楼价指数
func (s *priceIndexService) GetPriceIndexHistory(ctx context.Context, filter *map[string]interface{}) (*[]models.PriceIndex, error) {
	indexType, _ := (*filter)["index_type"].(string)
	districtID, _ := (*filter)["district_id"].(*uint)
	estateID, _ := (*filter)["estate_id"].(*uint)
	propertyType, _ := (*filter)["property_type"].(*string)
	years, _ := (*filter)["years"].(int)
	if years == 0 {
		years = 1
	}
	
	indices, err := s.repo.GetHistory(ctx, indexType, districtID, estateID, propertyType, years)
	if err != nil {
		s.logger.Error("failed to get price index history", zap.Error(err))
		return nil, err
	}
	
	if len(indices) == 0 {
		return nil, tools.ErrNotFound
	}
	
	result := make([]models.PriceIndex, 0, len(indices))
	for _, index := range indices {
		result = append(result, *convertToPriceIndexResponse(index))
	}
	
	return &result, nil
}

// 9. CreatePriceIndex 创建楼价指数
func (s *priceIndexService) CreatePriceIndex(ctx context.Context, req *models.PriceIndex) (*models.PriceIndex, error) {
	// 解析周期
	year, month, err := parsePeriod(req.Period)
	if err != nil {
		return nil, fmt.Errorf("invalid period format: %w", err)
	}
	
	index := &models.PriceIndex{
		IndexType:        req.IndexType,
		DistrictID:       req.DistrictID,
		EstateID:         req.EstateID,
		PropertyType:     req.PropertyType,
		IndexValue:       req.IndexValue,
		ChangeValue:      req.ChangeValue,
		ChangePercent:    req.ChangePercent,
		AvgPrice:         req.AvgPrice,
		AvgPricePerSqft:  req.AvgPricePerSqft,
		TransactionCount: req.TransactionCount,
		Period:           req.Period,
		Year:             year,
		Month:            month,
		DataSource:       req.DataSource,
		Notes:            req.Notes,
	}
	
	if err := s.repo.Create(ctx, index); err != nil {
		s.logger.Error("failed to create price index", zap.Error(err))
		return nil, err
	}
	
	return convertToPriceIndexResponse(index), nil
}

// 10. UpdatePriceIndex 更新楼价指数
func (s *priceIndexService) UpdatePriceIndex(ctx context.Context, id uint, req *models.PriceIndex) (*models.PriceIndex, error) {
	// 获取现有记录
	index, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get price index", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	if index == nil {
		return nil, tools.ErrNotFound
	}
	
	// 更新字段
	if req.IndexValue != 0 {
		index.IndexValue = req.IndexValue
	}
	if req.ChangeValue != 0 {
		index.ChangeValue = req.ChangeValue
	}
	if req.ChangePercent != 0 {
		index.ChangePercent = req.ChangePercent
	}
	if req.AvgPrice != nil {
		index.AvgPrice = req.AvgPrice
	}
	if req.AvgPricePerSqft != nil {
		index.AvgPricePerSqft = req.AvgPricePerSqft
	}
	if req.TransactionCount != 0 {
		index.TransactionCount = req.TransactionCount
	}
	if req.DataSource != "" {
		index.DataSource = req.DataSource
	}
	if req.Notes != nil {
		index.Notes = req.Notes
	}
	
	if err := s.repo.Update(ctx, index); err != nil {
		s.logger.Error("failed to update price index", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	
	return convertToPriceIndexResponse(index), nil
}

// ========== 辅助函数 ==========

// convertToPriceIndexResponse 转换为楼价指数响应
func convertToPriceIndexResponse(index *models.PriceIndex) *models.PriceIndex {
	resp := &models.PriceIndex{
		ID:               index.ID,
		IndexType:        index.IndexType,
		DistrictID:       index.DistrictID,
		EstateID:         index.EstateID,
		PropertyType:     index.PropertyType,
		IndexValue:       index.IndexValue,
		ChangeValue:      index.ChangeValue,
		ChangePercent:    index.ChangePercent,
		AvgPrice:         index.AvgPrice,
		AvgPricePerSqft:  index.AvgPricePerSqft,
		TransactionCount: index.TransactionCount,
		Period:           index.Period,
		Year:             index.Year,
		Month:            index.Month,
		Day:              index.Day,
		DataSource:       index.DataSource,
		Notes:            index.Notes,
		CreatedAt:        index.CreatedAt,
		UpdatedAt:        index.UpdatedAt,
		District:         index.District,
		Estate:           index.Estate,
	}
	
	return resp
}

// convertToPriceIndexListItemResponse 转换为楼价指数列表项响应
func convertToPriceIndexListItemResponse(index *models.PriceIndex) *models.PriceIndex {
	return &models.PriceIndex{
		ID:               index.ID,
		IndexType:        index.IndexType,
		IndexValue:       index.IndexValue,
		ChangeValue:      index.ChangeValue,
		ChangePercent:    index.ChangePercent,
		AvgPrice:         index.AvgPrice,
		AvgPricePerSqft:  index.AvgPricePerSqft,
		TransactionCount: index.TransactionCount,
		Period:           index.Period,
		Year:             index.Year,
		Month:            index.Month,
	}
}

// calculateTrendStatistics 计算走势统计信息
func calculateTrendStatistics(indices []*models.PriceIndex) map[string]interface{} {
	if len(indices) == 0 {
		return nil
	}
	
	stats := map[string]interface{}{
		"highest_value": indices[0].IndexValue,
		"lowest_value":  indices[0].IndexValue,
	}
	
	var sum float64
	for _, index := range indices {
		sum += index.IndexValue
		if index.IndexValue > stats["highest_value"].(float64) {
			stats["highest_value"] = index.IndexValue
		}
		if index.IndexValue < stats["lowest_value"].(float64) {
			stats["lowest_value"] = index.IndexValue
		}
	}
	
	avgValue := sum / float64(len(indices))
	stats["average_value"] = avgValue
	
	totalChange := indices[len(indices)-1].IndexValue - indices[0].IndexValue
	stats["total_change"] = totalChange
	
	if indices[0].IndexValue != 0 {
		stats["total_change_percent"] = (totalChange / indices[0].IndexValue) * 100
	} else {
		stats["total_change_percent"] = 0.0
	}
	
	// 计算波动率（标准差）
	var variance float64
	for _, index := range indices {
		diff := index.IndexValue - avgValue
		variance += diff * diff
	}
	variance = variance / float64(len(indices))
	stats["volatility_rate"] = math.Sqrt(variance)
	
	return stats
}

// calculateYearlyStatistics 计算年度统计
func calculateYearlyStatistics(indices []*models.PriceIndex) []map[string]interface{} {
	yearMap := make(map[int][]*models.PriceIndex)
	
	for _, index := range indices {
		yearMap[index.Year] = append(yearMap[index.Year], index)
	}
	
	stats := []map[string]interface{}{}
	for year, yearIndices := range yearMap {
		if len(yearIndices) == 0 {
			continue
		}
		
		var sum float64
		var totalTransactions int
		highest := yearIndices[0].IndexValue
		lowest := yearIndices[0].IndexValue
		
		for _, index := range yearIndices {
			sum += index.IndexValue
			totalTransactions += index.TransactionCount
			if index.IndexValue > highest {
				highest = index.IndexValue
			}
			if index.IndexValue < lowest {
				lowest = index.IndexValue
			}
		}
		
		avg := sum / float64(len(yearIndices))
		yearStart := yearIndices[0].IndexValue
		yearEnd := yearIndices[len(yearIndices)-1].IndexValue
		yearChange := yearEnd - yearStart
		yearChangePercent := float64(0)
		if yearStart != 0 {
			yearChangePercent = (yearChange / yearStart) * 100
		}
		
		stats = append(stats, map[string]interface{}{
			"year":                year,
			"average_value":       avg,
			"year_start_value":    yearStart,
			"year_end_value":      yearEnd,
			"year_change":         yearChange,
			"year_change_percent": yearChangePercent,
			"highest_value":       highest,
			"lowest_value":        lowest,
			"total_transactions":  totalTransactions,
		})
	}
	
	return stats
}

// groupIndicesByID 按ID分组指数数据
func groupIndicesByID(indices []*models.PriceIndex, groupType string) []map[string]interface{} {
	groupMap := make(map[uint][]*models.PriceIndex)
	nameMap := make(map[uint]string)
	
	for _, index := range indices {
		var id uint
		var name string
		
		switch groupType {
		case "district":
			if index.DistrictID != nil {
				id = *index.DistrictID
				if index.District != nil {
					name = index.District.NameZhHant
				}
			}
		case "estate":
			if index.EstateID != nil {
				id = *index.EstateID
				if index.Estate != nil {
					name = index.Estate.Name
				}
			}
		}
		
		if id > 0 {
			groupMap[id] = append(groupMap[id], index)
			if name != "" {
				nameMap[id] = name
			}
		}
	}
	
	series := []map[string]interface{}{}
	for id, groupIndices := range groupMap {
		dataPoints := make([]map[string]interface{}, 0, len(groupIndices))
		for _, index := range groupIndices {
			dataPoints = append(dataPoints, map[string]interface{}{
				"period":             index.Period,
				"index_value":        index.IndexValue,
				"change_value":       index.ChangeValue,
				"change_percent":     index.ChangePercent,
				"avg_price":          index.AvgPrice,
				"avg_price_per_sqft": index.AvgPricePerSqft,
				"transaction_count":  index.TransactionCount,
			})
		}
		
		series = append(series, map[string]interface{}{
			"id":          id,
			"name":        nameMap[id],
			"type":        groupType,
			"data_points": dataPoints,
		})
	}
	
	return series
}

// parsePeriod 解析周期字符串
func parsePeriod(period string) (year int, month int, err error) {
	parts := strings.Split(period, "-")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid period format, expected YYYY-MM")
	}
	
	year, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year: %w", err)
	}
	
	month, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month: %w", err)
	}
	
	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("month must be between 1 and 12")
	}
	
	return year, month, nil
}
