package services

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"go.uber.org/zap"
)

// SearchService Methods:
// 0. NewSearchService(repo databases.SearchRepository, logger *zap.Logger) -> 注入依赖
// 1. GlobalSearch(ctx context.Context, req *map[string]interface{}) -> 全局搜索
// 2. SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest) -> 搜索房产
// 3. SearchEstates(ctx context.Context, req *models.ListEstatesRequest) -> 搜索屋苑
// 4. SearchAgents(ctx context.Context, req *models.ListAgentsRequest) -> 搜索代理人
// 5. GetSearchSuggestions(ctx context.Context, req *map[string]interface{}) -> 获取搜索建议
// 6. GetSearchHistory(ctx context.Context, userID *uint, req *map[string]interface{}) -> 获取搜索历史

// SearchServiceInterface 定义搜索服务接口
type SearchServiceInterface interface {
	GlobalSearch(ctx context.Context, req *map[string]interface{}) (*map[string]interface{}, error)
	SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest) (*[]models.Property, error)
	SearchEstates(ctx context.Context, req *models.ListEstatesRequest) (*[]models.Estate, error)
	SearchAgents(ctx context.Context, req *models.ListAgentsRequest) (*[]models.Agent, error)
	GetSearchSuggestions(ctx context.Context, req *map[string]interface{}) (*[]string, error)
	GetSearchHistory(ctx context.Context, userID *uint, req *map[string]interface{}) (*[]models.SearchHistory, error)
}

// SearchService 搜索服务
type SearchService struct {
	repo   databases.SearchRepository
	logger *zap.Logger
}

// 0. NewSearchService 构造函数
func NewSearchService(repo databases.SearchRepository, logger *zap.Logger) *SearchService {
	return &SearchService{
		repo:   repo,
		logger: logger,
	}
}

// 1. GlobalSearch 全局搜索
func (s *SearchService) GlobalSearch(ctx context.Context, req *map[string]interface{}) (*map[string]interface{}, error) {
	keyword, _ := (*req)["keyword"].(string)
	
	// 并行搜索多个类型
	propertyChan := make(chan []*models.Property)
	estateChan := make(chan []*models.Estate)
	agentChan := make(chan []*models.Agent)
	newsChan := make(chan []*models.News)
	errChan := make(chan error, 4)

	// 搜索房产
	go func() {
		filters := &models.SearchPropertiesRequest{
			Keyword:  keyword,
			Page:     1,
			PageSize: 5,
		}
		properties, _, err := s.repo.SearchProperties(ctx, keyword, filters)
		if err != nil {
			errChan <- err
			return
		}
		propertyChan <- properties
	}()

	// 搜索屋苑
	go func() {
		estates, _, err := s.repo.SearchEstates(ctx, keyword, nil, 1, 5)
		if err != nil {
			errChan <- err
			return
		}
		estateChan <- estates
	}()

	// 搜索代理人
	go func() {
		agents, _, err := s.repo.SearchAgents(ctx, keyword, nil, 1, 5)
		if err != nil {
			errChan <- err
			return
		}
		agentChan <- agents
	}()

	// 搜索新闻
	go func() {
		news, err := s.repo.SearchNews(ctx, keyword, 5)
		if err != nil {
			errChan <- err
			return
		}
		newsChan <- news
	}()

	// 收集结果
	var properties []*models.Property
	var estates []*models.Estate
	var agents []*models.Agent
	var news []*models.News

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
	propertyResults := make([]*models.Property, len(properties))
	for i, prop := range properties {
		propertyResults[i] = convertToPropertySearchResult(prop)
	}

	estateResults := make([]*models.Estate, len(estates))
	for i, estate := range estates {
		estateResults[i] = convertToEstateSearchResult(estate)
	}

	agentResults := make([]*models.Agent, len(agents))
	for i, agent := range agents {
		agentResults[i] = convertToAgentSearchResult(agent)
	}

	newsResults := make([]*models.News, len(news))
	for i, n := range news {
		newsResults[i] = convertToNewsSearchResult(n)
	}

	totalCount := len(properties) + len(estates) + len(agents) + len(news)

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     keyword,
		SearchType:  models.SearchTypeGlobal,
		ResultCount: totalCount,
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	result := map[string]interface{}{
		"properties":  propertyResults,
		"estates":     estateResults,
		"agents":      agentResults,
		"news":        newsResults,
		"total_count": totalCount,
		"keyword":     keyword,
	}
	
	return &result, nil
}

