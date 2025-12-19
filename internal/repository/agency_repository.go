package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// AgencyRepository 代理公司数据访问接口
type AgencyRepository interface {
	// 代理公司相关
	List(ctx context.Context, filter *request.ListAgenciesRequest) ([]*model.AgencyDetail, int64, error)
	GetByID(ctx context.Context, id uint) (*model.AgencyDetail, error)
	GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) ([]*model.Property, int64, error)
	Search(ctx context.Context, filter *request.SearchAgenciesRequest) ([]*model.AgencyDetail, int64, error)
	
	// 代理人相关
	GetAgencyAgents(ctx context.Context, agencyID uint) ([]*model.Agent, error)
	GetTopAgents(ctx context.Context, agencyID uint, limit int) ([]*model.Agent, error)
	
	// 统计相关
	GetPropertyCount(ctx context.Context, agencyID uint) (int, error)
}

type agencyRepository struct {
	db *gorm.DB
}

// NewAgencyRepository 创建代理公司仓库
func NewAgencyRepository(db *gorm.DB) AgencyRepository {
	return &agencyRepository{db: db}
}

// List 获取代理公司列表
func (r *agencyRepository) List(ctx context.Context, filter *request.ListAgenciesRequest) ([]*model.AgencyDetail, int64, error) {
	var agencies []*model.AgencyDetail
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.AgencyDetail{})
	
	// 应用筛选条件
	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}
	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("company_name LIKE ? OR company_name_en LIKE ? OR description LIKE ?", keyword, keyword, keyword)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 设置默认值
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "rating"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	
	// 分页和排序
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	
	// 验证的代理公司优先
	query = query.Order("is_verified DESC")
	query = query.Order(sortBy + " " + sortOrder)
	
	// 预加载用户关联
	query = query.Preload("User")
	
	if err := query.Find(&agencies).Error; err != nil {
		return nil, 0, err
	}
	
	return agencies, total, nil
}

// GetByID 根据ID获取代理公司详情
func (r *agencyRepository) GetByID(ctx context.Context, id uint) (*model.AgencyDetail, error) {
	var agency model.AgencyDetail
	
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&agency, id).Error
	
	if err != nil {
		return nil, err
	}
	
	return &agency, nil
}

// GetAgencyProperties 获取代理公司的房源列表
func (r *agencyRepository) GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) ([]*model.Property, int64, error) {
	var properties []*model.Property
	var total int64
	
	// 先验证代理公司是否存在
	var agency model.AgencyDetail
	if err := r.db.WithContext(ctx).First(&agency, agencyID).Error; err != nil {
		return nil, 0, err
	}
	
	// 查询该代理公司的所有代理人
	var agentIDs []uint
	if err := r.db.WithContext(ctx).
		Model(&model.Agent{}).
		Where("agency_id = ?", agencyID).
		Pluck("id", &agentIDs).Error; err != nil {
		return nil, 0, err
	}
	
	if len(agentIDs) == 0 {
		return []*model.Property{}, 0, nil
	}
	
	// 查询这些代理人的房源
	query := r.db.WithContext(ctx).Model(&model.Property{}).
		Where("agent_id IN ?", agentIDs).
		Where("status = ?", "active")
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	
	// 排序：最新发布的在前
	query = query.Order("created_at DESC")
	
	// 预加载关联数据
	query = query.Preload("District").
		Preload("Estate").
		Preload("Images").
		Preload("Agent")
	
	if err := query.Find(&properties).Error; err != nil {
		return nil, 0, err
	}
	
	return properties, total, nil
}

// Search 搜索代理公司
func (r *agencyRepository) Search(ctx context.Context, filter *request.SearchAgenciesRequest) ([]*model.AgencyDetail, int64, error) {
	var agencies []*model.AgencyDetail
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.AgencyDetail{})
	
	// 搜索关键词
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("company_name LIKE ? OR company_name_en LIKE ? OR license_no LIKE ? OR description LIKE ?", 
			keyword, keyword, keyword, keyword)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 设置默认值
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	// 分页和排序
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	
	// 验证的优先，然后按评分排序
	query = query.Order("is_verified DESC, rating DESC NULLS LAST")
	
	// 预加载用户关联
	query = query.Preload("User")
	
	if err := query.Find(&agencies).Error; err != nil {
		return nil, 0, err
	}
	
	return agencies, total, nil
}

// GetAgencyAgents 获取代理公司的所有代理人
func (r *agencyRepository) GetAgencyAgents(ctx context.Context, agencyID uint) ([]*model.Agent, error) {
	var agents []*model.Agent
	
	err := r.db.WithContext(ctx).
		Where("agency_id = ?", agencyID).
		Where("status = ?", model.AgentStatusActive).
		Order("is_verified DESC, rating DESC NULLS LAST").
		Find(&agents).Error
	
	if err != nil {
		return nil, err
	}
	
	return agents, nil
}

// GetTopAgents 获取代理公司的优秀代理人
func (r *agencyRepository) GetTopAgents(ctx context.Context, agencyID uint, limit int) ([]*model.Agent, error) {
	var agents []*model.Agent
	
	if limit <= 0 {
		limit = 5
	}
	
	err := r.db.WithContext(ctx).
		Where("agency_id = ?", agencyID).
		Where("status = ?", model.AgentStatusActive).
		Where("rating IS NOT NULL").
		Order("rating DESC, review_count DESC").
		Limit(limit).
		Find(&agents).Error
	
	if err != nil {
		return nil, err
	}
	
	return agents, nil
}

// GetPropertyCount 获取代理公司的房源总数
func (r *agencyRepository) GetPropertyCount(ctx context.Context, agencyID uint) (int, error) {
	var count int64
	
	// 查询该代理公司的所有代理人
	var agentIDs []uint
	if err := r.db.WithContext(ctx).
		Model(&model.Agent{}).
		Where("agency_id = ?", agencyID).
		Pluck("id", &agentIDs).Error; err != nil {
		return 0, err
	}
	
	if len(agentIDs) == 0 {
		return 0, nil
	}
	
	// 统计这些代理人的活跃房源数
	err := r.db.WithContext(ctx).
		Model(&model.Property{}).
		Where("agent_id IN ?", agentIDs).
		Where("status = ?", "active").
		Count(&count).Error
	
	if err != nil {
		return 0, err
	}
	
	return int(count), nil
}
