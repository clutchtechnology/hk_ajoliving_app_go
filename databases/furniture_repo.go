package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// FurnitureRepo 家具仓储
type FurnitureRepo struct {
	db *gorm.DB
}

// NewFurnitureRepo 创建家具仓储
func NewFurnitureRepo(db *gorm.DB) *FurnitureRepo {
	return &FurnitureRepo{db: db}
}

// FindAll 查询家具列表
func (r *FurnitureRepo) FindAll(ctx context.Context, filter *models.ListFurnitureRequest) ([]models.Furniture, int64, error) {
	var furniture []models.Furniture
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

	if filter.Brand != nil && *filter.Brand != "" {
		query = query.Where("brand LIKE ?", "%"+*filter.Brand+"%")
	}

	if filter.Condition != nil {
		query = query.Where("condition = ?", *filter.Condition)
	}

	if filter.DeliveryDistrictID != nil {
		query = query.Where("delivery_district_id = ?", *filter.DeliveryDistrictID)
	}

	if filter.DeliveryMethod != nil {
		query = query.Where("delivery_method = ?", *filter.DeliveryMethod)
	}

	// 状态筛选
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	} else {
		query = query.Where("status = ?", "available")
	}

	// 关键词搜索
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR brand LIKE ?", keyword, keyword, keyword)
	}

	// 只显示未过期的
	query = query.Where("expires_at > ?", time.Now())

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	sortBy := "published_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "price":
			sortBy = "price"
		case "published_at":
			sortBy = "published_at"
		case "view_count":
			sortBy = "view_count"
		}
	}

	sortOrder := "desc"
	if filter.SortOrder == "asc" {
		sortOrder = "asc"
	}

	query = query.Order(sortBy + " " + sortOrder)

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// 预加载关联
	query = query.Preload("Category").Preload("DeliveryDistrict").Preload("Images")

	if err := query.Find(&furniture).Error; err != nil {
		return nil, 0, err
	}

	return furniture, total, nil
}

// FindByID 根据ID查询家具
func (r *FurnitureRepo) FindByID(ctx context.Context, id uint) (*models.Furniture, error) {
	var furniture models.Furniture
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("DeliveryDistrict").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&furniture, id).Error; err != nil {
		return nil, err
	}
	return &furniture, nil
}

// FindFeatured 查询精选家具
func (r *FurnitureRepo) FindFeatured(ctx context.Context, limit int) ([]models.Furniture, error) {
	var furniture []models.Furniture
	if err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at > ?", "available", time.Now()).
		Order("view_count DESC, published_at DESC").
		Limit(limit).
		Preload("Category").
		Preload("DeliveryDistrict").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Find(&furniture).Error; err != nil {
		return nil, err
	}
	return furniture, nil
}

// FindImagesByFurnitureID 查询家具图片
func (r *FurnitureRepo) FindImagesByFurnitureID(ctx context.Context, furnitureID uint) ([]models.FurnitureImage, error) {
	var images []models.FurnitureImage
	if err := r.db.WithContext(ctx).
		Where("furniture_id = ?", furnitureID).
		Order("sort_order ASC").
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// Create 创建家具
func (r *FurnitureRepo) Create(ctx context.Context, furniture *models.Furniture) error {
	return r.db.WithContext(ctx).Create(furniture).Error
}

// Update 更新家具
func (r *FurnitureRepo) Update(ctx context.Context, furniture *models.Furniture) error {
	return r.db.WithContext(ctx).Save(furniture).Error
}

// UpdateStatus 更新家具状态
func (r *FurnitureRepo) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Furniture{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除家具（软删除）
func (r *FurnitureRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Furniture{}, id).Error
}

// IncrementViewCount 增加浏览次数
func (r *FurnitureRepo) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Furniture{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GenerateFurnitureNo 生成家具编号
func (r *FurnitureRepo) GenerateFurnitureNo(ctx context.Context) (string, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Furniture{}).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("FUR%s%06d", time.Now().Format("20060102"), count+1), nil
}

// CreateImages 批量创建图片
func (r *FurnitureRepo) CreateImages(ctx context.Context, images []models.FurnitureImage) error {
	if len(images) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&images).Error
}

// DeleteImages 删除家具的所有图片
func (r *FurnitureRepo) DeleteImages(ctx context.Context, furnitureID uint) error {
	return r.db.WithContext(ctx).
		Where("furniture_id = ?", furnitureID).
		Delete(&models.FurnitureImage{}).Error
}

// UpdateImages 更新家具图片
func (r *FurnitureRepo) UpdateImages(ctx context.Context, furnitureID uint, imageURLs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧图片
		if err := tx.Where("furniture_id = ?", furnitureID).Delete(&models.FurnitureImage{}).Error; err != nil {
			return err
		}

		// 添加新图片
		if len(imageURLs) > 0 {
			images := make([]models.FurnitureImage, len(imageURLs))
			for i, url := range imageURLs {
				images[i] = models.FurnitureImage{
					FurnitureID: furnitureID,
					ImageURL:    url,
					IsCover:     i == 0, // 第一张为封面
					SortOrder:   i,
					CreatedAt:   time.Now(),
				}
			}
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ===== 分类相关 =====

// FindAllCategories 查询所有分类
func (r *FurnitureRepo) FindAllCategories(ctx context.Context) ([]models.FurnitureCategory, error) {
	var categories []models.FurnitureCategory
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND parent_id IS NULL", true).
		Order("sort_order ASC").
		Preload("SubCategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindCategoryByID 根据ID查询分类
func (r *FurnitureRepo) FindCategoryByID(ctx context.Context, id uint) (*models.FurnitureCategory, error) {
	var category models.FurnitureCategory
	if err := r.db.WithContext(ctx).
		Preload("SubCategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetFurnitureCountByCategory 获取分类下的家具数量
func (r *FurnitureRepo) GetFurnitureCountByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Furniture{}).
		Where("category_id = ? AND status = ? AND expires_at > ?", categoryID, "available", time.Now()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
