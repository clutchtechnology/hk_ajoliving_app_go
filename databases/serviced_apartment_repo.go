package databases

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// ServicedApartmentRepository 服务式公寓数据仓库接口
type ServicedApartmentRepository interface {
	Create(ctx context.Context, apartment *models.ServicedApartment) error
	GetByID(ctx context.Context, id uint) (*models.ServicedApartment, error)
	Update(ctx context.Context, apartment *models.ServicedApartment) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *models.ListServicedApartmentsRequest) ([]*models.ServicedApartment, int64, error)
	GetUnits(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentUnit, error)
	GetImages(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentImage, error)
	GetFeatured(ctx context.Context, limit int) ([]*models.ServicedApartment, error)
	IncrementViewCount(ctx context.Context, id uint) error
}

type servicedApartmentRepository struct {
	db *gorm.DB
}

func NewServicedApartmentRepository(db *gorm.DB) ServicedApartmentRepository {
	return &servicedApartmentRepository{db: db}
}

func (r *servicedApartmentRepository) Create(ctx context.Context, apartment *models.ServicedApartment) error {
	return r.db.WithContext(ctx).Create(apartment).Error
}

func (r *servicedApartmentRepository) GetByID(ctx context.Context, id uint) (*models.ServicedApartment, error) {
	var apartment models.ServicedApartment
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Units").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Facilities").
		First(&apartment, id).Error
	if err != nil {
		return nil, err
	}
	return &apartment, nil
}

func (r *servicedApartmentRepository) Update(ctx context.Context, apartment *models.ServicedApartment) error {
	return r.db.WithContext(ctx).Save(apartment).Error
}

func (r *servicedApartmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ServicedApartment{}, id).Error
}

func (r *servicedApartmentRepository) List(ctx context.Context, req *models.ListServicedApartmentsRequest) ([]*models.ServicedApartment, int64, error) {
	var apartments []*models.ServicedApartment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ServicedApartment{})

	// 筛选条件
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}
	if req.IsFeatured != nil {
		query = query.Where("is_featured = ?", *req.IsFeatured)
	}
	if req.MinStayDays != nil {
		query = query.Where("min_stay_days >= ?", *req.MinStayDays)
	}

	// 价格筛选需要通过 units 表
	if req.MinPrice != nil || req.MaxPrice != nil {
		subQuery := r.db.Model(&models.ServicedApartmentUnit{}).Select("serviced_apartment_id")
		if req.MinPrice != nil {
			subQuery = subQuery.Where("monthly_price >= ?", *req.MinPrice)
		}
		if req.MaxPrice != nil {
			subQuery = subQuery.Where("monthly_price <= ?", *req.MaxPrice)
		}
		query = query.Where("id IN (?)", subQuery)
	}

	// 统计
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页排序
	offset := (req.Page - 1) * req.PageSize
	sortColumn := req.SortBy
	if sortColumn == "" {
		sortColumn = "created_at"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	query = query.Offset(offset).Limit(req.PageSize).Order(sortColumn + " " + sortOrder)
	query = query.Preload("District").Preload("Units").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC").Limit(1)
	})

	if err := query.Find(&apartments).Error; err != nil {
		return nil, 0, err
	}

	return apartments, total, nil
}

func (r *servicedApartmentRepository) GetUnits(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentUnit, error) {
	var units []models.ServicedApartmentUnit
	err := r.db.WithContext(ctx).
		Where("serviced_apartment_id = ?", apartmentID).
		Order("monthly_price ASC").
		Find(&units).Error
	return units, err
}

func (r *servicedApartmentRepository) GetImages(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentImage, error) {
	var images []models.ServicedApartmentImage
	err := r.db.WithContext(ctx).
		Where("serviced_apartment_id = ?", apartmentID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

func (r *servicedApartmentRepository) GetFeatured(ctx context.Context, limit int) ([]*models.ServicedApartment, error) {
	var apartments []*models.ServicedApartment
	err := r.db.WithContext(ctx).
		Where("is_featured = ? AND status = ?", true, models.ServicedApartmentStatusActive).
		Preload("District").
		Preload("Units").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1)
		}).
		Order("created_at DESC").
		Limit(limit).
		Find(&apartments).Error
	return apartments, err
}

func (r *servicedApartmentRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.ServicedApartment{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
