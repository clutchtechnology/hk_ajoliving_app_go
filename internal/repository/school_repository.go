package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// SchoolRepository 校网和学校数据访问接口
type SchoolRepository interface {
	// 校网相关
	ListSchoolNets(ctx context.Context, filter *request.ListSchoolNetsRequest) ([]*model.SchoolNet, int64, error)
	GetSchoolNetByID(ctx context.Context, id uint) (*model.SchoolNet, error)
	GetSchoolNetByCode(ctx context.Context, code string) (*model.SchoolNet, error)
	SearchSchoolNets(ctx context.Context, keyword string, limit int, offset int) ([]*model.SchoolNet, int64, error)
	CountPropertiesInSchoolNet(ctx context.Context, schoolNetID uint) (int64, error)
	CountEstatesInSchoolNet(ctx context.Context, schoolNetID uint) (int64, error)
	
	// 学校相关
	ListSchools(ctx context.Context, filter *request.ListSchoolsRequest) ([]*model.School, int64, error)
	GetSchoolByID(ctx context.Context, id uint) (*model.School, error)
	GetSchoolsBySchoolNetID(ctx context.Context, schoolNetID uint) ([]*model.School, error)
	SearchSchools(ctx context.Context, keyword string, limit int, offset int) ([]*model.School, int64, error)
}

type schoolRepository struct {
	db *gorm.DB
}

// NewSchoolRepository 创建校网和学校仓库
func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolRepository{db: db}
}

// 校网相关

func (r *schoolRepository) ListSchoolNets(ctx context.Context, filter *request.ListSchoolNetsRequest) ([]*model.SchoolNet, int64, error) {
	var schoolNets []*model.SchoolNet
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.SchoolNet{}).Where("is_active = ?", true)
	
	// 应用筛选条件
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.Level != nil && *filter.Level != "" {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ? OR net_code LIKE ?", 
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
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "net_code"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "asc"
	}
	
	// 分页和排序
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	query = query.Order(sortBy + " " + sortOrder)
	
	// 预加载关联
	query = query.Preload("District")
	
	if err := query.Find(&schoolNets).Error; err != nil {
		return nil, 0, err
	}
	
	return schoolNets, total, nil
}

func (r *schoolRepository) GetSchoolNetByID(ctx context.Context, id uint) (*model.SchoolNet, error) {
	var schoolNet model.SchoolNet
	
	err := r.db.WithContext(ctx).
		Preload("District").
		First(&schoolNet, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &schoolNet, nil
}

func (r *schoolRepository) GetSchoolNetByCode(ctx context.Context, code string) (*model.SchoolNet, error) {
	var schoolNet model.SchoolNet
	
	err := r.db.WithContext(ctx).
		Preload("District").
		Where("net_code = ?", code).
		First(&schoolNet).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &schoolNet, nil
}

func (r *schoolRepository) SearchSchoolNets(ctx context.Context, keyword string, limit int, offset int) ([]*model.SchoolNet, int64, error) {
	var schoolNets []*model.SchoolNet
	var total int64
	
	searchPattern := "%" + keyword + "%"
	query := r.db.WithContext(ctx).Model(&model.SchoolNet{}).
		Where("is_active = ?", true).
		Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ? OR net_code LIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页
	query = query.Offset(offset).Limit(limit).Preload("District")
	
	if err := query.Find(&schoolNets).Error; err != nil {
		return nil, 0, err
	}
	
	return schoolNets, total, nil
}

func (r *schoolRepository) CountPropertiesInSchoolNet(ctx context.Context, schoolNetID uint) (int64, error) {
	var count int64
	
	// 通过地区关联统计（假设 properties 表有 district_id，并且 school_nets 也有 district_id）
	// 这里需要根据实际的关联关系调整
	err := r.db.WithContext(ctx).
		Model(&model.Property{}).
		Joins("JOIN school_nets ON properties.district_id = school_nets.district_id").
		Where("school_nets.id = ?", schoolNetID).
		Count(&count).Error
		
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *schoolRepository) CountEstatesInSchoolNet(ctx context.Context, schoolNetID uint) (int64, error) {
	var count int64
	
	// 通过地区关联统计
	err := r.db.WithContext(ctx).
		Model(&model.Estate{}).
		Joins("JOIN school_nets ON estates.district_id = school_nets.district_id").
		Where("school_nets.id = ?", schoolNetID).
		Count(&count).Error
		
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// 学校相关

func (r *schoolRepository) ListSchools(ctx context.Context, filter *request.ListSchoolsRequest) ([]*model.School, int64, error) {
	var schools []*model.School
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.School{}).Where("is_active = ?", true)
	
	// 应用筛选条件
	if filter.SchoolNetID != nil {
		query = query.Where("school_net_id = ?", *filter.SchoolNetID)
	}
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.Level != nil && *filter.Level != "" {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.Gender != nil && *filter.Gender != "" {
		query = query.Where("gender = ?", *filter.Gender)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ? OR school_code LIKE ?", 
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
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "name_zh_hant"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "asc"
	}
	
	// 分页和排序
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	query = query.Order(sortBy + " " + sortOrder)
	
	// 预加载关联
	query = query.Preload("District").Preload("SchoolNet")
	
	if err := query.Find(&schools).Error; err != nil {
		return nil, 0, err
	}
	
	return schools, total, nil
}

func (r *schoolRepository) GetSchoolByID(ctx context.Context, id uint) (*model.School, error) {
	var school model.School
	
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("SchoolNet").
		Preload("SchoolNet.District").
		First(&school, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &school, nil
}

func (r *schoolRepository) GetSchoolsBySchoolNetID(ctx context.Context, schoolNetID uint) ([]*model.School, error) {
	var schools []*model.School
	
	err := r.db.WithContext(ctx).
		Preload("District").
		Where("school_net_id = ? AND is_active = ?", schoolNetID, true).
		Order("name_zh_hant ASC").
		Find(&schools).Error
		
	if err != nil {
		return nil, err
	}
	
	return schools, nil
}

func (r *schoolRepository) SearchSchools(ctx context.Context, keyword string, limit int, offset int) ([]*model.School, int64, error) {
	var schools []*model.School
	var total int64
	
	searchPattern := "%" + keyword + "%"
	query := r.db.WithContext(ctx).Model(&model.School{}).
		Where("is_active = ?", true).
		Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ? OR school_code LIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页
	query = query.Offset(offset).Limit(limit).
		Preload("District").
		Preload("SchoolNet")
	
	if err := query.Find(&schools).Error; err != nil {
		return nil, 0, err
	}
	
	return schools, total, nil
}
