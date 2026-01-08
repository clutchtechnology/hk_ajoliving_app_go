package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

// ValuationRepo 估价仓储
type ValuationRepo struct {
	db *gorm.DB
}

// NewValuationRepo 创建估价仓储
func NewValuationRepo(db *gorm.DB) *ValuationRepo {
	return &ValuationRepo{db: db}
}

// FindAllValuations 查询屋苑估价列表
func (r *ValuationRepo) FindAllValuations(ctx context.Context, filter *models.ListValuationsRequest) ([]models.ValuationResponse, int64, error) {
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

	if filter.MinAvgPrice != nil {
		query = query.Where("avg_transaction_price >= ?", *filter.MinAvgPrice)
	}

	if filter.MaxAvgPrice != nil {
		query = query.Where("avg_transaction_price <= ?", *filter.MaxAvgPrice)
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
	sortBy := "avg_transaction_price"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "name":
			sortBy = "name"
		case "price":
			sortBy = "avg_transaction_price"
		case "yield":
			sortBy = "avg_transaction_price" // TODO: 增加租金回报率字段后修改
		case "transactions":
			sortBy = "recent_transactions_count"
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
	query = query.Preload("District")

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	// 转换为估价响应
	valuations := make([]models.ValuationResponse, len(estates))
	for i, estate := range estates {
		valuations[i] = r.estateToValuation(ctx, &estate)
	}

	return valuations, total, nil
}

// GetEstateValuation 获取指定屋苑估价详情
func (r *ValuationRepo) GetEstateValuation(ctx context.Context, estateID uint) (*models.EstateValuationDetail, error) {
	var estate models.Estate
	if err := r.db.WithContext(ctx).
		Preload("District").
		First(&estate, estateID).Error; err != nil {
		return nil, err
	}

	// 基础估价信息
	valuation := &models.EstateValuationDetail{
		EstateID:               estate.ID,
		EstateName:             estate.Name,
		EstateNameEn:           estate.NameEn,
		Address:                estate.Address,
		District:               estate.District,
		CompletionYear:         estate.CompletionYear,
		Developer:              estate.Developer,
		TotalBlocks:            estate.TotalBlocks,
		TotalUnits:             estate.TotalUnits,
		PrimarySchoolNet:       estate.PrimarySchoolNet,
		SecondarySchoolNet:     estate.SecondarySchoolNet,
		AvgPricePerSqft:        estate.AvgTransactionPrice,
		RecentTransactionCount: estate.RecentTransactionsCount,
		ForSaleCount:           estate.ForSaleCount,
		ForRentCount:           estate.ForRentCount,
		LastUpdated:            estate.UpdatedAt,
	}

	// 查询户型价格分布
	unitTypePrices, err := r.getUnitTypePrices(ctx, estate.Name)
	if err == nil {
		valuation.UnitTypePrices = unitTypePrices
	}

	// 计算平均售价和租金
	avgPrices := r.getAvgPrices(ctx, estate.Name)
	valuation.AvgSalePrice = avgPrices["sale"]
	valuation.AvgRentPrice = avgPrices["rent"]

	// 计算租金回报率
	if valuation.AvgSalePrice > 0 && valuation.AvgRentPrice > 0 {
		valuation.RentalYield = (valuation.AvgRentPrice * 12 / valuation.AvgSalePrice) * 100
	}

	// 获取价格范围
	priceRange := r.getPriceRange(ctx, estate.Name)
	valuation.MinPricePerSqft = priceRange["min"]
	valuation.MaxPricePerSqft = priceRange["max"]

	// 获取价格历史（最近12个月）
	// TODO: 实现价格历史查询
	valuation.PriceHistory = []models.PriceHistoryPoint{}

	// 获取近期成交
	// TODO: 实现成交记录查询
	valuation.RecentTransactions = []models.TransactionSummary{}

	return valuation, nil
}

// SearchValuations 搜索屋苑估价
func (r *ValuationRepo) SearchValuations(ctx context.Context, keyword string, page, pageSize int) ([]models.ValuationResponse, int64, error) {
	var estates []models.Estate
	var total int64

	searchKeyword := "%" + keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.Estate{}).
		Where("name LIKE ? OR name_en LIKE ? OR address LIKE ?", searchKeyword, searchKeyword, searchKeyword)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// 预加载关联
	query = query.Preload("District").Order("view_count DESC")

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	// 转换为估价响应
	valuations := make([]models.ValuationResponse, len(estates))
	for i, estate := range estates {
		valuations[i] = r.estateToValuation(ctx, &estate)
	}

	return valuations, total, nil
}

