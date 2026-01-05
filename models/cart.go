package models

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

// ============ Response DTO ============

// CartResponse 购物车响应
type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	TotalPrice float64    `json:"total_price"`
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	FurnitureID uint `json:"furniture_id" binding:"required"`
	Quantity    int  `json:"quantity" binding:"required,min=1"`
}

// AddToCartResponse 添加到购物车响应
type AddToCartResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Item    *CartItem `json:"item"`
}

// UpdateCartItemRequest 更新购物车项请求
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}
