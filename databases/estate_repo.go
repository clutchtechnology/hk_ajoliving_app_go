package databases

import (
	"context"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// EstateRepo 屋苑仓储
type EstateRepo struct {
	db *gorm.DB
}

// NewEstateRepo 创建屋苑仓储
func NewEstateRepo(db *gorm.DB) *EstateRepo {
	return &EstateRepo{db: db}
}

// FindAll 查询屋苑列表
func (r *EstateRepo) FindAll(ctx context.Context, filter *models.ListEstatesRequest) ([]models.Estate, int64, error) {
	var estates []models.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Estate{})

	// 应用筛选条件
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}

	if filter.PrimarySchoolNet != nil && *filter.PrimarySchoolNet != "" {
		query = query.Where("primary_school_net = ?", *filter.PrimarySchoolNet)
	}

	if filter.SecondarySchoolNet != nil && *filter.SecondarySchoolNet != "" {
		query = query.Where("secondary_school_net = ?", *filter.SecondarySchoolNet)
	}

	if filter.Developer != nil && *filter.Developer != "" {
		query = query.Where("developer LIKE ?", "%"+*filter.Developer+"%")
	}

	if filter.MinAvgPrice != nil {
		query = query.Where("avg_transaction_price >= ?", *filter.MinAvgPrice)
	}

	if filter.MaxAvgPrice != nil {
		query = query.Where("avg_transaction_price <= ?", *filter.MaxAvgPrice)
	}

	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}

	// 关键词搜索
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name LIKE ? OR name_en LIKE ? OR address LIKE ?", keyword, keyword, keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	sortBy := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "name":
			sortBy = "name"
		case "recent_transactions":
			sortBy = "recent_transactions_count"
		case "avg_price":
			sortBy = "avg_transaction_price"
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
	query = query.Preload("District").Preload("Images").Preload("Facilities")

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

// FindByID 根据ID查询屋苑
func (r *EstateRepo) FindByID(ctx context.Context, id uint) (*models.Estate, error) {
	var estate models.Estate
	if err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Images").
		Preload("Facilities").
		First(&estate, id).Error; err != nil {
		return nil, err
	}
	return &estate, nil
}

// FindFeatured 查询精选屋苑
func (r *EstateRepo) FindFeatured(ctx context.Context, limit int) ([]models.Estate, error) {
	var estates []models.Estate
	if err := r.db.WithContext(ctx).
		Where("is_featured = ?", true).
		Order("view_count DESC").
		Limit(limit).
		Preload("District").
		Preload("Images").
		Find(&estates).Error; err != nil {
		return nil, err
	}
	return estates, nil
}

