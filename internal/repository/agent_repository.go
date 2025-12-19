package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// AgentRepository 代理人数据访问接口
type AgentRepository interface {
	// 代理人相关
	List(ctx context.Context, filter *request.ListAgentsRequest) ([]*model.Agent, int64, error)
	GetByID(ctx context.Context, id uint) (*model.Agent, error)
	GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*model.Property, int64, error)
	
	// 联系请求相关
	CreateContactRequest(ctx context.Context, contactReq *model.AgentContactRequest) error
	GetContactRequestByID(ctx context.Context, id uint) (*model.AgentContactRequest, error)
}

type agentRepository struct {
	db *gorm.DB
}

// NewAgentRepository 创建代理人仓库
func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &agentRepository{db: db}
}

// 代理人相关

func (r *agentRepository) List(ctx context.Context, filter *request.ListAgentsRequest) ([]*model.Agent, int64, error) {
	var agents []*model.Agent
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.Agent{})
	
	// 默认只查询活跃的代理人
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	} else {
		query = query.Where("status = ?", model.AgentStatusActive)
	}
	
	// 应用筛选条件
	if filter.AgencyID != nil {
		query = query.Where("agency_id = ?", *filter.AgencyID)
	}
	if filter.DistrictID != nil {
		// 通过服务区域关联表筛选
		query = query.Joins("JOIN agent_service_areas ON agent_service_areas.agent_id = agents.id").
			Where("agent_service_areas.district_id = ?", *filter.DistrictID)
	}
	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}
	if filter.Specialization != "" {
		query = query.Where("specialization LIKE ?", "%"+filter.Specialization+"%")
	}
	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("agent_name LIKE ? OR agent_name_en LIKE ? OR license_no LIKE ?", keyword, keyword, keyword)
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
	
	// 验证的代理人优先
	query = query.Order("is_verified DESC")
	query = query.Order(sortBy + " " + sortOrder)
	
	// 预加载关联
	query = query.Preload("Agency").Preload("ServiceAreas").Preload("ServiceAreas.District")
	
	if err := query.Find(&agents).Error; err != nil {
		return nil, 0, err
	}
	
	return agents, total, nil
}

func (r *agentRepository) GetByID(ctx context.Context, id uint) (*model.Agent, error) {
	var agent model.Agent
	
	err := r.db.WithContext(ctx).
		Preload("Agency").
		Preload("ServiceAreas").
		Preload("ServiceAreas.District").
		First(&agent, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &agent, nil
}

func (r *agentRepository) GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*model.Property, int64, error) {
	var properties []*model.Property
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.Property{}).
		Where("agent_id = ? AND status = ?", agentID, "active")
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 设置默认值
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	
	// 排序
	query = query.Order("created_at DESC")
	
	// 预加载关联
	query = query.Preload("District").Preload("Estate").Preload("Images")
	
	if err := query.Find(&properties).Error; err != nil {
		return nil, 0, err
	}
	
	return properties, total, nil
}

// 联系请求相关

func (r *agentRepository) CreateContactRequest(ctx context.Context, contactReq *model.AgentContactRequest) error {
	return r.db.WithContext(ctx).Create(contactReq).Error
}

func (r *agentRepository) GetContactRequestByID(ctx context.Context, id uint) (*model.AgentContactRequest, error) {
	var contactReq model.AgentContactRequest
	
	err := r.db.WithContext(ctx).
		Preload("Agent").
		Preload("Property").
		First(&contactReq, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &contactReq, nil
}
