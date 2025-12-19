package repository

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"gorm.io/gorm"
)

// ServicedApartmentRepository 服务式公寓数据仓库接口
type ServicedApartmentRepository interface {
	Create(ctx context.Context, apartment *model.ServicedApartment) error
	GetByID(ctx context.Context, id uint) (*model.ServicedApartment, error)
	Update(ctx context.Context, apartment *model.ServicedApartment) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *request.ListServicedApartmentsRequest) ([]*model.ServicedApartment, int64, error)
	GetUnits(ctx context.Context, apartmentID uint) ([]model.ServicedApartmentUnit, error)
	GetImages(ctx context.Context, apartmentID uint) ([]model.ServicedApartmentImage, error)
	GetFeatured(ctx context.Context, limit int) ([]*model.ServicedApartment, error)
	IncrementViewCount(ctx context.Context, id uint) error
}

type servicedApartmentRepository struct {
	db *gorm.DB
}

func NewServicedApartmentRepository(db *gorm.DB) ServicedApartmentRepository {
	return &servicedApartmentRepository{db: db}
}

func (r *servicedApartmentRepository) Create(ctx context.Context, apartment *model.ServicedApartment) error {
	return r.db.WithContext(ctx).Create(apartment).Error
}

func (r *servicedApartmentRepository) GetByID(ctx context.Context, id uint) (*model.ServicedApartment, error) {
	var apartment model.ServicedApartment
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

func (r *servicedApartmentRepository) Update(ctx context.Context, apartment *model.ServicedApartment) error {
	return r.db.WithContext(ctx).Save(apartment).Error
}

func (r *servicedApartmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ServicedApartment{}, id).Error
}

func (r *servicedApartmentRepository) List(ctx context.Context, req *request.ListServicedApartmentsRequest) ([]*model.ServicedApartment, int64, error) {
	var apartments []*model.ServicedApartment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ServicedApartment{})

	// 筛选条件
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.IsFeatured != nil {
		query = query.Where("is_featured = ?", *req.IsFeatured)
	}
	if req.MinStayDays != nil {
		query = query.Where("min_stay_days >= ?", *req.MinStayDays)
	}

	// 价格筛选需要通过 units 表
	if req.MinPrice != nil || req.MaxPrice != nil {
		subQuery := r.db.Model(&model.ServicedApartmentUnit{}).Select("serviced_apartment_id")
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

func (r *servicedApartmentRepository) GetUnits(ctx context.Context, apartmentID uint) ([]model.ServicedApartmentUnit, error) {
	var units []model.ServicedApartmentUnit
	err := r.db.WithContext(ctx).
		Where("serviced_apartment_id = ?", apartmentID).
		Order("monthly_price ASC").
		Find(&units).Error
	return units, err
}

func (r *servicedApartmentRepository) GetImages(ctx context.Context, apartmentID uint) ([]model.ServicedApartmentImage, error) {
	var images []model.ServicedApartmentImage
	err := r.db.WithContext(ctx).
		Where("serviced_apartment_id = ?", apartmentID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

func (r *servicedApartmentRepository) GetFeatured(ctx context.Context, limit int) ([]*model.ServicedApartment, error) {
	var apartments []*model.ServicedApartment
	err := r.db.WithContext(ctx).
		Where("is_featured = ? AND status = ?", true, model.ServicedApartmentStatusActive).
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
		Model(&model.ServicedApartment{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
