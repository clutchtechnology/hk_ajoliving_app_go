package service

// FurnitureService Methods:
// 0. NewFurnitureService(furnitureRepo repository.FurnitureRepository) -> 注入依赖
// 1. ListFurniture(ctx context.Context, req *request.ListFurnitureRequest) -> 家具列表
// 2. GetFurniture(ctx context.Context, id uint) -> 家具详情
// 3. CreateFurniture(ctx context.Context, userID uint, req *request.CreateFurnitureRequest) -> 创建家具
// 4. UpdateFurniture(ctx context.Context, userID uint, id uint, req *request.UpdateFurnitureRequest) -> 更新家具
// 5. DeleteFurniture(ctx context.Context, userID uint, id uint) -> 删除家具
// 6. GetFurnitureCategories(ctx context.Context) -> 获取家具分类
// 7. GetFurnitureImages(ctx context.Context, id uint) -> 获取家具图片
// 8. UpdateFurnitureStatus(ctx context.Context, userID uint, id uint, req *request.UpdateFurnitureStatusRequest) -> 更新家具状态
// 9. GetFeaturedFurniture(ctx context.Context, limit int) -> 获取精选家具

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
	ErrFurnitureNotFound    = errors.New("furniture not found")
	ErrNotFurnitureOwner    = errors.New("you are not the owner of this furniture")
	ErrFurnitureNoInvalid   = errors.New("furniture number already exists")
	ErrFurnitureUnavailable = errors.New("furniture is not available")
	ErrCategoryNotFound     = errors.New("category not found")
)

// FurnitureServiceInterface 家具服务接口
type FurnitureServiceInterface interface {
	ListFurniture(ctx context.Context, req *request.ListFurnitureRequest) ([]*response.FurnitureListItemResponse, int64, error)
	GetFurniture(ctx context.Context, id uint) (*response.FurnitureResponse, error)
	CreateFurniture(ctx context.Context, userID uint, req *request.CreateFurnitureRequest) (*response.CreateFurnitureResponse, error)
	UpdateFurniture(ctx context.Context, userID uint, id uint, req *request.UpdateFurnitureRequest) (*response.UpdateFurnitureResponse, error)
	DeleteFurniture(ctx context.Context, userID uint, id uint) error
	GetFurnitureCategories(ctx context.Context) ([]*response.FurnitureCategoryResponse, error)
	GetFurnitureImages(ctx context.Context, id uint) ([]response.FurnitureImageResponse, error)
	UpdateFurnitureStatus(ctx context.Context, userID uint, id uint, req *request.UpdateFurnitureStatusRequest) (*response.UpdateFurnitureStatusResponse, error)
	GetFeaturedFurniture(ctx context.Context, limit int) ([]*response.FurnitureListItemResponse, error)
}

// FurnitureService 家具服务
type FurnitureService struct {
	furnitureRepo repository.FurnitureRepository
}

// 0. NewFurnitureService 注入依赖
func NewFurnitureService(furnitureRepo repository.FurnitureRepository) *FurnitureService {
	return &FurnitureService{
		furnitureRepo: furnitureRepo,
	}
}

// 1. ListFurniture 家具列表
func (s *FurnitureService) ListFurniture(ctx context.Context, req *request.ListFurnitureRequest) ([]*response.FurnitureListItemResponse, int64, error) {
	furniture, total, err := s.furnitureRepo.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return s.convertToListItems(furniture), total, nil
}

// 2. GetFurniture 家具详情
func (s *FurnitureService) GetFurniture(ctx context.Context, id uint) (*response.FurnitureResponse, error) {
	furniture, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if furniture == nil {
		return nil, ErrFurnitureNotFound
	}

	// 增加浏览量
	_ = s.furnitureRepo.IncrementViewCount(ctx, id)

	return s.convertToResponse(furniture), nil
}

