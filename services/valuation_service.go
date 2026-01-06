package services

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"go.uber.org/zap"
)

// ValuationService 物业估价服务接口
type ValuationService interface {
	ListValuations(ctx context.Context, req *models.ListValuationsRequest) ([]*models.ValuationListItemResponse, int64, error)
	GetEstateValuation(ctx context.Context, estateID uint) (*models.ValuationDetailResponse, error)
	SearchValuations(ctx context.Context, req *models.SearchValuationsRequest) ([]*models.ValuationListItemResponse, int64, error)
	GetDistrictValuations(ctx context.Context, districtID uint, page, pageSize int) (*models.DistrictValuationResponse, error)
}

type valuationService struct {
	repo   databases.ValuationRepository
	logger *zap.Logger
}

func NewValuationService(repo databases.ValuationRepository, logger *zap.Logger) ValuationService {
	return &valuationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *valuationService) ListValuations(ctx context.Context, req *models.ListValuationsRequest) ([]*models.ValuationListItemResponse, int64, error) {
	estates, total, err := s.repo.ListValuations(ctx, req)
	if err != nil {
		s.logger.Error("failed to list valuations", zap.Error(err))
		return nil, 0, tools.ErrInternalServer
	}

	result := make([]*models.ValuationListItemResponse, 0, len(estates))
	for _, estate := range estates {
		result = append(result, s.toListItemResponse(estate))
	}

	return result, total, nil
}

func (s *valuationService) GetEstateValuation(ctx context.Context, estateID uint) (*models.ValuationDetailResponse, error) {
	estate, err := s.repo.GetEstateValuation(ctx, estateID)
	if err != nil {
		s.logger.Error("failed to get estate valuation", zap.Uint("estate_id", estateID), zap.Error(err))
		return nil, tools.ErrNotFound
	}

	return s.toDetailResponse(estate), nil
}

func (s *valuationService) SearchValuations(ctx context.Context, req *models.SearchValuationsRequest) ([]*models.ValuationListItemResponse, int64, error) {
	estates, total, err := s.repo.SearchValuations(ctx, req)
	if err != nil {
		s.logger.Error("failed to search valuations", zap.String("keyword", req.Keyword), zap.Error(err))
		return nil, 0, tools.ErrInternalServer
	}

	result := make([]*models.ValuationListItemResponse, 0, len(estates))
	for _, estate := range estates {
		result = append(result, s.toListItemResponse(estate))
	}

	return result, total, nil
}

func (s *valuationService) GetDistrictValuations(ctx context.Context, districtID uint, page, pageSize int) (*models.DistrictValuationResponse, error) {
	// 获取地区统计数据
	statistics, err := s.repo.GetDistrictStatistics(ctx, districtID)
	if err != nil {
		s.logger.Error("failed to get district statistics", zap.Uint("district_id", districtID), zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	// 获取地区内的屋苑列表
	estates, _, err := s.repo.GetDistrictValuations(ctx, districtID, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get district valuations", zap.Uint("district_id", districtID), zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	// 构建响应
	districtName := ""
	if len(estates) > 0 && estates[0].District != nil {
		districtName = estates[0].District.NameZhHant
	}

	estateList := make([]models.ValuationListItemResponse, 0, len(estates))
	for _, estate := range estates {
		estateList = append(estateList, *s.toListItemResponse(estate))
	}

	resp := &models.DistrictValuationResponse{
		DistrictID:           districtID,
		DistrictName:         districtName,
		TotalEstates:         int(statistics["total_estates"].(int64)),
		AvgPricePerSqft:      statistics["avg_price"].(float64),
		MedianPricePerSqft:   statistics["avg_price"].(float64), // 简化处理，使用平均值
		MinPricePerSqft:      statistics["min_price"].(float64),
		MaxPricePerSqft:      statistics["max_price"].(float64),
		TotalTransactions:    int(statistics["total_transactions"].(int64)),
		PriceTrend:           "stable", // TODO: 根据历史数据计算趋势
		PriceTrendPercentage: 0.0,
		Estates:              estateList,
	}
	
	return resp, nil
}

// 转换为列表项响应
func (s *valuationService) toListItemResponse(estate *models.Estate) *models.ValuationListItemResponse {
	resp := &models.ValuationListItemResponse{
		EstateID:                estate.ID,
		EstateName:              estate.Name,
		Address:                 estate.Address,
		DistrictName:            "",
		AvgTransactionPrice:     0,
		RecentTransactionsCount: estate.RecentTransactionsCount,
		PriceTrend:              "stable", // TODO: 根据历史数据计算
		PriceTrendPercentage:    0.0,
	}

	if estate.District != nil {
		resp.DistrictName = estate.District.NameZhHant
	}

	if estate.NameEn != nil {
		resp.EstateNameEn = *estate.NameEn
	}

	if estate.CompletionYear != nil {
		resp.CompletionYear = *estate.CompletionYear
	}

	if estate.AvgTransactionPrice != nil {
		resp.AvgTransactionPrice = *estate.AvgTransactionPrice
		// 假设平均单位面积为 500 平方尺，计算每平方尺价格
		resp.AvgPricePerSqft = *estate.AvgTransactionPrice / 500.0
	}

	return resp
}

// 转换为详细响应
func (s *valuationService) toDetailResponse(estate *models.Estate) *models.ValuationDetailResponse {
	resp := &models.ValuationDetailResponse{
		EstateID:                estate.ID,
		EstateName:              estate.Name,
		Address:                 estate.Address,
		DistrictID:              estate.DistrictID,
		DistrictName:            "",
		RecentTransactionsCount: estate.RecentTransactionsCount,
		PriceUpdatedAt:          nil,
		PriceTrend:              "stable", // TODO: 根据历史数据计算
		PriceTrendPercentage:    0.0,
	}

	if estate.District != nil {
		resp.DistrictName = estate.District.NameZhHant
	}

	if estate.AvgTransactionPriceUpdatedAt != nil {
		updatedAt := estate.AvgTransactionPriceUpdatedAt.Format("2006-01-02 15:04:05")
		resp.PriceUpdatedAt = &updatedAt
	}

	if estate.NameEn != nil {
		resp.EstateNameEn = *estate.NameEn
	}

	if estate.TotalBlocks != nil {
		resp.TotalBlocks = *estate.TotalBlocks
	}

	if estate.TotalUnits != nil {
		resp.TotalUnits = *estate.TotalUnits
	}

	if estate.CompletionYear != nil {
		resp.CompletionYear = *estate.CompletionYear
	}

	if estate.PrimarySchoolNet != nil {
		resp.PrimarySchoolNet = *estate.PrimarySchoolNet
	}

	if estate.SecondarySchoolNet != nil {
		resp.SecondarySchoolNet = *estate.SecondarySchoolNet
	}

	if estate.AvgTransactionPrice != nil {
		resp.AvgTransactionPrice = *estate.AvgTransactionPrice
		// 假设平均单位面积为 500 平方尺，计算每平方尺价格
		resp.AvgPricePerSqft = *estate.AvgTransactionPrice / 500.0
	}

	return resp
}
