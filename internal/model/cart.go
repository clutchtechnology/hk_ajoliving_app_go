package model

import (
	"time"

	"gorm.io/gorm"
)

// CartItem 购物车项表模型
type CartItem struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index:idx_user_furniture,unique" json:"user_id"`
	FurnitureID uint           `gorm:"not null;index:idx_user_furniture,unique" json:"furniture_id"`
	Quantity    int            `gorm:"not null;default:1" json:"quantity"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Furniture *Furniture `gorm:"foreignKey:FurnitureID" json:"furniture,omitempty"`
}

// TableName 指定表名
func (CartItem) TableName() string {
	return "cart_items"
}

// GetTotalPrice 计算总价
func (c *CartItem) GetTotalPrice() float64 {
	if c.Furniture == nil {
		return 0
	}
	return c.Furniture.Price * float64(c.Quantity)
}

// IsAvailable 判断商品是否可用
func (c *CartItem) IsAvailable() bool {
	if c.Furniture == nil {
		return false
	}
	return c.Furniture.IsAvailable()
}

// BeforeCreate GORM hook - 创建前执行
func (c *CartItem) BeforeCreate(tx *gorm.DB) error {
	if c.Quantity <= 0 {
		c.Quantity = 1
	}
	return nil
}
