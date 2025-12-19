package request

// CalculateMortgageRequest 按揭计算请求
type CalculateMortgageRequest struct {
	PropertyPrice float64 `json:"property_price" binding:"required,gt=0"`   // 物业价格
	DownPayment   float64 `json:"down_payment" binding:"required,gte=0"`    // 首付金额
	LoanPeriod    int     `json:"loan_period" binding:"required,min=12,max=360"` // 贷款期限（月）
	InterestRate  float64 `json:"interest_rate" binding:"required,gt=0,lte=20"`  // 年利率（百分比，如 2.5）
}

// ApplyMortgageRequest 按揭申请请求
type ApplyMortgageRequest struct {
	PropertyID          *uint   `json:"property_id"`                                  // 物业ID（可选）
	BankID              uint    `json:"bank_id" binding:"required"`                   // 银行ID
	PropertyPrice       float64 `json:"property_price" binding:"required,gt=0"`       // 物业价格
	LoanAmount          float64 `json:"loan_amount" binding:"required,gt=0"`          // 贷款金额
	LoanPeriod          int     `json:"loan_period" binding:"required,min=12,max=360"` // 贷款期限（月）
	InterestRate        float64 `json:"interest_rate" binding:"required,gt=0,lte=20"` // 年利率（百分比）
	ApplicantName       string  `json:"applicant_name" binding:"required,max=100"`    // 申请人姓名
	ApplicantPhone      string  `json:"applicant_phone" binding:"required,max=20"`    // 申请人电话
	ApplicantEmail      string  `json:"applicant_email" binding:"required,email,max=100"` // 申请人邮箱
	ApplicantIncome     float64 `json:"applicant_income" binding:"required,gte=0"`    // 申请人月收入
	ApplicantOccupation *string `json:"applicant_occupation" binding:"omitempty,max=100"` // 申请人职业
	Remarks             *string `json:"remarks"`                                      // 备注
}

// ListMortgageApplicationsRequest 按揭申请列表请求
type ListMortgageApplicationsRequest struct {
	Status    *string `form:"status" json:"status"`                                              // 状态筛选
	BankID    *uint   `form:"bank_id" json:"bank_id"`                                            // 银行筛选
	SortBy    string  `form:"sort_by" json:"sort_by" binding:"omitempty,oneof=created_at submitted_at"` // 排序字段
	SortOrder string  `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc"`   // 排序方向
	Page      int     `form:"page" json:"page" binding:"omitempty,min=1"`                        // 页码
	PageSize  int     `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`      // 每页数量
}

// CompareMortgageRatesRequest 比较银行利率请求
type CompareMortgageRatesRequest struct {
	LoanAmount float64 `form:"loan_amount" json:"loan_amount" binding:"required,gt=0"` // 贷款金额
	LoanPeriod int     `form:"loan_period" json:"loan_period" binding:"required,min=12,max=360"` // 贷款期限（月）
	RateType   *string `form:"rate_type" json:"rate_type" binding:"omitempty,oneof=fixed floating hybrid"` // 利率类型
}
