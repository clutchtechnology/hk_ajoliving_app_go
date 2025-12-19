package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

// PriceIndexService 楼价指数服务接口
//
// PriceIndexService Methods:
// 0. NewPriceIndexService(repo repository.PriceIndexRepository, logger *zap.Logger) -> 注入依赖
// 1. GetPriceIndex(ctx context.Context, filter *request.GetPriceIndexRequest) -> 获取楼价指数列表
// 2. GetLatestPriceIndex(ctx context.Context) -> 获取最新楼价指数
// 3. GetDistrictPriceIndex(ctx context.Context, districtID uint, filter *request.GetDistrictPriceIndexRequest) -> 获取地区楼价指数
// 4. GetEstatePriceIndex(ctx context.Context, estateID uint, filter *request.GetEstatePriceIndexRequest) -> 获取屋苑楼价指数
// 5. GetPriceTrends(ctx context.Context, filter *request.GetPriceTrendsRequest) -> 获取价格走势
// 6. ComparePriceIndex(ctx context.Context, filter *request.ComparePriceIndexRequest) -> 对比楼价指数
// 7. ExportPriceData(ctx context.Context, filter *request.ExportPriceDataRequest) -> 导出价格数据
// 8. GetPriceIndexHistory(ctx context.Context, filter *request.GetPriceIndexHistoryRequest) -> 获取历史楼价指数
// 9. CreatePriceIndex(ctx context.Context, req *request.CreatePriceIndexRequest) -> 创建楼价指数
// 10. UpdatePriceIndex(ctx context.Context, id uint, req *request.UpdatePriceIndexRequest) -> 更新楼价指数
type PriceIndexService interface {
	GetPriceIndex(ctx context.Context, filter *request.GetPriceIndexRequest) ([]*response.PriceIndexListItemResponse, int64, error)
	GetLatestPriceIndex(ctx context.Context) (*response.LatestPriceIndexResponse, error)
	GetDistrictPriceIndex(ctx context.Context, districtID uint, filter *request.GetDistrictPriceIndexRequest) ([]*response.PriceIndexResponse, error)
	GetEstatePriceIndex(ctx context.Context, estateID uint, filter *request.GetEstatePriceIndexRequest) ([]*response.PriceIndexResponse, error)
	GetPriceTrends(ctx context.Context, filter *request.GetPriceTrendsRequest) (*response.PriceTrendResponse, error)
	ComparePriceIndex(ctx context.Context, filter *request.ComparePriceIndexRequest) (*response.ComparePriceIndexResponse, error)
	ExportPriceData(ctx context.Context, filter *request.ExportPriceDataRequest) (*response.ExportPriceDataResponse, error)
	GetPriceIndexHistory(ctx context.Context, filter *request.GetPriceIndexHistoryRequest) (*response.PriceIndexHistoryResponse, error)
	CreatePriceIndex(ctx context.Context, req *request.CreatePriceIndexRequest) (*response.CreatePriceIndexResponse, error)
	UpdatePriceIndex(ctx context.Context, id uint, req *request.UpdatePriceIndexRequest) (*response.UpdatePriceIndexResponse, error)
}

type priceIndexService struct {
	repo   repository.PriceIndexRepository
	logger *zap.Logger
}

// 0. NewPriceIndexService 创建楼价指数服务
func NewPriceIndexService(repo repository.PriceIndexRepository, logger *zap.Logger) PriceIndexService {
	return &priceIndexService{
		repo:   repo,
		logger: logger,
	}
}

// 1. GetPriceIndex 获取楼价指数列表
func (s *priceIndexService) GetPriceIndex(ctx context.Context, filter *request.GetPriceIndexRequest) ([]*response.PriceIndexListItemResponse, int64, error) {
	indices, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list price indices", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.PriceIndexListItemResponse, 0, len(indices))
	for _, index := range indices {
		result = append(result, convertToPriceIndexListItemResponse(index))
	}
	
	return result, total, nil
}

