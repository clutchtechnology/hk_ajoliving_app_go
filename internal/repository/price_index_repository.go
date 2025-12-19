package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// PriceIndexRepository 楼价指数数据访问接口
type PriceIndexRepository interface {
	// 查询
	List(ctx context.Context, filter *request.GetPriceIndexRequest) ([]*model.PriceIndex, int64, error)
	GetByID(ctx context.Context, id uint) (*model.PriceIndex, error)
	GetLatest(ctx context.Context, indexType string) (*model.PriceIndex, error)
	GetLatestByDistrict(ctx context.Context, districtID uint) (*model.PriceIndex, error)
	GetLatestByEstate(ctx context.Context, estateID uint) (*model.PriceIndex, error)
	GetDistrictPriceIndex(ctx context.Context, districtID uint, startPeriod, endPeriod *string, limit int) ([]*model.PriceIndex, error)
	GetEstatePriceIndex(ctx context.Context, estateID uint, startPeriod, endPeriod *string, limit int) ([]*model.PriceIndex, error)
	GetTrends(ctx context.Context, filter *request.GetPriceTrendsRequest) ([]*model.PriceIndex, error)
	GetForComparison(ctx context.Context, indexType string, ids []uint, startPeriod, endPeriod string) ([]*model.PriceIndex, error)
	GetHistory(ctx context.Context, indexType string, districtID, estateID *uint, propertyType *string, years int) ([]*model.PriceIndex, error)
	GetAllLatestByType(ctx context.Context, indexType string) ([]*model.PriceIndex, error)
	
	// 创建和更新
	Create(ctx context.Context, index *model.PriceIndex) error
	Update(ctx context.Context, index *model.PriceIndex) error
	Delete(ctx context.Context, id uint) error
}

type priceIndexRepository struct {
	db *gorm.DB
}

// NewPriceIndexRepository 创建楼价指数仓库
func NewPriceIndexRepository(db *gorm.DB) PriceIndexRepository {
	return &priceIndexRepository{db: db}
}

// List 获取楼价指数列表
func (r *priceIndexRepository) List(ctx context.Context, filter *request.GetPriceIndexRequest) ([]*model.PriceIndex, int64, error) {
	var indices []*model.PriceIndex
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.PriceIndex{})
	
	// 应用筛选条件
	if filter.IndexType != nil {
		query = query.Where("index_type = ?", *filter.IndexType)
	}
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.EstateID != nil {
		query = query.Where("estate_id = ?", *filter.EstateID)
	}
	if filter.PropertyType != nil {
		query = query.Where("property_type = ?", *filter.PropertyType)
	}
	if filter.StartPeriod != nil {
		query = query.Where("period >= ?", *filter.StartPeriod)
	}
	if filter.EndPeriod != nil {
		query = query.Where("period <= ?", *filter.EndPeriod)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)
	
	// 排序：最新的在前
	query = query.Order("year DESC, month DESC, day DESC")
	
	// 预加载关联
	query = query.Preload("District").Preload("Estate")
	
	if err := query.Find(&indices).Error; err != nil {
		return nil, 0, err
	}
	
	return indices, total, nil
}

// GetByID 根据ID获取楼价指数
func (r *priceIndexRepository) GetByID(ctx context.Context, id uint) (*model.PriceIndex, error) {
	var index model.PriceIndex
	
	err := r.db.WithContext(ctx).
		Preload("District").
		Preload("Estate").
		First(&index, id).Error
	
	if err != nil {
		return nil, err
	}
	
	return &index, nil
}

// GetLatest 获取最新的指定类型指数
func (r *priceIndexRepository) GetLatest(ctx context.Context, indexType string) (*model.PriceIndex, error) {
	var index model.PriceIndex
	
	err := r.db.WithContext(ctx).
		Where("index_type = ?", indexType).
		Order("year DESC, month DESC, day DESC").
		First(&index).Error
	
	if err != nil {
		return nil, err
	}
	
	return &index, nil
}

// GetLatestByDistrict 获取指定地区的最新指数
func (r *priceIndexRepository) GetLatestByDistrict(ctx context.Context, districtID uint) (*model.PriceIndex, error) {
	var index model.PriceIndex
	
	err := r.db.WithContext(ctx).
		Where("district_id = ?", districtID).
		Where("index_type = ?", model.IndexTypeDistrict).
		Order("year DESC, month DESC, day DESC").
		Preload("District").
		First(&index).Error
	
	if err != nil {
		return nil, err
	}
	
	return &index, nil
}

// GetLatestByEstate 获取指定屋苑的最新指数
func (r *priceIndexRepository) GetLatestByEstate(ctx context.Context, estateID uint) (*model.PriceIndex, error) {
	var index model.PriceIndex
	
	err := r.db.WithContext(ctx).
		Where("estate_id = ?", estateID).
		Where("index_type = ?", model.IndexTypeEstate).
		Order("year DESC, month DESC, day DESC").
		Preload("Estate").
		First(&index).Error
	
	if err != nil {
		return nil, err
	}
	
	return &index, nil
}

// GetDistrictPriceIndex 获取地区楼价指数
func (r *priceIndexRepository) GetDistrictPriceIndex(ctx context.Context, districtID uint, startPeriod, endPeriod *string, limit int) ([]*model.PriceIndex, error) {
	var indices []*model.PriceIndex
	
	query := r.db.WithContext(ctx).
		Where("district_id = ?", districtID).
		Where("index_type = ?", model.IndexTypeDistrict)
	
	if startPeriod != nil {
		query = query.Where("period >= ?", *startPeriod)
	}
	if endPeriod != nil {
		query = query.Where("period <= ?", *endPeriod)
	}
	
	query = query.Order("year DESC, month DESC, day DESC").
		Limit(limit).
		Preload("District")
	
	if err := query.Find(&indices).Error; err != nil {
		return nil, err
	}
	
	return indices, nil
}