// GetDistrictValuations 获取地区屋苑估价列表
func (r *ValuationRepo) GetDistrictValuations(ctx context.Context, districtID uint) (*models.DistrictValuationSummary, error) {
	var district models.District
	if err := r.db.WithContext(ctx).First(&district, districtID).Error; err != nil {
		return nil, err
	}

	var estates []models.Estate
	if err := r.db.WithContext(ctx).
		Where("district_id = ?", districtID).
		Order("avg_transaction_price DESC").
		Find(&estates).Error; err != nil {
		return nil, err
	}

	// 计算地区汇总数据
	summary := &models.DistrictValuationSummary{
		DistrictID:  districtID,
		District:    &district,
		EstateCount: len(estates),
	}

	if len(estates) > 0 {
		var sumPrice, minPrice, maxPrice float64
		var totalTransactions int

		minPrice = estates[0].AvgTransactionPrice
		maxPrice = estates[0].AvgTransactionPrice

		for _, estate := range estates {
			sumPrice += estate.AvgTransactionPrice
			totalTransactions += estate.RecentTransactionsCount

			if estate.AvgTransactionPrice < minPrice && estate.AvgTransactionPrice > 0 {
				minPrice = estate.AvgTransactionPrice
			}
			if estate.AvgTransactionPrice > maxPrice {
				maxPrice = estate.AvgTransactionPrice
			}
		}

		summary.AvgPricePerSqft = sumPrice / float64(len(estates))
		summary.MinPricePerSqft = minPrice
		summary.MaxPricePerSqft = maxPrice
		summary.TotalTransactions = totalTransactions
	}

	// 转换屋苑列表为估价响应
	summary.Estates = make([]models.ValuationResponse, len(estates))
	for i, estate := range estates {
		summary.Estates[i] = r.estateToValuation(ctx, &estate)
	}

	return summary, nil
}

// estateToValuation 转换屋苑为估价响应
func (r *ValuationRepo) estateToValuation(ctx context.Context, estate *models.Estate) models.ValuationResponse {
	valuation := models.ValuationResponse{
		EstateID:               estate.ID,
		EstateName:             estate.Name,
		EstateNameEn:           estate.NameEn,
		DistrictID:             estate.DistrictID,
		District:               estate.District,
		Address:                estate.Address,
		CompletionYear:         estate.CompletionYear,
		TotalUnits:             estate.TotalUnits,
		AvgPricePerSqft:        estate.AvgTransactionPrice,
		RecentTransactionCount: estate.RecentTransactionsCount,
		ForSaleCount:           estate.ForSaleCount,
		ForRentCount:           estate.ForRentCount,
		LastUpdated:            estate.UpdatedAt,
	}

	// 获取平均价格
	avgPrices := r.getAvgPrices(ctx, estate.Name)
	valuation.AvgSalePrice = avgPrices["sale"]
	valuation.AvgRentPrice = avgPrices["rent"]

	// 计算租金回报率
	if valuation.AvgSalePrice > 0 && valuation.AvgRentPrice > 0 {
		valuation.RentalYield = (valuation.AvgRentPrice * 12 / valuation.AvgSalePrice) * 100
	}

	// 获取价格范围
	priceRange := r.getPriceRange(ctx, estate.Name)
	valuation.MinPricePerSqft = priceRange["min"]
	valuation.MaxPricePerSqft = priceRange["max"]

	// TODO: 计算价格变化趋势（需要历史数据）
	valuation.PriceChange30d = 0
	valuation.PriceChange90d = 0

	return valuation
}

// getUnitTypePrices 获取户型价格分布
func (r *ValuationRepo) getUnitTypePrices(ctx context.Context, estateName string) ([]models.UnitTypePriceBreakdown, error) {
	var results []models.UnitTypePriceBreakdown

	err := r.db.WithContext(ctx).
		Model(&models.Property{}).
		Select(`
			bedrooms,
			AVG(area) as avg_area,
			AVG(price) as avg_price,
			AVG(price / area) as avg_price_per_sqft,
			MIN(price) as min_price,
			MAX(price) as max_price,
			COUNT(*) as available_count
		`).
		Where("building_name = ? AND status = ?", estateName, "available").
		Group("bedrooms").
		Order("bedrooms").
		Scan(&results).Error

	return results, err
}

// getAvgPrices 获取平均售价和租金
func (r *ValuationRepo) getAvgPrices(ctx context.Context, estateName string) map[string]float64 {
	result := make(map[string]float64)

	// 平均售价
	var avgSalePrice float64
	r.db.WithContext(ctx).
		Model(&models.Property{}).
		Select("AVG(price) as avg_sale_price").
		Where("building_name = ? AND listing_type = ? AND status = ?", estateName, "sale", "available").
		Scan(&avgSalePrice)
	result["sale"] = avgSalePrice

	// 平均租金
	var avgRentPrice float64
	r.db.WithContext(ctx).
		Model(&models.Property{}).
		Select("AVG(price) as avg_rent_price").
		Where("building_name = ? AND listing_type = ? AND status = ?", estateName, "rent", "available").
		Scan(&avgRentPrice)
	result["rent"] = avgRentPrice

	return result
}

// getPriceRange 获取价格范围
func (r *ValuationRepo) getPriceRange(ctx context.Context, estateName string) map[string]float64 {
	result := make(map[string]float64)

	var priceRange struct {
		MinPrice float64
		MaxPrice float64
	}

	r.db.WithContext(ctx).
		Model(&models.Property{}).
		Select("MIN(price / area) as min_price, MAX(price / area) as max_price").
		Where("building_name = ? AND status = ?", estateName, "available").
		Scan(&priceRange)

	result["min"] = priceRange.MinPrice
	result["max"] = priceRange.MaxPrice

	return result
}
