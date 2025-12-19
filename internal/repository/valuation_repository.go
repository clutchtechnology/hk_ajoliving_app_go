package repository

import (
	"context"
	"fmt"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"gorm.io/gorm"
)

// ValuationRepository 物业估价数据仓库接口
type ValuationRepository interface {
	ListValuations(ctx context.Context, req *request.ListValuationsRequest) ([]*model.Estate, int64, error)
	GetEstateValuation(ctx context.Context, estateID uint) (*model.Estate, error)
	SearchValuations(ctx context.Context, req *request.SearchValuationsRequest) ([]*model.Estate, int64, error)
	GetDistrictValuations(ctx context.Context, districtID uint, page, pageSize int) ([]*model.Estate, int64, error)
	GetDistrictStatistics(ctx context.Context, districtID uint) (map[string]interface{}, error)
}

type valuationRepository struct {
	db *gorm.DB
}

func NewValuationRepository(db *gorm.DB) ValuationRepository {
	return &valuationRepository{db: db}
}

func (r *valuationRepository) ListValuations(ctx context.Context, req *request.ListValuationsRequest) ([]*model.Estate, int64, error) {
	var estates []*model.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Estate{}).
		Where("avg_transaction_price IS NOT NULL")

	// 筛选条件
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.MinPrice != nil {
		query = query.Where("avg_transaction_price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("avg_transaction_price <= ?", *req.MaxPrice)
	}
	if req.SchoolNet != "" {
		query = query.Where("primary_school_net = ? OR secondary_school_net = ?", req.SchoolNet, req.SchoolNet)
	}

	// 统计
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页排序
	offset := (req.Page - 1) * req.PageSize
	sortColumn := req.SortBy
	if sortColumn == "" {
		sortColumn = "avg_transaction_price"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	query = query.Offset(offset).Limit(req.PageSize).Order(sortColumn + " " + sortOrder)
	query = query.Preload("District")

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

func (r *valuationRepository) GetEstateValuation(ctx context.Context, estateID uint) (*model.Estate, error) {
	var estate model.Estate
	err := r.db.WithContext(ctx).
		Preload("District").
		First(&estate, estateID).Error
	if err != nil {
		return nil, err
	}
	return &estate, nil
}

func (r *valuationRepository) SearchValuations(ctx context.Context, req *request.SearchValuationsRequest) ([]*model.Estate, int64, error) {
	var estates []*model.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Estate{}).
		Where("avg_transaction_price IS NOT NULL")

	// 搜索关键词（屋苑名称或地址）
	if req.Keyword != "" {
		searchPattern := fmt.Sprintf("%%%s%%", req.Keyword)
		query = query.Where("name LIKE ? OR name_en LIKE ? OR address LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// 筛选条件
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.MinPrice != nil {
		query = query.Where("avg_transaction_price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("avg_transaction_price <= ?", *req.MaxPrice)
	}

	// 统计
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize).
		Order("avg_transaction_price DESC").
		Preload("District")

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

func (r *valuationRepository) GetDistrictValuations(ctx context.Context, districtID uint, page, pageSize int) ([]*model.Estate, int64, error) {
	var estates []*model.Estate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Estate{}).
		Where("district_id = ?", districtID).
		Where("avg_transaction_price IS NOT NULL")

	// 统计
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).
		Order("avg_transaction_price DESC").
		Preload("District")

	if err := query.Find(&estates).Error; err != nil {
		return nil, 0, err
	}

	return estates, total, nil
}

func (r *valuationRepository) GetDistrictStatistics(ctx context.Context, districtID uint) (map[string]interface{}, error) {
	var result struct {
		TotalEstates       int64
		AvgPrice           float64
		MinPrice           float64
		MaxPrice           float64
		TotalTransactions  int64
	}

	err := r.db.WithContext(ctx).Model(&model.Estate{}).
		Select(`
			COUNT(*) as total_estates,
			AVG(avg_transaction_price) as avg_price,
			MIN(avg_transaction_price) as min_price,
			MAX(avg_transaction_price) as max_price,
			SUM(recent_transactions_count) as total_transactions
		`).
		Where("district_id = ?", districtID).
		Where("avg_transaction_price IS NOT NULL").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	statistics := map[string]interface{}{
		"total_estates":      result.TotalEstates,
		"avg_price":          result.AvgPrice,
		"min_price":          result.MinPrice,
		"max_price":          result.MaxPrice,
		"total_transactions": result.TotalTransactions,
	}

	return statistics, nil
}