// 3. CreateFurniture 创建家具
func (s *FurnitureService) CreateFurniture(ctx context.Context, userID uint, req *request.CreateFurnitureRequest) (*response.CreateFurnitureResponse, error) {
	// 验证分类是否存在
	category, err := s.furnitureRepo.GetCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	// 生成家具编号
	furnitureNo := s.generateFurnitureNo()

	// 设置默认过期时间（90天后）
	expiresAt := req.ExpiresAt
	if expiresAt == nil {
		expiry := time.Now().AddDate(0, 0, 90)
		expiresAt = &expiry
	}

	// 创建家具对象
	furniture := &model.Furniture{
		FurnitureNo:        furnitureNo,
		Title:              req.Title,
		Description:        req.Description,
		Price:              req.Price,
		CategoryID:         req.CategoryID,
		Brand:              req.Brand,
		Condition:          model.FurnitureCondition(req.Condition),
		PurchaseDate:       req.PurchaseDate,
		DeliveryDistrictID: req.DeliveryDistrictID,
		DeliveryTime:       req.DeliveryTime,
		DeliveryMethod:     model.DeliveryMethod(req.DeliveryMethod),
		Status:             model.FurnitureStatusAvailable,
		PublisherID:        userID,
		PublisherType:      model.PublisherTypeUser,
		ViewCount:          0,
		FavoriteCount:      0,
		PublishedAt:        time.Now(),
		ExpiresAt:          *expiresAt,
	}

	// 保存家具
	if err := s.furnitureRepo.Create(ctx, furniture); err != nil {
		return nil, err
	}

	// 保存图片
	if len(req.ImageURLs) > 0 {
		for i, url := range req.ImageURLs {
			image := &model.FurnitureImage{
				FurnitureID: furniture.ID,
				ImageURL:    url,
				SortOrder:   i + 1,
				IsCover:     i == 0, // 第一张为封面
			}
			// 这里需要实现 CreateImage 方法，暂时跳过
			_ = image
		}
	}

	return &response.CreateFurnitureResponse{
		ID:          furniture.ID,
		FurnitureNo: furniture.FurnitureNo,
		Title:       furniture.Title,
		Price:       furniture.Price,
		Status:      string(furniture.Status),
		PublishedAt: furniture.PublishedAt,
		ExpiresAt:   furniture.ExpiresAt,
		Message:     "家具发布成功",
	}, nil
}

// 4. UpdateFurniture 更新家具
func (s *FurnitureService) UpdateFurniture(ctx context.Context, userID uint, id uint, req *request.UpdateFurnitureRequest) (*response.UpdateFurnitureResponse, error) {
	// 获取家具
	furniture, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if furniture == nil {
		return nil, ErrFurnitureNotFound
	}

	// 验证权限
	if furniture.PublisherID != userID {
		return nil, ErrNotFurnitureOwner
	}

	// 更新字段
	if req.Title != nil {
		furniture.Title = *req.Title
	}
	if req.Description != nil {
		furniture.Description = req.Description
	}
	if req.Price != nil {
		furniture.Price = *req.Price
	}
	if req.CategoryID != nil {
		// 验证分类是否存在
		category, err := s.furnitureRepo.GetCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, ErrCategoryNotFound
		}
		furniture.CategoryID = *req.CategoryID
	}
	if req.Brand != nil {
		furniture.Brand = req.Brand
	}
	if req.Condition != nil {
		furniture.Condition = model.FurnitureCondition(*req.Condition)
	}
	if req.PurchaseDate != nil {
		furniture.PurchaseDate = req.PurchaseDate
	}
	if req.DeliveryDistrictID != nil {
		furniture.DeliveryDistrictID = *req.DeliveryDistrictID
	}
	if req.DeliveryTime != nil {
		furniture.DeliveryTime = req.DeliveryTime
	}
	if req.DeliveryMethod != nil {
		furniture.DeliveryMethod = model.DeliveryMethod(*req.DeliveryMethod)
	}
	if req.ExpiresAt != nil {
		furniture.ExpiresAt = *req.ExpiresAt
	}

	// 保存更新
	if err := s.furnitureRepo.Update(ctx, furniture); err != nil {
		return nil, err
	}

	return &response.UpdateFurnitureResponse{
		ID:          furniture.ID,
		FurnitureNo: furniture.FurnitureNo,
		Title:       furniture.Title,
		UpdatedAt:   furniture.UpdatedAt,
		Message:     "家具更新成功",
	}, nil
}

// 5. DeleteFurniture 删除家具
func (s *FurnitureService) DeleteFurniture(ctx context.Context, userID uint, id uint) error {
	// 获取家具
	furniture, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if furniture == nil {
		return ErrFurnitureNotFound
	}

	// 验证权限
	if furniture.PublisherID != userID {
		return ErrNotFurnitureOwner
	}

	// 删除家具
	return s.furnitureRepo.Delete(ctx, id)
}

