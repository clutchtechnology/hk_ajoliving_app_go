package services

import (
	"context"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

// NewsService 新闻服务接口
type NewsService interface {
	ListNews(ctx context.Context, filter *models.ListNewsRequest) ([]*models.News, int64, error)
	GetNews(ctx context.Context, id uint) (*models.News, error)
	GetHotNews(ctx context.Context) ([]*models.News, error)
	GetFeaturedNews(ctx context.Context) ([]*models.News, error)
	GetLatestNews(ctx context.Context) ([]*models.News, error)
	GetRelatedNews(ctx context.Context, id uint) ([]*[]models.News, error)
	GetNewsCategories(ctx context.Context) ([]*models.NewsCategory, error)
}

type newsService struct {
	repo   databases.NewsRepository
	logger *zap.Logger
}

// NewNewsService 创建新闻服务
func NewNewsService(repo databases.NewsRepository, logger *zap.Logger) NewsService {
	return &newsService{
		repo:   repo,
		logger: logger,
	}
}

// ListNews 获取新闻列表
func (s *newsService) ListNews(ctx context.Context, filter *models.ListNewsRequest) ([]*models.News, int64, error) {
	news, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list news", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.News, 0, len(news))
	for _, n := range news {
		result = append(result, convertToNewsListItemResponse(n))
	}
	
	return result, total, nil
}

// GetNews 获取新闻详情
func (s *newsService) GetNews(ctx context.Context, id uint) (*models.News, error) {
	news, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get news", zap.Error(err))
		return nil, err
	}
	if news == nil {
		return nil, tools.ErrNotFound
	}
	
	// 增加浏览量
	if err := s.repo.IncrementViewCount(ctx, id); err != nil {
		s.logger.Warn("failed to increment view count", zap.Uint("news_id", id), zap.Error(err))
	}
	
	return convertToNewsResponse(news), nil
}

// GetHotNews 获取热门新闻
func (s *newsService) GetHotNews(ctx context.Context) ([]*models.News, error) {
	news, err := s.repo.GetHotNews(ctx, 10)
	if err != nil {
		s.logger.Error("failed to get hot news", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.News, 0, len(news))
	for _, n := range news {
		result = append(result, convertToNewsListItemResponse(n))
	}
	
	return result, nil
}

// GetFeaturedNews 获取精选新闻
func (s *newsService) GetFeaturedNews(ctx context.Context) ([]*models.News, error) {
	news, err := s.repo.GetFeaturedNews(ctx, 10)
	if err != nil {
		s.logger.Error("failed to get featured news", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.News, 0, len(news))
	for _, n := range news {
		result = append(result, convertToNewsListItemResponse(n))
	}
	
	return result, nil
}

// GetLatestNews 获取最新新闻
func (s *newsService) GetLatestNews(ctx context.Context) ([]*models.News, error) {
	news, err := s.repo.GetLatestNews(ctx, 10)
	if err != nil {
		s.logger.Error("failed to get latest news", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.News, 0, len(news))
	for _, n := range news {
		result = append(result, convertToNewsListItemResponse(n))
	}
	
	return result, nil
}

// GetRelatedNews 获取相关新闻
func (s *newsService) GetRelatedNews(ctx context.Context, id uint) ([]*[]models.News, error) {
	// 先获取当前新闻
	currentNews, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get current news", zap.Error(err))
		return nil, err
	}
	if currentNews == nil {
		return nil, tools.ErrNotFound
	}
	
	// 获取同分类的相关新闻
	relatedNews, err := s.repo.GetRelatedNews(ctx, id, currentNews.CategoryID, 5)
	if err != nil {
		s.logger.Error("failed to get related news", zap.Error(err))
		return nil, err
	}
	
	// 转换为所需的返回类型
	newsSlice := make([]models.News, len(relatedNews))
	for i, n := range relatedNews {
		newsSlice[i] = *n
	}
	result := []*[]models.News{&newsSlice}
	
	return result, nil
}

// GetNewsCategories 获取新闻分类列表
func (s *newsService) GetNewsCategories(ctx context.Context) ([]*models.NewsCategory, error) {
	categories, err := s.repo.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error("failed to get news categories", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.NewsCategory, 0, len(categories))
	for _, c := range categories {
		result = append(result, convertToNewsCategoryResponse(c))
	}
	
	return result, nil
}

// 辅助函数

// convertToNewsListItemResponse 转换为新闻列表项响应（直接返回，预加载了关联数据）
func convertToNewsListItemResponse(news *models.News) *models.News {
	return news
}

// convertToNewsResponse 转换为新闻详情响应（直接返回，预加载了关联数据）
func convertToNewsResponse(news *models.News) *models.News {
	return news
}

// convertToRelatedNewsResponse 转换为相关新闻响应（已简化，保留供参考）
func convertToRelatedNewsResponse(news *models.News) *models.News {
	return news
}

// convertToNewsCategoryResponse 转换为新闻分类响应（直接返回）
func convertToNewsCategoryResponse(category *models.NewsCategory) *models.NewsCategory {
	return category
}
