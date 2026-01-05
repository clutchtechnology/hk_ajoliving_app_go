package databases

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// SearchRepository 搜索数据仓库接口
type SearchRepository interface {
	// 搜索相关
	SearchProperties(ctx context.Context, keyword string, filters *models.SearchPropertiesRequest) ([]*models.Property, int64, error)
	SearchEstates(ctx context.Context, keyword string, districtID *uint, page, pageSize int) ([]*models.Estate, int64, error)
	SearchAgents(ctx context.Context, keyword string, districtID *uint, page, pageSize int) ([]*models.Agent, int64, error)
	SearchNews(ctx context.Context, keyword string, limit int) ([]*models.News, error)
	
	// 搜索建议
	GetPropertySuggestions(ctx context.Context, keyword string, limit int) ([]string, error)
	GetEstateSuggestions(ctx context.Context, keyword string, limit int) ([]string, error)
	GetAgentSuggestions(ctx context.Context, keyword string, limit int) ([]string, error)
	
	// 搜索历史
	SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error
	GetSearchHistory(ctx context.Context, userID *uint, searchType *string, page, pageSize int) ([]*models.SearchHistory, int64, error)
	DeleteSearchHistory(ctx context.Context, userID uint, historyID uint) error
	ClearSearchHistory(ctx context.Context, userID uint) error
}

type searchRepository struct {
	db *gorm.DB
}

// NewSearchRepository 创建搜索仓库实例
func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

// SearchProperties 搜索房产
func (r *searchRepository) SearchProperties(ctx context.Context, keyword string, filters *models.SearchPropertiesRequest) ([]*models.Property, int64, error) {
	var properties []*models.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Property{})

	// 关键词搜索
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where(
			"title ILIKE ? OR description ILIKE ? OR address ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// 应用筛选条件
	if filters.ListingType != nil {
		query = query.Where("listing_type = ?", *filters.ListingType)
	}
	if filters.DistrictID != nil {
		query = query.Where("district_id = ?", *filters.DistrictID)
	}
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}
	if filters.Bedrooms != nil {
		query = query.Where("bedrooms = ?", *filters.Bedrooms)
	}
	if filters.PropertyType != nil {
		query = query.Where("property_type = ?", *filters.PropertyType)
	}

	// 只查询活跃状态
	query = query.Where("status = ?", "active")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (filters.Page - 1) * filters.PageSize
	if err := query.
		Preload("District").
		Preload("Estate").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1)
		}).
		Offset(offset).
		Limit(filters.PageSize).
		Order("created_at DESC").
		Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// SearchEstates 搜索屋苑
func (r *searchRepository) SearchEstates(ctx context.Context, keyword string, districtID *uint, page, pageSize int) ([]*models.Estate, int64, error) {
	var estates []*models.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Estate{})

	// 关键词搜索
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where(
			"name ILIKE ? OR name_en ILIKE ? OR address ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// 地区筛选
	if districtID != nil {
		query = query.Where("district_id = ?", *districtID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (page - 1) * pageSize
	if err := query.
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("image_type = ?", "exterior").Order("sort_order ASC").Limit(1)
		}).
		Offset(offset).
		Limit(pageSize).
		Order("view_count DESC, name ASC").
		Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

// SearchAgents 搜索代理人
func (r *searchRepository) SearchAgents(ctx context.Context, keyword string, districtID *uint, page, pageSize int) ([]*models.Agent, int64, error) {
	var agents []*models.Agent
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Agent{})

	// 关键词搜索
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where(
			"name ILIKE ? OR license_no ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// 地区筛选（通过服务区域）
	if districtID != nil {
		query = query.Joins("JOIN agent_service_areas ON agent_service_areas.agent_id = agents.id").
			Where("agent_service_areas.district_id = ?", *districtID).
			Distinct()
	}

	// 只查询活跃代理人
	query = query.Where("status = ?", "active")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (page - 1) * pageSize
	if err := query.
		Preload("Agency").
		Offset(offset).
		Limit(pageSize).
		Order("rating DESC, experience_years DESC").
		Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// SearchNews 搜索新闻
func (r *searchRepository) SearchNews(ctx context.Context, keyword string, limit int) ([]*models.News, error) {
	var news []*models.News

	query := r.db.WithContext(ctx).Model(&models.News{})

	// 关键词搜索
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where(
			"title ILIKE ? OR summary ILIKE ? OR content ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// 只查询已发布的新闻
	query = query.Where("status = ?", "published")

	if err := query.
		Preload("Category").
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

// GetPropertySuggestions 获取房产搜索建议
func (r *searchRepository) GetPropertySuggestions(ctx context.Context, keyword string, limit int) ([]string, error) {
	var suggestions []string

	searchPattern := "%" + keyword + "%"
	
	err := r.db.WithContext(ctx).
		Model(&models.Property{}).
		Where("status = ?", "active").
		Where("title ILIKE ? OR address ILIKE ?", searchPattern, searchPattern).
		Distinct("COALESCE(title, address)").
		Limit(limit).
		Pluck("COALESCE(title, address)", &suggestions).Error

	return suggestions, err
}

// GetEstateSuggestions 获取屋苑搜索建议
func (r *searchRepository) GetEstateSuggestions(ctx context.Context, keyword string, limit int) ([]string, error) {
	var suggestions []string

	searchPattern := "%" + keyword + "%"
	
	err := r.db.WithContext(ctx).
		Model(&models.Estate{}).
		Where("name ILIKE ? OR name_en ILIKE ?", searchPattern, searchPattern).
		Distinct("name").
		Limit(limit).
		Pluck("name", &suggestions).Error

	return suggestions, err
}

// GetAgentSuggestions 获取代理人搜索建议
func (r *searchRepository) GetAgentSuggestions(ctx context.Context, keyword string, limit int) ([]string, error) {
	var suggestions []string

	searchPattern := "%" + keyword + "%"
	
	err := r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("status = ?", "active").
		Where("name ILIKE ?", searchPattern).
		Distinct("name").
		Limit(limit).
		Pluck("name", &suggestions).Error

	return suggestions, err
}

// SaveSearchHistory 保存搜索历史
func (r *searchRepository) SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetSearchHistory 获取搜索历史
func (r *searchRepository) GetSearchHistory(ctx context.Context, userID *uint, searchType *string, page, pageSize int) ([]*models.SearchHistory, int64, error) {
	var histories []*models.SearchHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.SearchHistory{})

	// 用户筛选
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	} else {
		query = query.Where("user_id IS NULL")
	}

	// 类型筛选
	if searchType != nil {
		query = query.Where("search_type = ?", *searchType)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// DeleteSearchHistory 删除搜索历史
func (r *searchRepository) DeleteSearchHistory(ctx context.Context, userID uint, historyID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", historyID, userID).
		Delete(&models.SearchHistory{}).Error
}

// ClearSearchHistory 清空搜索历史
func (r *searchRepository) ClearSearchHistory(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.SearchHistory{}).Error
}
