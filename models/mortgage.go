package models

import (
	"time"

	"gorm.io/gorm"
)

// Bank 银行表模型
type Bank struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NameZhHant  string    `gorm:"type:varchar(100);not null" json:"name_zh_hant"`
	NameZhHans  *string   `gorm:"type:varchar(100)" json:"name_zh_hans,omitempty"`
	NameEn      *string   `gorm:"type:varchar(100)" json:"name_en,omitempty"`
	Code        string    `gorm:"type:varchar(20);uniqueIndex" json:"code"`
	Logo        *string   `gorm:"type:varchar(255)" json:"logo,omitempty"`
	Website     *string   `gorm:"type:varchar(255)" json:"website,omitempty"`
	Hotline     *string   `gorm:"type:varchar(50)" json:"hotline,omitempty"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Bank) TableName() string {
	return "banks"
}

// GetLocalizedName 根据语言获取本地化名称
func (b *Bank) GetLocalizedName(lang string) string {
	switch lang {
	case "zh-Hans", "zh_CN":
		if b.NameZhHans != nil {
			return *b.NameZhHans
		}
		return b.NameZhHant
	case "en":
		if b.NameEn != nil {
			return *b.NameEn
		}
		return b.NameZhHant
	default: // zh-Hant, zh_HK
		return b.NameZhHant
	}
}

// MortgageRateType 按揭利率类型
type MortgageRateType string

const (
	MortgageRateTypeFixed    MortgageRateType = "fixed"    // 固定利率
	MortgageRateTypeFloating MortgageRateType = "floating" // 浮动利率
	MortgageRateTypeHybrid   MortgageRateType = "hybrid"   // 混合利率
)

// MortgageRate 按揭利率表模型
type MortgageRate struct {
	ID               uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	BankID           uint             `gorm:"not null;index" json:"bank_id"`
	RateType         MortgageRateType `gorm:"type:varchar(20);not null;index" json:"rate_type"`
	InterestRate     float64          `gorm:"type:decimal(5,4);not null" json:"interest_rate"`     // 年利率（如 0.0250 表示 2.50%）
	MinLoanAmount    *float64         `gorm:"type:decimal(15,2)" json:"min_loan_amount,omitempty"` // 最低贷款额
	MaxLoanAmount    *float64         `gorm:"type:decimal(15,2)" json:"max_loan_amount,omitempty"` // 最高贷款额
	MinLoanPeriod    *int             `json:"min_loan_period,omitempty"`                           // 最短贷款期限（月）
	MaxLoanPeriod    *int             `json:"max_loan_period,omitempty"`                           // 最长贷款期限（月）
	LTV              *float64         `gorm:"type:decimal(5,4)" json:"ltv,omitempty"`              // 贷款成数 (Loan-to-Value)
	ProcessingFee    *float64         `gorm:"type:decimal(10,2)" json:"processing_fee,omitempty"`  // 手续费
	ProcessingFeeRate *float64        `gorm:"type:decimal(5,4)" json:"processing_fee_rate,omitempty"` // 手续费率
	Description      *string          `gorm:"type:text" json:"description,omitempty"`
	EffectiveDate    time.Time        `gorm:"not null;index" json:"effective_date"`
	ExpiryDate       *time.Time       `gorm:"index" json:"expiry_date,omitempty"`
	IsActive         bool             `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt        time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Bank *Bank `gorm:"foreignKey:BankID" json:"bank,omitempty"`
}

// TableName 指定表名
func (MortgageRate) TableName() string {
	return "mortgage_rates"
}

// IsEffective 判断利率是否有效
func (m *MortgageRate) IsEffective() bool {
	now := time.Now()
	if !m.IsActive {
		return false
	}
	if now.Before(m.EffectiveDate) {
		return false
	}
	if m.ExpiryDate != nil && now.After(*m.ExpiryDate) {
		return false
	}
	return true
}

// GetInterestRatePercent 获取百分比格式的利率
func (m *MortgageRate) GetInterestRatePercent() float64 {
	return m.InterestRate * 100
}

// MortgageApplicationStatus 按揭申请状态
type MortgageApplicationStatus string

const (
	MortgageApplicationStatusPending    MortgageApplicationStatus = "pending"    // 待审核
	MortgageApplicationStatusApproved   MortgageApplicationStatus = "approved"   // 已批准
	MortgageApplicationStatusRejected   MortgageApplicationStatus = "rejected"   // 已拒绝
	MortgageApplicationStatusWithdrawn  MortgageApplicationStatus = "withdrawn"  // 已撤回
	MortgageApplicationStatusCompleted  MortgageApplicationStatus = "completed"  // 已完成
	MortgageApplicationStatusCancelled  MortgageApplicationStatus = "cancelled"  // 已取消
)

