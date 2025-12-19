package service

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
	"go.uber.org/zap"
)

// SearchService Methods:
// 0. NewSearchService(repo repository.SearchRepository, logger *zap.Logger) -> 注入依赖
// 1. GlobalSearch(ctx context.Context, req *request.GlobalSearchRequest) -> 全局搜索
// 2. SearchProperties(ctx context.Context, req *request.SearchPropertiesRequest) -> 搜索房产
// 3. SearchEstates(ctx context.Context, req *request.SearchEstatesRequest) -> 搜索屋苑
// 4. SearchAgents(ctx context.Context, req *request.SearchAgentsRequest) -> 搜索代理人
// 5. GetSearchSuggestions(ctx context.Context, req *request.GetSearchSuggestionsRequest) -> 获取搜索建议
// 6. GetSearchHistory(ctx context.Context, userID *uint, req *request.GetSearchHistoryRequest) -> 获取搜索历史

// SearchServiceInterface 定义搜索服务接口
type SearchServiceInterface interface {
	GlobalSearch(ctx context.Context, req *request.GlobalSearchRequest) (*response.GlobalSearchResponse, error)
	SearchProperties(ctx context.Context, req *request.SearchPropertiesRequest) (*response.SearchPropertiesResponse, error)
	SearchEstates(ctx context.Context, req *request.SearchEstatesRequest) (*response.SearchEstatesResponse, error)
	SearchAgents(ctx context.Context, req *request.SearchAgentsRequest) (*response.SearchAgentsResponse, error)
	GetSearchSuggestions(ctx context.Context, req *request.GetSearchSuggestionsRequest) (*response.SearchSuggestionsResponse, error)
	GetSearchHistory(ctx context.Context, userID *uint, req *request.GetSearchHistoryRequest) (*response.SearchHistoryResponse, error)
}

// SearchService 搜索服务
type SearchService struct {
	repo   repository.SearchRepository
	logger *zap.Logger
}

// 0. NewSearchService 构造函数
func NewSearchService(repo repository.SearchRepository, logger *zap.Logger) *SearchService {
	return &SearchService{
		repo:   repo,
		logger: logger,
	}
}

// 1. GlobalSearch 全局搜索
func (s *SearchService) GlobalSearch(ctx context.Context, req *request.GlobalSearchRequest) (*response.GlobalSearchResponse, error) {
	// 并行搜索多个类型
	propertyChan := make(chan []*model.Property)
	estateChan := make(chan []*model.Estate)
	agentChan := make(chan []*model.Agent)
	newsChan := make(chan []*model.News)
	errChan := make(chan error, 4)

	// 搜索房产
	go func() {
		filters := &request.SearchPropertiesRequest{
			Keyword:  req.Keyword,
			Page:     1,
			PageSize: 5,
		}
		properties, _, err := s.repo.SearchProperties(ctx, req.Keyword, filters)
		if err != nil {
			errChan <- err
			return
		}
		propertyChan <- properties
	}()

	// 搜索屋苑
	go func() {
		estates, _, err := s.repo.SearchEstates(ctx, req.Keyword, nil, 1, 5)
		if err != nil {
			errChan <- err
			return
		}
		estateChan <- estates
	}()

	// 搜索代理人
	go func() {
		agents, _, err := s.repo.SearchAgents(ctx, req.Keyword, nil, 1, 5)
		if err != nil {
			errChan <- err
			return
		}
		agentChan <- agents
	}()

	// 搜索新闻
	go func() {
		news, err := s.repo.SearchNews(ctx, req.Keyword, 5)
		if err != nil {
			errChan <- err
			return
		}
		newsChan <- news
	}()

	// 收集结果
	var properties []*model.Property
	var estates []*model.Estate
	var agents []*model.Agent
	var news []*model.News

	for i := 0; i < 4; i++ {
		select {
		case p := <-propertyChan:
			properties = p
		case e := <-estateChan:
			estates = e
		case a := <-agentChan:
			agents = a
		case n := <-newsChan:
			news = n
		case err := <-errChan:
			s.logger.Error("Global search error", zap.Error(err))
			return nil, err
		}
	}

	// 转换为响应格式
	propertyResults := make([]*response.PropertySearchResult, len(properties))
	for i, prop := range properties {
		propertyResults[i] = convertToPropertySearchResult(prop)
	}

	estateResults := make([]*response.EstateSearchResult, len(estates))
	for i, estate := range estates {
		estateResults[i] = convertToEstateSearchResult(estate)
	}

	agentResults := make([]*response.AgentSearchResult, len(agents))
	for i, agent := range agents {
		agentResults[i] = convertToAgentSearchResult(agent)
	}

	newsResults := make([]*response.NewsSearchResult, len(news))
	for i, n := range news {
		newsResults[i] = convertToNewsSearchResult(n)
	}

	totalCount := len(properties) + len(estates) + len(agents) + len(news)

	// 保存搜索历史
	history := &model.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  model.SearchTypeGlobal,
		ResultCount: totalCount,
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &response.GlobalSearchResponse{
		Properties: propertyResults,
		Estates:    estateResults,
		Agents:     agentResults,
		News:       newsResults,
		TotalCount: totalCount,
		Keyword:    req.Keyword,
	}, nil
}

