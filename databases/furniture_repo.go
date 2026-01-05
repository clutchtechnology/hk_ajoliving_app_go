package databases

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// FurnitureRepository 家具数据访问接口
type FurnitureRepository interface {
	Create(ctx context.Context, furniture *models.Furniture) error
	GetByID(ctx context.Context, id uint) (*models.Furniture, error)
	GetByFurnitureNo(ctx context.Context, furnitureNo string) (*models.Furniture, error)
	Update(ctx context.Context, furniture *models.Furniture) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter *models.ListFurnitureRequest) ([]*models.Furniture, int64, error)
	GetFeatured(ctx context.Context, limit int) ([]*models.Furniture, error)
	GetByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]*models.Furniture, int64, error)
	IncrementViewCount(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status models.FurnitureStatus) error
	
	// 分类相关
	GetAllCategories(ctx context.Context) ([]*models.FurnitureCategory, error)
	GetCategoryByID(ctx context.Context, id uint) (*models.FurnitureCategory, error)
	
	// 图片相关
	GetImagesByFurnitureID(ctx context.Context, furnitureID uint) ([]*models.FurnitureImage, error)
}

type furnitureRepository struct {
	db *gorm.DB
}

// NewFurnitureRepository 创建家具仓库
func NewFurnitureRepository(db *gorm.DB) FurnitureRepository {
	return &furnitureRepository{db: db}
}

func (r *furnitureRepository) Create(ctx context.Context, furniture *models.Furniture) error {
	return r.db.WithContext(ctx).Create(furniture).Error
}

func (r *furnitureRepository) GetByID(ctx context.Context, id uint) (*models.Furniture, error) {
	var furniture models.Furniture
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Publisher").
		Preload("DeliveryDistrict").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&furniture, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &furniture, nil
}

func (r *furnitureRepository) GetByFurnitureNo(ctx context.Context, furnitureNo string) (*models.Furniture, error) {
	var furniture models.Furniture
	err := r.db.WithContext(ctx).
		Where("furniture_no = ?", furnitureNo).
		First(&furniture).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &furniture, nil
}

func (r *furnitureRepository) Update(ctx context.Context, furniture *models.Furniture) error {
	return r.db.WithContext(ctx).Save(furniture).Error
}

func (r *furnitureRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Furniture{}, id).Error
}

func (r *furnitureRepository) List(ctx context.Context, filter *models.ListFurnitureRequest) ([]*models.Furniture, int64, error) {
	var furniture []*models.Furniture
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Furniture{})

	// 应用筛选条件
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.Condition != nil && *filter.Condition != "" {
		query = query.Where("condition = ?", *filter.Condition)
	}
	if filter.Brand != nil && *filter.Brand != "" {
		query = query.Where("brand ILIKE ?", "%"+*filter.Brand+"%")
	}
	if filter.DeliveryDistrictID != nil {
		query = query.Where("delivery_district_id = ?", *filter.DeliveryDistrictID)
	}
	if filter.DeliveryMethod != nil && *filter.DeliveryMethod != "" {
		query = query.Where("delivery_method = ?", *filter.DeliveryMethod)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+ *filter.Keyword+"%", "%"+ *filter.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 设置默认值
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// 分页和排序
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	query = query.Order(sortBy + " " + sortOrder)

	// 预加载关联
	query = query.Preload("Category").
		Preload("DeliveryDistrict").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Order("sort_order ASC").Limit(1)
		})

	if err := query.Find(&furniture).Error; err != nil {
		return nil, 0, err
	}

	return furniture, total, nil
}

func (r *furnitureRepository) GetFeatured(ctx context.Context, limit int) ([]*models.Furniture, error) {
	var furniture []*models.Furniture
	
	if limit <= 0 {
		limit = 10
	}

	err := r.db.WithContext(ctx).
		Where("status = ?", models.FurnitureStatusAvailable).
		Where("expires_at > NOW()").
		Order("view_count DESC, favorite_count DESC").
		Limit(limit).
		Preload("Category").
		Preload("DeliveryDistrict").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Order("sort_order ASC").Limit(1)
		}).
		Find(&furniture).Error

	if err != nil {
		return nil, err
	}

	return furniture, nil
}

func (r *furnitureRepository) GetByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]*models.Furniture, int64, error) {
	var furniture []*models.Furniture
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Furniture{}).
		Where("category_id = ?", categoryID).
		Where("status = ?", models.FurnitureStatusAvailable)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 设置默认值
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// 预加载关联
	query = query.Preload("Category").
		Preload("DeliveryDistrict").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Order("sort_order ASC").Limit(1)
		})

	if err := query.Find(&furniture).Error; err != nil {
		return nil, 0, err
	}

	return furniture, total, nil
}

func (r *furnitureRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Furniture{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}

func (r *furnitureRepository) UpdateStatus(ctx context.Context, id uint, status models.FurnitureStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Furniture{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// 分类相关

func (r *furnitureRepository) GetAllCategories(ctx context.Context) ([]*models.FurnitureCategory, error) {
	var categories []*models.FurnitureCategory
	
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *furnitureRepository) GetCategoryByID(ctx context.Context, id uint) (*models.FurnitureCategory, error) {
	var category models.FurnitureCategory
	
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		First(&category, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &category, nil
}

// 图片相关

func (r *furnitureRepository) GetImagesByFurnitureID(ctx context.Context, furnitureID uint) ([]*models.FurnitureImage, error) {
	var images []*models.FurnitureImage
	
	err := r.db.WithContext(ctx).
		Where("furniture_id = ?", furnitureID).
		Order("sort_order ASC").
		Find(&images).Error
		
	if err != nil {
		return nil, err
	}
	
	return images, nil
}
