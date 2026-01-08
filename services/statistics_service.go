package services

import (
	"context"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// StatisticsService 统计服务
type StatisticsService struct {
	repo *databases.StatisticsRepo
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService(repo *databases.StatisticsRepo) *StatisticsService {
	return &StatisticsService{repo: repo}
}

// GetOverviewStatistics 获取总览统计
func (s *StatisticsService) GetOverviewStatistics(ctx context.Context, req *models.GetOverviewStatisticsRequest) (*models.OverviewStatisticsResponse, error) {
	var startDate, endDate *time.Time

	// 解析日期
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	stats, err := s.repo.GetOverviewStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 设置时间范围
	stats.StartDate = req.StartDate
	stats.EndDate = req.EndDate

	return stats, nil
}

// GetPropertyStatistics 获取房产统计
func (s *StatisticsService) GetPropertyStatistics(ctx context.Context, req *models.GetPropertyStatisticsRequest) (*models.PropertyStatisticsResponse, error) {
	var startDate, endDate *time.Time

	// 解析日期
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	stats, err := s.repo.GetPropertyStatistics(ctx, startDate, endDate, req.DistrictID)
	if err != nil {
		return nil, err
	}

	// 设置时间范围
	stats.StartDate = req.StartDate
	stats.EndDate = req.EndDate
	stats.DistrictID = req.DistrictID

	return stats, nil
}

// GetTransactionStatistics 获取成交统计
func (s *StatisticsService) GetTransactionStatistics(ctx context.Context, req *models.GetTransactionStatisticsRequest) (*models.TransactionStatisticsResponse, error) {
	var startDate, endDate *time.Time

	// 解析日期
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	stats, err := s.repo.GetTransactionStatistics(ctx, startDate, endDate, req.DistrictID)
	if err != nil {
		return nil, err
	}

	// 设置时间范围
	stats.StartDate = req.StartDate
	stats.EndDate = req.EndDate
	stats.DistrictID = req.DistrictID

	return stats, nil
}

// GetUserStatistics 获取用户统计
func (s *StatisticsService) GetUserStatistics(ctx context.Context, req *models.GetUserStatisticsRequest) (*models.UserStatisticsResponse, error) {
	var startDate, endDate *time.Time

	// 解析日期
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	stats, err := s.repo.GetUserStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 设置时间范围
	stats.StartDate = req.StartDate
	stats.EndDate = req.EndDate

	return stats, nil
}