// FindPropertiesByEstate 查询屋苑内的房源
func (r *EstateRepo) FindPropertiesByEstate(ctx context.Context, estateName string, filter *models.GetEstatePropertiesRequest) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Property{}).Where("building_name = ?", estateName)

	// 应用筛选条件
	if filter.ListingType != nil {
		query = query.Where("listing_type = ?", *filter.ListingType)
	}

	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	if filter.Bedrooms != nil {
		query = query.Where("bedrooms = ?", *filter.Bedrooms)
	}

	if filter.PropertyType != nil {
		query = query.Where("property_type = ?", *filter.PropertyType)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	} else {
		query = query.Where("status = ?", "available")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	// 预加载关联
	query = query.Preload("District").Preload("Images").Order("created_at DESC")

	if err := query.Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// FindImagesByEstateID 查询屋苑图片
func (r *EstateRepo) FindImagesByEstateID(ctx context.Context, estateID uint) ([]models.EstateImage, error) {
	var images []models.EstateImage
	if err := r.db.WithContext(ctx).
		Where("estate_id = ?", estateID).
		Order("sort_order ASC").
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// FindFacilitiesByEstateID 查询屋苑设施
func (r *EstateRepo) FindFacilitiesByEstateID(ctx context.Context, estateID uint) ([]models.Facility, error) {
	var facilities []models.Facility
	if err := r.db.WithContext(ctx).
		Joins("JOIN estate_facilities ON estate_facilities.facility_id = facilities.id").
		Where("estate_facilities.estate_id = ?", estateID).
		Order("facilities.sort_order ASC").
		Find(&facilities).Error; err != nil {
		return nil, err
	}
	return facilities, nil
}

// GetStatistics 获取屋苑统计数据
func (r *EstateRepo) GetStatistics(ctx context.Context, estateID uint) (*models.EstateStatisticsResponse, error) {
	var estate models.Estate
	if err := r.db.WithContext(ctx).First(&estate, estateID).Error; err != nil {
		return nil, err
	}

	// 查询价格统计
	var priceStats struct {
		MinPrice float64
		MaxPrice float64
		AvgArea  float64
	}

	r.db.WithContext(ctx).
		Model(&models.Property{}).
		Where("building_name = ? AND status = ?", estate.Name, "available").
		Select("MIN(price) as min_price, MAX(price) as max_price, AVG(area) as avg_area").
		Scan(&priceStats)

	stats := &models.EstateStatisticsResponse{
		EstateID:                estateID,
		EstateName:              estate.Name,
		TotalUnits:              estate.TotalUnits,
		ForSaleCount:            estate.ForSaleCount,
		ForRentCount:            estate.ForRentCount,
		RecentTransactionsCount: estate.RecentTransactionsCount,
		AvgTransactionPrice:     estate.AvgTransactionPrice,
		MinPrice:                priceStats.MinPrice,
		MaxPrice:                priceStats.MaxPrice,
		AvgArea:                 priceStats.AvgArea,
		ViewCount:               estate.ViewCount,
		FavoriteCount:           estate.FavoriteCount,
	}

	return stats, nil
}

// Create 创建屋苑
func (r *EstateRepo) Create(ctx context.Context, estate *models.Estate) error {
	return r.db.WithContext(ctx).Create(estate).Error
}

// Update 更新屋苑
func (r *EstateRepo) Update(ctx context.Context, estate *models.Estate) error {
	return r.db.WithContext(ctx).Save(estate).Error
}

// Delete 删除屋苑（软删除）
func (r *EstateRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Estate{}, id).Error
}

// IncrementViewCount 增加浏览次数
func (r *EstateRepo) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Estate{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// UpdateCounts 更新屋苑的统计数据（定时任务调用）
func (r *EstateRepo) UpdateCounts(ctx context.Context, estateID uint) error {
	var estate models.Estate
	if err := r.db.WithContext(ctx).First(&estate, estateID).Error; err != nil {
		return err
	}

	// 统计放盘数量
	var forSaleCount, forRentCount int64
	r.db.WithContext(ctx).
		Model(&models.Property{}).
		Where("building_name = ? AND status = ? AND listing_type = ?", estate.Name, "available", "sale").
		Count(&forSaleCount)

	r.db.WithContext(ctx).
		Model(&models.Property{}).
		Where("building_name = ? AND status = ? AND listing_type = ?", estate.Name, "available", "rent").
		Count(&forRentCount)

	// 更新统计数据
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Estate{}).
		Where("id = ?", estateID).
		Updates(map[string]interface{}{
			"for_sale_count": forSaleCount,
			"for_rent_count": forRentCount,
			"updated_at":     now,
		}).Error
}

// UpdateFacilities 更新屋苑设施
func (r *EstateRepo) UpdateFacilities(ctx context.Context, estateID uint, facilityIDs []uint) error {
	// 使用事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的设施关联
		if err := tx.Where("estate_id = ?", estateID).Delete(&models.EstateFacility{}).Error; err != nil {
			return err
		}

		// 添加新的设施关联
		if len(facilityIDs) > 0 {
			facilities := make([]models.EstateFacility, len(facilityIDs))
			for i, facilityID := range facilityIDs {
				facilities[i] = models.EstateFacility{
					EstateID:   estateID,
					FacilityID: facilityID,
					CreatedAt:  time.Now(),
				}
			}
			if err := tx.Create(&facilities).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