// 6. GetFurnitureCategories 获取家具分类
func (s *FurnitureService) GetFurnitureCategories(ctx context.Context) ([]*response.FurnitureCategoryResponse, error) {
	categories, err := s.furnitureRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	return s.convertToCategoryResponses(categories), nil
}

// 7. GetFurnitureImages 获取家具图片
func (s *FurnitureService) GetFurnitureImages(ctx context.Context, id uint) ([]response.FurnitureImageResponse, error) {
	// 验证家具是否存在
	furniture, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if furniture == nil {
		return nil, ErrFurnitureNotFound
	}

	images, err := s.furnitureRepo.GetImagesByFurnitureID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToImageResponses(images), nil
}

// 8. UpdateFurnitureStatus 更新家具状态
func (s *FurnitureService) UpdateFurnitureStatus(ctx context.Context, userID uint, id uint, req *request.UpdateFurnitureStatusRequest) (*response.UpdateFurnitureStatusResponse, error) {
	// 获取家具
	furniture, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if furniture == nil {
		return nil, ErrFurnitureNotFound
	}

	// 验证权限
	if furniture.PublisherID != userID {
		return nil, ErrNotFurnitureOwner
	}

	// 更新状态
	status := model.FurnitureStatus(req.Status)
	if err := s.furnitureRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return &response.UpdateFurnitureStatusResponse{
		ID:        id,
		Status:    req.Status,
		UpdatedAt: time.Now(),
		Message:   "状态更新成功",
	}, nil
}

// 9. GetFeaturedFurniture 获取精选家具
func (s *FurnitureService) GetFeaturedFurniture(ctx context.Context, limit int) ([]*response.FurnitureListItemResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	furniture, err := s.furnitureRepo.GetFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}

	return s.convertToListItems(furniture), nil
}

// 辅助方法

