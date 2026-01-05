package services

import (
	"time"
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"go.uber.org/zap"
)

// EstateService 屋苑服务接口
type EstateService interface {
	ListEstates(ctx context.Context, req *models.ListEstatesRequest) ([]*models.Estate, int64, error)
	GetEstate(ctx context.Context, id uint) (*models.Estate, error)
	GetEstateProperties(ctx context.Context, id uint, listingType string, page, pageSize int) ([]*models.Property, int64, error)
	GetEstateStatistics(ctx context.Context, id uint) (*map[string]interface{}, error)
	GetFeaturedEstates(ctx context.Context, limit int) ([]*models.Estate, error)
	CreateEstate(ctx context.Context, req *models.Estate) (*models.Estate, error)
	UpdateEstate(ctx context.Context, id uint, req *models.Estate) (*models.Estate, error)
	DeleteEstate(ctx context.Context, id uint) error
}

type estateService struct {
	repo   databases.EstateRepository
	logger *zap.Logger
}

func NewEstateService(repo databases.EstateRepository, logger *zap.Logger) EstateService {
	return &estateService{
		repo:   repo,
		logger: logger,
	}
}

func (s *estateService) ListEstates(ctx context.Context, req *models.ListEstatesRequest) ([]*models.Estate, int64, error) {
	estates, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("failed to list estates", zap.Error(err))
		return nil, 0, tools.ErrInternalServer
	}

	result := make([]*models.Estate, 0, len(estates))
	for _, estate := range estates {
		result = append(result, s.toListItemResponse(estate))
	}

	return result, total, nil
}

func (s *estateService) GetEstate(ctx context.Context, id uint) (*models.Estate, error) {
	estate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get estate", zap.Uint("id", id), zap.Error(err))
		return nil, tools.ErrNotFound
	}

	// 增加浏览次数
	go func() {
		if err := s.repo.IncrementViewCount(context.Background(), id); err != nil {
			s.logger.Warn("failed to increment view count", zap.Uint("id", id), zap.Error(err))
		}
	}()

	return s.toDetailResponse(estate), nil
}

func (s *estateService) GetEstateProperties(ctx context.Context, id uint, listingType string, page, pageSize int) ([]*models.Property, int64, error) {
	properties, total, err := s.repo.GetProperties(ctx, id, listingType, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get estate properties", zap.Uint("estate_id", id), zap.Error(err))
		return nil, 0, tools.ErrInternalServer
	}

	return properties, total, nil
}

func (s *estateService) GetEstateStatistics(ctx context.Context, id uint) (*map[string]interface{}, error) {
	estate, err := s.repo.GetStatistics(ctx, id)
	if err != nil {
		s.logger.Error("failed to get estate statistics", zap.Uint("id", id), zap.Error(err))
		return nil, tools.ErrNotFound
	}

	var avgPrice float64
	if estate.AvgTransactionPrice != nil {
		avgPrice = *estate.AvgTransactionPrice
	}

	return &map[string]interface{}{
		EstateID:                estate.ID,
		EstateName:              estate.Name,
		RecentTransactions:      estate.RecentTransactionsCount,
		ForSaleCount:            estate.ForSaleCount,
		ForRentCount:            estate.ForRentCount,
		AvgTransactionPrice:     avgPrice,
		LastTransactionDate:     estate.AvgTransactionPriceUpdatedAt,
	}, nil
}

