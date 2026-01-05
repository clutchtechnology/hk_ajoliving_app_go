package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

// MortgageService 按揭服务接口
type MortgageService interface {
	CalculateMortgage(ctx context.Context, req *request.CalculateMortgageRequest) (*response.MortgageCalculationResponse, error)
	GetMortgageRates(ctx context.Context) ([]*response.MortgageRateResponse, error)
	GetBankMortgageRate(ctx context.Context, bankID uint) ([]*response.MortgageRateResponse, error)
	CompareMortgageRates(ctx context.Context, req *request.CompareMortgageRatesRequest) (*response.MortgageRateComparisonResponse, error)
	ApplyMortgage(ctx context.Context, userID uint, req *request.ApplyMortgageRequest) (*response.MortgageApplicationResponse, error)
	GetMortgageApplications(ctx context.Context, userID uint, filter *request.ListMortgageApplicationsRequest) ([]*response.MortgageApplicationResponse, int64, error)
	GetMortgageApplication(ctx context.Context, userID uint, id uint) (*response.MortgageApplicationResponse, error)
}

type mortgageService struct {
	repo   repository.MortgageRepository
	logger *zap.Logger
}

// NewMortgageService 创建按揭服务
func NewMortgageService(repo repository.MortgageRepository, logger *zap.Logger) MortgageService {
	return &mortgageService{
		repo:   repo,
		logger: logger,
	}
}

// CalculateMortgage 计算按揭
func (s *mortgageService) CalculateMortgage(ctx context.Context, req *request.CalculateMortgageRequest) (*response.MortgageCalculationResponse, error) {
	// 计算贷款金额
	loanAmount := req.PropertyPrice - req.DownPayment
	
	// 年利率转月利率
	monthlyRate := req.InterestRate / 100 / 12
	
	// 计算月供（使用 PMT 公式）
	// PMT = P * [r(1+r)^n] / [(1+r)^n - 1]
	// P = 贷款本金
	// r = 月利率
	// n = 还款月数
	monthlyPayment := calculateMonthlyPayment(loanAmount, monthlyRate, req.LoanPeriod)
	
	// 计算总还款金额
	totalPayment := monthlyPayment * float64(req.LoanPeriod)
	
	// 计算总利息
	totalInterest := totalPayment - loanAmount
	
	// 生成还款计划（仅显示前12个月和最后1个月作为示例）
	schedule := generatePaymentSchedule(loanAmount, monthlyRate, req.LoanPeriod, monthlyPayment)
	
	ltv := (loanAmount / req.PropertyPrice) * 100
	loanPeriodYears := req.LoanPeriod / 12

	resp := &response.MortgageCalculationResponse{
		PropertyPrice:   req.PropertyPrice,
		DownPayment:     req.DownPayment,
		LoanAmount:      loanAmount,
		InterestRate:    req.InterestRate,
		LoanPeriod:      req.LoanPeriod,
		LoanPeriodYears: loanPeriodYears,
		LTV:             ltv,
		MonthlyPayment:  monthlyPayment,
		TotalPayment:    totalPayment,
		TotalInterest:   totalInterest,
		PaymentSchedule: convertPaymentScheduleSlice(schedule),
	}
	
	return resp, nil
}

// GetMortgageRates 获取按揭利率列表
func (s *mortgageService) GetMortgageRates(ctx context.Context) ([]*response.MortgageRateResponse, error) {
	rates, err := s.repo.GetEffectiveMortgageRates(ctx, nil)
	if err != nil {
		s.logger.Error("failed to get mortgage rates", zap.Error(err))
		return nil, err
	}
	
	result := make([]*response.MortgageRateResponse, 0, len(rates))
	for _, rate := range rates {
		result = append(result, convertToMortgageRateResponse(rate))
	}
	
	return result, nil
}

// GetBankMortgageRate 获取指定银行的按揭利率
func (s *mortgageService) GetBankMortgageRate(ctx context.Context, bankID uint) ([]*response.MortgageRateResponse, error) {
	// 检查银行是否存在
	bank, err := s.repo.GetBankByID(ctx, bankID)
	if err != nil {
		s.logger.Error("failed to get bank", zap.Error(err))
		return nil, err
	}
	if bank == nil {
		return nil, errors.ErrNotFound
	}
	
	rates, err := s.repo.GetMortgageRatesByBankID(ctx, bankID)
	if err != nil {
		s.logger.Error("failed to get bank mortgage rates", zap.Error(err))
		return nil, err
	}
	
	result := make([]*response.MortgageRateResponse, 0, len(rates))
	for _, rate := range rates {
		result = append(result, convertToMortgageRateResponse(rate))
	}
	
	return result, nil
}

