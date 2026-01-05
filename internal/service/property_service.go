package service

// PropertyService Methods:
// 0. NewPropertyService(propertyRepo repository.PropertyRepository) -> 注入依赖
// 1. ListProperties(ctx context.Context, req *request.ListPropertiesRequest) -> 房产列表
// 2. GetProperty(ctx context.Context, id uint) -> 房产详情
// 3. CreateProperty(ctx context.Context, userID uint, req *request.CreatePropertyRequest) -> 创建房产
// 4. UpdateProperty(ctx context.Context, userID uint, id uint, req *request.UpdatePropertyRequest) -> 更新房产
// 5. DeleteProperty(ctx context.Context, userID uint, id uint) -> 删除房产
// 6. GetSimilarProperties(ctx context.Context, id uint, limit int) -> 相似房源
// 7. GetFeaturedProperties(ctx context.Context, listingType string, limit int) -> 精选房源
// 8. GetHotProperties(ctx context.Context, listingType string, limit int) -> 热门房源

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

var (
	ErrPropertyNotFound  = errors.New("property not found")
	ErrNotPropertyOwner  = errors.New("you are not the owner of this property")
	ErrPropertyNoInvalid = errors.New("property number already exists")
)

// PropertyServiceInterface 房产服务接口
type PropertyServiceInterface interface {
	ListProperties(ctx context.Context, req *request.ListPropertiesRequest) ([]*response.PropertyListItemResponse, int64, error)
	GetProperty(ctx context.Context, id uint) (*response.PropertyResponse, error)
	CreateProperty(ctx context.Context, userID uint, req *request.CreatePropertyRequest) (*response.CreatePropertyResponse, error)
	UpdateProperty(ctx context.Context, userID uint, id uint, req *request.UpdatePropertyRequest) (*response.UpdatePropertyResponse, error)
	DeleteProperty(ctx context.Context, userID uint, id uint) error
	GetSimilarProperties(ctx context.Context, id uint, limit int) ([]*response.PropertyListItemResponse, error)
	GetFeaturedProperties(ctx context.Context, listingType string, limit int) ([]*response.PropertyListItemResponse, error)
	GetHotProperties(ctx context.Context, listingType string, limit int) ([]*response.PropertyListItemResponse, error)
}

// PropertyService 房产服务
type PropertyService struct {
	propertyRepo repository.PropertyRepository
}

// 0. NewPropertySvc 注入依赖
func NewPropertySvc(propertyRepo repository.PropertyRepository) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
	}
}

// 1. ListProperties 房产列表
func (s *PropertyService) ListProperties(ctx context.Context, req *request.ListPropertiesRequest) ([]*response.PropertyListItemResponse, int64, error) {
	properties, total, err := s.propertyRepo.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return s.convertToListItems(properties), total, nil
}

