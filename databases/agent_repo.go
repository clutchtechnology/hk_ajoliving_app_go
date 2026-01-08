package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type AgentRepo struct {
	db *gorm.DB
}

func NewAgentRepo(db *gorm.DB) *AgentRepo {
	return &AgentRepo{db: db}
}

// FindAll 查询所有代理人（分页+筛选）
func (r *AgentRepo) FindAll(ctx context.Context, req *models.ListAgentsRequest) ([]*models.Agent, int64, error) {
	var agents []*models.Agent
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Agent{})

	// 应用筛选条件
	if req.LicenseType != nil {
		query = query.Where("license_type = ?", *req.LicenseType)
	}
	if req.AgencyID != nil {
		query = query.Where("agency_id = ?", *req.AgencyID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.IsVerified != nil {
		query = query.Where("is_verified = ?", *req.IsVerified)
	}
	if req.Specialization != nil {
		query = query.Where("specialization LIKE ?", "%"+*req.Specialization+"%")
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("agent_name LIKE ? OR agent_name_en LIKE ? OR license_no LIKE ?",
			keyword, keyword, keyword)
	}

	// 如果指定了地区，通过服务区域关联查询
	if req.DistrictID != nil {
		query = query.Joins("INNER JOIN agent_service_areas ON agent_service_areas.agent_id = agents.id").
			Where("agent_service_areas.district_id = ?", *req.DistrictID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Offset(offset).
		Limit(req.PageSize).
		Order("rating DESC, properties_sold DESC").
		Find(&agents).Error

	return agents, total, err
}

// FindByID 根据ID查询代理人
func (r *AgentRepo) FindByID(ctx context.Context, id uint) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&agent, id).Error
	return &agent, err
}

// GetAgentProperties 获取代理人的房源列表
func (r *AgentRepo) GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*models.Property, int64, error) {
	var properties []*models.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("agent_id = ? AND status = ?", agentID, "available")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("District").
		Offset(offset).
		Limit(pageSize).
		Order("published_at DESC").
		Find(&properties).Error

	return properties, total, err
}

// GetServiceAreas 获取代理人服务区域
func (r *AgentRepo) GetServiceAreas(ctx context.Context, agentID uint) ([]*models.AgentServiceArea, error) {
	var serviceAreas []*models.AgentServiceArea
	err := r.db.WithContext(ctx).
		Preload("District").
		Where("agent_id = ?", agentID).
		Find(&serviceAreas).Error
	return serviceAreas, err
}

// CreateContact 创建联系代理人记录
func (r *AgentRepo) CreateContact(ctx context.Context, contact *models.AgentContact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

// IncrementPropertySold 增加已售物业数量
func (r *AgentRepo) IncrementPropertySold(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("id = ?", id).
		UpdateColumn("properties_sold", gorm.Expr("properties_sold + ?", 1)).
		Error
}

// IncrementPropertyRented 增加已租物业数量
func (r *AgentRepo) IncrementPropertyRented(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("id = ?", id).
		UpdateColumn("properties_rented", gorm.Expr("properties_rented + ?", 1)).
		Error
}
