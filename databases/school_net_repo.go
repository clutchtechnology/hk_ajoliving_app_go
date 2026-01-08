package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type SchoolNetRepo struct {
	db *gorm.DB
}

func NewSchoolNetRepo(db *gorm.DB) *SchoolNetRepo {
	return &SchoolNetRepo{db: db}
}

// FindAll 查询所有校网（分页+筛选）
func (r *SchoolNetRepo) FindAll(ctx context.Context, req *models.ListSchoolNetsRequest) ([]*models.SchoolNet, int64, error) {
	var schoolNets []*models.SchoolNet
	var total int64

	query := r.db.WithContext(ctx).Model(&models.SchoolNet{})

	// 应用筛选条件
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ? OR code LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Preload("District").
		Offset(offset).
		Limit(req.PageSize).
		Order("code ASC").
		Find(&schoolNets).Error

	return schoolNets, total, err
}

// FindByID 根据ID查询校网
func (r *SchoolNetRepo) FindByID(ctx context.Context, id uint) (*models.SchoolNet, error) {
	var schoolNet models.SchoolNet
	err := r.db.WithContext(ctx).
		Preload("District").
		First(&schoolNet, id).Error
	return &schoolNet, err
}

// GetSchoolsInNet 获取校网内的学校
func (r *SchoolNetRepo) GetSchoolsInNet(ctx context.Context, netID uint, page, pageSize int) ([]*models.School, int64, error) {
	var schools []*models.School
	var total int64

	query := r.db.WithContext(ctx).Model(&models.School{}).Where("school_net_id = ?", netID)

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
		Order("name_zh_hant ASC").
		Find(&schools).Error

	return schools, total, err
}

// GetPropertiesInNet 获取校网内的房源
func (r *SchoolNetRepo) GetPropertiesInNet(ctx context.Context, netID uint, netType string, page, pageSize int) ([]*models.Property, int64, error) {
	var properties []*models.Property
	var total int64

	// 先获取校网编码
	var schoolNet models.SchoolNet
	if err := r.db.WithContext(ctx).First(&schoolNet, netID).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&models.Property{})

	// 根据校网类型筛选
	if netType == "primary" {
		query = query.Where("primary_school_net = ?", schoolNet.Code)
	} else if netType == "secondary" {
		query = query.Where("secondary_school_net = ?", schoolNet.Code)
	}

	// 只查询可用状态
	query = query.Where("status = ?", "available")

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

// GetEstatesInNet 获取校网内的屋苑
func (r *SchoolNetRepo) GetEstatesInNet(ctx context.Context, netID uint, netType string, page, pageSize int) ([]*models.Estate, int64, error) {
	var estates []*models.Estate
	var total int64

	// 先获取校网编码
	var schoolNet models.SchoolNet
	if err := r.db.WithContext(ctx).First(&schoolNet, netID).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&models.Estate{})

	// 根据校网类型筛选
	if netType == "primary" {
		query = query.Where("primary_school_net = ?", schoolNet.Code)
	} else if netType == "secondary" {
		query = query.Where("secondary_school_net = ?", schoolNet.Code)
	}

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
		Order("name ASC").
		Find(&estates).Error

	return estates, total, err
}

// Search 搜索校网
func (r *SchoolNetRepo) Search(ctx context.Context, keyword string, page, pageSize int) ([]*models.SchoolNet, int64, error) {
	var schoolNets []*models.SchoolNet
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.SchoolNet{}).
		Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ? OR code LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)

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
		Order("code ASC").
		Find(&schoolNets).Error

	return schoolNets, total, err
}