// 2. GetLatestPriceIndex 获取最新楼价指数
func (s *priceIndexService) GetLatestPriceIndex(ctx context.Context) (*response.LatestPriceIndexResponse, error) {
	// 获取整体最新指数
	overall, err := s.repo.GetLatest(ctx, string(model.IndexTypeOverall))
	if err != nil {
		s.logger.Error("failed to get latest overall price index", zap.Error(err))
		// 不返回错误，继续获取其他数据
	}
	
	// 获取各地区最新指数
	districtIndices, err := s.repo.GetAllLatestByType(ctx, string(model.IndexTypeDistrict))
	if err != nil {
		s.logger.Error("failed to get latest district price indices", zap.Error(err))
		districtIndices = []*model.PriceIndex{}
	}
	
	// 获取各物业类型最新指数
	propertyTypeIndices, err := s.repo.GetAllLatestByType(ctx, string(model.IndexTypePropertyType))
	if err != nil {
		s.logger.Error("failed to get latest property type price indices", zap.Error(err))
		propertyTypeIndices = []*model.PriceIndex{}
	}
	
	// 构建响应
	resp := &response.LatestPriceIndexResponse{
		ByDistrict:     make([]response.PriceIndexResponse, 0, len(districtIndices)),
		ByPropertyType: make([]response.PriceIndexResponse, 0, len(propertyTypeIndices)),
		UpdatedAt:      time.Now(),
	}
	
	if overall != nil {
		overallResp := convertToPriceIndexResponse(overall)
		resp.Overall = overallResp
		resp.UpdatedAt = overall.UpdatedAt
	}
	
	for _, index := range districtIndices {
		resp.ByDistrict = append(resp.ByDistrict, *convertToPriceIndexResponse(index))
	}
	
	for _, index := range propertyTypeIndices {
		resp.ByPropertyType = append(resp.ByPropertyType, *convertToPriceIndexResponse(index))
	}
	
	return resp, nil
}

// 3. GetDistrictPriceIndex 获取地区楼价指数
func (s *priceIndexService) GetDistrictPriceIndex(ctx context.Context, districtID uint, filter *request.GetDistrictPriceIndexRequest) ([]*response.PriceIndexResponse, error) {
	indices, err := s.repo.GetDistrictPriceIndex(ctx, districtID, filter.StartPeriod, filter.EndPeriod, filter.Limit)
	if err != nil {
		s.logger.Error("failed to get district price index", zap.Error(err), zap.Uint("district_id", districtID))
		return nil, err
	}
	
	result := make([]*response.PriceIndexResponse, 0, len(indices))
	for _, index := range indices {
		result = append(result, convertToPriceIndexResponse(index))
	}
	
	return result, nil
}

// 4. GetEstatePriceIndex 获取屋苑楼价指数
func (s *priceIndexService) GetEstatePriceIndex(ctx context.Context, estateID uint, filter *request.GetEstatePriceIndexRequest) ([]*response.PriceIndexResponse, error) {
	indices, err := s.repo.GetEstatePriceIndex(ctx, estateID, filter.StartPeriod, filter.EndPeriod, filter.Limit)
	if err != nil {
		s.logger.Error("failed to get estate price index", zap.Error(err), zap.Uint("estate_id", estateID))
		return nil, err
	}
	
	result := make([]*response.PriceIndexResponse, 0, len(indices))
	for _, index := range indices {
		result = append(result, convertToPriceIndexResponse(index))
	}
	
	return result, nil
}