// GetEstatePriceIndex 获取屋苑楼价指数
func (r *priceIndexRepository) GetEstatePriceIndex(ctx context.Context, estateID uint, startPeriod, endPeriod *string, limit int) ([]*model.PriceIndex, error) {
	var indices []*model.PriceIndex
	
	query := r.db.WithContext(ctx).
		Where("estate_id = ?", estateID).
		Where("index_type = ?", model.IndexTypeEstate)
	
	if startPeriod != nil {
		query = query.Where("period >= ?", *startPeriod)
	}
	if endPeriod != nil {
		query = query.Where("period <= ?", *endPeriod)
	}
	
	query = query.Order("year DESC, month DESC, day DESC").
		Limit(limit).
		Preload("Estate")
	
	if err := query.Find(&indices).Error; err != nil {
		return nil, err
	}
	
	return indices, nil
}

// GetTrends 获取价格走势
func (r *priceIndexRepository) GetTrends(ctx context.Context, filter *request.GetPriceTrendsRequest) ([]*model.PriceIndex, error) {
	var indices []*model.PriceIndex
	
	query := r.db.WithContext(ctx).
		Where("index_type = ?", filter.IndexType).
		Where("period >= ?", filter.StartPeriod).
		Where("period <= ?", filter.EndPeriod)
	
	if filter.DistrictID != nil {
		query = query.Where("district_id = ?", *filter.DistrictID)
	}
	if filter.EstateID != nil {
		query = query.Where("estate_id = ?", *filter.EstateID)
	}
	if filter.PropertyType != nil {
		query = query.Where("property_type = ?", *filter.PropertyType)
	}
	
	query = query.Order("year ASC, month ASC, day ASC").
		Preload("District").
		Preload("Estate")
	
	if err := query.Find(&indices).Error; err != nil {
		return nil, err
	}
	
	return indices, nil
}

// GetForComparison 获取用于对比的数据
func (r *priceIndexRepository) GetForComparison(ctx context.Context, indexType string, ids []uint, startPeriod, endPeriod string) ([]*model.PriceIndex, error) {
	var indices []*model.PriceIndex
	
	query := r.db.WithContext(ctx).
		Where("index_type = ?", indexType).
		Where("period >= ?", startPeriod).
		Where("period <= ?", endPeriod)
	
	switch indexType {
	case string(model.IndexTypeDistrict):
		query = query.Where("district_id IN ?", ids)
	case string(model.IndexTypeEstate):
		query = query.Where("estate_id IN ?", ids)
	}
	
	query = query.Order("year ASC, month ASC, day ASC").
		Preload("District").
		Preload("Estate")
	
	if err := query.Find(&indices).Error; err != nil {
		return nil, err
	}
	
	return indices, nil
}

// GetHistory 获取历史楼价指数
func (r *priceIndexRepository) GetHistory(ctx context.Context, indexType string, districtID, estateID *uint, propertyType *string, years int) ([]*model.PriceIndex, error) {
	var indices []*model.PriceIndex
	
	// 计算起始年份
	currentYear := time.Now().Year()
	startYear := currentYear - years
	
	query := r.db.WithContext(ctx).
		Where("index_type = ?", indexType).
		Where("year >= ?", startYear)
	
	if districtID != nil {
		query = query.Where("district_id = ?", *districtID)
	}
	if estateID != nil {
		query = query.Where("estate_id = ?", *estateID)
	}
	if propertyType != nil {
		query = query.Where("property_type = ?", *propertyType)
	}
	
	query = query.Order("year ASC, month ASC, day ASC").
		Preload("District").
		Preload("Estate")
	
	if err := query.Find(&indices).Error; err != nil {
		return nil, err
	}
	
	return indices, nil
}

// GetAllLatestByType 获取指定类型的所有最新指数
func (r *priceIndexRepository) GetAllLatestByType(ctx context.Context, indexType string) ([]*model.PriceIndex, error) {
	var indices []*model.PriceIndex
	
	// 先找出最新的期数
	var latestPeriod string
	err := r.db.WithContext(ctx).
		Model(&model.PriceIndex{}).
		Where("index_type = ?", indexType).
		Order("year DESC, month DESC, day DESC").
		Limit(1).
		Pluck("period", &latestPeriod).Error
	
	if err != nil {
		return nil, err
	}
	
	// 获取该期数的所有指数
	err = r.db.WithContext(ctx).
		Where("index_type = ?", indexType).
		Where("period = ?", latestPeriod).
		Preload("District").
		Preload("Estate").
		Find(&indices).Error
	
	if err != nil {
		return nil, err
	}
	
	return indices, nil
}

// Create 创建楼价指数
func (r *priceIndexRepository) Create(ctx context.Context, index *model.PriceIndex) error {
	return r.db.WithContext(ctx).Create(index).Error
}

// Update 更新楼价指数
func (r *priceIndexRepository) Update(ctx context.Context, index *model.PriceIndex) error {
	return r.db.WithContext(ctx).Save(index).Error
}

// Delete 删除楼价指数
func (r *priceIndexRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.PriceIndex{}, id).Error
}
