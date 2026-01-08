package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// EstateService 屋苑服务
type EstateService struct {
	repo *databases.EstateRepo
}

// NewEstateService 创建屋苑服务
func NewEstateService(repo *databases.EstateRepo) *EstateService {
	return &EstateService{repo: repo}
}

// ListEstates 获取屋苑列表
func (s *EstateService) ListEstates(ctx context.Context, filter *models.ListEstatesRequest) (*models.PaginatedEstatesResponse, error) {
	estates, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	data := make([]models.EstateResponse, len(estates))
	for i, estate := range estates {
		data[i] = s.toEstateResponse(&estate)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedEstatesResponse{
		Data:       data,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetEstate 获取屋苑详情
func (s *EstateService) GetEstate(ctx context.Context, id uint) (*models.EstateResponse, error) {
	estate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 增加浏览次数
	_ = s.repo.IncrementViewCount(ctx, id)

	response := s.toEstateResponse(estate)
	return &response, nil
}

// GetFeaturedEstates 获取精选屋苑
func (s *EstateService) GetFeaturedEstates(ctx context.Context, limit int) ([]models.EstateResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	estates, err := s.repo.FindFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}

	response := make([]models.EstateResponse, len(estates))
	for i, estate := range estates {
		response[i] = s.toEstateResponse(&estate)
	}

	return response, nil
}

// GetEstateProperties 获取屋苑内房源列表
func (s *EstateService) GetEstateProperties(ctx context.Context, id uint, filter *models.GetEstatePropertiesRequest) (map[string]interface{}, error) {
	// 先查询屋苑是否存在
	estate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 查询屋苑内的房源
	properties, total, err := s.repo.FindPropertiesByEstate(ctx, estate.Name, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return map[string]interface{}{
		"estate_id":   id,
		"estate_name": estate.Name,
		"data":        properties,
		"total":       total,
		"page":        filter.Page,
		"page_size":   filter.PageSize,
		"total_pages": totalPages,
	}, nil
}

// GetEstateImages 获取屋苑图片
func (s *EstateService) GetEstateImages(ctx context.Context, id uint) ([]models.EstateImage, error) {
	// 先查询屋苑是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	images, err := s.repo.FindImagesByEstateID(ctx, id)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// GetEstateFacilities 获取屋苑设施
func (s *EstateService) GetEstateFacilities(ctx context.Context, id uint) ([]models.Facility, error) {
	// 先查询屋苑是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	facilities, err := s.repo.FindFacilitiesByEstateID(ctx, id)
	if err != nil {
		return nil, err
	}

	return facilities, nil
}

// GetEstateTransactions 获取屋苑成交记录
// TODO: 待实现 transactions 模块后补充
func (s *EstateService) GetEstateTransactions(ctx context.Context, id uint, page, pageSize int) (map[string]interface{}, error) {
	// 先查询屋苑是否存在
	estate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// TODO: 查询成交记录
	// transactions, total, err := s.transactionRepo.FindByEstateName(ctx, estate.Name, page, pageSize)

	return map[string]interface{}{
		"estate_id":   id,
		"estate_name": estate.Name,
		"data":        []interface{}{}, // TODO: 实际成交记录
		"total":       0,
		"page":        page,
		"page_size":   pageSize,
		"message":     "Transaction module not implemented yet",
	}, nil
}

// GetEstateStatistics 获取屋苑统计数据
func (s *EstateService) GetEstateStatistics(ctx context.Context, id uint) (*models.EstateStatisticsResponse, error) {
	// 先查询屋苑是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	stats, err := s.repo.GetStatistics(ctx, id)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// CreateEstate 创建屋苑
func (s *EstateService) CreateEstate(ctx context.Context, req *models.CreateEstateRequest) (*models.EstateResponse, error) {
	estate := &models.Estate{
		Name:               req.Name,
		NameEn:             req.NameEn,
		Address:            req.Address,
		DistrictID:         req.DistrictID,
		TotalBlocks:        req.TotalBlocks,
		TotalUnits:         req.TotalUnits,
		CompletionYear:     req.CompletionYear,
		Developer:          req.Developer,
		ManagementCompany:  req.ManagementCompany,
		PrimarySchoolNet:   req.PrimarySchoolNet,
		SecondarySchoolNet: req.SecondarySchoolNet,
		Description:        req.Description,
		IsFeatured:         req.IsFeatured,
	}

	if err := s.repo.Create(ctx, estate); err != nil {
		return nil, err
	}

	// 如果有设施，更新设施关联
	if len(req.FacilityIDs) > 0 {
		if err := s.repo.UpdateFacilities(ctx, estate.ID, req.FacilityIDs); err != nil {
			return nil, err
		}
	}

	// 重新查询以获取完整数据
	createdEstate, err := s.repo.FindByID(ctx, estate.ID)
	if err != nil {
		return nil, err
	}

	response := s.toEstateResponse(createdEstate)
	return &response, nil
}

// UpdateEstate 更新屋苑
func (s *EstateService) UpdateEstate(ctx context.Context, id uint, req *models.UpdateEstateRequest) (*models.EstateResponse, error) {
	// 查询屋苑是否存在
	estate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		estate.Name = *req.Name
	}
	if req.NameEn != nil {
		estate.NameEn = *req.NameEn
	}
	if req.Address != nil {
		estate.Address = *req.Address
	}
	if req.DistrictID != nil {
		estate.DistrictID = *req.DistrictID
	}
	if req.TotalBlocks != nil {
		estate.TotalBlocks = *req.TotalBlocks
	}
	if req.TotalUnits != nil {
		estate.TotalUnits = *req.TotalUnits
	}
	if req.CompletionYear != nil {
		estate.CompletionYear = *req.CompletionYear
	}
	if req.Developer != nil {
		estate.Developer = *req.Developer
	}
	if req.ManagementCompany != nil {
		estate.ManagementCompany = *req.ManagementCompany
	}
	if req.PrimarySchoolNet != nil {
		estate.PrimarySchoolNet = *req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != nil {
		estate.SecondarySchoolNet = *req.SecondarySchoolNet
	}
	if req.Description != nil {
		estate.Description = *req.Description
	}
	if req.IsFeatured != nil {
		estate.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.Update(ctx, estate); err != nil {
		return nil, err
	}

	// 更新设施关联
	if req.FacilityIDs != nil {
		if err := s.repo.UpdateFacilities(ctx, id, req.FacilityIDs); err != nil {
			return nil, err
		}
	}

	// 重新查询以获取完整数据
	updatedEstate, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.toEstateResponse(updatedEstate)
	return &response, nil
}

// DeleteEstate 删除屋苑
func (s *EstateService) DeleteEstate(ctx context.Context, id uint) error {
	// 查询屋苑是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tools.ErrNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

// toEstateResponse 转换为响应格式
func (s *EstateService) toEstateResponse(estate *models.Estate) models.EstateResponse {
	return models.EstateResponse{
		ID:                      estate.ID,
		Name:                    estate.Name,
		NameEn:                  estate.NameEn,
		Address:                 estate.Address,
		DistrictID:              estate.DistrictID,
		District:                estate.District,
		TotalBlocks:             estate.TotalBlocks,
		TotalUnits:              estate.TotalUnits,
		CompletionYear:          estate.CompletionYear,
		Developer:               estate.Developer,
		ManagementCompany:       estate.ManagementCompany,
		PrimarySchoolNet:        estate.PrimarySchoolNet,
		SecondarySchoolNet:      estate.SecondarySchoolNet,
		RecentTransactionsCount: estate.RecentTransactionsCount,
		ForSaleCount:            estate.ForSaleCount,
		ForRentCount:            estate.ForRentCount,
		AvgTransactionPrice:     estate.AvgTransactionPrice,
		Description:             estate.Description,
		ViewCount:               estate.ViewCount,
		FavoriteCount:           estate.FavoriteCount,
		IsFeatured:              estate.IsFeatured,
		Images:                  estate.Images,
		Facilities:              estate.Facilities,
		CreatedAt:               estate.CreatedAt,
		UpdatedAt:               estate.UpdatedAt,
	}
}
