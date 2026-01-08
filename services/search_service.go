package services

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// SearchService Methods:
// 0. NewSearchService(repo *databases.SearchRepo) -> 注入依赖
// 1. GlobalSearch(ctx context.Context, req *models.GlobalSearchRequest, userID *uint, ipAddress, userAgent string) -> 全局搜索
// 2. SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest, userID *uint, ipAddress, userAgent string) -> 搜索房产
// 3. SearchEstates(ctx context.Context, req *models.SearchEstatesRequest, userID *uint, ipAddress, userAgent string) -> 搜索屋苑
// 4. SearchAgents(ctx context.Context, req *models.SearchAgentsRequest, userID *uint, ipAddress, userAgent string) -> 搜索代理人
// 5. GetSearchSuggestions(ctx context.Context, req *models.GetSearchSuggestionsRequest) -> 获取搜索建议
// 6. GetSearchHistory(ctx context.Context, userID *uint, req *models.GetSearchHistoryRequest) -> 获取搜索历史

type SearchService struct {
	repo *databases.SearchRepo
}

// 0. NewSearchService 构造函数
func NewSearchService(repo *databases.SearchRepo) *SearchService {
	return &SearchService{repo: repo}
}

// 1. GlobalSearch 全局搜索
func (s *SearchService) GlobalSearch(ctx context.Context, req *models.GlobalSearchRequest, userID *uint, ipAddress, userAgent string) (*models.GlobalSearchResponse, error) {
	// 执行搜索
	results, err := s.repo.GlobalSearch(ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		UserID:      userID,
		Keyword:     req.Keyword,
		SearchType:  "global",
		ResultCount: results.TotalResults,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}
	s.repo.SaveSearchHistory(ctx, history)

	return results, nil
}

// 2. SearchProperties 搜索房产
func (s *SearchService) SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest, userID *uint, ipAddress, userAgent string) ([]models.PropertySearchResult, int64, error) {
	// 执行搜索
	results, total, err := s.repo.SearchProperties(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		UserID:      userID,
		Keyword:     req.Keyword,
		SearchType:  "property",
		ResultCount: int(total),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}
	s.repo.SaveSearchHistory(ctx, history)

	return results, total, nil
}

// 3. SearchEstates 搜索屋苑
func (s *SearchService) SearchEstates(ctx context.Context, req *models.SearchEstatesRequest, userID *uint, ipAddress, userAgent string) ([]models.EstateSearchResult, int64, error) {
	// 执行搜索
	results, total, err := s.repo.SearchEstates(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		UserID:      userID,
		Keyword:     req.Keyword,
		SearchType:  "estate",
		ResultCount: int(total),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}
	s.repo.SaveSearchHistory(ctx, history)

	return results, total, nil
}

// 4. SearchAgents 搜索代理人
func (s *SearchService) SearchAgents(ctx context.Context, req *models.SearchAgentsRequest, userID *uint, ipAddress, userAgent string) ([]models.AgentSearchResult, int64, error) {
	// 执行搜索
	results, total, err := s.repo.SearchAgents(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		UserID:      userID,
		Keyword:     req.Keyword,
		SearchType:  "agent",
		ResultCount: int(total),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}
	s.repo.SaveSearchHistory(ctx, history)

	return results, total, nil
}

// 5. GetSearchSuggestions 获取搜索建议
func (s *SearchService) GetSearchSuggestions(ctx context.Context, req *models.GetSearchSuggestionsRequest) ([]models.SearchSuggestion, error) {
	return s.repo.GetSearchSuggestions(ctx, req.Keyword, req.Limit)
}

// 6. GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(ctx context.Context, userID *uint, req *models.GetSearchHistoryRequest) (*models.PaginatedSearchHistoryResponse, error) {
	histories, total, err := s.repo.GetSearchHistory(ctx, userID, req.SearchType, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var historyResponses []models.SearchHistoryResponse
	for _, h := range histories {
		historyResponses = append(historyResponses, models.SearchHistoryResponse{
			ID:          h.ID,
			Keyword:     h.Keyword,
			SearchType:  h.SearchType,
			ResultCount: h.ResultCount,
			CreatedAt:   h.CreatedAt,
		})
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedSearchHistoryResponse{
		Histories:  historyResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// DeleteSearchHistory 删除搜索历史
func (s *SearchService) DeleteSearchHistory(ctx context.Context, userID uint) error {
	return s.repo.DeleteSearchHistory(ctx, userID)
}
