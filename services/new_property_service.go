package services

import (
	"context"
	"errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	pkgErrors "github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
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
	ListNewDevelopments(ctx context.Context, req *models.ListNewDevelopmentsRequest) ([]*models.NewProperty, int64, error)
	GetNewDevelopment(ctx context.Context, id uint) (*models.NewProperty, error)
	GetDevelopmentUnits(ctx context.Context, id uint) ([]models.NewPropertyLayout, error)
	GetFeaturedNewDevelopments(ctx context.Context, limit int) ([]*models.NewProperty, error)
}

// newPropertyService 新楼盘服务实现
type newPropertyService struct {
	repo   databases.NewPropertyRepository
	logger *zap.Logger
}

// NewNewPropertyService 创建新楼盘服务实例
func NewNewPropertyService(repo databases.NewPropertyRepository, logger *zap.Logger) NewPropertyService {
	return &newPropertyService{
		repo:   repo,
		logger: logger,
	}
}

// ListNewDevelopments 获取新楼盘列表
func (s *newPropertyService) ListNewDevelopments(ctx context.Context, req *models.ListNewDevelopmentsRequest) ([]*models.NewProperty, int64, error) {
	// 获取新楼盘列表
	newProperties, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list new developments", zap.Error(err))
		return nil, 0, err
	}

	// 转换为响应格式
	items := make([]*models.NewProperty, len(newProperties))
	for i, np := range newProperties {
		items[i] = s.toNewDevelopmentListItemResponse(np)
	}

	return items, total, nil
}

// GetNewDevelopment 获取新楼盘详情
func (s *newPropertyService) GetNewDevelopment(ctx context.Context, id uint) (*models.NewProperty, error) {
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
func (s *newPropertyService) GetDevelopmentUnits(ctx context.Context, id uint) ([]models.NewPropertyLayout, error) {
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
	result := make([]models.NewPropertyLayout, len(layouts))
	for i, layout := range layouts {
		result[i] = s.toLayoutResponse(&layout)
	}

	return result, nil
}

// GetFeaturedNewDevelopments 获取精选新楼盘
func (s *newPropertyService) GetFeaturedNewDevelopments(ctx context.Context, limit int) ([]*models.NewProperty, error) {
	newProperties, err := s.repo.GetFeatured(ctx, limit)
	if err != nil {
		s.logger.Error("Failed to get featured new developments", zap.Error(err))
		return nil, err
	}

	items := make([]*models.NewProperty, len(newProperties))
	for i, np := range newProperties {
		items[i] = s.toNewDevelopmentListItemResponse(np)
	}

	return items, nil
}

// toNewDevelopmentResponse 转换为新楼盘详情响应（直接返回，预加载了关联数据）
func (s *newPropertyService) toNewDevelopmentResponse(np *models.NewProperty) *models.NewProperty {
	return np
}

// toNewDevelopmentListItemResponse 转换为新楼盘列表项响应（直接返回，预加载了关联数据）
func (s *newPropertyService) toNewDevelopmentListItemResponse(np *models.NewProperty) *models.NewProperty {
	return np
}

// toLayoutResponse 转换为户型响应（直接返回模型）
func (s *newPropertyService) toLayoutResponse(layout *models.NewPropertyLayout) models.NewPropertyLayout {
	return *layout
}

// 辅助函数已不再需要，但保留以防其他地方使用
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