func (s *estateService) GetFeaturedEstates(ctx context.Context, limit int) ([]*models.Estate, error) {
	estates, err := s.repo.GetFeatured(ctx, limit)
	if err != nil {
		s.logger.Error("failed to get featured estates", zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	result := make([]*models.Estate, 0, len(estates))
	for _, estate := range estates {
		result = append(result, s.toListItemResponse(estate))
	}

	return result, nil
}

func (s *estateService) CreateEstate(ctx context.Context, req *models.Estate) (*models.Estate, error) {
	estate := &models.Estate{
		Name:       req.Name,
		Address:    req.Address,
		DistrictID: req.DistrictID,
		IsFeatured: req.IsFeatured,
	}

	// 处理可选字段
	if req.Description != "" {
		estate.Description = &req.Description
	}
	if req.TotalBlocks > 0 {
		estate.TotalBlocks = &req.TotalBlocks
	}
	if req.TotalUnits > 0 {
		estate.TotalUnits = &req.TotalUnits
	}
	if req.CompletionYear > 0 {
		estate.CompletionYear = &req.CompletionYear
	}
	if req.Developer != "" {
		estate.Developer = &req.Developer
	}
	if req.ManagementCompany != "" {
		estate.ManagementCompany = &req.ManagementCompany
	}
	if req.PrimarySchoolNet != "" {
		estate.PrimarySchoolNet = &req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != "" {
		estate.SecondarySchoolNet = &req.SecondarySchoolNet
	}

	if err := s.repo.Create(ctx, estate); err != nil {
		s.logger.Error("failed to create estate", zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	// TODO: 处理图片、设施上传

	return s.GetEstate(ctx, estate.ID)
}

func (s *estateService) UpdateEstate(ctx context.Context, id uint, req *models.Estate) (*models.Estate, error) {
	estate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, tools.ErrNotFound
	}

	// 更新字段
	if req.Name != "" {
		estate.Name = req.Name
	}
	if req.Description != "" {
		estate.Description = &req.Description
	}
	if req.Address != "" {
		estate.Address = req.Address
	}
	if req.DistrictID > 0 {
		estate.DistrictID = req.DistrictID
	}
	if req.TotalBlocks > 0 {
		estate.TotalBlocks = &req.TotalBlocks
	}
	if req.TotalUnits > 0 {
		estate.TotalUnits = &req.TotalUnits
	}
	if req.CompletionYear > 0 {
		estate.CompletionYear = &req.CompletionYear
	}
	if req.PrimarySchoolNet != "" {
		estate.PrimarySchoolNet = &req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != "" {
		estate.SecondarySchoolNet = &req.SecondarySchoolNet
	}
	if req.Developer != "" {
		estate.Developer = &req.Developer
	}
	if req.ManagementCompany != "" {
		estate.ManagementCompany = &req.ManagementCompany
	}
	if req.IsFeatured != nil {
		estate.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.Update(ctx, estate); err != nil {
		s.logger.Error("failed to update estate", zap.Uint("id", id), zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	return s.GetEstate(ctx, id)
}

func (s *estateService) DeleteEstate(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete estate", zap.Uint("id", id), zap.Error(err))
		return tools.ErrInternalServer
	}
	return nil
}

// 转换为列表项响应
func (s *estateService) toListItemResponse(estate *models.Estate) *models.Estate {
	var avgPrice float64
	if estate.AvgTransactionPrice != nil {
		avgPrice = *estate.AvgTransactionPrice
	}

	var completionYear, age int
	if estate.CompletionYear != nil {
		completionYear = *estate.CompletionYear
		age = time.Now().Year() - completionYear
	}

	resp := &models.Estate{
		ID:                      estate.ID,
		Name:                    estate.Name,
		Address:                 estate.Address,
		DistrictID:              estate.DistrictID,
		CompletionYear:          completionYear,
		Age:                     age,
		RecentTransactionsCount: estate.RecentTransactionsCount,
		ForSaleCount:            estate.ForSaleCount,
		ForRentCount:            estate.ForRentCount,
		AvgTransactionPrice:     avgPrice,
		ViewCount:               estate.ViewCount,
		FavoriteCount:           estate.FavoriteCount,
		IsFeatured:              estate.IsFeatured,
		CreatedAt:               estate.CreatedAt,
	}

	// 设置地区名称
	if estate.NameEn != nil {
		resp.NameEn = *estate.NameEn
	}

	// 设置封面图片
	for _, img := range estate.Images {
		// 使用第一张图片作为封面
		resp.CoverImage = img.ImageURL
		break
	}

	return resp
}

// 转换为详细响应
func (s *estateService) toDetailResponse(estate *models.Estate) *models.Estate {
	var avgPrice float64
	if estate.AvgTransactionPrice != nil {
		avgPrice = *estate.AvgTransactionPrice
	}

	var totalBlocks, totalUnits, completionYear, age int
	var nameEn, developer, managementCompany, primarySchoolNet, secondarySchoolNet, description string

	if estate.TotalBlocks != nil {
		totalBlocks = *estate.TotalBlocks
	}
	if estate.TotalUnits != nil {
		totalUnits = *estate.TotalUnits
	}
	if estate.CompletionYear != nil {
		completionYear = *estate.CompletionYear
		age = time.Now().Year() - completionYear
	}
	if estate.NameEn != nil {
		nameEn = *estate.NameEn
	}
	if estate.Developer != nil {
		developer = *estate.Developer
	}
	if estate.ManagementCompany != nil {
		managementCompany = *estate.ManagementCompany
	}
	if estate.PrimarySchoolNet != nil {
		primarySchoolNet = *estate.PrimarySchoolNet
	}
	if estate.SecondarySchoolNet != nil {
		secondarySchoolNet = *estate.SecondarySchoolNet
	}
	if estate.Description != nil {
		description = *estate.Description
	}

	resp := &models.Estate{
		ID:                           estate.ID,
		Name:                         estate.Name,
		NameEn:                       nameEn,
		Address:                      estate.Address,
		DistrictID:                   estate.DistrictID,
		TotalBlocks:                  totalBlocks,
		TotalUnits:                   totalUnits,
		CompletionYear:               completionYear,
		Age:                          age,
		Developer:                    developer,
		ManagementCompany:            managementCompany,
		PrimarySchoolNet:             primarySchoolNet,
		SecondarySchoolNet:           secondarySchoolNet,
		RecentTransactionsCount:      estate.RecentTransactionsCount,
		ForSaleCount:                 estate.ForSaleCount,
		ForRentCount:                 estate.ForRentCount,
		AvgTransactionPrice:          avgPrice,
		AvgTransactionPriceUpdatedAt: estate.AvgTransactionPriceUpdatedAt,
		Description:                  description,
		ViewCount:                    estate.ViewCount,
		FavoriteCount:                estate.FavoriteCount,
		IsFeatured:                   estate.IsFeatured,
		CreatedAt:                    estate.CreatedAt,
		UpdatedAt:                    estate.UpdatedAt,
		Images:                       []models.EstateImageResponse{},
		Facilities:                   []models.Facility{},
	}

	// 图片
	for _, img := range estate.Images {
		resp.Images = append(resp.Images, models.EstateImageResponse{
			ID:        img.ID,
			URL:       img.ImageURL,
			ImageType: string(img.ImageType),
			SortOrder: img.SortOrder,
		})
	}

	// 设施
	for _, facility := range estate.Facilities {
		resp.Facilities = append(resp.Facilities, models.Facility{
			ID:   facility.ID,
			NameZhHant: facility.NameZhHant,
			Icon: facility.Icon,
		})
	}

	return resp
}
