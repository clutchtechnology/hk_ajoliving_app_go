package response

import "time"

// CartItemResponse 购物车项响应
type CartItemResponse struct {
	ID          uint                      `json:"id"`
	FurnitureID uint                      `json:"furniture_id"`
	Furniture   *CartFurnitureResponse    `json:"furniture,omitempty"`
	Quantity    int                       `json:"quantity"`
	TotalPrice  float64                   `json:"total_price"`
	IsAvailable bool                      `json:"is_available"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

// CartFurnitureResponse 购物车中的家具信息响应
type CartFurnitureResponse struct {
	ID          uint    `json:"id"`
	FurnitureNo string  `json:"furniture_no"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	CoverImage  *string `json:"cover_image,omitempty"`
	IsAvailable bool    `json:"is_available"`
}

// CartResponse 购物车响应
type CartResponse struct {
	Items           []CartItemResponse `json:"items"`
	TotalItems      int                `json:"total_items"`
	TotalQuantity   int                `json:"total_quantity"`
	TotalPrice      float64            `json:"total_price"`
	AvailableItems  int                `json:"available_items"`
	UnavailableItems int               `json:"unavailable_items"`
}

// AddToCartResponse 添加到购物车响应
type AddToCartResponse struct {
	ID          uint      `json:"id"`
	FurnitureID uint      `json:"furniture_id"`
	Quantity    int       `json:"quantity"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message"`
}

// UpdateCartItemResponse 更新购物车项响应
type UpdateCartItemResponse struct {
	ID         uint      `json:"id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	UpdatedAt  time.Time `json:"updated_at"`
	Message    string    `json:"message"`
}
