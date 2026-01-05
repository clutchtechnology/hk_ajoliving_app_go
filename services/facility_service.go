package services

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	pkgErrors "github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FacilityService Methods:
// 0. NewFacilityService(repo databases.FacilityRepository, logger *zap.Logger) -> 注入依赖
// 1. ListFacilities(ctx context.Context, req *models.ListFacilitiesRequest) -> 获取设施列表
// 2. GetFacility(ctx context.Context, id uint) -> 获取单个设施详情
// 3. CreateFacility(ctx context.Context, req *models.Facility) -> 创建设施
// 4. UpdateFacility(ctx context.Context, id uint, req *models.Facility) -> 更新设施信息
// 5. DeleteFacility(ctx context.Context, id uint) -> 删除设施

// FacilityServiceInterface 定义设施服务接口
type FacilityServiceInterface interface {
	ListFacilities(ctx context.Context, req *models.ListFacilitiesRequest) (*[]models.Facility, error)
	GetFacility(ctx context.Context, id uint) (*models.Facility, error)
	CreateFacility(ctx context.Context, req *models.Facility) (*models.Facility, error)
	UpdateFacility(ctx context.Context, id uint, req *models.Facility) (*models.Facility, error)
	DeleteFacility(ctx context.Context, id uint) error
}

// FacilityService 设施服务
type FacilityService struct {
	repo   databases.FacilityRepository
	logger *zap.Logger
}

// 0. NewFacilityService 构造函数
func NewFacilityService(repo databases.FacilityRepository, logger *zap.Logger) *FacilityService {
	return &FacilityService{
		repo:   repo,
		logger: logger,
	}
}

// 1. ListFacilities 获取设施列表
func (s *FacilityService) ListFacilities(ctx context.Context, req *models.ListFacilitiesRequest) (*[]models.Facility, error) {
	facilities, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list facilities", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	facilityResponses := make([]*models.Facility, len(facilities))
	for i, facility := range facilities {
		facilityResponses[i] = convertToFacilityDetailResponse(facility)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	return &[]models.Facility{
		Facilities: facilityResponses,
		Pagination: &models.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

// 2. GetFacility 获取单个设施详情
func (s *FacilityService) GetFacility(ctx context.Context, id uint) (*models.Facility, error) {
	facility, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("Failed to get facility", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	return convertToFacilityDetailResponse(facility), nil
}

// 3. CreateFacility 创建设施
func (s *FacilityService) CreateFacility(ctx context.Context, req *models.Facility) (*models.Facility, error) {
	facility := &models.Facility{
		NameZhHant: req.NameZhHant,
		NameZhHans: req.NameZhHans,
		NameEn:     req.NameEn,
		Icon:       req.Icon,
		Category:   models.FacilityCategory(req.Category),
	}

	if err := s.repo.Create(ctx, facility); err != nil {
		s.logger.Error("Failed to create facility", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Facility created successfully", zap.Uint("id", facility.ID))
	return convertToFacilityDetailResponse(facility), nil
}

// 4. UpdateFacility 更新设施信息
func (s *FacilityService) UpdateFacility(ctx context.Context, id uint, req *models.Facility) (*models.Facility, error) {
	// 检查设施是否存在
	facility, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("Failed to get facility", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	// 更新字段
	if req.NameZhHant != nil {
		facility.NameZhHant = *req.NameZhHant
	}
	if req.NameZhHans != nil {
		facility.NameZhHans = req.NameZhHans
	}
	if req.NameEn != nil {
		facility.NameEn = req.NameEn
	}
	if req.Icon != nil {
		facility.Icon = req.Icon
	}
	if req.Category != nil {
		facility.Category = models.FacilityCategory(*req.Category)
	}
	if req.SortOrder != nil {
		facility.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(ctx, facility); err != nil {
		s.logger.Error("Failed to update facility", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	s.logger.Info("Facility updated successfully", zap.Uint("id", id))
	return convertToFacilityDetailResponse(facility), nil
}

// 5. DeleteFacility 删除设施
func (s *FacilityService) DeleteFacility(ctx context.Context, id uint) error {
	// 检查设施是否存在
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to check facility existence", zap.Error(err), zap.Uint("id", id))
		return err
	}
	if !exists {
		return pkgErrors.ErrNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete facility", zap.Error(err), zap.Uint("id", id))
		return err
	}

	s.logger.Info("Facility deleted successfully", zap.Uint("id", id))
	return nil
}

// convertToFacilityDetailResponse 将 Facility 模型转换为响应格式
func convertToFacilityDetailResponse(facility *models.Facility) *models.Facility {
	return &models.Facility{
		ID:         facility.ID,
		NameZhHant: facility.NameZhHant,
		NameZhHans: facility.NameZhHans,
		NameEn:     facility.NameEn,
		Icon:       facility.Icon,
		Category:   string(facility.Category),
	}
}
