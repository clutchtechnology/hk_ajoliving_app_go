package request

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	FurnitureID uint `json:"furniture_id" binding:"required,min=1"` // 家具ID
	Quantity    int  `json:"quantity" binding:"required,min=1"`      // 数量
}

// UpdateCartItemRequest 更新购物车项请求
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"` // 数量
}
