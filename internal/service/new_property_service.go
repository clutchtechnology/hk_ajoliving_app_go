package service

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	pkgErrors "github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewPropertyService 新楼盘服务接口
// Methods:
// 0. NewNewPropertyService(repo, logger) -> 注入依赖
// 1. ListNewDevelopments(ctx, req) -> 获取新楼盘列表
// 2. GetNewDevelopment(ctx, id) -> 获取新楼盘详情
// 3. GetDevelopmentUnits(ctx, id) -> 获取楼盘户型列表
// 4. GetFeaturedNewDevelopments(ctx, limit) -> 获取精选新楼盘
type NewPropertyService interface {
	ListNewDevelopments(ctx context.Context, req *request.ListNewDevelopmentsRequest) ([]*response.NewDevelopmentListItemResponse, int64, error)
	GetNewDevelopment(ctx context.Context, id uint) (*response.NewDevelopmentResponse, error)
	GetDevelopmentUnits(ctx context.Context, id uint) ([]response.NewDevelopmentLayoutResponse, error)
	GetFeaturedNewDevelopments(ctx context.Context, limit int) ([]*response.NewDevelopmentListItemResponse, error)
}

// newPropertyService 新楼盘服务实现
type newPropertyService struct {
	repo   repository.NewPropertyRepository
	logger *zap.Logger
}

// NewNewPropertyService 创建新楼盘服务实例
func NewNewPropertyService(repo repository.NewPropertyRepository, logger *zap.Logger) NewPropertyService {
	return &newPropertyService{
		repo:   repo,
		logger: logger,
	}
}

// ListNewDevelopments 获取新楼盘列表
func (s *newPropertyService) ListNewDevelopments(ctx context.Context, req *request.ListNewDevelopmentsRequest) ([]*response.NewDevelopmentListItemResponse, int64, error) {
	// 获取新楼盘列表
	newProperties, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list new developments", zap.Error(err))
		return nil, 0, err
	}

	// 转换为响应格式
	items := make([]*response.NewDevelopmentListItemResponse, len(newProperties))
	for i, np := range newProperties {
		items[i] = s.toNewDevelopmentListItemResponse(np)
	}

	return items, total, nil
}

