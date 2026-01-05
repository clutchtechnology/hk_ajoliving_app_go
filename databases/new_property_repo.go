package databases

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// NewPropertyRepository 新楼盘数据仓库接口
type NewPropertyRepository interface {
	// Create 创建新楼盘
	Create(ctx context.Context, newProperty *models.NewProperty) error
	// GetByID 根据ID获取新楼盘
	GetByID(ctx context.Context, id uint) (*models.NewProperty, error)
	// Update 更新新楼盘
	Update(ctx context.Context, newProperty *models.NewProperty) error
	// Delete 删除新楼盘
	Delete(ctx context.Context, id uint) error
	// List 获取新楼盘列表
	List(ctx context.Context, req *models.ListNewDevelopmentsRequest) ([]*models.NewProperty, int64, error)
	// GetLayouts 获取新楼盘户型列表
	GetLayouts(ctx context.Context, newPropertyID uint) ([]models.NewPropertyLayout, error)
	// GetImages 获取新楼盘图片列表
	GetImages(ctx context.Context, newPropertyID uint) ([]models.NewPropertyImage, error)
	// GetFeatured 获取精选新楼盘
	GetFeatured(ctx context.Context, limit int) ([]*models.NewProperty, error)
	// IncrementViewCount 增加浏览次数
	IncrementViewCount(ctx context.Context, id uint) error
	// GetPriceRange 获取新楼盘价格范围
	GetPriceRange(ctx context.Context, newPropertyID uint) (minPrice, maxPrice float64, err error)
	// GetCoverImage 获取新楼盘封面图
	GetCoverImage(ctx context.Context, newPropertyID uint) (string, error)
}

// newPropertyRepository 新楼盘数据仓库实现
type newPropertyRepository struct {
	db *gorm.DB
}

// NewNewPropertyRepository 创建新楼盘仓库实例
func NewNewPropertyRepository(db *gorm.DB) NewPropertyRepository {
	return &newPropertyRepository{db: db}
}

// Create 创建新楼盘
func (r *newPropertyRepository) Create(ctx context.Context, newProperty *models.NewProperty) error {
	return r.db.WithContext(ctx).Create(newProperty).Error
}

// GetByID 根据ID获取新楼盘
func (r *newPropertyRepository) GetByID(ctx context.Context, id uint) (*models.NewProperty, error) {
	var newProperty models.NewProperty
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Layouts").
		First(&newProperty, id).Error
	if err != nil {
		return nil, err
	}
	return &newProperty, nil
}

// Update 更新新楼盘
func (r *newPropertyRepository) Update(ctx context.Context, newProperty *models.NewProperty) error {
	return r.db.WithContext(ctx).Save(newProperty).Error
}

// Delete 删除新楼盘（软删除）
func (r *newPropertyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.NewProperty{}, id).Error
}

// List 获取新楼盘列表
func (r *newPropertyRepository) List(ctx context.Context, req *models.ListNewDevelopmentsRequest) ([]*models.NewProperty, int64, error) {
	var newProperties []*models.NewProperty
	var total int64

	query := r.db.WithContext(ctx).Model(&models.NewProperty{})

	// 应用筛选条件
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}

	if req.Developer != nil && *req.Developer != "" {
		query = query.Where("developer LIKE ?", "%" + *req.Developer + "%")
	}

	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}

	if req.IsFeatured != nil {
		query = query.Where("is_featured = ?", *req.IsFeatured)
	}

	if req.SchoolNet != nil && *req.SchoolNet != "" {
		query = query.Where("primary_school_net = ? OR secondary_school_net = ?", *req.SchoolNet, *req.SchoolNet)
	}

	// 价格筛选需要通过 layouts 表联合查询
	if req.MinPrice != nil || req.MaxPrice != nil || req.Bedrooms != nil {
		subQuery := r.db.Model(&models.NewPropertyLayout{}).Select("new_property_id")
		
		if req.MinPrice != nil {
			subQuery = subQuery.Where("min_price >= ?", *req.MinPrice)
		}
		if req.MaxPrice != nil {
			subQuery = subQuery.Where("min_price <= ?", *req.MaxPrice)
		}
		if req.Bedrooms != nil {
			subQuery = subQuery.Where("bedrooms = ?", *req.Bedrooms)
		}
		
		query = query.Where("id IN (?)", subQuery)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 排序
	sortColumn := req.SortBy
	if sortColumn == "" {
		sortColumn = "created_at"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	query = query.Order(sortColumn + " " + sortOrder)

	// 预加载
	query = query.Preload("District").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Where("image_type = ?", "cover").Order("sort_order ASC").Limit(1)
	}).Preload("Layouts")

	if err := query.Find(&newProperties).Error; err != nil {
		return nil, 0, err
	}

	return newProperties, total, nil
}

// GetLayouts 获取新楼盘户型列表
func (r *newPropertyRepository) GetLayouts(ctx context.Context, newPropertyID uint) ([]models.NewPropertyLayout, error) {
	var layouts []models.NewPropertyLayout
	err := r.db.WithContext(ctx).
		Where("new_property_id = ?", newPropertyID).
		Order("bedrooms ASC, saleable_area ASC").
		Find(&layouts).Error
	return layouts, err
}

// GetImages 获取新楼盘图片列表
func (r *newPropertyRepository) GetImages(ctx context.Context, newPropertyID uint) ([]models.NewPropertyImage, error) {
	var images []models.NewPropertyImage
	err := r.db.WithContext(ctx).
		Where("new_property_id = ?", newPropertyID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

// GetFeatured 获取精选新楼盘
func (r *newPropertyRepository) GetFeatured(ctx context.Context, limit int) ([]*models.NewProperty, error) {
	var newProperties []*models.NewProperty
	err := r.db.WithContext(ctx).
		Where("is_featured = ?", true).
		Where("status IN ?", []string{string(models.NewPropertyStatusPresale), string(models.NewPropertyStatusSelling)}).
		Preload("District").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("image_type = ?", "cover").Limit(1)
		}).
		Preload("Layouts").
		Order("created_at DESC").
		Limit(limit).
		Find(&newProperties).Error
	return newProperties, err
}

// IncrementViewCount 增加浏览次数
func (r *newPropertyRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.NewProperty{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetPriceRange 获取新楼盘价格范围
func (r *newPropertyRepository) GetPriceRange(ctx context.Context, newPropertyID uint) (minPrice, maxPrice float64, err error) {
	var result struct {
		MinPrice float64
		MaxPrice float64
	}
	err = r.db.WithContext(ctx).
		Model(&models.NewPropertyLayout{}).
		Select("MIN(min_price) as min_price, MAX(COALESCE(max_price, min_price)) as max_price").
		Where("new_property_id = ?", newPropertyID).
		Scan(&result).Error
	return result.MinPrice, result.MaxPrice, err
}

// GetCoverImage 获取新楼盘封面图
func (r *newPropertyRepository) GetCoverImage(ctx context.Context, newPropertyID uint) (string, error) {
	var image models.NewPropertyImage
	err := r.db.WithContext(ctx).
		Where("new_property_id = ? AND image_type = ?", newPropertyID, "cover").
		Order("sort_order ASC").
		First(&image).Error
	if err != nil {
		// 如果没有 cover 类型的图片，返回第一张图片
		err = r.db.WithContext(ctx).
			Where("new_property_id = ?", newPropertyID).
			Order("sort_order ASC").
			First(&image).Error
		if err != nil {
			return "", nil // 没有图片不返回错误
		}
	}
	return image.ImageURL, nil
}
