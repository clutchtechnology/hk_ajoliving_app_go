package request

import "time"

// ListFurnitureRequest 家具列表请求
type ListFurnitureRequest struct {
	CategoryID         *uint    `form:"category_id" json:"category_id"`                                  // 分类ID
	MinPrice           *float64 `form:"min_price" json:"min_price" binding:"omitempty,gte=0"`            // 最低价格
	MaxPrice           *float64 `form:"max_price" json:"max_price" binding:"omitempty,gte=0"`            // 最高价格
	Condition          *string  `form:"condition" json:"condition"`                                      // 新旧程度
	Brand              *string  `form:"brand" json:"brand"`                                              // 品牌
	DeliveryDistrictID *uint    `form:"delivery_district_id" json:"delivery_district_id"`                // 交收地区ID
	DeliveryMethod     *string  `form:"delivery_method" json:"delivery_method"`                          // 交收方法
	Status             *string  `form:"status" json:"status"`                                            // 状态
	Keyword            string   `form:"keyword" json:"keyword"`                                          // 关键字搜索
	SortBy             string   `form:"sort_by" json:"sort_by" binding:"omitempty,oneof=created_at price view_count favorite_count published_at"` // 排序字段
	SortOrder          string   `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc"` // 排序方向
	Page               int      `form:"page" json:"page" binding:"omitempty,min=1"`                      // 页码
	PageSize           int      `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`    // 每页数量
}

// GetFurnitureRequest 获取家具详情请求
type GetFurnitureRequest struct {
	ID uint `uri:"id" binding:"required,min=1"` // 家具ID
}

// CreateFurnitureRequest 创建家具请求
type CreateFurnitureRequest struct {
	Title              string    `json:"title" binding:"required,max=255"`                              // 标题
	Description        *string   `json:"description"`                                                   // 描述
	Price              float64   `json:"price" binding:"required,gt=0"`                                 // 价格
	CategoryID         uint      `json:"category_id" binding:"required"`                                // 分类ID
	Brand              *string   `json:"brand" binding:"omitempty,max=100"`                             // 品牌
	Condition          string    `json:"condition" binding:"required,oneof=new like_new good fair poor"` // 新旧程度
	PurchaseDate       *time.Time `json:"purchase_date"`                                                // 购买日期
	DeliveryDistrictID uint      `json:"delivery_district_id" binding:"required"`                       // 交收地区ID
	DeliveryTime       *string   `json:"delivery_time" binding:"omitempty,max=100"`                     // 交收时间
	DeliveryMethod     string    `json:"delivery_method" binding:"required,oneof=self_pickup delivery negotiable"` // 交收方法
	ImageURLs          []string  `json:"image_urls"`                                                    // 图片URLs
	ExpiresAt          *time.Time `json:"expires_at"`                                                   // 过期时间
}

// UpdateFurnitureRequest 更新家具请求
type UpdateFurnitureRequest struct {
	Title              *string    `json:"title" binding:"omitempty,max=255"`                             // 标题
	Description        *string    `json:"description"`                                                   // 描述
	Price              *float64   `json:"price" binding:"omitempty,gt=0"`                                // 价格
	CategoryID         *uint      `json:"category_id"`                                                   // 分类ID
	Brand              *string    `json:"brand" binding:"omitempty,max=100"`                             // 品牌
	Condition          *string    `json:"condition" binding:"omitempty,oneof=new like_new good fair poor"` // 新旧程度
	PurchaseDate       *time.Time `json:"purchase_date"`                                                // 购买日期
	DeliveryDistrictID *uint      `json:"delivery_district_id"`                                          // 交收地区ID
	DeliveryTime       *string    `json:"delivery_time" binding:"omitempty,max=100"`                     // 交收时间
	DeliveryMethod     *string    `json:"delivery_method" binding:"omitempty,oneof=self_pickup delivery negotiable"` // 交收方法
	ExpiresAt          *time.Time `json:"expires_at"`                                                   // 过期时间
}

// UpdateFurnitureStatusRequest 更新家具状态请求
type UpdateFurnitureStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=available reserved sold expired cancelled"` // 状态
}
