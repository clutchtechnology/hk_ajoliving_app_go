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
	
	resp := &response.MortgageCalculationResponse{
		PropertyPrice:      req.PropertyPrice,
		DownPayment:        req.DownPayment,
		LoanAmount:         loanAmount,
		InterestRate:       req.InterestRate,
		LoanPeriod:         req.LoanPeriod,
		MonthlyPayment:     monthlyPayment,
		TotalPayment:       totalPayment,
		TotalInterest:      totalInterest,
		FirstPaymentAmount: monthlyPayment,
		PaymentSchedule:    schedule,
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
	rates, err := s.repo.GetEffectiveMortgageRates(ctx, &req.RateType)
	if err != nil {
		s.logger.Error("failed to get effective mortgage rates", zap.Error(err))
		return nil, err
	}
	
	// 转换为响应对象
	rateResponses := make([]*response.MortgageRateResponse, 0, len(rates))
	for _, rate := range rates {
		rateResponses = append(rateResponses, convertToMortgageRateResponse(rate))
	}
	
	// 找出最低利率
	var lowestRate *response.MortgageRateResponse
	for _, rate := range rateResponses {
		if lowestRate == nil || rate.InterestRate < lowestRate.InterestRate {
			lowestRate = rate
		}
	}
	
	// 计算月供比较
	comparisons := make([]*response.MortgageRateItemComparison, 0, len(rateResponses))
	for _, rate := range rateResponses {
		// 计算贷款金额
		loanAmount := req.PropertyPrice - req.DownPayment
		monthlyRate := rate.InterestRate / 100 / 12
		
		// 计算月供
		monthlyPayment := calculateMonthlyPayment(loanAmount, monthlyRate, req.LoanPeriod)
		totalPayment := monthlyPayment * float64(req.LoanPeriod)
		totalInterest := totalPayment - loanAmount
		
		// 计算与最低利率的差额
		var savingsVsLowest float64
		if lowestRate != nil && rate.ID != lowestRate.ID {
			lowestMonthlyRate := lowestRate.InterestRate / 100 / 12
			lowestMonthlyPayment := calculateMonthlyPayment(loanAmount, lowestMonthlyRate, req.LoanPeriod)
			lowestTotalPayment := lowestMonthlyPayment * float64(req.LoanPeriod)
			savingsVsLowest = totalPayment - lowestTotalPayment
		}
		
		comparison := &response.MortgageRateItemComparison{
			RateID:              rate.ID,
			BankName:            rate.Bank.NameZhHant,
			InterestRate:        rate.InterestRate,
			RateType:            rate.RateType,
			MonthlyPayment:      monthlyPayment,
			TotalPayment:        totalPayment,
			TotalInterest:       totalInterest,
			SavingsVsLowest:     savingsVsLowest,
			IsLowestRate:        lowestRate != nil && rate.ID == lowestRate.ID,
		}
		
		comparisons = append(comparisons, comparison)
	}
	
	resp := &response.MortgageRateComparisonResponse{
		PropertyPrice: req.PropertyPrice,
		DownPayment:   req.DownPayment,
		LoanAmount:    req.PropertyPrice - req.DownPayment,
		LoanPeriod:    req.LoanPeriod,
		RateType:      req.RateType,
		Comparisons:   comparisons,
		LowestRate:    lowestRate,
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
	loanAmount := req.PropertyPrice - req.DownPayment
	monthlyRate := req.InterestRate / 100 / 12
	monthlyPayment := calculateMonthlyPayment(loanAmount, monthlyRate, req.LoanPeriod)
	
	// 创建申请
	application := &model.MortgageApplication{
		ApplicationNo:      applicationNo,
		UserID:             userID,
		PropertyID:         req.PropertyID,
		BankID:             req.BankID,
		PropertyPrice:      req.PropertyPrice,
		DownPayment:        req.DownPayment,
		LoanAmount:         loanAmount,
		InterestRate:       req.InterestRate,
		LoanPeriod:         req.LoanPeriod,
		MonthlyIncome:      req.MonthlyIncome,
		EmploymentStatus:   req.EmploymentStatus,
		EmployerName:       req.EmployerName,
		YearsEmployed:      req.YearsEmployed,
		Status:             "pending",
		ApplicantName:      req.ApplicantName,
		ApplicantPhone:     req.ApplicantPhone,
		ApplicantEmail:     req.ApplicantEmail,
		ApplicantIDCard:    req.ApplicantIDCard,
		Notes:              req.Notes,
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
func generatePaymentSchedule(principal float64, monthlyRate float64, periods int, monthlyPayment float64) []*response.PaymentScheduleItem {
	schedule := make([]*response.PaymentScheduleItem, 0)
	
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
		
		item := &response.PaymentScheduleItem{
			Month:              month,
			Payment:            monthlyPayment,
			Principal:          principalPayment,
			Interest:           interestPayment,
			RemainingBalance:   remainingBalance,
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
		
		item := &response.PaymentScheduleItem{
			Month:              periods,
			Payment:            monthlyPayment,
			Principal:          principalPayment,
			Interest:           interestPayment,
			RemainingBalance:   0,
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
		ID:            rate.ID,
		BankID:        rate.BankID,
		RateType:      rate.RateType,
		InterestRate:  rate.InterestRate,
		MinLoanAmount: rate.MinLoanAmount,
		MaxLoanAmount: rate.MaxLoanAmount,
		MinLoanPeriod: rate.MinLoanPeriod,
		MaxLoanPeriod: rate.MaxLoanPeriod,
		MaxLTV:        rate.MaxLTV,
		IsPromotional: rate.IsPromotional,
		EffectiveDate: rate.EffectiveDate,
		ExpiryDate:    rate.ExpiryDate,
		Description:   rate.Description,
		IsActive:      rate.IsActive,
		CreatedAt:     rate.CreatedAt,
		UpdatedAt:     rate.UpdatedAt,
	}
	
	if rate.Bank != nil {
		resp.Bank = &response.BankResponse{
			ID:         rate.Bank.ID,
			NameZhHant: rate.Bank.NameZhHant,
			NameZhHans: rate.Bank.NameZhHans,
			NameEn:     rate.Bank.NameEn,
			Code:       rate.Bank.Code,
			LogoURL:    rate.Bank.LogoURL,
			WebsiteURL: rate.Bank.WebsiteURL,
			IsActive:   rate.Bank.IsActive,
		}
	}
	
	return resp
}

// convertToMortgageApplicationResponse 转换为按揭申请响应
func convertToMortgageApplicationResponse(app *model.MortgageApplication) *response.MortgageApplicationResponse {
	resp := &response.MortgageApplicationResponse{
		ID:               app.ID,
		ApplicationNo:    app.ApplicationNo,
		UserID:           app.UserID,
		PropertyID:       app.PropertyID,
		BankID:           app.BankID,
		PropertyPrice:    app.PropertyPrice,
		DownPayment:      app.DownPayment,
		LoanAmount:       app.LoanAmount,
		InterestRate:     app.InterestRate,
		LoanPeriod:       app.LoanPeriod,
		MonthlyIncome:    app.MonthlyIncome,
		EmploymentStatus: app.EmploymentStatus,
		EmployerName:     app.EmployerName,
		YearsEmployed:    app.YearsEmployed,
		Status:           app.Status,
		ApplicantName:    app.ApplicantName,
		ApplicantPhone:   app.ApplicantPhone,
		ApplicantEmail:   app.ApplicantEmail,
		ApplicantIDCard:  app.ApplicantIDCard,
		Notes:            app.Notes,
		RejectionReason:  app.RejectionReason,
		SubmittedAt:      app.SubmittedAt,
		ReviewedAt:       app.ReviewedAt,
		ApprovedAt:       app.ApprovedAt,
		RejectedAt:       app.RejectedAt,
		WithdrawnAt:      app.WithdrawnAt,
		CreatedAt:        app.CreatedAt,
		UpdatedAt:        app.UpdatedAt,
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
			LogoURL:    app.Bank.LogoURL,
			WebsiteURL: app.Bank.WebsiteURL,
			IsActive:   app.Bank.IsActive,
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
