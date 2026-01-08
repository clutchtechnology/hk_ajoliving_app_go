package databases

import (
	"context"
	"fmt"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type AgencyRepo struct {
	db *gorm.DB
}

func NewAgencyRepo(db *gorm.DB) *AgencyRepo {
	return &AgencyRepo{db: db}
}

// FindAll 查询所有代理公司（带筛选和分页）
func (r *AgencyRepo) FindAll(ctx context.Context, filter *models.ListAgenciesRequest) ([]models.AgencyDetail, int64, error) {
	var agencies []models.AgencyDetail
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AgencyDetail{})

	// 应用筛选条件
	if filter.DistrictID != nil {
		// 查询该地区有代理人服务的公司
		query = query.Where("user_id IN (?)",
			r.db.Model(&models.Agent{}).
				Select("agency_id").
				Joins("INNER JOIN agent_service_areas ON agents.id = agent_service_areas.agent_id").
				Where("agent_service_areas.district_id = ? AND agents.agency_id IS NOT NULL", *filter.DistrictID),
		)
	}

	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}

	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("company_name LIKE ? OR company_name_en LIKE ? OR description LIKE ?",
			keyword, keyword, keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	if err := query.Find(&agencies).Error; err != nil {
		return nil, 0, err
	}

	return agencies, total, nil
}

// FindByID 根据ID查询代理公司
func (r *AgencyRepo) FindByID(ctx context.Context, id uint) (*models.AgencyDetail, error) {
	var agency models.AgencyDetail
	if err := r.db.WithContext(ctx).Preload("User").First(&agency, id).Error; err != nil {
		return nil, err
	}
	return &agency, nil
}

// FindByUserID 根据UserID查询代理公司
func (r *AgencyRepo) FindByUserID(ctx context.Context, userID uint) (*models.AgencyDetail, error) {
	var agency models.AgencyDetail
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&agency).Error; err != nil {
		return nil, err
	}
	return &agency, nil
}

// GetAgencyProperties 查询代理公司的房源列表
func (r *AgencyRepo) GetAgencyProperties(ctx context.Context, userID uint, page, pageSize int) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("publisher_id = ? AND publisher_type = ?", userID, "agency")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("District").
		Order("created_at DESC").
		Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// GetAgencyAgents 查询代理公司的代理人列表
func (r *AgencyRepo) GetAgencyAgents(ctx context.Context, userID uint) ([]models.Agent, error) {
	var agents []models.Agent
	if err := r.db.WithContext(ctx).
		Where("agency_id = ? AND status = ?", userID, "active").
		Order("created_at DESC").
		Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

// CreateContact 创建代理公司联系记录
func (r *AgencyRepo) CreateContact(ctx context.Context, contact *models.AgencyContact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

// Search 搜索代理公司
func (r *AgencyRepo) Search(ctx context.Context, keyword string, page, pageSize int) ([]models.AgencyDetail, int64, error) {
	var agencies []models.AgencyDetail
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.AgencyDetail{}).
		Where("company_name LIKE ? OR company_name_en LIKE ? OR license_no = ? OR description LIKE ?",
			searchPattern, searchPattern, keyword, searchPattern)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("is_verified DESC, rating DESC, agent_count DESC").
		Find(&agencies).Error; err != nil {
		return nil, 0, err
	}

	return agencies, total, nil
}

// IncrementAgentCount 增加代理人数量
func (r *AgencyRepo) IncrementAgentCount(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&models.AgencyDetail{}).
		Where("user_id = ?", userID).
		UpdateColumn("agent_count", gorm.Expr("agent_count + ?", 1)).Error
}

// DecrementAgentCount 减少代理人数量
func (r *AgencyRepo) DecrementAgentCount(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&models.AgencyDetail{}).
		Where("user_id = ?", userID).
		Where("agent_count > ?", 0).
		UpdateColumn("agent_count", gorm.Expr("agent_count - ?", 1)).Error
}
