package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// PropertyRepository 房产数据访问接口
type PropertyRepository interface {
	Create(ctx context.Context, property *model.Property) error
	GetByID(ctx context.Context, id uint) (*model.Property, error)
	GetByPropertyNo(ctx context.Context, propertyNo string) (*model.Property, error)
	Update(ctx context.Context, property *model.Property) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter *request.ListPropertiesRequest) ([]*model.Property, int64, error)
	ListByPublisher(ctx context.Context, publisherID uint, page, pageSize int) ([]*model.Property, int64, error)
	GetFeatured(ctx context.Context, listingType string, limit int) ([]*model.Property, error)
	GetHot(ctx context.Context, listingType string, limit int) ([]*model.Property, error)
	GetSimilar(ctx context.Context, property *model.Property, limit int) ([]*model.Property, error)
	IncrementViewCount(ctx context.Context, id uint) error
}

type propertyRepository struct {
	db *gorm.DB
}

// NewPropertyRepository 创建房产仓库
func NewPropertyRepository(db *gorm.DB) PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) Create(ctx context.Context, property *model.Property) error {
	return r.db.WithContext(ctx).Create(property).Error
}

func (r *propertyRepository) GetByID(ctx context.Context, id uint) (*model.Property, error) {
	var property model.Property
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Agent").
		Preload("Images").
		Preload("Facilities").
		First(&property, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &property, nil
}

func (r *propertyRepository) GetByPropertyNo(ctx context.Context, propertyNo string) (*model.Property, error) {
	var property model.Property
	err := r.db.WithContext(ctx).
		Where("property_no = ?", propertyNo).
		First(&property).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &property, nil
}

func (r *propertyRepository) Update(ctx context.Context, property *model.Property) error {
	return r.db.WithContext(ctx).Save(property).Error
}

func (r *propertyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Property{}, id).Error
}

func (r *propertyRepository) List(ctx context.Context, filter *request.ListPropertiesRequest) ([]*model.Property, int64, error) {
	var properties []*model.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Property{})

	// 应用筛选条件
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.BuildingName != "" {
		query = query.Where("building_name ILIKE ?", "%"+filter.BuildingName+"%")
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.MinArea != nil {
		query = query.Where("area >= ?", *filter.MinArea)
	}
	if filter.MaxArea != nil {
		query = query.Where("area <= ?", *filter.MaxArea)
	}
	if filter.Bedrooms != nil {
		query = query.Where("bedrooms = ?", *filter.Bedrooms)
	}
	if filter.PropertyType != "" {
		query = query.Where("property_type = ?", filter.PropertyType)
	}
	if filter.ListingType != "" {
		query = query.Where("listing_type = ?", filter.ListingType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	} else {
		// 默认只显示可用的房源
		query = query.Where("status = ?", model.PropertyStatusAvailable)
	}
	if filter.SchoolNet != "" {
		query = query.Where("primary_school_net = ? OR secondary_school_net = ?", filter.SchoolNet, filter.SchoolNet)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (filter.Page - 1) * filter.PageSize
	orderClause := fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder)
	query = query.Offset(offset).Limit(filter.PageSize).Order(orderClause)

	// 预加载关联
	query = query.Preload("District").Preload("Images", "is_cover = true")

	if err := query.Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *propertyRepository) ListByPublisher(ctx context.Context, publisherID uint, page, pageSize int) ([]*model.Property, int64, error) {
	var properties []*model.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Property{}).Where("publisher_id = ?", publisherID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Preload("District").
		Preload("Images", "is_cover = true").
		Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *propertyRepository) GetFeatured(ctx context.Context, listingType string, limit int) ([]*model.Property, error) {
	var properties []*model.Property

	query := r.db.WithContext(ctx).Model(&model.Property{}).
		Where("status = ?", model.PropertyStatusAvailable)

	if listingType != "" {
		query = query.Where("listing_type = ?", listingType)
	}

	// 精选房源按收藏数和浏览数排序
	if err := query.Order("favorite_count DESC, view_count DESC").
		Limit(limit).
		Preload("District").
		Preload("Images", "is_cover = true").
		Find(&properties).Error; err != nil {
		return nil, err
	}

	return properties, nil
}

func (r *propertyRepository) GetHot(ctx context.Context, listingType string, limit int) ([]*model.Property, error) {
	var properties []*model.Property

	query := r.db.WithContext(ctx).Model(&model.Property{}).
		Where("status = ?", model.PropertyStatusAvailable)

	if listingType != "" {
		query = query.Where("listing_type = ?", listingType)
	}

	// 热门房源按浏览数排序
	if err := query.Order("view_count DESC, created_at DESC").
		Limit(limit).
		Preload("District").
		Preload("Images", "is_cover = true").
		Find(&properties).Error; err != nil {
		return nil, err
	}

	return properties, nil
}

func (r *propertyRepository) GetSimilar(ctx context.Context, property *model.Property, limit int) ([]*model.Property, error) {
	var properties []*model.Property

	// 相似房源：同地区、同类型、价格相近的房源
	priceRange := property.Price * 0.3 // 30% 价格范围

	if err := r.db.WithContext(ctx).Model(&model.Property{}).
		Where("id != ?", property.ID).
		Where("status = ?", model.PropertyStatusAvailable).
		Where("listing_type = ?", property.ListingType).
		Where("district_id = ? OR property_type = ?", property.DistrictID, property.PropertyType).
		Where("price BETWEEN ? AND ?", property.Price-priceRange, property.Price+priceRange).
		Order("view_count DESC").
		Limit(limit).
		Preload("District").
		Preload("Images", "is_cover = true").
		Find(&properties).Error; err != nil {
		return nil, err
	}

	return properties, nil
}

func (r *propertyRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Property{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
