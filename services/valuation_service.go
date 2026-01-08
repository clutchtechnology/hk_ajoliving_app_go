package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// ValuationService 估价服务
type ValuationService struct {
	repo *databases.ValuationRepo
}

// NewValuationService 创建估价服务
func NewValuationService(repo *databases.ValuationRepo) *ValuationService {
	return &ValuationService{repo: repo}
}

// ListValuations 获取屋苑估价列表
func (s *ValuationService) ListValuations(ctx context.Context, filter *models.ListValuationsRequest) (*models.PaginatedValuationsResponse, error) {
	valuations, total, err := s.repo.FindAllValuations(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedValuationsResponse{
		Data:       valuations,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetEstateValuation 获取指定屋苑估价参考
func (s *ValuationService) GetEstateValuation(ctx context.Context, estateID uint) (*models.EstateValuationDetail, error) {
	valuation, err := s.repo.GetEstateValuation(ctx, estateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	return valuation, nil
}

// SearchValuations 搜索屋苑估价
func (s *ValuationService) SearchValuations(ctx context.Context, req *models.SearchValuationsRequest) (*models.PaginatedValuationsResponse, error) {
	valuations, total, err := s.repo.SearchValuations(ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedValuationsResponse{
		Data:       valuations,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetDistrictValuations 获取地区屋苑估价列表
func (s *ValuationService) GetDistrictValuations(ctx context.Context, districtID uint) (*models.DistrictValuationSummary, error) {
	summary, err := s.repo.GetDistrictValuations(ctx, districtID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	return summary, nil
}
