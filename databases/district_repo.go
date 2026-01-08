package databases

import (
	"context"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type DistrictRepo struct {
	db *gorm.DB
}

func NewDistrictRepo(db *gorm.DB) *DistrictRepo {
	return &DistrictRepo{db: db}
}

// FindAll 查询所有地区
func (r *DistrictRepo) FindAll(ctx context.Context, region string) ([]models.District, error) {
	var districts []models.District
	query := r.db.WithContext(ctx).Model(&models.District{})

	if region != "" {
		query = query.Where("region = ?", region)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

// FindByID 根据ID查询地区
func (r *DistrictRepo) FindByID(ctx context.Context, id uint) (*models.District, error) {
	var district models.District
	if err := r.db.WithContext(ctx).First(&district, id).Error; err != nil {
		return nil, err
	}
	return &district, nil
}

// GetDistrictProperties 查询地区内的房源
func (r *DistrictRepo) GetDistrictProperties(ctx context.Context, districtID uint, filter *models.GetDistrictPropertiesRequest) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND status = ?", districtID, "available")

	// 应用筛选条件
	if filter.ListingType != nil {
		query = query.Where("listing_type = ?", *filter.ListingType)
	}
	if filter.PropertyType != nil {
		query = query.Where("property_type = ?", *filter.PropertyType)
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

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).
		Preload("District").
		Order("created_at DESC").
		Find(&properties).Error; err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// GetDistrictEstates 查询地区内的屋苑
func (r *DistrictRepo) GetDistrictEstates(ctx context.Context, districtID uint, page, pageSize int) ([]models.Estate, int64, error) {
	var estates []models.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Estate{}).
		Where("district_id = ?", districtID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("District").
		Order("name ASC").
		Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

// GetDistrictStatistics 获取地区统计数据
func (r *DistrictRepo) GetDistrictStatistics(ctx context.Context, districtID uint) (*models.DistrictStatisticsResponse, error) {
	district, err := r.FindByID(ctx, districtID)
	if err != nil {
		return nil, err
	}

	stats := &models.DistrictStatisticsResponse{
		DistrictID:   district.ID,
		DistrictName: district.NameZhHant,
	}

	// 统计总房源数
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND status = ?", districtID, "available").
		Count(&stats.TotalProperties)

	// 统计出售房源数
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND listing_type = ? AND status = ?", districtID, "sale", "available").
		Count(&stats.PropertiesForSale)

	// 统计出租房源数
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND listing_type = ? AND status = ?", districtID, "rent", "available").
		Count(&stats.PropertiesForRent)

	// 统计屋苑数
	r.db.WithContext(ctx).Model(&models.Estate{}).
		Where("district_id = ?", districtID).
		Count(&stats.TotalEstates)

	// 平均售价
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND listing_type = ? AND status = ?", districtID, "sale", "available").
		Select("AVG(price)").
		Scan(&stats.AvgSalePrice)

	// 平均租金
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND listing_type = ? AND status = ?", districtID, "rent", "available").
		Select("AVG(price)").
		Scan(&stats.AvgRentPrice)

	// 平均每平方尺价格（出售）
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND listing_type = ? AND status = ? AND area > 0", districtID, "sale", "available").
		Select("AVG(price / area)").
		Scan(&stats.AvgPricePerSqft)

	// 本周新增房源数（最近7天）
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND created_at >= ?", districtID, oneWeekAgo).
		Count(&stats.NewPropertiesThisWeek)

	// 本月新增房源数（最近30天）
	oneMonthAgo := time.Now().AddDate(0, 0, -30)
	r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND created_at >= ?", districtID, oneMonthAgo).
		Count(&stats.NewPropertiesThisMonth)

	return stats, nil
}

// GetPropertyCount 获取地区房源数量
func (r *DistrictRepo) GetPropertyCount(ctx context.Context, districtID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND status = ?", districtID, "available").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetEstateCount 获取地区屋苑数量
func (r *DistrictRepo) GetEstateCount(ctx context.Context, districtID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Estate{}).
		Where("district_id = ?", districtID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAvgPropertyPrice 获取地区平均房价
func (r *DistrictRepo) GetAvgPropertyPrice(ctx context.Context, districtID uint) (float64, error) {
	var avgPrice float64
	if err := r.db.WithContext(ctx).Model(&models.Property{}).
		Where("district_id = ? AND status = ? AND listing_type = ?", districtID, "available", "sale").
		Select("AVG(price)").
		Scan(&avgPrice).Error; err != nil {
		return 0, err
	}
	return avgPrice, nil
}