// CompareMortgageRates 比较按揭利率
func (s *mortgageService) CompareMortgageRates(ctx context.Context, req *request.CompareMortgageRatesRequest) (*response.MortgageRateComparisonResponse, error) {
	// 获取有效利率
	rates, err := s.repo.GetEffectiveMortgageRates(ctx, req.RateType)
	if err != nil {
		s.logger.Error("failed to get effective mortgage rates", zap.Error(err))
		return nil, err
	}
	
	// 计算月供比较
	comparisons := make([]response.MortgageRateComparisonItem, 0, len(rates))
	for _, rate := range rates {
		monthlyRate := rate.InterestRate / 100 / 12
		
		// 计算月供
		monthlyPayment := calculateMonthlyPayment(req.LoanAmount, monthlyRate, req.LoanPeriod)
		totalPayment := monthlyPayment * float64(req.LoanPeriod)
		totalInterest := totalPayment - req.LoanAmount
		
		// 计算总成本（包括手续费）
		totalCost := totalPayment
		if rate.ProcessingFee != nil {
			totalCost += *rate.ProcessingFee
		}
		
		comparison := response.MortgageRateComparisonItem{
			Bank:           convertToBankResponse(rate.Bank),
			RateType:       string(rate.RateType),
			InterestRate:   rate.InterestRate,
			MonthlyPayment: monthlyPayment,
			TotalPayment:   totalPayment,
			TotalInterest:  totalInterest,
			ProcessingFee:  rate.ProcessingFee,
			TotalCost:      totalCost,
		}
		
		comparisons = append(comparisons, comparison)
	}
	
	resp := &response.MortgageRateComparisonResponse{
		LoanAmount:      req.LoanAmount,
		LoanPeriod:      req.LoanPeriod,
		RateComparisons: comparisons,
	}
	
	return resp, nil
}