// 5. GetPriceTrends 获取价格走势
func (s *priceIndexService) GetPriceTrends(ctx context.Context, filter *request.GetPriceTrendsRequest) (*response.PriceTrendResponse, error) {
	indices, err := s.repo.GetTrends(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get price trends", zap.Error(err))
		return nil, err
	}
	
	if len(indices) == 0 {
		return nil, errors.ErrNotFound
	}
	
	// 构建响应
	resp := &response.PriceTrendResponse{
		IndexType:   filter.IndexType,
		DistrictID:  filter.DistrictID,
		EstateID:    filter.EstateID,
		PropertyType: filter.PropertyType,
		StartPeriod: filter.StartPeriod,
		EndPeriod:   filter.EndPeriod,
		DataPoints:  make([]response.PriceTrendDataPoint, 0, len(indices)),
	}
	
	// 设置名称
	if len(indices) > 0 {
		if indices[0].District != nil {
			districtName := indices[0].District.NameZhHant
			resp.DistrictName = &districtName
		}
		if indices[0].Estate != nil {
			estateName := indices[0].Estate.Name
			resp.EstateName = &estateName
		}
	}
	
	// 转换数据点
	for _, index := range indices {
		resp.DataPoints = append(resp.DataPoints, response.PriceTrendDataPoint{
			Period:           index.Period,
			IndexValue:       index.IndexValue,
			ChangeValue:      index.ChangeValue,
			ChangePercent:    index.ChangePercent,
			AvgPrice:         index.AvgPrice,
			AvgPricePerSqft:  index.AvgPricePerSqft,
			TransactionCount: index.TransactionCount,
		})
	}
	
	// 计算统计信息
	resp.Statistics = calculateTrendStatistics(indices)
	
	return resp, nil
}

// 6. ComparePriceIndex 对比楼价指数
func (s *priceIndexService) ComparePriceIndex(ctx context.Context, filter *request.ComparePriceIndexRequest) (*response.ComparePriceIndexResponse, error) {
	resp := &response.ComparePriceIndexResponse{
		CompareType: filter.CompareType,
		StartPeriod: filter.StartPeriod,
		EndPeriod:   filter.EndPeriod,
		Series:      []response.CompareSeriesData{},
	}
	
	switch filter.CompareType {
	case "districts":
		indices, err := s.repo.GetForComparison(ctx, string(model.IndexTypeDistrict), filter.DistrictIDs, filter.StartPeriod, filter.EndPeriod)
		if err != nil {
			s.logger.Error("failed to get district comparison data", zap.Error(err))
			return nil, err
		}
		resp.Series = groupIndicesByID(indices, "district")
		
	case "estates":
		indices, err := s.repo.GetForComparison(ctx, string(model.IndexTypeEstate), filter.EstateIDs, filter.StartPeriod, filter.EndPeriod)
		if err != nil {
			s.logger.Error("failed to get estate comparison data", zap.Error(err))
			return nil, err
		}
		resp.Series = groupIndicesByID(indices, "estate")
		
	case "property_types":
		// 对于物业类型，需要特殊处理
		// TODO: 实现物业类型对比逻辑
		resp.Series = []response.CompareSeriesData{}
		
	default:
		return nil, fmt.Errorf("invalid compare type: %s", filter.CompareType)
	}
	
	return resp, nil
}

// 7. ExportPriceData 导出价格数据
func (s *priceIndexService) ExportPriceData(ctx context.Context, filter *request.ExportPriceDataRequest) (*response.ExportPriceDataResponse, error) {
	// 构建查询请求
	listFilter := &request.GetPriceIndexRequest{
		IndexType:    filter.IndexType,
		DistrictID:   filter.DistrictID,
		EstateID:     filter.EstateID,
		PropertyType: filter.PropertyType,
		StartPeriod:  &filter.StartPeriod,
		EndPeriod:    &filter.EndPeriod,
		Page:         1,
		PageSize:     10000, // 导出所有数据
	}
	
	indices, total, err := s.repo.List(ctx, listFilter)
	if err != nil {
		s.logger.Error("failed to get data for export", zap.Error(err))
		return nil, err
	}
	
	// TODO: 实现实际的文件生成和上传逻辑
	// 这里只是返回模拟数据
	fileName := fmt.Sprintf("price_index_%s_%s.%s", filter.StartPeriod, filter.EndPeriod, filter.Format)
	
	resp := &response.ExportPriceDataResponse{
		FileName:    fileName,
		DownloadURL: fmt.Sprintf("/downloads/%s", fileName),
		Format:      filter.Format,
		RecordCount: len(indices),
		ExportedAt:  time.Now(),
	}
	
	s.logger.Info("price data exported", zap.String("file_name", fileName), zap.Int64("record_count", total))
	
	return resp, nil
}