// 2. SearchProperties 搜索房产
func (s *SearchService) SearchProperties(ctx context.Context, req *request.SearchPropertiesRequest) (*response.SearchPropertiesResponse, error) {
	properties, total, err := s.repo.SearchProperties(ctx, req.Keyword, req)
	if err != nil {
		s.logger.Error("Failed to search properties", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	propertyResults := make([]*response.PropertySearchResult, len(properties))
	for i, prop := range properties {
		propertyResults[i] = convertToPropertySearchResult(prop)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	// 保存搜索历史
	history := &model.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  model.SearchTypeProperty,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &response.SearchPropertiesResponse{
		Properties: propertyResults,
		Pagination: &response.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
		Keyword: req.Keyword,
	}, nil
}

// 3. SearchEstates 搜索屋苑
func (s *SearchService) SearchEstates(ctx context.Context, req *request.SearchEstatesRequest) (*response.SearchEstatesResponse, error) {
	estates, total, err := s.repo.SearchEstates(ctx, req.Keyword, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search estates", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	estateResults := make([]*response.EstateSearchResult, len(estates))
	for i, estate := range estates {
		estateResults[i] = convertToEstateSearchResult(estate)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	// 保存搜索历史
	history := &model.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  model.SearchTypeEstate,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &response.SearchEstatesResponse{
		Estates: estateResults,
		Pagination: &response.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
		Keyword: req.Keyword,
	}, nil
}

// 4. SearchAgents 搜索代理人
func (s *SearchService) SearchAgents(ctx context.Context, req *request.SearchAgentsRequest) (*response.SearchAgentsResponse, error) {
	agents, total, err := s.repo.SearchAgents(ctx, req.Keyword, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search agents", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	agentResults := make([]*response.AgentSearchResult, len(agents))
	for i, agent := range agents {
		agentResults[i] = convertToAgentSearchResult(agent)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	// 保存搜索历史
	history := &model.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  model.SearchTypeAgent,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &response.SearchAgentsResponse{
		Agents: agentResults,
		Pagination: &response.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
		Keyword: req.Keyword,
	}, nil
}

// 5. GetSearchSuggestions 获取搜索建议
func (s *SearchService) GetSearchSuggestions(ctx context.Context, req *request.GetSearchSuggestionsRequest) (*response.SearchSuggestionsResponse, error) {
	var suggestions []*response.SearchSuggestion

	// 根据类型获取建议
	if req.Type == nil || *req.Type == "property" {
		propertySuggestions, err := s.repo.GetPropertySuggestions(ctx, req.Keyword, req.Limit)
		if err == nil {
			for _, text := range propertySuggestions {
				suggestions = append(suggestions, &response.SearchSuggestion{
					Text: text,
					Type: "property",
				})
			}
		}
	}

	if req.Type == nil || *req.Type == "estate" {
		estateSuggestions, err := s.repo.GetEstateSuggestions(ctx, req.Keyword, req.Limit)
		if err == nil {
			for _, text := range estateSuggestions {
				suggestions = append(suggestions, &response.SearchSuggestion{
					Text: text,
					Type: "estate",
				})
			}
		}
	}

	if req.Type == nil || *req.Type == "agent" {
		agentSuggestions, err := s.repo.GetAgentSuggestions(ctx, req.Keyword, req.Limit)
		if err == nil {
			for _, text := range agentSuggestions {
				suggestions = append(suggestions, &response.SearchSuggestion{
					Text: text,
					Type: "agent",
				})
			}
		}
	}

	// 限制总数
	if len(suggestions) > req.Limit {
		suggestions = suggestions[:req.Limit]
	}

	return &response.SearchSuggestionsResponse{
		Suggestions: suggestions,
		Keyword:     req.Keyword,
	}, nil
}

// 6. GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(ctx context.Context, userID *uint, req *request.GetSearchHistoryRequest) (*response.SearchHistoryResponse, error) {
	histories, total, err := s.repo.GetSearchHistory(ctx, userID, req.Type, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to get search history", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	historyItems := make([]*response.SearchHistoryItem, len(histories))
	for i, history := range histories {
		historyItems[i] = &response.SearchHistoryItem{
			ID:          history.ID,
			Keyword:     history.Keyword,
			SearchType:  string(history.SearchType),
			ResultCount: history.ResultCount,
			CreatedAt:   history.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	return &response.SearchHistoryResponse{
		Histories: historyItems,
		Pagination: &response.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

// Helper functions for converting models to response DTOs

func convertToPropertySearchResult(prop *model.Property) *response.PropertySearchResult {
	result := &response.PropertySearchResult{
		ID:          prop.ID,
		Title:       prop.Title,
		Price:       prop.Price,
		Area:        prop.Area,
		Bedrooms:    prop.Bedrooms,
		Bathrooms:   prop.Bathrooms,
		ListingType: prop.ListingType,
		Address:     prop.Address,
		CreatedAt:   prop.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if prop.District != nil {
		result.District = &prop.District.NameZhHant
	}

	if prop.Estate != nil {
		result.EstateName = &prop.Estate.Name
	}

	if len(prop.Images) > 0 {
		result.ImageURL = &prop.Images[0].URL
	}

	return result
}

func convertToEstateSearchResult(estate *model.Estate) *response.EstateSearchResult {
	result := &response.EstateSearchResult{
		ID:              estate.ID,
		Name:            estate.Name,
		NameEn:          estate.NameEn,
		Address:         estate.Address,
		BuildingCount:   estate.BuildingCount,
		UnitCount:       estate.UnitCount,
		PropertyCount:   0, // TODO: 从关联查询获取
	}

	if estate.District != nil {
		result.District = estate.District.NameZhHant
	}

	if len(estate.Images) > 0 {
		result.ImageURL = &estate.Images[0].ImageURL
	}

	return result
}

func convertToAgentSearchResult(agent *model.Agent) *response.AgentSearchResult {
	result := &response.AgentSearchResult{
		ID:            agent.ID,
		Name:          agent.Name,
		LicenseNo:     agent.LicenseNo,
		Phone:         agent.Phone,
		Email:         agent.Email,
		AvatarURL:     agent.AvatarURL,
		PropertyCount: 0, // TODO: 从关联查询获取
		Rating:        0, // TODO: 计算评分
	}

	if agent.Agency != nil {
		result.AgencyName = &agent.Agency.CompanyName
	}

	return result
}

func convertToNewsSearchResult(news *model.News) *response.NewsSearchResult {
	result := &response.NewsSearchResult{
		ID:        news.ID,
		Title:     news.Title,
		Summary:   news.Summary,
		ImageURL:  news.ImageURL,
		ViewCount: news.ViewCount,
	}

	if news.Category != nil {
		result.Category = news.Category.Name
	}

	if news.PublishedAt != nil {
		result.PublishedAt = news.PublishedAt.Format("2006-01-02 15:04:05")
	}

	return result
}