// ApplyMortgage 申请按揭
func (s *mortgageService) ApplyMortgage(ctx context.Context, userID uint, req *request.ApplyMortgageRequest) (*response.MortgageApplicationResponse, error) {
	// 检查银行是否存在
	bank, err := s.repo.GetBankByID(ctx, req.BankID)
	if err != nil {
		s.logger.Error("failed to get bank", zap.Error(err))
		return nil, err
	}
	if bank == nil {
		return nil, errors.ErrNotFound
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
	application := &model.MortgageApplication{
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
		Status:              model.MortgageApplicationStatusPending,
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
func (s *mortgageService) GetMortgageApplications(ctx context.Context, userID uint, filter *request.ListMortgageApplicationsRequest) ([]*response.MortgageApplicationResponse, int64, error) {
	applications, total, err := s.repo.GetApplicationsByUserID(ctx, userID, filter)
	if err != nil {
		s.logger.Error("failed to get mortgage applications", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.MortgageApplicationResponse, 0, len(applications))
	for _, app := range applications {
		result = append(result, convertToMortgageApplicationResponse(app))
	}
	
	return result, total, nil
}

// GetMortgageApplication 获取按揭申请详情
func (s *mortgageService) GetMortgageApplication(ctx context.Context, userID uint, id uint) (*response.MortgageApplicationResponse, error) {
	application, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get mortgage application", zap.Error(err))
		return nil, err
	}
	if application == nil {
		return nil, errors.ErrNotFound
	}
	
	// 检查所有权
	if application.UserID != userID {
		return nil, errors.ErrForbidden
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
func generatePaymentSchedule(principal float64, monthlyRate float64, periods int, monthlyPayment float64) []*response.MortgagePaymentSchedule {
	schedule := make([]*response.MortgagePaymentSchedule, 0)
	
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
		
		item := &response.MortgagePaymentSchedule{
			Period:            month,
			Payment:           monthlyPayment,
			Principal:         principalPayment,
			Interest:          interestPayment,
			RemainingBalance:  remainingBalance,
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
		
		item := &response.MortgagePaymentSchedule{
			Period:            periods,
			Payment:           monthlyPayment,
			Principal:         principalPayment,
			Interest:          interestPayment,
			RemainingBalance:  0,
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

// convertToMortgageRateResponse 转换为按揭利率响应
func convertToMortgageRateResponse(rate *model.MortgageRate) *response.MortgageRateResponse {
	resp := &response.MortgageRateResponse{
		ID:                rate.ID,
		BankID:            rate.BankID,
		RateType:          string(rate.RateType),
		InterestRate:      rate.InterestRate * 100, // 转换为百分比格式
		MinLoanAmount:     rate.MinLoanAmount,
		MaxLoanAmount:     rate.MaxLoanAmount,
		MinLoanPeriod:     rate.MinLoanPeriod,
		MaxLoanPeriod:     rate.MaxLoanPeriod,
		LTV:               rate.LTV,
		ProcessingFee:     rate.ProcessingFee,
		ProcessingFeeRate: rate.ProcessingFeeRate,
		Description:       rate.Description,
		EffectiveDate:     rate.EffectiveDate,
		ExpiryDate:        rate.ExpiryDate,
		IsEffective:       rate.IsEffective(),
	}
	
	if rate.Bank != nil {
		resp.Bank = &response.BankResponse{
			ID:         rate.Bank.ID,
			NameZhHant: rate.Bank.NameZhHant,
			NameZhHans: rate.Bank.NameZhHans,
			NameEn:     rate.Bank.NameEn,
			Code:       rate.Bank.Code,
			Logo:       rate.Bank.Logo,
			Website:    rate.Bank.Website,
			Hotline:    rate.Bank.Hotline,
		}
	}
	
	return resp
}

// convertToMortgageApplicationResponse 转换为按揭申请响应
func convertToMortgageApplicationResponse(app *model.MortgageApplication) *response.MortgageApplicationResponse {
	resp := &response.MortgageApplicationResponse{
		ID:                  app.ID,
		ApplicationNo:       app.ApplicationNo,
		UserID:              app.UserID,
		PropertyID:          app.PropertyID,
		BankID:              app.BankID,
		PropertyPrice:       app.PropertyPrice,
		DownPayment:         app.DownPayment,
		LoanAmount:          app.LoanAmount,
		InterestRate:        app.InterestRate * 100, // 转换为百分比
		LoanPeriod:          app.LoanPeriod,
		LoanPeriodYears:     app.LoanPeriod / 12,
		MonthlyPayment:      app.MonthlyPayment,
		TotalPayment:        app.TotalPayment,
		TotalInterest:       app.TotalInterest,
		LTV:                 app.LTV * 100, // 转换为百分比
		ApplicantName:       app.ApplicantName,
		ApplicantPhone:      app.ApplicantPhone,
		ApplicantEmail:      app.ApplicantEmail,
		ApplicantIncome:     app.ApplicantIncome,
		ApplicantOccupation: app.ApplicantOccupation,
		Remarks:             app.Remarks,
		Status:              string(app.Status),
		RejectionReason:     app.RejectionReason,
		ApprovedAt:          app.ApprovedAt,
		RejectedAt:          app.RejectedAt,
		CompletedAt:         app.CompletedAt,
		SubmittedAt:         app.SubmittedAt,
		CreatedAt:           app.CreatedAt,
		UpdatedAt:           app.UpdatedAt,
	}
	
	// 设置状态标志
	resp.CanUpdate = app.CanUpdate()
	resp.CanWithdraw = app.CanWithdraw()
	
	// 设置关联对象
	if app.Bank != nil {
		resp.Bank = &response.BankResponse{
			ID:         app.Bank.ID,
			NameZhHant: app.Bank.NameZhHant,
			NameZhHans: app.Bank.NameZhHans,
			NameEn:     app.Bank.NameEn,
			Code:       app.Bank.Code,
			Logo:       app.Bank.Logo,
			Website:    app.Bank.Website,
			Hotline:    app.Bank.Hotline,
		}
	}
	
	if app.Property != nil {
		resp.Property = &response.PropertyBasicResponse{
			ID:      app.Property.ID,
			Title:   app.Property.Title,
			Address: app.Property.Address,
			Price:   app.Property.Price,
		}
	}
	
	return resp
}


// convertPaymentScheduleSlice 转换还款计划切片
func convertPaymentScheduleSlice(schedules []*response.MortgagePaymentSchedule) []response.MortgagePaymentSchedule {
	result := make([]response.MortgagePaymentSchedule, 0, len(schedules))
	for _, s := range schedules {
		if s != nil {
			result = append(result, *s)
		}
	}
	return result
}

// convertToBankResponse 转换银行模型为响应
func convertToBankResponse(bank *model.Bank) *response.BankResponse {
	if bank == nil {
		return nil
	}
	return &response.BankResponse{
		ID:         bank.ID,
		NameZhHant: bank.NameZhHant,
		NameZhHans: bank.NameZhHans,
		NameEn:     bank.NameEn,
		Logo:       bank.Logo,
		Website:    bank.Website,
	}
}
