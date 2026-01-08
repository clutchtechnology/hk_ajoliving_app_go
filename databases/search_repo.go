package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type SearchRepo struct {
	db *gorm.DB
}

func NewSearchRepo(db *gorm.DB) *SearchRepo {
	return &SearchRepo{db: db}
}

// GlobalSearch 全局搜索
func (r *SearchRepo) GlobalSearch(ctx context.Context, keyword string, page, pageSize int) (*models.GlobalSearchResponse, error) {
	response := &models.GlobalSearchResponse{
		Properties: []models.PropertySearchResult{},
		Estates:    []models.EstateSearchResult{},
		Agents:     []models.AgentSearchResult{},
		Agencies:   []models.AgencySearchResult{},
	}

	searchPattern := "%" + keyword + "%"

	// 搜索房产（最多返回5条）
	var properties []models.Property
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("(title LIKE ? OR address LIKE ? OR building_name LIKE ?) AND status = ?",
			searchPattern, searchPattern, searchPattern, "available").
		Preload("District").
		Limit(5).
		Find(&properties)

	for _, p := range properties {
		districtName := ""
		if p.District != nil {
			districtName = p.District.NameZhHant
		}
		response.Properties = append(response.Properties, models.PropertySearchResult{
			ID:           p.ID,
			PropertyNo:   p.PropertyNo,
			Title:        p.Title,
			Price:        p.Price,
			Area:         p.Area,
			Bedrooms:     p.Bedrooms,
			PropertyType: p.PropertyType,
			ListingType:  p.ListingType,
			Address:      p.Address,
			DistrictName: districtName,
			Status:       p.Status,
		})
	}
	response.PropertyCount = len(response.Properties)

	// 搜索屋苑（最多返回5条）
	var estates []models.Estate
	r.db.WithContext(ctx).Model(&models.Estate{}).
		Where("name LIKE ? OR name_en LIKE ? OR address LIKE ?",
			searchPattern, searchPattern, searchPattern).
		Preload("District").
		Limit(5).
		Find(&estates)

	for _, e := range estates {
		districtName := ""
		if e.District != nil {
			districtName = e.District.NameZhHant
		}
		response.Estates = append(response.Estates, models.EstateSearchResult{
			ID:                  e.ID,
			Name:                e.Name,
			NameEn:              e.NameEn,
			Address:             e.Address,
			DistrictName:        districtName,
			TotalUnits:          e.TotalUnits,
			AvgTransactionPrice: e.AvgTransactionPrice,
			CompletionYear:      e.CompletionYear,
		})
	}
	response.EstateCount = len(response.Estates)

	// 搜索代理人（最多返回5条）
	var agents []models.Agent
	r.db.WithContext(ctx).Model(&models.Agent{}).
		Where("(agent_name LIKE ? OR agent_name_en LIKE ? OR license_no = ? OR specialization LIKE ?) AND status = ?",
			searchPattern, searchPattern, keyword, searchPattern, "active").
		Limit(5).
		Find(&agents)

	for _, a := range agents {
		response.Agents = append(response.Agents, models.AgentSearchResult{
			ID:             a.ID,
			AgentName:      a.AgentName,
			AgentNameEn:    a.AgentNameEn,
			LicenseNo:      a.LicenseNo,
			Phone:          a.Phone,
			Email:          a.Email,
			Specialization: a.Specialization,
			Rating:         a.Rating,
			PropertiesSold: a.PropertiesSold,
			ProfilePhoto:   a.ProfilePhoto,
		})
	}
	response.AgentCount = len(response.Agents)

	// 搜索代理公司（最多返回5条）
	var agencies []models.AgencyDetail
	r.db.WithContext(ctx).Model(&models.AgencyDetail{}).
		Where("company_name LIKE ? OR company_name_en LIKE ? OR license_no = ?",
			searchPattern, searchPattern, keyword).
		Limit(5).
		Find(&agencies)

	for _, a := range agencies {
		response.Agencies = append(response.Agencies, models.AgencySearchResult{
			ID:            a.ID,
			CompanyName:   a.CompanyName,
			CompanyNameEn: a.CompanyNameEn,
			LicenseNo:     a.LicenseNo,
			Phone:         a.Phone,
			Address:       a.Address,
			AgentCount:    a.AgentCount,
			Rating:        a.Rating,
			IsVerified:    a.IsVerified,
			LogoURL:       a.LogoURL,
		})
	}
	response.AgencyCount = len(response.Agencies)

	response.TotalResults = response.PropertyCount + response.EstateCount + response.AgentCount + response.AgencyCount

	return response, nil
}

// SearchProperties 搜索房产
func (r *SearchRepo) SearchProperties(ctx context.Context, req *models.SearchPropertiesRequest) ([]models.PropertySearchResult, int64, error) {
	var properties []models.Property
	var total int64

	searchPattern := "%" + req.Keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("(title LIKE ? OR address LIKE ? OR building_name LIKE ?) AND status = ?",
			searchPattern, searchPattern, searchPattern, "available")

	// 应用筛选条件
	if req.ListingType != nil {
		query = query.Where("listing_type = ?", *req.ListingType)
	}
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}
	if req.Bedrooms != nil {
		query = query.Where("bedrooms = ?", *req.Bedrooms)
	}
	if req.PropertyType != nil {
		query = query.Where("property_type = ?", *req.PropertyType)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Preload("District").
		Order("created_at DESC").
		Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	var results []models.PropertySearchResult
	for _, p := range properties {
		districtName := ""
		if p.District != nil {
			districtName = p.District.NameZhHant
		}
		results = append(results, models.PropertySearchResult{
			ID:           p.ID,
			PropertyNo:   p.PropertyNo,
			Title:        p.Title,
			Price:        p.Price,
			Area:         p.Area,
			Bedrooms:     p.Bedrooms,
			PropertyType: p.PropertyType,
			ListingType:  p.ListingType,
			Address:      p.Address,
			DistrictName: districtName,
			Status:       p.Status,
		})
	}

	return results, total, nil
}

