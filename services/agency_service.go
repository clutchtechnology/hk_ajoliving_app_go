package services

import (
	"context"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// AgencyService 代理公司服务接口
//
// AgencyService Methods:
// 0. NewAgencyService(repo databases.AgencyRepository, logger *zap.Logger) -> 注入依赖
// 1. ListAgencies(ctx context.Context, filter *models.ListAgenciesRequest) -> 获取代理公司列表
// 2. GetAgency(ctx context.Context, id uint) -> 获取代理公司详情
// 3. GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) -> 获取代理公司房源列表
// 4. ContactAgency(ctx context.Context, agencyID uint, req *models.ContactAgencyRequest) -> 联系代理公司
// 5. SearchAgencies(ctx context.Context, filter *models.SearchAgenciesRequest) -> 搜索代理公司
type AgencyService interface {
	ListAgencies(ctx context.Context, filter *models.ListAgenciesRequest) ([]*models.AgencyDetail, int64, error)
	GetAgency(ctx context.Context, id uint) (*models.AgencyDetail, error)
	GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) ([]*models.Property, int64, error)
	ContactAgency(ctx context.Context, agencyID uint, req *models.ContactAgencyRequest) (*models.AgencyDetail, error)
	SearchAgencies(ctx context.Context, filter *models.SearchAgenciesRequest) ([]*models.AgencyDetail, int64, error)
}

type agencyService struct {
	repo   databases.AgencyRepository
	logger *zap.Logger
}

// 0. NewAgencyService 创建代理公司服务
func NewAgencyService(repo databases.AgencyRepository, logger *zap.Logger) AgencyService {
	return &agencyService{
		repo:   repo,
		logger: logger,
	}
}

// 1. ListAgencies 获取代理公司列表
func (s *agencyService) ListAgencies(ctx context.Context, filter *models.ListAgenciesRequest) ([]*models.AgencyDetail, int64, error) {
	agencies, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list agencies", zap.Error(err))
		return nil, 0, err
	}
	
	return agencies, total, nil
}

// 2. GetAgency 获取代理公司详情
func (s *agencyService) GetAgency(ctx context.Context, id uint) (*models.AgencyDetail, error) {
	agency, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get agency", zap.Error(err), zap.Uint("agency_id", id))
		return nil, err
	}
	if agency == nil {
		return nil, tools.ErrNotFound
	}
	
	return agency, nil
}

// 3. GetAgencyProperties 获取代理公司房源列表
func (s *agencyService) GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) ([]*models.Property, int64, error) {
	// 先检查代理公司是否存在
	agency, err := s.repo.GetByID(ctx, agencyID)
	if err != nil {
		s.logger.Error("failed to get agency", zap.Error(err), zap.Uint("agency_id", agencyID))
		return nil, 0, err
	}
	if agency == nil {
		return nil, 0, tools.ErrNotFound
	}
	
	// 获取代理公司的房源
	properties, total, err := s.repo.GetAgencyProperties(ctx, agencyID, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get agency properties", zap.Error(err), zap.Uint("agency_id", agencyID))
		return nil, 0, err
	}
	
	return properties, total, nil
}

// 4. ContactAgency 联系代理公司
func (s *agencyService) ContactAgency(ctx context.Context, agencyID uint, req *models.ContactAgencyRequest) (*models.AgencyDetail, error) {
	// 检查代理公司是否存在
	agency, err := s.repo.GetByID(ctx, agencyID)
	if err != nil {
		s.logger.Error("failed to get agency", zap.Error(err), zap.Uint("agency_id", agencyID))
		return nil, err
	}
	if agency == nil {
		return nil, tools.ErrNotFound
	}
	
	// TODO: 这里可以实现联系请求的逻辑，比如：
	// 1. 发送邮件通知代理公司
	// 2. 记录联系请求到数据库
	// 3. 发送短信通知
	// 4. 创建待办事项
	
	s.logger.Info("agency contact request received",
		zap.Uint("agency_id", agencyID),
		zap.String("name", req.Name),
		zap.String("phone", req.Phone),
		)
	
	// 返回成功响应
	return agency, nil
}

// 5. SearchAgencies 搜索代理公司
func (s *agencyService) SearchAgencies(ctx context.Context, filter *models.SearchAgenciesRequest) ([]*models.AgencyDetail, int64, error) {
	agencies, total, err := s.repo.Search(ctx, filter)
	if err != nil {
		s.logger.Error("failed to search agencies", zap.Error(err), zap.String("keyword", filter.Keyword))
		return nil, 0, err
	}
	
	return agencies, total, nil
}