// 8. GetPriceIndexHistory 获取历史楼价指数
func (s *priceIndexService) GetPriceIndexHistory(ctx context.Context, filter *request.GetPriceIndexHistoryRequest) (*response.PriceIndexHistoryResponse, error) {
	indices, err := s.repo.GetHistory(ctx, filter.IndexType, filter.DistrictID, filter.EstateID, filter.PropertyType, filter.Years)
	if err != nil {
		s.logger.Error("failed to get price index history", zap.Error(err))
		return nil, err
	}
	
	if len(indices) == 0 {
		return nil, errors.ErrNotFound
	}
	
	resp := &response.PriceIndexHistoryResponse{
		IndexType:    filter.IndexType,
		DistrictID:   filter.DistrictID,
		EstateID:     filter.EstateID,
		PropertyType: filter.PropertyType,
		Years:        filter.Years,
		DataPoints:   make([]response.PriceTrendDataPoint, 0, len(indices)),
		YearlyStats:  []response.YearlyStatistics{},
	}
	
	// 设置名称
	if len(indices) > 0 {
		if indices[0].District != nil {
			districtName := indices[0].District.NameZhHant
			resp.DistrictName = &districtName
		}
		if indices[0].Estate != nil {
			estateName := indices[0].Estate.Name
			resp.EstateName = &estateName
		}
	}
	
	// 转换数据点
	for _, index := range indices {
		resp.DataPoints = append(resp.DataPoints, response.PriceTrendDataPoint{
			Period:           index.Period,
			IndexValue:       index.IndexValue,
			ChangeValue:      index.ChangeValue,
			ChangePercent:    index.ChangePercent,
			AvgPrice:         index.AvgPrice,
			AvgPricePerSqft:  index.AvgPricePerSqft,
			TransactionCount: index.TransactionCount,
		})
	}
	
	// 计算年度统计
	resp.YearlyStats = calculateYearlyStatistics(indices)
	
	return resp, nil
}

