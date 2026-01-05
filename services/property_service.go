package services

// PropertyService Methods:
// 0. NewPropertyService(propertyRepo databases.PropertyRepository) -> 注入依赖
// 1. ListProperties(ctx context.Context, req *models.ListPropertiesRequest) -> 房产列表
// 2. GetProperty(ctx context.Context, id uint) -> 房产详情
// 3. CreateProperty(ctx context.Context, userID uint, req *models.Property) -> 创建房产
// 4. UpdateProperty(ctx context.Context, userID uint, id uint, req *models.Property) -> 更新房产
// 5. DeleteProperty(ctx context.Context, userID uint, id uint) -> 删除房产
// 6. GetSimilarProperties(ctx context.Context, id uint, limit int) -> 相似房源
// 7. GetFeaturedProperties(ctx context.Context, listingType string, limit int) -> 精选房源
// 8. GetHotProperties(ctx context.Context, listingType string, limit int) -> 热门房源

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

var (
	ErrPropertyNotFound  = errors.New("property not found")
	ErrNotPropertyOwner  = errors.New("you are not the owner of this property")
	ErrPropertyNoInvalid = errors.New("property number already exists")
)

// PropertyServiceInterface 房产服务接口
type PropertyServiceInterface interface {
	ListProperties(ctx context.Context, req *models.ListPropertiesRequest) ([]*models.Property, int64, error)
	GetProperty(ctx context.Context, id uint) (*models.Property, error)
	CreateProperty(ctx context.Context, userID uint, req *models.Property) (*models.Property, error)
	UpdateProperty(ctx context.Context, userID uint, id uint, req *models.Property) (*models.Property, error)
	DeleteProperty(ctx context.Context, userID uint, id uint) error
	GetSimilarProperties(ctx context.Context, id uint, limit int) ([]*models.Property, error)
	GetFeaturedProperties(ctx context.Context, listingType string, limit int) ([]*models.Property, error)
	GetHotProperties(ctx context.Context, listingType string, limit int) ([]*models.Property, error)
}

// PropertyService 房产服务
type PropertyService struct {
	propertyRepo databases.PropertyRepository
}

// 0. NewPropertySvc 注入依赖
func NewPropertySvc(propertyRepo databases.PropertyRepository) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
	}
}

// 1. ListProperties 房产列表
func (s *PropertyService) ListProperties(ctx context.Context, req *models.ListPropertiesRequest) ([]*models.Property, int64, error) {
	properties, total, err := s.propertyRepo.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return s.convertToListItems(properties), total, nil
}

// 2. GetProperty 房产详情
func (s *PropertyService) GetProperty(ctx context.Context, id uint) (*models.Property, error) {
	property, err := s.propertyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, ErrPropertyNotFound
	}

	// 增加浏览量
	_ = s.propertyRepo.IncrementViewCount(ctx, id)

	return s.convertToResponse(property), nil
}

