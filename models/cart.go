package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// CartItem 购物车项
type CartItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index:idx_user_furniture" json:"user_id"`
	FurnitureID uint           `gorm:"not null;index:idx_user_furniture" json:"furniture_id"`
	Quantity    int            `gorm:"not null;default:1" json:"quantity"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Furniture *Furniture `gorm:"foreignKey:FurnitureID" json:"furniture,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

// ============ Request DTO ============

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	FurnitureID uint `json:"furniture_id" binding:"required"`
	Quantity    int  `json:"quantity" binding:"required,min=1,max=99"`
}

// UpdateCartItemRequest 更新购物车项请求
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=99"`
}

// ============ Response DTO ============

// CartItemResponse 购物车项响应
type CartItemResponse struct {
	ID          uint              `json:"id"`
	FurnitureID uint              `json:"furniture_id"`
	Quantity    int               `json:"quantity"`
	Furniture   *FurnitureInCart  `json:"furniture"`
	CreatedAt   time.Time         `json:"created_at"`
}

// FurnitureInCart 购物车中的家具信息
type FurnitureInCart struct {
	ID            uint    `json:"id"`
	FurnitureNo   string  `json:"furniture_no"`
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	CoverImageURL string  `json:"cover_image_url"`
	CategoryName  string  `json:"category_name"`
	DistrictName  string  `json:"district_name"`
}

// CartResponse 购物车响应
type CartResponse struct {
	Items      []*CartItemResponse `json:"items"`
	TotalItems int                 `json:"total_items"`
	TotalPrice float64             `json:"total_price"`
}
