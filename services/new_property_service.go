package services

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// NewDevelopmentService 新盘服务
// Methods:
// 1. ListNewProperties(ctx context.Context, filter *models.ListNewPropertiesRequest) -> 获取新盘列表
// 2. GetNewProperty(ctx context.Context, id uint) -> 获取新盘详情
// 3. GetNewPropertyLayouts(ctx context.Context, newPropertyID uint) -> 获取新盘户型列表
type NewDevelopmentService struct {
	repo *databases.NewDevelopmentRepo
}

// NewNewDevelopmentService 创建新盘服务
func NewNewDevelopmentService(repo *databases.NewDevelopmentRepo) *NewDevelopmentService {
	return &NewDevelopmentService{repo: repo}
}

// ListNewProperties 获取新盘列表
func (s *NewDevelopmentService) ListNewProperties(ctx context.Context, filter *models.ListNewPropertiesRequest) (*models.PaginatedNewPropertiesResponse, error) {
	// 默认参数
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	newProperties, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]models.NewPropertyResponse, len(newProperties))
	for i, np := range newProperties {
		items[i] = *np.ToNewPropertyResponse()
	}

	// 计算总页数
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedNewPropertiesResponse{
		Data:       items,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetNewProperty 获取新盘详情
func (s *NewDevelopmentService) GetNewProperty(ctx context.Context, id uint) (*models.NewPropertyDetailResponse, error) {
	newProperty, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 异步增加浏览次数
	go func() {
		_ = s.repo.IncrementViewCount(context.Background(), id)
	}()

	return newProperty.ToNewPropertyDetailResponse(), nil
}

// GetNewPropertyLayouts 获取新盘户型列表
func (s *NewDevelopmentService) GetNewPropertyLayouts(ctx context.Context, newPropertyID uint) ([]models.NewPropertyLayout, error) {
	// 先检查新盘是否存在
	_, err := s.repo.FindByID(ctx, newPropertyID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindLayouts(ctx, newPropertyID)
}
