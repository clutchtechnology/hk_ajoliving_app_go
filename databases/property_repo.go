package databases

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// PropertyRepo 房产仓储
type PropertyRepo struct {
	db *gorm.DB
}

// NewPropertyRepo 创建房产仓储
func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
	return &PropertyRepo{db: db}
}

// Create 创建房产
func (r *PropertyRepo) Create(ctx context.Context, property *models.Property) error {
	return r.db.WithContext(ctx).Create(property).Error
}

// FindByID 根据ID查找房产
func (r *PropertyRepo) FindByID(ctx context.Context, id uint) (*models.Property, error) {
	var property models.Property
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&property, id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property not found")
		}
		return nil, err
	}
	return &property, nil
}

// FindAll 查找所有房产（支持筛选和分页）
func (r *PropertyRepo) FindAll(ctx context.Context, filter *models.ListPropertiesRequest) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Property{})

	// 应用筛选条件
	if filter.ListingType != nil {
		query = query.Where("listing_type = ?", *filter.ListingType)
	}
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
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
	if filter.PropertyType != nil {
		query = query.Where("property_type = ?", *filter.PropertyType)
	}
	if filter.BuildingName != nil {
		query = query.Where("building_name LIKE ?", "%"+*filter.BuildingName+"%")
	}
	if filter.PrimarySchool != nil {
		query = query.Where("primary_school_net = ?", *filter.PrimarySchool)
	}
	if filter.SecondarySchool != nil {
		query = query.Where("secondary_school_net = ?", *filter.SecondarySchool)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	} else {
		// 默认只显示可用状态
		query = query.Where("status = ?", "available")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch filter.SortBy {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "area_asc":
		query = query.Order("area ASC")
	case "area_desc":
		query = query.Order("area DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// 预加载关联
	query = query.Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1) // 列表只加载第一张图
		})

	if err := query.Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// Update 更新房产
func (r *PropertyRepo) Update(ctx context.Context, property *models.Property) error {
	return r.db.WithContext(ctx).Save(property).Error
}

// Delete 删除房产（软删除）
func (r *PropertyRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Property{}, id).Error
}

// IncrementViewCount 增加浏览次数
func (r *PropertyRepo) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Property{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// FindSimilar 查找相似房源
func (r *PropertyRepo) FindSimilar(ctx context.Context, property *models.Property, limit int) ([]models.Property, error) {
	var similar []models.Property

	// 相似条件：同地区、同类型、价格相近
	priceMin := property.Price * 0.8
	priceMax := property.Price * 1.2

	err := r.db.WithContext(ctx).
		Where("id != ?", property.ID).
		Where("district_id = ?", property.DistrictID).
		Where("listing_type = ?", property.ListingType).
		Where("property_type = ?", property.PropertyType).
		Where("price BETWEEN ? AND ?", priceMin, priceMax).
		Where("status = ?", "available").
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1)
		}).
		Order("RANDOM()").
		Limit(limit).
		Find(&similar).Error

	return similar, err
}

// FindFeatured 查找精选房源
func (r *PropertyRepo) FindFeatured(ctx context.Context, limit int) ([]models.Property, error) {
	var featured []models.Property

	err := r.db.WithContext(ctx).
		Where("status = ?", "available").
		Where("published_at IS NOT NULL").
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1)
		}).
		Order("favorite_count DESC, view_count DESC").
		Limit(limit).
		Find(&featured).Error

	return featured, err
}

// FindHot 查找热门房源
func (r *PropertyRepo) FindHot(ctx context.Context, limit int) ([]models.Property, error) {
	var hot []models.Property

	err := r.db.WithContext(ctx).
		Where("status = ?", "available").
		Where("created_at >= NOW() - INTERVAL '30 days'"). // 最近30天
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1)
		}).
		Order("view_count DESC, created_at DESC").
		Limit(limit).
		Find(&hot).Error

	return hot, err
}

// FindByPublisher 根据发布者查找房产
func (r *PropertyRepo) FindByPublisher(ctx context.Context, publisherID uint, listingType *string) ([]models.Property, error) {
	var properties []models.Property

	query := r.db.WithContext(ctx).
		Where("publisher_id = ?", publisherID).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Limit(1)
		})

	if listingType != nil {
		query = query.Where("listing_type = ?", *listingType)
	}

	err := query.Order("created_at DESC").Find(&properties).Error
	return properties, err
}

// CreateImages 批量创建房产图片
func (r *PropertyRepo) CreateImages(ctx context.Context, images []models.PropertyImage) error {
	if len(images) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&images).Error
}

// GeneratePropertyNo 生成房产编号
func (r *PropertyRepo) GeneratePropertyNo(ctx context.Context) (string, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Property{}).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("P%08d", count+1), nil
}

// CalculateTotalPages 计算总页数
func CalculateTotalPages(total int64, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