// 3. CreateProperty 创建房产
func (s *PropertyService) CreateProperty(ctx context.Context, userID uint, req *models.Property) (*models.Property, error) {
	// 生成房产编号
	propertyNo := s.generatePropertyNo()

	// 创建房产对象
	property := &models.Property{
		PropertyNo:   propertyNo,
		Title:        req.Title,
		Area:         req.Area,
		Price:        req.Price,
		Address:      req.Address,
		DistrictID:   req.DistrictID,
		Bedrooms:     req.Bedrooms,
		PropertyType: models.PropertyType(req.PropertyType),
		ListingType:  models.ListingType(req.ListingType),
		PublisherID:  userID,
		PublisherType: models.PublisherTypeIndividual,
		Status:       models.PropertyStatusAvailable,
	}

	// 可选字段
	if req.Description != nil && *req.Description != "" {
		property.Description = req.Description
	}
	if req.BuildingName != nil && *req.BuildingName != "" {
		property.BuildingName = req.BuildingName
	}
	if req.Floor != nil && *req.Floor != "" {
		property.Floor = req.Floor
	}
	if req.Orientation != nil && *req.Orientation != "" {
		property.Orientation = req.Orientation
	}
	if req.Bathrooms != nil && *req.Bathrooms > 0 {
		property.Bathrooms = req.Bathrooms
	}
	if req.PrimarySchoolNet != nil && *req.PrimarySchoolNet != "" {
		property.PrimarySchoolNet = req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != nil && *req.SecondarySchoolNet != "" {
		property.SecondarySchoolNet = req.SecondarySchoolNet
	}

	// 设置发布时间
	now := time.Now()
	property.PublishedAt = &now

	if err := s.propertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	return &models.Property{
		ID:         property.ID,
		PropertyNo: property.PropertyNo,
		}, nil
}

// 4. UpdateProperty 更新房产
func (s *PropertyService) UpdateProperty(ctx context.Context, userID uint, id uint, req *models.Property) (*models.Property, error) {
	property, err := s.propertyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, ErrPropertyNotFound
	}

	// 验证所有权
	if property.PublisherID != userID {
		return nil, ErrNotPropertyOwner
	}

	// 更新字段
	if req.Title != "" {
		property.Title = req.Title
	}
	if req.Description != nil && *req.Description != "" {
		property.Description = req.Description
	}
	if req.Area > 0 {
		property.Area = req.Area
	}
	if req.Price > 0 {
		property.Price = req.Price
	}
	if req.Address != "" {
		property.Address = req.Address
	}
	if req.DistrictID > 0 {
		property.DistrictID = req.DistrictID
	}
	if req.BuildingName != nil && *req.BuildingName != "" {
		property.BuildingName = req.BuildingName
	}
	if req.Floor != nil && *req.Floor != "" {
		property.Floor = req.Floor
	}
	if req.Orientation != nil && *req.Orientation != "" {
		property.Orientation = req.Orientation
	}
	if req.Bedrooms >= 0 {
		property.Bedrooms = req.Bedrooms
	}
	if req.Bathrooms != nil && *req.Bathrooms >= 0 {
		property.Bathrooms = req.Bathrooms
	}
	if req.PrimarySchoolNet != nil && *req.PrimarySchoolNet != "" {
		property.PrimarySchoolNet = req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != nil && *req.SecondarySchoolNet != "" {
		property.SecondarySchoolNet = req.SecondarySchoolNet
	}
	if req.PropertyType != "" {
		property.PropertyType = models.PropertyType(req.PropertyType)
	}
	if req.Status != "" {
		property.Status = models.PropertyStatus(req.Status)
	}

	if err := s.propertyRepo.Update(ctx, property); err != nil {
		return nil, err
	}

	return &models.Property{
		ID:      property.ID,
		}, nil
}

// 5. DeleteProperty 删除房产
func (s *PropertyService) DeleteProperty(ctx context.Context, userID uint, id uint) error {
	property, err := s.propertyRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if property == nil {
		return ErrPropertyNotFound
	}

	// 验证所有权
	if property.PublisherID != userID {
		return ErrNotPropertyOwner
	}

	return s.propertyRepo.Delete(ctx, id)
}

// 6. GetSimilarProperties 相似房源
func (s *PropertyService) GetSimilarProperties(ctx context.Context, id uint, limit int) ([]*models.Property, error) {
	property, err := s.propertyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, ErrPropertyNotFound
	}

	properties, err := s.propertyRepo.GetSimilar(ctx, property, limit)
	if err != nil {
		return nil, err
	}

	return s.convertToListItems(properties), nil
}

// 7. GetFeaturedProperties 精选房源
func (s *PropertyService) GetFeaturedProperties(ctx context.Context, listingType string, limit int) ([]*models.Property, error) {
	properties, err := s.propertyRepo.GetFeatured(ctx, listingType, limit)
	if err != nil {
		return nil, err
	}

	return s.convertToListItems(properties), nil
}

// 8. GetHotProperties 热门房源
func (s *PropertyService) GetHotProperties(ctx context.Context, listingType string, limit int) ([]*models.Property, error) {
	properties, err := s.propertyRepo.GetHot(ctx, listingType, limit)
	if err != nil {
		return nil, err
	}

	return s.convertToListItems(properties), nil
}

// generatePropertyNo 生成房产编号
func (s *PropertyService) generatePropertyNo() string {
	return fmt.Sprintf("P%d", time.Now().UnixNano())
}

// convertToListItems 转换为列表项响应（直接返回，预加载了关联数据）
func (s *PropertyService) convertToListItems(properties []*models.Property) []*models.Property {
	return properties
}

// convertToResponse 转换为详情响应（直接返回，预加载了关联数据）
func (s *PropertyService) convertToResponse(p *models.Property) *models.Property {
	return p
}
