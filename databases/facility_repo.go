package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type FacilityRepo struct {
	db *gorm.DB
}

func NewFacilityRepo(db *gorm.DB) *FacilityRepo {
	return &FacilityRepo{db: db}
}

// FindAll 查询所有设施
func (r *FacilityRepo) FindAll(ctx context.Context, category string) ([]models.Facility, error) {
	var facilities []models.Facility
	query := r.db.WithContext(ctx).Model(&models.Facility{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("category ASC, sort_order ASC, id ASC").Find(&facilities).Error; err != nil {
		return nil, err
	}

	return facilities, nil
}

// FindByID 根据ID查询设施
func (r *FacilityRepo) FindByID(ctx context.Context, id uint) (*models.Facility, error) {
	var facility models.Facility
	if err := r.db.WithContext(ctx).First(&facility, id).Error; err != nil {
		return nil, err
	}
	return &facility, nil
}

// Create 创建设施
func (r *FacilityRepo) Create(ctx context.Context, facility *models.Facility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

// Update 更新设施
func (r *FacilityRepo) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Facility{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除设施
func (r *FacilityRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Facility{}, id).Error
}

// CheckNameExists 检查设施名称是否存在
func (r *FacilityRepo) CheckNameExists(ctx context.Context, nameZhHant string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Facility{}).
		Where("name_zh_hant = ?", nameZhHant)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
