package databases

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// MortgageRepository 按揭数据访问接口
type MortgageRepository interface {
	// 银行相关
	GetAllBanks(ctx context.Context) ([]*models.Bank, error)
	GetBankByID(ctx context.Context, id uint) (*models.Bank, error)
	
	// 利率相关
	GetAllMortgageRates(ctx context.Context) ([]*models.MortgageRate, error)
	GetMortgageRatesByBankID(ctx context.Context, bankID uint) ([]*models.MortgageRate, error)
	GetEffectiveMortgageRates(ctx context.Context, rateType *string) ([]*models.MortgageRate, error)
	GetMortgageRateByID(ctx context.Context, id uint) (*models.MortgageRate, error)
	
	// 申请相关
	CreateApplication(ctx context.Context, application *models.MortgageApplication) error
	GetApplicationByID(ctx context.Context, id uint) (*models.MortgageApplication, error)
	GetApplicationByNo(ctx context.Context, applicationNo string) (*models.MortgageApplication, error)
	GetApplicationsByUserID(ctx context.Context, userID uint, filter *models.ListMortgageApplicationsRequest) ([]*models.MortgageApplication, int64, error)
	UpdateApplication(ctx context.Context, application *models.MortgageApplication) error
}

type mortgageRepository struct {
	db *gorm.DB
}

// NewMortgageRepository 创建按揭仓库
func NewMortgageRepository(db *gorm.DB) MortgageRepository {
	return &mortgageRepository{db: db}
}

// 银行相关

func (r *mortgageRepository) GetAllBanks(ctx context.Context) ([]*models.Bank, error) {
	var banks []*models.Bank
	
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, name_zh_hant ASC").
		Find(&banks).Error
		
	if err != nil {
		return nil, err
	}
	
	return banks, nil
}

func (r *mortgageRepository) GetBankByID(ctx context.Context, id uint) (*models.Bank, error) {
	var bank models.Bank
	
	err := r.db.WithContext(ctx).First(&bank, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &bank, nil
}

// 利率相关

func (r *mortgageRepository) GetAllMortgageRates(ctx context.Context) ([]*models.MortgageRate, error) {
	var rates []*models.MortgageRate
	
	err := r.db.WithContext(ctx).
		Preload("Bank").
		Where("is_active = ?", true).
		Order("interest_rate ASC").
		Find(&rates).Error
		
	if err != nil {
		return nil, err
	}
	
	return rates, nil
}

func (r *mortgageRepository) GetMortgageRatesByBankID(ctx context.Context, bankID uint) ([]*models.MortgageRate, error) {
	var rates []*models.MortgageRate
	
	err := r.db.WithContext(ctx).
		Preload("Bank").
		Where("bank_id = ? AND is_active = ?", bankID, true).
		Order("interest_rate ASC").
		Find(&rates).Error
		
	if err != nil {
		return nil, err
	}
	
	return rates, nil
}

func (r *mortgageRepository) GetEffectiveMortgageRates(ctx context.Context, rateType *string) ([]*models.MortgageRate, error) {
	var rates []*models.MortgageRate
	
	query := r.db.WithContext(ctx).
		Preload("Bank").
		Where("is_active = ?", true).
		Where("effective_date <= NOW()").
		Where("(expiry_date IS NULL OR expiry_date > NOW())")
	
	if rateType != nil && *rateType != "" {
		query = query.Where("rate_type = ?", *rateType)
	}
	
	err := query.Order("interest_rate ASC").Find(&rates).Error
	if err != nil {
		return nil, err
	}
	
	return rates, nil
}

func (r *mortgageRepository) GetMortgageRateByID(ctx context.Context, id uint) (*models.MortgageRate, error) {
	var rate models.MortgageRate
	
	err := r.db.WithContext(ctx).
		Preload("Bank").
		First(&rate, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &rate, nil
}

// 申请相关

func (r *mortgageRepository) CreateApplication(ctx context.Context, application *models.MortgageApplication) error {
	return r.db.WithContext(ctx).Create(application).Error
}

func (r *mortgageRepository) GetApplicationByID(ctx context.Context, id uint) (*models.MortgageApplication, error) {
	var application models.MortgageApplication
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Property").
		Preload("Bank").
		First(&application, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &application, nil
}

func (r *mortgageRepository) GetApplicationByNo(ctx context.Context, applicationNo string) (*models.MortgageApplication, error) {
	var application models.MortgageApplication
	
	err := r.db.WithContext(ctx).
		Where("application_no = ?", applicationNo).
		First(&application).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &application, nil
}

func (r *mortgageRepository) GetApplicationsByUserID(ctx context.Context, userID uint, filter *models.ListMortgageApplicationsRequest) ([]*models.MortgageApplication, int64, error) {
	var applications []*models.MortgageApplication
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.MortgageApplication{}).Where("user_id = ?", userID)
	
	// 应用筛选条件
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.BankID != nil {
		query = query.Where("bank_id = ?", *filter.BankID)
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
	query = query.Preload("Bank").Preload("Property")
	
	if err := query.Find(&applications).Error; err != nil {
		return nil, 0, err
	}
	
	return applications, total, nil
}

func (r *mortgageRepository) UpdateApplication(ctx context.Context, application *models.MortgageApplication) error {
	return r.db.WithContext(ctx).Save(application).Error
}