// GetNewDevelopment 获取新楼盘详情
func (s *newPropertyService) GetNewDevelopment(ctx context.Context, id uint) (*response.NewDevelopmentResponse, error) {
	// 获取新楼盘
	np, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("Failed to get new development", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// 增加浏览次数
	if err := s.repo.IncrementViewCount(ctx, id); err != nil {
		s.logger.Warn("Failed to increment view count", zap.Uint("id", id), zap.Error(err))
	}

	return s.toNewDevelopmentResponse(np), nil
}

// GetDevelopmentUnits 获取楼盘户型列表
func (s *newPropertyService) GetDevelopmentUnits(ctx context.Context, id uint) ([]response.NewDevelopmentLayoutResponse, error) {
	// 检查新楼盘是否存在
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("Failed to get new development", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// 获取户型列表
	layouts, err := s.repo.GetLayouts(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get development units", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	result := make([]response.NewDevelopmentLayoutResponse, len(layouts))
	for i, layout := range layouts {
		result[i] = s.toLayoutResponse(&layout)
	}

	return result, nil
}

// GetFeaturedNewDevelopments 获取精选新楼盘
func (s *newPropertyService) GetFeaturedNewDevelopments(ctx context.Context, limit int) ([]*response.NewDevelopmentListItemResponse, error) {
	newProperties, err := s.repo.GetFeatured(ctx, limit)
	if err != nil {
		s.logger.Error("Failed to get featured new developments", zap.Error(err))
		return nil, err
	}

	items := make([]*response.NewDevelopmentListItemResponse, len(newProperties))
	for i, np := range newProperties {
		items[i] = s.toNewDevelopmentListItemResponse(np)
	}

	return items, nil
}

// toNewDevelopmentResponse 转换为新楼盘详情响应
func (s *newPropertyService) toNewDevelopmentResponse(np *model.NewProperty) *response.NewDevelopmentResponse {
	resp := &response.NewDevelopmentResponse{
		ID:                 np.ID,
		Name:               np.Name,
		NameEn:             derefString(np.NameEn),
		Address:            np.Address,
		DistrictID:         np.DistrictID,
		Status:             string(np.Status),
		UnitsForSale:       derefInt(np.UnitsForSale),
		UnitsSold:          derefInt(np.UnitsSold),
		Developer:          np.Developer,
		ManagementCompany:  derefString(np.ManagementCompany),
		TotalUnits:         np.TotalUnits,
		TotalBlocks:        np.TotalBlocks,
		MaxFloors:          np.MaxFloors,
		PrimarySchoolNet:   derefString(np.PrimarySchoolNet),
		SecondarySchoolNet: derefString(np.SecondarySchoolNet),
		WebsiteURL:         derefString(np.WebsiteURL),
		SalesOfficeAddress: derefString(np.SalesOfficeAddress),
		SalesPhone:         derefString(np.SalesPhone),
		ExpectedCompletion: np.ExpectedCompletion,
		OccupationDate:     np.OccupationDate,
		Description:        derefString(np.Description),
		ViewCount:          np.ViewCount,
		FavoriteCount:      np.FavoriteCount,
		IsFeatured:         np.IsFeatured,
		SalesProgress:      np.GetSalesProgress(),
		CreatedAt:          np.CreatedAt,
		UpdatedAt:          np.UpdatedAt,
	}

	// 地区信息
	if np.District != nil {
		resp.District = &response.DistrictResponse{
			ID:         np.District.ID,
			NameZhHant: np.District.NameZhHant,
			NameZhHans: derefString(np.District.NameZhHans),
			NameEn:     derefString(np.District.NameEn),
			Region:     string(np.District.Region),
		}
	}

	// 图片列表
	if len(np.Images) > 0 {
		resp.Images = make([]response.NewDevelopmentImageResponse, len(np.Images))
		for i, img := range np.Images {
			resp.Images[i] = response.NewDevelopmentImageResponse{
				ID:        img.ID,
				URL:       img.ImageURL,
				ImageType: string(img.ImageType),
				Title:     derefString(img.Title),
				SortOrder: img.SortOrder,
			}
		}
	}

	// 户型列表
	if len(np.Layouts) > 0 {
		resp.Layouts = make([]response.NewDevelopmentLayoutResponse, len(np.Layouts))
		for i, layout := range np.Layouts {
			resp.Layouts[i] = s.toLayoutResponse(&layout)
		}
	}

	return resp
}

// toNewDevelopmentListItemResponse 转换为新楼盘列表项响应
func (s *newPropertyService) toNewDevelopmentListItemResponse(np *model.NewProperty) *response.NewDevelopmentListItemResponse {
	resp := &response.NewDevelopmentListItemResponse{
		ID:                 np.ID,
		Name:               np.Name,
		NameEn:             derefString(np.NameEn),
		Address:            np.Address,
		DistrictID:         np.DistrictID,
		Status:             string(np.Status),
		Developer:          np.Developer,
		TotalUnits:         np.TotalUnits,
		UnitsForSale:       derefInt(np.UnitsForSale),
		ExpectedCompletion: np.ExpectedCompletion,
		ViewCount:          np.ViewCount,
		FavoriteCount:      np.FavoriteCount,
		IsFeatured:         np.IsFeatured,
		SalesProgress:      np.GetSalesProgress(),
		CreatedAt:          np.CreatedAt,
	}

	// 地区信息
	if np.District != nil {
		resp.District = &response.DistrictResponse{
			ID:         np.District.ID,
			NameZhHant: np.District.NameZhHant,
			NameZhHans: derefString(np.District.NameZhHans),
			NameEn:     derefString(np.District.NameEn),
			Region:     string(np.District.Region),
		}
	}

	// 封面图 - 从 Images 获取
	if len(np.Images) > 0 {
		resp.CoverImage = np.Images[0].ImageURL
	}

	// 价格范围 - 从 Layouts 计算
	if len(np.Layouts) > 0 {
		var minPrice, maxPrice float64
		for _, layout := range np.Layouts {
			if minPrice == 0 || layout.MinPrice < minPrice {
				minPrice = layout.MinPrice
			}
			layoutMax := derefFloat64(layout.MaxPrice)
			if layoutMax == 0 {
				layoutMax = layout.MinPrice
			}
			if layoutMax > maxPrice {
				maxPrice = layoutMax
			}
		}
		resp.MinPrice = minPrice
		resp.MaxPrice = maxPrice
	}

	return resp
}

// toLayoutResponse 转换为户型响应
func (s *newPropertyService) toLayoutResponse(layout *model.NewPropertyLayout) response.NewDevelopmentLayoutResponse {
	pricePerSqft := float64(0)
	if layout.SaleableArea > 0 {
		pricePerSqft = layout.MinPrice / layout.SaleableArea
	}

	return response.NewDevelopmentLayoutResponse{
		ID:             layout.ID,
		UnitType:       layout.UnitType,
		Bedrooms:       layout.Bedrooms,
		Bathrooms:      derefInt(layout.Bathrooms),
		SaleableArea:   layout.SaleableArea,
		GrossArea:      derefFloat64(layout.GrossArea),
		MinPrice:       layout.MinPrice,
		MaxPrice:       derefFloat64(layout.MaxPrice),
		PricePerSqft:   pricePerSqft,
		AvailableUnits: layout.AvailableUnits,
		FloorplanURL:   derefString(layout.FloorplanURL),
	}
}

// derefString 解引用 string 指针，返回值或空字符串
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// derefInt 解引用 int 指针，返回值或 0
func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// derefFloat64 解引用 float64 指针，返回值或 0
func derefFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