// SearchEstates 搜索屋苑
func (r *SearchRepo) SearchEstates(ctx context.Context, req *models.SearchEstatesRequest) ([]models.EstateSearchResult, int64, error) {
	var estates []models.Estate
	var total int64

	searchPattern := "%" + req.Keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.Estate{}).
		Where("name LIKE ? OR name_en LIKE ? OR address LIKE ?",
			searchPattern, searchPattern, searchPattern)

	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Preload("District").
		Order("name ASC").
		Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	var results []models.EstateSearchResult
	for _, e := range estates {
		districtName := ""
		if e.District != nil {
			districtName = e.District.NameZhHant
		}
		results = append(results, models.EstateSearchResult{
			ID:                  e.ID,
			Name:                e.Name,
			NameEn:              e.NameEn,
			Address:             e.Address,
			DistrictName:        districtName,
			TotalUnits:          e.TotalUnits,
			AvgTransactionPrice: e.AvgTransactionPrice,
			CompletionYear:      e.CompletionYear,
		})
	}

	return results, total, nil
}

// SearchAgents 搜索代理人
func (r *SearchRepo) SearchAgents(ctx context.Context, req *models.SearchAgentsRequest) ([]models.AgentSearchResult, int64, error) {
	var agents []models.Agent
	var total int64

	searchPattern := "%" + req.Keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.Agent{}).
		Where("(agent_name LIKE ? OR agent_name_en LIKE ? OR license_no = ? OR specialization LIKE ?) AND status = ?",
			searchPattern, searchPattern, req.Keyword, searchPattern, "active")

	// 应用筛选条件
	if req.DistrictID != nil {
		// 查询服务该地区的代理人
		query = query.Where("id IN (?)",
			r.db.Model(&models.AgentServiceArea{}).
				Select("agent_id").
				Where("district_id = ?", *req.DistrictID),
		)
	}
	if req.Specialization != nil {
		query = query.Where("specialization LIKE ?", "%"+*req.Specialization+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).
		Order("rating DESC, properties_sold DESC").
		Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	var results []models.AgentSearchResult
	for _, a := range agents {
		results = append(results, models.AgentSearchResult{
			ID:             a.ID,
			AgentName:      a.AgentName,
			AgentNameEn:    a.AgentNameEn,
			LicenseNo:      a.LicenseNo,
			Phone:          a.Phone,
			Email:          a.Email,
			Specialization: a.Specialization,
			Rating:         a.Rating,
			PropertiesSold: a.PropertiesSold,
			ProfilePhoto:   a.ProfilePhoto,
		})
	}

	return results, total, nil
}

// GetSearchSuggestions 获取搜索建议
func (r *SearchRepo) GetSearchSuggestions(ctx context.Context, keyword string, limit int) ([]models.SearchSuggestion, error) {
	var suggestions []models.SearchSuggestion
	searchPattern := "%" + keyword + "%"

	// 房产名称建议
	var propertyTitles []string
	r.db.WithContext(ctx).Model(&models.Property{}).
		Select("DISTINCT title").
		Where("title LIKE ? AND status = ?", searchPattern, "available").
		Limit(limit / 4).
		Pluck("title", &propertyTitles)

	for _, title := range propertyTitles {
		suggestions = append(suggestions, models.SearchSuggestion{
			Keyword: title,
			Type:    "property",
			Label:   "房产",
		})
	}

	// 屋苑名称建议
	var estateNames []string
	r.db.WithContext(ctx).Model(&models.Estate{}).
		Select("DISTINCT name").
		Where("name LIKE ?", searchPattern).
		Limit(limit / 4).
		Pluck("name", &estateNames)

	for _, name := range estateNames {
		suggestions = append(suggestions, models.SearchSuggestion{
			Keyword: name,
			Type:    "estate",
			Label:   "屋苑",
		})
	}

	// 地区名称建议
	var districtNames []string
	r.db.WithContext(ctx).Model(&models.District{}).
		Select("DISTINCT name_zh_hant").
		Where("name_zh_hant LIKE ? OR name_en LIKE ?", searchPattern, searchPattern).
		Limit(limit / 4).
		Pluck("name_zh_hant", &districtNames)

	for _, name := range districtNames {
		suggestions = append(suggestions, models.SearchSuggestion{
			Keyword: name,
			Type:    "district",
			Label:   "地区",
		})
	}

	// 代理人名称建议
	var agentNames []string
	r.db.WithContext(ctx).Model(&models.Agent{}).
		Select("DISTINCT agent_name").
		Where("agent_name LIKE ? AND status = ?", searchPattern, "active").
		Limit(limit / 4).
		Pluck("agent_name", &agentNames)

	for _, name := range agentNames {
		suggestions = append(suggestions, models.SearchSuggestion{
			Keyword: name,
			Type:    "agent",
			Label:   "代理人",
		})
	}

	// 限制返回数量
	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

// GetSearchHistory 获取搜索历史
func (r *SearchRepo) GetSearchHistory(ctx context.Context, userID *uint, searchType string, page, pageSize int) ([]models.SearchHistory, int64, error) {
	var histories []models.SearchHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.SearchHistory{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if searchType != "" {
		query = query.Where("search_type = ?", searchType)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// SaveSearchHistory 保存搜索历史
func (r *SearchRepo) SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// DeleteSearchHistory 删除搜索历史
func (r *SearchRepo) DeleteSearchHistory(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.SearchHistory{}).Error
}
