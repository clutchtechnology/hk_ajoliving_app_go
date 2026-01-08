package databases

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// ServicedApartmentRepo 服务式住宅仓储
type ServicedApartmentRepo struct {
	db *gorm.DB
}

// NewServicedApartmentRepo 创建服务式住宅仓储
func NewServicedApartmentRepo(db *gorm.DB) *ServicedApartmentRepo {
	return &ServicedApartmentRepo{db: db}
}

// Create 创建服务式住宅
func (r *ServicedApartmentRepo) Create(ctx context.Context, sa *models.ServicedApartment) error {
	return r.db.WithContext(ctx).Create(sa).Error
}

// FindAll 查找所有服务式住宅（支持筛选和分页）
func (r *ServicedApartmentRepo) FindAll(ctx context.Context, filter *models.ListServicedApartmentsRequest) ([]models.ServicedApartment, int64, error) {
	var apartments []models.ServicedApartment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ServicedApartment{})

	// 应用筛选条件
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	} else {
		// 默认只显示营业中
		query = query.Where("status = ?", "active")
	}
	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序：精选优先，然后按评分降序
	query = query.Order("is_featured DESC, rating DESC, created_at DESC")

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// 预加载关联
	query = query.Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("serviced_apartment_id IS NOT NULL").Order("sort_order ASC").Limit(1)
		}).
		Preload("Units", func(db *gorm.DB) *gorm.DB {
			return db.Order("monthly_price ASC").Limit(1) // 只加载最低价房型
		})

	if err := query.Find(&apartments).Error; err != nil {
		return nil, 0, err
	}

	return apartments, total, nil
}

// FindByID 根据ID查找服务式住宅
func (r *ServicedApartmentRepo) FindByID(ctx context.Context, id uint) (*models.ServicedApartment, error) {
	var apartment models.ServicedApartment
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Units", func(db *gorm.DB) *gorm.DB {
			return db.Order("monthly_price ASC")
		}).
		First(&apartment, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("serviced apartment not found")
		}
		return nil, err
	}
	return &apartment, nil
}

// Update 更新服务式住宅
func (r *ServicedApartmentRepo) Update(ctx context.Context, sa *models.ServicedApartment) error {
	return r.db.WithContext(ctx).Save(sa).Error
}

// Delete 删除服务式住宅（软删除）
func (r *ServicedApartmentRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ServicedApartment{}, id).Error
}

// IncrementViewCount 增加浏览次数
func (r *ServicedApartmentRepo) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.ServicedApartment{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// FindUnits 查找服务式住宅的所有房型
func (r *ServicedApartmentRepo) FindUnits(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentUnit, error) {
	var units []models.ServicedApartmentUnit
	err := r.db.WithContext(ctx).
		Where("serviced_apartment_id = ?", apartmentID).
		Order("monthly_price ASC").
		Find(&units).Error
	return units, err
}

// FindImages 查找服务式住宅的所有图片
func (r *ServicedApartmentRepo) FindImages(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentImage, error) {
	var images []models.ServicedApartmentImage
	err := r.db.WithContext(ctx).
		Where("serviced_apartment_id = ?", apartmentID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

// FindByCompanyID 根据公司ID查找服务式住宅
func (r *ServicedApartmentRepo) FindByCompanyID(ctx context.Context, companyID uint) ([]models.ServicedApartment, error) {
	var apartments []models.ServicedApartment
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("serviced_apartment_id IS NOT NULL").Order("sort_order ASC").Limit(1)
		}).
		Order("created_at DESC").
		Find(&apartments).Error
	return apartments, err
}