// 2. GetProperty 房产详情
func (s *PropertyService) GetProperty(ctx context.Context, id uint) (*response.PropertyResponse, error) {
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
func (s *PropertyService) CreateProperty(ctx context.Context, userID uint, req *request.CreatePropertyRequest) (*response.CreatePropertyResponse, error) {
	// 生成房产编号
	propertyNo := s.generatePropertyNo()

	// 创建房产对象
	property := &model.Property{
		PropertyNo:   propertyNo,
		Title:        req.Title,
		Area:         req.Area,
		Price:        req.Price,
		Address:      req.Address,
		DistrictID:   req.DistrictID,
		Bedrooms:     req.Bedrooms,
		PropertyType: model.PropertyType(req.PropertyType),
		ListingType:  model.ListingType(req.ListingType),
		PublisherID:  userID,
		PublisherType: model.PublisherTypeIndividual,
		Status:       model.PropertyStatusAvailable,
	}

	// 可选字段
	if req.Description != "" {
		property.Description = &req.Description
	}
	if req.BuildingName != "" {
		property.BuildingName = &req.BuildingName
	}
	if req.Floor != "" {
		property.Floor = &req.Floor
	}
	if req.Orientation != "" {
		property.Orientation = &req.Orientation
	}
	if req.Bathrooms > 0 {
		property.Bathrooms = &req.Bathrooms
	}
	if req.PrimarySchoolNet != "" {
		property.PrimarySchoolNet = &req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != "" {
		property.SecondarySchoolNet = &req.SecondarySchoolNet
	}

	// 设置发布时间
	now := time.Now()
	property.PublishedAt = &now

	if err := s.propertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	return &response.CreatePropertyResponse{
		ID:         property.ID,
		PropertyNo: property.PropertyNo,
		Message:    "Property created successfully",
	}, nil
}

// 4. UpdateProperty 更新房产
func (s *PropertyService) UpdateProperty(ctx context.Context, userID uint, id uint, req *request.UpdatePropertyRequest) (*response.UpdatePropertyResponse, error) {
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
	if req.Description != "" {
		property.Description = &req.Description
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
	if req.BuildingName != "" {
		property.BuildingName = &req.BuildingName
	}
	if req.Floor != "" {
		property.Floor = &req.Floor
	}
	if req.Orientation != "" {
		property.Orientation = &req.Orientation
	}
	if req.Bedrooms >= 0 {
		property.Bedrooms = req.Bedrooms
	}
	if req.Bathrooms >= 0 {
		property.Bathrooms = &req.Bathrooms
	}
	if req.PrimarySchoolNet != "" {
		property.PrimarySchoolNet = &req.PrimarySchoolNet
	}
	if req.SecondarySchoolNet != "" {
		property.SecondarySchoolNet = &req.SecondarySchoolNet
	}
	if req.PropertyType != "" {
		property.PropertyType = model.PropertyType(req.PropertyType)
	}
	if req.Status != "" {
		property.Status = model.PropertyStatus(req.Status)
	}

	if err := s.propertyRepo.Update(ctx, property); err != nil {
		return nil, err
	}

	return &response.UpdatePropertyResponse{
		ID:      property.ID,
		Message: "Property updated successfully",
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
func (s *PropertyService) GetSimilarProperties(ctx context.Context, id uint, limit int) ([]*response.PropertyListItemResponse, error) {
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
func (s *PropertyService) GetFeaturedProperties(ctx context.Context, listingType string, limit int) ([]*response.PropertyListItemResponse, error) {
	properties, err := s.propertyRepo.GetFeatured(ctx, listingType, limit)
	if err != nil {
		return nil, err
	}

	return s.convertToListItems(properties), nil
}

// 8. GetHotProperties 热门房源
func (s *PropertyService) GetHotProperties(ctx context.Context, listingType string, limit int) ([]*response.PropertyListItemResponse, error) {
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

// convertToListItems 转换为列表项响应
func (s *PropertyService) convertToListItems(properties []*model.Property) []*response.PropertyListItemResponse {
	items := make([]*response.PropertyListItemResponse, 0, len(properties))
	for _, p := range properties {
		item := &response.PropertyListItemResponse{
			ID:           p.ID,
			PropertyNo:   p.PropertyNo,
			ListingType:  string(p.ListingType),
			Title:        p.Title,
			Area:         p.Area,
			Price:        p.Price,
			Address:      p.Address,
			DistrictID:   p.DistrictID,
			Bedrooms:     p.Bedrooms,
			PropertyType: string(p.PropertyType),
			Status:       string(p.Status),
			ViewCount:    p.ViewCount,
			FavoriteCount: p.FavoriteCount,
			CreatedAt:    p.CreatedAt,
		}

		if p.BuildingName != nil {
			item.BuildingName = *p.BuildingName
		}
		if p.Bathrooms != nil {
			item.Bathrooms = *p.Bathrooms
		}

		// 获取封面图
		for _, img := range p.Images {
			if img.IsCover {
				item.CoverImage = img.URL
				break
			}
		}
		if item.CoverImage == "" && len(p.Images) > 0 {
			item.CoverImage = p.Images[0].URL
		}

		// 地区信息
		if p.District != nil {
			item.District = &response.DistrictResponse{
				ID:         p.District.ID,
				NameZhHant: p.District.NameZhHant,
				Region:     string(p.District.Region),
			}
			if p.District.NameZhHans != nil {
				item.District.NameZhHans = *p.District.NameZhHans
			}
			if p.District.NameEn != nil {
				item.District.NameEn = *p.District.NameEn
			}
		}

		items = append(items, item)
	}
	return items
}

// convertToResponse 转换为详情响应
func (s *PropertyService) convertToResponse(p *model.Property) *response.PropertyResponse {
	resp := &response.PropertyResponse{
		ID:            p.ID,
		PropertyNo:    p.PropertyNo,
		ListingType:   string(p.ListingType),
		Title:         p.Title,
		Area:          p.Area,
		Price:         p.Price,
		Address:       p.Address,
		DistrictID:    p.DistrictID,
		Bedrooms:      p.Bedrooms,
		PropertyType:  string(p.PropertyType),
		Status:        string(p.Status),
		PublisherID:   p.PublisherID,
		PublisherType: string(p.PublisherType),
		ViewCount:     p.ViewCount,
		FavoriteCount: p.FavoriteCount,
		PublishedAt:   p.PublishedAt,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	// 可选字段
	if p.EstateNo != nil {
		resp.EstateNo = *p.EstateNo
	}
	if p.Description != nil {
		resp.Description = *p.Description
	}
	if p.BuildingName != nil {
		resp.BuildingName = *p.BuildingName
	}
	if p.Floor != nil {
		resp.Floor = *p.Floor
	}
	if p.Orientation != nil {
		resp.Orientation = *p.Orientation
	}
	if p.Bathrooms != nil {
		resp.Bathrooms = *p.Bathrooms
	}
	if p.PrimarySchoolNet != nil {
		resp.PrimarySchoolNet = *p.PrimarySchoolNet
	}
	if p.SecondarySchoolNet != nil {
		resp.SecondarySchoolNet = *p.SecondarySchoolNet
	}
	if p.AgentID != nil {
		resp.AgentID = *p.AgentID
	}

	// 地区信息
	if p.District != nil {
		resp.District = &response.DistrictResponse{
			ID:         p.District.ID,
			NameZhHant: p.District.NameZhHant,
			Region:     string(p.District.Region),
		}
		if p.District.NameZhHans != nil {
			resp.District.NameZhHans = *p.District.NameZhHans
		}
		if p.District.NameEn != nil {
			resp.District.NameEn = *p.District.NameEn
		}
	}

	// 代理信息
	if p.Agent != nil {
		avatar := ""
		if p.Agent.ProfilePhoto != nil {
			avatar = *p.Agent.ProfilePhoto
		}
		resp.Agent = &response.AgentBriefResponse{
			ID:        p.Agent.ID,
			Name:      p.Agent.AgentName,
			Phone:     p.Agent.Phone,
			Avatar:    avatar,
			LicenseNo: p.Agent.LicenseNo,
		}
	}

	// 图片
	if len(p.Images) > 0 {
		resp.Images = make([]response.PropertyImageResponse, 0, len(p.Images))
		for _, img := range p.Images {
			imgResp := response.PropertyImageResponse{
				ID:        img.ID,
				URL:       img.URL,
				SortOrder: img.SortOrder,
				IsCover:   img.IsCover,
			}
			if img.Caption != nil {
				imgResp.Caption = *img.Caption
			}
			resp.Images = append(resp.Images, imgResp)
		}
	}

	// 设施
	if len(p.Facilities) > 0 {
		resp.Facilities = make([]response.FacilityResponse, 0, len(p.Facilities))
		for _, f := range p.Facilities {
			facResp := response.FacilityResponse{
				ID:         f.ID,
				NameZhHant: f.NameZhHant,
				Category:   string(f.Category),
				NameZhHans: f.NameZhHans,
				NameEn:     f.NameEn,
				Icon:       f.Icon,
			}
			resp.Facilities = append(resp.Facilities, facResp)
		}
	}

	return resp
}
