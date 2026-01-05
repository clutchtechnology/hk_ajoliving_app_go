package databases

import (
	"context"
	"fmt"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// FacilityRepository 设施数据仓库接口
type FacilityRepository interface {
	List(ctx context.Context, req *models.ListFacilitiesRequest) ([]*models.Facility, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Facility, error)
	Create(ctx context.Context, facility *models.Facility) error
	Update(ctx context.Context, facility *models.Facility) error
	Delete(ctx context.Context, id uint) error
	ExistsByID(ctx context.Context, id uint) (bool, error)
}

type facilityRepository struct {
	db *gorm.DB
}

// NewFacilityRepository 创建设施仓库实例
func NewFacilityRepository(db *gorm.DB) FacilityRepository {
	return &facilityRepository{db: db}
}

// List 获取设施列表
func (r *facilityRepository) List(ctx context.Context, req *models.ListFacilitiesRequest) ([]*models.Facility, int64, error) {
	var facilities []*models.Facility
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Facility{})

	// 应用筛选条件
	if req.Category != nil {
		query = query.Where("category = ?", *req.Category)
	}

	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where(
			"name_zh_hant ILIKE ? OR name_zh_hans ILIKE ? OR name_en ILIKE ?",
			keyword, keyword, keyword,
		)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (req.Page - 1) * req.PageSize
	orderBy := fmt.Sprintf("%s %s", req.SortBy, req.SortOrder)
	
	if err := query.
		Offset(offset).
		Limit(req.PageSize).
		Order(orderBy).
		Find(&facilities).Error; err != nil {
		return nil, 0, err
	}

	return facilities, total, nil
}

// GetByID 根据ID获取设施
func (r *facilityRepository) GetByID(ctx context.Context, id uint) (*models.Facility, error) {
	var facility models.Facility
	if err := r.db.WithContext(ctx).First(&facility, id).Error; err != nil {
		return nil, err
	}
	return &facility, nil
}

// Create 创建设施
func (r *facilityRepository) Create(ctx context.Context, facility *models.Facility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

// Update 更新设施
func (r *facilityRepository) Update(ctx context.Context, facility *models.Facility) error {
	return r.db.WithContext(ctx).Save(facility).Error
}

// Delete 删除设施
func (r *facilityRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Facility{}, id).Error
}

// ExistsByID 检查设施是否存在
func (r *facilityRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Facility{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
