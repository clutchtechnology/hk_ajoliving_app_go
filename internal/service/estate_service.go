package service

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
	"go.uber.org/zap"
)

// EstateService 屋苑服务接口
type EstateService interface {
	ListEstates(ctx context.Context, req *request.ListEstatesRequest) ([]*response.EstateListItemResponse, int64, error)
	GetEstate(ctx context.Context, id uint) (*response.EstateResponse, error)
	GetEstateProperties(ctx context.Context, id uint, listingType string, page, pageSize int) ([]*model.Property, int64, error)
	GetEstateStatistics(ctx context.Context, id uint) (*response.EstateStatisticsResponse, error)
	GetFeaturedEstates(ctx context.Context, limit int) ([]*response.EstateListItemResponse, error)
	CreateEstate(ctx context.Context, req *request.CreateEstateRequest) (*response.EstateResponse, error)
	UpdateEstate(ctx context.Context, id uint, req *request.UpdateEstateRequest) (*response.EstateResponse, error)
	DeleteEstate(ctx context.Context, id uint) error
}

type estateService struct {
	repo   repository.EstateRepository
	logger *zap.Logger
}

func NewEstateService(repo repository.EstateRepository, logger *zap.Logger) EstateService {
	return &estateService{
		repo:   repo,
		logger: logger,
	}
}

func (s *estateService) ListEstates(ctx context.Context, req *request.ListEstatesRequest) ([]*response.EstateListItemResponse, int64, error) {
	estates, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("failed to list estates", zap.Error(err))
		return nil, 0, errors.ErrInternalServer
	}

	result := make([]*response.EstateListItemResponse, 0, len(estates))
	for _, estate := range estates {
		result = append(result, s.toListItemResponse(estate))
	}

	return result, total, nil
}

func (s *estateService) GetEstate(ctx context.Context, id uint) (*response.EstateResponse, error) {
	estate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get estate", zap.Uint("id", id), zap.Error(err))
		return nil, errors.ErrNotFound
	}

	// 增加浏览次数
	go func() {
		if err := s.repo.IncrementViewCount(context.Background(), id); err != nil {
			s.logger.Warn("failed to increment view count", zap.Uint("id", id), zap.Error(err))
		}
	}()

	return s.toDetailResponse(estate), nil
}

func (s *estateService) GetEstateProperties(ctx context.Context, id uint, listingType string, page, pageSize int) ([]*model.Property, int64, error) {
	properties, total, err := s.repo.GetProperties(ctx, id, listingType, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get estate properties", zap.Uint("estate_id", id), zap.Error(err))
		return nil, 0, errors.ErrInternalServer
	}

	return properties, total, nil
}

func (s *estateService) GetEstateStatistics(ctx context.Context, id uint) (*response.EstateStatisticsResponse, error) {
	estate, err := s.repo.GetStatistics(ctx, id)
	if err != nil {
		s.logger.Error("failed to get estate statistics", zap.Uint("id", id), zap.Error(err))
		return nil, errors.ErrNotFound
	}

	return &response.EstateStatisticsResponse{
		EstateID:                   estate.ID,
		EstateName:                 estate.Name,
		RecentTransactionsCount:    estate.RecentTransactionsCount,
		ForSaleCount:               estate.ForSaleCount,
		ForRentCount:               estate.ForRentCount,
		AvgTransactionPrice:        estate.AvgTransactionPrice,
		AvgTransactionPriceUpdated: estate.AvgTransactionPriceUpdatedAt,
	}, nil
}

