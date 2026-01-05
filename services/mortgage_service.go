package services

import (
	"context"
	"fmt"
	"math"
	"time"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

// MortgageService 按揭服务接口
type MortgageService interface {
	CalculateMortgage(ctx context.Context, req *map[string]interface{}) (*map[string]interface{}, error)
	GetMortgageRates(ctx context.Context) ([]*models.MortgageRate, error)
	GetBankMortgageRate(ctx context.Context, bankID uint) ([]*models.MortgageRate, error)
	CompareMortgageRates(ctx context.Context, req *map[string]interface{}) (*map[string]interface{}, error)
	ApplyMortgage(ctx context.Context, userID uint, req *models.MortgageApplication) (*models.MortgageApplication, error)
	GetMortgageApplications(ctx context.Context, userID uint, filter *models.ListMortgageApplicationsRequest) ([]*models.MortgageApplication, int64, error)
	GetMortgageApplication(ctx context.Context, userID uint, id uint) (*models.MortgageApplication, error)
}

type mortgageService struct {
	repo   databases.MortgageRepository
	logger *zap.Logger
}

// NewMortgageService 创建按揭服务
func NewMortgageService(repo databases.MortgageRepository, logger *zap.Logger) MortgageService {
	return &mortgageService{
		repo:   repo,
		logger: logger,
	}
}

// CalculateMortgage 计算按揭
func (s *mortgageService) CalculateMortgage(ctx context.Context, req *map[string]interface{}) (*map[string]interface{}, error) {
	// 从 map 中获取参数
	propertyPrice := (*req)["property_price"].(float64)
	downPayment := (*req)["down_payment"].(float64)
	interestRate := (*req)["interest_rate"].(float64)
	loanPeriod := int((*req)["loan_period"].(float64))
	
	// 计算贷款金额
	loanAmount := propertyPrice - downPayment
	
	// 年利率转月利率
	monthlyRate := interestRate / 100 / 12
	
	// 计算月供（使用 PMT 公式）
	// PMT = P * [r(1+r)^n] / [(1+r)^n - 1]
	// P = 贷款本金
	// r = 月利率
	// n = 还款月数
	monthlyPayment := calculateMonthlyPayment(loanAmount, monthlyRate, loanPeriod)
	
	// 计算总还款金额
	totalPayment := monthlyPayment * float64(loanPeriod)
	
	// 计算总利息
	totalInterest := totalPayment - loanAmount
	
	// 生成还款计划（仅显示前12个月和最后1个月作为示例）
	schedule := generatePaymentSchedule(loanAmount, monthlyRate, loanPeriod, monthlyPayment)
	
	ltv := (loanAmount / propertyPrice) * 100
	loanPeriodYears := loanPeriod / 12

	resp := &map[string]interface{}{
		"property_price":    propertyPrice,
		"down_payment":      downPayment,
		"loan_amount":       loanAmount,
		"interest_rate":     interestRate,
		"loan_period":       loanPeriod,
		"loan_period_years": loanPeriodYears,
		"ltv":               ltv,
		"monthly_payment":   monthlyPayment,
		"total_payment":     totalPayment,
		"total_interest":    totalInterest,
		"payment_schedule":  convertPaymentScheduleSlice(schedule),
	}
	
	return resp, nil
}

// GetMortgageRates 获取按揭利率列表
func (s *mortgageService) GetMortgageRates(ctx context.Context) ([]*models.MortgageRate, error) {
	rates, err := s.repo.GetEffectiveMortgageRates(ctx, nil)
	if err != nil {
		s.logger.Error("failed to get mortgage rates", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.MortgageRate, 0, len(rates))
	for _, rate := range rates {
		result = append(result, convertToMortgageRateResponse(rate))
	}
	
	return result, nil
}

// GetBankMortgageRate 获取指定银行的按揭利率
func (s *mortgageService) GetBankMortgageRate(ctx context.Context, bankID uint) ([]*models.MortgageRate, error) {
	// 检查银行是否存在
	bank, err := s.repo.GetBankByID(ctx, bankID)
	if err != nil {
		s.logger.Error("failed to get bank", zap.Error(err))
		return nil, err
	}
	if bank == nil {
		return nil, tools.ErrNotFound
	}
	
	rates, err := s.repo.GetMortgageRatesByBankID(ctx, bankID)
	if err != nil {
		s.logger.Error("failed to get bank mortgage rates", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.MortgageRate, 0, len(rates))
	for _, rate := range rates {
		result = append(result, convertToMortgageRateResponse(rate))
	}
	
	return result, nil
}

// CompareMortgageRates 比较按揭利率
func (s *mortgageService) CompareMortgageRates(ctx context.Context, req *map[string]interface{}) (*map[string]interface{}, error) {
	// 从 map 中获取参数
	loanAmount := (*req)["loan_amount"].(float64)
	loanPeriod := int((*req)["loan_period"].(float64))
	var rateType *string
	if rt, ok := (*req)["rate_type"].(string); ok {
		rateType = &rt
	}
	
	// 获取有效利率
	rates, err := s.repo.GetEffectiveMortgageRates(ctx, rateType)
	if err != nil {
		s.logger.Error("failed to get effective mortgage rates", zap.Error(err))
		return nil, err
	}
	
	// 计算月供比较
	comparisons := make([]map[string]interface{}, 0, len(rates))
	for _, rate := range rates {
		monthlyRate := rate.InterestRate / 100 / 12
		
		// 计算月供
		monthlyPayment := calculateMonthlyPayment(loanAmount, monthlyRate, loanPeriod)
		totalPayment := monthlyPayment * float64(loanPeriod)
		totalInterest := totalPayment - loanAmount
		
		// 计算总成本（包括手续费）
		totalCost := totalPayment
		if rate.ProcessingFee != nil {
			totalCost += *rate.ProcessingFee
		}
		
		comparison := map[string]interface{}{
			"bank":            convertToBankResponse(rate.Bank),
			"rate_type":       string(rate.RateType),
			"interest_rate":   rate.InterestRate,
			"monthly_payment": monthlyPayment,
			"total_payment":   totalPayment,
			"total_interest":  totalInterest,
			"processing_fee":  rate.ProcessingFee,
			"total_cost":      totalCost,
		}
		
		comparisons = append(comparisons, comparison)
	}
	
	resp := &map[string]interface{}{
		"loan_amount":       loanAmount,
		"loan_period":       loanPeriod,
		"rate_comparisons":  comparisons,
	}
	
	return resp, nil
}

// ApplyMortgage 申请按揭
func (s *mortgageService) ApplyMortgage(ctx context.Context, userID uint, req *models.MortgageApplication) (*models.MortgageApplication, error) {
	// 检查银行是否存在
	bank, err := s.repo.GetBankByID(ctx, req.BankID)
	if err != nil {
		s.logger.Error("failed to get bank", zap.Error(err))
		return nil, err
	}
	if bank == nil {
		return nil, tools.ErrNotFound
	}
	
	// 生成申请编号
	applicationNo := generateApplicationNo()
	
	// 计算贷款金额和月供
	loanAmount := req.LoanAmount
	monthlyRate := req.InterestRate / 100 / 12
	monthlyPayment := calculateMonthlyPayment(loanAmount, monthlyRate, req.LoanPeriod)
	
	// 计算首付
	downPayment := req.PropertyPrice - loanAmount
	// 计算贷款价值比
	ltv := loanAmount / req.PropertyPrice
	// 计算总还款
	totalPayment := monthlyPayment * float64(req.LoanPeriod)
	// 计算总利息
	totalInterest := totalPayment - loanAmount
	
	// 创建申请
	application := &models.MortgageApplication{
		ApplicationNo:       applicationNo,
		UserID:              userID,
		PropertyID:          req.PropertyID,
		BankID:              req.BankID,
		PropertyPrice:       req.PropertyPrice,
		DownPayment:         downPayment,
		LoanAmount:          loanAmount,
		InterestRate:        req.InterestRate,
		LoanPeriod:          req.LoanPeriod,
		MonthlyPayment:      monthlyPayment,
		TotalPayment:        totalPayment,
		TotalInterest:       totalInterest,
		LTV:                 ltv,
		ApplicantName:       req.ApplicantName,
		ApplicantPhone:      req.ApplicantPhone,
		ApplicantEmail:      req.ApplicantEmail,
		ApplicantIncome:     req.ApplicantIncome,
		ApplicantOccupation: req.ApplicantOccupation,
		Remarks:             req.Remarks,
		Status:              models.MortgageApplicationStatusPending,
		SubmittedAt:         time.Now(),
	}
	
	if err := s.repo.CreateApplication(ctx, application); err != nil {
		s.logger.Error("failed to create mortgage application", zap.Error(err))
		return nil, err
	}
	
	// 重新查询以获取关联数据
	application, err = s.repo.GetApplicationByID(ctx, application.ID)
	if err != nil {
		s.logger.Error("failed to get mortgage application", zap.Error(err))
		return nil, err
	}
	
	return convertToMortgageApplicationResponse(application), nil
}

// GetMortgageApplications 获取用户的按揭申请列表
func (s *mortgageService) GetMortgageApplications(ctx context.Context, userID uint, filter *models.ListMortgageApplicationsRequest) ([]*models.MortgageApplication, int64, error) {
	applications, total, err := s.repo.GetApplicationsByUserID(ctx, userID, filter)
	if err != nil {
		s.logger.Error("failed to get mortgage applications", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.MortgageApplication, 0, len(applications))
	for _, app := range applications {
		result = append(result, convertToMortgageApplicationResponse(app))
	}
	
	return result, total, nil
}

// GetMortgageApplication 获取按揭申请详情
func (s *mortgageService) GetMortgageApplication(ctx context.Context, userID uint, id uint) (*models.MortgageApplication, error) {
	application, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get mortgage application", zap.Error(err))
		return nil, err
	}
	if application == nil {
		return nil, tools.ErrNotFound
	}
	
	// 检查所有权
	if application.UserID != userID {
		return nil, tools.ErrForbidden
	}
	
	return convertToMortgageApplicationResponse(application), nil
}

// 辅助函数

// calculateMonthlyPayment 计算月供（PMT 公式）
func calculateMonthlyPayment(principal float64, monthlyRate float64, periods int) float64 {
	if monthlyRate == 0 {
		return principal / float64(periods)
	}
	
	// PMT = P * [r(1+r)^n] / [(1+r)^n - 1]
	pow := math.Pow(1+monthlyRate, float64(periods))
	return principal * (monthlyRate * pow) / (pow - 1)
}

// generatePaymentSchedule 生成还款计划（仅显示前12个月和最后1个月）
func generatePaymentSchedule(principal float64, monthlyRate float64, periods int, monthlyPayment float64) []*map[string]interface{} {
	schedule := make([]*map[string]interface{}, 0)
	
	remainingBalance := principal
	
	// 显示前12个月
	monthsToShow := 12
	if periods < 12 {
		monthsToShow = periods
	}
	
	for month := 1; month <= monthsToShow; month++ {
		interestPayment := remainingBalance * monthlyRate
		principalPayment := monthlyPayment - interestPayment
		remainingBalance -= principalPayment
		
		item := &map[string]interface{}{
			"period":             month,
			"payment":            monthlyPayment,
			"principal":          principalPayment,
			"interest":           interestPayment,
			"remaining_balance":  remainingBalance,
		}
		
		schedule = append(schedule, item)
	}
	
	// 如果还款期超过12个月，添加最后一个月
	if periods > 12 {
		// 快速计算到最后一个月的余额
		remainingBalance = principal
		for month := 1; month < periods; month++ {
			interestPayment := remainingBalance * monthlyRate
			principalPayment := monthlyPayment - interestPayment
			remainingBalance -= principalPayment
		}
		
		// 最后一个月
		interestPayment := remainingBalance * monthlyRate
		principalPayment := monthlyPayment - interestPayment
		
		item := &map[string]interface{}{
			"period":             periods,
			"payment":            monthlyPayment,
			"principal":          principalPayment,
			"interest":           interestPayment,
			"remaining_balance":  0,
		}
		
		schedule = append(schedule, item)
	}
	
	return schedule
}

// generateApplicationNo 生成申请编号
func generateApplicationNo() string {
	now := time.Now()
	return fmt.Sprintf("MTG%s%06d", now.Format("20060102"), now.Unix()%1000000)
}

// convertToMortgageRateResponse 转换为按揭利率响应（直接返回，预加载了关联数据）
func convertToMortgageRateResponse(rate *models.MortgageRate) *models.MortgageRate {
	return rate
}

// convertToMortgageApplicationResponse 转换为按揭申请响应（直接返回，预加载了关联数据）
func convertToMortgageApplicationResponse(app *models.MortgageApplication) *models.MortgageApplication {
	return app
}


// convertPaymentScheduleSlice 转换还款计划切片
func convertPaymentScheduleSlice(schedules []*map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(schedules))
	for _, s := range schedules {
		if s != nil {
			result = append(result, *s)
		}
	}
	return result
}

// convertToBankResponse 转换银行模型为响应
func convertToBankResponse(bank *models.Bank) *map[string]interface{} {
	if bank == nil {
		return nil
	}
	return &map[string]interface{}{
		"id":           bank.ID,
		"name_zh_hant": bank.NameZhHant,
		"name_zh_hans": bank.NameZhHans,
		"name_en":      bank.NameEn,
		"logo":         bank.Logo,
		"website":      bank.Website,
	}
}
