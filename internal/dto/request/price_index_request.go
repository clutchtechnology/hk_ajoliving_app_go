package request

// GetPriceIndexRequest 获取楼价指数请求
type GetPriceIndexRequest struct {
	IndexType    *string `form:"index_type" binding:"omitempty,oneof=overall district estate property_type"` // 指数类型
	DistrictID   *uint   `form:"district_id"`                                                                 // 地区ID
	EstateID     *uint   `form:"estate_id"`                                                                   // 屋苑ID
	PropertyType *string `form:"property_type"`                                                               // 物业类型
	StartPeriod  *string `form:"start_period"`                                                                // 开始周期 YYYY-MM
	EndPeriod    *string `form:"end_period"`                                                                  // 结束周期 YYYY-MM
	Page         int     `form:"page,default=1" binding:"min=1"`                                              // 页码
	PageSize     int     `form:"page_size,default=20" binding:"min=1,max=100"`                                // 每页数量
}

// GetDistrictPriceIndexRequest 获取地区楼价指数请求
type GetDistrictPriceIndexRequest struct {
	StartPeriod *string `form:"start_period"` // 开始周期 YYYY-MM
	EndPeriod   *string `form:"end_period"`   // 结束周期 YYYY-MM
	Limit       int     `form:"limit,default=12" binding:"min=1,max=120"` // 返回记录数
}

// GetEstatePriceIndexRequest 获取屋苑楼价指数请求
type GetEstatePriceIndexRequest struct {
	StartPeriod *string `form:"start_period"` // 开始周期 YYYY-MM
	EndPeriod   *string `form:"end_period"`   // 结束周期 YYYY-MM
	Limit       int     `form:"limit,default=12" binding:"min=1,max=120"` // 返回记录数
}

// GetPriceTrendsRequest 获取价格走势请求
type GetPriceTrendsRequest struct {
	IndexType    string  `form:"index_type" binding:"required,oneof=overall district estate property_type"` // 指数类型
	DistrictID   *uint   `form:"district_id"`                                                                // 地区ID（当index_type=district时必填）
	EstateID     *uint   `form:"estate_id"`                                                                  // 屋苑ID（当index_type=estate时必填）
	PropertyType *string `form:"property_type"`                                                              // 物业类型
	StartPeriod  string  `form:"start_period" binding:"required"`                                            // 开始周期 YYYY-MM
	EndPeriod    string  `form:"end_period" binding:"required"`                                              // 结束周期 YYYY-MM
}

// ComparePriceIndexRequest 对比楼价指数请求
type ComparePriceIndexRequest struct {
	CompareType  string  `form:"compare_type" binding:"required,oneof=districts estates property_types"` // 对比类型
	DistrictIDs  []uint  `form:"district_ids"`                                                            // 地区ID列表（当compare_type=districts时使用）
	EstateIDs    []uint  `form:"estate_ids"`                                                              // 屋苑ID列表（当compare_type=estates时使用）
	PropertyTypes []string `form:"property_types"`                                                         // 物业类型列表（当compare_type=property_types时使用）
	StartPeriod  string  `form:"start_period" binding:"required"`                                         // 开始周期 YYYY-MM
	EndPeriod    string  `form:"end_period" binding:"required"`                                           // 结束周期 YYYY-MM
}

// ExportPriceDataRequest 导出价格数据请求
type ExportPriceDataRequest struct {
	IndexType    *string `form:"index_type" binding:"omitempty,oneof=overall district estate property_type"` // 指数类型
	DistrictID   *uint   `form:"district_id"`                                                                 // 地区ID
	EstateID     *uint   `form:"estate_id"`                                                                   // 屋苑ID
	PropertyType *string `form:"property_type"`                                                               // 物业类型
	StartPeriod  string  `form:"start_period" binding:"required"`                                             // 开始周期 YYYY-MM
	EndPeriod    string  `form:"end_period" binding:"required"`                                               // 结束周期 YYYY-MM
	Format       string  `form:"format,default=csv" binding:"omitempty,oneof=csv json excel"`                 // 导出格式
}

// GetPriceIndexHistoryRequest 获取历史楼价指数请求
type GetPriceIndexHistoryRequest struct {
	IndexType    string  `form:"index_type" binding:"required,oneof=overall district estate property_type"` // 指数类型
	DistrictID   *uint   `form:"district_id"`                                                                // 地区ID
	EstateID     *uint   `form:"estate_id"`                                                                  // 屋苑ID
	PropertyType *string `form:"property_type"`                                                              // 物业类型
	Years        int     `form:"years,default=5" binding:"min=1,max=20"`                                     // 查询最近几年的数据
}

// CreatePriceIndexRequest 创建楼价指数请求
type CreatePriceIndexRequest struct {
	IndexType        string   `json:"index_type" binding:"required,oneof=overall district estate property_type"` // 指数类型
	DistrictID       *uint    `json:"district_id"`                                                                // 地区ID
	EstateID         *uint    `json:"estate_id"`                                                                  // 屋苑ID
	PropertyType     *string  `json:"property_type"`                                                              // 物业类型
	IndexValue       float64  `json:"index_value" binding:"required,gte=0"`                                       // 指数值
	ChangeValue      float64  `json:"change_value"`                                                               // 变化值
	ChangePercent    float64  `json:"change_percent"`                                                             // 变化百分比
	AvgPrice         *float64 `json:"avg_price" binding:"omitempty,gte=0"`                                        // 平均价格
	AvgPricePerSqft  *float64 `json:"avg_price_per_sqft" binding:"omitempty,gte=0"`                               // 平均每平方尺价格
	TransactionCount int      `json:"transaction_count" binding:"gte=0"`                                          // 成交数量
	Period           string   `json:"period" binding:"required"`                                                  // 周期 YYYY-MM
	DataSource       string   `json:"data_source" binding:"required"`                                             // 数据来源
	Notes            *string  `json:"notes"`                                                                      // 备注
}

// UpdatePriceIndexRequest 更新楼价指数请求
type UpdatePriceIndexRequest struct {
	IndexValue       *float64 `json:"index_value" binding:"omitempty,gte=0"`         // 指数值
	ChangeValue      *float64 `json:"change_value"`                                  // 变化值
	ChangePercent    *float64 `json:"change_percent"`                                // 变化百分比
	AvgPrice         *float64 `json:"avg_price" binding:"omitempty,gte=0"`           // 平均价格
	AvgPricePerSqft  *float64 `json:"avg_price_per_sqft" binding:"omitempty,gte=0"`  // 平均每平方尺价格
	TransactionCount *int     `json:"transaction_count" binding:"omitempty,gte=0"`   // 成交数量
	DataSource       *string  `json:"data_source"`                                   // 数据来源
	Notes            *string  `json:"notes"`                                         // 备注
}
