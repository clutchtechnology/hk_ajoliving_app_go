package services

import (
	"context"
	"errors"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// FurnitureService 家具服务
type FurnitureService struct {
	repo *databases.FurnitureRepo
}

// NewFurnitureService 创建家具服务
func NewFurnitureService(repo *databases.FurnitureRepo) *FurnitureService {
	return &FurnitureService{repo: repo}
}

// ListFurniture 获取家具列表
func (s *FurnitureService) ListFurniture(ctx context.Context, filter *models.ListFurnitureRequest) (*models.PaginatedFurnitureResponse, error) {
	furniture, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	data := make([]models.FurnitureResponse, len(furniture))
	for i, f := range furniture {
		data[i] = s.toFurnitureResponse(&f)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedFurnitureResponse{
		Data:       data,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetFurniture 获取家具详情
func (s *FurnitureService) GetFurniture(ctx context.Context, id uint) (*models.FurnitureResponse, error) {
	furniture, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 增加浏览次数
	_ = s.repo.IncrementViewCount(ctx, id)

	response := s.toFurnitureResponse(furniture)
	return &response, nil
}

// GetFeaturedFurniture 获取精选家具
func (s *FurnitureService) GetFeaturedFurniture(ctx context.Context, limit int) ([]models.FurnitureResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	furniture, err := s.repo.FindFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}

	response := make([]models.FurnitureResponse, len(furniture))
	for i, f := range furniture {
		response[i] = s.toFurnitureResponse(&f)
	}

	return response, nil
}

// GetFurnitureImages 获取家具图片
func (s *FurnitureService) GetFurnitureImages(ctx context.Context, id uint) ([]models.FurnitureImage, error) {
	// 先查询家具是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	images, err := s.repo.FindImagesByFurnitureID(ctx, id)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// CreateFurniture 创建家具
func (s *FurnitureService) CreateFurniture(ctx context.Context, userID uint, userType string, req *models.CreateFurnitureRequest) (*models.FurnitureResponse, error) {
	// 生成家具编号
	furnitureNo, err := s.repo.GenerateFurnitureNo(ctx)
	if err != nil {
		return nil, err
	}

	// 设置过期时间（90天后）
	now := time.Now()
	expiresAt := now.AddDate(0, 0, 90)

	furniture := &models.Furniture{
		FurnitureNo:        furnitureNo,
		Title:              req.Title,
		Description:        req.Description,
		Price:              req.Price,
		CategoryID:         req.CategoryID,
		Brand:              req.Brand,
		Condition:          req.Condition,
		PurchaseDate:       req.PurchaseDate,
		DeliveryDistrictID: req.DeliveryDistrictID,
		DeliveryTime:       req.DeliveryTime,
		DeliveryMethod:     req.DeliveryMethod,
		Status:             "available",
		PublisherID:        userID,
		PublisherType:      userType,
		PublishedAt:        now,
		ExpiresAt:          expiresAt,
	}

	if err := s.repo.Create(ctx, furniture); err != nil {
		return nil, err
	}

	// 创建图片
	if len(req.ImageURLs) > 0 {
		images := make([]models.FurnitureImage, len(req.ImageURLs))
		for i, url := range req.ImageURLs {
			images[i] = models.FurnitureImage{
				FurnitureID: furniture.ID,
				ImageURL:    url,
				IsCover:     i == 0, // 第一张为封面
				SortOrder:   i,
				CreatedAt:   now,
			}
		}
		if err := s.repo.CreateImages(ctx, images); err != nil {
			return nil, err
		}
	}

	// 重新查询以获取完整数据
	created, err := s.repo.FindByID(ctx, furniture.ID)
	if err != nil {
		return nil, err
	}

	response := s.toFurnitureResponse(created)
	return &response, nil
}

// UpdateFurniture 更新家具
func (s *FurnitureService) UpdateFurniture(ctx context.Context, id uint, userID uint, req *models.UpdateFurnitureRequest) (*models.FurnitureResponse, error) {
	// 查询家具是否存在
	furniture, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 验证权限：只有发布者可以更新
	if furniture.PublisherID != userID {
		return nil, tools.ErrForbidden
	}

	// 更新字段
	if req.Title != nil {
		furniture.Title = *req.Title
	}
	if req.Description != nil {
		furniture.Description = *req.Description
	}
	if req.Price != nil {
		furniture.Price = *req.Price
	}
	if req.CategoryID != nil {
		furniture.CategoryID = *req.CategoryID
	}
	if req.Brand != nil {
		furniture.Brand = *req.Brand
	}
	if req.Condition != nil {
		furniture.Condition = *req.Condition
	}
	if req.PurchaseDate != nil {
		furniture.PurchaseDate = req.PurchaseDate
	}
	if req.DeliveryDistrictID != nil {
		furniture.DeliveryDistrictID = *req.DeliveryDistrictID
	}
	if req.DeliveryTime != nil {
		furniture.DeliveryTime = *req.DeliveryTime
	}
	if req.DeliveryMethod != nil {
		furniture.DeliveryMethod = *req.DeliveryMethod
	}

	if err := s.repo.Update(ctx, furniture); err != nil {
		return nil, err
	}

	// 更新图片
	if req.ImageURLs != nil {
		if err := s.repo.UpdateImages(ctx, id, req.ImageURLs); err != nil {
			return nil, err
		}
	}

	// 重新查询以获取完整数据
	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.toFurnitureResponse(updated)
	return &response, nil
}

// UpdateFurnitureStatus 更新家具状态
func (s *FurnitureService) UpdateFurnitureStatus(ctx context.Context, id uint, userID uint, status string) error {
	// 查询家具是否存在
	furniture, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tools.ErrNotFound
		}
		return err
	}

	// 验证权限：只有发布者可以更新状态
	if furniture.PublisherID != userID {
		return tools.ErrForbidden
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// DeleteFurniture 删除家具
func (s *FurnitureService) DeleteFurniture(ctx context.Context, id uint, userID uint) error {
	// 查询家具是否存在
	furniture, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tools.ErrNotFound
		}
		return err
	}

	// 验证权限：只有发布者可以删除
	if furniture.PublisherID != userID {
		return tools.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

// GetFurnitureCategories 获取家具分类列表
func (s *FurnitureService) GetFurnitureCategories(ctx context.Context) ([]models.FurnitureCategoryResponse, error) {
	categories, err := s.repo.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]models.FurnitureCategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = s.toCategoryResponse(ctx, &category)
	}

	return response, nil
}

// toFurnitureResponse 转换为响应格式
func (s *FurnitureService) toFurnitureResponse(furniture *models.Furniture) models.FurnitureResponse {
	return models.FurnitureResponse{
		ID:                 furniture.ID,
		FurnitureNo:        furniture.FurnitureNo,
		Title:              furniture.Title,
		Description:        furniture.Description,
		Price:              furniture.Price,
		Category:           furniture.Category,
		Brand:              furniture.Brand,
		Condition:          furniture.Condition,
		PurchaseDate:       furniture.PurchaseDate,
		DeliveryDistrictID: furniture.DeliveryDistrictID,
		DeliveryDistrict:   furniture.DeliveryDistrict,
		DeliveryTime:       furniture.DeliveryTime,
		DeliveryMethod:     furniture.DeliveryMethod,
		Status:             furniture.Status,
		PublisherID:        furniture.PublisherID,
		PublisherType:      furniture.PublisherType,
		ViewCount:          furniture.ViewCount,
		FavoriteCount:      furniture.FavoriteCount,
		Images:             furniture.Images,
		PublishedAt:        furniture.PublishedAt,
		UpdatedAt:          furniture.UpdatedAt,
		ExpiresAt:          furniture.ExpiresAt,
	}
}

// toCategoryResponse 转换分类为响应格式
func (s *FurnitureService) toCategoryResponse(ctx context.Context, category *models.FurnitureCategory) models.FurnitureCategoryResponse {
	response := models.FurnitureCategoryResponse{
		ID:         category.ID,
		ParentID:   category.ParentID,
		NameZhHant: category.NameZhHant,
		NameZhHans: category.NameZhHans,
		NameEn:     category.NameEn,
		Icon:       category.Icon,
		SortOrder:  category.SortOrder,
		IsActive:   category.IsActive,
	}

	// 获取该分类下的家具数量
	count, _ := s.repo.GetFurnitureCountByCategory(ctx, category.ID)
	response.FurnitureCount = int(count)

	// 处理子分类
	if len(category.SubCategories) > 0 {
		response.SubCategories = make([]models.FurnitureCategoryResponse, len(category.SubCategories))
		for i, subCategory := range category.SubCategories {
			response.SubCategories[i] = s.toCategoryResponse(ctx, &subCategory)
		}
	}

	return response
}
