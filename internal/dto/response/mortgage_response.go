package response

import "time"

// MortgageCalculationResponse 按揭计算响应
type MortgageCalculationResponse struct {
	PropertyPrice       float64                       `json:"property_price"`        // 物业价格
	DownPayment         float64                       `json:"down_payment"`          // 首付金额
	LoanAmount          float64                       `json:"loan_amount"`           // 贷款金额
	LoanPeriod          int                           `json:"loan_period"`           // 贷款期限（月）
	LoanPeriodYears     int                           `json:"loan_period_years"`     // 贷款期限（年）
	InterestRate        float64                       `json:"interest_rate"`         // 年利率（百分比）
	MonthlyPayment      float64                       `json:"monthly_payment"`       // 月供
	TotalPayment        float64                       `json:"total_payment"`         // 总还款额
	TotalInterest       float64                       `json:"total_interest"`        // 总利息
	LTV                 float64                       `json:"ltv"`                   // 贷款成数（百分比）
	PaymentSchedule     []MortgagePaymentSchedule     `json:"payment_schedule"`      // 还款计划（可选）
}

// MortgagePaymentSchedule 还款计划
type MortgagePaymentSchedule struct {
	Period            int     `json:"period"`             // 期数
	Payment           float64 `json:"payment"`            // 还款额
	Principal         float64 `json:"principal"`          // 本金
	Interest          float64 `json:"interest"`           // 利息
	RemainingBalance  float64 `json:"remaining_balance"`  // 剩余本金
}

// BankResponse 银行响应
type BankResponse struct {
	ID         uint    `json:"id"`
	NameZhHant string  `json:"name_zh_hant"`
	NameZhHans *string `json:"name_zh_hans,omitempty"`
	NameEn     *string `json:"name_en,omitempty"`
	Code       string  `json:"code"`
	Logo       *string `json:"logo,omitempty"`
	Website    *string `json:"website,omitempty"`
	Hotline    *string `json:"hotline,omitempty"`
}

// MortgageRateResponse 按揭利率响应
type MortgageRateResponse struct {
	ID                uint         `json:"id"`
	BankID            uint         `json:"bank_id"`
	Bank              *BankResponse `json:"bank,omitempty"`
	RateType          string       `json:"rate_type"`
	InterestRate      float64      `json:"interest_rate"`       // 百分比格式（如 2.50）
	MinLoanAmount     *float64     `json:"min_loan_amount,omitempty"`
	MaxLoanAmount     *float64     `json:"max_loan_amount,omitempty"`
	MinLoanPeriod     *int         `json:"min_loan_period,omitempty"`
	MaxLoanPeriod     *int         `json:"max_loan_period,omitempty"`
	LTV               *float64     `json:"ltv,omitempty"`
	ProcessingFee     *float64     `json:"processing_fee,omitempty"`
	ProcessingFeeRate *float64     `json:"processing_fee_rate,omitempty"`
	Description       *string      `json:"description,omitempty"`
	EffectiveDate     time.Time    `json:"effective_date"`
	ExpiryDate        *time.Time   `json:"expiry_date,omitempty"`
	IsEffective       bool         `json:"is_effective"`
}

// MortgageRateComparisonResponse 利率比较响应
type MortgageRateComparisonResponse struct {
	LoanAmount      float64                            `json:"loan_amount"`
	LoanPeriod      int                                `json:"loan_period"`
	RateComparisons []MortgageRateComparisonItem       `json:"rate_comparisons"`
}

// MortgageRateComparisonItem 利率比较项
type MortgageRateComparisonItem struct {
	Bank            *BankResponse `json:"bank"`
	RateType        string        `json:"rate_type"`
	InterestRate    float64       `json:"interest_rate"`
	MonthlyPayment  float64       `json:"monthly_payment"`
	TotalPayment    float64       `json:"total_payment"`
	TotalInterest   float64       `json:"total_interest"`
	ProcessingFee   *float64      `json:"processing_fee,omitempty"`
	TotalCost       float64       `json:"total_cost"`
}

// MortgageApplicationResponse 按揭申请响应
type MortgageApplicationResponse struct {
	ID                  uint                         `json:"id"`
	ApplicationNo       string                       `json:"application_no"`
	UserID              uint                         `json:"user_id"`
	PropertyID          *uint                        `json:"property_id,omitempty"`
	Property            *PropertyBasicResponse       `json:"property,omitempty"`
	BankID              uint                         `json:"bank_id"`
	Bank                *BankResponse                `json:"bank,omitempty"`
	PropertyPrice       float64                      `json:"property_price"`
	LoanAmount          float64                      `json:"loan_amount"`
	LoanPeriod          int                          `json:"loan_period"`
	LoanPeriodYears     int                          `json:"loan_period_years"`
	InterestRate        float64                      `json:"interest_rate"`
	MonthlyPayment      float64                      `json:"monthly_payment"`
	TotalPayment        float64                      `json:"total_payment"`
	TotalInterest       float64                      `json:"total_interest"`
	DownPayment         float64                      `json:"down_payment"`
	LTV                 float64                      `json:"ltv"`
	ApplicantName       string                       `json:"applicant_name"`
	ApplicantPhone      string                       `json:"applicant_phone"`
	ApplicantEmail      string                       `json:"applicant_email"`
	ApplicantIncome     float64                      `json:"applicant_income"`
	ApplicantOccupation *string                      `json:"applicant_occupation,omitempty"`
	Remarks             *string                      `json:"remarks,omitempty"`
	Status              string                       `json:"status"`
	RejectionReason     *string                      `json:"rejection_reason,omitempty"`
	ApprovedAt          *time.Time                   `json:"approved_at,omitempty"`
	RejectedAt          *time.Time                   `json:"rejected_at,omitempty"`
	CompletedAt         *time.Time                   `json:"completed_at,omitempty"`
	SubmittedAt         time.Time                    `json:"submitted_at"`
	CreatedAt           time.Time                    `json:"created_at"`
	UpdatedAt           time.Time                    `json:"updated_at"`
	CanUpdate           bool                         `json:"can_update"`
	CanWithdraw         bool                         `json:"can_withdraw"`
}

// PropertyBasicResponse 物业基本信息响应
type PropertyBasicResponse struct {
	ID          uint    `json:"id"`
	PropertyNo  string  `json:"property_no"`
	Title       string  `json:"title"`
	Address     string  `json:"address"`
	Price       float64 `json:"price"`
}

// MortgageApplicationListItemResponse 按揭申请列表项响应
type MortgageApplicationListItemResponse struct {
	ID              uint           `json:"id"`
	ApplicationNo   string         `json:"application_no"`
	BankID          uint           `json:"bank_id"`
	Bank            *BankResponse  `json:"bank,omitempty"`
	PropertyPrice   float64        `json:"property_price"`
	LoanAmount      float64        `json:"loan_amount"`
	MonthlyPayment  float64        `json:"monthly_payment"`
	Status          string         `json:"status"`
	SubmittedAt     time.Time      `json:"submitted_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// ApplyMortgageResponse 按揭申请响应
type ApplyMortgageResponse struct {
	ID            uint      `json:"id"`
	ApplicationNo string    `json:"application_no"`
	Status        string    `json:"status"`
	SubmittedAt   time.Time `json:"submitted_at"`
	Message       string    `json:"message"`
}