// 9. CreatePriceIndex 创建楼价指数
func (s *priceIndexService) CreatePriceIndex(ctx context.Context, req *request.CreatePriceIndexRequest) (*response.CreatePriceIndexResponse, error) {
	// 解析周期
	year, month, err := parsePeriod(req.Period)
	if err != nil {
		return nil, fmt.Errorf("invalid period format: %w", err)
	}
	
	index := &model.PriceIndex{
		IndexType:        model.IndexType(req.IndexType),
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
	
	return &response.CreatePriceIndexResponse{
		ID:      index.ID,
		Period:  index.Period,
		Message: "楼价指数创建成功",
	}, nil
}

// 10. UpdatePriceIndex 更新楼价指数
func (s *priceIndexService) UpdatePriceIndex(ctx context.Context, id uint, req *request.UpdatePriceIndexRequest) (*response.UpdatePriceIndexResponse, error) {
	// 获取现有记录
	index, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get price index", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	if index == nil {
		return nil, errors.ErrNotFound
	}
	
	// 更新字段
	if req.IndexValue != nil {
		index.IndexValue = *req.IndexValue
	}
	if req.ChangeValue != nil {
		index.ChangeValue = *req.ChangeValue
	}
	if req.ChangePercent != nil {
		index.ChangePercent = *req.ChangePercent
	}
	if req.AvgPrice != nil {
		index.AvgPrice = req.AvgPrice
	}
	if req.AvgPricePerSqft != nil {
		index.AvgPricePerSqft = req.AvgPricePerSqft
	}
	if req.TransactionCount != nil {
		index.TransactionCount = *req.TransactionCount
	}
	if req.DataSource != nil {
		index.DataSource = *req.DataSource
	}
	if req.Notes != nil {
		index.Notes = req.Notes
	}
	
	if err := s.repo.Update(ctx, index); err != nil {
		s.logger.Error("failed to update price index", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	
	return &response.UpdatePriceIndexResponse{
		ID:      index.ID,
		Message: "楼价指数更新成功",
	}, nil
}

// ========== 辅助函数 ==========

// convertToPriceIndexResponse 转换为楼价指数响应
func convertToPriceIndexResponse(index *model.PriceIndex) *response.PriceIndexResponse {
	resp := &response.PriceIndexResponse{
		ID:               index.ID,
		IndexType:        string(index.IndexType),
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
	}
	
	if index.District != nil {
		resp.District = &response.DistrictResponse{
			ID:         index.District.ID,
			NameZhHant: index.District.NameZhHant,
			Region:     string(index.District.Region),
		}
		if index.District.NameZhHans != nil {
			resp.District.NameZhHans = *index.District.NameZhHans
		}
		if index.District.NameEn != nil {
			resp.District.NameEn = *index.District.NameEn
		}
	}
	
	if index.Estate != nil {
		resp.Estate = &response.EstateBasicInfo{
			ID:     index.Estate.ID,
			Name:   index.Estate.Name,
			NameEn: index.Estate.NameEn,
		}
	}
	
	return resp
}

// convertToPriceIndexListItemResponse 转换为楼价指数列表项响应
func convertToPriceIndexListItemResponse(index *model.PriceIndex) *response.PriceIndexListItemResponse {
	return &response.PriceIndexListItemResponse{
		ID:               index.ID,
		IndexType:        string(index.IndexType),
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
func calculateTrendStatistics(indices []*model.PriceIndex) *response.TrendStatistics {
	if len(indices) == 0 {
		return nil
	}
	
	stats := &response.TrendStatistics{
		HighestValue: indices[0].IndexValue,
		LowestValue:  indices[0].IndexValue,
	}
	
	var sum float64
	for _, index := range indices {
		sum += index.IndexValue
		if index.IndexValue > stats.HighestValue {
			stats.HighestValue = index.IndexValue
		}
		if index.IndexValue < stats.LowestValue {
			stats.LowestValue = index.IndexValue
		}
	}
	
	stats.AverageValue = sum / float64(len(indices))
	stats.TotalChange = indices[len(indices)-1].IndexValue - indices[0].IndexValue
	if indices[0].IndexValue != 0 {
		stats.TotalChangePercent = (stats.TotalChange / indices[0].IndexValue) * 100
	}
	
	// 计算波动率（标准差）
	var variance float64
	for _, index := range indices {
		diff := index.IndexValue - stats.AverageValue
		variance += diff * diff
	}
	variance = variance / float64(len(indices))
	stats.VolatilityRate = math.Sqrt(variance)
	
	return stats
}

// calculateYearlyStatistics 计算年度统计
func calculateYearlyStatistics(indices []*model.PriceIndex) []response.YearlyStatistics {
	yearMap := make(map[int][]*model.PriceIndex)
	
	for _, index := range indices {
		yearMap[index.Year] = append(yearMap[index.Year], index)
	}
	
	stats := []response.YearlyStatistics{}
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
		
		stats = append(stats, response.YearlyStatistics{
			Year:              year,
			AverageValue:      avg,
			YearStartValue:    yearStart,
			YearEndValue:      yearEnd,
			YearChange:        yearChange,
			YearChangePercent: yearChangePercent,
			HighestValue:      highest,
			LowestValue:       lowest,
			TotalTransactions: totalTransactions,
		})
	}
	
	return stats
}

// groupIndicesByID 按ID分组指数数据
func groupIndicesByID(indices []*model.PriceIndex, groupType string) []response.CompareSeriesData {
	groupMap := make(map[uint][]*model.PriceIndex)
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
	
	series := []response.CompareSeriesData{}
	for id, groupIndices := range groupMap {
		dataPoints := make([]response.PriceTrendDataPoint, 0, len(groupIndices))
		for _, index := range groupIndices {
			dataPoints = append(dataPoints, response.PriceTrendDataPoint{
				Period:           index.Period,
				IndexValue:       index.IndexValue,
				ChangeValue:      index.ChangeValue,
				ChangePercent:    index.ChangePercent,
				AvgPrice:         index.AvgPrice,
				AvgPricePerSqft:  index.AvgPricePerSqft,
				TransactionCount: index.TransactionCount,
			})
		}
		
		series = append(series, response.CompareSeriesData{
			ID:         id,
			Name:       nameMap[id],
			Type:       groupType,
			DataPoints: dataPoints,
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