// 2. SearchProperties 搜索房产
func (s *SearchService) SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest) (*[]models.Property, error) {
	properties, total, err := s.repo.SearchProperties(ctx, req.Keyword, req)
	if err != nil {
		s.logger.Error("Failed to search properties", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	propertyResults := make([]models.Property, len(properties))
	for i, prop := range properties {
		propertyResults[i] = *convertToPropertySearchResult(prop)
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  models.SearchTypeProperty,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &propertyResults, nil
}

// 3. SearchEstates 搜索屋苑
func (s *SearchService) SearchEstates(ctx context.Context, req *models.ListEstatesRequest) (*[]models.Estate, error) {
	// ListEstatesRequest doesn't have Keyword, use Name for search
	keyword := ""
	if req.Name != nil {
		keyword = *req.Name
	}
	
	estates, total, err := s.repo.SearchEstates(ctx, keyword, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search estates", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	estateResults := make([]models.Estate, len(estates))
	for i, estate := range estates {
		estateResults[i] = *convertToEstateSearchResult(estate)
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     keyword,
		SearchType:  models.SearchTypeEstate,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &estateResults, nil
}

// 4. SearchAgents 搜索代理人
func (s *SearchService) SearchAgents(ctx context.Context, req *models.ListAgentsRequest) (*[]models.Agent, error) {
	keyword := ""
	if req.Keyword != nil {
		keyword = *req.Keyword
	}
	
	agents, total, err := s.repo.SearchAgents(ctx, keyword, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search agents", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	agentResults := make([]models.Agent, len(agents))
	for i, agent := range agents {
		agentResults[i] = *convertToAgentSearchResult(agent)
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     keyword,
		SearchType:  models.SearchTypeAgent,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &agentResults, nil
}

// 5. GetSearchSuggestions 获取搜索建议
func (s *SearchService) GetSearchSuggestions(ctx context.Context, req *map[string]interface{}) (*[]string, error) {
	keyword, _ := (*req)["keyword"].(string)
	reqType, _ := (*req)["type"].(*string)
	limit, _ := (*req)["limit"].(int)
	if limit == 0 {
		limit = 10
	}
	
	var suggestions []string

	// 根据类型获取建议
	if reqType == nil || *reqType == "property" {
		propertySuggestions, err := s.repo.GetPropertySuggestions(ctx, keyword, limit)
		if err == nil {
			suggestions = append(suggestions, propertySuggestions...)
		}
	}

	if reqType == nil || *reqType == "estate" {
		estateSuggestions, err := s.repo.GetEstateSuggestions(ctx, keyword, limit)
		if err == nil {
			suggestions = append(suggestions, estateSuggestions...)
		}
	}

	if reqType == nil || *reqType == "agent" {
		agentSuggestions, err := s.repo.GetAgentSuggestions(ctx, keyword, limit)
		if err == nil {
			suggestions = append(suggestions, agentSuggestions...)
		}
	}

	// 限制总数
	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return &suggestions, nil
}

// 6. GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(ctx context.Context, userID *uint, req *map[string]interface{}) (*[]models.SearchHistory, error) {
	reqType, _ := (*req)["type"].(*string)
	page, _ := (*req)["page"].(int)
	pageSize, _ := (*req)["page_size"].(int)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	histories, _, err := s.repo.GetSearchHistory(ctx, userID, reqType, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get search history", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	historyResults := make([]models.SearchHistory, len(histories))
	for i, history := range histories {
		historyResults[i] = *history
	}

	return &historyResults, nil
}

// Helper functions for converting models to response DTOs

func convertToPropertySearchResult(prop *models.Property) *models.Property {
	result := &models.Property{
		ID:          prop.ID,
		Title:       prop.Title,
		Price:       prop.Price,
		Area:        prop.Area,
		Bedrooms:    prop.Bedrooms,
		Bathrooms:   prop.Bathrooms,
		ListingType: prop.ListingType,
		Address:     prop.Address,
		CreatedAt:   prop.CreatedAt,
		District:    prop.District,
		Images:      prop.Images,
	}

	return result
}

func convertToEstateSearchResult(estate *models.Estate) *models.Estate {
	result := &models.Estate{
		ID:         estate.ID,
		Name:       estate.Name,
		NameEn:     estate.NameEn,
		Address:    estate.Address,
		TotalBlocks: estate.TotalBlocks,
		TotalUnits:  estate.TotalUnits,
		District:   estate.District,
		Images:     estate.Images,
	}

	return result
}

func convertToAgentSearchResult(agent *models.Agent) *models.Agent {
	result := &models.Agent{
		ID:          agent.ID,
		AgentName:   agent.AgentName,
		LicenseNo:   agent.LicenseNo,
		Phone:       agent.Phone,
		Email:       agent.Email,
		ProfilePhoto: agent.ProfilePhoto,
		Agency:      agent.Agency,
	}

	return result
}

func convertToNewsSearchResult(news *models.News) *models.News {
	result := &models.News{
		ID:            news.ID,
		Title:         news.Title,
		Summary:       news.Summary,
		CoverImageURL: news.CoverImageURL,
		ViewCount:     news.ViewCount,
		Category:      news.Category,
		PublishedAt:   news.PublishedAt,
	}

	return result
}