// generateFurnitureNo 生成家具编号
func (s *FurnitureService) generateFurnitureNo() string {
	return fmt.Sprintf("FUR%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}

// convertToListItems 转换为列表项响应
func (s *FurnitureService) convertToListItems(furniture []*model.Furniture) []*response.FurnitureListItemResponse {
	items := make([]*response.FurnitureListItemResponse, 0, len(furniture))
	for _, f := range furniture {
		items = append(items, s.convertToListItem(f))
	}
	return items
}

// convertToListItem 转换为列表项响应
func (s *FurnitureService) convertToListItem(f *model.Furniture) *response.FurnitureListItemResponse {
	item := &response.FurnitureListItemResponse{
		ID:                 f.ID,
		FurnitureNo:        f.FurnitureNo,
		Title:              f.Title,
		Price:              f.Price,
		CategoryID:         f.CategoryID,
		Brand:              f.Brand,
		Condition:          string(f.Condition),
		DeliveryDistrictID: f.DeliveryDistrictID,
		DeliveryMethod:     string(f.DeliveryMethod),
		Status:             string(f.Status),
		ViewCount:          f.ViewCount,
		FavoriteCount:      f.FavoriteCount,
		PublishedAt:        f.PublishedAt,
		ExpiresAt:          f.ExpiresAt,
		DaysUntilExpiry:    f.GetDaysUntilExpiry(),
	}

	// 分类名称
	if f.Category != nil {
		item.CategoryName = f.Category.NameZhHant
	}

	// 交收地区
	if f.DeliveryDistrict != nil {
		item.DeliveryDistrict = f.DeliveryDistrict.NameZhHant
	}

	// 封面图片
	if len(f.Images) > 0 {
		item.CoverImage = &f.Images[0].ImageURL
	}

	return item
}

// convertToResponse 转换为详情响应
func (s *FurnitureService) convertToResponse(f *model.Furniture) *response.FurnitureResponse {
	resp := &response.FurnitureResponse{
		ID:                 f.ID,
		FurnitureNo:        f.FurnitureNo,
		Title:              f.Title,
		Description:        f.Description,
		Price:              f.Price,
		CategoryID:         f.CategoryID,
		Brand:              f.Brand,
		Condition:          string(f.Condition),
		PurchaseDate:       f.PurchaseDate,
		Age:                f.GetAge(),
		DeliveryDistrictID: f.DeliveryDistrictID,
		DeliveryTime:       f.DeliveryTime,
		DeliveryMethod:     string(f.DeliveryMethod),
		SupportsDelivery:   f.SupportsDelivery(),
		SupportsSelfPickup: f.SupportsSelfPickup(),
		Status:             string(f.Status),
		PublisherID:        f.PublisherID,
		ViewCount:          f.ViewCount,
		FavoriteCount:      f.FavoriteCount,
		PublishedAt:        f.PublishedAt,
		UpdatedAt:          f.UpdatedAt,
		ExpiresAt:          f.ExpiresAt,
		DaysUntilExpiry:    f.GetDaysUntilExpiry(),
		IsAvailable:        f.IsAvailable(),
		IsExpired:          f.IsExpired(),
		CreatedAt:          f.CreatedAt,
	}

	// 分类信息
	if f.Category != nil {
		resp.Category = &response.FurnitureCategoryResponse{
			ID:         f.Category.ID,
			ParentID:   f.Category.ParentID,
			NameZhHant: f.Category.NameZhHant,
			NameZhHans: f.Category.NameZhHans,
			NameEn:     f.Category.NameEn,
			Icon:       f.Category.Icon,
			SortOrder:  f.Category.SortOrder,
			IsActive:   f.Category.IsActive,
			IsTopLevel: f.Category.IsTopLevel(),
		}
	}

	// 交收地区信息
	if f.DeliveryDistrict != nil {
		resp.DeliveryDistrict = &response.DistrictBasicResponse{
			ID:         f.DeliveryDistrict.ID,
			NameZhHant: f.DeliveryDistrict.NameZhHant,
			NameZhHans: f.DeliveryDistrict.NameZhHans,
			NameEn:     f.DeliveryDistrict.NameEn,
		}
	}

	// 发布者信息
	if f.Publisher != nil {
		resp.Publisher = &response.PublisherBasicResponse{
			ID:            f.Publisher.ID,
			Name:          f.Publisher.Name,
			Avatar:        f.Publisher.Avatar,
			PublisherType: string(f.PublisherType),
		}
	}

	// 图片列表
	resp.Images = s.convertToImageResponses(f.Images)

	return resp
}

// convertToCategoryResponses 转换为分类响应列表
func (s *FurnitureService) convertToCategoryResponses(categories []*model.FurnitureCategory) []*response.FurnitureCategoryResponse {
	responses := make([]*response.FurnitureCategoryResponse, 0, len(categories))
	for _, c := range categories {
		// 只返回顶级分类（带子分类）
		if c.IsTopLevel() {
			responses = append(responses, s.convertToCategoryResponse(c))
		}
	}
	return responses
}

// convertToCategoryResponse 转换为分类响应
func (s *FurnitureService) convertToCategoryResponse(c *model.FurnitureCategory) *response.FurnitureCategoryResponse {
	resp := &response.FurnitureCategoryResponse{
		ID:         c.ID,
		ParentID:   c.ParentID,
		NameZhHant: c.NameZhHant,
		NameZhHans: c.NameZhHans,
		NameEn:     c.NameEn,
		Icon:       c.Icon,
		SortOrder:  c.SortOrder,
		IsActive:   c.IsActive,
		IsTopLevel: c.IsTopLevel(),
	}

	// 子分类
	if len(c.Subcategories) > 0 {
		subcategories := make([]response.FurnitureCategoryResponse, 0, len(c.Subcategories))
		for _, sub := range c.Subcategories {
			subcategories = append(subcategories, response.FurnitureCategoryResponse{
				ID:         sub.ID,
				ParentID:   sub.ParentID,
				NameZhHant: sub.NameZhHant,
				NameZhHans: sub.NameZhHans,
				NameEn:     sub.NameEn,
				Icon:       sub.Icon,
				SortOrder:  sub.SortOrder,
				IsActive:   sub.IsActive,
				IsTopLevel: sub.IsTopLevel(),
			})
		}
		resp.Subcategories = subcategories
	}

	return resp
}

// convertToImageResponses 转换为图片响应列表
func (s *FurnitureService) convertToImageResponses(images []model.FurnitureImage) []response.FurnitureImageResponse {
	responses := make([]response.FurnitureImageResponse, 0, len(images))
	for _, img := range images {
		responses = append(responses, response.FurnitureImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			SortOrder: img.SortOrder,
			IsCover:   img.IsCover,
		})
	}
	return responses
}