func (s *estateService) GetFeaturedEstates(ctx context.Context, limit int) ([]*response.EstateListItemResponse, error) {
	estates, err := s.repo.GetFeatured(ctx, limit)
	if err != nil {
		s.logger.Error("failed to get featured estates", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	result := make([]*response.EstateListItemResponse, 0, len(estates))
	for _, estate := range estates {
		result = append(result, s.toListItemResponse(estate))
	}

	return result, nil
}

func (s *estateService) CreateEstate(ctx context.Context, req *request.CreateEstateRequest) (*response.EstateResponse, error) {
	estate := &model.Estate{
		Name:                req.Name,
		Description:         req.Description,
		Address:             req.Address,
		DistrictID:          req.DistrictID,
		TotalBlocks:         req.TotalBlocks,
		TotalUnits:          req.TotalUnits,
		CompletionYear:      req.CompletionYear,
		PrimarySchoolNet:    req.PrimarySchoolNet,
		SecondarySchoolNet:  req.SecondarySchoolNet,
		Developer:           req.Developer,
		ManagementCompany:   req.ManagementCompany,
		ManagementFee:       req.ManagementFee,
	}

	if err := s.repo.Create(ctx, estate); err != nil {
		s.logger.Error("failed to create estate", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// TODO: 处理图片、设施上传

	return s.GetEstate(ctx, estate.ID)
}

func (s *estateService) UpdateEstate(ctx context.Context, id uint, req *request.UpdateEstateRequest) (*response.EstateResponse, error) {
	estate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	// 更新字段
	if req.Name != nil {
		estate.Name = *req.Name
	}
	if req.Description != nil {
		estate.Description = *req.Description
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
	if req.PrimarySchoolNet != nil {
		estate.PrimarySchoolNet = *req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != nil {
		estate.SecondarySchoolNet = *req.SecondarySchoolNet
	}
	if req.Developer != nil {
		estate.Developer = *req.Developer
	}
	if req.ManagementCompany != nil {
		estate.ManagementCompany = *req.ManagementCompany
	}
	if req.ManagementFee != nil {
		estate.ManagementFee = *req.ManagementFee
	}
	if req.IsFeatured != nil {
		estate.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.Update(ctx, estate); err != nil {
		s.logger.Error("failed to update estate", zap.Uint("id", id), zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return s.GetEstate(ctx, id)
}

func (s *estateService) DeleteEstate(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete estate", zap.Uint("id", id), zap.Error(err))
		return errors.ErrInternalServer
	}
	return nil
}

// 转换为列表项响应
func (s *estateService) toListItemResponse(estate *model.Estate) *response.EstateListItemResponse {
	resp := &response.EstateListItemResponse{
		ID:                      estate.ID,
		Name:                    estate.Name,
		Address:                 estate.Address,
		DistrictName:            "",
		TotalBlocks:             estate.TotalBlocks,
		TotalUnits:              estate.TotalUnits,
		CompletionYear:          estate.CompletionYear,
		PrimarySchoolNet:        estate.PrimarySchoolNet,
		SecondarySchoolNet:      estate.SecondarySchoolNet,
		RecentTransactionsCount: estate.RecentTransactionsCount,
		ForSaleCount:            estate.ForSaleCount,
		ForRentCount:            estate.ForRentCount,
		AvgTransactionPrice:     estate.AvgTransactionPrice,
		ViewCount:               estate.ViewCount,
		IsFeatured:              estate.IsFeatured,
		CoverImage:              "",
	}

	if estate.District != nil {
		resp.DistrictName = estate.District.Name
	}

	if len(estate.Images) > 0 {
		resp.CoverImage = estate.Images[0].ImageURL
	}

	return resp
}

// 转换为详细响应
func (s *estateService) toDetailResponse(estate *model.Estate) *response.EstateResponse {
	resp := &response.EstateResponse{
		ID:                      estate.ID,
		Name:                    estate.Name,
		Description:             estate.Description,
		Address:                 estate.Address,
		DistrictID:              estate.DistrictID,
		DistrictName:            "",
		TotalBlocks:             estate.TotalBlocks,
		TotalUnits:              estate.TotalUnits,
		CompletionYear:          estate.CompletionYear,
		PrimarySchoolNet:        estate.PrimarySchoolNet,
		SecondarySchoolNet:      estate.SecondarySchoolNet,
		Developer:               estate.Developer,
		ManagementCompany:       estate.ManagementCompany,
		ManagementFee:           estate.ManagementFee,
		RecentTransactionsCount: estate.RecentTransactionsCount,
		ForSaleCount:            estate.ForSaleCount,
		ForRentCount:            estate.ForRentCount,
		AvgTransactionPrice:     estate.AvgTransactionPrice,
		ViewCount:               estate.ViewCount,
		IsFeatured:              estate.IsFeatured,
		CreatedAt:               estate.CreatedAt,
		UpdatedAt:               estate.UpdatedAt,
		Images:                  []response.EstateImageResponse{},
		Facilities:              []string{},
	}

	if estate.District != nil {
		resp.DistrictName = estate.District.Name
	}

	// 图片
	for _, img := range estate.Images {
		resp.Images = append(resp.Images, response.EstateImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			ImageType: img.ImageType,
			SortOrder: img.SortOrder,
		})
	}

	// 设施
	for _, facility := range estate.Facilities {
		resp.Facilities = append(resp.Facilities, facility.Name)
	}

	return resp
}