// MortgageApplication 按揭申请表模型
type MortgageApplication struct {
	ID                  uint                      `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationNo       string                    `gorm:"type:varchar(50);not null;uniqueIndex" json:"application_no"`
	UserID              uint                      `gorm:"not null;index" json:"user_id"`
	PropertyID          *uint                     `gorm:"index" json:"property_id,omitempty"`
	BankID              uint                      `gorm:"not null;index" json:"bank_id"`
	PropertyPrice       float64                   `gorm:"type:decimal(15,2);not null" json:"property_price"`
	LoanAmount          float64                   `gorm:"type:decimal(15,2);not null" json:"loan_amount"`
	LoanPeriod          int                       `gorm:"not null" json:"loan_period"` // 贷款期限（月）
	InterestRate        float64                   `gorm:"type:decimal(5,4);not null" json:"interest_rate"`
	MonthlyPayment      float64                   `gorm:"type:decimal(10,2);not null" json:"monthly_payment"`
	TotalPayment        float64                   `gorm:"type:decimal(15,2);not null" json:"total_payment"`
	TotalInterest       float64                   `gorm:"type:decimal(15,2);not null" json:"total_interest"`
	DownPayment         float64                   `gorm:"type:decimal(15,2);not null" json:"down_payment"`
	LTV                 float64                   `gorm:"type:decimal(5,4);not null" json:"ltv"`
	ApplicantName       string                    `gorm:"type:varchar(100);not null" json:"applicant_name"`
	ApplicantPhone      string                    `gorm:"type:varchar(20);not null" json:"applicant_phone"`
	ApplicantEmail      string                    `gorm:"type:varchar(100);not null" json:"applicant_email"`
	ApplicantIncome     float64                   `gorm:"type:decimal(12,2);not null" json:"applicant_income"`
	ApplicantOccupation *string                   `gorm:"type:varchar(100)" json:"applicant_occupation,omitempty"`
	Remarks             *string                   `gorm:"type:text" json:"remarks,omitempty"`
	Status              MortgageApplicationStatus `gorm:"type:varchar(20);not null;index;default:'pending'" json:"status"`
	RejectionReason     *string                   `gorm:"type:text" json:"rejection_reason,omitempty"`
	ApprovedAt          *time.Time                `json:"approved_at,omitempty"`
	RejectedAt          *time.Time                `json:"rejected_at,omitempty"`
	CompletedAt         *time.Time                `json:"completed_at,omitempty"`
	SubmittedAt         time.Time                 `gorm:"not null;index" json:"submitted_at"`
	CreatedAt           time.Time                 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time                 `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Property *Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Bank     *Bank     `gorm:"foreignKey:BankID" json:"bank,omitempty"`
}

// TableName 指定表名
func (MortgageApplication) TableName() string {
	return "mortgage_applications"
}

// ============ Request DTO ============

// ListMortgageApplicationsRequest 获取按揭申请列表请求
type ListMortgageApplicationsRequest struct {
	BankID    *uint   `form:"bank_id"`
	Status    *string `form:"status"`
	Page      int     `form:"page,default=1" binding:"min=1"`
	PageSize  int     `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy    string  `form:"sort_by,default=created_at"`
	SortOrder string  `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// IsPending 判断是否待审核
func (m *MortgageApplication) IsPending() bool {
	return m.Status == MortgageApplicationStatusPending
}

// IsApproved 判断是否已批准
func (m *MortgageApplication) IsApproved() bool {
	return m.Status == MortgageApplicationStatusApproved
}

// IsRejected 判断是否已拒绝
func (m *MortgageApplication) IsRejected() bool {
	return m.Status == MortgageApplicationStatusRejected
}

// IsCompleted 判断是否已完成
func (m *MortgageApplication) IsCompleted() bool {
	return m.Status == MortgageApplicationStatusCompleted
}

// CanUpdate 判断是否可以更新
func (m *MortgageApplication) CanUpdate() bool {
	return m.Status == MortgageApplicationStatusPending
}

// CanWithdraw 判断是否可以撤回
func (m *MortgageApplication) CanWithdraw() bool {
	return m.Status == MortgageApplicationStatusPending || m.Status == MortgageApplicationStatusApproved
}

// GetLTVPercent 获取百分比格式的贷款成数
func (m *MortgageApplication) GetLTVPercent() float64 {
	return m.LTV * 100
}

// GetInterestRatePercent 获取百分比格式的利率
func (m *MortgageApplication) GetInterestRatePercent() float64 {
	return m.InterestRate * 100
}

// BeforeCreate GORM hook - 创建前执行
func (m *MortgageApplication) BeforeCreate(tx *gorm.DB) error {
	if m.Status == "" {
		m.Status = MortgageApplicationStatusPending
	}
	if m.SubmittedAt.IsZero() {
		m.SubmittedAt = time.Now()
	}
	return nil
}
