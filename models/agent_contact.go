package models

import (
	"time"

	"gorm.io/gorm"
)

// AgentContactRequest 代理人联系请求
type AgentContactRequest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AgentID     uint           `gorm:"index;not null" json:"agent_id"`                // 代理人ID
	UserID      *uint          `gorm:"index" json:"user_id"`                          // 用户ID（可选，未登录用户为空）
	PropertyID  *uint          `gorm:"index" json:"property_id"`                      // 房源ID（可选）
	Name        string         `gorm:"size:100;not null" json:"name"`                 // 联系人姓名
	Phone       string         `gorm:"size:20;not null" json:"phone"`                 // 联系电话
	Email       string         `gorm:"size:255" json:"email"`                         // 联系邮箱
	Message     string         `gorm:"type:text" json:"message"`                      // 留言内容
	ContactType string         `gorm:"size:20;default:'inquiry'" json:"contact_type"` // 联系类型: inquiry, viewing, valuation, other
	Status      string         `gorm:"size:20;default:'pending';index" json:"status"` // 状态: pending, contacted, completed, cancelled
	ContactedAt *time.Time     `json:"contacted_at"`                                  // 联系时间
	Notes       string         `gorm:"type:text" json:"notes"`                        // 备注（代理人添加）
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Agent    *Agent    `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Property *Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (AgentContactRequest) TableName() string {
	return "agent_contact_requests"
}

// 联系类型常量
const (
	ContactTypeInquiry   = "inquiry"   // 咨询
	ContactTypeViewing   = "viewing"   // 看房
	ContactTypeValuation = "valuation" // 估价
	ContactTypeOther     = "other"     // 其他
)

// 状态常量
const (
	ContactStatusPending   = "pending"   // 待处理
	ContactStatusContacted = "contacted" // 已联系
	ContactStatusCompleted = "completed" // 已完成
	ContactStatusCancelled = "cancelled" // 已取消
)

// BeforeCreate GORM 钩子：创建前设置默认值
func (c *AgentContactRequest) BeforeCreate(tx *gorm.DB) error {
	if c.ContactType == "" {
		c.ContactType = ContactTypeInquiry
	}
	if c.Status == "" {
		c.Status = ContactStatusPending
	}
	return nil
}
