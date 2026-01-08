package services

import (
	"context"
	"errors"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// PropertyService 房产服务
type PropertyService struct {
	propertyRepo *databases.PropertyRepo
}

// NewPropertyService 创建房产服务
func NewPropertyService(propertyRepo *databases.PropertyRepo) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
	}
}

// ListProperties 获取房产列表
func (s *PropertyService) ListProperties(ctx context.Context, filter *models.ListPropertiesRequest) (*models.PaginatedPropertiesResponse, error) {
	properties, total, err := s.propertyRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	data := make([]models.PropertyResponse, len(properties))
	for i, p := range properties {
		data[i] = *p.ToPropertyResponse()
	}

	return &models.PaginatedPropertiesResponse{
		Data:       data,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: databases.CalculateTotalPages(total, filter.PageSize),
	}, nil
}

// GetProperty 获取房产详情
func (s *PropertyService) GetProperty(ctx context.Context, id uint) (*models.PropertyDetailResponse, error) {
	property, err := s.propertyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 增加浏览次数（异步执行，不影响主流程）
	go s.propertyRepo.IncrementViewCount(context.Background(), id)

	return property.ToPropertyDetailResponse(), nil
}

// CreateProperty 创建房产
func (s *PropertyService) CreateProperty(ctx context.Context, userID uint, userType string, req *models.CreatePropertyRequest) (*models.PropertyDetailResponse, error) {
	// 生成房产编号
	propertyNo, err := s.propertyRepo.GeneratePropertyNo(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiredAt := now.AddDate(0, 3, 0) // 默认3个月后过期

	property := &models.Property{
		PropertyNo:      propertyNo,
		EstateNo:        req.EstateNo,
		ListingType:     req.ListingType,
		Title:           req.Title,
		Description:     req.Description,
		Area:            req.Area,
		Price:           req.Price,
		Address:         req.Address,
		DistrictID:      req.DistrictID,
		BuildingName:    req.BuildingName,
		Floor:           req.Floor,
		Orientation:     req.Orientation,
		Bedrooms:        req.Bedrooms,
		Bathrooms:       req.Bathrooms,
		PrimarySchool:   req.PrimarySchool,
		SecondarySchool: req.SecondarySchool,
		PropertyType:    req.PropertyType,
		Status:          "available",
		PublisherID:     userID,
		PublisherType:   userType,
		AgentID:         req.AgentID,
		ViewCount:       0,
		FavoriteCount:   0,
		PublishedAt:     &now,
		ExpiredAt:       &expiredAt,
	}

	// 创建房产
	if err := s.propertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	// 创建图片
	if len(req.ImageURLs) > 0 {
		images := make([]models.PropertyImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			imageType := "interior"
			if i == 0 {
				imageType = "cover" // 第一张作为封面
			}
			images[i] = models.PropertyImage{
				PropertyID: property.ID,
				ImageURL:   url,
				ImageType:  imageType,
				SortOrder:  i,
			}
		}
		if err := s.propertyRepo.CreateImages(ctx, images); err != nil {
			return nil, err
		}
	}

	// 重新查询完整信息
	return s.GetProperty(ctx, property.ID)
}

// UpdateProperty 更新房产
func (s *PropertyService) UpdateProperty(ctx context.Context, id uint, userID uint, req *models.UpdatePropertyRequest) (*models.PropertyDetailResponse, error) {
	// 查找房产
	property, err := s.propertyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 权限检查：只能更新自己发布的房产
	if property.PublisherID != userID {
		return nil, errors.New("permission denied")
	}

	// 更新字段
	if req.Title != nil {
		property.Title = *req.Title
	}
	if req.Description != nil {
		property.Description = *req.Description
	}
	if req.Price != nil {
		property.Price = *req.Price
	}
	if req.Area != nil {
		property.Area = *req.Area
	}
	if req.Floor != nil {
		property.Floor = *req.Floor
	}
	if req.Orientation != nil {
		property.Orientation = *req.Orientation
	}
	if req.Bathrooms != nil {
		property.Bathrooms = *req.Bathrooms
	}
	if req.Status != nil {
		property.Status = *req.Status
	}
	if req.AgentID != nil {
		property.AgentID = req.AgentID
	}

	// 保存更新
	if err := s.propertyRepo.Update(ctx, property); err != nil {
		return nil, err
	}

	// 重新查询完整信息
	return s.GetProperty(ctx, id)
}

// DeleteProperty 删除房产
func (s *PropertyService) DeleteProperty(ctx context.Context, id uint, userID uint) error {
	// 查找房产
	property, err := s.propertyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 权限检查：只能删除自己发布的房产
	if property.PublisherID != userID {
		return errors.New("permission denied")
	}

	return s.propertyRepo.Delete(ctx, id)
}

// GetSimilarProperties 获取相似房源
func (s *PropertyService) GetSimilarProperties(ctx context.Context, id uint, limit int) ([]models.PropertyResponse, error) {
	// 获取原房产信息
	property, err := s.propertyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 查找相似房源
	similar, err := s.propertyRepo.FindSimilar(ctx, property, limit)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	data := make([]models.PropertyResponse, len(similar))
	for i, p := range similar {
		data[i] = *p.ToPropertyResponse()
	}

	return data, nil
}

// GetFeaturedProperties 获取精选房源
func (s *PropertyService) GetFeaturedProperties(ctx context.Context, limit int) ([]models.PropertyResponse, error) {
	featured, err := s.propertyRepo.FindFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	data := make([]models.PropertyResponse, len(featured))
	for i, p := range featured {
		data[i] = *p.ToPropertyResponse()
	}

	return data, nil
}

// GetHotProperties 获取热门房源
func (s *PropertyService) GetHotProperties(ctx context.Context, limit int) ([]models.PropertyResponse, error) {
	hot, err := s.propertyRepo.FindHot(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	data := make([]models.PropertyResponse, len(hot))
	for i, p := range hot {
		data[i] = *p.ToPropertyResponse()
	}

	return data, nil
}
