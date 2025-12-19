package repository

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"gorm.io/gorm"
)

// EstateRepository 屋苑数据仓库接口
type EstateRepository interface {
	Create(ctx context.Context, estate *model.Estate) error
	GetByID(ctx context.Context, id uint) (*model.Estate, error)
	Update(ctx context.Context, estate *model.Estate) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *request.ListEstatesRequest) ([]*model.Estate, int64, error)
	GetProperties(ctx context.Context, estateID uint, listingType string, page, pageSize int) ([]*model.Property, int64, error)
	GetImages(ctx context.Context, estateID uint) ([]model.EstateImage, error)
	GetFacilities(ctx context.Context, estateID uint) ([]model.Facility, error)
	GetStatistics(ctx context.Context, estateID uint) (*model.Estate, error)
	GetFeatured(ctx context.Context, limit int) ([]*model.Estate, error)
	IncrementViewCount(ctx context.Context, id uint) error
}

type estateRepository struct {
	db *gorm.DB
}

func NewEstateRepository(db *gorm.DB) EstateRepository {
	return &estateRepository{db: db}
}

func (r *estateRepository) Create(ctx context.Context, estate *model.Estate) error {
	return r.db.WithContext(ctx).Create(estate).Error
}

func (r *estateRepository) GetByID(ctx context.Context, id uint) (*model.Estate, error) {
	var estate model.Estate
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Facilities").
		First(&estate, id).Error
	if err != nil {
		return nil, err
	}
	return &estate, nil
}

func (r *estateRepository) Update(ctx context.Context, estate *model.Estate) error {
	return r.db.WithContext(ctx).Save(estate).Error
}

func (r *estateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Estate{}, id).Error
}

func (r *estateRepository) List(ctx context.Context, req *request.ListEstatesRequest) ([]*model.Estate, int64, error) {
	var estates []*model.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Estate{})

	// 筛选条件
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.SchoolNet != "" {
		query = query.Where("primary_school_net = ? OR secondary_school_net = ?", req.SchoolNet, req.SchoolNet)
	}
	if req.MinCompletionYear != nil {
		query = query.Where("completion_year >= ?", *req.MinCompletionYear)
	}
	if req.MaxCompletionYear != nil {
		query = query.Where("completion_year <= ?", *req.MaxCompletionYear)
	}
	if req.MinAvgPrice != nil {
		query = query.Where("avg_transaction_price >= ?", *req.MinAvgPrice)
	}
	if req.MaxAvgPrice != nil {
		query = query.Where("avg_transaction_price <= ?", *req.MaxAvgPrice)
	}
	if req.HasListings != nil && *req.HasListings {
		query = query.Where("(for_sale_count > 0 OR for_rent_count > 0)")
	}
	if req.HasTransactions != nil && *req.HasTransactions {
		query = query.Where("recent_transactions_count > 0")
	}
	if req.IsFeatured != nil {
		query = query.Where("is_featured = ?", *req.IsFeatured)
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
	query = query.Preload("District").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Where("image_type = ?", "exterior").Order("sort_order ASC").Limit(1)
	})

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

func (r *estateRepository) GetProperties(ctx context.Context, estateID uint, listingType string, page, pageSize int) ([]*model.Property, int64, error) {
	var properties []*model.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Property{}).Where("estate_id = ?", estateID)
	
	if listingType != "" {
		query = query.Where("listing_type = ?", listingType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Limit(1)
		}).
		Order("created_at DESC")

	if err := query.Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *estateRepository) GetImages(ctx context.Context, estateID uint) ([]model.EstateImage, error) {
	var images []model.EstateImage
	err := r.db.WithContext(ctx).
		Where("estate_id = ?", estateID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

func (r *estateRepository) GetFacilities(ctx context.Context, estateID uint) ([]model.Facility, error) {
	var facilities []model.Facility
	err := r.db.WithContext(ctx).
		Joins("JOIN estate_facilities ON estate_facilities.facility_id = facilities.id").
		Where("estate_facilities.estate_id = ?", estateID).
		Find(&facilities).Error
	return facilities, err
}

func (r *estateRepository) GetStatistics(ctx context.Context, estateID uint) (*model.Estate, error) {
	var estate model.Estate
	err := r.db.WithContext(ctx).
		Select("id, name, recent_transactions_count, for_sale_count, for_rent_count, avg_transaction_price, avg_transaction_price_updated_at").
		First(&estate, estateID).Error
	if err != nil {
		return nil, err
	}
	return &estate, nil
}

func (r *estateRepository) GetFeatured(ctx context.Context, limit int) ([]*model.Estate, error) {
	var estates []*model.Estate
	err := r.db.WithContext(ctx).
		Where("is_featured = ?", true).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("image_type = ?", "exterior").Order("sort_order ASC").Limit(1)
		}).
		Order("created_at DESC").
		Limit(limit).
		Find(&estates).Error
	return estates, err
}

func (r *estateRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Estate{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
