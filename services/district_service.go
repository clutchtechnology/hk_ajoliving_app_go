package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// DistrictService Methods:
// 0. NewDistrictService(repo *databases.DistrictRepo) -> 注入依赖
// 1. ListDistricts(ctx context.Context, region string) -> 获取地区列表
// 2. GetDistrict(ctx context.Context, id uint) -> 获取地区详情
// 3. GetDistrictProperties(ctx context.Context, id uint, filter *models.GetDistrictPropertiesRequest) -> 获取地区房源
// 4. GetDistrictEstates(ctx context.Context, id uint, page, pageSize int) -> 获取地区屋苑
// 5. GetDistrictStatistics(ctx context.Context, id uint) -> 获取地区统计数据

type DistrictService struct {
	repo *databases.DistrictRepo
}

// 0. NewDistrictService 构造函数
func NewDistrictService(repo *databases.DistrictRepo) *DistrictService {
	return &DistrictService{repo: repo}
}

// 1. ListDistricts 获取地区列表
func (s *DistrictService) ListDistricts(ctx context.Context, region string) ([]models.DistrictResponse, error) {
	districts, err := s.repo.FindAll(ctx, region)
	if err != nil {
		return nil, err
	}

	var responses []models.DistrictResponse
	for _, district := range districts {
		responses = append(responses, models.DistrictResponse{
			ID:         district.ID,
			NameZhHant: district.NameZhHant,
			NameZhHans: district.NameZhHans,
			NameEn:     district.NameEn,
			Region:     district.Region,
			SortOrder:  district.SortOrder,
		})
	}

	return responses, nil
}

// 2. GetDistrict 获取地区详情
func (s *DistrictService) GetDistrict(ctx context.Context, id uint) (*models.DistrictDetailResponse, error) {
	district, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 获取统计数据
	propertyCount, _ := s.repo.GetPropertyCount(ctx, id)
	estateCount, _ := s.repo.GetEstateCount(ctx, id)
	avgPrice, _ := s.repo.GetAvgPropertyPrice(ctx, id)

	return &models.DistrictDetailResponse{
		ID:               district.ID,
		NameZhHant:       district.NameZhHant,
		NameZhHans:       district.NameZhHans,
		NameEn:           district.NameEn,
		Region:           district.Region,
		PropertyCount:    propertyCount,
		EstateCount:      estateCount,
		AvgPropertyPrice: avgPrice,
		SortOrder:        district.SortOrder,
	}, nil
}

// 3. GetDistrictProperties 获取地区房源
func (s *DistrictService) GetDistrictProperties(ctx context.Context, id uint, filter *models.GetDistrictPropertiesRequest) ([]models.Property, int64, error) {
	// 验证地区存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, 0, tools.ErrNotFound
		}
		return nil, 0, err
	}

	properties, total, err := s.repo.GetDistrictProperties(ctx, id, filter)
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// 4. GetDistrictEstates 获取地区屋苑
func (s *DistrictService) GetDistrictEstates(ctx context.Context, id uint, page, pageSize int) ([]models.Estate, int64, error) {
	// 验证地区存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, 0, tools.ErrNotFound
		}
		return nil, 0, err
	}

	estates, total, err := s.repo.GetDistrictEstates(ctx, id, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

// 5. GetDistrictStatistics 获取地区统计数据
func (s *DistrictService) GetDistrictStatistics(ctx context.Context, id uint) (*models.DistrictStatisticsResponse, error) {
	stats, err := s.repo.GetDistrictStatistics(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	return stats, nil
}
