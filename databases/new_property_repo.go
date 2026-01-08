package databases

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// NewDevelopmentRepo 新盘仓储
type NewDevelopmentRepo struct {
	db *gorm.DB
}

// NewNewDevelopmentRepo 创建新盘仓储
func NewNewDevelopmentRepo(db *gorm.DB) *NewDevelopmentRepo {
	return &NewDevelopmentRepo{db: db}
}

// FindAll 查找所有新盘（支持筛选和分页）
func (r *NewDevelopmentRepo) FindAll(ctx context.Context, filter *models.ListNewPropertiesRequest) ([]models.NewProperty, int64, error) {
	var newProperties []models.NewProperty
	var total int64

	query := r.db.WithContext(ctx).Model(&models.NewProperty{})

	// 应用筛选条件
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Developer != nil {
		query = query.Where("developer LIKE ?", "%"+*filter.Developer+"%")
	}
	if filter.PrimarySchoolNet != nil {
		query = query.Where("primary_school_net = ?", *filter.PrimarySchoolNet)
	}
	if filter.SecondarySchoolNet != nil {
		query = query.Where("secondary_school_net = ?", *filter.SecondarySchoolNet)
	}
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序：精选优先，然后按创建时间降序
	query = query.Order("is_featured DESC, sort_order ASC, created_at DESC")

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// 预加载关联
	query = query.Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1) // 列表只加载第一张图
		})

	if err := query.Find(&newProperties).Error; err != nil {
		return nil, 0, err
	}

	return newProperties, total, nil
}

// FindByID 根据ID查找新盘
func (r *NewDevelopmentRepo) FindByID(ctx context.Context, id uint) (*models.NewProperty, error) {
	var newProperty models.NewProperty
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Layouts", func(db *gorm.DB) *gorm.DB {
			return db.Order("min_price ASC")
		}).
		First(&newProperty, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("new property not found")
		}
		return nil, err
	}
	return &newProperty, nil
}

// IncrementViewCount 增加浏览次数
func (r *NewDevelopmentRepo) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.NewProperty{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// FindLayouts 查找新盘的所有户型
func (r *NewDevelopmentRepo) FindLayouts(ctx context.Context, newPropertyID uint) ([]models.NewPropertyLayout, error) {
	var layouts []models.NewPropertyLayout
	err := r.db.WithContext(ctx).
		Where("new_property_id = ?", newPropertyID).
		Order("min_price ASC").
		Find(&layouts).Error
	return layouts, err
}
