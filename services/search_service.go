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
	// 并行搜索多个类型
	propertyChan := make(chan []*models.Property)
	estateChan := make(chan []*models.Estate)
	agentChan := make(chan []*models.Agent)
	newsChan := make(chan []*models.News)
	errChan := make(chan error, 4)

	// 搜索房产
	go func() {
		filters := &models.SearchPropertiesRequest{
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
		Keyword:     req.Keyword,
		SearchType:  models.SearchTypeGlobal,
		ResultCount: totalCount,
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &map[string]interface{}{
		Properties: propertyResults,
		Estates:    estateResults,
		Agents:     agentResults,
		News:       newsResults,
		TotalCount: totalCount,
		Keyword:    req.Keyword,
	}, nil
}

// 2. SearchProperties 搜索房产
func (s *SearchService) SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest) (*[]models.Property, error) {
	properties, total, err := s.repo.SearchProperties(ctx, req.Keyword, req)
	if err != nil {
		s.logger.Error("Failed to search properties", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	propertyResults := make([]*models.Property, len(properties))
	for i, prop := range properties {
		propertyResults[i] = convertToPropertySearchResult(prop)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  models.SearchTypeProperty,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &[]models.Property{
		Properties: propertyResults,
		Pagination: &models.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
		Keyword: req.Keyword,
	}, nil
}

// 3. SearchEstates 搜索屋苑
func (s *SearchService) SearchEstates(ctx context.Context, req *models.ListEstatesRequest) (*[]models.Estate, error) {
	estates, total, err := s.repo.SearchEstates(ctx, req.Keyword, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search estates", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	estateResults := make([]*models.Estate, len(estates))
	for i, estate := range estates {
		estateResults[i] = convertToEstateSearchResult(estate)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  models.SearchTypeEstate,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &[]models.Estate{
		Estates: estateResults,
		Pagination: &models.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
		Keyword: req.Keyword,
	}, nil
}

// 4. SearchAgents 搜索代理人
func (s *SearchService) SearchAgents(ctx context.Context, req *models.ListAgentsRequest) (*[]models.Agent, error) {
	agents, total, err := s.repo.SearchAgents(ctx, req.Keyword, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to search agents", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	agentResults := make([]*models.Agent, len(agents))
	for i, agent := range agents {
		agentResults[i] = convertToAgentSearchResult(agent)
	}

	// 计算总页数
	totalPage := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPage++
	}

	// 保存搜索历史
	history := &models.SearchHistory{
		Keyword:     req.Keyword,
		SearchType:  models.SearchTypeAgent,
		ResultCount: int(total),
	}
	_ = s.repo.SaveSearchHistory(ctx, history)

	return &[]models.Agent{
		Agents: agentResults,
		Pagination: &models.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
		Keyword: req.Keyword,
	}, nil
}

// 5. GetSearchSuggestions 获取搜索建议
func (s *SearchService) GetSearchSuggestions(ctx context.Context, req *map[string]interface{}) (*[]string, error) {
	var suggestions []*models.SearchSuggestion

	// 根据类型获取建议
	if req.Type == nil || *req.Type == "property" {
		propertySuggestions, err := s.repo.GetPropertySuggestions(ctx, req.Keyword, req.Limit)
		if err == nil {
			for _, text := range propertySuggestions {
				suggestions = append(suggestions, &models.SearchSuggestion{
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
				suggestions = append(suggestions, &models.SearchSuggestion{
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
				suggestions = append(suggestions, &models.SearchSuggestion{
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

	return &[]string{
		Suggestions: suggestions,
		Keyword:     req.Keyword,
	}, nil
}

// 6. GetSearchHistory 获取搜索历史
func (s *SearchService) GetSearchHistory(ctx context.Context, userID *uint, req *map[string]interface{}) (*[]models.SearchHistory, error) {
	histories, total, err := s.repo.GetSearchHistory(ctx, userID, req.Type, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to get search history", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	historyItems := make([]*models.SearchHistoryItem, len(histories))
	for i, history := range histories {
		historyItems[i] = &models.SearchHistoryItem{
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

	return &[]models.SearchHistory{
		Histories: historyItems,
		Pagination: &models.Pagination{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
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

func convertToEstateSearchResult(estate *models.Estate) *models.Estate {
	result := &models.Estate{
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

func convertToAgentSearchResult(agent *models.Agent) *models.Agent {
	result := &models.Agent{
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

func convertToNewsSearchResult(news *models.News) *models.News {
	result := &models.News{
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
